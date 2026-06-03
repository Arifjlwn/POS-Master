<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import api from "../../../../api.js";
import { useBukaKasir } from "../../composables/useBukaKasir.js";

const router = useRouter();
const {
    role,
    name,
    storeName,
    stationNumber,
    modalAwal,
    loading,
    checkExistingSession,
    handleInputModal,
    handleBukaKasir,
} = useBukaKasir();

// 🚀 STATE BARU BUAT MODAL & LOADING PEMBAYARAN
const showUpgradeModal = ref(false);
const quotaTerminal = ref(1);
const loadingPayment = ref(false);

// 🚀 WRAPPER FUNGSI BUKA KASIR BIAR BISA NANGKEP ERROR DARI GOLANG
const submitBukaKasir = async () => {
    try {
        // Panggil fungsi asli dari composable lu
        await handleBukaKasir();
    } catch (error) {
        // 🛡️ TANGKEP ERROR KALO DITENDANG SATPAM (QUOTA FULL)
        if (
            error.response &&
            error.response.status === 403 &&
            error.response.data.error_code === "QUOTA_FULL"
        ) {
            showUpgradeModal.value = true;
        } else {
            // Error lain (misal belum absen, dll)
            alert(
                error.response?.data?.error ||
                    "Terjadi kesalahan saat membuka kasir",
            );
        }
    }
};

// 💸 LOGIKA MIDTRANS SNAP BUAT BELI LISENSI
const beliLisensiTambahan = async () => {
    loadingPayment.value = true;
    try {
        // 2. Ganti axios.post menjadi api.post, dan hapus '/api' di depannya karena biasanya instance api udah otomatis nambahin itu
        const response = await api.post("/retail/subscription/upgrade", {
            plan_name: "Terminal Tambahan",
            price: 50000,
        });

        const snapToken = response.data.token;

        window.snap.pay(snapToken, {
            onSuccess: function (result) {
                alert(
                    "Pembayaran Berhasil! Kuota Mesin Kasir Anda telah bertambah.",
                );
                showUpgradeModal.value = false;
                loadingPayment.value = false;
                handleBukaKasir();
            },
            onPending: function (result) {
                alert("Menunggu pembayaran Anda diselesaikan...");
                loadingPayment.value = false;
            },
            onError: function (result) {
                alert("Pembayaran Gagal!");
                loadingPayment.value = false;
            },
            onClose: function () {
                loadingPayment.value = false;
            },
        });
    } catch (error) {
        alert(
            error.response?.data?.error || "Gagal menghubungi Payment Gateway.",
        );
        loadingPayment.value = false;
    }
};

onMounted(async () => {
    checkExistingSession();

    // 🚀 TAMBAHKAN BLOK INI UNTUK MENARIK DATA QUOTA DARI BACKEND
    try {
        const res = await api.get("/retail/store/settings");
        if (res.data && res.data.data) {
            // Mengambil quota_terminal dari response database
            quotaTerminal.value = res.data.data.quota_terminal || 1;
        }
    } catch (error) {
        console.error("Gagal menarik pengaturan toko:", error);
    }

    // 🚀 INJECT SCRIPT MIDTRANS SNAP KE DALAM VUE
    const script = document.createElement("script");
    script.src = "https://app.sandbox.midtrans.com/snap/snap.js"; // Pakai sandbox buat testing

    // ⚠️ WAJIB GANTI PAKE CLIENT KEY MIDTRANS LU SENDIRI DI SINI:
    script.setAttribute("data-client-key", "SB-Mid-client-XXXXXXXXXXXXXXXX");
    document.head.appendChild(script);
});
</script>

