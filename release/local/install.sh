#!/usr/bin/env bash

set -e -o pipefail

if [ -d /usr/local/go ]; then
  export PATH="$PATH:/usr/local/go/bin"
fi

DIR=$(dirname "$0")
PROJECT=$DIR/../..

pushd $PROJECT
go install -v -trimpath -ldflags "-s -w -buildid=" -tags "with_acme" ./cmd/serenity
popd

sudo cp $(go env GOPATH)/bin/serenity /usr/local/bin/
sudo mkdir -p /usr/local/etc/serenity
sudo cp $PROJECT/release/config/config.json /usr/local/etc/serenity/config.json
sudo cp $DIR/serenity.service /etc/systemd/system
sudo systemctl daemon-reload
