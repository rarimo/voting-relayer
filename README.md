# voting-relayer

## Description

Relayer service that makes calls to the smart contracts. Now it is used to work with such endpoints:
  - POST `/integrations/voting-relayer/v1/register` - it takes calldata in hex encoided format from Rarime mobile applications and makes call to the smart contract.

    Request body example:
    ```json
    {
      "data": {
        "tx_data": "hex_encoded_call_data",
        "destination": "contract_address,optional"
      }
    }
    ```

  - GET `/integrations/voting-relayer/v1/operations/latest`  - it returnes info about latest transited root
    Response body example 

    ```json
    {
       "data": {
       "id": "",
       "type": "",
       "attributes": {
          "block_height": 0,
          "destination_chain": "",
          "operation_id": "",
          "proof": "",
          "tx_hash": ""
       }
    },
    "included": []
    }
    ```

  - GET `/integrations/voting-relayer/v1/operations/{root}`  - it returnes info about transited root, where root is lower case hex string with 0x prefix
    Response body example

     ```json
     {
        "data": {
        "id": "",
        "type": "",
        "attributes": {
           "block_height": 0,
           "destination_chain": "",
           "operation_id": "",
           "proof": "",
           "tx_hash": ""
        }
     },
     "included": []
     }
     ```

## Install

  ```
  git clone github.com/rarimo/voting-relayer
  cd voting-relayer
  go build main.go
  export KV_VIPER_FILE=./config.yaml
  ./main run service
  ```

## Documentation

We do use openapi:json standard for API. We use swagger for documenting our API.

To open online documentation, go to [swagger editor](http://localhost:8080/swagger-editor/) here is how you can start it
```
  cd docs
  npm install
  npm start
```
To build documentation use `npm run build` command,
that will create open-api documentation in `web_deploy` folder.

To generate resources for Go models run `./generate.sh` script in root folder.
use `./generate.sh --help` to see all available options.

Note: if you are using Gitlab for building project `docs/spec/paths` folder must not be
empty, otherwise only `Build and Publish` job will be passed.  

## Config
```yaml
  network:
    rpc: "" # (url) RPC API endpoint | required
    registration: "" # (hex) target contract address | required
    private_key: "" # (hex without 0x) ECDSA secp256k1 private key for sign transactions
    vault_address: "http://127.0.0.1:8200" # (url) vault address
    vault_mount_path: "secret_data" # (string)
    whitelist: # (list of hex addresses) specify which contracts can be passed in `destination` field in request
      - "0x123...123"
      - "0x123...123"
```
ENV
```
  VAULT_TOKEN (will be cleared after start service)
```
There must be specified or private key, or vault address and path. If specified both, then will be used private_key from config.

## Running from docker 
  
Make sure that docker installed.

use `docker run ` with `-p 8080:80` to expose port 80 to 8080

  ```
  docker build -t github.com/rarimo/voting-relayer .
  docker run -e KV_VIPER_FILE=/config.yaml github.com/rarimo/voting-relayer
  ```

## Running from Source

* Set up environment value with config file path `KV_VIPER_FILE=./config.yaml`
* Provide valid config file
* Launch the service with `run service` command to run full mode (voting and relayer)

### Voting mode 
To run service only in voting mode should use `run voting`

### Relayer mode 
To run service only in relayer mode should use `run relayer`


### Third-party services


## Contact

Responsible 
The primary contact for this project is  [//]: # (TODO: place link to your telegram and email)
