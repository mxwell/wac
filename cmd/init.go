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
	"os"
	"path/filepath"

	"github.com/mxwell/wac/model"
	"github.com/mxwell/wac/platforms/atcoder"
	"github.com/spf13/cobra"
)

func determineRootDirectory(args []string) (string, error) {
	if len(args) == 2 {
		if root_dirname, err := filepath.Abs(args[1]); err != nil {
			return "", err
		} else {
			return root_dirname, nil
		}
	}
	if root_dirname, err := os.Getwd(); err != nil {
		return "", err
	} else {
		return root_dirname, nil
	}
}

func initContestDirectory(contest *model.Contest) error {
	/* Compose path for metadata and check it's not taken */
	root_file := model.GetRootFile(contest)
	if _, err := os.Stat(root_file); err == nil {
		return fmt.Errorf("File %s already exists -- please, remove to initialize the directory for the contest\n", root_file)
	}
	/* Create root directory if it doesn't exist */
	if _, err := os.Stat(contest.RootDir); os.IsNotExist(err) {
		err = os.MkdirAll(contest.RootDir, 0777)
		if err != nil {
			return err
		}
	}
	/* Create root file */
	if err := model.SaveContest(contest); err != nil {
		return err
	}
	/* Create subdirs for tasks */
	for token, _ := range contest.Tasks {
		path := filepath.Join(contest.RootDir, token)
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
	}
	return nil
}

var initCmd = &cobra.Command{
	Use:   "init URL [DIRECTORY]",
	Short: "Initialize contest in directory",
	Long: `Initialize DIRECTORY with metadata of a contest at URL.
Current directory if DIRECTORY is not specified.
Directory is created if it doesn't exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 || len(args) > 2 {
			fmt.Printf("wrong number of arguments - %d\n", len(args))
			return
		}

		root_dirname, err := determineRootDirectory(args)
		if err != nil {
			fmt.Printf("ERROR can't determine working directory: %s\n", err)
			return
		}

		fmt.Println(root_dirname)

		/*
		 * TODO there should be a filter collection, to which platform handlers should register
		 * Handlers receives an URL one after another. If a handler is able to process the URL,
		 * it processes it and returns *Contest. Otherwise, it returns nil and the next
		 * handler proceeds.
		 */
		contest, err := atcoder.FetchContest(args[0], root_dirname)
		if err != nil {
			fmt.Printf("ERROR can't fetch contest: %s\n", err)
			return
		}

		if err := initContestDirectory(contest); err != nil {
			fmt.Printf("ERROR can't init contest directory: %s\n", err)
		}

		fmt.Printf("Root directory: %s\n", contest.RootDir)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
