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

type SupabaseUploadResponse struct { Key string `json:"Key"` }

func UploadToSupabase(file io.Reader, originalFilename string, contentType string, bucketName string, remotePath string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL"); supabaseToken := os.Getenv("SUPABASE_TOKEN")
	if supabaseURL == "" { return "", fmt.Errorf("SUPABASE_URL belum diatur") }
	if supabaseToken == "" { return "", fmt.Errorf("SUPABASE_TOKEN belum diatur") }
	ext := strings.ToLower(filepath.Ext(originalFilename))
	if ext != "" && !strings.HasSuffix(strings.ToLower(remotePath), ext) { remotePath += ext }
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucketName, remotePath)
	buf := new(bytes.Buffer); if _, err := io.Copy(buf, file); err != nil { return "", err }
	req, err := http.NewRequest(http.MethodPost, uploadURL, buf); if err != nil { return "", err }
	req.Header.Set("Authorization", "Bearer "+supabaseToken); req.Header.Set("apikey", supabaseToken)
	req.Header.Set("Content-Type", contentType); req.Header.Set("x-upsert", "true")
	client := &http.Client{}; resp, err := client.Do(req); if err != nil { return "", err }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body); return "", fmt.Errorf("supabase upload failed (%d): %s", resp.StatusCode, string(bodyBytes))
	}
	var uploadResp SupabaseUploadResponse; _ = json.NewDecoder(resp.Body).Decode(&uploadResp)
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucketName, remotePath), nil
}