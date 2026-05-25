DOTFILES_BIN := ./bin/dotfiles
INSTALL_BIN := $(HOME)/.local/bin/dotfiles

.PHONY: bootstrap build test validate status install-core install-editor install-terminal uninstall-core

bootstrap: build
	mkdir -p "$(HOME)/.local/bin"
	cp "$(DOTFILES_BIN)" "$(INSTALL_BIN)"

build:
	mkdir -p bin
	go build -o "$(DOTFILES_BIN)" ./cmd/dotfiles

test:
	go test ./...

validate: build
	$(DOTFILES_BIN) validate

status: build
	$(DOTFILES_BIN) status

install-core: build
	$(DOTFILES_BIN) install --group core

install-editor: build
	$(DOTFILES_BIN) install --group editor

install-terminal: build
	$(DOTFILES_BIN) install --group terminal

uninstall-core: build
	$(DOTFILES_BIN) uninstall --group core
