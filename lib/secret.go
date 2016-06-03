package lib

type Secret struct {
	Object    Object
	Rule      Rule
	Nline     int
	Line      string
	Exception bool
}

func (s *Seekret) ListSecrets() []Secret {
	return s.secretList
}