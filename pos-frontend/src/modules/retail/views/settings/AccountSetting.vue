<script setup>
import Swal from "sweetalert2";
import { onMounted, ref } from "vue";
import api from "../../../../api.js";
import Sidebar from "../../components/Sidebar.vue";
import AccountProfile from "../../components/settings/account/AccountProfile.vue";
import AccountSecurity from "../../components/settings/account/AccountSecurity.vue";
import { useAccount } from "../../composables/useAccount.js";

const { isLoading, isSaving, activeTab, role, profileForm, passwordForm, fotoPreview, handleFileChange, saveProfile, updatePassword } = useAccount();

const tabs = [
  {
    id: "profile",
    label: "Profil Saya",
    icon: "M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z",
  },
  {
    id: "security",
    label: "Keamanan",
    icon: "M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z",
  },
];

// ==========================================
// 🚀 STATE KHUSUS BILLING SAAS
// ==========================================
const subPlan = ref("Basic");
const subStatus = ref("");
const subEnd = ref("");
const isBillingLoading = ref(true);
const showUpgradeModal = ref(false);
const isExpired = ref(false);
const loadingPayment = ref(false);
const quotaTerminal = ref(1);

// ==========================================
// 🚀 JURUS SINKRONISASI OTOMATIS (ANTI MANUAL)
// ==========================================
const sinkronisasiStatusBerlangganan = async () => {
  isBillingLoading.value = true;
  try {
    const res = await api.get("/retail/store/settings");
    const storeData = res.data.data;

    // Amankan state lokal
    subPlan.value = storeData.subscription_plan || "Basic";
    subStatus.value = storeData.subscription_status || "inactive";
    quotaTerminal.value = storeData.quota_terminal || 1;

    if (storeData.subscription_plan) {
      localStorage.setItem("subscriptionPlan", storeData.subscription_plan.toLowerCase());
    }

    if (storeData.subscription_end) {
      const dateObj = new Date(storeData.subscription_end);
      subEnd.value = dateObj.toLocaleDateString("id-ID", {
        day: "numeric",
        month: "long",
        year: "numeric",
      });

      const today = new Date();
      const diffDays = Math.ceil((dateObj.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));

      if (diffDays > 0 && subStatus.value === "active") {
        isExpired.value = false; // Buka gembok layar full screen otomatis!
      } else {
        isExpired.value = true;
      }
    } else {
      subEnd.value = "Tidak Terbatas";
      isExpired.value = subStatus.value !== "active";
    }

    // 🚀 Teriak ke seluruh aplikasi agar Sidebar & Composable ikut ter-update realtime!
    window.dispatchEvent(new Event("store-updated"));
    window.dispatchEvent(new Event("storage"));
  } catch (error) {
    console.error("Gagal melakukan verifikasi otomatis data langganan:", error);
  } finally {
    isBillingLoading.value = false;
  }
};

onMounted(async () => {
  // Suntik Script Midtrans
  if (!document.getElementById("midtrans-script-owner")) {
    const midtransEnv = import.meta.env.VITE_MIDTRANS_ENV || "sandbox";
    const snapUrl = midtransEnv === "production" ? "https://app.midtrans.com/snap/snap.js" : "https://app.sandbox.midtrans.com/snap/snap.js";

    const script = document.createElement("script");
    script.id = "midtrans-script-owner";
    script.src = snapUrl;
    script.setAttribute("data-client-key", import.meta.env.VITE_MIDTRANS_CLIENT_KEY);
    document.head.appendChild(script);
  }

  // Jalankan load data awal
  await sinkronisasiStatusBerlangganan();
});

