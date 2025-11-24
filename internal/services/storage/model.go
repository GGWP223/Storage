package storage

type Request struct {
	FileID string `json:"file_id"`
	UserID string `json:"user_id"`
}

type TokenRequest struct {
	Token  string `json:"token"`
	FileID string `json:"file_id"`
}

type FileMetadata struct {
	FileID    string `json:"file_id"`
	UserID    string `json:"user_id"`
	FileName  string `json:"filename"`
	MimeType  string `json:"mime_type"`
	Size      int64  `json:"size"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
