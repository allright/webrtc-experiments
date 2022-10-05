#!/bin/sh
 cd ./webrtc-min-client
 npm run build
 cd ..
 ./deploy.sh
 