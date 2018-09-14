#!/bin/bash
set -eu

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi

mkdir -p /usr/local/bin
cp dmswitch /usr/local/bin

mkdir -p /etc/systemd/system
cp -r etc/systemd/system/* /etc/systemd/system

cp -r etc/dmswitch /etc

echo "Install finished."
