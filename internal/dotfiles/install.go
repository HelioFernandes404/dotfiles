package dotfiles

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type InstallOptions struct {
	DryRun bool
	Yes    bool
	Input  io.Reader
	Out    io.Writer
}

func Install(repo string, links []Link, opts InstallOptions) error {
	statuses := CheckStatuses(repo, links)
	var conflicts []LinkStatus
	for _, status := range statuses {
		if status.Status == StatusConflict || status.Status == StatusBroken {
			conflicts = append(conflicts, status)
		}
		if status.Status == StatusInvalid {
			return fmt.Errorf("invalid link %s: %s", status.Link.Source, status.Detail)
		}
	}

	backupDir := filepath.Join(repo, ".dotfiles-backups", time.Now().Format("2006-01-02_150405"))
	printInstallPlan(repo, statuses, backupDir, opts.Out)
	if opts.DryRun {
		return nil
	}
	if len(conflicts) > 0 && !opts.Yes {
		fmt.Fprint(opts.Out, "Continue? [y/N] ")
		reader := bufio.NewReader(opts.Input)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			return fmt.Errorf("cancelled")
		}
	}

	for _, status := range statuses {
		if status.Status == StatusOK {
			continue
		}
		if err := installOne(repo, status.Link, status, backupDir); err != nil {
			return err
		}
	}
	return nil
}

func installOne(repo string, link Link, status LinkStatus, backupDir string) error {
	source := sourcePath(repo, link)
	target, err := expandTarget(link.Target)
	if err != nil {
		return err
	}
	parent := filepath.Dir(target)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	if link.ParentMode != "" {
		if err := os.Chmod(parent, parseMode(link.ParentMode)); err != nil {
			return err
		}
	}
	if link.Mode != "" {
		if err := os.Chmod(source, parseMode(link.Mode)); err != nil {
			return err
		}
	}
	if status.Status == StatusConflict || status.Status == StatusBroken {
		backupTarget := filepath.Join(backupDir, strings.TrimPrefix(target, string(filepath.Separator)))
		if err := os.MkdirAll(filepath.Dir(backupTarget), 0o755); err != nil {
			return err
		}
		if err := os.Rename(target, backupTarget); err != nil {
			return err
		}
	}
	return os.Symlink(source, target)
}

func printInstallPlan(repo string, statuses []LinkStatus, backupDir string, out io.Writer) {
	for _, status := range statuses {
		source := sourcePath(repo, status.Link)
		switch status.Status {
		case StatusOK:
			fmt.Fprintf(out, "ok      %s -> %s\n", status.Link.Target, source)
		case StatusMissing:
			fmt.Fprintf(out, "create  %s -> %s\n", status.Link.Target, source)
		case StatusConflict, StatusBroken:
			target, _ := expandTarget(status.Link.Target)
			backupTarget := filepath.Join(backupDir, strings.TrimPrefix(target, string(filepath.Separator)))
			fmt.Fprintf(out, "backup  %s -> %s\n", status.Link.Target, backupTarget)
			fmt.Fprintf(out, "create  %s -> %s\n", status.Link.Target, source)
		}
	}
}
