// The main file to run for execution
package main

import (
	"os"

	"github.com/urfave/cli/v2"

	dbg "github.com/alanpjohn/go-cdcl/pkg/debug"
)

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func solve(cCtx *cli.Context) error {
	if cCtx.Bool("verbose") {
		dbg.Info("Verbose flag detected")
	}
	if isInputFromPipe() {
		dbg.Info("Input from Standard Pipe")
	}
	if cCtx.String("file") != "" {
		dbg.Info("Input from flag")
	}
	return dbg.ThrowSolverError("No input provided", nil)
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
		dbg.Error(err)
	}
}
