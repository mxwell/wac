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

	"github.com/PuerkitoBio/goquery"
	"github.com/mxwell/wac/model"
	"github.com/spf13/cobra"
)

func readFile(path string) (*goquery.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %s", err)
	}
	return goquery.NewDocumentFromReader(f)
}

var fetchCmd = &cobra.Command{
	Use:   "fetch [arg]",
	Short: "TODO",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("1 arg is required")
			return
		}
		doc, err := readFile(args[0])
		if err != nil {
			fmt.Printf("ERROR %s\n", err)
			return
		}
		tasks := make(map[string]model.Task)
		doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
			link := s.Find("td.center a.linkwrapper")
			href, ok := link.Attr("href")
			if ok {
				fmt.Printf("%s - %s\n", link.Text(), href)
				tasks[link.Text()] = model.Task{href, link.Text(), link.Text()}
			}
		})
		contest := model.Contest{"http://example.com", "Example Contest", tasks}
		err = model.SaveContest(&contest, "bb.json")
		if err != nil {
			fmt.Printf("ERROR %s\n", err)
			return
		}
		c, err := model.LoadContest("aa.json")
		if err != nil {
			fmt.Printf("ERROR %s\n", err)
			return
		}
		err = model.SaveContest(c, "cc.json")
		if err != nil {
			fmt.Printf("ERROR %s\n", err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(fetchCmd)
}
