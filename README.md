# rollup-node

Rollup-node is an optional module for DappLink Appchain. It integrates five DA interfaces, which are Anytrust, Celestia,
EigenDA, NearDA, and EIP-4844.

Users can roll up data to a specified DA and receive a receipt. They can also retrieve corresponding data from the
specified DA based on the receipt.

Users can set configuration options for all five DAs simultaneously, or configure only the DAs that they intend to use
among the five.

To use the rollup module, the corresponding DA needs to be prepared in advance. For example, Anytrust requires an
off-chain DA service to be run beforehand, and NearDA requires a Near account, contract, and private key to be prepared
in advance.

We provide APIs and SDKs to offer users a convenient way to interact.

## Run Rollup Node

- Submodule

  run `git submodule update --init --recursive --remote ` to update submodule
- Additional Considerations for NearDA

  Because NearDA uses a C code library `near_da_rpc_sys.a`, we need to place this library file in the specified location
  before compilation. The library file and related documentation are located in the `./c-lib` directory.
- Compilation & Run

  `make build` and `./rollupNode rollup-node --rpcAddress localhost:9000 --apiAddress localhost:9001`

  or

  `go build` and `./rollup-node --rpcAddress localhost:9000 --apiAddress localhost:9001`

## API & SDK

- API

  When starting the rollup-node, you need to set `--apiAddress` as the listening address for the web server.

    - rollup & retrieve

      | route | type | args                                       | comment                                             |
      |:----- |:-----|:-------------------------------------------|:----------------------------------------------------|
      |`/api/v1/rollup-with-type`| post | `{"da_type": 4,"data":"base64 string"}`    | Rollup data to a specified DA |
      |`/api/v1/retrieve-with-type` | post |  `{"da_type": 4, "args":"rollup receipt"}` | Retrieve data from specified DA with rollup receipt |



- SDK 
    When starting the rollup-node, you need to set `--rpcAddress` as the listening address for the web server.
  - new a sdk: `rollupSdk, err := sdk.NewRollupSdk(rpcAddress)`
  - rollup: `rollupSdk.RollupWithType(dataByte, daType)`
  - retrieve: `rollupSdk.RetrieveWithType(daType, rollupReceipt)`


## Configs & Envs

- Anytrust

    config file: `./config/anytrust.toml` and all fields can be set by env.
- Celestia
- EigenDA

    config file: `./config/eigenda.toml` and all fields can be set by env.
- Eip-4844
- NearDA

    config file: `./config/nearda.toml` and all fields can be set by env.