package web

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/layer-3/clearsync/pkg/signer"
)

var rpcStateABIArgs abi.Arguments

type RPCState struct {
	RequestID uint64
	Timestamp uint64
	method    string
	params    []byte
	result    []byte
}

type SignedRPCState struct {
	RPCState
	ClientSig signer.Signature
	ServerSig signer.Signature
}

func (s *Server) Sign(state RPCState) (signer.Signature, error) {
	packed, err := rpcStateABIArgs.Pack(
		new(big.Int).SetUint64(state.RequestID),
		new(big.Int).SetUint64(state.Timestamp),
		state.method,
		state.params,
		state.result,
	)
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
	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		panic(err)
	}

	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		panic(err)
	}

	bytesType, err := abi.NewType("bytes", "", nil)
	if err != nil {
		panic(err)
	}

	rpcStateABIArgs = abi.Arguments{
		{Type: uint256Type},
		{Type: uint256Type},
		{Type: stringType},
		{Type: bytesType},
		{Type: bytesType},
	}
}
