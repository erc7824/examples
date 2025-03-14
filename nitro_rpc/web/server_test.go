package web

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/erc7824/examples/nitro_rpc/proto"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/layer-3/clearsync/pkg/signer"
	"github.com/stretchr/testify/require"
)

type MethodHandlerMock struct{}

func (m *MethodHandlerMock) HandleCall(ctx context.Context, params []byte) ([]byte, error) {
	return params, nil
}

func TestServerCall(t *testing.T) {
	mh := &MethodHandlerMock{}

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	ss := signer.NewLocalSigner(privateKey)

	handler := proto.NewNitroRPCServer(NewServer(mh, ss))
	srv := http.Server{
		Addr:    ":8089",
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("server error: %v", err)
		}
	}()
	defer srv.Shutdown(context.Background())

	client := proto.NewNitroRPCProtobufClient("http://localhost:8089", &http.Client{})
	req := &proto.Request{
		Req: &proto.RequestPayload{
			RequestId: 1,
			Method:    "test",
			Params:    []byte("test"),
			Timestamp: uint64(time.Now().UTC().UnixNano()),
		},
	}
	res, err := client.Call(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, req.Req.GetRequestId(), res.Res.GetRequestId())
	require.Equal(t, req.Req.GetMethod(), res.Res.GetMethod())
	require.Equal(t, req.Req.GetTimestamp(), res.Res.GetTimestamp())
	require.Equal(t, req.Req.GetParams(), res.Res.GetResults())
}
