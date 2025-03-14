package web

import (
	"math/big"
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
		Timestamp: new(big.Int).SetUint64(1000),
		RequestID: new(big.Int).SetUint64(42),
		Method:    "testMethod",
		Params:    []byte("testParams"),
		Result:    []byte("testResult"),
	}

	sig, err := srv.Sign(state)
	require.NoError(t, err)

	expectedSig := signer.NewSignatureFromBytes(hexutil.MustDecode("0x4a9aa9cafc9e804a17cf56ba8c6ff82a5dffacb00997d2f95005636ed8b40de32fbaaf2361b62dfd1bbecfd6603422d5497eda127fdfb902951e2acff17c9c361c"))
	require.Equal(t, expectedSig, sig)
}
