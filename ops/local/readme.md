
Linux requirements Docker: BuildKit (https://docs.docker.com/develop/develop-images/build_enhancements/#to-enable-buildkit-builds)
COMPOSE_DOCKER_CLI_BUILD=1
DOCKER_BUILDKIT=1

# Running the full blockchain locally

1. Navigate to repo root/ops/local folder
2. Run ./run_blockchain.sh (WSL )or sudo sh run_blockchain.sh (Non root user on Linux)

# Cleanup

1. Navigate to repo root/ops/local folder
2. Run ./clear_blockchain.sh (WSL )or sudo sh clear_blockchain.sh (Non root user on Linux)

# Postman collection

After setting up a local blockchain, you can play around with requests between two proxies by following the steps:

1. Import proxy_app/misc/baseledger demo proxy.postman_collection
2. make sure endpoints have correct user and pass for basic auth previously set, and are set to target localhost
3. make sure workgroup, organization id and organization nats endpoints have correct values set (as provided in docker compose)