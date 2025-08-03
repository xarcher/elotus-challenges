package handler

import (
	"github.com/xarcher/backend/internal/domain"
	"github.com/xarcher/backend/pkg/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type UploadHandler struct {
	uploadUsecase domain.UploadUsecase
}

func NewUploadHandler(uploadUsecase domain.UploadUsecase) *UploadHandler {
	return &UploadHandler{
		uploadUsecase: uploadUsecase,
	}
}

func (h *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse multipart form (32MB max memory)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	// Get file from form
	file, handler, err := r.FormFile("data")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// Check file size (8MB limit)
	const maxSize = 8 * 1024 * 1024
	if handler.Size > maxSize {
		utils.RespondError(w, http.StatusBadRequest, "File size exceeds 8MB limit")
		return
	}

	// Check content type
	contentType := handler.Header.Get("Content-Type")
	if contentType == "" {
		// Try to detect content type
		buffer := make([]byte, 512)
		_, err := file.Read(buffer)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Unable to read file")
			return
		}
		contentType = http.DetectContentType(buffer)
		file.Seek(0, 0) // Reset file pointer
	}

	// Create temp file
	tempFile, err := os.CreateTemp("/tmp", "upload_*"+filepath.Ext(handler.Filename))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Unable to create temp file")
		return
	}
	defer tempFile.Close()

	// Copy file content
	_, err = io.Copy(tempFile, file)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Unable to save file")
		return
	}

	// Save metadata to database
	upload, err := h.uploadUsecase.UploadFile(
		userID,
		handler.Filename,
		contentType,
		handler.Size,
		tempFile.Name(),
		r.UserAgent(),
		r.RemoteAddr,
	)
	if err != nil {
		// Clean up temp file on error
		os.Remove(tempFile.Name())
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, upload)
}

func (h *UploadHandler) ServeUploadForm(w http.ResponseWriter, r *http.Request) {
	html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>File Upload</title>
    </head>
    <body>
        <h2>Upload Image File</h2>
        <form action="/upload" method="post" enctype="multipart/form-data">
            <input type="file" name="data" accept="image/*" required>
            <br><br>
            <input type="submit" value="Upload">
        </form>
    </body>
    </html>
    `
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
