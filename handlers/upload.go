package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/whaleLogic/googlecloud/storage"
)

// UploadHandler handles file upload requests
type UploadHandler struct {
	storageClient storage.StorageClient
}

// Response represents API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(storageClient storage.StorageClient) *UploadHandler {
	return &UploadHandler{
		storageClient: storageClient,
	}
}

// HandleUpload handles POST /upload requests
func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST method
	if r.Method != http.MethodPost {
		h.sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Get the file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		h.sendErrorResponse(w, "No file provided or invalid file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file
	if header.Size == 0 {
		h.sendErrorResponse(w, "Empty file not allowed", http.StatusBadRequest)
		return
	}

	// Check file size (max 100MB)
	maxSize := int64(100 << 20) // 100MB
	if header.Size > maxSize {
		h.sendErrorResponse(w, "File too large (max 100MB)", http.StatusBadRequest)
		return
	}

	// Validate file extension
	filename := header.Filename
	if !h.isAllowedFileType(filename) {
		h.sendErrorResponse(w, "File type not allowed", http.StatusBadRequest)
		return
	}

	// Upload to Google Cloud Storage
	result, err := h.storageClient.UploadFile(r.Context(), filename, file)
	if err != nil {
		h.sendErrorResponse(w, fmt.Sprintf("Upload failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	h.sendSuccessResponse(w, result)
}

// isAllowedFileType checks if the file type is allowed
func (h *UploadHandler) isAllowedFileType(filename string) bool {
	allowedExtensions := []string{
		".pdf", ".doc", ".docx", ".txt", ".rtf",
		".jpg", ".jpeg", ".png", ".gif", ".bmp",
		".zip", ".rar", ".tar", ".gz",
		".csv", ".xls", ".xlsx",
		".ppt", ".pptx",
	}

	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

// sendSuccessResponse sends a success response
func (h *UploadHandler) sendSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Success: true,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}

// sendErrorResponse sends an error response
func (h *UploadHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := Response{
		Success: false,
		Error:   message,
	}
	json.NewEncoder(w).Encode(response)
}