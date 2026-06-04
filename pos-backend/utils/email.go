package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendOTPEmail(targetEmail string, otp string) error {
	from := os.Getenv("SMTP_FROM"); if from == "" { from = "no-reply@nexapos.com" }
	username := os.Getenv("SMTP_USERNAME"); if username == "" { username = "da1bed098977f5" }
	password := os.Getenv("SMTP_PASSWORD"); if password == "" { password = "3f5706f134b6eb" }
	smtpHost := os.Getenv("SMTP_HOST"); if smtpHost == "" { smtpHost = "sandbox.smtp.mailtrap.io" }
	smtpPort := os.Getenv("SMTP_PORT"); if smtpPort == "" { smtpPort = "2525" }
	subject := "Subject: [NEXAPOS] Kode Verifikasi Akun\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`<div style="font-family: sans-serif; max-width: 400px; margin: auto; border: 1px solid #eee; padding: 20px; border-radius: 20px;"><h2 style="color: #2563eb; text-align: center;">NEXA<span style="color: #1e293b;">POS</span></h2><p style="text-align: center; color: #64748b;">Gunakan kode OTP di bawah untuk verifikasi:</p><div style="font-size: 40px; font-weight: 900; letter-spacing: 10px; text-align: center; color: #2563eb; padding: 20px; background: #f8fafc; border-radius: 15px; margin: 20px 0;">%s</div><p style="font-size: 10px; color: #94a3b8; text-align: center; text-transform: uppercase; letter-spacing: 1px;">Kode berlaku 5 menit</p></div>`, otp)
	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{targetEmail}, []byte(subject+mime+body))
	if err != nil { fmt.Println("Gagal kirim email SMTP:", err); return err }
	fmt.Println("OTP Berhasil terkirim via SMTP Provider! 🚀"); return nil
}