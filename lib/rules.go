package lib

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type ruleYaml struct {
	Match   string
	Unmatch []string
}

type Rule struct {
	Name    string
	Match   *regexp.Regexp
	Unmatch []*regexp.Regexp
}

func (s *Seekret) LoadRulesFromFile(file string) error {
	var ruleYamlMap map[string]ruleYaml

	if file == "" {
		return nil
	}

	filename, _ := filepath.Abs(file)
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
		s.ruleList = append(s.ruleList, rule)
	}

	return nil
}

func (s *Seekret) LoadRulesFromDir(dir string) error {
	fileList, err := filepath.Glob(dir + "/*")
	if err != nil {
		return err
	}
	for _, file := range fileList {
		if strings.HasSuffix(file, ".rule") == true {
			err := s.LoadRulesFromFile(file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Seekret) LoadRulesFromPath(path string) error {
	dirList := strings.Split(path, ":")
	for _, dir := range dirList {
		err := s.LoadRulesFromDir(dir)
		if err != nil {
			return err
		}
	}
	return nil
}
