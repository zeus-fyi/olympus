package structs

type Selector struct {
	MatchLabels ChildClassMultiValue
}

func NewSelector() Selector {
	s := Selector{MatchLabels: NewChildClassMultiValues("selector")}
	return s
}
