// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "forge-std/Test.sol";
import "forge-std/console.sol";
import "../src/NitroRPC.sol";

// A helper contract to expose the internal recoverPayloadSigner function.
contract TestableNitroRPC is NitroRPC {
    function publicRecoverPayloadSigner(
        Payload memory payload, 
        INitroTypes.Signature memory signature
    ) public pure returns (address) {
        return recoverPayloadSigner(payload, signature);
    }
}

contract NitroRPCTest is Test {
    TestableNitroRPC testContract;
    // Fixed private key for testing.
    uint256 constant fixedPrivateKey = 1;

    function setUp() public {
        testContract = new TestableNitroRPC();
    }

    function testFixedInputsOutputs() public view {
        // Use fixed values for payload so the test output is deterministic.
        NitroRPC.Payload memory payload = NitroRPC.Payload({
            requestId: 42,
            timestamp: 1000,
            method: "testMethod",
            params: bytes("testParams"),
            result: bytes("testResult")
        });

        // 1. Compute the message hash from the payload.
        bytes32 messageHash = keccak256(
            abi.encode(
                payload.requestId,
                payload.timestamp,
                payload.method,
                payload.params,
                payload.result
            )
        );
        
        // 2. Apply the Ethereum Signed Message prefix.
        bytes32 ethSignedMessageHash = keccak256(
            abi.encodePacked("\x19Ethereum Signed Message:\n32", messageHash)
        );

        // 3. Sign the hash using Foundry's cheatcode.
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(fixedPrivateKey, ethSignedMessageHash);
        address expectedSigner = vm.addr(fixedPrivateKey);

        // Build the signature struct.
        INitroTypes.Signature memory signature = INitroTypes.Signature({
            v: v,
            r: r,
            s: s
        });

        // Recover the signer using the testable function.
        address recovered = testContract.publicRecoverPayloadSigner(payload, signature);

        // Log the values in hexadecimal representations.
        console.log("privateKey:");
        console.logBytes32(bytes32(fixedPrivateKey));

        console.log("v:");
        console.logUint(uint256(v));

        console.log("r:");
        console.logBytes32(r);

        console.log("s:");
        console.logBytes32(s);

        // Verify that the recovered address matches the expected signer.
        assertEq(recovered, expectedSigner, "Recovered signer does not match expected signer");
    }
}
