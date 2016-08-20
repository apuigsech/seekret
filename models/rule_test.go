package models

import (
	//	"fmt"
	"testing"
)

type NewRuleSample struct {
	name    string
	match   string
	unmatch []string
	ok      bool
}

func TestNewRule(t *testing.T) {
	testSamples := []NewRuleSample{
		{
			name:  "rule_1",
			match: "match_1",
			ok:    true,
		},
		{
			name:  "rule_2",
			match: "*",
			ok:    false,
		},
		{
			name:  "rule_3",
			match: "match_3",
			unmatch: []string{
				"unmatch_3_1",
				"unmatch_3_2",
				"unmatch_3_3",
			},
			ok: true,
		},
		{
			name:  "rule_4",
			match: "match_4",
			unmatch: []string{
				"unmatch_4_1",
				"*",
				"unmatch_4_3",
			},
			ok: false,
		},
	}

	for _, ts := range testSamples {
		if !testNewRuleSample(ts) {
			t.Error("unexpected new rule")
		}
	}
}

func testNewRuleSample(ts NewRuleSample) bool {
	r, err := NewRule(ts.name, ts.match)

	ok_1 := (err == nil)

	ok_2 := true
	for _, unmatch := range ts.unmatch {
		err := r.AddUnmatch(unmatch)

		if err != nil {
			ok_2 = false
		}
	}

	ok := (ok_1 && ok_2)

	if ok == ts.ok {
		return true
	} else {
		return false
	}

	return true
}

type RunRuleSample struct {
	match      string
	unmatch    []string
	content    []byte
	expResults []RunResult
	ok         bool
}

func TestRunRule(t *testing.T) {
	testSamples := []RunRuleSample{
		{
			match: ".*TEST_1.*",
			content: []byte(
				"xxx\n" +
					"yyy\n" +
					"xxx TEST_1 yyy\n" +
					"TEST_1",
			),
			expResults: []RunResult{
				{
					Line:  "xxx TEST_1 yyy",
					nLine: 3,
				},
				{
					Line:  "TEST_1",
					nLine: 4,
				},
			},
			ok: true,
		},
		{
			match: ".*TEST_2.*",
			unmatch: []string{
				".*yyy.*",
				".*zzz.*",
			},
			content: []byte(
				"xxx\n" +
					"yyy\n" +
					"xxx TEST_2 yyy\n" +
					"xxx TEST_2 zzz\n" +
					"xxx TEST_2 www\n" +
					"TEST_2",
			),
			expResults: []RunResult{
				{
					Line:  "xxx TEST_2 www",
					nLine: 5,
				},
				{
					Line:  "TEST_2",
					nLine: 6,
				},
			},
			ok: true,
		},
	}

	for _, ts := range testSamples {
		if !testRunRuleSample(ts) {
			t.Error("unexpected run rule results")
		}
	}
}

func testRunRuleSample(ts RunRuleSample) bool {
	r, err := NewRule("rule", ts.match)
	if err != nil {
		return false
	}

	for _, unmatch := range ts.unmatch {
		err := r.AddUnmatch(unmatch)
		if err != nil {
			return false
		}
	}

	results := r.Run(ts.content)

	ok := true
	for _, res := range results {
		ok_s := false
		for _, expRes := range ts.expResults {
			if res.nLine == expRes.nLine && res.Line == expRes.Line {
				ok_s = true
			}
		}

		if ok_s == false {
			ok = false
		}
	}

	if ok == ts.ok {
		return true
	} else {
		return false
	}
}
