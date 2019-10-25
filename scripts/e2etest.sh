#!/bin/bash
COMPOSE_PROJECT_NAME=dash_supply
export CHANNEL_NAME=supplychannel
LANGUAGE="golang"
DELAY=3
TIMEOUT=20
LANGUAGE=`echo "$LANGUAGE" | tr [:upper:] [:lower:]`
COUNTER=1
MAX_RETRY=5

ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/supply.com/orderers/orderer.supply.com/msp/tlscacerts/tlsca.supply.com-cert.pem
CC_SRC_PATH="github.com/chaincode/supplyContract"
# verify the result of test
verifyResult () {
	if [ $1 -ne 0 ] ; then
		echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
    echo "========= ERROR !!! FAILED to execute generate.sh ==========="
		echo
   		exit 1
	fi
}

# PEER0 for consumer, PEER1 for retailer
# PEER2 for logistics, PEER3 for wholesaler
# PEER4 for manufacture processor, PEER5 for product producer
setGlobals () {
	ORG=$2
	if [ $ORG -eq 1 ] ; then
		CORE_PEER_LOCALMSPID="Org1MSP"
		CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.supply.com/peers/peer0.org1.supply.com/tls/ca.crt
		CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.supply.com/users/Admin@org1.supply.com/msp
		#consumer and retailer on shopping site 
		if [ $1 -eq 0 ]; then
			#consumer
			CORE_PEER_ADDRESS=peer0.org1.supply.com:7051
			PEER=PEER0
		else
			#retailer
			CORE_PEER_ADDRESS=peer1.org1.supply.com:7051
			PEER=PEER1
		fi
	elif [ $ORG -eq 2 ] ; then
		CORE_PEER_LOCALMSPID="Org2MSP"
		CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.supply.com/peers/peer0.org2.supply.com/tls/ca.crt
		CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.supply.com/users/Admin@org2.supply.com/msp
		#logistics and wholesaler
		if [ $1 -eq 0 ]; then
			#logistics
			CORE_PEER_ADDRESS=peer0.org2.supply.com:7051
			PEER=PEER2
		else
			#wholesaler
			CORE_PEER_ADDRESS=peer1.org2.supply.com:7051
			PEER=PEER3
		fi
	elif [ $ORG -eq 3 ] ; then
		CORE_PEER_LOCALMSPID="Org3MSP"
		CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.supply.com/peers/peer0.org3.supply.com/tls/ca.crt
		CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.supply.com/users/Admin@org3.supply.com/msp
		#manufacture processor and raw food producer
		if [ $1 -eq 0 ]; then
			#manufacture processor
			CORE_PEER_ADDRESS=peer0.org3.supply.com:7051
			PEER=PEER4
		else
			#raw food producer
			CORE_PEER_ADDRESS=peer1.org3.supply.com:7051
			PEER=PEER5
		fi		
	else
		echo "================== ERROR !!! ORG OR PEER Unknown =================="
	fi

	env |grep CORE
}

setGlobals 0 1

peer chaincode invoke -o orderer.supply.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n supplycc -c '{"Args":["setProduct", "fac1_pd1", "fac1", "pd1", "3", "5", "created", "NULL"]}' >&log.txt
peer chaincode invoke -o orderer.supply.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n supplycc -c '{"Args":["setProduct", "fac2_pd1", "fac2", "pd1", "1", "5", "created", "NULL"]}' >&log.txt
cat log.txt
echo "Invoke:setProduct transaction on PEER $PEER on channel '$CHANNEL_NAME' is successful. "
sleep 5
peer chaincode invoke -o orderer.supply.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n supplycc -c '{"Args":["useProduct", "fac1_pd1", "1"]}' >&log.txt
cat log.txt
echo "Invoke:useProduct transaction on PEER $PEER on channel '$CHANNEL_NAME' is successful. "
sleep 5
peer chaincode invoke -o orderer.supply.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n supplycc -c '{"Args":["moveProduct", "fac1_pd1", "fac2_pd1", "2"]}' >&log.txt
cat log.txt
echo "Invoke:moveProduct transaction on PEER $PEER on channel '$CHANNEL_NAME' is successful. "
sleep 5
peer chaincode invoke -o orderer.supply.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n supplycc -c '{"Args":["query", "fac1_pd1"]}' >&log.txt
cat log.txt
echo "Invoke:queary transaction on PEER $PEER on channel '$CHANNEL_NAME' is successful. "

./scripts/query.sh 0 1 fac2_pd1