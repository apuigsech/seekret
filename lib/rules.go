package lib

import (
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
	Name    string
	ObjectMatch *regexp.Regexp
	Match   *regexp.Regexp
	Unmatch []*regexp.Regexp
}

func (s *Seekret) AddRule(rule Rule) {
	s.ruleList = append(s.ruleList, rule)
}

func (s *Seekret) LoadRulesFromFile(file string) error {
	var ruleYamlMap map[string]ruleYaml

	if file == "" {
		return nil
	}

	filename, _ := filepath.Abs(file)
//	x := filepath.Ext(filename)

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
			Name:  k,
			Match: regexp.MustCompile("(?i)" + v.Match),
		}
		for _, e := range v.Unmatch {
			rule.Unmatch = append(rule.Unmatch, regexp.MustCompile("(?i)"+e))
		}
		s.AddRule(rule)
	}

	return nil
}

func (s *Seekret) LoadRulesFromDir(dir string) error {
	fileList, err := filepath.Glob(dir + "/*")
	if err != nil {
		return err
	}
	for _, file := range fileList {
		if strings.HasSuffix(file, ".rule") {
			err := s.LoadRulesFromFile(file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Seekret) LoadRulesFromPath(path string) error {
	if path == "" {
		path = os.ExpandEnv(DefaultRulesPath)
	}
	dirList := strings.Split(path, ":")
	for _, dir := range dirList {
		err := s.LoadRulesFromDir(dir)
		if err != nil {
			return err
		}
	}
	return nil
}
