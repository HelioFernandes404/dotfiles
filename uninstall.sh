#!/bin/bash
# Script de desinstalação de dotfiles
# Remove os symlinks criados

set -e

DOTFILES_DIR="$HOME/dotfiles"

echo "🗑️  Removendo symlinks dos dotfiles..."

remove_symlink() {
    local target="$1"

    if [ -L "$target" ]; then
        echo "  ❌ Removendo: $target"
        rm "$target"
    else
        echo "  ⚠️  Não é um symlink (ignorando): $target"
    fi
}

# ZSH
remove_symlink "$HOME/.zshrc"
remove_symlink "$HOME/.zprofile"

# Git
remove_symlink "$HOME/.gitconfig"

# Kitty
remove_symlink "$HOME/.config/kitty/kitty.conf"
remove_symlink "$HOME/.config/kitty/current-theme.conf"

# SSH
remove_symlink "$HOME/.ssh/config"

echo ""
echo "✅ Desinstalação concluída!"
echo "💡 Restaure seus backups de ~/dotfiles_backup_* se necessário"
