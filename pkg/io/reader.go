package io

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	base "github.com/alanpjohn/go-cdcl/pkg/base"
	handler "github.com/alanpjohn/go-cdcl/pkg/error"
	logger "github.com/alanpjohn/go-cdcl/pkg/logger"
)

func Process(f *os.File) (sat base.SATFile, err error) {

	fileScanner := bufio.NewScanner(f)

	fileScanner.Split(bufio.ScanLines)

	var atomCount int
	var clauseCount int
	var clauses []base.Disjunction

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
			var cl base.Disjunction
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
				cl = append(cl, base.Literal(val))
			}
			logger.Info(fmt.Sprintf("%v", cl))
			clauses = append(clauses, cl)
		}

	}
	sat.AtomCount = uint(atomCount)
	sat.ClauseCount = uint(clauseCount)
	sat.Clauses = clauses
	logger.Info("Processed SAT file")

	defer f.Close()

	return base.SATFile{}, err
}

func ReadFile(filename string) (out base.SATFile, err error) {
	if filename == "" {
		return base.SATFile{}, handler.Throw("Please input a file", nil)
	}
	exists, e := fileExists(filename)
	if !exists {
		return base.SATFile{}, handler.Throw("The file provided does not exist", e)
	}
	file, e := os.Open(filename)
	if e != nil {
		return base.SATFile{}, handler.Throw("File could not be read", err)
	}
	return Process(file)
}

func fileExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}
