package models

import (
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

/*
type RunExceptionSample struct {
	object string
	content []byte
	rule *Rule
	exception *Exception
}

func TestRunException(t *testing.T) {
	rule_1,_ := NewRule("rule_1", ".*TEST_1.*")
	exception_1 := NewException()
	exception_1.SetRule("rule_1")
	
	rule_2,_ := NewRule("rule_2", ".*TEST_2.*")
	exception_2 := NewException()
	exception_2.SetObject("object_2")
	
	rule_3,_ := NewRule("rule_3", ".*TEST_3.*")
	exception_3 := NewException()
	exception_3.SetNline(3)

	rule_4,_ := NewRule("rule_4", "TEST_4")
	exception_4 := NewException()
	exception_4.SetContent(".*xxx.*")

	testSamples := []RunExceptionSample{
		{	
			object: "object_1",
			content: []byte(
				"xxx\n" +
				"yyy\n" +
				"xxx TEST_1 yyy\n" +
				"xxx TEST_1 zzz\n" +
				"xxx TEST_1 www\n" +
				"TEST_1",
			),
			rule: rule_1,
			exception: exception_1,
		},
		{	
			object: "object_2",
			content: []byte(
				"xxx\n" +
				"yyy\n" +
				"xxx TEST_2 yyy\n" +
				"xxx TEST_2 zzz\n" +
				"xxx TEST_2 www\n" +
				"TEST_2",
			),
			rule: rule_2,
			exception: exception_2,
		},
		{	
			object: "object_3",
			content: []byte(
				"xxx\n" +
				"yyy\n" +
				"xxx TEST_3 yyy\n" +
				"xxx TEST_3 zzz\n" +
				"xxx TEST_3 www\n" +
				"TEST_3",
			),
			rule: rule_3,
			exception: exception_3,
		},
		{	
			object: "object_4",
			content: []byte(
				"xxx\n" +
				"yyy\n" +
				"xxx TEST_4 yyy\n" +
				"xxx TEST_4 zzz\n" +
				"xxx TEST_4 www\n" +
				"TEST_4",
			),
			rule: rule_4,
			exception: exception_4,
		},		
	}

	for _,ts := range testSamples {
		if !testRunExceptionSample(ts) {
			t.Error("unexpected run exception result")
		}
	}
}


func testRunExceptionSample(ts RunExceptionSample) bool {
	return true
}
*/
