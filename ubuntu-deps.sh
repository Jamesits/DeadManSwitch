#!/bin/bash
set -euv

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi

apt-get update -y
apt-get install golang-1.10-go
ln -s /usr/lib/go-1.10/bin/go /usr/local/bin/go
