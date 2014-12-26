package core

func NewRange(val uint) Range {
	return Range{val, val}
}

type Range struct {
	Min, Max uint
}

func (r *Range) Update(val uint) {
	r.Max += val
	if r.Min < val {
		r.Min = val
	}
}
