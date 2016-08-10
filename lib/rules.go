package lib

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"os"
)


type ruleYaml struct {
	ObjectMatch string
	Match   string
	Unmatch []string
}

type Rule struct {
	Enabled bool
	Name    string
	ObjectMatch *regexp.Regexp
	Match   *regexp.Regexp
	Unmatch []*regexp.Regexp
}

const DefaultRulesDir = "$GOPATH/src/github.com/apuigsech/seekret/rules"

func DefaultRulesPath() string {
	rulesPath := os.Getenv("SEEKRET_RULES_PATH")
	if rulesPath == "" {
		rulesPath = os.ExpandEnv(DefaultRulesDir)
	}
	return rulesPath
}

func (s *Seekret) ListRules() []Rule {
	return s.ruleList
}

func (s *Seekret) AddRule(rule Rule, enabled bool) {
	rule.Enabled = enabled
	s.ruleList = append(s.ruleList, rule)
}

func (s *Seekret) LoadRulesFromFile(file string, defaulEnabled bool) error {
	var ruleYamlMap map[string]ruleYaml

	if file == "" {
		return nil
	}

	filename, _ := filepath.Abs(file)

	ruleBase := filepath.Base(filename)
	if filepath.Ext(ruleBase) == ".rule" {
		ruleBase = ruleBase[0:len(ruleBase)-5]
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
		rule := Rule{
			Name:  ruleBase + "." + k,
			Match: regexp.MustCompile("(?i)" + v.Match),
		}
		for _, e := range v.Unmatch {
			rule.Unmatch = append(rule.Unmatch, regexp.MustCompile("(?i)"+e))
		}
		s.AddRule(rule, defaulEnabled)
	}

	return nil
}

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

func (s *Seekret) LoadRulesFromPath(path string, defaulEnabled bool) error {
	if path == "" {
		path = os.ExpandEnv(DefaultRulesDir)
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

func (s *Seekret) EnableRule(name string) (error) {
	return setRuleEnabled(s.ruleList, name, true)
}

func (s *Seekret) DisableRule(name string) (error) {
	return setRuleEnabled(s.ruleList, name, false)
}

func setRuleEnabled(ruleList []Rule, name string, enabled bool) (error) {
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