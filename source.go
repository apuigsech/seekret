// Copyright 2016 - Authors included on AUTHORS file.
//
// Use of this source code is governed by a Apache License
// that can be found in the LICENSE file.

package seekret

import (
	"github.com/apuigsech/seekret/models"
)

type LoadOptions map[string]interface{}

type SourceType interface {
	LoadObjects(source string, opt LoadOptions) ([]models.Object, error)
}

func (s *Seekret) LoadObjects(st SourceType, source string, opt LoadOptions) error {
	objectList, err := st.LoadObjects(source, opt)
	if err != nil {
		return err
	}
	s.objectList = append(s.objectList, objectList...)
	return nil
}
