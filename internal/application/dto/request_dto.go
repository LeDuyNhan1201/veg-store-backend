package dto

type FindByIdOption struct {
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
	Operator string // =, >, >=, <, <=, LIKE, ILIKE, !=, IN
	Value    any
}

type SortCondition struct {
	Field     string
	Direction string // ASC or DESC
}
