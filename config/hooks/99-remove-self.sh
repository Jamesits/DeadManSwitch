#!/bin/bash
set -u

systemctl disable dmswitch
rm /etc/systemd/system/dmswitch.service
systemctl daemon-reload

rm -rf /etc/dmswitch
rm -rf /usr/local/bin/dmswitch

