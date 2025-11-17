package auth

import (
	proto "File_Storage/internal/api/auth"
	"File_Storage/internal/services/auth"
	"context"
)

type server struct {
	service auth.Service
	proto.UnimplementedAuthServiceServer
}

func NewGRPCAuthHandler(service auth.Service) proto.AuthServiceServer {
	return &server{service: service}
}

func (s *server) Register(ctx context.Context, req *proto.AuthRequest) (*proto.RegisterResponse, error) {
	request := auth.Request{
		Login:    req.Login,
		Password: req.Password,
	}

	if err := s.service.Register(request); err != nil {
		return &proto.RegisterResponse{Success: false}, err
	}

	return &proto.RegisterResponse{Success: true}, nil
}

func (s *server) Login(ctx context.Context, req *proto.AuthRequest) (*proto.LoginResponse, error) {
	request := auth.Request{
		Login:    req.Login,
		Password: req.Password,
	}

	token, err := s.service.Login(request)

	if err != nil {
		return &proto.LoginResponse{Success: false}, err
	}

	return &proto.LoginResponse{Success: true, Token: token}, nil
}

func (s *server) ValidateToken(ctx context.Context, req *proto.TokenRequest) (*proto.ValidateResponse, error) {
	if _, err := s.service.ValidateToken(req.Token); err != nil {
		return &proto.ValidateResponse{Success: false}, err
	}

	return &proto.ValidateResponse{Success: true}, nil
}

func (s *server) RefreshToken(ctx context.Context, req *proto.TokenRequest) (*proto.RefreshResponse, error) {
	token, err := s.service.RefreshToken(req.Token)

	if err != nil {
		return &proto.RefreshResponse{Success: false}, err
	}

	return &proto.RefreshResponse{Success: true, Token: token}, nil
}

func (s *server) ParseToken(ctx context.Context, req *proto.TokenRequest) (*proto.ParseResponse, error) {
	claims, err := s.service.ValidateToken(req.Token)

	if err != nil {
		return &proto.ParseResponse{}, err
	}

	return &proto.ParseResponse{Uid: claims.Uid, Login: claims.Login}, nil
}
