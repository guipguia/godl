package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"sort"

	"github.com/guipguia/godl/internal/commands"
	"github.com/guipguia/godl/internal/flags"
	"github.com/guipguia/godl/internal/util"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "godl",
		Usage:       "CLI tool to download files from the internet.",
		Description: "Powerful CLI that enables asynchronous file downloads with ease and speed.",
		Flags: []cli.Flag{
			flags.Filename(),
			flags.Directory(),
			flags.Concurrency(),
			flags.Version(),
			flags.Force(),
		},
		Version: util.VersionNumber,
		Action:  commands.Run(),
	}

	if info, ok := debug.ReadBuildInfo(); !ok {
		app.Version = "Could not get version information. Please take a look into: https://github.com/guipguia/godl"
	} else {
		app.Version = info.Main.Version
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
