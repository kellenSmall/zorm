package demo

import (
	"database/sql"
	"zorm/demo/internal/model"
	"zorm/demo/internal/valuer"
)

type DBOption func(*DB)

type DB struct {
	r          model.Registry
	db         *sql.DB
	valCreator valuer.Creator
}

func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func DBUseReflectValuer() DBOption {
	return func(db *DB) {
		db.valCreator = valuer.NewValueUnsafe
	}
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		db:         db,
		r:          model.NewRegistry(),
		valCreator: valuer.NewValueReflect,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}
