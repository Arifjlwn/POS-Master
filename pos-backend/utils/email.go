package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Inisialisasi HTTP Client secara global dengan timeout 10 detik.
// Ini mencegah kebocoran soket TCP (saturasi port) saat traffic pendaftaran tinggi.
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// SendOTPEmail mengirimkan kode OTP menggunakan HTTP API Resend dengan Template Corporate SaaS
func SendOTPEmail(targetEmail string, otp string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("konfigurasi RESEND_API_KEY di Environment Variable belum diisi")
	}

	// URL logo default
	logoURL := "https://arzura-pos-dev.vercel.app/logo.png"
	currentYear := time.Now().Year()

	requestBody, err := json.Marshal(map[string]interface{}{
		"from":    "Arzura POS Official <no-reply@arzura-pos.my.id>",
		"to":      []string{targetEmail},
		"subject": "Kode Keamanan: Verifikasi Login Arzura POS",
		"html": fmt.Sprintf(`
            <!DOCTYPE html>
            <html>
            <head>
                <meta charset="UTF-8">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
            </head>
            <body style="margin: 0; padding: 0; background-color: #f1f5f9; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;">
                <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background-color: #f1f5f9; padding: 40px 10px;">
                    <tr>
                        <td align="center">
                            <table role="presentation" width="100%%" style="max-width: 500px; background-color: #ffffff; border: 1px solid #e2e8f0; border-top: 6px solid #4f46e5; border-radius: 16px; overflow: hidden; box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.05);">
                                
                                <tr>
                                    <td style="padding: 40px 32px 10px 32px; text-align: center;">
                                        <img src="%s" alt="ARZURA POS" style="max-height: 48px; width: auto; border: 0; outline: none; text-decoration: none;" onerror="this.style.display='none'; document.getElementById('text-logo').style.display='block';">
                                        <div id="text-logo" style="display: none; font-size: 26px; font-weight: 900; letter-spacing: -0.5px; color: #0f172a;">
                                            ARZURA<span style="color: #4f46e5; font-style: italic;">POS</span>
                                        </div>
                                    </td>
                                </tr>

                                <tr>
                                    <td style="padding: 10px 32px 20px 32px; color: #334155; text-align: center;">
                                        <h2 style="font-size: 22px; font-weight: 800; color: #0f172a; margin: 0 0 12px 0; text-transform: uppercase; letter-spacing: 0.5px;">Kode Verifikasi Anda</h2>
                                        <p style="font-size: 14px; line-height: 24px; color: #64748b; margin: 0 0 24px 0;">
                                            Kami menerima permintaan login ke sistem kasir Anda. Gunakan kode <strong>One-Time Password (OTP)</strong> berikut untuk memvalidasi akses:
                                        </p>
                                        
                                        <div style="background: linear-gradient(135deg, #4f46e5 0%%, #7c3aed 100%%); padding: 24px 16px; border-radius: 16px; margin-bottom: 24px; box-shadow: 0 4px 15px rgba(79, 70, 229, 0.2);">
                                            <div style="font-size: 42px; font-weight: 900; letter-spacing: 14px; color: #ffffff; font-family: 'Courier New', Courier, monospace; text-shadow: 0 2px 4px rgba(0,0,0,0.1); padding-left: 14px;">
                                                %s
                                            </div>
                                        </div>

                                        <div style="background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 12px; margin-bottom: 20px;">
                                            <p style="font-size: 13px; font-weight: 700; color: #ef4444; margin: 0; display: flex; justify-content: center; align-items: center;">
                                            Kode ini akan kedaluwarsa dalam waktu 3 menit.
                                            </p>
                                        </div>

                                        <p style="font-size: 13px; line-height: 20px; color: #94a3b8; margin: 0;">
                                            Demi keamanan, <strong>jangan pernah membagikan kode ini</strong> kepada siapa pun, termasuk staf operasional maupun tim teknis Arzura.
                                        </p>
                                    </td>
                                </tr>

                                <tr>
                                    <td style="padding: 24px 32px; background-color: #f8fafc; border-top: 1px solid #f1f5f9; text-align: center;">
                                        <p style="font-size: 12px; font-weight: 600; color: #64748b; margin: 0 0 8px 0; text-transform: uppercase; letter-spacing: 1px;">
                                            Keamanan Transaksi Anda Prioritas Kami
                                        </p>
                                        <p style="font-size: 11px; line-height: 18px; color: #cbd5e1; margin: 0;">
                                            Email ini dihasilkan secara otomatis oleh gateway keamanan.<br>
                                            &copy; %d Arzura POS Premium Ecosystem. All rights reserved.
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

	// Menggunakan httpClient global yang aman dari memory leak / TCP saturation
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("gagal mengeksekusi request ke Resend (Network Timeout/Drop): %v", err)
	}
	defer resp.Body.Close()

	// Membaca response body untuk mengantisipasi silent error dari Resend
	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("resend API error (Status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	fmt.Println("OTP Berhasil dieksekusi via Resend API !")
	return nil
}
