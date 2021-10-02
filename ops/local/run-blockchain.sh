# Docker prereqs
export COMPOSE_DOCKER_CLI_BUILD=1 && export DOCKER_BUILDKIT=1

first_node_tendermint_grpc_port=26657
second_node_tendermint_grpc_port=26658

first_node_tendermint_p2p_port=26655
second_node_tendermint_p2p_port=26656

tendermint_internal_grpc_port=26657

# Sets up environment for first node and runs it
export POSTGRES_EXPOSED_PORT=5432 && 
export NATS_EXPOSED_PORT=4222 && 
export BLOCKCHAIN_APP_API_PORT=1317 && 
export TENDERMINT_NODE_GRPC_PORT=$first_node_tendermint_grpc_port && 
export TENDERMINT_NODE_PORT=$first_node_tendermint_p2p_port && 
export PROXY_APP_PORT=8081 &&
export ORGANIZATION_ID=d9c102eb-b173-45ac-b640-24d28a3c9f0c && # unique identifier of the organization, currently hardcoded in seed data
export API_CONCIRCLE_URL=s4h.rp.concircle.com
docker-compose -p first_node up -d

# Sets up environment for second node and runs it
export POSTGRES_EXPOSED_PORT=5433 && 
export NATS_EXPOSED_PORT=4223 && 
export BLOCKCHAIN_APP_API_PORT=1318 && 
export TENDERMINT_NODE_GRPC_PORT=$second_node_tendermint_grpc_port && 
export TENDERMINT_NODE_PORT=$second_node_tendermint_p2p_port && 
export PROXY_APP_PORT=8082 &&
export ORGANIZATION_ID=969e989c-bb61-4180-928c-0d48afd8c6a3 && # unique identifier of the organization, currently hardcoded in seed data
export API_CONCIRCLE_URL=s4p.rp.concircle.com

docker-compose -p second_node up -d

# Initialize first node
docker exec first_node_blockchain_app_1 baseledgerd init node1 --chain-id baseledger

# Initialize first node validator account
docker exec first_node_blockchain_app_1 baseledgerd keys add node1_validator --keyring-backend test

# Add first node validator account as genesis
node1_validator_address=$(docker exec first_node_blockchain_app_1 baseledgerd keys show node1_validator -a --keyring-backend test)
docker exec first_node_blockchain_app_1 baseledgerd add-genesis-account ${node1_validator_address} 100000000000stake,100000000000token

# Set up faucet account
# docker exec -ti first_node_blockchain_app_1 sh
# baseledgerd keys add faucet --recover=true --keyring-backend test
# #enter mnemonic
# exit

# Add faucet account as genesis
docker exec first_node_blockchain_app_1 baseledgerd add-genesis-account baseledger1ecskwhmn9l0j7zt98lfps3dtv4y4fuq8lh67ls 100000000000stake,100000000000token

# Initialize second node validator account
docker exec second_node_blockchain_app_1 baseledgerd keys add node2_validator --keyring-backend test

# Add second node validator account as genesis
node2_validator_address=$(docker exec second_node_blockchain_app_1 baseledgerd keys show node2_validator -a --keyring-backend test)
docker exec first_node_blockchain_app_1 baseledgerd add-genesis-account ${node2_validator_address} 100000000000stake,100000000000token

# Copy genesis from first to second node via host machine for gentx generation
docker cp first_node_blockchain_app_1:/root/.baseledger/config/genesis.json .
docker cp ./genesis.json second_node_blockchain_app_1:/root/.baseledger/config/genesis.json

# Generate genensis transaction on the second node
docker exec second_node_blockchain_app_1 baseledgerd gentx node2_validator 100000000stake --chain-id baseledger --keyring-backend test

# Copy gentx from second to first node via host machine for genesis collection
docker cp second_node_blockchain_app_1:/root/.baseledger/config/gentx .
docker cp ./gentx first_node_blockchain_app_1:/root/.baseledger/config

