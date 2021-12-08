# Drop a proof on baseledger lakewood via local node

Highlevel overview of the steps needed:

1. Run local node and join lakewood testnet
2. Get some work tokens
3. Trigger transaction to store proof

First, you need to setup a local node and join lakewood. In order to do that, you need to checkout the repo https://github.com/Baseledger/baseledger-lakewood and have Docker installed. 

Once ready, navigate to the <repo root>/docs/dropping_proofs_via_local_node and run:

    docker run -d -p 1317:1317 -p 26655:26656 --name baseledger_lakewood_node baseledger/blockchain_app

This command will create and start a container baseledger_lakewood_node. Make sure your network allows incoming traffic through the 26655 port, which is being used for synchronization. 1317 is the localhost port you will be using to talk to the node, and should be protected from outside access. 

Next step is to setup your node to talk to the lakewood testnet. To do this, open the script *add-node-to-lakewood.sh* and give a nice name to your node by replacing all occurences of <your_node_name> with your name. Then run the script from the root of the repo.

The output of the script should, among other things, contain something like the following:

    name: <your_node_name>_replicator_address_1
    type: local
    address: baseledger1xvcmvc9ufkacfpr4ulzuz2jm050r9l50w2ry6s
    pubkey: baseledgerpub1addwnpepq2044540zp36szem9plt032kq6urv9lttlxdz7lynzdu0rthr3kh2ujev6z
    mnemonic: ""
    threshold: 0
    pubkeys: []

Make sure to store this information somewhere, as you will be needing it to get the test tokens later.

Your node is now configured to join lakewood and you just need to start it. Execute the following command:

    docker exec baseledger_lakewood_node baseledgerd start

You should see the log output of the node, showing progress as it starts to sync with the network. This process can take some time. You can close the window and the node will continue to run in the background.

While waiting for the sync, take the address your previously stored and send it in an email to worktokens@baseledger.net in order to get some work tokens in it. A testnet faucet is under development.

Once the node has synced and you've got your work tokens, you can trigger a transaction against your node by triggering the following endpoint:

    POST http://localhost:1317/signAndBroadcast
    {
        "payload": "your proof",
        "op_code": 9
    }

in CURL:

    curl --location --request POST 'http://localhost:1317/signAndBroadcast' --header 'Content-Type: text/plain' --data-raw '{
    "payload": "your proof",
    "op_code": 9
    }'

which will return 200 with a transaction hash.

That is it! You can verify your proof is stored by visiting https://lakewood.baseledger.net/transactions/<transaction_hash>

In case you try to drop proof without having acquired tokens beforehand, you will get an error stating that the account cannot be found.


To stop the node and cleanup, just use these two commands to stop and remove the container:

    docker stop baseledger_lakewood_node
    docker rm baseledger_lakewood_node

followed with this one to remove the image:

    docker rmi baseledger/blockchain_app