const handleUpgrade = async (planName, price) => {
  showUpgradeModal.value = false;

  Swal.fire({
    title: "Menyiapkan Tagihan...",
    text: "Mohon tunggu sebentar",
    allowOutsideClick: false,
    didOpen: () => {
      Swal.showLoading();
    },
  });

  try {
    const res = await api.post("/retail/subscription/upgrade", {
      plan_name: planName,
      price: price,
    });

    const snapToken = res.data.token;
    Swal.close();

    window.snap.pay(snapToken, {
      onSuccess: async function (result) {
        Swal.fire({
          title: "Memverifikasi Pembayaran...",
          text: "Menghubungkan ke server finansial, mohon jangan tutup halaman ini.",
          icon: "info",
          allowOutsideClick: false,
          didOpen: () => {
            Swal.showLoading();
          },
        });

        // 🚀 Kasih delay toleransi 1.5 detik agar webhook backend golang kelar memproses invoice
        await new Promise((resolve) => setTimeout(resolve, 1500));

        // Panggil fungsi sinkronisasi otomatis
        await sinkronisasiStatusBerlangganan();
        Swal.close();

        Swal.fire({
          title: "Pembayaran Sukses!",
          text: `Selamat! Toko Anda kini telah aktif menggunakan paket ${planName.toUpperCase()}.`,
          icon: "success",
          timer: 2500,
          showConfirmButton: false,
        });

        activeTab.value = "billing";
      },
      onPending: function (result) {
        Swal.fire("Menunggu Pembayaran", "Silakan selesaikan pembayaran Anda.", "info");
      },
      onError: function (result) {
        Swal.fire("Pembayaran Gagal", "Terjadi kesalahan saat memproses pembayaran.", "error");
      },
      onClose: function () {
        if (isExpired.value) {
          Swal.fire("Pembayaran Tertunda", "Silakan selesaikan invoice Anda kapan saja untuk membuka gembok toko.", "warning");
        }
      },
    });
  } catch (error) {
    Swal.fire("Gagal!", "Server sedang sibuk, coba lagi nanti.", "error");
  }
};

// 🚀 FUNGSI BELI LISENSI KASIR (Add-on)
const beliLisensiTambahan = async () => {
  loadingPayment.value = true;
  try {
    const response = await api.post("/retail/subscription/upgrade", {
      plan_name: "Terminal Tambahan",
      price: 50000,
    });

    const snapToken = response.data.token;

    window.snap.pay(snapToken, {
      onSuccess: async function (result) {
        Swal.fire({
          title: "Memproses Penambahan Perangkat...",
          allowOutsideClick: false,
          didOpen: () => {
            Swal.showLoading();
          },
        });

        await new Promise((resolve) => setTimeout(resolve, 1500));
        await sinkronisasiStatusBerlangganan();

        Swal.close();
        Swal.fire("Berhasil!", "Pembayaran sukses! Kuota Mesin Kasir Anda telah bertambah.", "success");
        loadingPayment.value = false;
      },
      onPending: function (result) {
        Swal.fire("Pending", "Menunggu pembayaran diselesaikan...", "info");
        loadingPayment.value = false;
      },
      onError: function (result) {
        Swal.fire("Gagal!", "Pembayaran gagal diproses.", "error");
        loadingPayment.value = false;
      },
      onClose: function () {
        loadingPayment.value = false;
      },
    });
  } catch (error) {
    Swal.fire("Gagal!", error.response?.data?.error || "Gagal menghubungi Payment Gateway.", "error");
    loadingPayment.value = false;
  }
};
</script>

