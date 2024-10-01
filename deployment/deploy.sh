#!/bin/bash

set -o nounset

REPO_DIR=$HOME/property-viewer
DEPL_DIR=$HOME/property-viewer-depl

rm -rf $DEPL_DIR
mkdir $DEPL_DIR
mkdir $DEPL_DIR/rss

cp -r $REPO_DIR/src/assets $DEPL_DIR/assets
cp $REPO_DIR/src/rss/template.xml $DEPL_DIR/rss/.

cd $REPO_DIR/src
env GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -o $DEPL_DIR/property-viewer

sudo systemctl daemon-reload
sudo systemctl restart property-viewer
