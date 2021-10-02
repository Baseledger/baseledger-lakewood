export POSTGRES_EXPOSED_PORT=5432 && 
export NATS_EXPOSED_PORT=4222 && 
export BLOCKCHAIN_APP_API_PORT=1317 && 
export TENDERMINT_NODE_GRPC_PORT=26657 && 
export TENDERMINT_NODE_PORT=26655 && 
export PROXY_APP_PORT=8081 &&
export ORGANIZATION_ID=<unique identifier of the organization>
export API_CONCIRCLE_URL=<system of record endpoint>

docker-compose -p baseledger-node up -d