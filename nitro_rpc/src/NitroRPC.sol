// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

/**
 * @title NitroRPC
 * @dev Solidity implementation of Nitro RPC protocol structures
 */
interface NitroRPC {

    struct Payload {
        uint256 requestId;
        string method;
        bytes params;
        bytes result;
        uint256 timestamp;
    }

    struct PayloadSigned {
      Payload rpcMessage;
      INitroTypes.Signature clientSig
      INitroTypes.Signature serverSig
    }
}
