package flags

import (
	"fmt"
	"os"

	"github.com/guipguia/godl/internal/util"
	"github.com/urfave/cli/v2"
)

// Set app version
func SetVersion() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print current godl version.",
	}

	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Printf("version=%s \n", ctx.App.Version)
	}
}

func Version() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print current godl version.",
		Action: func(ctx *cli.Context, b bool) error {
			if b {
				fmt.Printf("godl version %s\n", util.VersionNumber)
			}
			os.Exit(0)
			return nil
		},
	}
}
