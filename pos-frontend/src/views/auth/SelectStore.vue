<script setup>
import Swal from "sweetalert2";
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import api from "../../api.js";

const router = useRouter();
const stores = ref([]);
const userName = ref("");
const isLoading = ref(false);

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

// ==============================================================
// 🚀 STATE & DATA REAKTIF MODAL PRICING
// ==============================================================
const showPricingModal = ref(false);
const activePricingTab = ref("retail");

const industries = [
  {
    id: "retail",
    title: "Retail & Distribusi",
    desc: "Supermarket, Butik, Elektronik",
  },
  { id: "fnb", title: "Food & Beverage", desc: "Cafe, Restoran, Franchise" },
  { id: "jasa", title: "Layanan & Jasa", desc: "Laundry, Barbershop, Bengkel" },
];

const pricingPlans = {
  retail: [
    {
      id: "trial",
      name: "Starter Trial",
      price: "Rp 0",
      duration: "14 Hari",
      desc: "Validasi kesesuaian sistem dengan ekosistem bisnis Anda.",
      features: [
        "POS Kasir Retail",
        "Master Data Produk",
        "Scan Barcode Reader",
        "Tanpa Kartu Kredit",
      ],
    },
    {
      id: "basic",
      name: "Retail Basic",
      price: "49k",
      duration: "/Bulan",
      desc: "Solusi solid untuk toko dengan skala operasional ringan.",
      features: [
        "Kasir Tanpa Batas",
        "Manajemen Stok Dasar",
        "Struk Thermal Bluetooth",
        "Laporan Penjualan Harian",
      ],
    },
    {
      id: "pro",
      name: "Retail Pro",
      price: "149k",
      duration: "/Bulan",
      desc: "Cocok untuk minimarket yang mulai mengelola staf.",
      features: [
        "Semua Fitur Basic",
        "Manajemen Absensi & Shift",
        "Audit Stock Opname",
        "Laporan Ekspor (Excel/PDF)",
      ],
    },
    {
      id: "premium",
      name: "Retail Premium",
      price: "299k",
      duration: "/Bulan",
      desc: "Kendali penuh untuk bisnis multi-cabang & manajemen limbah.",
      features: [
        "Semua Fitur Pro",
        "Sistem Multi-Cabang (HO)",
        "Manajemen Retur & Waste",
        "Notifikasi WhatsApp",
      ],
    },
  ],
  fnb: [
    {
      id: "trial",
      name: "Starter Trial",
      price: "Rp 0",
      duration: "14 Hari",
      desc: "Validasi kesesuaian sistem dengan alur dapur Anda.",
      features: [
        "POS Kasir F&B",
        "Manajemen Menu & Kategori",
        "Hold/Simpan Pesanan",
        "Tanpa Kartu Kredit",
      ],
    },
    {
      id: "basic",
      name: "F&B Basic",
      price: "59k",
      duration: "/Bulan",
      desc: "Sistem operasional efisien untuk kedai atau coffee shop.",
      features: [
        "Manajemen Layout Meja",
        "Cetak Tiket Dapur (Kitchen)",
        "Pajak & Service Charge",
        "Struk Thermal Bluetooth",
      ],
    },
    {
      id: "pro",
      name: "F&B Pro",
      price: "169k",
      duration: "/Bulan",
      desc: "Untuk restoran dengan kontrol bahan baku ketat.",
      features: [
        "Semua Fitur Basic",
        "Resep Bahan Baku (BOM)",
        "Split Bill & Gabung Meja",
        "Manajemen Absensi & Shift",
      ],
    },
    {
      id: "premium",
      name: "F&B Premium",
      price: "349k",
      duration: "/Bulan",
      desc: "Skalabilitas franchise dengan analitik terpusat.",
      features: [
        "Semua Fitur Pro",
        "Manajemen Franchise (HO)",
        "Self-Order QR Menu",
        "Notifikasi WhatsApp",
      ],
    },
  ],
  jasa: [
    {
      id: "trial",
      name: "Starter Trial",
      price: "Rp 0",
      duration: "14 Hari",
      desc: "Coba modul layanan untuk bengkel, salon, atau laundry.",
      features: [
        "POS Layanan Jasa",
        "Database Pelanggan Dasar",
        "Penerimaan Pesanan",
        "Tanpa Kartu Kredit",
      ],
    },
    {
      id: "basic",
      name: "Service Basic",
      price: "49k",
      duration: "/Bulan",
      desc: "Sistem tracking pesanan yang rapi untuk bisnis jasa kecil.",
      features: [
        "Tracking Status Pesanan",
        "Cetak Nota / Resi Barcode",
        "Manajemen Layanan & Tarif",
        "Laporan Pendapatan",
      ],
    },
    {
      id: "pro",
      name: "Service Pro",
      price: "159k",
      duration: "/Bulan",
      desc: "Sistem otomatisasi performa teknisi dan staf.",
      features: [
        "Semua Fitur Basic",
        "Bagi Hasil / Komisi Karyawan",
        "Manajemen Absensi & Shift",
        "Audit Laporan Layanan",
      ],
    },
    {
      id: "premium",
      name: "Service Premium",
      price: "329k",
      duration: "/Bulan",
      desc: "Manajemen tingkat lanjut dengan pengingat otomatis.",
      features: [
        "Semua Fitur Pro",
        "Sistem Booking & Reservasi",
        "Sistem Multi-Cabang (HO)",
        "Notifikasi WhatsApp",
      ],
    },
  ],
};

