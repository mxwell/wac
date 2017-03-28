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
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mxwell/wac/model"
	"github.com/spf13/cobra"
)

/* Provided two abs path, the function returns their difference */
func getRelativePath(src string, dest string) string {
	prefix := ""
	for len(src) > 0 {
		if strings.HasPrefix(dest, src) {
			break
		}
		prefix = filepath.Join(prefix, "..")
		src = filepath.Dir(src)
	}
	branch := strings.TrimPrefix(dest, src)
	if len(branch) == 0 {
		branch = "."
	} else {
		branch = branch[1:]
	}
	return filepath.Join(prefix, branch)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show info about the working tree",
	Long:  `Show details of the contest associated with the working tree`,
	Run: func(cmd *cobra.Command, args []string) {
		contest, err := model.LocateContest()
		if err != nil {
			log.Fatalf("ERROR %s\n", err)
		}
		fmt.Printf("Contest: %s -- %s\n", contest.Name, contest.Link)
		if len(contest.Tasks) == 0 {
			fmt.Println("No tasks.")
		} else {
			wd, err := os.Getwd()
			if err != nil {
				log.Fatalf("ERROR %s\n", err)
			}
			fmt.Println("Tasks:")
			/* Order tokens lexicographically */
			tokens := make([]string, 0, len(contest.Tasks))
			for token, _ := range contest.Tasks {
				tokens = append(tokens, token)
			}
			sort.Strings(tokens)
			/* Process tasks ordered by token */
			for _, token := range tokens {
				task, _ := contest.Tasks[token]
				task_path := filepath.Join(contest.RootDir, token)
				rel_path := getRelativePath(wd, task_path)
				var token_s string
				if rel_path == "." {
					token_s = "* [" + token + "]"
				} else {
					token_s = "  [" + token + "]"
				}
				fmt.Printf("\n%s %s -- %s\n\tpath:  %s\n", token_s, task.Name, task.Link, rel_path)
				if len(task.TestTokens) > 0 {
					fmt.Printf("\ttests:")
					for _, testToken := range task.TestTokens {
						fmt.Printf(" %s", testToken)
					}
					fmt.Println()
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
