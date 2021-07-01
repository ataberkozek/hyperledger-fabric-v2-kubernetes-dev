packageAndInstall() {
CCURL=$1
CCNAME=$2
# CCURL="github.com/blockchaind/hyperledger-fabric-v2-kubernetes-dev/key-value-chaincode"
# CCNAME="keyval"
LANG=golang
LABEL=${CCNAME}v1
cat <<EOF
echo "getting chaincode for ${ORG} ${PEER}"
go get -d -x ${CCURL}
echo "packaging chaincode for ${ORG} ${PEER}"
peer lifecycle chaincode package ${CCNAME}.tar.gz --path ${CCURL} --lang ${LANG} --label ${LABEL}
echo "installing chaincode on ${ORG} ${PEER}"
peer lifecycle chaincode install ${CCNAME}.tar.gz
EOF
}

approve() {
CCNAME=$1
CHANNEL_ID=$2
LABEL=${CCNAME}v1

cat <<EOF
PACKAGE_ID=\$(peer lifecycle chaincode queryinstalled | awk '/${LABEL}/ {print substr(\$3, 1, length(\$3)-1)}')
echo "Package ID: \${PACKAGE_ID}"
peer lifecycle chaincode approveformyorg --package-id \${PACKAGE_ID} \
  --signature-policy "AND('Org1MSP.member','Org2MSP.member','Org3MSP.member')" \
  -C ${CHANNEL_ID} -n ${CCNAME} -v 1.0  --sequence 1 \
  --tls true --cafile \$ORDERER_TLS_ROOTCERT_FILE --waitForEvent
EOF
}

checkCommitReadiness() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
peer lifecycle chaincode checkcommitreadiness \
--name ${CCNAME} --channelID ${CHANNEL_ID} \
--signature-policy "AND('Org1MSP.member', 'Org2MSP.member', 'Org3MSP.member')" \
--version 1.0 --sequence 1
EOF
}

commit() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
echo "Committing smart contract"
peer lifecycle chaincode commit \
  --channelID ${CHANNEL_ID} \
  --name ${CCNAME} \
  --version 1.0 \
  --signature-policy "AND('Org1MSP.member','Org2MSP.member', 'Org3MSP.member')" \
  --sequence 1 --waitForEvent \
  --peerAddresses peer0.org1:7051 \
  --peerAddresses peer0.org2:7051 \
  --peerAddresses peer0.org3:7051  \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org1-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org2-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org3-cert.pem  \
  --tls true --cafile \$ORDERER_TLS_ROOTCERT_FILE \
  -o orderer.org1:7050
EOF
}

init() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
echo "Submitting invoketransaction to smart contract on ${CHANNEL_ID}"
peer chaincode invoke \
  --channelID ${CHANNEL_ID} \
  --name ${CCNAME} \
  --ctor '{"Args":["InitLedger"]}' \
  --waitForEvent \
  --waitForEventTimeout 300s \
  --cafile \$ORDERER_TLS_ROOTCERT_FILE \
  --tls true -o orderer.org1:7050 \
  --peerAddresses peer0.org1:7051 \
  --peerAddresses peer0.org2:7051 \
  --peerAddresses peer0.org3:7051  \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org1-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org2-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org3-cert.pem
EOF
}

queryAllCars() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
peer chaincode query --name ${CCNAME} \
--channelID ${CHANNEL_ID} \
--ctor '{"Args":["QueryAllCars"]}' \
--tls --cafile \$ORDERER_TLS_ROOTCERT_FILE
EOF
}

queryCar() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
peer chaincode query --name ${CCNAME} \
--channelID ${CHANNEL_ID} \
--ctor '{"Args":["QueryCar", "CAR1"]}' \
--tls --cafile \$ORDERER_TLS_ROOTCERT_FILE
EOF
}

changeCarOwner() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
echo "Submitting invoketransaction to smart contract on ${CHANNEL_ID}"
peer chaincode invoke \
  --channelID ${CHANNEL_ID} \
  --name ${CCNAME} \
  --ctor '{"Args":["ChangeCarOwner", "CAR1", "Robert"]}' \
  --waitForEvent \
  --waitForEventTimeout 300s \
  --cafile \$ORDERER_TLS_ROOTCERT_FILE \
  --tls true -o orderer.org1:7050 \
  --peerAddresses peer0.org1:7051 \
  --peerAddresses peer0.org2:7051 \
  --peerAddresses peer0.org3:7051  \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org1-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org2-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org3-cert.pem 
