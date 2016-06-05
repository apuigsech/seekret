package main

import (
	"os"
	"bufio"
	"path/filepath"
	"regexp"
	seekret "github.com/apuigsech/seekret/lib"
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
		s.AddRule(rule)
    }

    if err := scanner.Err(); err != nil {
        return err
    }

	return nil
}