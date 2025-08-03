package usecase

import (
	"errors"
	"github.com/xarcher/backend/config"
	"strings"
	"time"

	"github.com/xarcher/backend/internal/domain"
)

type uploadUsecase struct {
	uploadRepo domain.UploadRepository
	uploadCfg  config.UploadConfig
	timeout    time.Duration
}

func NewUploadUsecase(uploadRepo domain.UploadRepository, uploadCfg config.UploadConfig, timeout time.Duration) domain.UploadUsecase {
	return &uploadUsecase{
		uploadRepo: uploadRepo,
		uploadCfg:  uploadCfg,
		timeout:    timeout,
	}
}

func (u *uploadUsecase) UploadFile(userID int, filename string, contentType string,
	size int64, filePath string, userAgent string,
	remoteAddr string) (*domain.FileUpload, error) {

	// Validate content type
	if !strings.HasPrefix(contentType, "image/") {
		return nil, errors.New("file must be an image")
	}

	// Validate size (8MB limit)
	const maxSize = 8 * 1024 * 1024
	if size > maxSize {
		return nil, errors.New("file size exceeds 8MB limit")
	}

	upload := &domain.FileUpload{
		Filename:    filename,
		ContentType: contentType,
		Size:        size,
		FilePath:    filePath,
		UserAgent:   userAgent,
		RemoteAddr:  remoteAddr,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	if err := u.uploadRepo.Create(upload); err != nil {
		return nil, err
	}

	return upload, nil
}
