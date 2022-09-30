#!/bin/sh
sudo su
cd /opt/yourteams
docker-compose down
cd /opt
rm -rf /opt/yourteams
mv /home/debian/yourteams /opt/
cd ./yourteams
docker-compose up -d --remove-orphans
