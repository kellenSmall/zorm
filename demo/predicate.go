package demo

// op 代表操作符
type op string

func (o op) String() string {
	return string(o)
}

// 后面可以每次支持新的操作符就加一个
const (
	opEQ op = "="
	opLT op = "<"
	opGT op = ">"

	opNOT = "NOT"
	opAND = "AND"
	opOR  = "OR"
)

// exprOf 把any 转换成 Expression
func exprOf(e any) Expression {
	switch expr := e.(type) {
	case Expression:
		return expr
	default:
		return valueOf(expr)
	}
}

// Expression 代表语句，或者语句的部分
// 暂时没想好怎么设计方法，所以直接做成标记接口
type Expression interface {
	expr()
}

// Predicate 代表一个查询条件
// Predicate 可以通过和 Predicate 组合构成复杂的查询条件
type Predicate struct {
	left  Expression
	Op    op
	right Expression
}

func (Predicate) expr() {}

func Not(p Predicate) Predicate {
	return Predicate{
		Op:    opNOT,
		right: p,
	}
}

func (p Predicate) And(r Predicate) Predicate {
	return Predicate{
		left:  p,
		Op:    opAND,
		right: r,
	}
}

func (p Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  p,
		Op:    opOR,
		right: right,
	}
}
