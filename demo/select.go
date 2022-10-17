package demo

import (
	"context"
	"fmt"
	"strings"
	"zorm/demo/internal/errs"
	model2 "zorm/demo/internal/model"
)

type Selectable interface {
	selectable()
}

type Selector[T any] struct {
	sb      strings.Builder
	table   string
	where   []Predicate
	args    []any
	columns []Selectable
	model   *model2.Model
	db      *DB
}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	rows, err := s.db.db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}
	t := new(T)
	val := s.db.valCreator(t, s.model)

	return t, val.SetColumns(rows)
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) Build() (*Query, error) {
	t := new(T)
	var err error
	s.model, err = s.db.r.Get(t)
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("SELECT ")
	if len(s.columns) == 0 {
		s.sb.WriteByte('*')
	} else {
		for i, c := range s.columns {
			if i > 0 {
				s.sb.WriteByte(',')
			}
			switch col := c.(type) {
			case Column:
				fd, ok := s.model.FieldMap[col.name]
				if !ok {
					return nil, errs.NewErrUnknownColumn(col.name)
				}
				s.sb.WriteByte('`')
				s.sb.WriteString(fd.ColName)
				s.sb.WriteByte('`')
			case Aggregate:
				s.sb.WriteString(col.fn)
				s.sb.WriteByte('(')
				fd, ok := s.model.FieldMap[col.arg]
				if !ok {
					return nil, errs.NewErrUnknownColumn(col.arg)
				}
				s.sb.WriteByte('`')
				s.sb.WriteString(fd.ColName)
				s.sb.WriteByte('`')
				s.sb.WriteByte(')')
			case RawExpr:
				s.sb.WriteString(col.raw)
				if len(col.args) > 0 {
					s.args = append(s.args, col.args...)
				}
			}
		}
	}
	s.sb.WriteString(" FROM ")

	if s.table == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.model.TableName)
		s.sb.WriteByte('`')
	} else {
		s.sb.WriteString(s.table)
	}
	// 构造 WHERE
	if len(s.where) > 0 {
		// 类似这种可有可无的部分，都要在前面加一个空格
		s.sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpression(p); err != nil {
			return nil, err
		}
	}
	s.sb.WriteString(";")
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) buildExpression(p Expression) error {
	if p == nil {
		return nil
	}
	switch exp := p.(type) {
	case Column:
		s.sb.WriteByte('`')
		s.sb.WriteString(exp.name)
		s.sb.WriteByte('`')
	case value:
		s.sb.WriteByte('?')
		s.args = append(s.args, exp.val)
	case Predicate:
		_, lp := exp.left.(Predicate)
		if lp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if lp {
			s.sb.WriteByte(')')
		}

		s.sb.WriteByte(' ')
		s.sb.WriteString(exp.Op.String())
		s.sb.WriteByte(' ')

		_, rp := exp.right.(Predicate)
		if rp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if rp {
			s.sb.WriteByte(')')
		}

	default:
		return fmt.Errorf("orm: 不支持的表达式 %v", exp)
	}
	return nil
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db: db,
	}
}
