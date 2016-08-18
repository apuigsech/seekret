package main

import (
	"fmt"
	"github.com/apuigsech/seekret"
)

func FormatOutput(secretList []seekret.Secret, format string) string {
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

func formatOutputHuman(secretList []seekret.Secret) (string, error) {
	var out string
	for _, s := range secretList {
		out = out + fmt.Sprintf("%s\n\t%d: [%s] %s %s\n", s.Object.Name, s.Nline, s.Rule.Name, s.Line, s.Exception)
	}
	return out, nil
}

// TODO: Implement
func formatOutputJSON(secretList []seekret.Secret) (string, error) {
	return "Not implemented", nil
}

func formatOutputXML(secretList []seekret.Secret) (string, error) {
	return "Not implemented", nil
}

func formatOutputCSV(secretList []seekret.Secret) (string, error) {
	return "Not implemented", nil
}
