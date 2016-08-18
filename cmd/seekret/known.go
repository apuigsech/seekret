package main

import (
	"bufio"
	"github.com/apuigsech/seekret"
	"os"
	"path/filepath"
	"regexp"
)

func LoadKnownFromFile(s *seekret.Seekret, file string) error {
	if file == "" {
		return nil
	}

	filename, _ := filepath.Abs(file)

	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		rule := seekret.Rule{
			Name:  "known",
			Match: regexp.MustCompile("(?i)" + scanner.Text()),
		}
		s.AddRule(rule, true)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
