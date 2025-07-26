package utils

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

func FetchFilesFromGitHub(owner, repo, branch, folderPath, targetDir string) error {
    os.MkdirAll(targetDir, os.ModePerm)

    apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, folderPath, branch)
    resp, err := http.Get(apiURL)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var files []struct {
        Name string `json:"name"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
        return err
    }

    for _, file := range files {
        rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s/%s", owner, repo, branch, folderPath, file.Name)
        outPath := fmt.Sprintf("%s/%s", targetDir, file.Name)
        fmt.Printf("Downloading %s ...\n", file.Name)

        out, err := os.Create(outPath)
        if err != nil {
            return err
        }

        res, err := http.Get(rawURL)
        if err != nil {
            return err
        }
        defer res.Body.Close()

        _, err = io.Copy(out, res.Body)
        out.Close()
        if err != nil {
            return err
        }
    }

    fmt.Println("Download complete.")
    return nil
}
