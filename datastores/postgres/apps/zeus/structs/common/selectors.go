package common

type Selector struct {
	MatchLabels ChildClassMultiValue
}

func NewSelector() Selector {
	s := Selector{MatchLabels: NewChildClassAndValues("selector")}
	return s
}
