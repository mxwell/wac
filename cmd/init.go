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

	"github.com/mxwell/wac/model"
	"github.com/mxwell/wac/platforms/atcoder"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [url]",
	Short: "Initialize contest in directory",
	Long:  `Initialize current or specified directory with metadata of a contest at [url]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Printf("at least 1 arg is required")
			return
		}
		/*
		 * TODO there should be a filter collection, to which platform handlers should register
		 * Handlers receives an URL one after another. If a handler is able to process the URL,
		 * it processes it and returns *Contest. Otherwise, it returns nil and the next
		 * handler proceeds.
		 */
		contest, err := atcoder.FetchContest(args[0])
		if err != nil {
			fmt.Printf("ERROR can't fetch contest: %s\n", err)
			return
		}
		err = model.SaveContest(contest, "contest.json")
		if err != nil {
			fmt.Printf("ERROR can't save the contest: %s\n", err)
			return
		}
		fmt.Println("OK")
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
