package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type SupabaseUploadResponse struct {
	Key string `json:"Key"`
}

func UploadToSupabase(file io.Reader, originalFilename string, contentType string, bucketName string, remotePath string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" {
		return "", fmt.Errorf("SUPABASE_URL belum diatur di konfigurasi .env")
	}
	if supabaseKey == "" {
		return "", fmt.Errorf("SUPABASE_KEY belum diatur di konfigurasi .env")
	}

	ext := strings.ToLower(filepath.Ext(originalFilename))
	if ext != "" && !strings.HasSuffix(strings.ToLower(remotePath), ext) {
		remotePath += ext
	}

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucketName, remotePath)

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", fmt.Errorf("gagal menyalin stream file ke buffer: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, uploadURL, buf)
	if err != nil {
		return "", fmt.Errorf("gagal merakit request HTTP: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal mengeksekusi request ke server storage: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload storage ditolak (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var uploadResp SupabaseUploadResponse
	_ = json.NewDecoder(resp.Body).Decode(&uploadResp)

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucketName, remotePath)
	return publicURL, nil
}
