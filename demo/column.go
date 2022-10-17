package demo

type Column struct {
	name string
}

func (c Column) selectable() {}

func (Column) expr() {}

type value struct {
	val any
}

// valueOf 把any转换成 value
func valueOf(e any) value {
	return value{val: e}
}

func (value) expr() {}

func C(name string) Column {
	return Column{name: name}
}

// Eq 例如：C("id").Eq(12)
func (c Column) EQ(arg any) Predicate {
	return Predicate{
		left:  c,
		Op:    opEQ,
		right: exprOf(arg),
	}
}

func (c Column) LT(arg any) Predicate {
	return Predicate{
		left:  c,
		Op:    opLT,
		right: exprOf(arg),
	}
}

func (c Column) GT(arg any) Predicate {
	return Predicate{
		left:  c,
		Op:    opGT,
		right: exprOf(arg),
	}
}
