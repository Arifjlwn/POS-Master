// 🚀 FIX MUTLAK: Samakan nama package dengan database.go agar sirkuit compiler Go lu normal kembali!
package src

// BusinessSubType mendefinisikan kesiapan sub-modul POS bray
type BusinessSubType struct {
	Name    string
	IsReady bool // 🚀 Gerbang LOCK/UNLOCK utama platform lu ada di sini!
}

// BusinessCluster menampung klaster industri global
type BusinessCluster struct {
	Name     string
	SubTypes map[string]BusinessSubType
}

// MasterAllowedIndustries adalah Single-Source-of-Truth klaster industri platform Arzura POS bray.
// Di sini kita kunci email user agar konsisten 1 industri dan amankan dari exploit lintas modul.
var MasterAllowedIndustries = map[string]BusinessCluster{
	"RETAIL": {
		Name: "Retail & Distribusi",
		SubTypes: map[string]BusinessSubType{
			"MINIMARKET": {Name: "Minimarket / Toko Kelontong", IsReady: true},
			"BOUTIQUE":   {Name: "Butik / Fashion", IsReady: true},
		},
	},
	"FNB": {
		Name: "Food & Beverage",
		SubTypes: map[string]BusinessSubType{
			"CAFE":       {Name: "Cafeteria / Kedai Kopi", IsReady: false},
			"RESTAURANT": {Name: "Restoran / Franchise", IsReady: false},
		},
	},
	"JASA": {
		Name: "Layanan Jasa",
		SubTypes: map[string]BusinessSubType{
			"LAUNDRY":    {Name: "Jasa Laundry Ruko", IsReady: true}, // 🧺 MODUL LAUNDRY LU LIVE KASTA TERTINGGI!
			"BARBER":     {Name: "Barbershop / Salon", IsReady: false},
			"BENGKEL":    {Name: "Bengkel / Car Wash", IsReady: false},
		},
	},
}