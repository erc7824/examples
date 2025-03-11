// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;


import {INitroTypes} from "./TempINitroTypes.sol";
import {NitroUtils} from "./TempNitroUtils.sol";

/**
 * @title NitroRPC
 * @dev Solidity implementation of Nitro RPC protocol structures
 */
interface NitroRPC {
    struct Payload {
        uint256 requestId;
        uint256 timestamp;
        string method;
        bytes params;
        bytes result;
    }

    struct PayloadSigned {
      Payload rpcMessage;
      Signature clientSig;
      Signature serverSig;
    }

    enum AllocationIndices {
        Client,
        Server
    }

    /**
     * @notice Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @dev Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @param fixedPart Fixed part of the state channel.
     * @param proof Array of recovered variable parts which constitutes a support proof for the candidate.
     * @param candidate Recovered variable part the proof was supplied for.
     */
    function stateIsSupported(
        FixedPart calldata fixedPart,
        RecoveredVariablePart[] calldata proof,
        RecoveredVariablePart calldata candidate
    ) external pure override returns (bool, string memory) {
        require(fixedPart.participants.length == uint256(AllocationIndices.Server) + 1, "bad number of participants");
        
        PayloadSigned memory payloadSigned = abi.decode(candidate.variablePart.appData, (PayloadSigned));
        requireValidPayload(payloadSigned);

        return (true, "");
    }

    function requireValidPayload(PayloadSigned memory payloadSigned) internal pure {
        require(recoverPayloadSigner(payloadSigned.rpcMessage, payloadSigned.clientSig) == fixedPart.participants[AllocationIndices.Client], "bad client signature");
        require(recoverPayloadSigner(payloadSigned.rpcMessage, payloadSigned.serverSig) == fixedPart.participants[AllocationIndices.Server], "bad server signature");

        // TODO: verify timestamp and requestId
    }

    // This pure internal function recovers the signer address from the payload and its signature.
    function recoverPayloadSigner(Payload memory payload, Signature memory signature) internal pure returns (address) {
        // Encode and hash the payload data.
        // Using abi.encode ensures proper padding and decoding, avoiding potential ambiguities with dynamic types.
        bytes32 messageHash = keccak256(
            abi.encode(
                payload.requestId,
                payload.timestamp,
                payload.method,
                payload.params,
                payload.result
            )
        );
        
        // Apply the Ethereum Signed Message prefix.
        bytes32 ethSignedMessageHash = keccak256(
            abi.encodePacked("\x19Ethereum Signed Message:\n32", messageHash)
        );
        
        // Recover the signer address using ecrecover.
        return ecrecover(ethSignedMessageHash, signature.v, signature.r, signature.s);
    }
}
