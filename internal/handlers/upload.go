package handlers

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadHandler struct {
	DB *sql.DB
}

type CloudinaryResponse struct {
	SecureURL string `json:"secure_url"`
	PublicID  string `json:"public_id"`
	Error     *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	if productID == "" {
		http.Error(w, `{"error":"Product ID required"}`, http.StatusBadRequest)
		return
	}

	// 5 MB max upload
	maxSize := int64(5 * 1024 * 1024)
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	if err := r.ParseMultipartForm(maxSize); err != nil {
		http.Error(w, `{"error":"File size exceeds 5 MB limit"}`, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, `{"error":"Image file is required"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		http.Error(w, `{"error":"Invalid file format. Only JPG, PNG, and WebP are allowed"}`, http.StatusBadRequest)
		return
	}

	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	var imageURL string

	if cloudName != "" && apiKey != "" && apiSecret != "" {
		// Upload to Cloudinary
		cloudURL, err := uploadToCloudinary(file, header.Filename, cloudName, apiKey, apiSecret)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Cloudinary upload failed: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		imageURL = cloudURL
	} else {
		// Local storage fallback for development
		filename := fmt.Sprintf("%s-%d%s", productID, time.Now().UnixNano(), ext)
		destPath := filepath.Join("frontend", "assets", "images", filename)
		os.MkdirAll(filepath.Dir(destPath), 0755)

		out, err := os.Create(destPath)
		if err != nil {
			http.Error(w, `{"error":"Failed to save image locally"}`, http.StatusInternalServerError)
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			http.Error(w, `{"error":"Failed to write image data"}`, http.StatusInternalServerError)
			return
		}

		// Also copy to root assets/images if present
		rootPath := filepath.Join("assets", "images", filename)
		os.MkdirAll(filepath.Dir(rootPath), 0755)
		if rf, err := os.Create(rootPath); err == nil {
			file.Seek(0, io.SeekStart)
			io.Copy(rf, file)
			rf.Close()
		}

		imageURL = "/assets/images/" + filename
	}

	// Store ONLY the URL in product_images
	var maxOrder int
	h.DB.QueryRow("SELECT COALESCE(MAX(display_order), -1) FROM product_images WHERE product_id = $1", productID).Scan(&maxOrder)

	_, err = h.DB.Exec(`
		INSERT INTO product_images (product_id, url, display_order)
		VALUES ($1, $2, $3)
	`, productID, imageURL, maxOrder+1)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Failed to record image in database: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"url":     imageURL,
	})
}

func uploadToCloudinary(file multipart.File, filename, cloudName, apiKey, apiSecret string) (string, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	// Generate SHA1 signature: "timestamp=...<api_secret>"
	sigStr := fmt.Sprintf("timestamp=%s%s", timestamp, apiSecret)
	h := sha1.New()
	h.Write([]byte(sigStr))
	signature := hex.EncodeToString(h.Sum(nil))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	writer.WriteField("api_key", apiKey)
	writer.WriteField("timestamp", timestamp)
	writer.WriteField("signature", signature)
	writer.Close()

	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", cloudName)
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var cloudRes CloudinaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&cloudRes); err != nil {
		return "", err
	}

	if cloudRes.Error != nil {
		return "", fmt.Errorf(cloudRes.Error.Message)
	}

	if cloudRes.SecureURL == "" {
		return "", fmt.Errorf("no secure URL returned from Cloudinary")
	}

	return cloudRes.SecureURL, nil
}
