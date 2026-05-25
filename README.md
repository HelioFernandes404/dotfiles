# Dotfiles

Personal Linux dotfiles managed by a small Go CLI.

## Contents

- `git/gitconfig` -> `~/.gitconfig`
- `ssh/config` -> `~/.ssh/config`
- `nvim` -> `~/.config/nvim`
- `yazi` -> `~/.config/yazi`
- `dotfiles.toml` declares managed symlinks and groups
- `cmd/dotfiles` contains the CLI entrypoint
- `internal/dotfiles` contains the implementation

## Bootstrap

Clone the repository and install the CLI binary:

```bash
git clone https://github.com/HelioFernandes404/dotfiles.git ~/dotfiles
cd ~/dotfiles
make bootstrap
```

`make bootstrap` builds the CLI and copies it to `~/.local/bin/dotfiles`. It does not install any symlinks.

## Usage

Run commands from the repository root. The CLI uses `dotfiles.toml` as the root marker and manifest.

```bash
dotfiles validate
dotfiles status
dotfiles install --group core
dotfiles install --group editor
dotfiles install --group terminal
dotfiles install --all
dotfiles uninstall --group core
dotfiles doctor
```

The Makefile uses the freshly built `./bin/dotfiles` binary for development shortcuts:

```bash
make validate
make status
make install-core
make install-editor
make install-terminal
make uninstall-core
make test
```

## Manifest

Links are explicit in `dotfiles.toml`:

```toml
[[links]]
source = "git/gitconfig"
target = "~/.gitconfig"
groups = ["core", "git"]
```

Supported fields:

- `source`: path inside this repository.
- `target`: absolute target path after expanding `~`, `$HOME`, or `${HOME}`.
- `groups`: one or more groups used for filtering.
- `mode`: optional octal mode applied to the source file.
- `parent_mode`: optional octal mode applied to the target parent directory.

Initial groups:

- `core`: Git and SSH.
- `editor`: Neovim.
- `terminal`: Yazi.
- Specific groups also exist: `git`, `ssh`, `nvim`, `yazi`.

Multiple groups use OR semantics:

```bash
dotfiles install --group core,editor
```

## Behavior Contract

- Linux only for now.
- `install` and `uninstall` require `--group` or `--all`.
- `status` shows all links by default and accepts `--group` as a filter.
- `validate` always validates the full manifest.
- Duplicate targets are invalid.
- Missing parent directories are created idempotently.
- Symlinks use absolute source paths.
- Existing correct symlinks are left unchanged.
- Existing files, directories, broken symlinks, or symlinks pointing elsewhere are conflicts.
- Conflicts are backed up inside `.dotfiles-backups/<timestamp>/...` before replacement.
- If a backup is needed, `install` prints the plan and asks for confirmation.
- `--yes` skips backup confirmation.
- `--dry-run` prints planned actions without changing files.
- `uninstall` removes only symlinks that point to the expected repository source.
- `uninstall` never removes parent directories.
- The CLI stores no state. The manifest and filesystem are the source of truth.

## Doctor

`dotfiles doctor` is read-only. It checks:

- Linux support.
- Manifest presence and validity.
- `go` and `make` availability as warnings.
- Whether `~/.local/bin` is in `PATH`.
- SSH directory permissions when `~/.ssh` exists.
- Link status summary.

Use verbose mode for full link status:

```bash
dotfiles doctor --verbose
```

Warnings do not make `doctor` fail. Errors do.

## Exit Codes

- `0`: success.
- `1`: general error, failed status, cancelled confirmation, or doctor error.
- `2`: invalid manifest.
- `3`: installation conflict or install failure.
- `4`: invalid usage.

## Design Rationale

- Go was chosen to produce a single binary with no runtime dependency.
- TOML was chosen because this repo already uses TOML and it is easy to edit by hand.
- The manifest is explicit instead of inferred from folders to avoid unsafe path guessing.
- Groups allow both broad installs, such as `core`, and specific installs, such as `nvim`.
- No local state is stored because symlink targets can be derived from the manifest and filesystem.
- Backups live inside the repo under `.dotfiles-backups/` but are ignored by Git.
- Make is only a convenience layer; the CLI is the source of truth.

## Security

Never commit secrets, tokens, private SSH keys, `known_hosts`, or machine-specific credentials. The repository intentionally tracks only `ssh/config`, not SSH keys.
