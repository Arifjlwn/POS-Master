package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// UploadToSupabase mengirimkan file gambar langsung ke Cloud Storage Supabase
func UploadToSupabase(file io.Reader, headerFilename string, contentType string, bucketName string, customFileName string) (string, error) {
	// 🎯 AMBIL DARI .ENV (Jauh lebih aman dan profesional!)
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseToken := os.Getenv("SUPABASE_TOKEN")

	// Validasi jaga-jaga kalau .env lupa diisi
	if supabaseURL == "" || supabaseToken == "" {
		return "", fmt.Errorf("konfigurasi SUPABASE_URL atau SUPABASE_TOKEN belum diatur di .env")
	}

	// Tentukan nama file unik di cloud
	ext := filepath.Ext(headerFilename)
	remotePath := fmt.Sprintf("%s%s", customFileName, ext)
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucketName, remotePath)

	// Baca file buffer
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	// Siapkan HTTP Request
	req, err := http.NewRequest("POST", uploadURL, buf)
	if err != nil {
		return "", err
	}

	// Inject Header Keamanan Supabase
	req.Header.Set("Authorization", "Bearer "+supabaseToken)
	req.Header.Set("apikey", supabaseToken)
	req.Header.Set("Content-Type", contentType)

	// Eksekusi Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("supabase error status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Return URL Publik Gambar
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucketName, remotePath)
	return publicURL, nil
}