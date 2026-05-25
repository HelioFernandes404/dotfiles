package dotfiles

const (
	ExitOK              = 0
	ExitGeneralError    = 1
	ExitManifestInvalid = 2
	ExitInstallConflict = 3
	ExitUsageInvalid    = 4
)

type Manifest struct {
	Links []Link `toml:"links"`
}

type Link struct {
	Source     string   `toml:"source"`
	Target     string   `toml:"target"`
	Groups     []string `toml:"groups"`
	Mode       string   `toml:"mode"`
	ParentMode string   `toml:"parent_mode"`
}

type Status string

const (
	StatusOK       Status = "ok"
	StatusMissing  Status = "missing"
	StatusConflict Status = "conflict"
	StatusBroken   Status = "broken"
	StatusInvalid  Status = "invalid"
)

type LinkStatus struct {
	Link   Link
	Status Status
	Detail string
}
