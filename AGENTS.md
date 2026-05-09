# AGENTS.md

- Treat `install.sh` and `uninstall.sh` as the source of truth for symlinks; update both together when adding or removing a dotfile.
- `install.sh` backs up replaced files to `~/dotfiles_backup_*` and hardens SSH permissions with `700` on `~/.ssh` and `600` on `~/.ssh/config`.
- The root `README.md` tree is stale; trust the tracked filesystem and scripts over the docs when they disagree.
- Current top-level config lives in `config/`, `git/`, `ssh/`, `yazi/`, `claude/`, and `nvim/`.
- `config/bashrc` is an Omarchy overlay: it sources `~/.local/share/omarchy/default/bash/rc`, so keep custom exports and aliases below that line.
- This repo is used on two Omarchy machines: desktop `192.168.100.112` and notebook `192.168.100.119`.
- SSH works both ways with `helio@192.168.100.112` and `helio@192.168.100.119`.
- `ssh/` is config-only. Never add private keys, `known_hosts`, or `.pem` / `.ppk` files.
- `claude/` mixes real config with generated state (`todos/`, `statsig/`, `shell-snapshots/`). Avoid editing the generated parts unless that is the task.
- Before touching `nvim/`, read `nvim/CLAUDE.md`. Neovim starts from `nvim/init.lua` -> `lua/config/lazy.lua`, and custom plugin overrides live in `nvim/lua/plugins/`.
- For Neovim changes, validate with `nvim`; use `:Lazy`, `:LspInfo`, and `:LspLog` when debugging.
- Before touching `claude/`, read `claude/CLAUDE.md`.
