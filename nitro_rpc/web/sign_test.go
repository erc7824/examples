package web

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/layer-3/clearsync/pkg/signer"
	"github.com/stretchr/testify/require"
)

func TestServerSign(t *testing.T) {
	mh := &MethodHandlerMock{}

	privateKey, err := crypto.HexToECDSA("0000000000000000000000000000000000000000000000000000000000000001")
	require.NoError(t, err)
	ss := signer.NewLocalSigner(privateKey)

	srv := NewServer(mh, ss)

	state := RPCState{
		Timestamp: 1000,
		RequestID: 42,
		method:    "testMethod",
		params:    []byte("testParams"),
		result:    []byte("testResult"),
	}

	sig, err := srv.Sign(state)
	require.NoError(t, err)

	expectedSig := signer.NewSignature(
		hexutil.MustDecode("0x8884c2fa5269d26ae4645f776f0a1a41f6dc45aa8b8fea947392b7629bc645ec"),
		hexutil.MustDecode("0x73a4a1a2b8b5775bb42956ff666df93089c3ad984db4f7870b69841df75d1bde"),
		27,
	)
	require.Equal(t, expectedSig, sig)
}
