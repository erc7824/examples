# Nitro RPC

## Protobuf & Twirp Setup

This guide walks you through setting up `protoc`, adding the required plugins, and configuring your development environment to work seamlessly with Twirp and protobuf.
Twirp uses your protobuf service definitions to generate an HTTP/1.1 API. It’s lightweight and simpler than gRPC, which can be a good fit if you prefer HTTP/1.1 over HTTP/2.

### Prerequisites

Before generating the Go files from your protobuf definitions, ensure that you have installed the Protocol Buffers compiler (`protoc`) and the necessary Go plugins for Twirp and protobuf generation. Detailed installation instructions are available in the [Twirp Installation Documentation](https://twitchtv.github.io/twirp/docs/install.html).

### Usage

All protobuf definitions are maintained inside the `proto/` directory. You can generate the corresponding Go files using the following command:

```bash
protoc -I proto \
  --go_out=proto --go_opt=paths=source_relative \
  --twirp_out=proto --twirp_opt=paths=source_relative \
  proto/nitro_rpc.proto
```

This command will create both the `.pb.go` and `.twirp.go` files in the `proto/` directory, keeping your service definitions and generated code neatly organized.

## Foundry

**Foundry is a blazing fast, portable and modular toolkit for Ethereum application development written in Rust.**

Foundry consists of:

-   **Forge**: Ethereum testing framework (like Truffle, Hardhat and DappTools).
-   **Cast**: Swiss army knife for interacting with EVM smart contracts, sending transactions and getting chain data.
-   **Anvil**: Local Ethereum node, akin to Ganache, Hardhat Network.
-   **Chisel**: Fast, utilitarian, and verbose solidity REPL.

## Documentation

https://book.getfoundry.sh/

## Usage

### Build

```shell
$ forge build
```

### Test

```shell
$ forge test
```

### Format

```shell
$ forge fmt
```

### Gas Snapshots

```shell
$ forge snapshot
```

### Anvil

```shell
$ anvil
```

### Deploy

```shell
$ forge script script/Counter.s.sol:CounterScript --rpc-url <your_rpc_url> --private-key <your_private_key>
```

### Cast

```shell
$ cast <subcommand>
```

### Help

```shell
$ forge --help
$ anvil --help
$ cast --help
```
