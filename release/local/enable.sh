#!/usr/bin/env bash

set -e -o pipefail

sudo systemctl enable serenity
sudo systemctl start serenity
sudo journalctl -u serenity --output cat -f
