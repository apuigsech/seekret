package seekret

import (
	"github.com/apuigsech/seekret/models"
)

type Secret struct {
	Object    models.Object
	Rule      models.Rule
	Nline     int
	Line      string
	Exception bool
}

func (s *Seekret) ListSecrets() []Secret {
	return s.secretList
}
