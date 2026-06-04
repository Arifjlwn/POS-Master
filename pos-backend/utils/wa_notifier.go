package utils

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func SendWhatsAppFonnte(token, targetPhone, message string) {
	go func() {
		if token == "" || targetPhone == "" { return }
		payload := new(bytes.Buffer); writer := multipart.NewWriter(payload)
		_ = writer.WriteField("target", targetPhone); _ = writer.WriteField("message", message); _ = writer.Close()
		req, err := http.NewRequest("POST", "https://api.fonnte.com/send", payload)
		if err != nil { log.Println("Gagal membuat request WA:", err); return }
		req.Header.Set("Authorization", token); req.Header.Set("Content-Type", writer.FormDataContentType())
		client := &http.Client{}; resp, err := client.Do(req); if err != nil { log.Println("Gagal mengirim WA:", err); return }
		defer resp.Body.Close(); body, _ := io.ReadAll(resp.Body); log.Println("Respon Fonnte:", string(body))
	}()
}

func SendSystemWhatsApp(targetPhone, message string) {
	systemToken := os.Getenv("WA_SYSTEM_TOKEN")
	if systemToken == "" { log.Println("⚠️ Peringatan: WA_SYSTEM_TOKEN di .env belum diisi!"); return }
	SendWhatsAppFonnte(systemToken, targetPhone, message)
}