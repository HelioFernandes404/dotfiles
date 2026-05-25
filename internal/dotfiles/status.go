package dotfiles

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CheckStatuses(repo string, links []Link) []LinkStatus {
	statuses := make([]LinkStatus, 0, len(links))
	for _, link := range links {
		statuses = append(statuses, CheckStatus(repo, link))
	}
	return statuses
}

func CheckStatus(repo string, link Link) LinkStatus {
	source := sourcePath(repo, link)
	if _, err := os.Lstat(source); err != nil {
		return LinkStatus{Link: link, Status: StatusInvalid, Detail: err.Error()}
	}
	target, err := expandTarget(link.Target)
	if err != nil {
		return LinkStatus{Link: link, Status: StatusInvalid, Detail: err.Error()}
	}
	info, err := os.Lstat(target)
	if os.IsNotExist(err) {
		return LinkStatus{Link: link, Status: StatusMissing}
	}
	if err != nil {
		return LinkStatus{Link: link, Status: StatusInvalid, Detail: err.Error()}
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return LinkStatus{Link: link, Status: StatusConflict, Detail: "target exists and is not a symlink"}
	}
	dest, err := os.Readlink(target)
	if err != nil {
		return LinkStatus{Link: link, Status: StatusInvalid, Detail: err.Error()}
	}
	if !filepath.IsAbs(dest) {
		dest = filepath.Join(filepath.Dir(target), dest)
	}
	dest = filepath.Clean(dest)
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		return LinkStatus{Link: link, Status: StatusBroken, Detail: dest}
	}
	if samePath(dest, source) {
		return LinkStatus{Link: link, Status: StatusOK}
	}
	return LinkStatus{Link: link, Status: StatusConflict, Detail: fmt.Sprintf("points to %s", dest)}
}

func samePath(a, b string) bool {
	absA, errA := filepath.Abs(a)
	absB, errB := filepath.Abs(b)
	if errA == nil {
		a = absA
	}
	if errB == nil {
		b = absB
	}
	return filepath.Clean(a) == filepath.Clean(b)
}

func anyNotOK(statuses []LinkStatus) bool {
	for _, status := range statuses {
		if status.Status != StatusOK {
			return true
		}
	}
	return false
}

func joinGroups(groups []string) string {
	return strings.Join(groups, ",")
}
