#!/bin/bash
set -eu

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi

export GOPATH=/tmp/go
export GOBIN=$GOPATH/bin

pushd src
make deps
make
chmod +x dmswitch
mkdir -p /usr/local/bin
mv dmswitch /usr/local/bin
#rm -rf /tmp/go
popd

pushd config
mkdir -p /etc/dmswitch
cp config.toml /etc/dmswitch
cp -r hooks /etc/dmswitch
cp dmswitch.service /etc/systemd/system
popd
