#!/bin/bash

set -o nounset

if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

APP_NAME=property-viewer
REPO_DIR=$HOME/property-viewer
DEPL_DIR="$HOME/$APP_NAME-dist"
OVERRIDE_CONF=/etc/systemd/system/$APP_NAME.service.d/override.conf

echo "Repo dir ${REPO_DIR}"
echo "Deployment dir ${DEPL_DIR}"

sudo ln -s $REPO_DIR/deployment/$APP_NAME.service /etc/systemd/system
        
sudo touch $OVERRIDE_CONF

if [ -d '/etc/systemd/system/$APP_NAME.service.d' ]; then
    echo 'Env dir already created'
else
    sudo mkdir /etc/systemd/system/$APP_NAME.service.d
    sudo echo "[Service]" >> $OVERRIDE_CONF
    sudo echo "Environment=\"PROPERTY_VIEWER_DB_HOST=<db-host>\"" >> $OVERRIDE_CONF
    sudo echo "Environment=\"PROPERTY_VIEWER_DB_PASSWORD=<db-password>\"" >> $OVERRIDE_CONF
    echo "Created file '$OVERRIDE_CONF'. Set the environment variables in this file."
fi

