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

	"github.com/mxwell/wac/model"
	"github.com/mxwell/wac/platforms"
	"github.com/spf13/cobra"
)

var fetchAll bool

func saveStringToFile(s *string, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("can't create file '%s': %s", path, err)
	}
	defer f.Close()
	_, err = f.WriteString(*s)
	if err != nil {
		return fmt.Errorf("error while writing to file '%s': %s", path, err)
	}
	return nil
}

func contains(arr *[]string, value string) bool {
	for _, e := range *arr {
		if e == value {
			return true
		}
	}
	return false
}

// The function fetches samples from a platform, saves them into task directory
// and adds info on the samples into contest struct
func fetchForTask(platform model.Platform, contest *model.Contest, token string) error {
	task_path := filepath.Join(contest.RootDir, token)
	if _, err := os.Stat(task_path); os.IsNotExist(err) {
		err = os.MkdirAll(task_path, 0777)
		if err != nil {
			return fmt.Errorf("can't create a subdir '%s' for task: %s", task_path, err)
		}
	}
	log.Printf("Processing task '%s' ...\n", token)
	task, ok := contest.Tasks[token]
	if !ok {
		return fmt.Errorf("no task with token '%s' in contest '%s'", token, contest.Name)
	}
	tests, err := platform.GetTests(&task)
	if err != nil {
		return fmt.Errorf("unable to get tests for task with token '%s': %s", token, err)
	}
	for _, test := range tests {
		sample_path := filepath.Join(task_path, test.Token)
		if contains(&task.TestTokens, test.Token) {
			log.Printf("Test '%s' was already present, re-writing...", test.Token)
		} else {
			task.TestTokens = append(task.TestTokens, test.Token)
		}
		input_path := sample_path + ".in"
		output_path := sample_path + ".out"
		err = saveStringToFile(&test.Input, input_path)
		if err != nil {
			return err
		}
		err = saveStringToFile(&test.Output, output_path)
		if err != nil {
			return err
		}
		log.Printf("%s saved to %s and %s\n", test.Token, filepath.Base(input_path), filepath.Base(output_path))
	}
	contest.Tasks[token] = task
	return nil
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch tests for task(s)",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		contest, err := model.LocateContest()
		if err != nil {
			log.Fatalf("ERROR %s\n", err)
		}

		platform := platforms.FindPlatform(contest.Link)
		if platform == nil {
			log.Fatalf("ERROR unable to find platform for contest url %s\n", contest.Link)
		}

		if fetchAll {
			/* Order tokens lexicographically */
			tokens := make([]string, 0, len(contest.Tasks))
			for token, _ := range contest.Tasks {
				tokens = append(tokens, token)
			}
			sort.Strings(tokens)
			/* Process tasks ordered by token */
			for _, token := range tokens {
				err := fetchForTask(platform, contest, token)
				if err != nil {
					log.Printf("ERROR can't fetch task: %s\n", err)
				}
			}
		} else {
			token, err := model.DetermineCurrentTask(contest)
			if err != nil {
				log.Fatalf("ERROR can't determine current task: %s\n", err)
			}
			err = fetchForTask(platform, contest, token)
			if err != nil {
				log.Fatalf("ERROR can't fetch task: %s\n", err)
			}
		}

		err = model.SaveContest(contest)
		if err != nil {
			log.Fatalf("ERROR failed to save contest metadata.")
		}
	},
}

func init() {
	fetchCmd.Flags().BoolVarP(&fetchAll, "all", "a", false, "Fetch tests for all tasks")
	RootCmd.AddCommand(fetchCmd)
}