onMounted(() => {
  const tempStores = localStorage.getItem("temp_stores");
  const name = localStorage.getItem("temp_name");

  // Jika memori kosong, tendang balik ke login
  if (!tempStores && !name) {
    router.push("/login");
    return;
  }

  // 🎯 REVISI: Cukup petakan datanya saja, JANGAN dipaksa showPricingModal.value = true
  stores.value = tempStores ? JSON.parse(tempStores) : [];
  userName.value = name || "Owner";
});

// 🚀 FILTER PUSAT VS CABANG
const pusatStore = computed(() => {
  if (stores.value.length === 0) return null;
  return stores.value.reduce(
    (min, store) => (store.id < min.id ? store : min),
    stores.value[0],
  );
});

const cabangStores = computed(() => {
  if (!pusatStore.value) return [];
  return stores.value.filter((store) => store.id !== pusatStore.value.id);
});

const selectBranch = async (storeId) => {
  isLoading.value = true;
  try {
    const res = await api.post("/auth/select-store", { store_id: storeId });

    localStorage.setItem("token", res.data.token);
    localStorage.setItem("store_id", res.data.store_id);
    localStorage.setItem("storeName", res.data.store_name || "POS UMKM");
    localStorage.setItem(
      "storeLogo",
      res.data.store_logo || res.data.logo_url || "",
    );
    localStorage.setItem(
      "subscriptionPlan",
      res.data.subscription_plan || "basic",
    );

    let finalRole = "owner";

    if (res.data.role) {
      finalRole = res.data.role.toLowerCase();
    } else if (res.data.user && res.data.user.role) {
      finalRole = res.data.user.role.toLowerCase();
    } else {
      // Jaga-jaga kalau backend tidak kirim, baca dari role pendaftaran awal
      const savedRole = localStorage.getItem("role");
      if (savedRole) finalRole = savedRole.toLowerCase();
    }
    localStorage.setItem("role", finalRole);

    localStorage.setItem(
      "name",
      res.data.name || res.data.user?.name || "Owner",
    );
    localStorage.setItem(
      "foto_url",
      res.data.foto_url || res.data.user?.foto_url || "",
    );

    router.push("/retail/pos/riwayat");
  } catch (error) {
    console.error("Gagal inisialisasi sesi cabang:", error);
    Swal.fire({
      icon: "error",
      title: "Akses Ditolak",
      text: error.response?.data?.error || "Gagal masuk ke cabang ini.",
      confirmButtonColor: "#ef4444",
      customClass: { popup: "rounded-[32px]" },
    });
  } finally {
    isLoading.value = false;
  }
};

