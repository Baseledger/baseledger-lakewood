export BLOCKCHAIN_APP_API_PORT=1317 && 
export TENDERMINT_NODE_GRPC_PORT=26657 && 
export TENDERMINT_NODE_PORT=26655

docker-compose -p baseledger-node up -d