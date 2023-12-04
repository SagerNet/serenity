#!/usr/bin/env bash

sudo systemctl stop serenity
sudo rm -rf /var/lib/serenity
sudo rm -rf /usr/local/bin/serenity
sudo rm -rf /usr/local/etc/serenity
sudo rm -rf /etc/systemd/system/serenity.service
sudo systemctl daemon-reload
