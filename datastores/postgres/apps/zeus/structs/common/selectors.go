package common

type Selector struct {
	MatchLabels ChildClassAndValues
}

func NewSelector() Selector {
	s := Selector{MatchLabels: NewChildClassAndValues("selector")}
	return s
}
