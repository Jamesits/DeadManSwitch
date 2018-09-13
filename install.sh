#!/bin/bash
set -eu

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi

pushd src
make
chmod +x dmswitch
mkdir -p /usr/local/bin
mv dmswitch /usr/local/bin
popd

pushd config
mkdir -p /etc/dmswitch
cp config.toml /etc/dmswitch
cp dmswitch.service /etc/systemd/system
popd