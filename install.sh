#!/bin/bash
set -e

INSTALL_DIR="$HOME/.local/bin"
PATH_ENTRY="export PATH=\"$INSTALL_DIR:\$PATH\""

if ! command -v go &> /dev/null; then
    echo "Erro: Go não está instalado."
    exit 1
fi

if ! command -v git &> /dev/null; then
    echo "Erro: Git não está instalado."
    exit 1
fi

if [ ! -d "$INSTALL_DIR" ]; then
    mkdir -p "$INSTALL_DIR"
fi

TEMP_PATH=$(mktemp -d)
cd "$TEMP_PATH"

git clone https://github.com/atendi9/capivara.git .
go build -o capivara main.go
mv capivara "$INSTALL_DIR/"

cd "$HOME"
rm -rf "$TEMP_PATH"

for CONFIG_FILE in "$HOME/.bashrc" "$HOME/.zshrc"; do
    if [ -f "$CONFIG_FILE" ]; then
        if ! grep -q "$INSTALL_DIR" "$CONFIG_FILE"; then
            echo "" >> "$CONFIG_FILE"
            echo "$PATH_ENTRY" >> "$CONFIG_FILE"
        fi
    fi
done

echo "=> Installation complete."
echo "=> Please restart your terminal or run: source ~/.bashrc (or ~/.zshrc)"