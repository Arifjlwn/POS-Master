package usecase

import (
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pos-backend/models"
	"pos-backend/src/modules/jasalayanan/laundry/domain"
	"pos-backend/src/modules/jasalayanan/laundry/repository"
)

type LaundryUseCase interface {
	ProcessCheckout(storeID, userID uint, input domain.CheckoutLaundryInput) (string, string, error)
	ProcessPelunasan(trxID, storeID uint, input domain.PelunasanInput) error
	GetLaporanRingkasan(storeID uint) (domain.ReportSummaryResponse, error)
	UpdateStatusCucian(detailID uint, status string) error
}

type laundryUseCase struct {
	repo repository.LaundryRepository
}

func NewLaundryUseCase(repo repository.LaundryRepository) LaundryUseCase {
	return &laundryUseCase{repo: repo}
}

func (u *laundryUseCase) saveBase64Image(base64Data, folder, filename string) (string, error) {
	if base64Data == "" {
		return "", nil
	}
	parts := strings.Split(base64Data, ",")
	var pureBase64 string
	if len(parts) > 1 {
		pureBase64 = parts[1]
	} else {
		pureBase64 = parts[0]
	}
	decodedData, err := base64.StdEncoding.DecodeString(pureBase64)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return "", err
	}
	filePath := filepath.Join(folder, filename)
	if err := os.WriteFile(filePath, decodedData, 0644); err != nil {
		return "", err
	}
	return strings.ReplaceAll(filePath, "\\", "/"), nil
}

func (u *laundryUseCase) ProcessCheckout(storeID, userID uint, input domain.CheckoutLaundryInput) (string, string, error) {
	estimasiTime, err := time.Parse("2006-01-02", input.EstimasiSelesai)
	if err != nil {
		estimasiTime = time.Now().Add(time.Hour * 48)
	}

	invoiceCode := fmt.Sprintf("INV/LD/%s/%s", time.Now().Format("20060102"), time.Now().Format("150405"))

	var buktiPath, fotoBarangPath string
	if input.PaymentMethod == "QRIS" && input.BuktiTransferBase64 != "" {
		buktiPath, _ = u.saveBase64Image(input.BuktiTransferBase64, "public/uploads/qris", strings.ReplaceAll(invoiceCode, "/", "")+".jpg")
	}
	if input.FotoBarangBase64 != "" {
		fotoBarangPath, _ = u.saveBase64Image(input.FotoBarangBase64, "public/uploads/items", strings.ReplaceAll(invoiceCode, "/", "")+".jpg")
	}

	db := u.repo.GetDB()
	tx := db.Begin()

	newTx := models.Transaction{
		SessionID:     1,
		StoreID:       storeID,
		UserID:        userID,
		NoInvoice:     invoiceCode,
		SubTotal:      input.TotalAmount,
		Pajak:         0,
		Pembulatan:    0,
		TotalHarga:    input.TotalAmount,
		MetodeBayar:   input.PaymentMethod,
		StatusBayar:   input.PaymentStatus,
		NominalBayar:  input.TotalAmount,
		Kembalian:     0,
		BuktiTransfer: buktiPath,
	}

	if err := u.repo.CreateTransactionTx(tx, &newTx); err != nil {
		tx.Rollback()
		return "", "", fmt.Errorf("gagal membuat invoice induk: %v", err)
	}

	for _, item := range input.Items {
		laundryDetail := domain.TransactionLaundryDetail{
			TransactionID: newTx.ID,
			StoreID:       storeID,
			ProductID:     item.ProductID,
			NamaPelanggan: input.CustomerName,
			NoWhatsapp:    input.CustomerPhone,
			BeratKg:       item.BeratKg,
			HargaPerKg:    item.HargaPerKg,
			SubTotal:      item.SubTotal,
			StatusCucian:  "ANTRI",
			StatusBayar:   input.PaymentStatus,
			MetodeBayar:   input.PaymentMethod,
			EstimasiWaktu: estimasiTime,
			FotoBarang:    fotoBarangPath,
			NamaParfum:    item.NamaParfum,
			HargaParfum:   item.HargaParfum,
		}

		if err := u.repo.CreateLaundryDetailTx(tx, &laundryDetail); err != nil {
			tx.Rollback()
			return "", "", fmt.Errorf("gagal menyimpan rincian cucian: %v", err)
		}
	}

	existingCustomer, err := u.repo.FindCustomerByPhone(storeID, input.CustomerPhone)
	if err != nil {
		newCustomer := models.Customer{
			StoreID:    storeID,
			Nama:       input.CustomerName,
			NoWhatsapp: input.CustomerPhone,
		}
		tx.Create(&newCustomer)
	} else {
		existingCustomer.Nama = input.CustomerName
		tx.Save(existingCustomer)
	}

	tx.Commit()

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	fotoPublicUrl := ""
	if fotoBarangPath != "" {
		fotoPublicUrl = baseURL + "/" + fotoBarangPath
	}

	return invoiceCode, fotoPublicUrl, nil
}

