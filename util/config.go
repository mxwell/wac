package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type BuildMethodRaw struct {
	Language string
	Command  string
}

type ExecMethod struct {
	Command string
}

type Configuration struct {
	TemplatesDir       string
	DefaultTemplate    string
	SolutionName       string
	Extensions         map[string]string
	BuildMethods       map[string]BuildMethodRaw
	DefaultBuildMethod string
	RunMethods         map[string]ExecMethod
	DefaultRunMethod   string
}

func GetDefaultLocation() string {
	if runtime.GOOS == "linux" {
		return filepath.Join(os.Getenv("HOME"), ".config", "wac")
	} else {
		panic(fmt.Errorf("OS %s is not supported", runtime.GOOS))
	}
}

func initConfiguration() *Configuration {
	var conf = &Configuration{
		TemplatesDir:    filepath.Join(GetDefaultLocation(), "templates"),
		DefaultTemplate: "gcc",
		SolutionName:    "main",
		Extensions: map[string]string{
			"c++11":   "cpp",
			"ocaml":   "ml",
			"python3": "py",
		},
		BuildMethods: map[string]BuildMethodRaw{
			"gcc": BuildMethodRaw{
				"c++11",
				"g++ --std=c++11 -pedantic -Wshadow -Wformat=2 -Wfloat-equal -Wconversion -Wlogical-op -fwhole-program -g -fsanitize=address -fstack-protector -Wall -Werror -Wextra  $INPUT -o $OUTPUT",
			},
			"gcc_fast": BuildMethodRaw{
				"c++11",
				"g++ --std=c++11 -O2 -Wall $INPUT -o $OUTPUT",
			},
			"ocaml": BuildMethodRaw{
				"ocaml",
				"ocamlopt $INPUT -o $OUTPUT",
			},
			"python3": BuildMethodRaw{
				"python3",
				"cp $INPUT $OUTPUT",
			},
		},
		DefaultBuildMethod: "gcc",
		RunMethods: map[string]ExecMethod{
			"gcc":     ExecMethod{"./$OUTPUT"},
			"ocaml":   ExecMethod{"./$OUTPUT"},
			"python3": ExecMethod{"python3 $OUTPUT"},
		},
		DefaultRunMethod: "gcc",
	}
	return conf
}

func createIfNotPresent(path string, bytes []byte) error {
	if PathExists(path) {
		return nil
	}
	return ioutil.WriteFile(path, bytes, 0600)
}

var GCC_TEMPLATE = `#include <bits/stdc++.h>

using namespace std;

int main() {
  // solution comes here
}`

var OCAML_TEMPLATE = `open Printf

let () =
  (* solution comes here *)`

var PY3_TEMPLATE = `import sys


def main():
  return 0


if __name__ == "__main__":
  sys.exit(main())`

func saveTemplates(dir string) error {
	err := createIfNotPresent(filepath.Join(dir, "gcc.cpp"), []byte(GCC_TEMPLATE))
	if err != nil {
		return err
	}
	err = createIfNotPresent(filepath.Join(dir, "ocaml.ml"), []byte(OCAML_TEMPLATE))
	if err != nil {
		return err
	}
	err = createIfNotPresent(filepath.Join(dir, "py3.py"), []byte(PY3_TEMPLATE))
	if err != nil {
		return err
	}
	// TODO add more here
	return err
}

func saveConfiguration(conf *Configuration, path string) error {
	b, err := json.Marshal(*conf)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0600)
}

func loadConfiguration(path string) (*Configuration, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf Configuration
	err = json.Unmarshal(b, &conf)
	return &conf, err
}

// Return configuration or panic
func CheckConfiguration() *Configuration {
	location := GetDefaultLocation()
	confPath := filepath.Join(location, "wac.json")
	templatesPath := filepath.Join(location, "templates")
	if PathExists(confPath) {
		conf, err := loadConfiguration(confPath)
		if err != nil {
			panic(fmt.Errorf("Failed to load configuration: %s", err))
		}
		return conf
	} else {
		conf := initConfiguration()
		err := os.MkdirAll(templatesPath, 0777)
		if err != nil {
			panic(fmt.Errorf("Failed to create dirs: %s", err))
		}
		err = saveTemplates(templatesPath)
		if err != nil {
			panic(fmt.Errorf("Failed to create templates: %s", err))
		}
		err = saveConfiguration(conf, confPath)
		if err != nil {
			panic(fmt.Errorf("Failed to create a new configuration: %s", err))
		}
		log.Printf("A new config is put into %s\n", location)
		return conf
	}
}
