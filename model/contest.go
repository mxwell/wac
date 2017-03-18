package model

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

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
