// The main file to run for execution
package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	err "github.com/alanpjohn/go-cdcl/pkg/error"
)

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func solve(cCtx *cli.Context) error {
	if cCtx.Bool("verbose") {
		log.Println("Verbose flag detected")
	}
	if isInputFromPipe() {
		log.Println("Input from Standard Pipe")
	}
	if cCtx.String("file") != "" {
		log.Print("Input from flag", cCtx.String("file"))
	}
	return err.SolverError{Message: "No input provided"}
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
		log.Fatal(err)
	}
}
