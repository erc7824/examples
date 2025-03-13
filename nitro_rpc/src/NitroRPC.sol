// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import {INitroTypes} from "nitro/interfaces/INitroTypes.sol";
import {NitroUtils} from "nitro/libraries/NitroUtils.sol";
import {IForceMoveApp} from "nitro/interfaces/IForceMoveApp.sol";

/**
 * @title NitroRPC
 * @dev Solidity implementation of Nitro RPC protocol structures
 */
contract NitroRPC is IForceMoveApp {
    struct Payload {
        uint256 requestId;
        uint256 timestamp;
        string method;
        bytes params;
        bytes result;
    }

    struct PayloadSigned {
      Payload rpcMessage;
      INitroTypes.Signature clientSig;
      INitroTypes.Signature serverSig;
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
        INitroTypes.FixedPart calldata fixedPart,
        INitroTypes.RecoveredVariablePart[] calldata proof,
        INitroTypes.RecoveredVariablePart calldata candidate
    ) external pure override returns (bool, string memory) {
        require(fixedPart.participants.length == uint256(type(AllocationIndices).max) + 1, "bad number of participants");
        
        PayloadSigned memory payloadSigned = abi.decode(candidate.variablePart.appData, (PayloadSigned));
        requireValidPayload(fixedPart, payloadSigned);

        return (true, "");
    }

    function requireValidPayload(INitroTypes.FixedPart calldata fixedPart, PayloadSigned memory payloadSigned) internal pure {
        require(recoverPayloadSigner(payloadSigned.rpcMessage, payloadSigned.clientSig) == fixedPart.participants[uint256(AllocationIndices.Client)], "bad client signature");
        require(recoverPayloadSigner(payloadSigned.rpcMessage, payloadSigned.serverSig) == fixedPart.participants[uint256(AllocationIndices.Server)], "bad server signature");

        // TODO: verify timestamp and requestId
    }

    // This pure internal function recovers the signer address from the payload and its signature.
    function recoverPayloadSigner(Payload memory payload, INitroTypes.Signature memory signature) internal pure returns (address) {
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
        
        return ECDSA.recover(MessageHashUtils.toEthSignedMessageHash(messageHash), signature.v, signature.r, signature.s);
    }
}
