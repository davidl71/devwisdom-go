#!/bin/bash
# Installation script for devwisdom zsh plugin

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLUGIN_DIR="$SCRIPT_DIR"

echo -e "${GREEN}Installing devwisdom zsh plugin...${NC}"

# Detect zsh framework
if [[ -d "$HOME/.oh-my-zsh" ]]; then
    echo -e "${YELLOW}Detected oh-my-zsh${NC}"
    TARGET_DIR="$HOME/.oh-my-zsh/custom/plugins/devwisdom"
    FRAMEWORK="oh-my-zsh"
elif [[ -d "$HOME/.zsh" ]]; then
    echo -e "${YELLOW}Detected standard zsh directory${NC}"
    TARGET_DIR="$HOME/.zsh/plugins/devwisdom"
    FRAMEWORK="standard"
else
    echo -e "${YELLOW}No zsh framework detected, using ~/.zsh/plugins${NC}"
    TARGET_DIR="$HOME/.zsh/plugins/devwisdom"
    FRAMEWORK="standard"
fi

# Create target directory
mkdir -p "$TARGET_DIR"

# Copy plugin files
echo -e "Copying plugin files to $TARGET_DIR..."
cp "$PLUGIN_DIR/devwisdom.plugin.zsh" "$TARGET_DIR/"
cp "$PLUGIN_DIR/_devwisdom" "$TARGET_DIR/"

# Make completion file executable
chmod +x "$TARGET_DIR/_devwisdom"

# Add to zshrc if not already present
ZSH_RC="$HOME/.zshrc"
if [[ "$FRAMEWORK" == "oh-my-zsh" ]]; then
    if ! grep -q "devwisdom" "$ZSH_RC" 2>/dev/null; then
        echo ""
        echo -e "${YELLOW}Add the following to your ~/.zshrc:${NC}"
        echo -e "${GREEN}plugins=(... devwisdom)${NC}"
    else
        echo -e "${GREEN}Plugin already configured in ~/.zshrc${NC}"
    fi
else
    # For standard zsh, add source line
    SOURCE_LINE="source $TARGET_DIR/devwisdom.plugin.zsh"
    if ! grep -q "$SOURCE_LINE" "$ZSH_RC" 2>/dev/null; then
        echo "" >> "$ZSH_RC"
        echo "# devwisdom plugin" >> "$ZSH_RC"
        echo "$SOURCE_LINE" >> "$ZSH_RC"
        echo -e "${GREEN}Added plugin to ~/.zshrc${NC}"
    else
        echo -e "${GREEN}Plugin already configured in ~/.zshrc${NC}"
    fi
    
    # Add completion directory to fpath if not present
    COMPLETION_DIR="$TARGET_DIR"
    if ! grep -q "fpath=($COMPLETION_DIR" "$ZSH_RC" 2>/dev/null; then
        echo "fpath=($COMPLETION_DIR \$fpath)" >> "$ZSH_RC"
        echo "autoload -Uz compinit && compinit" >> "$ZSH_RC"
        echo -e "${GREEN}Added completion to ~/.zshrc${NC}"
    fi
fi

echo ""
echo -e "${GREEN}✓ Installation complete!${NC}"
echo ""
echo "To use the plugin:"
echo "  1. Restart your shell or run: source ~/.zshrc"
echo "  2. Use commands: devwisdom-daily, devwisdom-quote, devwisdom-consult"
echo ""
echo "Optional: Enable auto-daily wisdom on shell startup:"
echo "  export DEVWISDOM_AUTO_DAILY=true"
echo "  (add to ~/.zshrc)"
