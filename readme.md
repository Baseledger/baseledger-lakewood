# baseledger-lakewood
_In the spirit of developing a proof of concept implementation to experiment with network validation in tendermint (including staking and delegation), native opcodes and a community block explorer, the team built the ["lakewood" testnet](https://github.com/baseledger/lakewood). This testnet was created using Cosmos SDK._

## Explore

To explore the current state of baseledger-lakewood, open the [Lakewood Baseledger Explorer](https://lakewood.baseledger.net).

## Join lakewood as replicator node and drop proofs

To join lakewood as a replicator node, perform the steps from docs/dropping_proofs_via_local_node/dropping_proofs_via_local_node.md

## Drop proofs via predeloyed nodes

Perform the steps from docs/dropping_proofs_via_predeployed_nodes.md

## Install your own block explorer for lakewood

1. Install prerequisites - NodeJS LTS and meteor
2. Clone the repo to your server
3. Navigate to explorer folder
4. Create/modify settings.json
5. Run *npm install --save* 
6. Run *meteor --settings settings.json*

## Join lakewood as a validator node

As the Lakewood testnet may be frequently reset, it is currently not possible to become an external user or validator on lakewood.

## baseledger-core and Peachtree Testnet

The ["peachtree" testnet](https://explorer.peachtree.baseledger.net) was created from scratch using tendermint for BFT consensus and the [Provide stack](https://docs.provide.services) for subscribing to events emitted by the Baseledger governance and staking contracts, broadcasting _baseline proofs_ to the network and otherwise interacting with the [Baseline Protocol](https://github.com/eea-oasis/baseline). As a result of this design, [`baseledger-core`](https://github.com/Baseledger/baseledger-core) can be built as a single container and added to existing deployments of the Provide stack for increased security. `baseledger-core` can also run standalone (i.e., outside the context of a Provide stack). Baseledger nodes running outside the context of a Provide stack are not restricted from operating as validator, full or seed nodes. Organizations implementing the _baseline_ pattern in commercial multiparty workflows benefit from running a local Baseledger node because it provides additional security to the cryptographic commitments (proofs) stored within the Provide stack without sacrificing any privacy guarantees inherent to _baselining_.


