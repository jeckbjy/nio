#!/bin/bash
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}")" && pwd )
ROOT=$DIR/../..
#FILE=server-linux

echo 'Step1:Build for ubuntu'
#cd $ROOT/src/fire

GOOS=linux GOARCH=amd64 go build -o ./echo_linux .

echo 'Step3:Sync to server'
# rsync -avcuR ${EXCLUDE} server ${SERVER}:${TARGET_DIR}/${TARGET_FOLDER_NAME} || exit 1
#cd $DIR/context
HOST=fire
DEST=/home/jeckbjy/app/
# cd ./context
# rsync -avcuR --progress ./* fire:~/app/
rsync -avcuR --progress ./echo_linux ${HOST}:${DEST}

#echo 'Step4:Start server'
#ssh $HOST "cd ~/app/ && sh ./start.sh"
# ssh $HOST 'bash -l -c "~/app/start.sh"'
# ssh $HOST 'bash -c "pkill server | ~/app/server &"'
# ssh $HOST "pkill server | ~/app/server &&"
#echo 'Step5:Clean'
#cd $DIR
#rm -rf ./context