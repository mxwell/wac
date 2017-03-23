package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Test struct {
	Token  string
	Input  string
	Output string
}

type Task struct {
	Link  string
	Name  string
	Token string
}

type Contest struct {
	Link    string
	Name    string
	Tasks   map[string]Task
	RootDir string
}

type Platform interface {
	ValidUrl(url string) bool
	GetContest(url string, root_dirname string) (*Contest, error)
	GetTests(task *Task) ([]Test, error)
}

const root_file = ".contest.json"

func GetRootFile(contest *Contest) string {
	return filepath.Join(contest.RootDir, root_file)
}

func SaveContest(contest *Contest) error {
	b, err := json.Marshal(*contest)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(GetRootFile(contest), b, 0644)
}

func LoadContest(path string) (*Contest, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var contest Contest
	err = json.Unmarshal(b, &contest)
	return &contest, err
}

func LocateContest() (*Contest, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dir := wd
	for len(dir) > 0 {
		/* check if root file exists inside dir */
		path := filepath.Join(dir, root_file)
		if _, err := os.Stat(path); err == nil {
			return LoadContest(path)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return nil, fmt.Errorf("unable to locate contest metadata in the current directory %s or its parents", wd)
}
