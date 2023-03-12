// package main contains the main execution function that starts the CLI
package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2" // CLI framework for a better user experience

	handler "github.com/alanpjohn/go-cdcl/pkg/error"
	reader "github.com/alanpjohn/go-cdcl/pkg/io"
	logger "github.com/alanpjohn/go-cdcl/pkg/logger"
	solver "github.com/alanpjohn/go-cdcl/pkg/solver"
	types "github.com/alanpjohn/go-cdcl/pkg/types"
)

// version string for versioning program externally using compiler options
var version string

/*
Checks if the input is coming from stdin pipe
*/
func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

/*
Takes configurable parameters from the CLI and starts the solver

Tha Main command
*/
func solve(cCtx *cli.Context) error {
	logger.Verbosity = cCtx.Bool("verbose")

	filename := cCtx.String("file")

	if filename == "" && !isInputFromPipe() {
		return handler.Throw("No input was provided", nil)
	}

	var (
		sol      types.Solver  // The Solver class with the methods implemented for CDCL
		sat      types.SATFile // Contains all the information extracted as DIMCAS Format
		err      error
		solution types.Solution // SATISFIABLE or UNSATISFIABLE or UNKNOWN
	)

	if isInputFromPipe() {
		// The input is coming for stdin, our Process function returns an instance of SATFile
		logger.Info("Recieved Input for stdin pipe")
		if sat, err = reader.Process(os.Stdin); err != nil {
			return err
		}
	}

	if filename != "" {
		logger.Info("Input from flag")
		// The input has to read from the file, our Readfile function reads the file
		// and then calls Process (internally) to return an instance of SATFile
		if sat, err = reader.ReadFile(filename); err != nil {
			return err
		}
	}

	// Initalize the Solver with the SATFile
	if sol, err = solver.InitializeBaseSolver(sat, cCtx.Bool("experimental")); err != nil {
		return err
	}
	logger.Info("Solver initialized")
	solution, err = sol.Solve() // Get Solution
	fmt.Print(solution.String())
	return err
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
				Usage:    "Switches on detailed logging for cdcl solver",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "experimental",
				Aliases:  []string{"e"},
				Value:    false,
				Usage:    "use experimental features",
				Required: false,
			},
		},
		Action: solve,
	})

	if err := app.Run(os.Args); err != nil {
		logger.Error(err)
	}
}
