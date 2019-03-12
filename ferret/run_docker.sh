#!/bin/bash

echo 'Remove alpeware/chrome-headless-stable container including the volume'
docker stop chrome-headless
sleep 5
docker rm -v chrome-headless
sudo rm -rf /tmp/chromedata

echo 'Starting alpeware/chrome-headless-stable Docker container'
docker pull alpeware/chrome-headless-stable
docker run -d -p=127.0.0.1:9222:9222 --name=chrome-headless -v /tmp/chromedata/:/data alpeware/chrome-headless-stable
