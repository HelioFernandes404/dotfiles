package dotfiles

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/tabwriter"
)

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printHelp(stdout)
		return ExitOK
	}
	if len(args) > 1 && wantsCommandHelp(args[1:]) {
		printCommandHelp(stdout, args[0])
		return ExitOK
	}
	if runtime.GOOS != "linux" {
		fmt.Fprintln(stderr, "error: dotfiles CLI currently supports Linux only")
		return ExitGeneralError
	}

	repo, err := requireRoot()
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitUsageInvalid
	}
	manifest, err := LoadManifest(repo)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitManifestInvalid
	}

	switch args[0] {
	case "validate":
		return runValidate(repo, manifest, stdout, stderr)
	case "status":
		return runStatus(repo, manifest, args[1:], stdout, stderr)
	case "install":
		return runInstall(repo, manifest, args[1:], stdout, stderr)
	case "uninstall":
		return runUninstall(repo, manifest, args[1:], stdout, stderr)
	case "doctor":
		return runDoctor(repo, manifest, args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "error: unknown command %q\n", args[0])
		return ExitUsageInvalid
	}
}

func runValidate(repo string, manifest Manifest, stdout, stderr io.Writer) int {
	if err := ValidateManifest(repo, manifest); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitManifestInvalid
	}
	fmt.Fprintln(stdout, "manifest valid")
	return ExitOK
}

func runStatus(repo string, manifest Manifest, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.SetOutput(stderr)
	group := fs.String("group", "", "filter by comma-separated groups")
	if err := fs.Parse(args); err != nil {
		return ExitUsageInvalid
	}
	if err := ValidateManifest(repo, manifest); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitManifestInvalid
	}
	links := manifest.Links
	if *group != "" {
		selected, err := selectForStatus(manifest, *group)
		if err != nil {
			fmt.Fprintf(stderr, "error: %v\navailable groups: %s\n", err, strings.Join(Groups(manifest), ", "))
			return ExitUsageInvalid
		}
		links = selected
	}
	statuses := CheckStatuses(repo, links)
	printStatuses(stdout, statuses)
	if anyNotOK(statuses) {
		return ExitGeneralError
	}
	return ExitOK
}

func runInstall(repo string, manifest Manifest, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	fs.SetOutput(stderr)
	group := fs.String("group", "", "comma-separated groups to install")
	all := fs.Bool("all", false, "install all links")
	dryRun := fs.Bool("dry-run", false, "print actions without changing files")
	yes := fs.Bool("yes", false, "skip confirmation when backups are needed")
	if err := fs.Parse(args); err != nil {
		return ExitUsageInvalid
	}
	if err := ValidateManifest(repo, manifest); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitManifestInvalid
	}
	links, err := SelectLinks(manifest, *group, *all)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\navailable groups: %s\n", err, strings.Join(Groups(manifest), ", "))
		return ExitUsageInvalid
	}
	if err := Install(repo, links, InstallOptions{DryRun: *dryRun, Yes: *yes, Input: os.Stdin, Out: stdout}); err != nil {
		if strings.Contains(err.Error(), "cancelled") {
			return ExitGeneralError
		}
		fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitInstallConflict
	}
	return ExitOK
}

func runUninstall(repo string, manifest Manifest, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("uninstall", flag.ContinueOnError)
	fs.SetOutput(stderr)
	group := fs.String("group", "", "comma-separated groups to uninstall")
	all := fs.Bool("all", false, "uninstall all managed links")
	dryRun := fs.Bool("dry-run", false, "print actions without changing files")
	if err := fs.Parse(args); err != nil {
		return ExitUsageInvalid
	}
	if err := ValidateManifest(repo, manifest); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitManifestInvalid
	}
	links, err := SelectLinks(manifest, *group, *all)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\navailable groups: %s\n", err, strings.Join(Groups(manifest), ", "))
		return ExitUsageInvalid
	}
	if err := Uninstall(repo, links, *dryRun, stdout); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitGeneralError
	}
	return ExitOK
}

