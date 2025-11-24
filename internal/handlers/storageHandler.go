package handlers

import (
	proto "File_Storage/internal/api/storage"
	"File_Storage/internal/services/storage"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type storageServer struct {
	service storage.Service
	proto.UnimplementedStorageServiceServer
}

func NewStorageServer(service storage.Service) proto.StorageServiceServer {
	return &storageServer{service: service}
}

func (s *storageServer) UploadFile(ctx context.Context, req *proto.UploadRequest) (*proto.SuccessResponse, error) {
	if req.Token == "" || req.FileName == "" || req.MimeType == "" || req.Data == nil {
		return &proto.SuccessResponse{Success: false}, status.Error(codes.InvalidArgument, "Upload Request Error")
	}

	if err := s.service.UploadFile(req.Token, req.FileName, req.MimeType, req.Data); err != nil {
		return &proto.SuccessResponse{Success: false}, status.Error(codes.Internal, "Service Error: "+err.Error())
	}

	return &proto.SuccessResponse{Success: true}, nil
}

func (s *storageServer) GetFile(ctx context.Context, req *proto.FileRequest) (*proto.GetFileResponse, error) {
	if req.Token == "" {
		return &proto.GetFileResponse{}, status.Error(codes.InvalidArgument, "Token is empty")
	}

	if req.FileID == "" {
		return &proto.GetFileResponse{}, status.Error(codes.InvalidArgument, "File ID is empty")
	}

	meta, data, err := s.service.GetFile(req.Token, req.FileID)

	if err != nil {
		return &proto.GetFileResponse{}, status.Error(codes.Internal, "service: "+err.Error())
	}

	var metadata = &proto.FileMetadata{
		FileID:    meta.FileID,
		UserID:    meta.UserID,
		FileName:  meta.FileName,
		MimeType:  meta.MimeType,
		Size:      meta.Size,
		Location:  meta.Location,
		CreatedAt: meta.CreatedAt,
		UpdatedAt: meta.UpdatedAt,
	}

	return &proto.GetFileResponse{Meta: metadata, Data: data}, nil
}

func (s *storageServer) DeleteFile(ctx context.Context, req *proto.FileRequest) (*proto.SuccessResponse, error) {
	if req.Token == "" || req.FileID == "" {
		return &proto.SuccessResponse{Success: false}, status.Error(codes.InvalidArgument, "Token or id is empty")
	}

	if err := s.service.DeleteFile(req.Token, req.FileID); err != nil {
		return &proto.SuccessResponse{Success: false}, status.Error(codes.Internal, err.Error())
	}

	return &proto.SuccessResponse{Success: true}, nil
}

func (s *storageServer) GetAllFiles(ctx context.Context, req *proto.TokenRequest) (*proto.GetFilesResponse, error) {
	if req.Token == "" {
		return &proto.GetFilesResponse{}, status.Error(codes.InvalidArgument, "Token is empty")
	}

	meta, err := s.service.GetAllFiles(req.Token)

	if err != nil {
		return &proto.GetFilesResponse{}, status.Error(codes.Internal, err.Error())
	}

	metadata := make([]*proto.FileMetadata, 0, len(meta))

	for _, meta := range meta {
		metadata = append(metadata, &proto.FileMetadata{
			FileID:    meta.FileID,
			UserID:    meta.UserID,
			FileName:  meta.FileName,
			MimeType:  meta.MimeType,
			Size:      meta.Size,
			Location:  meta.Location,
			CreatedAt: meta.CreatedAt,
			UpdatedAt: meta.UpdatedAt,
		})
	}

	return &proto.GetFilesResponse{Meta: metadata}, nil
}
