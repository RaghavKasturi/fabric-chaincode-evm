mkdir output
cp -rp ../evmcc/*.go ./output
cp -rp ../evmcc/go.mod ./output
cp -rp ../evmcc/vendor ./output

mkdir -p ./output/vendor/github.com/hyperledger/fabric-chaincode-evm/evmcc
cp -rp ../evmcc/address ./output/vendor/github.com/hyperledger/fabric-chaincode-evm/evmcc
cp -rp ../evmcc/event ./output/vendor/github.com/hyperledger/fabric-chaincode-evm/evmcc
cp -rp ../evmcc/eventmanager ./output/vendor/github.com/hyperledger/fabric-chaincode-evm/evmcc
cp -rp ../evmcc/mocks ./output/vendor/github.com/hyperledger/fabric-chaincode-evm/evmcc
cp -rp ../evmcc/statemanager ./output/vendor/github.com/hyperledger/fabric-chaincode-evm/evmcc

cd output
zip -r evmcc.zip *.go go.mod vendor/
mv evmcc.zip ..
cd ..
rm -rf output