// 🚀 ALUR UTAMA PILIH PAKET: Terbuka untuk Trial maupun komersial
const handlePilihPaketEkspansi = (industry, planId) => {
  localStorage.setItem("pendingIndustry", industry);
  localStorage.setItem("pendingPlan", planId);

  showPricingModal.value = false;

  // Kirim data langsung ke halaman setup pembuatan infrastruktur toko pertama / cabang baru
  router.push("/setup-toko?is_expansion=true");
};

const getPlanStyle = (plan) => {
  const p = plan ? plan.toLowerCase() : "basic";
  if (p === "premium")
    return "bg-amber-50 text-amber-700 border-amber-200 ring-amber-500";
  if (p === "pro")
    return "bg-indigo-50 text-indigo-700 border-indigo-200 ring-indigo-500";
  return "bg-sky-50 text-sky-700 border-sky-200 ring-sky-500";
};

const cleanLogoUrl = (url) => {
  if (!url) return "";
  let cleanPath = url
    .replace(/http:\/\/localhost:8080/g, "")
    .replace(new RegExp(API_BASE_URL, "g"), "");
  return `${API_BASE_URL}${cleanPath}`;
};

const getBranchLabel = (index) => {
  if (index === 0) return "Cabang Utama";
  return `Cabang Ke-${index + 1}`;
};
</script>

