# Join lakewood as replicator node

1. Copy ops folder to the server
3. Run *setup-node-infrastructure-for-testnet.sh*
4. Open only ports in the firewall that are necessary for external communication, 26656. Be aware of docker\ufw incompatibility https://github.com/docker/for-linux/issues/777
5. Open *add-node-to-running-blokchain.sh* and add your prefered node name instead of <node_id>
6. Run *add-node-to-running-blokchain.sh*
6. docker exec baseledger-node_blockchain_app_1 baseledgerd start

# Setup baseledger lakewood explorer

1. Install prerequisites - NodeJS LTS and meteor
2. Clone the repo
2. Create/modify settings.json
3. Run *npm install --save* 
4. Run *meteor --settings settings.json*