// Copyright 2016 - Authors included on AUTHORS file.
//
// Use of this source code is governed by a Apache License
// that can be found in the LICENSE file.

package models

// Represents a found secret.
type Secret struct {
	// Object in witch the secret is found.
	Object *Object

	// Rule that matches.
	Rule *Rule

	// Number of line in the content that contains the secret.
	Nline int

	// Content of the specific line.
	Line string

	// Specifies if this matches an exception too.
	Exception bool
}

// NewSecret creates a new secret.
func NewSecret(object *Object, rule *Rule, nLine int, line string) *Secret {
	s := &Secret{
		Object: object,
		Rule:   rule,
		Nline:  nLine,
		Line:   line,
	}
	return s
}

// SetException specifies that a found secret is an exception (of false positive).
func (s *Secret) SetException(exception bool) {
	s.Exception = exception
}
