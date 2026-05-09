#!/bin/bash
# Script de instalação de dotfiles
# Cria symlinks dos dotfiles para o diretório home

set -e

# Always resolve the dotfiles directory from this script's location.
# Avoids relying on ~/dotfiles existing or being the current directory.
DOTFILES_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_DIR="$HOME/dotfiles_backup_$(date +%Y%m%d_%H%M%S)"

echo "🔧 Instalando dotfiles..."

# Criar diretório de backup se necessário
if [ ! -d "$BACKUP_DIR" ]; then
    mkdir -p "$BACKUP_DIR"
fi

# Função para criar symlink com backup
create_symlink() {
    local source="$1"
    local target="$2"

    if [ -e "$target" ] || [ -L "$target" ]; then
        if [ ! -L "$target" ]; then
            echo "  📦 Fazendo backup: $target"
            # Avoid collisions when backing up different paths with the same basename.
            local backup_name
            backup_name="${target#/}"
            backup_name="${backup_name//\//__}"

            local backup_target="$BACKUP_DIR/$backup_name"
            if [ -e "$backup_target" ] || [ -L "$backup_target" ]; then
                backup_target="$backup_target.$(date +%s%N)"
            fi

            mv "$target" "$backup_target"
        else
            echo "  🔗 Removendo symlink antigo: $target"
            rm "$target"
        fi
    fi

    echo "  ✅ Criando symlink: $target -> $source"
    ln -s "$source" "$target"
}

# ZSH
echo "📝 Configurando ZSH..."
create_symlink "$DOTFILES_DIR/zsh/zshrc" "$HOME/.zshrc"
create_symlink "$DOTFILES_DIR/zsh/zprofile" "$HOME/.zprofile"

# Git
echo "🌿 Configurando Git..."
create_symlink "$DOTFILES_DIR/git/gitconfig" "$HOME/.gitconfig"

# Kitty Terminal
if [ -d "$DOTFILES_DIR/kitty" ]; then
    echo "🐱 Configurando Kitty..."
    mkdir -p "$HOME/.config/kitty"
    create_symlink "$DOTFILES_DIR/kitty/kitty.conf" "$HOME/.config/kitty/kitty.conf"
    if [ -f "$DOTFILES_DIR/kitty/current-theme.conf" ]; then
        create_symlink "$DOTFILES_DIR/kitty/current-theme.conf" "$HOME/.config/kitty/current-theme.conf"
    fi
fi

# SSH Config (com permissões corretas)
if [ -f "$DOTFILES_DIR/ssh/config" ]; then
    echo "🔐 Configurando SSH..."
    mkdir -p "$HOME/.ssh"
    create_symlink "$DOTFILES_DIR/ssh/config" "$HOME/.ssh/config"
    chmod 700 "$HOME/.ssh"
    chmod 600 "$HOME/.ssh/config"
fi

# Yazi File Manager
if [ -d "$DOTFILES_DIR/config/yazi" ]; then
    echo "📁 Configurando Yazi..."
    mkdir -p "$HOME/.config"
    create_symlink "$DOTFILES_DIR/config/yazi" "$HOME/.config/yazi"
fi

# Omarchy (user overrides)
if [ -d "$DOTFILES_DIR/system/themes/github-dark/omarchy/theme" ]; then
    echo "🎨 Configurando Omarchy (theme github-dark)..."
    mkdir -p "$HOME/.config/omarchy/themes" "$HOME/.config/omarchy/backgrounds"
    create_symlink "$DOTFILES_DIR/system/themes/github-dark/omarchy/theme" "$HOME/.config/omarchy/themes/github-dark"
    create_symlink "$DOTFILES_DIR/system/themes/github-dark/omarchy/backgrounds" "$HOME/.config/omarchy/backgrounds/github-dark"
fi

if [ -d "$DOTFILES_DIR/config/omarchy" ]; then
    echo "🎨 Configurando Omarchy (hooks/templates)..."
    mkdir -p "$HOME/.config/omarchy/themed" "$HOME/.config/omarchy/hooks"
    if [ -f "$DOTFILES_DIR/config/omarchy/themed/starship.toml.tpl" ]; then
        create_symlink "$DOTFILES_DIR/config/omarchy/themed/starship.toml.tpl" "$HOME/.config/omarchy/themed/starship.toml.tpl"
    fi
    if [ -f "$DOTFILES_DIR/config/omarchy/hooks/theme-set" ]; then
        create_symlink "$DOTFILES_DIR/config/omarchy/hooks/theme-set" "$HOME/.config/omarchy/hooks/theme-set"
        chmod +x "$HOME/.config/omarchy/hooks/theme-set" || true
    fi
fi

# Fish
if [ -d "$DOTFILES_DIR/config/fish" ]; then
    echo "🐟 Configurando Fish..."
    mkdir -p "$HOME/.config/fish/conf.d"
    if [ -f "$DOTFILES_DIR/config/fish/conf.d/omarchy-starship-config.fish" ]; then
        create_symlink "$DOTFILES_DIR/config/fish/conf.d/omarchy-starship-config.fish" "$HOME/.config/fish/conf.d/omarchy-starship-config.fish"
    fi
fi

# Neovim
if [ -d "$DOTFILES_DIR/config/nvim" ]; then
    echo "✏️  Configurando Neovim..."
    mkdir -p "$HOME/.config"
    create_symlink "$DOTFILES_DIR/config/nvim" "$HOME/.config/nvim"
fi

# Claude Code
if [ -d "$DOTFILES_DIR/claude" ]; then
    echo "🤖 Configurando Claude Code..."
    create_symlink "$DOTFILES_DIR/claude" "$HOME/.claude"
fi

echo ""
echo "✨ Instalação concluída!"
echo "🔄 Execute 'source ~/.zshrc' ou reinicie o terminal para aplicar as mudanças"

if [ -d "$BACKUP_DIR" ] && [ "$(ls -A $BACKUP_DIR)" ]; then
    echo "📦 Backup dos arquivos originais em: $BACKUP_DIR"
else
    rmdir "$BACKUP_DIR" 2>/dev/null || true
fi