<template>
  <div
    class="min-h-screen bg-[#F8FAFC] flex flex-col items-center py-12 px-4 md:px-8 font-sans relative overflow-x-hidden select-none"
  >
    <!-- Ornamen -->
    <div
      class="absolute top-[-15%] left-[-15%] w-[45rem] h-[45rem] bg-gradient-to-br from-indigo-300/20 to-purple-300/10 rounded-full filter blur-[140px] pointer-events-none"
    ></div>
    <div
      class="absolute bottom-[-15%] right-[-15%] w-[45rem] h-[45rem] bg-gradient-to-tr from-blue-300/20 to-sky-300/10 rounded-full filter blur-[140px] pointer-events-none"
    ></div>

    <div class="w-full max-w-6xl z-10 flex flex-col items-center">
      <!-- HEADER -->
      <div class="text-center mb-16">
        <div
          class="inline-flex items-center justify-center px-4 py-2 bg-white rounded-2xl shadow-sm border border-slate-100 mb-6"
        >
          <div
            class="font-black text-2xl text-slate-900 tracking-tighter leading-none flex items-center gap-1.5"
          >
            NEXA
            <span
              class="text-indigo-600 bg-indigo-50 px-2 py-0.5 rounded-lg border border-indigo-100"
              >POS</span
            >
          </div>
        </div>
        <h1
          class="text-3xl md:text-5xl font-black text-slate-800 tracking-tight mb-4"
        >
          Selamat Datang, {{ userName.split(" ")[0] }}!
        </h1>
        <p
          class="text-slate-400 font-black text-[11px] md:text-xs uppercase tracking-[0.3em] bg-slate-100 border border-slate-200/60 px-4 py-1.5 rounded-full inline-block"
        >
          {{
            stores.length > 0
              ? "Topografi Infrastruktur Sesi & Kendali Operasional"
              : "Sesi Inisialisasi Akun SaaS Baru"
          }}
        </p>
      </div>

      <!-- STATE 1: JIKA USER SUDAH MEMILIKI TOKO (TAMPILKAN GRID TOPOLOGI) -->
      <div
        v-if="stores.length > 0"
        class="w-full flex flex-col items-center relative"
      >
        <!-- HQ CARD -->
        <div
          v-if="pusatStore"
          class="w-full max-w-md relative z-10 mb-20 group"
        >
          <div
            @click="selectBranch(pusatStore.id)"
            class="bg-slate-900 rounded-[36px] p-7 border-2 border-slate-800 shadow-2xl shadow-slate-900/20 hover:shadow-indigo-500/10 hover:-translate-y-2 transition-all duration-500 cursor-pointer flex flex-col relative overflow-hidden"
          >
            <div
              class="absolute top-0 right-0 bg-gradient-to-r from-amber-500 via-orange-500 to-amber-600 text-white font-black text-[9px] uppercase tracking-widest px-5 py-2.5 rounded-bl-3xl shadow-md border-b border-l border-white/10 animate-pulse"
            >
              HEADQUARTER / KANTOR PUSAT
            </div>
            <div class="flex items-center gap-5 mb-5">
              <div
                class="w-16 h-16 rounded-[20px] bg-white/10 border border-white/10 flex items-center justify-center overflow-hidden shrink-0 group-hover:rotate-6 transition-transform duration-500"
              >
                <img
                  v-if="pusatStore.logo_url"
                  :src="cleanLogoUrl(pusatStore.logo_url)"
                  class="w-full h-full object-cover"
                />
                <svg
                  v-else
                  class="w-7 h-7 text-slate-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                  />
                </svg>
              </div>
              <div class="flex flex-col">
                <span
                  class="text-[9px] font-black uppercase tracking-widest text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded-md border border-amber-500/20 self-start"
                >
                  {{ pusatStore.subscription_plan || "BASIC" }} ENTERPRISE
                </span>
                <h3
                  class="text-2xl font-black text-white mt-1.5 line-clamp-1 uppercase tracking-tight"
                >
                  {{ pusatStore.nama_toko }}
                </h3>
              </div>
            </div>
            <div
              class="flex items-center justify-between border-t border-white/5 pt-4 mt-2"
            >
              <span
                class="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center gap-1.5"
              >
                <div
                  class="w-1.5 h-1.5 bg-emerald-500 rounded-full animate-ping"
                ></div>
                {{ pusatStore.kota || "LOKASI BELUM DIATUR" }}
              </span>
              <span
                class="text-xs font-black text-amber-400 uppercase tracking-widest flex items-center gap-1 group-hover:text-white transition-colors"
              >
                KONTROL UTAMA
                <svg
                  class="w-4 h-4 group-hover:translate-x-1 transition-transform"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="3"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M14 5l7 7m0 0l-7 7m7-7H3"
                  />
                </svg>
              </span>
            </div>
          </div>
        </div>

        <!-- LINES ARCH -->
        <div
          v-if="cabangStores.length > 0"
          class="absolute top-32 bottom-28 w-1 bg-slate-200 pointer-events-none hidden lg:block z-0 overflow-hidden"
        >
          <div
            class="w-full h-24 bg-gradient-to-b from-transparent via-indigo-600 to-transparent absolute top-0 left-0 animate-data-flow"
          ></div>
        </div>

        <!-- GRID BANJAR CABANG -->
        <div
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 w-full mt-4 relative z-10"
        >
          <div
            v-for="(store, idx) in cabangStores"
            :key="store.id"
            @click="selectBranch(store.id)"
            class="group bg-white rounded-[32px] p-6 border border-slate-100 shadow-xs hover:shadow-2xl hover:-translate-y-2 transition-all duration-500 cursor-pointer flex flex-col relative overflow-hidden"
          >
            <div
              class="absolute top-0 left-0 w-full h-1.5 transition-all duration-300"
              :class="
                getPlanStyle(store.subscription_plan)
                  .split(' ')[3]
                  .replace('ring-', 'bg-')
              "
            ></div>
            <div class="flex items-start justify-between mb-5">
              <div
                class="w-14 h-14 rounded-2xl bg-slate-50 border border-slate-100 flex items-center justify-center overflow-hidden shrink-0"
              >
                <img
                  v-if="store.logo_url"
                  :src="cleanLogoUrl(store.logo_url)"
                  class="w-full h-full object-cover"
                />
                <svg
                  v-else
                  class="w-6 h-6 text-slate-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                  />
                </svg>
              </div>
              <span
                class="text-[10px] font-black text-indigo-600 uppercase tracking-widest bg-indigo-50 border border-indigo-100/50 px-3 py-1 rounded-xl"
              >
                {{ getBranchLabel(idx) }}
              </span>
            </div>
            <div class="flex-1">
              <h3
                class="text-xl font-black text-slate-800 mb-1 line-clamp-1 uppercase tracking-tight"
              >
                {{ store.nama_toko }}
              </h3>
              <p
                class="text-xs font-bold text-slate-400 uppercase tracking-wider line-clamp-1 flex items-center gap-1"
              >
                <svg
                  class="w-3.5 h-3.5 text-slate-300"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="2.5"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                  />
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                </svg>
                {{ store.kota || "Lokasi Belum Diatur" }}
              </p>
            </div>
            <div
              class="mt-6 flex items-center justify-between border-t border-slate-50 pt-4"
            >
              <div class="flex flex-col">
                <span
                  class="text-[8px] font-black text-slate-300 uppercase tracking-widest leading-none"
                  >Infrastruktur</span
                >
                <span
                  class="text-[10px] font-black text-slate-500 uppercase tracking-wider mt-0.5"
                  >Terminal Active</span
                >
              </div>
              <div
                class="text-xs font-black text-slate-700 uppercase tracking-widest flex items-center gap-0.5 group-hover:text-indigo-600 transition-colors"
              >
                Buka Sesi
                <svg
                  class="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="3"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M14 5l7 7m0 0l-7 7m7-7H3"
                  />
                </svg>
              </div>
            </div>
          </div>

          <!-- BUTTON JIKA MAU EKSPANSI LAGI -->
          <div
            @click="showPricingModal = true"
            class="bg-white rounded-[32px] p-6 border-2 border-dashed border-slate-200 hover:border-indigo-400 hover:bg-indigo-50/10 transition-all duration-500 cursor-pointer flex flex-col items-center justify-center text-center min-h-[200px] group shadow-2xs"
          >
            <div
              class="w-12 h-12 rounded-full bg-slate-50 border border-slate-200/60 group-hover:bg-indigo-100 flex items-center justify-center mb-3 transition-all"
            >
              <svg
                class="w-5 h-5 text-slate-400 group-hover:text-indigo-600"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="3"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M12 4v16m8-8H4"
                />
              </svg>
            </div>
            <h3
              class="text-sm font-black text-slate-600 group-hover:text-indigo-800 uppercase tracking-wider mb-0.5"
            >
              Ekspansi Jaringan
            </h3>
            <p
              class="text-[11px] font-bold text-slate-400 px-4 leading-relaxed"
            >
              Daftarkan node cabang baru ke dalam ekosistem topografi bisnis
              Anda.
            </p>
          </div>
        </div>
      </div>

      <!-- STATE 2: JIKA TOKO KOSONG (USER BARU DAFTAR ACC) -->
      <div
        v-else
        class="w-full max-w-xl text-center py-12 px-6 bg-white rounded-[40px] border border-slate-100 shadow-xl animate-fade-in"
      >
        <div
          class="w-20 h-20 bg-indigo-50 border border-indigo-100 text-indigo-600 rounded-[28px] flex items-center justify-center mx-auto mb-6 shadow-xs"
        >
          <svg
            class="w-10 h-10 animate-bounce"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
            />
          </svg>
        </div>
        <h2
          class="text-2xl font-black text-slate-800 uppercase tracking-tight mb-2"
        >
          Infrastruktur Toko Belum Tersedia
        </h2>
        <p
          class="text-slate-400 font-bold text-xs md:text-sm leading-relaxed mb-8 px-4"
        >
          Akun Anda berhasil diaktifkan secara global. Silakan tentukan kluster
          industri dan pilih paket modul untuk menginisialisasi server toko
          pertama Anda.
        </p>
        <button
          @click="showPricingModal = true"
          class="w-full py-4 bg-indigo-600 hover:bg-slate-900 text-white rounded-2xl font-black text-xs uppercase tracking-widest shadow-xl shadow-indigo-100 transition-all duration-300 active:scale-95"
        >
          Buat Toko Pertama Anda
        </button>
      </div>
    </div>

    <!-- ============================================================== -->
    <!-- 🚀 OVERLAY POPUP MODAL: PRICING AKTIF (SINKRON SAMA LANDING PAGE) -->
    <!-- ============================================================== -->
    <div
      v-if="showPricingModal"
      class="fixed inset-0 bg-slate-900/60 backdrop-blur-md z-[100] flex items-center justify-center p-4 md:p-6 overflow-y-auto animate-fade-in"
    >
      <div
        class="bg-white w-full max-w-6xl rounded-[40px] shadow-2xl p-6 md:p-10 relative border border-slate-100 my-auto max-h-[92vh] overflow-y-auto custom-scrollbar flex flex-col items-center"
      >
        <!-- Tombol Tutup (Hanya nampil kalau udah punya toko biar gak kejebak blank screen) -->
        <button
          v-if="stores.length > 0"
          @click="showPricingModal = false"
          class="absolute top-6 right-6 w-10 h-10 bg-slate-50 hover:bg-rose-50 text-slate-400 hover:text-rose-500 rounded-full flex items-center justify-center transition-all"
        >
          <svg
            class="w-5 h-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2.5"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>

        <div class="text-center mb-10">
          <h2
            class="text-3xl md:text-4xl font-black text-slate-800 tracking-tight"
          >
            INVESTASI <span class="text-indigo-600">TRANSPARAN</span>
          </h2>
          <p
            class="text-slate-400 font-bold text-xs md:text-sm mt-2 max-w-2xl uppercase tracking-wider"
          >
            PILIH MODUL INDUSTRI ANDA, DAN TEMUKAN SKALABILITAS YANG DIRANCANG
            KHUSUS UNTUK BISNIS ANDA.
          </p>
        </div>

        <!-- FILTER TAB MODUL INDUSTRI REAKTIF -->
        <div
          class="bg-white/80 backdrop-blur-md p-1.5 rounded-[24px] border border-slate-200/80 flex w-full max-w-xl shadow-xl shadow-slate-100 mb-12 shrink-0"
        >
          <button
            v-for="ind in industries"
            :key="ind.id"
            @click="activePricingTab = ind.id"
            :class="
              activePricingTab === ind.id
                ? 'bg-slate-900 text-white shadow-lg shadow-slate-900/20'
                : 'text-slate-500 hover:text-slate-800 hover:bg-slate-50'
            "
            class="flex-1 px-2 py-3.5 rounded-2xl font-black text-[10px] uppercase tracking-widest transition-all text-center duration-300"
          >
            {{ ind.title }}
          </button>
        </div>

        <!-- GRID PLAN PAKET -->
        <div
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 w-full items-stretch"
        >
          <div
            v-for="plan in pricingPlans[activePricingTab]"
            :key="plan.id"
            class="p-8 rounded-[32px] border transition-all duration-300 flex flex-col relative group"
            :class="[
              plan.id === 'trial'
                ? 'bg-white border-slate-200 shadow-sm border-2 border-dashed'
                : '',
              plan.id === 'basic'
                ? 'bg-white border-slate-200 shadow-sm hover:border-sky-300 hover:shadow-xl shadow-sky-50/50'
                : '',
              plan.id === 'pro'
                ? 'bg-white border-2 border-indigo-600 shadow-2xl shadow-indigo-100 lg:scale-105 z-10'
                : '',
              plan.id === 'premium'
                ? 'bg-slate-900 border-slate-800 shadow-2xl shadow-slate-950/20 z-0'
                : '',
            ]"
          >
            <div
              v-if="plan.id === 'pro'"
              class="absolute -top-4 left-1/2 -translate-x-1/2 bg-gradient-to-r from-indigo-600 to-blue-600 text-white px-5 py-2 rounded-full text-[9px] font-black uppercase tracking-widest shadow-xl whitespace-nowrap animate-pulse"
            >
              REKOMENDASI UMKM
            </div>
            <div
              v-if="plan.id === 'premium'"
              class="absolute -top-4 left-1/2 -translate-x-1/2 bg-gradient-to-r from-amber-500 to-amber-600 text-slate-950 px-5 py-2 rounded-full text-[9px] font-black uppercase tracking-widest shadow-xl whitespace-nowrap"
            >
              FITUR ENTERPRISE
            </div>

            <div class="mb-8 mt-2 flex-1">
              <h3
                class="font-black text-[11px] uppercase tracking-[0.2em] mb-4"
                :class="[
                  plan.id === 'trial' ? 'text-indigo-500' : '',
                  plan.id === 'basic' ? 'text-sky-500' : '',
                  plan.id === 'pro' ? 'text-indigo-600' : '',
                  plan.id === 'premium' ? 'text-amber-400' : '',
                ]"
              >
                {{ plan.name }}
              </h3>
              <div class="flex items-baseline gap-1">
                <span
                  class="text-4xl lg:text-5xl font-black tracking-tighter"
                  :class="
                    plan.id === 'premium' ? 'text-white' : 'text-slate-900'
                  "
                  >{{ plan.price }}</span
                >
                <span
                  class="font-bold text-[10px] uppercase tracking-widest mb-1 ml-1"
                  :class="
                    plan.id === 'premium' ? 'text-slate-400' : 'text-slate-400'
                  "
                  >{{ plan.duration }}</span
                >
              </div>
              <p
                class="text-xs font-bold mt-5 h-12 leading-relaxed"
                :class="
                  plan.id === 'premium' ? 'text-slate-400' : 'text-slate-500'
                "
              >
                {{ plan.desc }}
              </p>

              <ul
                class="space-y-4 mb-2 mt-6 border-t pt-6"
                :class="
                  plan.id === 'premium'
                    ? 'border-slate-800'
                    : 'border-slate-100'
                "
              >
                <li
                  v-for="feat in plan.features"
                  :key="feat"
                  class="flex items-start gap-3 text-xs font-bold leading-tight"
                  :class="
                    plan.id === 'premium' ? 'text-slate-300' : 'text-slate-700'
                  "
                >
                  <svg
                    class="w-4 h-4 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="3"
                    :class="[
                      plan.id === 'premium' ? 'text-amber-400' : '',
                      plan.id === 'pro' ? 'text-indigo-500' : '',
                      plan.id === 'basic' ? 'text-sky-500' : '',
                      plan.id === 'trial' ? 'text-indigo-500' : '',
                    ]"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                  {{ feat }}
                </li>
              </ul>
            </div>

            <!-- 🚀 BUTTON SEKARANG FULL AKTIF TERMASUK UNTUK TRIAL -->
            <button
              @click="handlePilihPaketEkspansi(activePricingTab, plan.id)"
              :class="[
                plan.id === 'trial'
                  ? 'bg-indigo-50 hover:bg-indigo-600 text-indigo-600 hover:text-white border border-indigo-100'
                  : '',
                plan.id === 'basic'
                  ? 'bg-slate-900 text-white hover:bg-sky-600 shadow-lg'
                  : '',
                plan.id === 'pro'
                  ? 'bg-gradient-to-r from-indigo-600 to-blue-600 text-white hover:from-slate-900 hover:to-slate-900 shadow-xl'
                  : '',
                plan.id === 'premium'
                  ? 'bg-gradient-to-r from-amber-400 to-amber-500 text-slate-950 font-extrabold hover:from-white hover:to-white shadow-xl'
                  : '',
              ]"
              class="block w-full text-center py-4 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all duration-300 transform active:scale-95"
            >
              {{
                plan.id === "trial" ? "MULAI EKSPLORASI" : "PILIH " + plan.name
              }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- OVERLAY LOADING -->
    <div
      v-if="isLoading"
      class="fixed inset-0 z-[150] bg-slate-900/50 backdrop-blur-sm flex items-center justify-center"
    >
      <div
        class="bg-white p-7 rounded-[24px] shadow-2xl flex flex-col items-center border border-slate-100"
      >
        <div
          class="w-10 h-10 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin mb-4"
        ></div>
        <div
          class="text-xs font-black text-slate-600 uppercase tracking-widest animate-pulse"
        >
          Menghubungkan Sesi Node...
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.98) translateY(10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}
@keyframes dataFlow {
  0% {
    top: -10%;
  }
  100% {
    top: 110%;
  }
}
.animate-data-flow {
  animation: dataFlow 3s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}
.custom-scrollbar::-webkit-scrollbar {
  width: 5px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 10px;
}
</style>
