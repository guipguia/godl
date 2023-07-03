package flags

import (
	"os"

	"github.com/urfave/cli/v2"
)

func Directory() *cli.StringFlag {
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &cli.StringFlag{
		Name:    "dir",
		Aliases: []string{"d"},
		Value:   curDir,
		Usage:   "Specify where you want to save the file.",
		EnvVars: []string{"GL_DIRECTORY"},
	}
}
