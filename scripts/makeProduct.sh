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

	#env |grep CORE
}

chaincodeInvokeMakeProduct() {
	PEER=$1
	ORG=$2
    ORDERID=$3
    ID=$4
    PRODUCTID=$5
    AMOUNT=$6
    PRICE=$7
    DETAIL=$8
    SUB=$9
    CODE={\"Args\":[\"setProduct\",\"$3\",\"$4\",\"$5\",\"$6\",\"$7\",\"$8\",\"$9\"]}
    #echo $CODE
	setGlobals $PEER $ORG
	if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
		#peer chaincode invoke -o orderer.supply.com:7050 -C $CHANNEL_NAME -n supplycc -c '{"Args":["setProduct", "fac1_pd1", "fac1", "pd1", "3", "5", "created", "NULL"]}' >&log.txt
        peer chaincode invoke -o orderer.supply.com:7050 -C $CHANNEL_NAME -n supplycc -c $CODE >&log.txt
        #echo "peer chaincode invoke -o orderer.supply.com:7050 -C $CHANNEL_NAME -n supplycc -c $CODE"
	else
		#peer chaincode invoke -o orderer.supply.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n supplycc -c '{"Args":["setProduct", "fac1_pd1", "fac1", "pd1", "3", "5", "created", "NULL"]}' >&log.txt
        peer chaincode invoke -o orderer.supply.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n supplycc -c $CODE >&log.txt
        #echo "peer chaincode invoke -o orderer.supply.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n supplycc -c $CODE"
	fi
	res=$?
	cat log.txt
	verifyResult $res "Invoke:setProductexecution on PEER$PEER failed "
	echo "Invoke:setProduct transaction on PEER $PEER on channel '$CHANNEL_NAME' is successful. "
}

chaincodeInvokeMakeProduct $1 $2 $3 $4 $5 $6 $7 $8 $9
sleep 3