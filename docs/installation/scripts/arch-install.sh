#!/bin/bash

set -e -o pipefail

ARCH_RAW=$(uname -m)
case "${ARCH_RAW}" in
    'x86_64')    ARCH='amd64';;
    'x86' | 'i686' | 'i386')     ARCH='386';;
    'aarch64' | 'arm64') ARCH='arm64';;
    'armv7l')   ARCH='armv7';;
    's390x')    ARCH='s390x';;
    *)          echo "Unsupported architecture: ${ARCH_RAW}"; exit 1;;
esac

VERSION=$(curl -s https://api.github.com/repos/SagerNet/serenity/releases/latest \
    | grep tag_name \
    | cut -d ":" -f2 \
    | sed 's/\"//g;s/\,//g;s/\ //g;s/v//')

curl -Lo serenity.pkg.tar.zst "https://github.com/SagerNet/serenity/releases/download/v${VERSION}/serenity_${VERSION}_linux_${ARCH}.pkg.tar.zst"
sudo pacman -U serenity.pkg.tar.zst
rm serenity.pkg.tar.zst
