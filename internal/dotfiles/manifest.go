package dotfiles

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

func LoadManifest(repo string) (Manifest, error) {
	data, err := os.ReadFile(filepath.Join(repo, "dotfiles.toml"))
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	if err := toml.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

func ValidateManifest(repo string, manifest Manifest) error {
	if len(manifest.Links) == 0 {
		return errors.New("manifest has no links")
	}

	targets := map[string]string{}
	for i, link := range manifest.Links {
		prefix := fmt.Sprintf("links[%d]", i)
		if link.Source == "" {
			return fmt.Errorf("%s source is required", prefix)
		}
		if link.Target == "" {
			return fmt.Errorf("%s target is required", prefix)
		}
		if len(link.Groups) == 0 {
			return fmt.Errorf("%s groups is required", prefix)
		}
		if filepath.IsAbs(link.Source) || strings.HasPrefix(link.Source, "..") || strings.Contains(link.Source, string(filepath.Separator)+".."+string(filepath.Separator)) {
			return fmt.Errorf("%s source must stay inside the repository", prefix)
		}
		sourcePath := filepath.Join(repo, filepath.Clean(link.Source))
		if _, err := os.Lstat(sourcePath); err != nil {
			return fmt.Errorf("%s source does not exist: %s", prefix, link.Source)
		}
		if prev, exists := targets[link.Target]; exists {
			return fmt.Errorf("duplicate target %s: %s and %s", link.Target, prev, link.Source)
		}
		targets[link.Target] = link.Source
		if link.Mode != "" && parseMode(link.Mode) == 0 {
			return fmt.Errorf("%s mode must be an octal file mode", prefix)
		}
		if link.ParentMode != "" && parseMode(link.ParentMode) == 0 {
			return fmt.Errorf("%s parent_mode must be an octal file mode", prefix)
		}
	}
	return nil
}

func Groups(manifest Manifest) []string {
	seen := map[string]bool{}
	var groups []string
	for _, link := range manifest.Links {
		for _, group := range link.Groups {
			if !seen[group] {
				seen[group] = true
				groups = append(groups, group)
			}
		}
	}
	return groups
}

func SelectLinks(manifest Manifest, groupCSV string, all bool) ([]Link, error) {
	if all && groupCSV != "" {
		return nil, errors.New("use either --all or --group, not both")
	}
	if all {
		return manifest.Links, nil
	}
	if groupCSV == "" {
		return nil, errors.New("--group or --all is required")
	}

	available := map[string]bool{}
	for _, group := range Groups(manifest) {
		available[group] = true
	}
	wanted := map[string]bool{}
	for _, group := range strings.Split(groupCSV, ",") {
		group = strings.TrimSpace(group)
		if group == "" {
			continue
		}
		if !available[group] {
			return nil, fmt.Errorf("unknown group %q", group)
		}
		wanted[group] = true
	}
	if len(wanted) == 0 {
		return nil, errors.New("--group cannot be empty")
	}

	var selected []Link
	for _, link := range manifest.Links {
		for _, group := range link.Groups {
			if wanted[group] {
				selected = append(selected, link)
				break
			}
		}
	}
	return selected, nil
}
