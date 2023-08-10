package commands

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/guipguia/godl/internal/download"
	"github.com/guipguia/godl/internal/util"
	"github.com/urfave/cli/v2"
)

// doVerifications will do argument checks
func doInitialVerifications(ctx *cli.Context) {
	if ctx.Args().Len() == 0 {
		fmt.Println("Please provide a URL to download something.")
		os.Exit(1)
	}

	urlString := ctx.Args().Get(0)

	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		fmt.Printf("Could not parse given url: %v", err)
		os.Exit(1)
	}
}

// getFileName will check if any parameter was given
// to modify file name.
func getFileName(ctx *cli.Context) string {
	var outputFileName string = util.GetBaseName(ctx.Args().Get(0))

	if len(ctx.String("filename")) != 0 && ctx.String("filename") != util.GetBaseName(ctx.Args().Get(0)) {
		outputFileName = ctx.String("filename")
	}

	return outputFileName
}

// doForceVerification will check force flag
// and if file exists and if we should ovewrite it
func doForceVerification(f bool, path string) {
	if !f {
		_, err := os.Stat(path)
		if !os.IsNotExist(err) {
			fmt.Printf("File: %s already exists, if you want to overwrite it, please use flag --force", path)
			os.Exit(1)
		}
	}
}

// Download is the command function to gether arguments
// so we can start download
func MainCmd() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		doInitialVerifications(ctx)

		downloadLink := ctx.Args().First()
		downloadDir := ctx.String("dir")
		fileName := getFileName(ctx)
		concurrency := ctx.Int("concurrency")
		force := ctx.Bool("force")
		fullPath := path.Join(downloadDir, fileName)
		fullPath = path.Clean(fullPath)

		doForceVerification(force, fullPath)

		download.StartDownload(downloadLink, fullPath, concurrency)
		return nil
	}
}