func (u *laundryUseCase) ProcessPelunasan(trxID, storeID uint, input domain.PelunasanInput) error {
	trx, err := u.repo.GetTransactionByID(trxID, storeID)
	if err != nil {
		return fmt.Errorf("transaksi tidak ditemukan")
	}
	if trx.StatusBayar == "LUNAS" {
		return fmt.Errorf("transaksi ini sudah lunas sebelumnya")
	}
	if input.MetodeBayar == "QRIS" && input.BuktiTransferBase64 != "" {
		buktiPath, _ := u.saveBase64Image(input.BuktiTransferBase64, "public/uploads/qris", strings.ReplaceAll(trx.NoInvoice, "/", "")+"_lunas.jpg")
		trx.BuktiTransfer = buktiPath
	}
	trx.StatusBayar = "LUNAS"
	trx.MetodeBayar = input.MetodeBayar
	return u.repo.UpdateTransaction(trx)
}

func (u *laundryUseCase) GetLaporanRingkasan(storeID uint) (domain.ReportSummaryResponse, error) {
	transactions, err := u.repo.GetAllTransactions(storeID)
	if err != nil {
		return domain.ReportSummaryResponse{}, err
	}

	var reportData []domain.TransactionReportResponse
	var tunai, qris, debit, piutang, omset float64
	totalOrder := len(transactions)

	for _, trx := range transactions {
		if trx.StatusBayar == "BELUM_LUNAS" || trx.StatusBayar == "BELUM_BAYAR" {
			piutang += trx.TotalHarga
		} else {
			omset += trx.TotalHarga
			switch strings.ToUpper(trx.MetodeBayar) {
			case "TUNAI", "CASH":
				tunai += trx.TotalHarga
			case "QRIS":
				qris += trx.TotalHarga
			case "DEBIT":
				debit += trx.TotalHarga
			}
		}

		detail, err := u.repo.GetLaundryDetailByTxID(trx.ID)
		var namaPelanggan, noWhatsapp string
		var productID uint
		var beratKg, subTotalDetail float64
		var estimasiWaktu time.Time

		if err == nil {
			namaPelanggan = detail.NamaPelanggan
			noWhatsapp = detail.NoWhatsapp
			productID = detail.ProductID
			beratKg = detail.BeratKg
			subTotalDetail = detail.SubTotal
			estimasiWaktu = detail.EstimasiWaktu
		}

		layananName := "Paket Laundry"
		if productID > 0 {
			if prod, err := u.repo.GetProductByIDSimple(productID); err == nil {
				layananName = prod.NamaProduk
			}
		}

		reportData = append(reportData, domain.TransactionReportResponse{
			Transaction:   trx,
			Invoice:       trx.NoInvoice,
			Pelanggan:     namaPelanggan,
			Whatsapp:      noWhatsapp,
			Layanan:       layananName,
			BeratKg:       beratKg,
			SubTotal:      subTotalDetail,
			EstimasiWaktu: estimasiWaktu,
		})
	}

	avg := 0.0
	if totalOrder > 0 {
		avg = omset / float64(totalOrder)
	}

	return domain.ReportSummaryResponse{
		Ringkasan: domain.ReportSummary{
			TotalOmset: omset,
			TotalOrder: totalOrder,
			RataRata:   math.Round(avg*100) / 100,
			Tunai:      tunai,
			Qris:       qris,
			Debit:      debit,
			Piutang:    piutang,
		},
		Transaksi: reportData,
	}, nil
}

func (u *laundryUseCase) UpdateStatusCucian(detailID uint, status string) error {
	return u.repo.UpdateStatusDetailCucian(detailID, status)
}