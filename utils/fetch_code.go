package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/IshaanNene/EliteCode-brew/problems"
)

const (
	GHOwner = "IshaanNene"
	GHRepo  = "AlgoRank"
	GHBranch = "main"
)

func FetchStarterCode(p problems.SelectedProblem) (string, error) {
	dir := filepath.Join(".", p.ID+"_"+strings.ToLower(p.Language))
	os.MkdirAll(dir, 0755)

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/Solutions/%s?ref=%s", GHOwner, GHRepo, p.ID, GHBranch)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch GitHub files list: %v", err)
	}
	defer resp.Body.Close()

	var files []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}
	json.NewDecoder(resp.Body).Decode(&files)

	for _, f := range files {
		if strings.Contains(f.Name, p.Language) || strings.Contains(f.Name, "testcases") {
			rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", GHOwner, GHRepo, GHBranch, f.Path)
			destPath := filepath.Join(dir, f.Name)

			fileResp, err := http.Get(rawURL)
			if err != nil {
				return "", fmt.Errorf("failed to fetch %s: %v", f.Name, err)
			}
			defer fileResp.Body.Close()

			out, _ := os.Create(destPath)
			io.Copy(out, fileResp.Body)
			out.Close()
		}
	}
	return dir, nil
}
