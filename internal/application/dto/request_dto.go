package dto

type FindByIDOption struct {
	Preloads []string
}

type OffsetPageOption struct {
	Page    int8
	Size    int8
	Where   []WhereCondition
	Sort    []SortCondition
	Preload []string
}

type WhereCondition struct {
	Field    string
	Operator Operator
	Value    any
}

type SortCondition struct {
	Field     string
	Direction Direction
}

type Operator int

const (
	OpEqual Operator = iota
	OpGreaterThan
	OpGreaterThanOrEqual
	OpLessThan
	OpLessThanOrEqual
	OpLike
	OpILike
	OpNotEqual
	OpIn
)

func (o Operator) String() string {
	return []string{"=", ">", ">=", "<", "<=", "LIKE", "ILIKE", "!=", "IN"}[o]
}

func (o Operator) IsValid() bool {
	switch o {
	case OpEqual, OpGreaterThan, OpGreaterThanOrEqual, OpLessThan,
		OpLessThanOrEqual, OpLike, OpILike, OpNotEqual, OpIn:
		return true
	default:
		return false
	}
}

type Direction string

const (
	Asc  Direction = "ASC"
	Desc Direction = "DESC"
)

func (d Direction) IsValid() bool {
	return d == Asc || d == Desc
}
