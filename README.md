# Google Cloud Storage Upload API

A simple REST API built in Go for uploading documents to Google Cloud Storage.

## Features

- Upload files to Google Cloud Storage bucket
- Support for multiple file types (PDF, images, documents, archives)
- File size validation (max 100MB)
- File type validation
- Automatic file naming with timestamps
- CORS support for web applications
- Health check endpoint
- Graceful shutdown

## Prerequisites

1. Google Cloud Project with Storage API enabled
2. Google Cloud Storage bucket
3. Service account credentials with Storage Admin role

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `PORT` | Server port (default: 8080) | No |
| `GCS_BUCKET_NAME` | Google Cloud Storage bucket name | Yes |
| `GOOGLE_CLOUD_PROJECT` | Google Cloud Project ID | Yes |
| `GOOGLE_APPLICATION_CREDENTIALS` | Path to service account JSON file | Yes |

## Installation

1. Clone the repository:
```bash
git clone https://github.com/whaleLogic/googlecloud.git
cd googlecloud
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build -o gcs-upload-api .
```

## Configuration

### Google Cloud Setup

1. Create a Google Cloud Storage bucket:
```bash
gsutil mb gs://your-bucket-name
```

2. Create a service account:
```bash
gcloud iam service-accounts create gcs-uploader \
    --description="GCS Upload API Service Account" \
    --display-name="GCS Uploader"
```

3. Grant Storage Admin role:
```bash
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
    --member="serviceAccount:gcs-uploader@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/storage.admin"
```

4. Create and download service account key:
```bash
gcloud iam service-accounts keys create key.json \
    --iam-account=gcs-uploader@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

### Environment Setup

Create a `.env` file or set environment variables:
```bash
export GCS_BUCKET_NAME="your-bucket-name"
export GOOGLE_CLOUD_PROJECT="your-project-id"
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/key.json"
export PORT="8080"
```

## Usage

### Start the Server

```bash
./gcs-upload-api
```

The server will start on the configured port (default: 8080).

### API Endpoints

#### Upload File
- **URL**: `/upload`
- **Method**: `POST`
- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `file`: The file to upload (form field)

**Example using curl:**
```bash
curl -X POST \
  -F "file=@document.pdf" \
  http://localhost:8080/upload
```

**Success Response:**
```json
{
  "success": true,
  "data": {
    "fileName": "20240101-120000-document.pdf",
    "url": "https://storage.googleapis.com/your-bucket/20240101-120000-document.pdf",
    "size": 1234567
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "File type not allowed"
}
```

#### Health Check
- **URL**: `/health`
- **Method**: `GET`

**Response:**
```json
{
  "status": "healthy",
  "service": "gcs-upload-api"
}
```

#### API Information
- **URL**: `/`
- **Method**: `GET`

**Response:**
```json
{
  "service": "Google Cloud Storage Upload API",
  "version": "1.0.0",
  "endpoints": {
    "upload": "POST /upload",
    "health": "GET /health"
  },
  "usage": {
    "upload": "Send multipart form with 'file' field to /upload endpoint"
  }
}
```

### Supported File Types

- Documents: `.pdf`, `.doc`, `.docx`, `.txt`, `.rtf`
- Images: `.jpg`, `.jpeg`, `.png`, `.gif`, `.bmp`
- Archives: `.zip`, `.rar`, `.tar`, `.gz`
- Spreadsheets: `.csv`, `.xls`, `.xlsx`
- Presentations: `.ppt`, `.pptx`

### File Size Limits

- Maximum file size: 100MB
- Empty files are not allowed

## Example Client Code

### JavaScript (Browser)
```html
<input type="file" id="fileInput" />
<button onclick="uploadFile()">Upload</button>

<script>
async function uploadFile() {
    const fileInput = document.getElementById('fileInput');
    const file = fileInput.files[0];
    
    if (!file) {
        alert('Please select a file');
        return;
    }
    
    const formData = new FormData();
    formData.append('file', file);
    
    try {
        const response = await fetch('http://localhost:8080/upload', {
            method: 'POST',
            body: formData
        });
        
        const result = await response.json();
        
        if (result.success) {
            console.log('Upload successful:', result.data);
            alert(`File uploaded: ${result.data.url}`);
        } else {
            console.error('Upload failed:', result.error);
            alert(`Upload failed: ${result.error}`);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Upload error');
    }
}
</script>
```

### Python
```python
import requests

def upload_file(file_path):
    url = 'http://localhost:8080/upload'
    
    with open(file_path, 'rb') as f:
        files = {'file': f}
        response = requests.post(url, files=files)
    
    if response.status_code == 200:
        result = response.json()
        if result['success']:
            print(f"Upload successful: {result['data']['url']}")
        else:
            print(f"Upload failed: {result['error']}")
    else:
        print(f"HTTP error: {response.status_code}")

# Usage
upload_file('document.pdf')
```

## Docker Usage

### Build Docker Image
```bash
docker build -t gcs-upload-api .
```

### Run Container
```bash
docker run -p 8080:8080 \
  -e GCS_BUCKET_NAME="your-bucket-name" \
  -e GOOGLE_CLOUD_PROJECT="your-project-id" \
  -e GOOGLE_APPLICATION_CREDENTIALS="/app/key.json" \
  -v /path/to/key.json:/app/key.json:ro \
  gcs-upload-api
```

## Security Considerations

1. **Authentication**: This API doesn't include authentication. Consider adding API keys or OAuth for production use.
2. **CORS**: The API allows all origins (`*`). Restrict this in production.
3. **File Validation**: The API validates file types and sizes, but consider additional security scans.
4. **Rate Limiting**: Consider implementing rate limiting for production use.
5. **HTTPS**: Use HTTPS in production environments.

## Error Handling

The API returns appropriate HTTP status codes:
- `200`: Success
- `400`: Bad Request (invalid file, file too large, etc.)
- `405`: Method Not Allowed
- `500`: Internal Server Error

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the CC0 1.0 Universal License - see the LICENSE file for details.