func runDoctor(repo string, manifest Manifest, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(stderr)
	verbose := fs.Bool("verbose", false, "print full link status")
	if err := fs.Parse(args); err != nil {
		return ExitUsageInvalid
	}
	hasError := false
	printCheck(stdout, "ok", "linux", "")
	printCheck(stdout, "ok", "dotfiles.toml found", "")
	if err := ValidateManifest(repo, manifest); err != nil {
		printCheck(stdout, "error", "manifest invalid", err.Error())
		hasError = true
	} else {
		printCheck(stdout, "ok", "manifest valid", "")
	}
	checkCommand(stdout, "go", "install Go to build the dotfiles CLI from source")
	checkCommand(stdout, "make", "install make to use repository shortcuts")
	checkLocalBin(stdout)
	checkSSH(stdout)
	statuses := CheckStatuses(repo, manifest.Links)
	if *verbose {
		printStatuses(stdout, statuses)
	} else {
		fmt.Fprintf(stdout, "links: %s\n", summarizeStatuses(statuses))
	}
	if hasError {
		return ExitGeneralError
	}
	return ExitOK
}

func selectForStatus(manifest Manifest, group string) ([]Link, error) {
	return SelectLinks(manifest, group, false)
}

func printStatuses(out io.Writer, statuses []LinkStatus) {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "STATUS\tGROUPS\tTARGET\tSOURCE\tDETAIL")
	for _, status := range statuses {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", status.Status, joinGroups(status.Link.Groups), status.Link.Target, status.Link.Source, status.Detail)
	}
	w.Flush()
}

func printCheck(out io.Writer, status, name, suggestion string) {
	fmt.Fprintf(out, "%s\t%s\n", status, name)
	if suggestion != "" {
		fmt.Fprintf(out, "\t%s\n", suggestion)
	}
}

func checkCommand(out io.Writer, name, suggestion string) {
	if _, err := exec.LookPath(name); err == nil {
		printCheck(out, "ok", name+" found", "")
		return
	}
	printCheck(out, "warn", name+" not found in PATH", suggestion)
}

func wantsCommandHelp(args []string) bool {
	return len(args) > 0 && (args[0] == "--help" || args[0] == "-h" || args[0] == "help")
}

func checkLocalBin(out io.Writer) {
	home, _ := os.UserHomeDir()
	localBin := home + "/.local/bin"
	for _, part := range strings.Split(os.Getenv("PATH"), ":") {
		if part == localBin {
			printCheck(out, "ok", "~/.local/bin in PATH", "")
			return
		}
	}
	printCheck(out, "warn", "~/.local/bin not in PATH", "add 'export PATH=\"$HOME/.local/bin:$PATH\"' to your shell config")
}

func checkSSH(out io.Writer) {
	home, _ := os.UserHomeDir()
	info, err := os.Stat(home + "/.ssh")
	if os.IsNotExist(err) {
		printCheck(out, "warn", "~/.ssh missing", "run dotfiles install --group ssh")
		return
	}
	if err != nil {
		printCheck(out, "warn", "~/.ssh not readable", err.Error())
		return
	}
	if info.Mode().Perm() != 0o700 {
		printCheck(out, "warn", "~/.ssh mode is not 0700", "run chmod 700 ~/.ssh")
		return
	}
	printCheck(out, "ok", "~/.ssh mode is 0700", "")
}

func summarizeStatuses(statuses []LinkStatus) string {
	counts := map[Status]int{}
	for _, status := range statuses {
		counts[status.Status]++
	}
	parts := []string{}
	for _, status := range []Status{StatusOK, StatusMissing, StatusConflict, StatusBroken, StatusInvalid} {
		if counts[status] > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", counts[status], status))
		}
	}
	return strings.Join(parts, ", ")
}

func printHelp(out io.Writer) {
	fmt.Fprintln(out, `dotfiles manages symlinks for this dotfiles repository.

Usage:
  dotfiles <command> [flags]

Commands:
  validate              validate dotfiles.toml
  status [--group g]    show link status
  install               install selected links
  uninstall             remove managed symlinks
  doctor                diagnose the CLI environment
  help                  show this help

Install and uninstall require --group <name[,name]> or --all.
Run commands from the repository root.`)
}

func printCommandHelp(out io.Writer, command string) {
	switch command {
	case "validate":
		fmt.Fprintln(out, "Usage: dotfiles validate")
	case "status":
		fmt.Fprintln(out, "Usage: dotfiles status [--group name[,name]]")
	case "install":
		fmt.Fprintln(out, "Usage: dotfiles install (--group name[,name] | --all) [--dry-run] [--yes]")
	case "uninstall":
		fmt.Fprintln(out, "Usage: dotfiles uninstall (--group name[,name] | --all) [--dry-run]")
	case "doctor":
		fmt.Fprintln(out, "Usage: dotfiles doctor [--verbose]")
	default:
		fmt.Fprintf(out, "Unknown command %q\n", command)
	}
}