EOF
}

changeCarPrice() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
echo "Submitting invoketransaction to smart contract on ${CHANNEL_ID}"
peer chaincode invoke \
  --channelID ${CHANNEL_ID} \
  --name ${CCNAME} \
  --ctor '{"Args":["ChangeCarPrice", "CAR1", "200000"]}' \
  --waitForEvent \
  --waitForEventTimeout 300s \
  --cafile \$ORDERER_TLS_ROOTCERT_FILE \
  --tls true -o orderer.org1:7050 \
  --peerAddresses peer0.org1:7051 \
  --peerAddresses peer0.org2:7051 \
  --peerAddresses peer0.org3:7051  \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org1-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org2-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org3-cert.pem 
EOF
}

createCar() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
echo "Submitting invoketransaction to smart contract on ${CHANNEL_ID}"
peer chaincode invoke \
  --channelID ${CHANNEL_ID} \
  --name ${CCNAME} \
  --ctor '{"Args":["CreateCar", "CAR10", "Volkswagen", "Polo", "white", "Alice"]}' \
  --waitForEvent \
  --waitForEventTimeout 300s \
  --cafile \$ORDERER_TLS_ROOTCERT_FILE \
  --tls true -o orderer.org1:7050 \
  --peerAddresses peer0.org1:7051 \
  --peerAddresses peer0.org2:7051 \
  --peerAddresses peer0.org3:7051  \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org1-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org2-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org3-cert.pem 
EOF
}

createParty() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
echo "Submitting invoketransaction to smart contract on ${CHANNEL_ID}"
peer chaincode invoke \
  --channelID ${CHANNEL_ID} \
  --name ${CCNAME} \
  --ctor '{"Args":["CreateParty", "Greens"]}' \
  --waitForEvent \
  --waitForEventTimeout 300s \
  --cafile \$ORDERER_TLS_ROOTCERT_FILE \
  --tls true -o orderer.org1:7050 \
  --peerAddresses peer0.org1:7051 \
  --peerAddresses peer0.org2:7051 \
  --peerAddresses peer0.org3:7051  \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org1-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org2-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org3-cert.pem 
EOF
}

castVoteDemocrats() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
echo "Submitting invoketransaction to smart contract on ${CHANNEL_ID}"
peer chaincode invoke \
  --channelID ${CHANNEL_ID} \
  --name ${CCNAME} \
  --ctor '{"Args":["CastVote", "Democrats"]}' \
  --waitForEvent \
  --waitForEventTimeout 300s \
  --cafile \$ORDERER_TLS_ROOTCERT_FILE \
  --tls true -o orderer.org1:7050 \
  --peerAddresses peer0.org1:7051 \
  --peerAddresses peer0.org2:7051 \
  --peerAddresses peer0.org3:7051  \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org1-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org2-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org3-cert.pem 
EOF
}

castVoteRepublicans() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
echo "Submitting invoketransaction to smart contract on ${CHANNEL_ID}"
peer chaincode invoke \
  --channelID ${CHANNEL_ID} \
  --name ${CCNAME} \
  --ctor '{"Args":["CastVote", "Republicans"]}' \
  --waitForEvent \
  --waitForEventTimeout 300s \
  --cafile \$ORDERER_TLS_ROOTCERT_FILE \
  --tls true -o orderer.org1:7050 \
  --peerAddresses peer0.org1:7051 \
  --peerAddresses peer0.org2:7051 \
  --peerAddresses peer0.org3:7051  \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org1-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org2-cert.pem \
  --tlsRootCertFiles /etc/hyperledger/fabric-peer/client-root-tlscas/tlsca.org3-cert.pem 
EOF
}

queryAllPartys() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
peer chaincode query --name ${CCNAME} \
--channelID ${CHANNEL_ID} \
--ctor '{"Args":["GetAllPartys"]}' \
--tls --cafile \$ORDERER_TLS_ROOTCERT_FILE
EOF
}
