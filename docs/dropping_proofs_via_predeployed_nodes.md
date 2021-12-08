# Drop a proof on baseledger via predeployed test nodes

Currently available nodes for this purpose are:

alice.lakewood.baseledger.net

bob.lakewood.baseledger.net

In order to drop a proof via one of these nodes, you first need to register a user email and password with the node. This can be acomplished by triggering the following endpoint:

    POST http://alice.lakewood.baseledger.net/dev/users
    {
        "email": "<your_email>",
        "password": "<your_pass>"
    }

in CURL:

    curl --location --request POST 'http://alice.lakewood.baseledger.net/dev/    users' --header 'Content-Type: application/json' --data-raw '{
        "email": "<your_email>",
        "password": "<your_pass>"
    }'

which will return 200. Note that per user, you are limited to 10 proofs in a period of 24 hours.

Then you login with the provided user:

    POST http://alice.lakewood.baseledger.net/dev/auth
    {
        "email": "<your_email>",
        "password": "<your_pass>"
    }

in CURL:

    curl --location --request POST 'http://alice.lakewood.baseledger.net/dev/    auth' --header 'Content-Type: application/json' --data-raw '{
        "email": "<your_email>",
        "password": "<your_pass>"
    }'

which will return 200 with a JWT token.

Final step is to trigger the transaction endpoint:

    Authorization header: Bearer <jwt_token>
    POST http://alice.lakewood.baseledger.net/dev/tx
    {
        "payload": "<your_payload>", // Legth must be <= 30
        "op_code": 9 // this is the only supported op_code for now, representing simple proof storage
    }

in CURL:

    curl --location --request POST 'http://alice.lakewood.baseledger.net/dev/tx'     --header 'Authorization: Bearer <jwt_token>' --header 'Content-Type: application/json' --data-raw '{
        "payload": "<your_payload>",
        "op_code": 9
    }'

which will return 200 with a transaction hash.

That is it! You can verify your proof is stored by visiting https://lakewood.baseledger.net/transactions/<transaction_hash>