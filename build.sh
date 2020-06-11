#!/bin/bash -xe

BASEDIR=$(dirname "$0")
pushd $BASEDIR
SCRIPT_DIR=`pwd`
popd

export GO111MODULE=on
unameOut="$(uname -s)"
case "${unameOut}" in
Linux*) MACHINE=linux
      ;;
Darwin*) MACHINE=darwin
      ;;
CYGWIN*) MACHINE=Cygwin
      ;;
MINGW*) MACHINE=MinGw
      ;;
*) MACHINE="UNKNOWN:${unameOut}" ;;
esac

GOOS=$1
GOARCH=$2

if [ -z $GOOS ]; then
  GO_OS=$MACHINE
fi

if [ -z $GOARCH ]; then
  GOARCH=amd64
fi

echo "OS=${GOOS}, ARCH=${GOARCH}"

export GOOS=${GOOS}
export GOARCH=${GOARCH}
export GOARM=5


CONTROLLER_DIR=${SCRIPT_DIR}/controller
CONTROLLER_GENERATED_DIR=${CONTROLLER_DIR}/generated

CONTROLLER_MAIN_DIR=${CONTROLLER_DIR}/cmd

CONTROLLER_MAIN_FILE=${CONTROLLER_DIR}/cli/cmd/main.go

CONTROLLER_UI_DIR=${CONTROLLER_DIR}/ui

CONTROLLER_SWAGGERUI_DIR=${CONTROLLER_UI_DIR}/swaggerui

if [ -d "${CONTROLLER_GENERATED_DIR}" ]; then
  rm -rf "${CONTROLLER_GENERATED_DIR}"
fi


mkdir -p ${CONTROLLER_GENERATED_DIR}

BIN_OUT_DIR=${SCRIPT_DIR}/build
mkdir -p ${BIN_OUT_DIR}

rm -rf ${BIN_OUT_DIR}/controller

go get github.com/rakyll/statik

# Generate controller swagger client code
swagger generate client -t ${CONTROLLER_GENERATED_DIR} -f ${CONTROLLER_DIR}/controller-swagger.json -A controller

swagger generate server -t ${CONTROLLER_GENERATED_DIR} -f ${CONTROLLER_DIR}/controller-swagger.json -A controller

GOOS=$MACHINE GOARCH=amd64 go run ${SCRIPT_DIR}/tool/cmd/main.go replace-middleware -f ${CONTROLLER_GENERATED_DIR}/restapi/configure_controller.go

cp -rf ${CONTROLLER_DIR}/controller-swagger.json ${CONTROLLER_SWAGGERUI_DIR}/swagger.json

statik -src=$CONTROLLER_UI_DIR -dest=$CONTROLLER_GENERATED_DIR

# Compile controller CLI
go build -ldflags '-w -s'  -o ${BIN_OUT_DIR}/controller ${CONTROLLER_MAIN_FILE}

# For linux arm (RPI)
if [[ "$GOOS" == "linux" ]] && [[ "$GOARCH" == "arm" ]]; then
 scp ${BIN_OUT_DIR}/controller pi@192.168.1.20:/home/pi/controller
 ssh pi@192.168.1.20 'sudo systemctl stop controlboard.service && sudo cp /home/pi/controller /controlboard/bin/ && sudo systemctl start controlboard.service'
fi