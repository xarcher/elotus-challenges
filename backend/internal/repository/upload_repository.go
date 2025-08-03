package repository

import (
	"database/sql"

	"github.com/xarcher/backend/internal/domain"
)

type uploadRepository struct {
	db *sql.DB
}

func NewUploadRepository(db *sql.DB) domain.UploadRepository {
	return &uploadRepository{db: db}
}

func (r *uploadRepository) Create(upload *domain.FileUpload) error {
	query := `INSERT INTO file_uploads (filename, content_type, size, file_path, user_agent, remote_addr, user_id, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return r.db.QueryRow(query, upload.Filename, upload.ContentType, upload.Size,
		upload.FilePath, upload.UserAgent, upload.RemoteAddr,
		upload.UserID, upload.CreatedAt).Scan(&upload.ID)
}

func (r *uploadRepository) GetByID(id int) (*domain.FileUpload, error) {
	upload := &domain.FileUpload{}
	query := `SELECT id, filename, content_type, size, file_path, user_agent, remote_addr, user_id, created_at 
              FROM file_uploads WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&upload.ID, &upload.Filename, &upload.ContentType,
		&upload.Size, &upload.FilePath, &upload.UserAgent,
		&upload.RemoteAddr, &upload.UserID, &upload.CreatedAt)
	if err != nil {
		return nil, err
	}
	return upload, nil
}
