#!/bin/bash

REPO_DIR="../"

cd $REPO_DIR/src/
templ generate
env GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -o $REPO_DIR/property-viewer

sudo systemctl daemon-reload
sudo systemctl restart property-viewer
