package flags

import "github.com/urfave/cli/v2"

func Force() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:     "force",
		Usage:    "If file already exists, it will overwrite it.",
		Required: false,
		Value:    false,
		// Aliases:  []string{""},
	}
}
