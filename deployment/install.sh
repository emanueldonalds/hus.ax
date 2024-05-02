#!/bin/bash

if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

REPO_DIR=../
DEPLOYMENT_DIR=../
OVERRIDE_CONF=/etc/systemd/system/property-viewer.service.d/override.conf

sudo ln -s $REPO_DIR/deployment/property-viewer.service /etc/systemd/system
        
sudo touch $OVERRIDE_CONF

if [ -d '/etc/systemd/system/property-viewer.service.d' ]; then
        echo 'Env dir already created'
else
        sudo mkdir /etc/systemd/system/property-viewer.service.d
        sudo echo "[Service]" >> $OVERRIDE_CONF
        sudo echo "Environment=\"PROPERTY_VIEWER_ASSETS_DIR=<assets-dir>\"" >> $OVERRIDE_CONF
        sudo echo "Environment=\"PROPERTY_VIEWER_DB_HOST=<db-host>\"" >> $OVERRIDE_CONF
        sudo echo "Environment=\"PROPERTY_VIEWER_DB_PASSWORD=<db-password>\"" >> $OVERRIDE_CONF
        echo "Created file '$OVERRIDE_CONF'. Set the environment variables in this file."
fi

