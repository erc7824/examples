package web

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/layer-3/clearsync/pkg/signer"
)

var rpcStateABIArgs abi.Arguments

type RPCState struct {
	RequestID *big.Int `abi:"requestId"`
	Timestamp *big.Int `abi:"timestamp"`
	Method    string   `abi:"method"`
	Params    []byte   `abi:"params"`
	Result    []byte   `abi:"result"`
}

type SignedRPCState struct {
	RPCState
	ClientSig signer.Signature
	ServerSig signer.Signature
}

func (s *Server) Sign(state RPCState) (signer.Signature, error) {
	packed, err := rpcStateABIArgs.Pack(state)
	if err != nil {
		return signer.Signature{}, err
	}

	messageHash := crypto.Keccak256(packed)
	sig, err := signer.SignEthMessage(s.stateSigner, messageHash)
	if err != nil {
		return signer.Signature{}, err
	}

	return sig, nil
}

func init() {
	tupleArgs := []abi.ArgumentMarshaling{
		{Name: "requestId", Type: "uint256"},
		{Name: "timestamp", Type: "uint256"},
		{Name: "method", Type: "string"},
		{Name: "params", Type: "bytes"},
		{Name: "result", Type: "bytes"},
	}

	tupleType, err := abi.NewType("tuple", "", tupleArgs)
	if err != nil {
		panic(err)
	}

	rpcStateABIArgs = abi.Arguments{
		{
			Type: tupleType,
		},
	}
}
