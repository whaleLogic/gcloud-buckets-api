package handlers

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/whaleLogic/googlecloud/storage"
)

// MockStorageClient implements a mock storage client for testing
type MockStorageClient struct {
	shouldFail bool
}

func (m *MockStorageClient) UploadFile(ctx context.Context, fileName string, reader io.Reader) (*storage.UploadResult, error) {
	if m.shouldFail {
		return nil, os.ErrInvalid
	}
	
	return &storage.UploadResult{
		FileName: "test-file.txt",
		URL:      "https://storage.googleapis.com/test-bucket/test-file.txt",
		Size:     100,
	}, nil
}

func (m *MockStorageClient) Close() error {
	return nil
}

func TestUploadHandler_HandleUpload(t *testing.T) {
	tests := []struct {
		name           string
		fileName       string
		fileContent    string
		shouldFail     bool
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid file upload",
			fileName:       "test.txt",
			fileContent:    "Hello, World!",
			shouldFail:     false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid file type",
			fileName:       "test.xyz",
			fileContent:    "Hello, World!",
			shouldFail:     false,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "File type not allowed",
		},
		{
			name:           "Empty file",
			fileName:       "test.txt",
			fileContent:    "",
			shouldFail:     false,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Empty file not allowed",
		},
		{
			name:           "Storage error",
			fileName:       "test.txt",
			fileContent:    "Hello, World!",
			shouldFail:     true,
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Upload failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock storage client
			mockClient := &MockStorageClient{shouldFail: tt.shouldFail}
			handler := NewUploadHandler(mockClient)

			// Create multipart form data
			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			
			if tt.name != "Empty file" {
				if tt.fileContent != "" {
					part, err := writer.CreateFormFile("file", tt.fileName)
					if err != nil {
						t.Fatal(err)
					}
					part.Write([]byte(tt.fileContent))
				}
			} else {
				// For empty file test, create a form file but with zero bytes
				part, err := writer.CreateFormFile("file", tt.fileName)
				if err != nil {
					t.Fatal(err)
				}
				// Don't write anything to simulate an empty file
				_ = part
			}
			writer.Close()

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// Create response recorder
			rec := httptest.NewRecorder()

			// Call handler
			handler.HandleUpload(rec, req)

			// Check status code
			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			// Check error message if expected
			if tt.expectedError != "" {
				body := rec.Body.String()
				if !contains(body, tt.expectedError) {
					t.Errorf("Expected error message containing '%s', got '%s'", tt.expectedError, body)
				}
			}
		})
	}
}

func TestUploadHandler_HandleUpload_MethodNotAllowed(t *testing.T) {
	mockClient := &MockStorageClient{}
	handler := NewUploadHandler(mockClient)

	req := httptest.NewRequest(http.MethodGet, "/upload", nil)
	rec := httptest.NewRecorder()

	handler.HandleUpload(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestUploadHandler_HandleUpload_Options(t *testing.T) {
	mockClient := &MockStorageClient{}
	handler := NewUploadHandler(mockClient)

	req := httptest.NewRequest(http.MethodOptions, "/upload", nil)
	rec := httptest.NewRecorder()

	handler.HandleUpload(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestIsAllowedFileType(t *testing.T) {
	handler := &UploadHandler{}

	tests := []struct {
		fileName string
		expected bool
	}{
		{"test.pdf", true},
		{"test.txt", true},
		{"test.jpg", true},
		{"test.png", true},
		{"test.doc", true},
		{"test.docx", true},
		{"test.zip", true},
		{"test.xyz", false},
		{"test.exe", false},
		{"test", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			result := handler.isAllowedFileType(tt.fileName)
			if result != tt.expected {
				t.Errorf("Expected %v for %s, got %v", tt.expected, tt.fileName, result)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}