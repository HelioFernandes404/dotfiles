# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a LazyVim configuration for Neovim. LazyVim is a Neovim setup powered by lazy.nvim plugin manager.

### Core Structure

- `init.lua`: Entry point that bootstraps the configuration via `require("config.lazy")`
- `lua/config/`: Core configuration files
  - `lazy.lua`: Lazy.nvim bootstrap and plugin specification loader
  - `options.lua`: Vim options (extends LazyVim defaults)
  - `keymaps.lua`: Custom keymaps (extends LazyVim defaults)
  - `autocmds.lua`: Autocommands (extends LazyVim defaults)
- `lua/plugins/`: Plugin specifications and overrides
  - Each `.lua` file defines plugin specs that override or extend LazyVim defaults
  - All files are auto-loaded by lazy.nvim except `example.lua` (disabled with `if true then return {} end`)

### Plugin Loading System

The configuration uses lazy.nvim's spec system defined in `lua/config/lazy.lua`:
1. Loads base LazyVim distribution
2. Imports LazyVim extras: `lazyvim.plugins.extras.dap.core` and `lazyvim.plugins.extras.lang.python`
3. Auto-imports all plugin specs from `lua/plugins/` directory
4. Plugin checker enabled but notifications disabled

## Python Development Setup

This configuration uses a modern Python toolchain with **ty** as the primary LSP:

### Language Servers (configured in `lua/plugins/python.lua` and related files)

- **ty** (`lua/plugins/ty.lua`): Primary type checker, LSP, and IntelliSense provider
  - Provides completions, hover, diagnostics, and inlay hints
  - Workspace-wide type checking enabled
  - Auto-imports enabled
  - Custom lspconfig registration with root pattern detection

- **ruff** (`lua/plugins/lsp_ruff.lua`): Linting and formatting (complementary to ty)
  - Hover provider disabled to prioritize ty
  - Custom keybindings:
    - `<leader>co`: Organize imports
    - `<leader>cf`: Fix all Ruff issues
  - Line length: 88 characters
  - Extends select: I (isort), E (pycodestyle errors), F (pyflakes)

- **pyright/basedpyright**: Explicitly disabled (redundant with ty)

### Debugging (configured in `lua/plugins/dap.lua`)

- Uses nvim-dap with debugpy backend
- Auto-loads `.vscode/launch.json` configurations when starting debug sessions
- Hook on `dap.continue()` automatically loads launch.json if present

## Key Custom Behaviors

### Theme System

- Symlink `lua/plugins/theme.lua` points to external theme manager: `/home/helio/.config/omarchy/current/theme/neovim.lua`
- Theme hot-reload capability configured in `lua/plugins/omarchy-theme-hotreload.lua`
- All theme collection available via `lua/plugins/all-themes.lua`

### Disabled Features

- Relative line numbers: `vim.opt.relativenumber = false` (in `lua/config/options.lua`)
- Snacks animated scrolling: disabled via `lua/plugins/snacks-animated-scrolling-off.lua`
- News alerts: disabled via `lua/plugins/disable-news-alert.lua`

## Working with This Configuration

### Testing Changes

Start Neovim to test configuration:
```bash
nvim
```

Check lazy.nvim status and installed plugins:
```
:Lazy
```

### LSP Commands

Check LSP status:
```
:LspInfo
```

View LSP logs:
```
:LspLog
```

### DAP (Debugging)

Place debug configurations in `.vscode/launch.json` in your project root. They will be auto-loaded when you start debugging.

### Plugin Development Pattern

Plugin specs in `lua/plugins/` follow this structure:
- Return a table (or array of tables) with plugin specifications
- Use `opts` to merge with default LazyVim configuration
- Use `setup` function when custom initialization is needed (see ty.lua for custom LSP registration)
- Use `config` function for plugin-specific setup (see dap.lua for hooks)

### Important File Patterns

When working with LSP configurations:
- Root directory detection uses: `pyproject.toml`, `setup.py`, `setup.cfg`, `requirements.txt`, `Pipfile`, `.git`
- Launch.json path: `.vscode/launch.json`

## Code Style

- Lua formatting: Uses stylua (config in `stylua.toml`)
- Python formatting: Handled by ruff (88 char line length)
