// Copyright 2016 - Authors included on AUTHORS file.
//
// Use of this source code is governed by a Apache License
// that can be found in the LICENSE file.

package models

import (
	"fmt"
	"regexp"
)

// Represents a Rule.
type Rule struct {
	// Contains the name of the rule.
	Name string

	// Specifies if the rule is enabled or not.
	Enabled bool

	// All lines of the content are analised separatelly.
	// For a line to be considered a secret it should match the Match regular
	// expression and not match any of the regular expressions contained on the
	// Unmacth array.
	Match   *regexp.Regexp
	Unmatch []*regexp.Regexp
}

type RunResult struct {
	Line  string
	Nline int
}

// NewRule creates a new rule.
func NewRule(name string, match string) (*Rule, error) {
	matchRegexp, err := regexp.Compile("(?i)" + match)
	if err != nil {
		return nil, err
	}
	if err != nil {
		fmt.Println(err)
	}

	r := &Rule{
		Enabled: false,
		Name:    name,
		Match:   matchRegexp,
	}
	return r, nil
}

// Enable marks the rule as enabled.
func (r *Rule) Enable() {
	r.Enabled = true
}

// Enable marks the rule as disabled.
func (r *Rule) Disable() {
	r.Enabled = false
}

// AddUnmatch adds a refular expression into the unmatch list.
func (r *Rule) AddUnmatch(unmatch string) error {
	unmatchRegexp, err := regexp.Compile("(?i)" + unmatch)
	if err != nil {
		return err
	}

	r.Unmatch = append(r.Unmatch, unmatchRegexp)

	return nil
}

// Run executes the rule on a line passed in and returns a boolean
func (r *Rule) Run(line string) bool {

	result := false

	if r.Match.MatchString(line) {
		unmatch := false
		for _, Unmatch := range r.Unmatch {
			if Unmatch.MatchString(line) {
				unmatch = true
			}
		}

		if !unmatch {
			result = true
		}
	}

	return result
}
