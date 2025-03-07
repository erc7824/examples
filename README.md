# Nitro Examples
Nitro Protocol Examples

## Nitro RPC

Those application leverage the NitroRPC Asynchronous protocol

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
