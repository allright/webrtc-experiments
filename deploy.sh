#!/bin/sh
rm -rf ./yourteams
mkdir yourteams
cp ./nginx.conf.prod ./yourteams/nginx.conf
cp ./docker-compose-prod.yml ./yourteams/docker-compose.yml
cp -a ./webrtc-min-client/dist ./yourteams/html
scp -r ./yourteams vk:
ssh vk 'bash /dev/stdin' < on_host_script.sh
rm -rf ./yourteams