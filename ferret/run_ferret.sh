#!/bin/bash

# Manually start Chrome
# Tested on portable version 72.0.3626.0 (Developer Build) (64-bit)
# DL: https://www.googleapis.com/download/storage/v1/b/chromium-browser-snapshots/o/Linux_x64%2F612434%2Fchrome-linux.zip?generation=1543535498055804&alt=media
# cd /home/lmaly/Downloads/chrome-linux;CHROME_DEVEL_SANDBOX="$PWD/chrome_sandbox" /home/lmaly/Downloads/chrome-linux/chrome --remote-debugging-port=9222

echo 'Running dzone-refcardz-dl.fql via headless Chrome'
PATH=/home/lmaly/Downloads/chrome-linux/:$PATH ferret --param=userid:\"3590306\" --param=username:\"dzone-refcardz@mailcatch.com\" --param=password:\"password123456\" --cdp http://127.0.0.1:9222 dzone-refcardz-not-direct-dl.fql > dzone-refcardz.json

echo 'Generating the JSON file...'
cat dzone-refcardz.json | jq 'sort_by(.name)' > dzone-refcardz$(date +"%Y%m%d").json
rm dzone-refcardz.json
