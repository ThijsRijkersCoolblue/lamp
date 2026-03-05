#!/usr/bin/env bash
# install-macos.sh — Package Lamp as a native macOS .app (Fyne, no Terminal)
set -e

APP_NAME="Lamp"
BINARY_NAME="lamp"
BINARY_SRC="./$BINARY_NAME"
APP_DIR="./$APP_NAME.app"
CONTENTS="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS/MacOS"
RESOURCES_DIR="$CONTENTS/Resources"

if [ ! -f "$BINARY_SRC" ]; then
  echo "Building $BINARY_NAME..."
  go build -o "$BINARY_NAME" .
fi

echo "Creating $APP_NAME.app bundle..."
rm -rf "$APP_DIR"
mkdir -p "$MACOS_DIR"
mkdir -p "$RESOURCES_DIR"

cp "$BINARY_SRC" "$MACOS_DIR/$BINARY_NAME"
chmod +x "$MACOS_DIR/$BINARY_NAME"

if [ -f "./lamp.icns" ]; then
  cp "./lamp.icns" "$RESOURCES_DIR/lamp.icns"
  ICON_LINE="<key>CFBundleIconFile</key><string>lamp</string>"
else
  echo "No lamp.icns found — skipping icon."
  ICON_LINE=""
fi

cat > "$CONTENTS/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleName</key>
  <string>$APP_NAME</string>
  <key>CFBundleDisplayName</key>
  <string>$APP_NAME</string>
  <key>CFBundleIdentifier</key>
  <string>com.yourname.lamp</string>
  <key>CFBundleVersion</key>
  <string>1.0.0</string>
  <key>CFBundleExecutable</key>
  <string>$BINARY_NAME</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>NSHighResolutionCapable</key>
  <true/>
  <key>NSPrincipalClass</key>
  <string>NSApplication</string>
  <key>NSSupportsAutomaticGraphicsSwitching</key>
  <true/>
  <key>LSMinimumSystemVersion</key>
  <string>10.13</string>
  <key>LSUIElement</key>
  <false/>
  $ICON_LINE
</dict>
</plist>
PLIST

# Remove any quarantine/extended attributes that cause Terminal to open
xattr -cr "$APP_DIR" 2>/dev/null || true

echo ""
echo "Done! $APP_NAME.app created in the current directory."
echo "Run: sudo cp -r $APP_NAME.app /Applications/"