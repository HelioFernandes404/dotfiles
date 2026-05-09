# 🏠 Dotfiles

Configurações pessoais para sistemas Unix-like (Linux/macOS).

## 📋 Conteúdo

Este repositório contém configurações para:

- **ZSH** - Shell configuration (`.zshrc`, `.zprofile`)
- **Git** - Version control settings (`.gitconfig`)
- **Kitty** - Terminal emulator config (`kitty.conf`, `current-theme.conf`)
- **SSH** - SSH client configuration (`.ssh/config`) ⚠️ **sem chaves privadas**

## 🚀 Instalação

### Primeira instalação

```bash
# Clone o repositório
git clone https://github.com/SEU_USUARIO/dotfiles.git ~/dotfiles

# Entre no diretório
cd ~/dotfiles

# Execute o script de instalação
./install.sh
```

### Atualizando

```bash
cd ~/dotfiles
git pull origin main
./install.sh
```

## 🔄 Como funciona

O script `install.sh` cria **symlinks** dos arquivos de configuração do diretório `~/dotfiles` para o seu `$HOME`. Isso permite que você:

- ✅ Versione suas configurações com Git
- ✅ Mantenha tudo organizado em um único lugar
- ✅ Sincronize entre diferentes máquinas
- ✅ Restaure configurações rapidamente em sistemas novos

### Estrutura

```
~/dotfiles/
├── zsh/
│   ├── zshrc       → ~/.zshrc
│   └── zprofile    → ~/.zprofile
├── git/
│   └── gitconfig   → ~/.gitconfig
├── kitty/
│   ├── kitty.conf          → ~/.config/kitty/kitty.conf
│   └── current-theme.conf  → ~/.config/kitty/current-theme.conf
├── ssh/
│   ├── config      → ~/.ssh/config
│   └── README.md   (avisos de segurança)
├── install.sh      (script de instalação)
├── uninstall.sh    (script de remoção)
└── README.md
```

## 🗑️ Desinstalação

Para remover os symlinks:

```bash
cd ~/dotfiles
./uninstall.sh
```

## 📝 Adicionando novos dotfiles

1. Copie o arquivo para a pasta apropriada em `~/dotfiles/`
2. Adicione uma linha no `install.sh` para criar o symlink
3. Adicione uma linha no `uninstall.sh` para remover o symlink
4. Commit e push

```bash
# Exemplo: adicionar .tmux.conf
mkdir -p ~/dotfiles/tmux
cp ~/.tmux.conf ~/dotfiles/tmux/tmux.conf

# Edite install.sh para adicionar:
# create_symlink "$DOTFILES_DIR/tmux/tmux.conf" "$HOME/.tmux.conf"

# Commit
git add .
git commit -m "Add tmux configuration"
git push
```

## ⚠️ Importante

- Os arquivos originais são **automaticamente backupeados** em `~/dotfiles_backup_*` antes da instalação
- **Nunca** commite arquivos com senhas, tokens ou informações sensíveis
- Revise o `.gitignore` para garantir que arquivos sensíveis não sejam versionados

### ⚠️ Segurança SSH

**ATENÇÃO:** Este repositório contém apenas o arquivo `ssh/config`.

**NUNCA adicione ao Git:**
- ❌ Chaves privadas SSH (`id_rsa`, `id_ed25519`, etc.)
- ❌ Arquivos `known_hosts`
- ❌ Chaves `.pem`, `.ppk`
- ❌ Tokens ou certificados

O `.gitignore` já está configurado para bloquear esses arquivos, mas sempre revise antes de commitar!

## 🔧 Personalização

Antes de fazer push para um repositório público:

1. Revise os arquivos e remova informações pessoais
2. Configure seu Git com suas informações:
   ```bash
   git config user.name "Seu Nome"
   git config user.email "seu@email.com"
   ```

## 📦 Backup automático

Você pode criar um script para fazer commit automático das mudanças:

```bash
#!/bin/bash
cd ~/dotfiles
git add -A
git commit -m "Update dotfiles - $(date +'%Y-%m-%d %H:%M')"
git push
```

## 🌟 Repositórios de inspiração

- [Awesome Dotfiles](https://raw.githubusercontent.com/webpro/awesome-dotfiles/refs/heads/master/README.md)
- [GitHub Dotfiles](https://dotfiles.github.io/)

## 📄 Licença

Use livremente para suas próprias configurações!
