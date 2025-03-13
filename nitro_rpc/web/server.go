package web

import (
	"context"

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
		RequestID: req.Req.GetRequestId(),
		Timestamp: req.Req.GetTimestamp(),
		method:    req.Req.GetMethod(),
		params:    req.Req.GetParams(),
		result:    result,
	}
	sig, err := s.Sign(state)
	if err != nil {
		return nil, err
	}

	return &proto.Response{
		Res: &proto.ResponsePayload{
			RequestId: state.RequestID,
			Method:    state.method,
			Timestamp: state.Timestamp,
			Results:   state.result,
		},
		Sig: sig.String(),
	}, nil
}
