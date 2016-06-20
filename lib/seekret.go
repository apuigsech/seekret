package lib

type Seekret struct {
	ruleList      []Rule
	objectList    []Object
	secretList    []Secret
	exceptionList []Exception
}

func NewSeekret() *Seekret {
	s := &Seekret{}
	return s
}