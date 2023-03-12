// The io package handles the processing of DIMCAS format file
package io

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	handler "github.com/alanpjohn/go-cdcl/pkg/error"
	logger "github.com/alanpjohn/go-cdcl/pkg/logger"
	types "github.com/alanpjohn/go-cdcl/pkg/types"
)

// Utility function to check if file exists
func fileExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

// Extracts SATFile from the input buffer
func Process(f *os.File) (sat types.SATFile, err error) {

	fileScanner := bufio.NewScanner(f)

	fileScanner.Split(bufio.ScanLines)

	var atomCount int
	var clauseCount int
	var clauses []types.Disjunction

	for fileScanner.Scan() {
		line := fileScanner.Text()
		line = strings.TrimSpace(line)
		items := strings.Split(line, " ")

		if items[0] == "c" {
			logger.Info("Comment: " + line)
		} else if items[0] == "p" {
			logger.Info("Atom Count :" + items[2])
			if atomCount, err = strconv.Atoi(items[2]); err != nil {
				return
			}

			logger.Info("Clause Count :" + items[3])
			if clauseCount, err = strconv.Atoi(items[3]); err != nil {
				return
			}
		} else {
			var val int
			var arr []int
			for i := 0; i < len(items); i++ {
				if val, err = strconv.Atoi(items[i]); err != nil {
					return
				}
				if (val > atomCount && val > 0) || (val < 0 && val < -atomCount) {
					return sat, handler.Throw("Invalid Literal found: "+items[i], nil)
				}
				if val == 0 {
					break
				}
				arr = append(arr, val)
			}
			var d types.Disjunction
			sort.Ints(arr)
			for _, v := range arr {
				d = append(d, types.Literal(v))
			}
			logger.Info(fmt.Sprintf("%v", d))
			clauses = append(clauses, d)
		}

	}
	sat.AtomCount = uint(atomCount)
	sat.ClauseCount = uint(clauseCount)
	sat.Clauses = clauses
	logger.Info("Processed SAT file")

	defer f.Close()

	return sat, err
}

// Open filename provided by user to process file contents into a SATFile
func ReadFile(filename string) (out types.SATFile, err error) {
	if filename == "" {
		return types.SATFile{}, handler.Throw("Please input a file", nil)
	}
	exists, e := fileExists(filename)
	if !exists {
		return types.SATFile{}, handler.Throw("The file provided does not exist", e)
	}
	file, e := os.Open(filename)
	if e != nil {
		return types.SATFile{}, handler.Throw("File could not be read", err)
	}
	return Process(file)
}
