// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "forge-std/Test.sol";
import "forge-std/console.sol";
import "../src/NitroRPC.sol";

// A helper contract to expose the internal recoverPayloadSigner function.
contract NitroRPCHarness is NitroRPC {
    function exposed_recoverPayloadSigner(
        Payload memory payload, 
        bytes memory signature
    ) public pure returns (address) {
        return recoverPayloadSigner(payload, signature);
    }
}

contract NitroRPCTest is Test {
    NitroRPCHarness testContract;
    // Fixed private key for testing.
    uint256 constant fixedPrivateKey = 1;

    function setUp() public {
        testContract = new NitroRPCHarness();
    }

    function test_recoverPayloadSigner_fixedPayloadAndSigner() public view {
        // Use fixed values for payload so the test output is deterministic.
        NitroRPC.Payload memory payload = NitroRPC.Payload({
            requestId: 42,
            timestamp: 1000,
            method: "testMethod",
            params: bytes("testParams"),
            result: bytes("testResult")
        });

        // 1. Compute the message hash from the payload.
        bytes32 messageHash = keccak256(abi.encode(payload));
        
        // 2. Apply the Ethereum Signed Message prefix.
        bytes32 ethSignedMessageHash = keccak256(
            abi.encodePacked("\x19Ethereum Signed Message:\n32", messageHash)
        );

        // 3. Sign the hash using Foundry's cheatcode.
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(fixedPrivateKey, ethSignedMessageHash);
        address expectedSigner = vm.addr(fixedPrivateKey);

        // 4. Create the signature object.
        bytes memory signature = abi.encodePacked(r, s, v);
        // Recover the signer using the testable function.
        address recovered = testContract.exposed_recoverPayloadSigner(payload, signature);

        // Log the values in hexadecimal representations.
        console.log("privateKey:");
        console.logBytes32(bytes32(fixedPrivateKey));

        console.log("signature:");
        console.logBytes(signature);

        // Verify that the recovered address matches the expected signer.
        assertEq(recovered, expectedSigner, "Recovered signer does not match expected signer");
    }
}
