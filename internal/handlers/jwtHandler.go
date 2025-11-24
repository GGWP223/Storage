package handlers

import (
	proto "File_Storage/internal/api/jwt"
	"File_Storage/internal/services/jwt"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type jwtServer struct {
	service jwt.Service
	proto.UnimplementedJwtServiceServer
}

func NewJwtServer(service jwt.Service) proto.JwtServiceServer {
	return &jwtServer{service: service}
}

func (s *jwtServer) RefreshToken(ctx context.Context, req *proto.TokenRequest) (*proto.TokenResponse, error) {
	if req.Token == "" {
		return &proto.TokenResponse{}, status.Error(codes.InvalidArgument, "Token is empty")
	}

	token, err := s.service.RefreshToken(req.Token)

	if err != nil {
		return &proto.TokenResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.TokenResponse{Token: token}, nil
}

func (s *jwtServer) ValidateToken(ctx context.Context, req *proto.TokenRequest) (*proto.ValidateResponse, error) {
	if req.Token == "" {
		return &proto.ValidateResponse{}, status.Error(codes.InvalidArgument, "Token is empty")
	}

	claims, err := s.service.ValidateToken(req.Token)

	if err != nil {
		return &proto.ValidateResponse{Success: false}, status.Error(codes.Internal, "Token is invalid")
	}

	return &proto.ValidateResponse{
		Success: true,
		Claims:  &proto.Claims{Uid: claims.Uid, Login: claims.Login},
	}, nil
}

func (s *jwtServer) GenerateToken(ctx context.Context, req *proto.GenerateRequest) (*proto.TokenResponse, error) {
	if req.Uid == "" || req.Login == "" {
		return &proto.TokenResponse{}, status.Error(codes.InvalidArgument, "uid or login is empty")
	}

	if req.TokenType != jwt.RefreshToken && req.TokenType != jwt.AccessToken {
		return &proto.TokenResponse{}, status.Error(codes.InvalidArgument, "token type wrong")
	}

	token, err := s.service.GenerateToken(req.Login, req.Uid, req.TokenType)

	if err != nil {
		return &proto.TokenResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.TokenResponse{Token: token}, nil
}
