// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mxwell/wac/model"
	"github.com/spf13/cobra"
)

type Outcome struct {
	exec_time      time.Duration
	output_differs bool
}

var SolutionCommand string

func getSolutionCommand() *exec.Cmd {
	tokens := strings.Split(SolutionCommand, " ")
	if len(tokens) == 1 {
		return exec.Command(tokens[0])
	} else if len(tokens) > 1 {
		return exec.Command(tokens[0], tokens[1:]...)
	} else {
		panic("Empty command to run solution!")
	}
}

func doRun(inputPath string, resultPath string) (error, time.Duration) {
	inputReader, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open test input: %s", err), 0
	}
	defer inputReader.Close()
	resultWriter, err := os.OpenFile(resultPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file to write output: %s", err), 0
	}
	defer resultWriter.Close()
	var stderr bytes.Buffer

	command := getSolutionCommand()
	command.Stdin = inputReader
	command.Stdout = resultWriter
	command.Stderr = &stderr

	start := time.Now()
	err = command.Run()
	elapsed := time.Since(start)

	if stderr.Len() > 0 {
		fmt.Println("<stderr>")
		stderr.WriteTo(os.Stdout)
	}

	return err, elapsed
}

func readWholeLine(r *bufio.Reader) (string, error) {
	result := make([]byte, 0)
	for {
		if bytes, isPrefix, err := r.ReadLine(); err == nil {
			result = append(result, bytes...)
			if !isPrefix {
				break
			}
		} else {
			return "", err
		}
	}
	return strings.TrimRight(string(result), " "), nil
}

/* true if result is different from what was expected */
func checkOutput(expectedPath string, resultPath string) (bool, error) {
	a0, err := os.Open(expectedPath)
	if err != nil {
		return true, fmt.Errorf("failed to open file with expected output: %s", err)
	}
	defer a0.Close()
	a := bufio.NewReader(a0)

	b0, err := os.Open(resultPath)
	if err != nil {
		return true, fmt.Errorf("failed to open file with solution output: %s", err)
	}
	defer b0.Close()
	b := bufio.NewReader(b0)

	adone := false
	bdone := false
	for !adone && !bdone {
		aline, aerr := readWholeLine(a)
		bline, berr := readWholeLine(b)
		if aerr == io.EOF {
			adone = true
			aerr = nil
		}
		if berr == io.EOF {
			bdone = true
			berr = nil
		}
		if aerr != nil {
			return true, fmt.Errorf("error when reading from '%s': %s", expectedPath, aerr)
		}
		if berr != nil {
			return true, fmt.Errorf("error when reading from '%s': %s", resultPath, berr)
		}
		if aline != bline {
			return true, nil
		}
	}
	return false, nil
}

func runSingleTest(taskDir string, testToken string) (*Outcome, error) {
	testPathPrefix := filepath.Join(taskDir, testToken)
	inputPath := testPathPrefix + ".in"
	outputPath := testPathPrefix + ".out"
	resultPath := testPathPrefix + ".result"

	err, elapsed := doRun(inputPath, resultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to run solution: %s", err)
	}

	diff, err := checkOutput(outputPath, resultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check output: %s", err)
	}
	return &Outcome{elapsed, diff}, nil
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run built solution on test cases",
	Long: `Run built solution on test cases until an error happens or force to go on through all test cases.
Set of test cases could be specified in command arguments. If no arguments are given, then all available tests are used.`,
	Run: func(cmd *cobra.Command, args []string) {
		contest, err := model.LocateContest()
		if err != nil {
			log.Fatalf("ERROR %s\n", err)
		}
		taskToken, err := model.DetermineCurrentTask(contest)
		if err != nil {
			log.Fatalf("ERROR can't determine current task: %s\n", err)
		}
		task, ok := contest.Tasks[taskToken]
		if !ok {
			log.Fatalf("ERROR contest have no task for the working directory")
		}
		if len(task.TestTokens) == 0 {
			fmt.Println("No tests.")
			return
		}
		for _, testToken := range task.TestTokens {
			fmt.Printf("[%s] ... ", testToken)
			outc, err := runSingleTest(filepath.Join(contest.RootDir, taskToken), testToken)
			if err != nil {
				log.Fatalf("ERROR failed to run test '%s': %s", testToken, err)
			}
			var msg string
			if outc.output_differs {
				msg = "Differs"
			} else {
				msg = "Ok"
			}
			fmt.Printf("%s -- %dms\n", msg, int(outc.exec_time/1000000))
			if outc.output_differs {
				break
			}
		}
	},
}

func init() {
	runCmd.Flags().StringVarP(&SolutionCommand, "command", "c", "./main", "Command to execute solution")
	RootCmd.AddCommand(runCmd)
}
