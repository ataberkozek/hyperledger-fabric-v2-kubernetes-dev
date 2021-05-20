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

queryAllRes() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
peer chaincode query --name ${CCNAME} \
--channelID ${CHANNEL_ID} \
--ctor '{"Args":["QueryAllRes"]}' \
--tls --cafile \$ORDERER_TLS_ROOTCERT_FILE
EOF
}

queryCar() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
peer chaincode query --name ${CCNAME} \
--channelID ${CHANNEL_ID} \
--ctor '{"Args":["QueryCar", "CAR10"]}' \
--tls --cafile \$ORDERER_TLS_ROOTCERT_FILE
EOF
}

queryRe() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
peer chaincode query --name ${CCNAME} \
--channelID ${CHANNEL_ID} \
--ctor '{"Args":["QueryRe", "RE1"]}' \
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
  --ctor '{"Args":["ChangeCarOwner", "CAR0", "Ataberk"]}' \
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

changeReOwner() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
echo "Submitting invoketransaction to smart contract on ${CHANNEL_ID}"
peer chaincode invoke \
  --channelID ${CHANNEL_ID} \
  --name ${CCNAME} \
  --ctor '{"Args":["ChangeReOwner", "RE1", "Ataberk"]}' \
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

changeRePrice() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
echo "Submitting invoketransaction to smart contract on ${CHANNEL_ID}"
peer chaincode invoke \
  --channelID ${CHANNEL_ID} \
  --name ${CCNAME} \
  --ctor '{"Args":["ChangeRePrice", "RE1", "$200,000"]}' \
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
  --ctor '{"Args":["CreateCar", "CAR12", "Mercedes", "SLS-AMG", "blue", "Batuhan"]}' \
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

createRe() {
CCNAME=$1
CHANNEL_ID=$2
cat <<EOF
echo "Submitting invoketransaction to smart contract on ${CHANNEL_ID}"
peer chaincode invoke \
  --channelID ${CHANNEL_ID} \
  --name ${CCNAME} \
  --ctor '{"Args":["AddRe", "RE10", "Istanbul", "4", "1", "$110,000", "130m2"]}' \
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
