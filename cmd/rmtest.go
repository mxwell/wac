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
	"fmt"
	"log"

	"github.com/mxwell/wac/model"
	"github.com/spf13/cobra"
)

var rmtestCmd = &cobra.Command{
	Use:   "rmtest <token>",
	Short: "Remove test from task",
	Long:  `Remove test case specified by the token from a task`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("ERROR single argument is required for the command")
		}
		testToken := args[0]

		contest, err := model.LocateContest()
		if err != nil {
			log.Fatalf("ERROR %s\n", err)
		}
		taskToken, err := model.DetermineCurrentTask(contest)
		if err != nil {
			log.Fatalf("ERROR can't determine current task: %s\n", err)
		}
		task, _ := contest.Tasks[taskToken]
		pos := -1
		tokens := task.TestTokens
		for i, token := range tokens {
			if token == testToken {
				pos = i
				break
			}
		}
		if pos == -1 {
			log.Fatalf("ERROR test %s not found", testToken)
		}
		task.TestTokens = append(tokens[:pos], tokens[pos+1:]...)
		contest.Tasks[taskToken] = task
		err = model.SaveContest(contest)
		if err != nil {
			log.Fatalf("ERROR failed to save contest metadata.")
		}
		fmt.Printf("Test %s is removed from task %s\n", testToken, taskToken)
	},
}

func init() {
	RootCmd.AddCommand(rmtestCmd)
}
