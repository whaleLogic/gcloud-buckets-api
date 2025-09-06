package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/whaleLogic/googlecloud/config"
	"github.com/whaleLogic/googlecloud/handlers"
	"github.com/whaleLogic/googlecloud/storage"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate required configuration
	if cfg.BucketName == "" {
		log.Fatal("GCS_BUCKET_NAME environment variable is required")
	}

	if cfg.ProjectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT environment variable is required")
	}

	// Create context
	ctx := context.Background()

	// Initialize Google Cloud Storage client
	storageClient, err := storage.NewClient(ctx, cfg.ProjectID, cfg.BucketName)
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	defer storageClient.Close()

	// Initialize handlers
	uploadHandler := handlers.NewUploadHandler(storageClient)

	// Set up HTTP routes
	mux := http.NewServeMux()
	
	// Upload endpoint
	mux.HandleFunc("/upload", uploadHandler.HandleUpload)
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "healthy", "service": "gcs-upload-api"}`)
	})

	// Root endpoint with API information
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"service": "Google Cloud Storage Upload API",
			"version": "1.0.0",
			"endpoints": {
				"upload": "POST /upload",
				"health": "GET /health"
			},
			"usage": {
				"upload": "Send multipart form with 'file' field to /upload endpoint"
			}
		}`)
	})

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		log.Printf("Using GCS bucket: %s", cfg.BucketName)
		log.Printf("Using project ID: %s", cfg.ProjectID)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}