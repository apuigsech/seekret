// Copyright 2016 - Authors included on AUTHORS file.
//
// Use of this source code is governed by a Apache License
// that can be found in the LICENSE file.

package models

import (
	"regexp"
)

// Represents an Exception. In order for a secret to be considered as exception
// all non-nill attributes should match with the secret information. That means
// it's considered like and AND statement.
type Exception struct {
	Name string

	// Regular expresion that should match the name of the rule.
	Rule *regexp.Regexp

	// Regular expresion that should match the name of the object.
	Object *regexp.Regexp

	// Number of line where the secret is found in the contect of the object.
	Nline *int

	// Regular expresion that should match the content of the object.
	Content *regexp.Regexp
}

// NewException creates a new exception.
func NewException() *Exception {
	x := &Exception{
		Name: "anonymous",
	}
	return x
}

// SetRule sets the regular expresion that should match the name of the rule.
func (x *Exception) SetRule(rule string) error {
	ruleRegexp, err := regexp.Compile("(?i)" + rule)
	if err != nil {
		return err
	}
	x.Rule = ruleRegexp
	return nil
}

// SetObject sets the regular expresion that should match the name of the
// object.
func (x *Exception) SetObject(object string) error {
	objectRegexp, err := regexp.Compile("(?i)" + object)
	if err != nil {
		return err
	}
	x.Object = objectRegexp
	return nil
}

// SetNline sets the number of line where secret should be found.
func (x *Exception) SetNline(nLine int) error {
	x.Nline = &nLine
	return nil
}

// SetContent sets the regular expresion that should match the content of the
// object.
func (x *Exception) SetContent(content string) error {
	contentRegexp, err := regexp.Compile("(?i)" + content)
	if err != nil {
		return err
	}
	x.Content = contentRegexp
	return nil
}

// Run executes the exception into a secret to determine if it's an exception
// or not.
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