# Generate genensis transaction on the first node
docker exec first_node_blockchain_app_1 baseledgerd gentx node1_validator 100000000stake --chain-id baseledger --keyring-backend test

# Collect both genesis transactions
docker exec first_node_blockchain_app_1 baseledgerd collect-gentxs

# Copy fully formed genesis from first to second node via host machine
docker cp first_node_blockchain_app_1:/root/.baseledger/config/genesis.json . 
docker cp ./genesis.json second_node_blockchain_app_1:/root/.baseledger/config/genesis.json

# Setup config toml
node1_id=$(docker exec first_node_blockchain_app_1 baseledgerd tendermint show-node-id)
node2_id=$(docker exec second_node_blockchain_app_1 baseledgerd tendermint show-node-id)

internal_host_ip=$(docker exec first_node_blockchain_app_1 getent hosts host.docker.internal | awk '{print $1}')

# this adds peers to each other
docker exec first_node_blockchain_app_1 sed -i 's/persistent_peers = ".*/persistent_peers = "'${node2_id}'@'${internal_host_ip}':'${second_node_tendermint_p2p_port}'"/' ~/.baseledger/config/config.toml
docker exec second_node_blockchain_app_1 sed -i 's/persistent_peers = ""/persistent_peers = "'${node1_id}'@'${internal_host_ip}':'${first_node_tendermint_p2p_port}'"/' ~/.baseledger/config/config.toml

# this enables grpc
docker exec first_node_blockchain_app_1 sed -i 's@laddr = "tcp://127.0.0.1:'${tendermint_internal_grpc_port}'"@laddr = "tcp://0.0.0.0:'${tendermint_internal_grpc_port}'"@' ~/.baseledger/config/config.toml
docker exec second_node_blockchain_app_1 sed -i 's@laddr = "tcp://127.0.0.1:'${tendermint_internal_grpc_port}'"@laddr = "tcp://0.0.0.0:'${tendermint_internal_grpc_port}'"@' ~/.baseledger/config/config.toml

# this enables rest api, it is only enable = false entry, maybe we can make it a bit more precise?
docker exec first_node_blockchain_app_1 sed -i 's/enable = false/enable = true/' ~/.baseledger/config/app.toml
docker exec second_node_blockchain_app_1 sed -i 's/enable = false/enable = true/' ~/.baseledger/config/app.toml

# this allows connecting peers not in the address book
docker exec first_node_blockchain_app_1 sed -i 's/addr_book_strict = true/addr_book_strict = false/' ~/.baseledger/config/config.toml
docker exec second_node_blockchain_app_1 sed -i 's/addr_book_strict = true/addr_book_strict = false/' ~/.baseledger/config/config.toml

# This allows connections from localhost to tendermint API
docker exec first_node_blockchain_app_1 sed -i 's/allow_duplicate_ip = false/allow_duplicate_ip = true/' ~/.baseledger/config/config.toml
docker exec second_node_blockchain_app_1 sed -i 's/allow_duplicate_ip = false/allow_duplicate_ip = true/' ~/.baseledger/config/config.toml

# This increases the timeout between blocks to 30s
# docker exec first_node_blockchain_app_1 sed -i 's/timeout_commit = "5s"/timeout_commit = "30s"/' ~/.baseledger/config/config.toml
# docker exec second_node_blockchain_app_1 sed -i 's/timeout_commit = "5s"/timeout_commit = "30s"/' ~/.baseledger/config/config.toml


# start first node - TODO: Has to  be executed in a separate window after running this script in order to have logs
# docker exec first_node_blockchain_app_1 baseledgerd start
# node2_adress = docker exec second_node_blockchain_app_1 baseledgerd keys show node2_validator -a
# docker exec first_node_blockchain_app_1 baseledgerd tx bank send node1_validator node2_adress 1000token --yes

# start second node - TODO: Has to  be executed in a separate window after running this script in order to have logs
# docker exec second_node_blockchain_app_1 baseledgerd start

# cleanup

rm ./genesis.json
rm -rf ./gentx

