#!/bin/sh

# crontab example
# * 5 * * * sh path/to/project/cleaner.sh

# project directory
DIR=$(cd $(dirname $0)/.. && pwd)
cd ${DIR}

# remove images more than 180 days.
find ./images/ -mtime +180 -exec rm -f {} \;
