package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// SendOTPEmail mengirimkan kode OTP menggunakan HTTP API Resend dengan Template Corporate SaaS
func SendOTPEmail(targetEmail string, otp string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("konfigurasi RESEND_API_KEY di Environment Variable belum diisi")
	}

	// TAHAP LOGO: Silakan upload file PNG logo Anda ke hosting/CDN, lalu masukkan URL-nya di bawah ini.
	// Jika belum diupload, kode di bawah otomatis menggunakan teks teks fallback yang elegan.
	logoURL := "https://arzura-pos-dev.vercel.app/logo.png"

	currentYear := time.Now().Year()

	requestBody, err := json.Marshal(map[string]interface{}{
		"from":    "Arzura POS Official <no-reply@arzura-pos.my.id>",
		"to":      []string{targetEmail},
		"subject": "Keamanan Akun: Kode OTP Verifikasi Arzura POS",
		"html": fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
			</head>
			<body style="margin: 0; padding: 0; background-color: #f8fafc; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;">
				<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background-color: #f8fafc; padding: 40px 10px;">
					<tr>
						<td align="center">
							<table role="presentation" width="100%%" style="max-width: 480px; background-color: #ffffff; border: 1px solid #e2e8f0; border-radius: 16px; overflow: hidden; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);">
								
								<!-- HEADER / LOGO SECTION -->
								<tr>
									<td style="padding: 32px 32px 16px 32px; text-align: center;">
										<!-- Jika logoURL aktif, gambar ini akan muncul. Jika gagal load, teks ARZURA POS yang rapi akan muncul -->
										<img src="%s" alt="ARZURA POS" style="max-height: 45px; width: auto; border: 0; outline: none; text-decoration: none;" onerror="this.style.display='none'; document.getElementById('text-logo').style.display='block';">
										<div id="text-logo" style="display: none; font-size: 24px; font-weight: 800; letter-spacing: -0.5px; color: #0f172a;">
											ARZURA<span style="color: #475569; font-weight: 400;">POS</span>
										</div>
									</td>
								</tr>

								<!-- CONTENT SECTION -->
								<tr>
									<td style="padding: 16px 32px; color: #334155;">
										<h3 style="font-size: 20px; font-weight: 700; color: #0f172a; margin: 0 0 12px 0; text-align: center;">Verifikasi Kode Keamanan</h3>
										<p style="font-size: 14px; line-height: 24px; color: #475569; margin: 0 0 24px 0; text-align: center;">
											Kami mendeteksi permintaan verifikasi untuk masuk ke sistem kasir Anda. Gunakan kode One-Time Password (OTP) berikut untuk melanjutkan proses:
										</p>
										
										<!-- KODE OTP PREMIUM -->
										<div style="background-color: #f1f5f9; padding: 24px; text-align: center; border-radius: 12px; margin-bottom: 24px; border: 1px dashed #cbd5e1;">
											<span style="font-size: 36px; font-weight: 800; letter-spacing: 12px; color: #0f172a; font-family: 'Courier New', Courier, monospace; padding-left: 12px;">%s</span>
										</div>

										<p style="font-size: 13px; line-height: 20px; color: #64748b; margin: 0; text-align: center;">
											Demi keamanan sistem data kasir Anda, jangan bagikan kode ini kepada siapa pun termasuk tim teknis Arzura POS. <strong>Kode ini akan kedaluwarsa dalam waktu 5 menit.</strong>
										</p>
									</td>
								</tr>

								<!-- FOOTER / LEGAL SECTION -->
								<tr>
									<td style="padding: 24px 32px 32px 32px; background-color: #fafafa; border-top: 1px solid #f1f5f9; text-align: center;">
										<p style="font-size: 12px; line-height: 18px; color: #94a3b8; margin: 0 0 8px 0;">
											Pesan otomatis ini dikirim oleh sistem keamanan Arzura POS. Mohon untuk tidak membalas email ini secara langsung.
										</p>
										<p style="font-size: 11px; color: #cbd5e1; margin: 0;">
											&copy; %d Arzura POS Ecosystem. All rights reserved.<br>
											Infrastructure powered by Resend Enterprise API.
										</p>
									</td>
								</tr>

							</table>
						</td>
					</tr>
				</table>
			</body>
			</html>
		`, logoURL, otp, currentYear),
	})

	if err != nil {
		return fmt.Errorf("gagal membuat request body JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("gagal membuat request ke API Resend: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("gagal mengeksekusi request ke Resend: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("resend API mengembalikan status error: %d", resp.StatusCode)
	}

	fmt.Println("OTP Berhasil dikirim dengan standard Corporate SaaS! 🚀 Business is ready!")
	return nil
}
