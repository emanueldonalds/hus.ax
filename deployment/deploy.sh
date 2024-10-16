#!/bin/bash

set -o nounset

BASE_DIR=/home/properties
APP_NAME=property-viewer

REPO_DIR=$BASE_DIR/$APP_NAME
DEPL_DIR="$REPO_DIR-depl"

rm -rf $DEPL_DIR
mkdir $DEPL_DIR
mkdir $DEPL_DIR/rss

cp -r $REPO_DIR/src/assets $DEPL_DIR/assets
cp $REPO_DIR/src/rss/template.xml $DEPL_DIR/rss/.

cd $REPO_DIR/src
env GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -o $DEPL_DIR/$APP_NAME

sudo systemctl daemon-reload
sudo systemctl restart $APP_NAME
