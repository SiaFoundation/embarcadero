#!/bin/bash
# Must be run on macOS, does not support any other OS due to Apple Gatekeeper
# notarization
rm -rf release dist
mkdir release

cleanup() {
	rm -rf dist
}
trap cleanup EXIT

# use either the tag name or short commit hash
RELEASE=$(git describe --tags --exact-match --abbrev=0)
if [ $? -ne 0 ]; then
	RELEASE=$(git log -1 --pretty=format:%h)
fi

for OS in linux windows darwin; do
	for ARCH in amd64 arm64; do
		echo "Building $RELEASE $OS/$ARCH"
		rm -rf dist
		mkdir -p dist/embarcadero
		BIN=embc
		if [ $OS = "windows" ]; then
			BIN=embc.exe
		fi
		GOOS=$OS GOARCH=$ARCH go build -trimpath -ldflags='-s -w' -o dist/embarcadero/$bin .
		cp README.md dist/embarcadero/
		ZIP_OUTPUT="release/embarcadero_${RELEASE}_${OS}_${ARCH}.zip"
		if [ "$OS" = "darwin" ]; then
			codesign --deep -f -v --timestamp -o runtime,library -s $APPLE_CERT_ID dist/embarcadero/embc
			ditto -ck dist/embarcadero $ZIP_OUTPUT
			xcrun notarytool submit -k ~/private_keys/AuthKey_$APPLE_API_KEY.p8 -d $APPLE_API_KEY -i $APPLE_API_ISSUER --wait --timeout 10m $ZIP_OUTPUT
		else
			zip -qj $ZIP_OUTPUT dist/embarcadero/*
		fi
	done
done
