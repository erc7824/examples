package web

import (
	"context"
	"math/big"

	"github.com/erc7824/examples/nitro_rpc/proto"
	"github.com/layer-3/clearsync/pkg/signer"
)

var _ proto.NitroRPC = (*Server)(nil)

type Server struct {
	mh          MethodHandler
	stateSigner signer.Signer
}

type MethodHandler interface {
	HandleCall(ctx context.Context, params []byte) ([]byte, error)
}

func NewServer(handler MethodHandler, stateSigner signer.Signer) *Server {
	return &Server{
		mh:          handler,
		stateSigner: stateSigner,
	}
}

func (s *Server) Call(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	result, err := s.mh.HandleCall(ctx, req.Req.GetParams())
	if err != nil {
		return nil, err
	}

	state := RPCState{
		RequestID: new(big.Int).SetUint64(req.Req.GetRequestId()),
		Timestamp: new(big.Int).SetUint64(req.Req.GetTimestamp()),
		Method:    req.Req.GetMethod(),
		Params:    req.Req.GetParams(),
		Result:    result,
	}
	sig, err := s.Sign(state)
	if err != nil {
		return nil, err
	}

	return &proto.Response{
		Res: &proto.ResponsePayload{
			RequestId: state.RequestID.Uint64(),
			Timestamp: state.Timestamp.Uint64(),
			Method:    state.Method,
			Results:   state.Result,
		},
		Sig: sig.String(),
	}, nil
}
