#!/bin/sh

APP_FILE=gyazo
PID_FILE=g.pid
DIR_PID=pid
DIR_BIN=bin

# project directory
DIR=$(cd $(dirname $0)/.. && pwd)
cd ${DIR}

# check directory
__check_dir() {
  if [ ! -d ${DIR}/$1 ]; then
      mkdir -p ${DIR}/$1
      echo "make $1 directory.."
  fi
}

__build() {
  __check_dir "${DIR_BIN}"

  # remove old application
  if [ -e ${DIR_BIN}/${APP_FILE} ]; then
    rm -f ${DIR_BIN}/${APP_FILE}
  fi

  # build application
  go build -o ${DIR_BIN}/${APP_FILE} src/main.go
  echo "application was built successfully."
}

__start() {
  __check_dir "${DIR_PID}"

  # check application
  if [ ! -e ${DIR_BIN}/${APP_FILE} ]; then
    __build
  fi

  # start application
  exec nohup gyazo > /tmp/gyazo.out 2>&1&
  echo $! > ${DIR_PID}/${PID_FILE}
  disown

  echo "gyazo service started."
}

__stop() {
  if [ -e ${DIR_PID}/${PID_FILE} ]; then
    kill `cat ${DIR_PID}/${PID_FILE}`
    echo "gyazo service stopped."
  fi
}

# call function
case $1 in
  build)
    __build
    ;;
  start)
    __start
    ;;
  stop)
    __stop
    ;;
  restart)
    __stop
    __start
    ;;
  *)
  echo "usage: run.sh {build|start|stop|restart}" ;;
esac
exit 0
