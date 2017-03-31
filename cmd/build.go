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
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Language struct {
	name      string
	extension string
}

type BuildMethod struct {
	language *Language
	command  string
}

var InputName string
var OutputName string
var LanguageByName = map[string]*Language{}
var MethodByName = map[string]*BuildMethod{}

func readBuildMethod(name string) *BuildMethod {
	subtree := viper.Sub("BuildMethods." + name)
	langName := subtree.GetString("Language")
	language, ok := LanguageByName[langName]
	if !ok {
		panic(fmt.Errorf("Bad config: build method '%s' uses unknown language '%s'", name, langName))
	}
	return &BuildMethod{language, subtree.GetString("Command")}
}

func readConfig() {
	languages := viper.GetStringMapString("Extensions")
	for name, ext := range languages {
		LanguageByName[name] = &Language{name, ext}
	}
	methods := viper.GetStringMap("BuildMethods")
	for name, _ := range methods {
		MethodByName[name] = readBuildMethod(name)
	}
}

func getInput(method *BuildMethod) (string, error) {
	name := InputName
	if strings.HasSuffix(name, ".*") {
		name = name[:len(name)-1] + method.language.extension
	}
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return name, fmt.Errorf("File '%s' does not exist.", name)
	} else {
		return name, nil
	}
}

func getCommand(method *BuildMethod, input, output string) *exec.Cmd {
	if len(method.command) == 0 {
		return nil
	}
	tokens := strings.Split(method.command, " ")
	for i, _ := range tokens {
		if strings.Contains(tokens[i], "$INPUT") {
			tokens[i] = strings.Replace(tokens[i], "$INPUT", input, -1)
		}
		if strings.Contains(tokens[i], "$OUTPUT") {
			tokens[i] = strings.Replace(tokens[i], "$OUTPUT", output, -1)
		}
	}
	return exec.Command(tokens[0], tokens[1:]...)
}

var buildCmd = &cobra.Command{
	Use:   "build [build method]",
	Short: "Build solution",
	Long:  `Build solution from a source file into an executable using [build method], if applicable`,
	Run: func(cmd *cobra.Command, args []string) {
		readConfig()
		methodName := viper.GetString("DefaultBuildMethod")
		if len(args) == 1 {
			methodName = args[0]
		}
		method, ok := MethodByName[methodName]
		if !ok {
			fmt.Printf("ERROR build method '%s' not found in config\n", methodName)
			return
		}
		if len(OutputName) == 0 {
			OutputName = viper.GetString("SolutionName")
		}
		input, err := getInput(method)
		if err != nil {
			fmt.Printf("ERROR bad input: %s", err)
			return
		}
		if input == OutputName {
			fmt.Printf("ERROR equal input and output - '%s'", input)
			return
		}
		command := getCommand(method, input, OutputName)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		if command != nil {
			command.Stdout = &stdout
			command.Stderr = &stderr
			err = command.Run()
		}
		if err != nil {
			fmt.Printf("ERROR Build failed: %s\n", err)
		} else {
			fmt.Println("OK")
		}
		if stdout.Len() > 0 {
			fmt.Println("<stdout>")
			stdout.WriteTo(os.Stdout)
		}
		if stderr.Len() > 0 {
			fmt.Println("<stderr>")
			stderr.WriteTo(os.Stdout)
		}
	},
}

func init() {
	buildCmd.Flags().StringVarP(&InputName, "input", "i", "main.*", "Build input file")
	buildCmd.Flags().StringVarP(&OutputName, "output", "o", "", "Build output file (default is set by config as SolutionName)")
	RootCmd.AddCommand(buildCmd)
}
