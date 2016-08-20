package models

import (
	"regexp"
)

type Exception struct {
	Name    string
	Rule    *regexp.Regexp
	Object  *regexp.Regexp
	Nline   *int
	Content *regexp.Regexp
}

func NewException() *Exception {
	x := &Exception{
		Name: "anonymous",
	}
	return x
}

func (x *Exception) SetRule(rule string) error {
	ruleRegexp, err := regexp.Compile("(?i)" + rule)
	if err != nil {
		return err
	}
	x.Rule = ruleRegexp
	return nil
}

func (x *Exception) SetObject(object string) error {
	objectRegexp, err := regexp.Compile("(?i)" + object)
	if err != nil {
		return err
	}
	x.Object = objectRegexp
	return nil
}

func (x *Exception) SetNline(nLine int) error {
	x.Nline = &nLine
	return nil
}

func (x *Exception) SetContent(content string) error {
	contentRegexp, err := regexp.Compile("(?i)" + content)
	if err != nil {
		return err
	}
	x.Content = contentRegexp
	return nil
}

func (x *Exception) Run(s *Secret) bool {
	match := true

	if match && x.Rule != nil && !x.Rule.MatchString(s.Rule.Name) {
		match = false
	}

	if match && x.Object != nil && !x.Object.MatchString(s.Object.Name) {
		match = false
	}

	if match && x.Nline != nil && *x.Nline != s.Nline {
		match = false
	}

	if match && x.Content != nil && !x.Content.MatchString(s.Line) {
		match = false
	}

	return match
}
