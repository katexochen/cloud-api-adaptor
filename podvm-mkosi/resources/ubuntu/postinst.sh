#!/usr/bin/env bash

set -euxo pipefail

# Somehow systemd-networkd is not enabled during boot, even when the preset
# is enabled. Linking the file manually as a workaround.
ln -s /usr/lib/systemd/system/systemd-networkd.service \
    "${BUILDROOT}/etc/systemd/system/multi-user.target.wants/systemd-networkd.service"
