#!/bin/bash

# Example script to test the Google Cloud Storage Upload API
# This script demonstrates how to upload a file using curl

API_URL="http://localhost:8080"

echo "=== Google Cloud Storage Upload API Test ==="
echo

# Test 1: Check if server is running
echo "1. Testing health endpoint..."
curl -s "${API_URL}/health" | jq '.' || echo "API not responding or jq not installed"
echo

# Test 2: Get API information
echo "2. Getting API information..."
curl -s "${API_URL}/" | jq '.' || echo "API not responding or jq not installed"
echo

# Test 3: Create a sample file and upload it
echo "3. Creating a sample file and uploading..."
TEMP_FILE="/tmp/sample_upload.txt"
echo "This is a sample file for testing the upload API." > "$TEMP_FILE"
echo "Created by: $(date)" >> "$TEMP_FILE"
echo "File size: $(wc -c < "$TEMP_FILE") bytes" >> "$TEMP_FILE"

echo "Uploading file: $TEMP_FILE"
RESPONSE=$(curl -s -X POST -F "file=@${TEMP_FILE}" "${API_URL}/upload")
echo "Response: $RESPONSE" | jq '.' || echo "$RESPONSE"

# Clean up
rm -f "$TEMP_FILE"
echo

echo "=== Test completed ==="
echo
echo "To run this script:"
echo "1. Make sure the API server is running: ./gcs-upload-api"
echo "2. Install jq for JSON formatting: sudo apt-get install jq"
echo "3. Run this script: ./test_api.sh"