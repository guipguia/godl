package flags

import (
	"github.com/urfave/cli/v2"
)

func Output() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "Change the name of the file, (if you want to change the directory, please also use -d or --dir)",
	}
}
