package main

import (
	"bufio"
	"github.com/apuigsech/seekret"
	"github.com/apuigsech/seekret/models"
	"os"
	"path/filepath"
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
		rule,err  := models.NewRule("known", scanner.Text())
		if err != nil {
			return err
		}
		s.AddRule(*rule, true)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
