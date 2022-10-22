package common

type Table struct {
	Name string
	Columns
}

type Columns []Column

type Column struct {
	Name     string
	DataType string
}
