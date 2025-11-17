package storage

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	SaveFile(file FileMetadata) error
	DeleteFile(file FileMetadata) error
	GetFile(request Request) (FileMetadata, error)
	GetFiles(userID string) ([]FileMetadata, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) SaveFile(file FileMetadata) error {
	request := Request{
		file.FileID,
		file.UserID,
	}

	_, err := r.GetFile(request)

	if err != nil {
		return r.db.Create(&file).Error
	}

	return r.db.Save(&file).Error
}

func (r *repository) DeleteFile(file FileMetadata) error {
	if file.FileID == "" || file.UserID == "" {
		return errors.New("missing file ID or user ID for delete")
	}

	return r.db.Where("file_id = ? AND user_id = ?", file.FileID, file.UserID).Delete(&FileMetadata{}).Error
}

func (r *repository) GetFile(request Request) (FileMetadata, error) {
	var file FileMetadata

	err := r.db.Where("file_id = ? AND user_id = ?", request.FileID, request.UserID).First(&file).Error

	return file, err
}

func (r *repository) GetFiles(userID string) ([]FileMetadata, error) {
	var files []FileMetadata

	err := r.db.Where("user_id = ?", userID).Find(&files).Error

	if err != nil {
		return nil, errors.New("could not find files: " + err.Error())
	}

	return files, nil
}
