#!/bin/bash
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}")" && pwd )
ROOT=$DIR/../..
FILE=gotest

echo 'Step1:Build for ubuntu'
GOOS=linux GOARCH=amd64 go build -o ./${FILE} .

echo 'Step3:Sync to server'
HOST=fire
DEST=/home/jeckbjy/app/

rsync -avcuR --progress ./${FILE} ${HOST}:${DEST}
