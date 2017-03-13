package model

import (
	"encoding/json"
	"io/ioutil"
)

type Task struct {
	Link  string
	Name  string
	Token string
}

type Contest struct {
	Link  string
	Name  string
	Tasks map[string]Task
}

func SaveContest(contest *Contest, path string) error {
	b, err := json.Marshal(*contest)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, b, 0644)
	return err
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
