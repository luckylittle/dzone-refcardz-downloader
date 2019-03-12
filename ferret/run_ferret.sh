#!/bin/bash

# echo 'Running dzone-refcardz.fql via headless Chrome'
# ferret -time --cdp http://127.0.0.1:9222 dzone-refcardz.fql

ferret --cdp http://127.0.0.1:9222 dzone-refcardz.fql > dzone-refcardz.json
cat dzone-refcardz.json | jq 'sort_by(.name)' > dzone-refcardz$(date +"%Y%m%d").json
rm dzone-refcardz.json
