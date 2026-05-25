package dotfiles

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateManifestRejectsDuplicateTargets(t *testing.T) {
	repo := t.TempDir()
	mustWrite(t, filepath.Join(repo, "a"), "a")
	mustWrite(t, filepath.Join(repo, "b"), "b")
	manifest := Manifest{Links: []Link{
		{Source: "a", Target: "~/.x", Groups: []string{"one"}},
		{Source: "b", Target: "~/.x", Groups: []string{"two"}},
	}}
	if err := ValidateManifest(repo, manifest); err == nil {
		t.Fatal("expected duplicate target error")
	}
}

func TestSelectLinksUsesOrForCommaSeparatedGroups(t *testing.T) {
	manifest := Manifest{Links: []Link{
		{Source: "git/gitconfig", Target: "~/.gitconfig", Groups: []string{"core", "git"}},
		{Source: "nvim", Target: "~/.config/nvim", Groups: []string{"editor", "nvim"}},
		{Source: "yazi", Target: "~/.config/yazi", Groups: []string{"terminal", "yazi"}},
	}}
	links, err := SelectLinks(manifest, "core,editor", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(links))
	}
}

func TestSelectLinksRejectsUnknownGroup(t *testing.T) {
	manifest := Manifest{Links: []Link{{Source: "git/gitconfig", Target: "~/.gitconfig", Groups: []string{"core"}}}}
	if _, err := SelectLinks(manifest, "missing", false); err == nil {
		t.Fatal("expected unknown group error")
	}
}

func TestCheckStatusClassifiesSymlinks(t *testing.T) {
	repo := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)
	mustWrite(t, filepath.Join(repo, "gitconfig"), "[user]\n")
	link := Link{Source: "gitconfig", Target: "$HOME/.gitconfig", Groups: []string{"core"}}

	if got := CheckStatus(repo, link).Status; got != StatusMissing {
		t.Fatalf("expected missing, got %s", got)
	}
	if err := os.Symlink(filepath.Join(repo, "gitconfig"), filepath.Join(home, ".gitconfig")); err != nil {
		t.Fatal(err)
	}
	if got := CheckStatus(repo, link).Status; got != StatusOK {
		t.Fatalf("expected ok, got %s", got)
	}
	if err := os.Remove(filepath.Join(home, ".gitconfig")); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(home, ".gitconfig"), "local")
	if got := CheckStatus(repo, link).Status; got != StatusConflict {
		t.Fatalf("expected conflict, got %s", got)
	}
}

func TestInstallAndUninstallUseTempHome(t *testing.T) {
	repo := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)
	mustWrite(t, filepath.Join(repo, "nvim"), "not a dir but linkable")
	link := Link{Source: "nvim", Target: "~/.config/nvim", Groups: []string{"editor"}}

	if err := Install(repo, []Link{link}, InstallOptions{Yes: true, Input: nil, Out: discard{}}); err != nil {
		t.Fatal(err)
	}
	if got := CheckStatus(repo, link).Status; got != StatusOK {
		t.Fatalf("expected ok after install, got %s", got)
	}
	if err := Uninstall(repo, []Link{link}, false, discard{}); err != nil {
		t.Fatal(err)
	}
	if got := CheckStatus(repo, link).Status; got != StatusMissing {
		t.Fatalf("expected missing after uninstall, got %s", got)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

type discard struct{}

func (discard) Write(p []byte) (int, error) { return len(p), nil }