<template>
    <div
        class="min-h-screen bg-slate-950 flex items-center justify-center p-6 relative overflow-hidden font-sans selection:bg-blue-100"
    >
        <div
            class="absolute top-0 left-0 w-full h-full opacity-10 pointer-events-none"
        >
            <div
                class="absolute -top-24 -left-24 w-96 h-96 bg-blue-600 rounded-full blur-[120px]"
            ></div>
            <div
                class="absolute -bottom-24 -right-24 w-96 h-96 bg-indigo-600 rounded-full blur-[120px]"
            ></div>
        </div>

        <div class="w-full max-w-lg relative">
            <div
                class="bg-white rounded-[48px] p-8 md:p-12 shadow-2xl border-[12px] border-slate-900/5 relative overflow-hidden"
            >
                <div class="text-center mb-10">
                    <div
                        class="w-20 h-20 bg-slate-900 rounded-[28px] flex items-center justify-center mx-auto mb-6 shadow-xl shadow-blue-500/20 transform -rotate-3 hover:rotate-0 transition-all duration-500"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-10 h-10 text-blue-500"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2.5"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <path
                                d="M20 20H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2z"
                            />
                            <path d="M2 12h20" />
                            <path d="M7 16h.01" />
                            <path d="M11 16h.01" />
                            <path d="M15 16h.01" />
                            <path d="M7 8h.01" />
                            <path d="M11 8h.01" />
                            <path d="M15 8h.01" />
                        </svg>
                    </div>
                    <h1
                        class="text-3xl font-black text-slate-900 tracking-tighter uppercase"
                    >
                        Point of <span class="text-blue-600">Sale</span>
                    </h1>
                    <p
                        class="text-slate-400 font-bold text-[10px] uppercase tracking-[0.4em] mt-1"
                    >
                        {{ storeName }}
                    </p>
                </div>

                <div class="space-y-8">
                    <div
                        class="bg-slate-50 p-5 rounded-[32px] border border-slate-100 flex items-center justify-between"
                    >
                        <div class="flex items-center gap-4">
                            <div
                                class="w-12 h-12 rounded-2xl bg-white border border-slate-200 flex items-center justify-center text-lg shadow-sm"
                            >
                                👤
                            </div>
                            <div>
                                <label
                                    class="text-[9px] font-black text-slate-400 uppercase tracking-widest block"
                                    >Logged Operator</label
                                >
                                <div
                                    class="text-sm font-black text-slate-800 uppercase flex items-center gap-2"
                                >
                                    <span
                                        class="w-2 h-2 bg-green-500 rounded-full animate-pulse"
                                    ></span>
                                    {{ name }}
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="space-y-3">
                        <label
                            class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-2 italic"
                            >Select Device Station</label
                        >
                        <div class="grid grid-cols-3 gap-3">
                            <button
                                v-for="n in ['01', '02', '03']"
                                :key="n"
                                @click="stationNumber = n"
                                :class="
                                    stationNumber === n
                                        ? 'bg-slate-900 text-white shadow-xl shadow-slate-200 scale-105 border-slate-900'
                                        : 'bg-slate-50 text-slate-400 grayscale border-transparent'
                                "
                                class="flex flex-col items-center gap-2 py-4 rounded-[28px] border-2 transition-all duration-300 outline-none"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="w-5 h-5"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2.5"
                                >
                                    <rect
                                        width="18"
                                        height="12"
                                        x="3"
                                        y="4"
                                        rx="2"
                                        ry="2"
                                    />
                                    <line x1="2" x2="22" y1="20" y2="20" />
                                    <line x1="12" x2="12" y1="16" y2="20" />
                                </svg>
                                <span
                                    class="text-[10px] font-black tracking-tighter"
                                    >POS {{ n }}</span
                                >
                            </button>
                        </div>
                    </div>

                    <div class="space-y-3">
                        <label
                            class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-2"
                            >Floating Capital (Modal Awal)</label
                        >
                        <div class="relative group">
                            <div
                                class="absolute left-6 top-1/2 -translate-y-1/2 flex flex-col items-center select-none pointer-events-none"
                            >
                                <span
                                    class="text-[10px] font-black text-blue-400 uppercase leading-none"
                                    >IDR</span
                                >
                                <span class="text-xl font-black text-blue-600"
                                    >Rp</span
                                >
                            </div>
                            <input
                                type="text"
                                placeholder="0"
                                :value="
                                    modalAwal === 0
                                        ? ''
                                        : modalAwal.toLocaleString('id-ID')
                                "
                                @input="handleInputModal"
                                class="w-full bg-blue-50/30 border-2 border-blue-100 p-8 pl-24 rounded-[36px] font-black text-4xl text-slate-900 focus:border-blue-600 focus:bg-white outline-none transition-all placeholder:text-slate-200"
                            />
                        </div>
                    </div>

                    <div class="pt-4">
                        <button
                            @click="submitBukaKasir"
                            :disabled="loading || loadingPayment"
                            class="w-full bg-blue-600 hover:bg-slate-900 text-white p-6 rounded-[32px] font-black text-sm uppercase tracking-[0.2em] shadow-2xl shadow-blue-200 transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-4"
                        >
                            <span>{{
                                loading || loadingPayment
                                    ? "Accessing Server..."
                                    : "Initialize Session"
                            }}</span>
                            <svg
                                v-if="!loading && !loadingPayment"
                                xmlns="http://www.w3.org/2000/svg"
                                class="w-5 h-5"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                                stroke-width="3"
                            >
                                <path d="M5 12h14" />
                                <path d="m12 5 7 7-7 7" />
                            </svg>
                        </button>

                        <div class="mt-8 text-center" v-if="role === 'owner'">
                            <button
                                @click="router.push('/retail/dashboard')"
                                class="text-[10px] font-black text-slate-300 hover:text-blue-600 uppercase tracking-[0.3em] transition-colors flex items-center justify-center gap-2 mx-auto outline-none"
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="w-3 h-3"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="3"
                                >
                                    <path d="m15 18-6-6 6-6" />
                                </svg>
                                Back to Dashboard
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- MODAL UPGRADE / JUALAN LISENSI KASIR -->
        <div
            v-if="showUpgradeModal"
            class="fixed inset-0 bg-slate-950/80 backdrop-blur-sm flex items-center justify-center z-[100] px-4 transition-all"
        >
            <div
                class="bg-white p-8 rounded-[40px] shadow-2xl max-w-sm w-full text-center border-8 border-slate-100 relative overflow-hidden transform scale-100"
            >
                <div
                    class="absolute top-0 left-0 w-full h-2 bg-gradient-to-r from-red-500 to-rose-500"
                ></div>

                <div
                    class="w-20 h-20 bg-rose-50 text-rose-500 rounded-full flex items-center justify-center mx-auto mb-6"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="w-10 h-10"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        stroke-width="2.5"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                        />
                    </svg>
                </div>

                <h3 class="text-2xl font-black text-slate-900 tracking-tight">
                    Kuota Penuh!
                </h3>
                <p
                    class="text-xs font-bold text-slate-500 mt-3 leading-relaxed"
                >
                    Lisensi toko Anda saat ini hanya mengizinkan
                    {{ quotaTerminal }} Kasir beroperasi. Karyawan lain sedang
                    menggunakan laci kasir di perangkat lain.
                </p>

                <!-- LOGIKA UNTUK OWNER (MUNCUL TOMBOL BAYAR) -->
                <div
                    v-if="role === 'owner'"
                    class="mt-8 bg-indigo-50/50 border-2 border-indigo-100 p-5 rounded-[24px]"
                >
                    <h4
                        class="font-black text-indigo-900 text-sm uppercase tracking-wider"
                    >
                        Upgrade Sistem Kasir
                    </h4>
                    <p class="text-[11px] font-bold text-indigo-600 mt-2">
                        Buka 2 Terminal Kasir sekaligus untuk mempercepat
                        antrean pelanggan Anda.
                    </p>

                    <button
                        @click="beliLisensiTambahan"
                        :disabled="loadingPayment"
                        class="mt-5 w-full bg-indigo-600 hover:bg-indigo-700 text-white py-4 rounded-2xl font-black text-xs uppercase tracking-widest transition-all shadow-lg shadow-indigo-200 active:scale-95 disabled:opacity-50"
                    >
                        {{
                            loadingPayment
                                ? "Memproses..."
                                : "Beli Kuota (Rp 50.000)"
                        }}
                    </button>
                </div>

                <!-- LOGIKA UNTUK KASIR (MUNCUL INFO HUBUNGI OWNER) -->
                <div
                    v-else
                    class="mt-8 bg-amber-50 border-2 border-amber-100 p-5 rounded-[24px]"
                >
                    <h4
                        class="font-black text-amber-900 text-sm uppercase tracking-wider"
                    >
                        Akses Terbatas
                    </h4>
                    <p class="text-[11px] font-bold text-amber-700 mt-2">
                        Silakan hubungi <b>Owner / Manajer</b> toko untuk
                        melakukan penambahan lisensi terminal Kasir.
                    </p>
                </div>

                <button
                    @click="showUpgradeModal = false"
                    class="mt-6 text-[10px] font-black text-slate-400 hover:text-slate-600 uppercase tracking-widest transition-colors outline-none"
                >
                    Tutup & Kembali
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.transition-all {
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}
input[type="number"]::-webkit-inner-spin-button,
input[type="number"]::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
}
</style>
