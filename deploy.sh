#!/bin/sh
scp -r ./yourteams vk:
ssh vk 'bash /dev/stdin' < on_host_script.sh
