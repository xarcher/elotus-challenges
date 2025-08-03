package domain

import "time"

type FileUpload struct {
	ID          int       `json:"id" db:"id"`
	Filename    string    `json:"filename" db:"filename"`
	ContentType string    `json:"content_type" db:"content_type"`
	Size        int64     `json:"size" db:"size"`
	FilePath    string    `json:"file_path" db:"file_path"`
	UserAgent   string    `json:"user_agent" db:"user_agent"`
	RemoteAddr  string    `json:"remote_addr" db:"remote_addr"`
	UserID      int       `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UploadRepository interface {
	Create(upload *FileUpload) error
	GetByID(id int) (*FileUpload, error)
}

type UploadUsecase interface {
	UploadFile(userID int, filename string, contentType string, size int64,
		filePath string, userAgent string, remoteAddr string) (*FileUpload, error)
}
