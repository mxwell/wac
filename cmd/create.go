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
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Template struct {
	name      string
	extension string
}

func fullName(t Template) string {
	return t.name + "." + t.extension
}

func listTemplates(dirpath string) ([]Template, error) {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return nil, fmt.Errorf("unable to read templates directory '%s': %s", dirpath, err)
	}
	var result []Template
	for _, file := range files {
		name := file.Name()
		dot := strings.LastIndex(name, ".")
		if file.Mode().IsRegular() && dot > 0 {
			result = append(result, Template{name[:dot], name[dot+1:]})
		}
	}
	return result, nil
}

func findTemplate(name string) (Template, error) {
	templates, err := listTemplates(viper.GetString("TemplatesDir"))
	if err != nil {
		return Template{}, err
	}
	for _, template := range templates {
		if name == template.name {
			return template, nil
		}
	}
	return Template{}, fmt.Errorf("not found template '%s'", name)
}

func checkDestination(template Template, destination string) (string, error) {
	pattern := "." + template.extension
	if len(destination) < len(pattern) || destination[len(destination)-len(pattern):] != pattern {
		destination = destination + pattern
	}
	if _, err := os.Stat(destination); os.IsNotExist(err) {
		return destination, nil
	}
	return "", fmt.Errorf("File %s already exists. Please, remove to proceed.", destination)
}

func copyTemplate(template Template, destination string) error {
	source := viper.GetString("TemplatesDir") + "/" + fullName(template)

	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	cerr := output.Close()

	if err != nil {
		return err
	}

	return cerr
}

var Filename string

var createCmd = &cobra.Command{
	Use:   "create [template name]",
	Short: "Copy code template into current directory",
	Long:  `Copy a code template [template name into the current directory`,
	Run: func(cmd *cobra.Command, args []string) {
		template_name := viper.GetString("DefaultTemplate")
		if len(args) == 1 {
			template_name = args[0]
		}
		template, err := findTemplate(template_name)
		if err != nil {
			fmt.Printf("ERROR %s\n", err)
			return
		}
		destination, err := checkDestination(template, Filename)
		if err != nil {
			fmt.Printf("ERROR %s\n", err)
			return
		}
		if err := copyTemplate(template, destination); err != nil {
			fmt.Printf("ERROR when copying template '%s' into '%s': %s\n", template.name, destination, err)
			return
		}
		fmt.Printf("Created %s\n", destination)
	},
}

func init() {
	createCmd.Flags().StringVarP(&Filename, "filename", "f", "main", "Destination filename")
	RootCmd.AddCommand(createCmd)
}
