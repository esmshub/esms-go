package models

type Condition interface {
	Validate() error
}

type Conditional struct {
	Action     string
	Values     []any
	Conditions []Condition
}
