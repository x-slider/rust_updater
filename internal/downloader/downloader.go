package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetGitHubReleaseAssetURL(api, name string) (string, error) {
	res, err := http.Get(api)
	if err != nil {
		return "", fmt.Errorf("GitHub API error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API status: %s", res.Status)
	}

	var data struct {
		Assets []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		}
	}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("decode error: %w", err)
	}

	for _, a := range data.Assets {
		if strings.EqualFold(a.Name, name) {
			return a.URL, nil
		}
	}
	return "", fmt.Errorf("asset not found: %s", name)
}

func DownloadFile(url, path string) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	return err
}
