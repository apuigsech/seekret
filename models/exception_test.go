package models

import (
	"regexp"
	"testing"
)

type NewExceptionSample struct {
	rule    string
	object  string
	nline   int
	content string
	ok		bool
}

func TestNewException(t *testing.T) {
	testSamples := []NewExceptionSample{
		{
			rule: "rule_1",
			object: "object_1",
			nline: 0,
			content: "content_1",
			ok: true,
		},
		{
			rule: "rule_2",
			object: "object_2",
			nline: 0,
			content: "*",
			ok: false,
		},
	}

	for _,ts := range testSamples {
		if !testNewExceptionSample(ts) {
			t.Error("unexpected new exception")
		}
	}
}

func testNewExceptionSample(ts NewExceptionSample) bool {
	x := NewException()

	ok := true

	err := x.SetRule(ts.rule)
	if err != nil {
		ok = false
	}

	err = x.SetObject(ts.object)
	if err != nil {
		ok = false
	}

	err = x.SetNline(ts.nline)
	if err != nil {
		ok = false
	}

	err = x.SetContent(ts.content)
	if err != nil {
		ok = false
	}

	if ok == ts.ok {
		return true
	} else {
		return false
	}	
}

type RunExceptionSample struct {
	secret *Secret
	exception *Exception
	expResult bool
}


func TestRunException(t *testing.T) {
	testSamples := []RunExceptionSample{
		{	
			secret: &Secret{
				Object: &Object{
					Name: "object_1",			
				},
				Rule: &Rule{
					Name: "rule_1",
				},
				Nline: 1,
				Line: "secret_1",
			},
			exception: &Exception{
				Rule: regexp.MustCompile("rule_1"),
			},
			expResult: true,
		},
		{	
			secret: &Secret{
				Object: &Object{
					Name: "object_2",			
				},
				Rule: &Rule{
					Name: "rule_2",
				},
				Nline: 2,
				Line: "secret_2",
			},
			exception: &Exception{
				Rule: regexp.MustCompile("rule_X"),
			},
			expResult: false,
		},
		{	
			secret: &Secret{
				Object: &Object{
					Name: "object_3",			
				},
				Rule: &Rule{
					Name: "rule_3",
				},
				Nline: 3,
				Line: "secret_3",
			},
			exception: &Exception{
				Rule: regexp.MustCompile("rule_3"),
				Object: regexp.MustCompile("object_3"),
			},
			expResult: true,
		},
		{	
			secret: &Secret{
				Object: &Object{
					Name: "object_4",			
				},
				Rule: &Rule{
					Name: "rule_4",
				},
				Nline: 4,
				Line: "secret_4",
			},
			exception: &Exception{
				Rule: regexp.MustCompile("rule_4"),
				Object: regexp.MustCompile("object_X"),
			},
			expResult: false,
		},
		{	
			secret: &Secret{
				Object: &Object{
					Name: "object_5",			
				},
				Rule: &Rule{
					Name: "rule_5",
				},
				Nline: 5,
				Line: "secret_5",
			},
			exception: &Exception{
				Rule: regexp.MustCompile("rule_5"),
				Object: regexp.MustCompile("object_5"),
				Nline: intPtr(5),
				Content: regexp.MustCompile("secret_5"),
			},
			expResult: true,
		},
		{	
			secret: &Secret{
				Object: &Object{
					Name: "object_6",			
				},
				Rule: &Rule{
					Name: "rule_6",
				},
				Nline: 6,
				Line: "secret_6",
			},
			exception: &Exception{
				Rule: regexp.MustCompile("rule_6"),
				Object: regexp.MustCompile("object_X"),
				Nline: intPtr(6),
				Content: regexp.MustCompile("secret_X"),
			},
			expResult: false,
		},
	}

	for _,ts := range testSamples {
		if !testRunExceptionSample(ts) {
			t.Error("unexpected run exception result")
		}
	}
}

func testRunExceptionSample(ts RunExceptionSample) bool {
	result := ts.exception.Run(ts.secret)

	if result == ts.expResult {
		return true
	} else {
		return false
	}
}

func intPtr(i int) *int {
	return &i
}