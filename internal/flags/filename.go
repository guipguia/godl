package flags

import (
	"github.com/urfave/cli/v2"
)

func Filename() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "filename",
		Aliases: []string{"f"},
		Usage:   "Change the name of the file, (if you want to change the directory, please also use -d or --dir)",
		EnvVars: []string{"GL_FILENAME"},
	}
}
