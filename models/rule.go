package models

import (
	"fmt"
	"regexp"
	"bufio"
	"bytes"
)


type Rule struct {
	Name    string
//	ObjectMatch *regexp.Regexp

	Enabled bool
	Match   *regexp.Regexp
	Unmatch []*regexp.Regexp
}

type RunResult struct {
	Line string
	nLine int
}

func NewRule(name string, match string) (*Rule,error) {
	matchRegexp,err := regexp.Compile("(?i)" + match)
	if err != nil {
		return nil,err
	}
	if err != nil {
		fmt.Println(err)
	}

	r := &Rule{
		Enabled: false,
		Name: name,
		Match: matchRegexp,
	}
	return r,nil
}

func (r *Rule)Enable() {
	r.Enabled = true
}

func (r *Rule)Disable() {
	r.Enabled = false
}

func (r *Rule)AddUnmatch(unmatch string) error {
	unmatchRegexp,err := regexp.Compile("(?i)" + unmatch)
	if err != nil {
		return err
	}

	r.Unmatch = append(r.Unmatch, unmatchRegexp)

	return nil
}


func (r *Rule)Run(content []byte) []RunResult {
	var results []RunResult 

	b := bufio.NewScanner(bytes.NewReader(content))

	nLine := 0
	for b.Scan() {
		nLine = nLine + 1
		line := b.Text()

		if r.Match.MatchString(line) {
			unmatch := false
			for _, Unmatch := range r.Unmatch {
				if Unmatch.MatchString(line) {
					unmatch = true
				}
			}

			if !unmatch {
				results = append(results, RunResult{
					Line: line,
					nLine: nLine,
				})
			}
		}		
	}

	return results
}
