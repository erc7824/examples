# Nitro Examples
Nitro Protocol Examples

## Nitro RPC

Those application leverage the NitroRPC Asynchronous protocol, it describe a data format, that must be understood and readable by both backend, frontend and smart-contract NitroRPCApp (adjudication)

Here is the format:
```solidity
struct NitroRPC {
    uint64 req_id;   // Unique request ID (non-zero)
    string method;   // RPC method name (e.g., "substract")
    string[] params; // Array of parameters for the RPC call
    uint64 ts;       // Milisecond unix timestamp provided by the server api
}

NitroRPC memory rpc = NitroRPC({
    req_id: 123,
    method: "subtract",
    params: new string  ts: 1710474422
});

bytes memory encoded1 = ABI.encode(rpc);
bytes memory encoded2 = ABI.encode(rpc.req_id, rpc.method, rpc.params, rpc.ts);

require(keccak256(encoded1) == keccak256(encoded2), "Mismatch in encoding!");
```

### The RPCHash

Client and server must sign the Nitro RPC Hash as followed

### Solidity

```solidity
rpc_hash = keccak256(
  abi.encode(
    rpc.req_id,
    rpc.method,
    rpc.params,
    rpc.ts
  )
);

# rpc_hash can be used to erecover the public key
```

```go
package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

type NitroRPC struct {
	ReqID  uint64
	Method string
	Params []string
	TS     uint64
}

func main() {
	rpc := NitroRPC{
		ReqID:  123,
		Method: "subtract",
		Params: []string{"param1", "param2"},
		TS:     1710474422,
	}

	packedData, err := abi.Arguments{
		{Type: abi.UintTy},
		{Type: abi.StringTy},
		{Type: abi.StringSliceTy},
		{Type: abi.UintTy},
	}.Pack(rpc.ReqID, rpc.Method, rpc.Params, rpc.TS)
	if err != nil {
		log.Fatal("ABI encoding failed:", err)
	}

	hash := crypto.Keccak256(packedData)
	fmt.Println("Keccak256 Hash:", hexutil.Encode(hash))
}
```

## NitroRPC Transport

In NitroRPC the transport layer is agnostic to describe in any syntax as long as the NitroRPC Type, RPCHash and signature convention are valid.
You can use json-rpc, msgpack, protobug, gRPC or any custom marshalling and transport layer.

### NitroRPC Request

In this example, the client is calling a method named `"subtract"` with positional parameters:

```json
{
  "req": [1001, "subtract", [42, 23], 1741344819012],
  "sig": "0xa0ad67f51cc73aee5b874ace9bc2e2053488bde06de257541e05fc58fd8c4f149cca44f1c702fcbdbde0aa09bcd24456f465e5c3002c011a3bc0f317df7777d2"
}
```

- req: rpc message payload `[request_id, method, params, ts]`
- sig: payload client signature

The millisecond timestamp was returned previously from the server, it is used as a height for proof of history

### NitroRPC Response Example

For a successful invocation, the server might respond like this:

```json
{
  "res": [1001, "subtract", [19], 1741344819814],
  "sig": "0xd73268362b04516451ec52170f5c8ca189d35d9ac5e9041c156c9f0faf9aebd2891309e3b2b5d8788578ab3449c96f7aa81aefb25482b53f02bac42c65f806e5"
}
```

- res: rpc message payload `[request_id, method, params, ts]`
- sig: payload server response signature

**ts**: server response with latest timestamp

### NitroCalc

Objectives of the prototype is to define opening, checkpoint, dispute and finalize of channel.

### State storage example

```sql
-- Table to store RPC messages
CREATE TABLE rpc_states (
    id SERIAL PRIMARY KEY,
    ts BIGINT,
    req_id INTEGER NOT NULL,
    method VARCHAR(255) NOT NULL,
    params JSONB NOT NULL,
    client_sig VARCHAR(128),
    server_sig VARCHAR(128),
    UNIQUE (request_id)
);
```

#### TODO:
- Create AppData in solidity
- Create NitroApp to validate states including AppData
- Create Protobuf description of NitroRPC Request / Response
- Include NitroRPC types into gRPC and generate client and server
- Implement methods: `add(int), sub(int), mul(int), mod(int)` starting from `0`
- Create a simplified NitroCalc Service using sqlite for state (in go)
- Finally, we can import nitro pkg, and client or server can pack and sqlite state and call nitro of on-chain channel operations.

### Questions

#### Where state-management is happening?

  Like the TicTacToe example, the Logic and App State is managed by the Nitro.App,
  But also the state is created in compliance by implementing correctly Nitro.App interface See ADR-002
  nitro package responsability is to unpack the state and submit it on-chain.
  
#### Communication diagram between BusinessApp, NitroRPC and NitroCore

NitroRPC is embeded into the BusinessApp and it's only a norm, expect for the smart-contract NitroRPCApp

nitro go package is providing helpers and functions to abstract the blockchain level things,
it will take your Nitro.App state and execute the blockchain operation you request on the NitroAdjudicator (`prefunding, postfunding, finalize, checkpoint, challenge`)

#### Who is responsible for the state signing? (One signer on Client, One signer on Server???)

Client Nitro.App signs requests
Server Nitro.App signs reponses

Channel service or nitro package does sign for you, the private key is obviously not part of the package.
But nitro pkg will help you sign and simplify all the nitro state creation.
  
#### Do we have 2-step communication like?
  - Client -> Server: Request
  - Server -> Client: Response (Server signs?)
  - Client -> Server: Acknowledge (Client signs?)
  - Server -> Client: Acknowledge
 
I would say 1-step is Request-Response pair.

- Request is signed by client
- Response is signed by server

anything else is an invalidate state (request signed without response, signature missing)

#### Do we implement nitro specific methods (open, checkpoint, dispute, finalize) in the NitroRPC service?

You only need to implement the transport, for example json-rpc handlers, or grpc, those specific method will be standardized
nitro pkg will provide the implementation of those methods, you just need to provide the correct prefunding postfunding state in nitro.App

A nitro.App is a state container, which can hold 1 or many state indexed by TurnNum, serialized and passed to nitro pkg/svc for execution.
  
#### Does NitroRPC server have BusinessApp logic?

NitroRPC is just a convention, the Application has the business logic and implement an RPC protocol which comply with the state convention
  
#### Does NitroRPC server stores rpc states?

It's high recommended, in the event of answering to a challenge, but most of the time you need only the recent state,
but like I provided an SQL table, Application should keep track of state in some degree, it could be in memory and in a custom format
as long as it's able to form an RPCHash again.

#### Markdown Table Example

| ID   | Timestamp (ts)     | Method         | Request Arguments                  | Response Arguments                |
|------|--------------------|----------------|------------------------------------|-----------------------------------|
| 1002 | 1741344820000      | open_channel   | [100]          | [chanId]                      |
| 1003 | 1741344821000      | add    | [50]              | [50]                      |
| 1004 | 1741344822000      | mul   | [2]                      | [100]                  |
| 1005 | 1741344823000      | sub | [10]                     | [90]        |
| 1006 | 1741344820000      | close_channel   | [nitro_state]          | [mutually_signed_state]                      |

And create a unhappy case

### NitroBet

An idea of an app with external event, could be 
Alice and Bob start both with $100
they continuous add money to their bet
and a outcome of the bet would determine the winner

To explore
