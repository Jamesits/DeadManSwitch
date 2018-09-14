#!/bin/bash
set -u

systemctl disable dmswitch
rm /etc/systemd/system/dmswitch.service
systemctl daemon-reload

rm -rf /usr/local/bin/dmswitch
rm -rf /etc/dmswitch

fstrim -a

exit 0
