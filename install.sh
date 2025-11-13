#!/bin/bash
# Script de instalação de dotfiles
# Cria symlinks dos dotfiles para o diretório home

set -e

DOTFILES_DIR="$HOME/dotfiles"
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
            mv "$target" "$BACKUP_DIR/"
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

echo ""
echo "✨ Instalação concluída!"
echo "🔄 Execute 'source ~/.zshrc' ou reinicie o terminal para aplicar as mudanças"

if [ -d "$BACKUP_DIR" ] && [ "$(ls -A $BACKUP_DIR)" ]; then
    echo "📦 Backup dos arquivos originais em: $BACKUP_DIR"
else
    rmdir "$BACKUP_DIR" 2>/dev/null || true
fi
