#!/bin/bash
set -eu

export GOPATH=/tmp/go
export GOBIN=$GOPATH/bin

rm -rf release
mkdir -p release

pushd src
make deps
make
chmod +x dmswitch
popd

mv src/dmswitch release

mkdir -p release/etc/dmswitch
cp -r config/config.toml config/hooks release/etc/dmswitch

mkdir -p release/etc/systemd/system
cp config/dmswitch.service release/etc/systemd/system


mv release dmswitch
tar -cvzf dmswitch-release.tar.gz dmswitch
rm -rf dmswitch

