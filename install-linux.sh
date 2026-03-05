#!/usr/bin/env bash
# install.sh — Install Lamp as a desktop app on Linux
set -e

BINARY_NAME="lamp"
BINARY_SRC="./$BINARY_NAME"
INSTALL_BIN="/usr/local/bin/$BINARY_NAME"
INSTALL_DESKTOP="$HOME/.local/share/applications/lamp.desktop"
INSTALL_ICON="$HOME/.local/share/icons/lamp.png"

# Build the binary if not already built
if [ ! -f "$BINARY_SRC" ]; then
  echo "Building $BINARY_NAME..."
  go build -o "$BINARY_NAME" .
fi

# Install binary
echo "Installing binary to $INSTALL_BIN..."
sudo cp "$BINARY_SRC" "$INSTALL_BIN"
sudo chmod +x "$INSTALL_BIN"

# Install icon (optional — replace lamp.png with your own)
mkdir -p "$(dirname "$INSTALL_ICON")"
if [ -f "./lamp.png" ]; then
  cp "./lamp.png" "$INSTALL_ICON"
else
  echo "No lamp.png found — skipping icon install. You can add one later."
fi

# Install .desktop file
mkdir -p "$(dirname "$INSTALL_DESKTOP")"
cat > "$INSTALL_DESKTOP" <<EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=Lamp
Comment=Lamp Terminal
Exec=$INSTALL_BIN
Icon=$INSTALL_ICON
Terminal=false
Categories=System;TerminalEmulator;
StartupNotify=true
EOF

# Refresh desktop database
if command -v update-desktop-database &>/dev/null; then
  update-desktop-database "$HOME/.local/share/applications"
fi

echo ""
echo "Done! Lamp is now installed."
echo "You can find it in your application launcher, or run: $BINARY_NAME"