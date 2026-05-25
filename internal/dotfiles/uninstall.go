package dotfiles

import (
	"fmt"
	"io"
	"os"
)

func Uninstall(repo string, links []Link, dryRun bool, out io.Writer) error {
	for _, status := range CheckStatuses(repo, links) {
		if status.Status != StatusOK {
			fmt.Fprintf(out, "skip    %s (%s)\n", status.Link.Target, status.Status)
			continue
		}
		fmt.Fprintf(out, "remove  %s\n", status.Link.Target)
		if dryRun {
			continue
		}
		target, err := expandTarget(status.Link.Target)
		if err != nil {
			return err
		}
		if err := os.Remove(target); err != nil {
			return err
		}
	}
	return nil
}
