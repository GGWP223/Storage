package handlers

import (
	proto "File_Storage/internal/api/auth"
	"File_Storage/internal/services/auth"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authServer struct {
	service auth.Service
	proto.UnimplementedAuthServiceServer
}

func NewAuthHandler(service auth.Service) proto.AuthServiceServer {
	return &authServer{service: service}
}

func (s *authServer) Register(ctx context.Context, req *proto.AuthRequest) (*proto.RegisterResponse, error) {
	request := auth.Request{
		Login:    req.Login,
		Password: req.Password,
	}

	if err := s.service.Register(request); err != nil {
		return &proto.RegisterResponse{Success: false}, status.Error(codes.Internal, err.Error())
	}

	return &proto.RegisterResponse{Success: true}, nil
}

func (s *authServer) Login(ctx context.Context, req *proto.AuthRequest) (*proto.LoginResponse, error) {
	request := auth.Request{
		Login:    req.Login,
		Password: req.Password,
	}

	access, refresh, err := s.service.Login(request)

	if err != nil {
		return &proto.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.LoginResponse{Access: access, Refresh: refresh}, nil
}
