package main

import (
	"fmt"
	seekret "github.com/apuigsech/seekret/lib"
	"github.com/urfave/cli"
	"os"
)

var s *seekret.Seekret

func main() {
	s = seekret.NewSeekret()

	app := cli.NewApp()

	app.Name = "seekret"
	app.Version = "0.0.1"
	app.Usage = "seek for secrets on various sources."

	app.Author = "Albert Puigsech Galicia"
	app.Email = "albert@puigsech.com"

	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "exception, x",
			Usage: "load exceptions from `FILE`.",
		},
		cli.StringFlag{
			Name:   "rules",
			Usage:  "`PATH` with rules.",
			EnvVar: "SEEKRET_RULES_PATH",
		},
		cli.StringFlag{
			Name:  "format, f",
			Usage: "specify the output format.",
			Value: "human",
		},
		// TODO: To be implemented.
		/*
			cli.StringFlag{
				Name: "groupby, g",
				Usage: "Group output by specific field",
			},
		*/
		cli.StringFlag{
			Name:  "known, k",
			Usage: "load known secrets from `FILE`.",
		},
		cli.IntFlag{
			Name:  "workers, w",
			Usage: "number of workers used for the inspection",
			Value: 4,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:     "git",
			Usage:    "seek for seecrets on a git repository",
			Category: "seek",
			Action:   seekretGit,

			Flags: []cli.Flag{
				// TODO: To be implemented.
				/*
					cli.BoolFlag{
						Name: "recursive, r",
					},
					cli.BoolFlag{
						Name: "all, a",
					},
					cli.StringFlag{
						Name: "branches, b",
					},
				*/
				cli.IntFlag{
					Name:  "count, c",
					Usage: "number of commits to inspect (0 = all)",
					Value: 0,
				},
			},
		},
		{
			Name:     "dir",
			Usage:    "seek for seecrets on a directory",
			Category: "seek",
			Action:   seekretDir,

			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "recursive, r",
				},
				cli.BoolFlag{
					Name: "hidden",
				},
			},
		},
	}

	app.Before = seekretBefore
	app.After = seekretAfter

	app.Run(os.Args)
}

func seekretBefore(c *cli.Context) error {
	var err error

	rulesPath := c.String("rules")

	err = s.LoadRulesFromPath(rulesPath)
	if err != nil {
		return err
	}

	LoadKnownFromFile(s, c.String("known"))

	err = s.LoadExceptionsFromFile(c.String("exception"))
	if err != nil {
		return err

	}

	return nil
}

func seekretDir(c *cli.Context) error {
	source := c.Args().Get(0)
	if source == ""  {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	options := map[string]interface{}{
		"hidden":    c.Bool("hidden"),
		"recursive": c.Bool("recursive"),
	}

	err := s.LoadObjects(seekret.SourceTypeDir, source, options)
	if err != nil {
		return err
	}

	return nil
}

func seekretGit(c *cli.Context) error {
	source := c.Args().Get(0)
	if source == ""  {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	options := map[string]interface{}{
		"count": c.Int("count"),
	}

	err := s.LoadObjects(seekret.SourceTypeGit, source, options)
	if err != nil {
		return err
	}

	return nil
}

func seekretAfter(c *cli.Context) error {
	s.Inspect(c.Int("workers"))

	fmt.Println(FormatOutput(s.ListSecrets(), c.String("format")))

	return nil
}