<template>
  <Sidebar>
    <div v-if="isExpired" class="fixed inset-0 z-[9999] bg-slate-900 flex flex-col items-center justify-center p-4 md:p-8 overflow-y-auto">
      <div class="text-center mt-32 md:mt-10 mb-10 max-w-2xl mx-auto animate-pulse">
        <div class="w-20 h-20 bg-rose-500/20 text-rose-500 rounded-full flex items-center justify-center mx-auto mb-6 shadow-[0_0_40px_rgba(244,63,94,0.3)]">
          <svg class="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
            />
          </svg>
        </div>

        <h1 class="text-3xl md:text-5xl text-white font-black uppercase tracking-tight mb-4">Masa Aktif Berakhir</h1>

        <p class="text-slate-400 font-bold tracking-widest text-xs md:text-sm uppercase">
          Operasional sistem POS dihentikan sementara.

          <br />

          Silakan perpanjang paket untuk membuka kembali akses toko Anda.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl w-full mb-20">
        <div class="bg-slate-800 rounded-[24px] p-6 border border-slate-700 flex flex-col text-center">
          <div class="text-[10px] font-black text-sky-400 uppercase tracking-widest mb-1">Retail Basic</div>

          <h3 class="text-2xl font-black text-white uppercase mb-3">Basic</h3>

          <div class="flex items-baseline justify-center gap-1 mb-6">
            <span class="text-sm font-black text-slate-400">Rp</span>

            <span class="text-5xl font-black text-white">
              49

              <span class="text-lg text-slate-500">k</span>
            </span>
          </div>

          <button
            @click="handleUpgrade('Basic', 49000)"
            class="w-full py-4 bg-sky-500 text-white rounded-xl font-black text-[10px] uppercase tracking-widest shadow-lg hover:bg-sky-400 transition-all active:scale-95"
          >
            Perpanjang Basic
          </button>
        </div>

        <div
          class="bg-slate-800 rounded-[24px] p-6 border-2 border-indigo-500 shadow-xl shadow-indigo-500/20 relative flex flex-col text-center transform md:-translate-y-4"
        >
          <div
            class="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-indigo-500 text-white text-[9px] font-black uppercase tracking-widest px-4 py-1.5 rounded-full"
          >
            Paling Laris
          </div>

          <div class="text-[10px] font-black text-indigo-400 uppercase tracking-widest mb-1 mt-2">Retail Pro</div>

          <h3 class="text-2xl font-black text-white uppercase mb-3">Pro</h3>

          <div class="flex items-baseline justify-center gap-1 mb-6">
            <span class="text-sm font-black text-slate-400">Rp</span>

            <span class="text-5xl font-black text-white">
              149

              <span class="text-lg text-slate-500">k</span>
            </span>
          </div>

          <button
            @click="handleUpgrade('Pro', 149000)"
            class="w-full py-4 bg-indigo-500 text-white rounded-xl font-black text-[10px] uppercase tracking-widest shadow-lg hover:bg-indigo-400 transition-all active:scale-95"
          >
            Perpanjang Pro
          </button>
        </div>

        <div class="bg-slate-800 rounded-[24px] p-6 border border-slate-700 flex flex-col text-center">
          <div class="text-[10px] font-black text-amber-400 uppercase tracking-widest mb-1">Retail Premium</div>

          <h3 class="text-2xl font-black text-white uppercase mb-3">Premium</h3>

          <div class="flex items-baseline justify-center gap-1 mb-6">
            <span class="text-sm font-black text-slate-400">Rp</span>

            <span class="text-5xl font-black text-white">
              299

              <span class="text-lg text-slate-500">k</span>
            </span>
          </div>

          <button
            @click="handleUpgrade('Premium', 299000)"
            class="w-full py-4 bg-amber-500 text-slate-900 rounded-xl font-black text-[10px] uppercase tracking-widest shadow-lg hover:bg-amber-400 transition-all active:scale-95"
          >
            Perpanjang Premium
          </button>
        </div>
      </div>
    </div>

    <div v-if="!isExpired" class="p-4 md:p-8 lg:p-10 max-w-5xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 class="text-2xl md:text-3xl font-black text-slate-800 tracking-tighter uppercase">Akun Saya</h1>

          <p class="text-[10px] md:text-xs font-black text-slate-400 uppercase tracking-widest mt-1">Kelola Informasi Pribadi dan Keamanan Akun</p>
        </div>
      </div>

      <div v-if="isLoading" class="py-20 flex flex-col items-center justify-center bg-white rounded-[32px] border border-slate-100 shadow-sm">
        <div class="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>

        <p class="text-slate-400 font-black text-[10px] uppercase tracking-widest animate-pulse">Menyiapkan Profil...</p>
      </div>

      <div v-else class="bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden flex flex-col lg:flex-row min-h-[500px]">
        <div class="w-full lg:w-64 bg-slate-50/50 border-r border-slate-100 p-6 flex flex-row lg:flex-col gap-2 overflow-x-auto custom-scrollbar shrink-0">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'flex items-center gap-3 p-3.5 rounded-2xl font-black text-[10px] uppercase tracking-widest transition-all whitespace-nowrap',

              activeTab === tab.id ? 'bg-white text-indigo-600 shadow-sm border border-slate-200' : 'text-slate-500 hover:bg-slate-100',
            ]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" :d="tab.icon" />
            </svg>

            {{ tab.label }}
          </button>

          <button
            v-if="role === 'owner'"
            @click="activeTab = 'billing'"
            :class="[
              'flex items-center gap-3 p-3.5 rounded-2xl font-black text-[10px] uppercase tracking-widest transition-all whitespace-nowrap mt-4 border-t border-slate-200 pt-4',

              activeTab === 'billing' ? 'bg-white text-amber-600 shadow-sm border border-slate-200' : 'text-amber-500 hover:bg-amber-50',
            ]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
            </svg>

            Paket Berlangganan
          </button>
        </div>

        <div class="flex-1 p-6 md:p-8 overflow-y-auto relative custom-scrollbar">
          <div v-if="isSaving" class="absolute inset-0 bg-white/60 backdrop-blur-sm z-10 flex items-center justify-center">
            <div class="w-8 h-8 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin"></div>
          </div>

          <AccountProfile
            v-show="activeTab === 'profile'"
            :form="profileForm"
            :fotoPreview="fotoPreview"
            @update-file="handleFileChange"
            @submit="saveProfile"
          />

          <AccountSecurity v-show="activeTab === 'security'" :form="passwordForm" @submit="updatePassword" />

          <div v-show="activeTab === 'billing'" class="animate-fade-in-up w-full">
            <div v-if="isBillingLoading" class="text-center py-10">
              <div class="w-8 h-8 border-4 border-amber-200 border-t-amber-500 rounded-full animate-spin mx-auto mb-4"></div>

              <p class="text-slate-400 font-black text-[10px] uppercase tracking-widest animate-pulse">Mengecek Data Berlangganan...</p>
            </div>

            <div v-else>
              <div class="bg-slate-50 border border-slate-100 p-8 rounded-[32px] text-center shadow-sm">
                <div class="w-16 h-16 bg-indigo-100 text-indigo-600 rounded-2xl flex items-center justify-center mx-auto mb-4 border border-indigo-200">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                </div>

                <h3 class="text-xs font-black text-slate-400 uppercase tracking-widest mb-1">Paket Saat Ini</h3>

                <div class="text-2xl font-black text-slate-800 uppercase tracking-widest mb-6">
                  PAKET {{ subPlan }}

                  <span
                    v-if="subStatus === 'active'"
                    class="inline-block align-middle ml-2 px-3 py-1 bg-emerald-100 text-emerald-600 text-[10px] rounded-lg border border-emerald-200"
                  >
                    AKTIF
                  </span>

                  <span v-else class="inline-block align-middle ml-2 px-3 py-1 bg-rose-100 text-rose-600 text-[10px] rounded-lg border border-rose-200">
                    TIDAK AKTIF
                  </span>
                </div>

                <div class="bg-white p-5 rounded-[20px] border border-slate-200 text-sm mb-6 max-w-sm mx-auto">
                  <p class="font-bold text-slate-500 text-xs">Masa aktif berlangganan Anda sampai dengan:</p>

                  <p class="text-lg font-black text-indigo-600 mt-1">
                    {{ subEnd }}
                  </p>
                </div>

                <button
                  @click="showUpgradeModal = true"
                  class="max-w-sm w-full mx-auto bg-slate-900 text-white py-4 rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-indigo-600 transition-all shadow-lg active:scale-95 flex items-center justify-center gap-2"
                >
                  Perpanjang / Upgrade Paket
                </button>
              </div>

              <div class="mt-8 border-t border-slate-100 pt-8">
                <div class="mb-6">
                  <h3 class="text-xl font-black text-slate-800 uppercase tracking-tighter italic">Lisensi & Add-ons</h3>

                  <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">Kelola Kapasitas Operasional Toko Anda</p>
                </div>

                <div class="max-w-md">
                  <!-- CARD TERMINAL KASIR -->

                  <div
                    class="bg-indigo-50/50 border-2 border-indigo-100 p-6 rounded-[24px] flex flex-col justify-between transition-all hover:border-indigo-300"
                  >
                    <div>
                      <div class="flex items-center justify-between">
                        <h4 class="font-black text-indigo-900 text-sm uppercase tracking-wider">Terminal Kasir</h4>

                        <span class="bg-indigo-200 text-indigo-800 text-[9px] font-black px-2 py-1 rounded-md uppercase tracking-widest">Active</span>
                      </div>

                      <p class="text-[11px] font-bold text-indigo-600 mt-2 leading-relaxed">
                        Menentukan jumlah perangkat (HP/Tablet/PC) yang bisa membuka laci kasir secara bersamaan.
                      </p>

                      <div class="mt-4 flex items-end gap-1">
                        <span class="text-4xl font-black text-indigo-900 leading-none">{{ quotaTerminal }}</span>

                        <span class="text-[10px] font-bold text-indigo-600 uppercase tracking-widest pb-1">Kuota Saat Ini</span>
                      </div>
                    </div>

                    <button
                      @click="beliLisensiTambahan"
                      :disabled="loadingPayment"
                      class="mt-6 w-full bg-indigo-600 hover:bg-slate-900 text-white py-3 rounded-xl font-black text-xs uppercase tracking-widest transition-all shadow-lg active:scale-95 disabled:opacity-50"
                    >
                      {{ loadingPayment ? "Memproses..." : "+ Tambah Kuota (Rp 50.000)" }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="showUpgradeModal"
      class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[200] p-4 md:p-8 backdrop-blur-md overflow-y-auto print:hidden"
    >
      <div class="bg-slate-50 w-full max-w-5xl rounded-[32px] shadow-2xl relative my-auto border border-slate-200 flex flex-col max-h-[90vh]">
        <div class="p-6 md:p-8 border-b border-slate-200 flex justify-between items-center bg-white rounded-t-[32px] shrink-0">
          <div>
            <h2 class="text-2xl font-black text-slate-800 uppercase tracking-tighter">Pilih Paket Bisnis</h2>

            <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mt-1">Tingkatkan fitur aplikasi sesuai kebutuhan toko Anda</p>
          </div>

          <button
            @click="showUpgradeModal = false"
            class="w-10 h-10 bg-slate-100 hover:bg-rose-100 text-slate-500 hover:text-rose-500 rounded-2xl flex items-center justify-center transition-colors"
          >
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="p-6 md:p-8 overflow-y-auto custom-scrollbar flex-1">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div class="bg-white rounded-[24px] p-6 border border-slate-200 flex flex-col relative">
              <div class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Starter Trial</div>

              <h3 class="text-xl font-black text-slate-800 uppercase mb-3">Trial 14 Hari</h3>

              <p class="text-[11px] font-bold text-slate-500 mb-5 min-h-[35px] leading-relaxed">Validasi kesesuaian sistem dengan ekosistem bisnis Anda.</p>

              <div class="flex items-baseline gap-1 mb-6">
                <span class="text-3xl font-black text-slate-800">Rp 0</span>
              </div>

              <ul class="space-y-3 mb-8 flex-1">
                <li class="flex items-start gap-2 text-xs font-bold text-slate-500">
                  <svg class="w-4 h-4 text-emerald-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  POS Kasir Retail
                </li>

                <li class="flex items-start gap-2 text-xs font-bold text-slate-500">
                  <svg class="w-4 h-4 text-emerald-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Master Data Produk
                </li>

                <li class="flex items-start gap-2 text-xs font-bold text-slate-500">
                  <svg class="w-4 h-4 text-emerald-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Scan Barcode Reader
                </li>

                <li class="flex items-start gap-2 text-xs font-bold text-slate-500">
                  <svg class="w-4 h-4 text-emerald-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Tanpa Kartu Kredit
                </li>
              </ul>

              <button
                v-if="subPlan?.toLowerCase() === 'trial'"
                disabled
                class="w-full py-3.5 bg-slate-100 text-slate-400 rounded-xl font-black text-[10px] uppercase tracking-widest cursor-not-allowed"
              >
                Paket Saat Ini
              </button>

              <button
                v-else
                disabled
                class="w-full py-3.5 bg-slate-50 border border-slate-200 text-slate-400 rounded-xl font-black text-[10px] uppercase tracking-widest cursor-not-allowed opacity-70"
              >
                Khusus Pengguna Baru
              </button>
            </div>

            <div class="bg-white rounded-[24px] p-6 border border-slate-200 flex flex-col">
              <div class="text-[10px] font-black text-sky-500 uppercase tracking-widest mb-1">Retail Basic</div>

              <h3 class="text-xl font-black text-slate-800 uppercase mb-3">Basic</h3>

              <p class="text-[11px] font-bold text-slate-500 mb-5 min-h-[35px] leading-relaxed">Solusi solid untuk toko kelontong dengan 1 titik kasir.</p>

              <div class="flex items-baseline gap-1 mb-6">
                <span class="text-sm font-black text-slate-400">Rp</span>

                <span class="text-4xl font-black text-slate-800">
                  49

                  <span class="text-base text-slate-500">k</span>
                </span>

                <span class="text-[10px] font-bold text-slate-400">/Bulan</span>
              </div>

              <ul class="space-y-3 mb-8 flex-1">
                <li class="flex items-start gap-2 text-xs font-bold text-slate-600">
                  <svg class="w-4 h-4 text-sky-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Kasir Tanpa Batas
                </li>

                <li class="flex items-start gap-2 text-xs font-bold text-slate-600">
                  <svg class="w-4 h-4 text-sky-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Manajemen Stok Dasar
                </li>

                <li class="flex items-start gap-2 text-xs font-bold text-slate-600">
                  <svg class="w-4 h-4 text-sky-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Struk Thermal Bluetooth
                </li>

                <li class="flex items-start gap-2 text-xs font-bold text-slate-600">
                  <svg class="w-4 h-4 text-sky-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Laporan Penjualan Harian
                </li>
              </ul>

              <button
                v-if="subPlan?.toLowerCase() === 'basic'"
                disabled
                class="w-full py-3.5 bg-slate-100 text-slate-400 rounded-xl font-black text-[10px] uppercase tracking-widest cursor-not-allowed"
              >
                Paket Saat Ini
              </button>

              <button
                v-else
                @click="handleUpgrade('Basic', 49000)"
                class="w-full py-3.5 bg-slate-800 text-white rounded-xl font-black text-[10px] uppercase tracking-widest shadow-lg hover:bg-slate-700 transition-all active:scale-95"
              >
                Upgrade Basic
              </button>
            </div>

            <div class="bg-white rounded-[24px] p-6 border-2 border-indigo-500 shadow-xl relative flex flex-col transform lg:-translate-y-4">
              <div
                class="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-indigo-500 text-white text-[9px] font-black uppercase tracking-widest px-4 py-1.5 rounded-full shadow-lg whitespace-nowrap"
              >
                Paling Laris
              </div>

              <div class="text-[10px] font-black text-indigo-500 uppercase tracking-widest mb-1 mt-2">Retail Pro</div>

              <h3 class="text-xl font-black text-slate-800 uppercase mb-3">Pro</h3>

              <p class="text-[11px] font-bold text-slate-500 mb-5 min-h-[35px] leading-relaxed">Cocok untuk minimarket yang mulai mengelola karyawan.</p>

              <div class="flex items-baseline gap-1 mb-6">
                <span class="text-sm font-black text-slate-400">Rp</span>

                <span class="text-4xl font-black text-slate-800">
                  149

                  <span class="text-base text-slate-500">k</span>
                </span>

                <span class="text-[10px] font-bold text-slate-400">/Bulan</span>
              </div>

              <ul class="space-y-3 mb-8 flex-1">
                <li class="flex items-start gap-2 text-xs font-bold text-slate-700">
                  <svg class="w-4 h-4 text-indigo-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Semua Fitur Basic
                </li>

                <li class="flex items-start gap-2 text-xs font-bold text-slate-700">
                  <svg class="w-4 h-4 text-indigo-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Manajemen Akses Kasir
                </li>

                <li class="flex items-start gap-2 text-xs font-bold text-slate-700">
                  <svg class="w-4 h-4 text-indigo-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Smart Attendance & Shift
                </li>

                <li class="flex items-start gap-2 text-xs font-bold text-slate-700">
                  <svg class="w-4 h-4 text-indigo-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Laporan Ekspor (Excel/PDF)
                </li>
              </ul>

              <button
                v-if="subPlan?.toLowerCase() === 'pro'"
                disabled
                class="w-full py-3.5 bg-slate-100 text-slate-400 rounded-xl font-black text-[10px] uppercase tracking-widest cursor-not-allowed"
              >
                Paket Saat Ini
              </button>

              <button
                v-else
                @click="handleUpgrade('Pro', 149000)"
                class="w-full py-3.5 bg-indigo-600 text-white rounded-xl font-black text-[10px] uppercase tracking-widest shadow-lg shadow-indigo-500/30 hover:bg-indigo-700 transition-all active:scale-95"
              >
                Upgrade Pro
              </button>
            </div>

            <div class="bg-slate-900 rounded-[24px] p-6 border border-slate-700 flex flex-col">
              <div class="text-[10px] font-black text-amber-500 uppercase tracking-widest mb-1">Retail Premium</div>

              <h3 class="text-xl font-black text-white uppercase mb-3">Premium</h3>

              <p class="text-[11px] font-bold text-slate-400 mb-5 min-h-[35px] leading-relaxed">Kendali penuh untuk bisnis multi-cabang & gudang.</p>

              <div class="flex items-baseline gap-1 mb-6">
                <span class="text-sm font-black text-slate-400">Rp</span>

                <span class="text-4xl font-black text-white">
                  299

                  <span class="text-base text-slate-400">k</span>
                </span>

                <span class="text-[10px] font-bold text-slate-400">/Bulan</span>
              </div>

              <ul class="space-y-3 mb-8 flex-1">
                <li class="flex items-start gap-2 text-xs font-bold text-slate-300">
                  <svg class="w-4 h-4 text-amber-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Semua Fitur Pro
                </li>

                <li class="flex items-start gap-2 text-xs font-bold text-slate-300">
                  <svg class="w-4 h-4 text-amber-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Sistem Multi-Cabang (HO)
                </li>

                <li class="flex items-start gap-2 text-xs font-bold text-slate-300">
                  <svg class="w-4 h-4 text-amber-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Manajemen Multi-Gudang
                </li>

                <li class="flex items-start gap-2 text-xs font-bold text-slate-300">
                  <svg class="w-4 h-4 text-amber-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>

                  Dedicated Support 24/7
                </li>
              </ul>

              <button
                v-if="['Premium', 'premium'].includes(subPlan?.toLowerCase())"
                disabled
                class="w-full py-3.5 bg-slate-800 text-slate-400 rounded-xl font-black text-[10px] uppercase tracking-widest cursor-not-allowed"
              >
                Paket Saat Ini
              </button>

              <button
                v-else
                @click="handleUpgrade('Premium', 299000)"
                class="w-full py-3.5 bg-amber-500 text-slate-900 rounded-xl font-black text-[10px] uppercase tracking-widest shadow-lg shadow-amber-500/20 hover:bg-amber-400 transition-all active:scale-95"
              >
                Upgrade Premium
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  height: 4px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e1;

  border-radius: 10px;
}

.animate-fade-in-up {
  animation: fadeInUp 0.4s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;

    transform: translateY(10px);
  }

  to {
    opacity: 1;

    transform: translateY(0);
  }
}
</style>
