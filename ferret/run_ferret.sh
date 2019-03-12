#!/bin/bash

echo 'Running dzone-refcardz.fql via headless Chrome'
ferret -time --cdp http://127.0.0.1:9222 dzone-refcardz.fql
