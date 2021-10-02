
Linux requirements Docker: BuildKit (https://docs.docker.com/develop/develop-images/build_enhancements/#to-enable-buildkit-builds)
COMPOSE_DOCKER_CLI_BUILD=1
DOCKER_BUILDKIT=1

# Running the full blockchain locally

1. Navigate to repo root/ops folder
2. Run ./local/run_blockchain.sh (WSL )or sudo sh run_blockchain.sh (Non root user on Linux)

# Cleanup

1. Navigate to repo root/ops folder
2. Run ./clear_blockchain.sh (WSL )or sudo sh clear_blockchain.sh (Non root user on Linux)
