.PHONY: elm go exe

APP=Counter.app
APPEXE=counter

clear:
	rm -rf bin

setup: clear
	mkdir -p bin/macos/$(APP)/Contents/{MacOS,Resources}

# Package Mac OS assets
asset: setup
	cp go/assets/icons/icon.icns bin/macos/$(APP)/Contents/Resources/icon.icns
	cp go/assets/info.plist bin/macos/$(APP)/Contents/Info.plist

# Generate javascript from Elm code
js:
	cd elm; elm make src/Main.elm --output=../go/www/main.js

# Convert resources to Go code
generate: asset js
	cd go; go run -tags generate gen.go

# Create an executable for MacOS
exe: generate
	cd go; go build -o ../bin/macos/$(APP)/Contents/MacOS/$(APPEXE)

# --- Extra
beautify:
	cd elm; elm-format src/ --yes

test:
	cd elm; elm-test

# Download Elm and Go dependencies
deps:
	cd go; go get -u github.com/zserge/lorca
