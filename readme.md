# baseledger-lakewood
_In the spirit of developing a proof of concept implementation to experiment with network validation in tendermint (including staking and delegation), native opcodes and a community block explorer, the team built the ["lakewood" testnet](https://github.com/baseledger/lakewood). This testnet was created using Cosmos SDK._

## Explore

To explore the current state of baseledger-lakewood, open the [Lakewood Baseledger Explorer](https://lakewood.baseledger.net).

## Join lakewood as replicator node and drop proofs

To join lakewood as a replicator node, perform the steps from https://docs.baseledger.net/howtos-1/how-to-drop-a-proof-on-baseledger-lakewood/drop-a-proof-via-a-local-node

## Drop proofs via preinstalled nodes

Perform the steps from https://docs.baseledger.net/howtos-1/how-to-drop-a-proof-on-baseledger-lakewood/drop-a-proof-via-preinstalled-nodes

## Install your own block explorer for lakewood

1. Install prerequisites - NodeJS LTS and meteor
2. Clone the repo to your server
3. Navigate to explorer folder
4. Create/modify settings.json
5. Run *npm install --save* 
6. Run *meteor --settings settings.json*

## Join lakewood as a validator node

As the Lakewood testnet may be frequently reset, it is currently not possible to become an external user or validator on lakewood.
