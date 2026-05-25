package dotfiles

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func requireRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(filepath.Join(cwd, "dotfiles.toml")); err != nil {
		return "", fmt.Errorf("dotfiles.toml not found in current directory\nrun this command from the dotfiles repository root")
	}
	return cwd, nil
}

func expandTarget(target string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if target == "~" {
		return home, nil
	}
	if strings.HasPrefix(target, "~/") {
		return filepath.Join(home, target[2:]), nil
	}
	target = strings.ReplaceAll(target, "${HOME}", home)
	target = strings.ReplaceAll(target, "$HOME", home)
	if !filepath.IsAbs(target) {
		return "", fmt.Errorf("target must expand to an absolute path: %s", target)
	}
	return filepath.Clean(target), nil
}

func sourcePath(repo string, link Link) string {
	return filepath.Join(repo, filepath.Clean(link.Source))
}

func parseMode(value string) os.FileMode {
	parsed, err := strconv.ParseUint(value, 8, 32)
	if err != nil || parsed == 0 {
		return 0
	}
	return os.FileMode(parsed)
}
