package storage

import "time"

type Request struct {
	FileID string `json:"file_id"`
	UserID string `json:"user_id"`
}

type TokenRequest struct {
	Token  string `json:"token"`
	FileID string `json:"file_id"`
}

type FileMetadata struct {
	FileID    string    `json:"file_id"`
	UserID    string    `json:"user_id"`
	FileName  string    `json:"filename"`
	MimeType  string    `json:"mime_type"`
	Size      int64     `json:"size"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
