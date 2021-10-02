export POSTGRES_EXPOSED_PORT=5432 && export NATS_EXPOSED_PORT=4222 && export BLOCKCHAIN_APP_API_PORT=1317 && export TENDERMINT_NODE_GRPC_PORT=26657 && export TENDERMINT_NODE_PORT=26656 && export PROXY_APP_PORT=8081
docker-compose -p first_node down -v
docker-compose -p second_node down -v