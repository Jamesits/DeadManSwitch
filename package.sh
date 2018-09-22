#!/bin/bash
set -eu

export GOPATH=/tmp/go
export GOBIN=$GOPATH/bin
NOAUTOARCHIVE=0

while test $# -gt 0
do
    case "$1" in
        --noautoarchive) NOAUTOARCHIVE=1
            ;;
    esac
    shift
done

rm -rf release
mkdir -p release

pushd src
make deps
make
chmod +x dmswitch
popd

mv src/dmswitch release
cp install-dist.sh release/install.sh

mkdir -p release/etc/dmswitch
cp -r config/config.toml config/hooks release/etc/dmswitch

mkdir -p release/etc/systemd/system
cp config/dmswitch.service release/etc/systemd/system

mv release dmswitch

if [ $NOAUTOARCHIVE = 0 ]; then
    tar -cvzf dmswitch-release.tar.gz dmswitch
    rm -rf dmswitch
fi
