#!/bin/bash

if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

DEPLOYMENT_DIR=/home/properties/property-viewer

sudo ln -s $DEPLOYMENT_DIR/property-viewer.service /etc/systemd/system

sudo touch /etc/systemd/system/property-viewer.service.d/override.conf

if [ -d '/etc/systemd/system/property-viewer.service.d' ]; then
        echo 'Env dir already created'
else
        echo "Creating env file at '/etc/systemd/system/property-viewer.service.d/override.conf'. Enter environment variables here."
        sudo mkdir /etc/systemd/system/property-viewer.service.d
        sudo echo "[Service]" >> /etc/systemd/system/property-viewer.service.d/override.conf;
        sudo echo "Environment=\"PROPERTY_VIEWER_ASSETS_DIR=<assets-dir>\"" >> /etc/systemd/system/property-viewer.service.d/override.conf;
        sudo echo "Environment=\"PROPERTY_VIEWER_DB_PASSWORD=<db-password>\"" >> /etc/systemd/system/property-viewer.service.d/override.conf;
fi

