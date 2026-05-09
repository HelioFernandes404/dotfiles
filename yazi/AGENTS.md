# Repository Guidelines

## Project Structure & Module Organization

This repository contains personal Yazi configuration files.

- `yazi.toml`: main Yazi settings, manager behavior, preview limits, openers, file-opening rules, task workers, input dialogs, confirmations, and picker layout.
- `keymap.toml`: Yazi keybindings grouped by mode, such as `[mgr]`, `[tasks]`, `[spot]`, `[pick]`, and `[input]`.

There is no source tree, build output, or test fixture directory. Keep new Yazi config files at repository root unless Yazi expects a specific subdirectory.

## Build, Test, and Development Commands

- `yazi`: launch Yazi with this config from the normal config path.
- `taplo fmt yazi.toml keymap.toml`: format TOML if Taplo is installed.
- `taplo lint yazi.toml keymap.toml`: validate TOML syntax and schema hints if Taplo is installed.
- `rg 'pattern'`: search bindings, openers, or setting names quickly.

After changing keybindings or open rules, restart Yazi and test the affected action directly.

## Coding Style & Naming Conventions

Use TOML only. Preserve Yazi section names and command strings exactly, for example `[opener]`, `[open]`, and `shell --block`. Keep comments short and place them above related blocks. Align repeated assignment blocks when it improves scanability, as in `yazi.toml`. In `keymap.toml`, keep entries grouped by mode and workflow with existing comment headers such as `# Navigation`, `# Operation`, and `# Goto`.

Use descriptive opener names with lowercase snake case, such as `nvim_dir`.

## Testing Guidelines

No automated test suite exists. Validate changes by:

- running `taplo lint yazi.toml keymap.toml`;
- launching `yazi`;
- testing changed keys, opener rules, and platform-specific commands.

For shell commands inside TOML, check quoting carefully. Prefer existing patterns like `"$1"` for one file and `"$@"` for selected files.

## Commit & Pull Request Guidelines

No local Git history is available in this directory, so no project-specific commit convention can be inferred. Use concise imperative commits, for example `Add nvim directory opener` or `Adjust Yazi preview limits`.

Pull requests should include a short purpose, changed files, manual validation steps, and any platform assumptions such as Linux-only `xdg-open`, `mpv`, `fd`, `rg`, or `zoxide` usage.

## Agent-Specific Instructions

Read existing TOML before editing. Keep changes narrow, preserve user keybinding intent, and do not rewrite unrelated formatting. Use Context7 to consult current Yazi documentation before changing Yazi-specific options, schemas, commands, plugin calls, or keymap syntax.
