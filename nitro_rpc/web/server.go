package web

import (
	"context"

	"github.com/erc7824/examples/nitro_rpc/proto"
)

var _ proto.NitroRPC = (*Server)(nil)

type Server struct {
	mh MethodHandler
}

type MethodHandler interface {
	HandleCall(ctx context.Context, params []byte) ([]byte, error)
}

func NewServer(handler MethodHandler) *Server {
	return &Server{
		mh: handler,
	}
}

func (s *Server) Call(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	results, err := s.mh.HandleCall(ctx, req.Req.GetParams())
	if err != nil {
		return nil, err
	}

	return &proto.Response{
		Res: &proto.ResponsePayload{
			RequestId: req.Req.GetRequestId(),
			Method:    req.Req.GetMethod(),
			Timestamp: req.Req.GetTimestamp(),
			Results:   results,
		},
	}, nil
}
