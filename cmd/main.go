// The main file to run for execution
package main

import (
	"os"

	"github.com/urfave/cli/v2"

	handler "github.com/alanpjohn/go-cdcl/pkg/error"
	reader "github.com/alanpjohn/go-cdcl/pkg/io"
	logger "github.com/alanpjohn/go-cdcl/pkg/logger"
)

var version string

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func solve(cCtx *cli.Context) error {
	logger.Verbosity = cCtx.Bool("verbose")

	filename := cCtx.String("file")

	if filename == "" && !isInputFromPipe() {
		return handler.Throw("No input was provided", nil)
	}

	if isInputFromPipe() {
		logger.Info("Recieved Input for stdin pipe")
		if _, err := reader.Process(os.Stdin); err != nil {
			return err
		}
	}

	if filename != "" {
		logger.Info("Input from flag")
		if _, err := reader.ReadFile(filename); err != nil {
			return err
		}
	}

	return nil
}

// Run CLI application which reads SAT file from standard input pipe and returns solution
func main() {
	app := (&cli.App{
		Name:  "gocdcl",
		Usage: "Pass SAT file as stdin pipe or using the -f/--file flag to run SAT solver",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Value:    "",
				Usage:    ".SAT file to be processed. This option is overridden if input provided by stdin pipe",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "verbose",
				Aliases:  []string{"v"},
				Value:    false,
				Usage:    "language for the greeting",
				Required: false,
			},
		},
		Action: solve,
	})

	if err := app.Run(os.Args); err != nil {
		logger.Error(err)
	}
}
