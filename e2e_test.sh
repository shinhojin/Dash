#!/bin/bash
echo "Dankook Univ, 2019-2 Capstone Design Project"
echo "           ____             __   "
echo "          / __ \____ ______/ /_  "
echo '         / / / / __ `/ ___/ __ \ '
echo "        / /_/ / /_/ (__  ) / / / "
echo "       /_____/\__,_/____/_/ /_/  "
echo "                           "
echo "========================================="
echo " ____    _____      _      ____    _____ "
echo "/ ___|  |_   _|    / \    |  _ \  |_   _|"
echo "\___ \    | |     / _ \   | |_) |   | |  "
echo " ___) |   | |    / ___ \  |  _ <    | |  "
echo "|____/    |_|   /_/   \_\ |_| \_\   |_|  "
echo "========================================="
echo "      Build Supply Chain Network         "
echo "========================================="

export PATH=${PWD}/bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}
export COMPOSE_PROJECT_NAME="dash"
export CHANNEL_NAME=supplychannel
DELAY="3"
LANGUAGE="golang"
TIMEOUT=20
LANGUAGE=`echo "$LANGUAGE" | tr [:upper:] [:lower:]`

echo "Channel name : "$CHANNEL_NAME

CLI_TIMEOUT=100
CLI_DELAY=3
CHANNEL_NAME="supplychannel"
COMPOSE_FILE=docker-compose-cli.yaml
LANGUAGE=golang
echo "================================"
echo "            set up       "
echo "================================"
./init.sh
echo "================================"
echo "      docker container up       "
echo "================================"
#CHANNEL_NAME=$CHANNEL_NAME TIMEOUT=$CLI_TIMEOUT DELAY=$CLI_DELAY docker-compose -f $COMPOSE_FILE up -d
docker-compose -f ./base/docker-compose-couch.yaml up -d
docker-compose -f $COMPOSE_FILE up -d
echo ""
echo "======================================="
echo "  waiting for booting all containers"
echo "======================================="
echo ""
# from now process will start in containers
docker exec cli /bin/bash -c "sleep 20;./scripts/generate.sh"


echo ""
echo "======================================="
echo "  Ready for Invoke test for chaincode"
echo "======================================="
echo ""
sleep 10
echo ""
echo "======================================="
echo "   start Invoke test for chaincode"
echo "======================================="
echo ""
docker exec cli /bin/bash -c "./scripts/e2etest.sh"
echo ""
echo "======================================="
echo "         Remove Containers"
echo "======================================="
echo ""
./stop.sh

echo
echo "========= All GOOD, Supplychain execution completed =========== "
echo

echo
echo " _____   _   _   ____   "
echo "| ____| | \ | | |  _ \  "
echo "|  _|   |  \| | | | | | "
echo "| |___  | |\  | | |_| | "
echo "|_____| |_| \_| |____/  "
echo

exit 0
