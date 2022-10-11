package demo

type DBOption func(*DB)

type DB struct {
	r *registry
}

func NewDB(opts ...DBOption) (*DB, error) {
	res := &DB{
		r: &registry{},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}
