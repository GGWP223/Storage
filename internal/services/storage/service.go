package storage

import (
	proto "File_Storage/internal/api/jwt"
	"context"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	UploadFile(token, filename, mimeType string, data []byte) error
	GetFile(token, fileID string) (FileMetadata, []byte, error)
	DeleteFile(token, fileID string) error
	GetAllFiles(token string) ([]FileMetadata, error)
}

type service struct {
	jwtClient proto.JwtServiceClient
	repo      Repository
}

func NewService(repo Repository, client proto.JwtServiceClient) Service {
	return &service{repo: repo, jwtClient: client}
}

func (s *service) UploadFile(token, filename, mimeType string, data []byte) error {
	response, err := s.ParseToken(context.Background(), token)

	if err != nil {
		return errors.New("Error parsing token: " + token + ", err: " + err.Error())
	}

	path := "./uploads/" + uuid.NewString() + "_" + filename

	if err := os.MkdirAll("./uploads", 0755); err != nil {
		return errors.New("Error creating uploads directory: " + err.Error())
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return errors.New("Error writing file: " + err.Error())
	}

	meta := FileMetadata{
		FileID:    uuid.NewString(),
		UserID:    response.Claims.Uid,
		FileName:  filename,
		MimeType:  mimeType,
		Size:      int64(len(data)),
		Location:  path,
		CreatedAt: time.Now().Format(time.DateTime),
		UpdatedAt: time.Now().Format(time.DateTime),
	}

	return s.repo.SaveFile(meta)
}

func (s *service) GetFile(token, fileID string) (FileMetadata, []byte, error) {
	response, err := s.ParseToken(context.Background(), token)

	if err != nil {
		return FileMetadata{}, nil, errors.New("Error parsing token: " + token + ", err: " + err.Error())
	}

	request := Request{
		UserID: response.Claims.Uid,
		FileID: fileID,
	}

	meta, err := s.repo.GetFile(request)

	if err != nil {
		return FileMetadata{}, nil, errors.New("Error getting file: " + err.Error())
	}

	data, err := os.ReadFile(meta.Location)

	if err != nil {
		return FileMetadata{}, nil, errors.New("Error reading file: " + err.Error())
	}

	return meta, data, nil
}

func (s *service) DeleteFile(token, fileID string) error {
	response, err := s.ParseToken(context.Background(), token)

	if err != nil {
		return errors.New("Error parsing token: " + token + ", err: " + err.Error())
	}

	request := Request{
		UserID: response.Claims.Uid,
		FileID: fileID,
	}

	file, err := s.repo.GetFile(request)

	if err != nil {
		return errors.New("Error getting file: " + err.Error())
	}

	if err := s.repo.DeleteFile(file); err != nil {
		return errors.New("Error deleting file: " + err.Error())
	}

	return os.Remove(file.Location)
}

func (s *service) GetAllFiles(token string) ([]FileMetadata, error) {
	response, err := s.ParseToken(context.Background(), token)

	if err != nil {
		return nil, errors.New("could not parse token")
	}

	return s.repo.GetFiles(response.Claims.Uid)
}

func (s *service) ParseToken(ctx context.Context, token string) (*proto.ValidateResponse, error) {
	return s.jwtClient.ValidateToken(ctx, &proto.TokenRequest{Token: token})
}
