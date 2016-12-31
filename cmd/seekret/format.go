// Copyright 2016 - Authors included on AUTHORS file.
//
// Use of this source code is governed by a Apache License
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/apuigsech/seekret/models"
)

func FormatOutput(secretList []models.Secret, format string) string {
	var out string

	switch format {
	case "human":
		out, _ = formatOutputHuman(secretList)
	case "json":
		out, _ = formatOutputJSON(secretList)
	case "xml":
		out, _ = formatOutputJSON(secretList)
	case "csv":
		out, _ = formatOutputCSV(secretList)
	default:
		out, _ = formatOutputHuman(secretList)
	}

	return out
}

func formatOutputHuman(secretList []models.Secret) (string, error) {
	var out string
	for _, s := range secretList {
		out = out + fmt.Sprintf("%s\n\t%d: [%s] %s %t\n", s.Object.Name, s.Nline, s.Rule.Name, s.Line, s.Exception)
	}
	return out, nil
}

// TODO: Implement
func formatOutputJSON(secretList []models.Secret) (string, error) {
	return "Not implemented", nil
}

func formatOutputXML(secretList []models.Secret) (string, error) {
	return "Not implemented", nil
}

func formatOutputCSV(secretList []models.Secret) (string, error) {
	return "Not implemented", nil
}
