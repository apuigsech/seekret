// Copyright 2016 - Authors included on AUTHORS file.
//
// Use of this source code is governed by a Apache License
// that can be found in the LICENSE file.

package seekret

import (
	"fmt"
	"github.com/apuigsech/seekret/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Seekret contains a seekret context and exposes the API to manipulate it.
type Seekret struct {
	// List of rules loaded into the context.
	ruleList []models.Rule

	// List of objects loaded into the context.
	objectList []models.Object

	// List of exceptions loaded into the context.
	exceptionList []models.Exception

	// List of secrets detected after the inspection.
	secretList []models.Secret
}

// NewSeekret returns a new seekret context.
func NewSeekret() *Seekret {
	s := &Seekret{}
	return s
}

// AddRule adds a new rule into the context.
func (s *Seekret) AddRule(rule models.Rule, enabled bool) {
	if enabled {
		rule.Enable()
	}
	s.ruleList = append(s.ruleList, rule)
}

type ruleYaml struct {
	ObjectMatch string
	Match       string
	Unmatch     []string
}

// LoadRulesFromFile loads rules from a YAML file.
func (s *Seekret) LoadRulesFromFile(file string, defaulEnabled bool) error {
	var ruleYamlMap map[string]ruleYaml

	if file == "" {
		return nil
	}

	filename, _ := filepath.Abs(file)

	ruleBase := filepath.Base(filename)
	if filepath.Ext(ruleBase) == ".rule" {
		ruleBase = ruleBase[0 : len(ruleBase)-5]
	}

	yamlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlData, &ruleYamlMap)
	if err != nil {
		return err
	}

	for k, v := range ruleYamlMap {
		rule, err := models.NewRule(ruleBase+"."+k, v.Match)
		if err != nil {
			return err
		}

		for _, e := range v.Unmatch {
			rule.AddUnmatch(e)
		}
		s.AddRule(*rule, defaulEnabled)
	}

	return nil
}

// LoadRulesFromFile loads rules from all YAML files inside a directory.
func (s *Seekret) LoadRulesFromDir(dir string, defaulEnabled bool) error {
	fi, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		err := fmt.Errorf("%s is not a directory", dir)
		return err
	}

	fileList, err := filepath.Glob(dir + "/*")
	if err != nil {
		return err
	}
	for _, file := range fileList {
		if strings.HasSuffix(file, ".rule") {
			err := s.LoadRulesFromFile(file, defaulEnabled)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// LoadRulesFromFile loads rules from all YAML files inside different
// directories separated by ':'.
func (s *Seekret) LoadRulesFromPath(path string, defaulEnabled bool) error {
	if path == "" {
		path = DefaultRulesPath()
	}
	dirList := strings.Split(path, ":")
	for _, dir := range dirList {
		err := s.LoadRulesFromDir(dir, defaulEnabled)
		if err != nil {
			return err
		}
	}
	return nil
}

const defaultRulesDir = "$GOPATH/src/github.com/apuigsech/seekret/rules"

// DefaultRulesPath return the default PATH that contains rules.
func DefaultRulesPath() string {
	rulesPath := os.Getenv("SEEKRET_RULES_PATH")
	if rulesPath == "" {
		rulesPath = os.ExpandEnv(defaultRulesDir)
	}
	return rulesPath
}

// ListRules return an array with all loaded rules.
func (s *Seekret) ListRules() []models.Rule {
	return s.ruleList
}

// EnableRule enables rules that match with a regular expression.
func (s *Seekret) EnableRule(name string) error {
	return setRuleEnabled(s.ruleList, name, true)
}

// DisableRule disables rules that match with a regular expression.
func (s *Seekret) DisableRule(name string) error {
	return setRuleEnabled(s.ruleList, name, false)
}

func setRuleEnabled(ruleList []models.Rule, name string, enabled bool) error {
	// TODO: implement regular expression.
	found := false
	for i, r := range ruleList {
		if r.Name == name {
			found = true
			ruleList[i].Enabled = enabled
		}
	}
	if !found {
		err := fmt.Errorf("%s rule not found", name)
		return err
	}

	return nil
}

// LoadObjects loads objects form an specific source. It can load objects from
// different source types, that are implemented following the SourceType
// interface.
func (s *Seekret) LoadObjects(st SourceType, source string, opt LoadOptions) error {
	objectList, err := st.LoadObjects(source, opt)
	if err != nil {
		return err
	}
	s.objectList = append(s.objectList, objectList...)
	return nil
}

// GroupObjectsByMetadata returns a map with all objects grouped by specific
// metadata key.
func (s *Seekret) GroupObjectsByMetadata(k string) map[string][]models.Object {
	return models.GroupObjectsByMetadata(s.objectList, k)
}

// GroupObjectsByPrimaryKeyHash returns a map with all objects grouped by
// the primary key hash, that is calculated from all metadata keys with the
// primary attribute.
// All returned objects could have the same content, even if are not the same.
func (s *Seekret) GroupObjectsByPrimaryKeyHash() map[string][]models.Object {
	return models.GroupObjectsByPrimaryKeyHash(s.objectList)
}

type exceptionYaml struct {
	Rule    *string
	Object  *string
	Line    *int
	Content *string
}

// AddException adds a new exception into the context.
func (s *Seekret) AddException(exception models.Exception) {
	s.exceptionList = append(s.exceptionList, exception)
}

// LoadExceptionsFromFile loads exceptions from a YAML file.
func (s *Seekret) LoadExceptionsFromFile(file string) error {
	var exceptionYamlList []exceptionYaml

	if file == "" {
		return nil
	}

	filename, _ := filepath.Abs(file)
	yamlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlData, &exceptionYamlList)
	if err != nil {
		return err
	}

	for _, v := range exceptionYamlList {
		x := models.NewException()

		if v.Rule != nil {
			err := x.SetRule(*v.Rule)
			if err != nil {
				return err
			}
		}

		if v.Object != nil {
			err := x.SetObject(*v.Object)
			if err != nil {
				return err
			}
		}

		if v.Line != nil {
			err := x.SetNline(*v.Line)
			if err != nil {
				return err
			}
		}

		if v.Content != nil {
			err := x.SetContent(*v.Content)
			if err != nil {
				return err
			}
		}

		s.AddException(*x)
	}

	return nil
}

// ListSecrets return an array with all found secrets after the inspection.
func (s *Seekret) ListSecrets() []models.Secret {
	return s.secretList
}
