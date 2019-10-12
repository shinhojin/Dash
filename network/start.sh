#! /bin/bash

echo "starting mq container"
docker-compose -f docker-compose-mq.yaml up -d

echo "starting orderer container"
docker-compose -f docker-compose-orderer.yaml up -d

echo "starting couchdb container"
docker-compose -f docker-compose-couch.yaml up -d

echo "starting peer container"
docker-compose -f docker-compose-peer.yaml up -d

echo "starting cli container"
docker-compose -f docker-compose-cli.yaml up -d

export CHANNEL_NAME=channel.first
export ARTIFACTS_DIR=/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts
docker exec cli peer channel create -o orderer0.example.com:7050 -c $CHANNEL_NAME -f $ARTIFACTS_DIR/channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
docker exec cli peer channel join -b $CHANNEL_NAME.block
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer channel join -b $CHANNEL_NAME.block 
docker exec -e "CORE_PEER_LOCALMSPID=Org3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org3.example.com:7051" cli peer channel join -b $CHANNEL_NAME.block