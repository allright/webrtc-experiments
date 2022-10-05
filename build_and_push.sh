#!/bin/sh
docker build ./WebRTCSignal -t registry.timephone.org:5000/wrtsig
docker push registry.timephone.org:5000/wrtsig
