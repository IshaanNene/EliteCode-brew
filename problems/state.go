package problems

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getSelectedStatePath() string {
	return filepath.Join(os.Getenv("HOME"), ".elitecode", "selected.json")
}

func GetSelectedProblemAndLang() (Problem, string, error) {
	var p Problem
	var lang string
	path := getSelectedStatePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return p, lang, err
	}

	var stored struct {
		Problem Problem `json:"problem"`
		Lang    string  `json:"lang"`
	}

	err = json.Unmarshal(data, &stored)
	if err != nil {
		return p, lang, err
	}

	return stored.Problem, stored.Lang, nil
}

func GetProblemDirectory() string {
	p, _, err := GetSelectedProblemAndLang()
	if err != nil {
		return ""
	}
	return p.ID
}

func GetSelectedLanguage() string {
	_, lang, err := GetSelectedProblemAndLang()
	if err != nil {
		return ""
	}
	return lang
}
