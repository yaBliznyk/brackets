package server

import (
	"context"
	"github.com/yaBliznyk/brackets/internal/brackets"
	"github.com/yaBliznyk/brackets/proto"
)

type server struct {
	bktSvc brackets.Brackets
	proto.UnimplementedBracketsServer
}

func (s *server) Validate(ctx context.Context, request *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	isValid := s.bktSvc.Validate(ctx, request.Str)
	return &proto.ValidateResponse{IsValid: isValid}, nil
}

func (s *server) Fix(ctx context.Context, request *proto.FixRequest) (*proto.FixResponse, error) {
	err, result := s.bktSvc.Fix(ctx, request.Str)
	return &proto.FixResponse{Result: result}, err
}

func NewBracketServer(bktSvc brackets.Brackets) proto.BracketsServer {
	return &server{
		bktSvc: bktSvc,
	}
}
