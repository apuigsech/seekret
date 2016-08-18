package models

type Secret struct {
	Object    *Object
	Rule      *Rule
	Nline     int
	Line      string
	Exception bool
}

func NewSecret(object *Object, rule *Rule, nLine int, line string) *Secret {
	s := &Secret{
		Object: object,
		Rule: rule,
		Nline: nLine,
		Line: line,
	}
	return s
}

func (s *Secret) SetException(exception bool) {
	s.Exception = exception
}