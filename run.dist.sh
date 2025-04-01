#!/usr/bin/env bash

Usage="Usage: run.sh -e [execute command] -t [target]"
exeCmd=""
targetName=""

while getopts "e:t:" opt
do
  case $opt in
    e)
      exeCmd=$OPTARG
      ;;
    t)
      targetName=$OPTARG
      ;;
    ?)
      echo $Usage
      exit 1
      ;;
  esac
done

shift $(($OPTIND - 1))
if [ -z "$exeCmd" ]; then
  echo $Usage
  exit 1
fi

cmd_start() {
  echo "executing command: start"
  GIN_MODE=release pm2 start ./server --name ${targetName}-uploader-prod -i 1
}

cmd_stop() {
  echo "executing command: stop"
  pm2 stop ${targetName}-uploader-prod
}

cmd_restart() {
  echo "executing command: restart"
  pm2 restart ${targetName}-uploader-prod
}

cmd_delete() {
  echo "executing command: delete"
  pm2 delete ${targetName}-uploader-prod
}

cmd_logcat() {
  echo "executing command: logcat"
  pm2 logs ${targetName}-uploader-prod
}

if [ "$exeCmd" = "start" ]; then
  cmd_start
elif [ "$exeCmd" = "stop" ]; then
  cmd_stop
elif [ "$exeCmd" = "restart" ]; then
  cmd_restart
elif [ "$exeCmd" = "delete" ]; then
  cmd_delete
elif [ "$exeCmd" = "logcat" ]; then
  cmd_logcat
else
  echo "unknown executing command: ${exeCmd}"
fi