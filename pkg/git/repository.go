package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetGitDir(repoDir string) (string, error) {
	gitPath := filepath.Join(repoDir, ".git")

	info, err := os.Stat(gitPath)
	if err != nil {
		return "", fmt.Errorf("stat .git: %w", err)
	}

	if info.IsDir() {
		return gitPath, nil
	}

	data, err := os.ReadFile(gitPath)
	if err != nil {
		return "", fmt.Errorf("read .git file: %w", err)
	}

	const prefix = "gitdir: "

	line := strings.TrimSpace(string(data))
	if !strings.HasPrefix(line, prefix) {
		return "", fmt.Errorf("invalid .git file")
	}

	gitDir := strings.TrimSpace(strings.TrimPrefix(line, prefix))

	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(repoDir, gitDir)
	}

	return filepath.Clean(gitDir), nil
}

func ParseOriginURL(config string) (string, error) {
	lines := strings.Split(config, "\n")

	inOrigin := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			inOrigin = line == `[remote "origin"]`
			continue
		}

		if !inOrigin {
			continue
		}

		if !strings.HasPrefix(line, "url") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		url := strings.TrimSpace(parts[1])
		if url == "" {
			return "", fmt.Errorf("origin remote has an empty URL")
		}

		return url, nil
	}

	return "", fmt.Errorf("origin remote not found")
}
