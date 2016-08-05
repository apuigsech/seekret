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

const DefaultRulesPath = "$GOPATH/src/github.com/apuigsech/seekret/rules"

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
		path = os.ExpandEnv(DefaultRulesPath)
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

func (s *Seekret) EnableRule(name string) {
	setRuleEnabled(s.ruleList, name, true)
}

func (s *Seekret) DisableRule(name string) {
	setRuleEnabled(s.ruleList, name, false)
}

func setRuleEnabled(ruleList []Rule, name string, enabled bool) {
	for i, r := range ruleList {
		fmt.Println(r.Name, name)
		if r.Name == name {
			ruleList[i].Enabled = enabled
		}
	}
}