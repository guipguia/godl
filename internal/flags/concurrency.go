package flags

import (
	"github.com/urfave/cli/v2"
)

func Concurrency() *cli.IntFlag {
	return &cli.IntFlag{
		Name:     "concurrency",
		Usage:    "Specify number of goroutines, that will divide the file and download at the same time.",
		Required: false,
		Value:    4,
		Aliases:  []string{"c"},
		EnvVars:  []string{"GL_CONCURRENCY"},
	}
}
