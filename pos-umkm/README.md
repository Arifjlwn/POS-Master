# 🚀 POS-UMKM (SaaS Point of Sale)
Sistem Kasir berbasis web yang dirancang khusus untuk efisiensi operasional UMKM, dibangun menggunakan stack modern Laravel 12 dan Vue 3.

---

## 🛠️ Tech Stack
- **Framework:** Laravel 12
- **Frontend:** Vue 3 (Composition API)
- **State Management:** Inertia.js
- **Styling:** Tailwind CSS
- **Database:** MySQL
- **Charts:** Chart.js (untuk Analitik Dashboard)

---

## 📌 Fitur yang Sudah Selesai (DONE)

### 🏗️ Arsitektur & Keamanan
- **RBAC (Role Based Access Control):** Pembedaan akses antara **Owner** dan **Kasir**.
- **Smart Redirect:** - Owner otomatis diarahkan ke Dashboard Analitik setelah login.
  - Kasir otomatis diarahkan ke halaman Absensi setelah login.
- **Store-Specific Data:** Filter data berdasarkan `store_id` agar data antar cabang/toko tidak bercampur.

### 💰 Operasional Toko
- **POS (Point of Sale):** Keranjang belanja dinamis, kalkulasi pembayaran, dan cetak struk.
- **Master Produk:** Manajemen katalog barang lengkap dengan sistem **Barcode/SKU**.
- **Mesin Absensi:** Sistem kehadiran karyawan mandiri dengan log waktu masuk dan pulang.

### 📊 Administrasi & Analisis (Owner Only)
- **Dashboard Analitik:** Visualisasi tren penjualan harian dan 5 produk paling laris (Fast Moving).
- **Manajemen Karyawan:** Input data kasir otomatis menghasilkan NIK sebagai ID login.
- **Export Data:** Fitur tarik laporan absensi dalam format CSV.

---

## 🚀 Rencana Pengembangan (TO-DO)

### 🏗️ Fase 1: Data Integrity
- [ ] Sinkronisasi angka dashboard dengan data riil dari tabel `transactions`.
- [ ] Implementasi Database Trigger/Logic untuk pengurangan stok otomatis saat checkout sukses.

### 📦 Fase 2: Inventori & Logistik
- [ ] Fitur **Stock-In**: Pencatatan barang masuk dari supplier.
- [ ] **Stock Alert**: Notifikasi di dashboard untuk produk dengan stok di bawah ambang batas (Stok Kritis).

### 📈 Fase 3: Financial Insights
- [ ] Laporan Laba/Rugi (Gross Profit vs Net Profit).
- [ ] Filter performa penjualan berdasarkan periode waktu tertentu.

### ⚙️ Fase 4: Optimization
- [ ] Integrasi Printer Thermal via USB/Bluetooth.
- [ ] Sistem manajemen diskon (Persentase atau nominal tetap).

---

## 👨‍💻 Author
**Arif Juliawan**
*Student ID: 055756684*
*Sistem Informasi - Universitas Terbuka*

---
*Last Updated: 8 Mei 2026*
