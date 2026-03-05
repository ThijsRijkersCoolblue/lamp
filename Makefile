BINARY_NAME := lamp
BUILD_FLAGS := -ldflags="-s -w"

.PHONY: build install-linux install-macos clean

build:
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) .

install-linux: build
	bash install-linux.sh

install-macos: build
	bash install-macos.sh
	# Ad-hoc codesign so macOS treats it as a proper app, not a CLI tool
	codesign --force --deep --sign - Lamp.app
	sudo cp -r Lamp.app /Applications/
	# Clear quarantine on installed app
	sudo xattr -cr /Applications/Lamp.app

clean:
	rm -f $(BINARY_NAME)
	rm -rf Lamp.app