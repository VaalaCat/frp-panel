package upgrade

import (
	"fmt"
	"strings"
)

func buildDownloadURL(opt Options) (string, error) {
	if u := strings.TrimSpace(opt.DownloadURL); len(u) > 0 {
		return u, nil
	}
	version := strings.TrimSpace(opt.Version)
	if len(version) == 0 {
		version = "latest"
	}

	asset, err := detectAssetName()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://github.com/VaalaCat/frp-panel/releases/download/%s/%s", version, asset)
	if opt.UseGithubProxy && len(strings.TrimSpace(opt.GithubProxy)) > 0 {
		url = fmt.Sprintf("%s/%s", strings.TrimRight(strings.TrimSpace(opt.GithubProxy), "/"), url)
	}
	return url, nil
}


