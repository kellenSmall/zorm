package demo

type Aggregate struct {
	fn  string
	arg string
}

func (a Aggregate) selectable() {}

func Avg(col string) Aggregate {
	return Aggregate{
		fn:  "AVG",
		arg: col,
	}
}

func Max(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "MAX",
	}
}
func Count(col string) Aggregate {
	return Aggregate{
		arg: col,
		fn:  "COUNT",
	}
}
