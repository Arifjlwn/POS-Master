import { ref } from "vue";
import { useRouter } from "vue-router";
import Swal from "sweetalert2";
import { posService } from "../services/posService.js";

export function useBukaKasir() {
    const router = useRouter();

    // Metadata User & Store dari LocalStorage
    const role = ref(localStorage.getItem("role"));
    const name = ref(localStorage.getItem("name") || "Operator");
    const storeName = ref(localStorage.getItem("storeName") || "POS UMKM");

    // Sesi Input States - Ambil dari localStorage jika sebelumnya kasir sengaja/tidak sengaja melakukan refresh halaman 
    const stationNumber = ref(localStorage.getItem("active_station") || "01");
    const modalAwal = ref(0);
    const loading = ref(false);

    // 🔍 Cek Sesi Aktif pas Halaman Pertama Kali Dibuka
    const checkExistingSession = async () => {
        const token = localStorage.getItem("token");
        if (!token) return;
        
        try {
            const res = await posService.checkSession(token);
            // FIX INTEGRASI: Jika server bilang sesi kasir masih aktif merayap, amankan nomor stasiunnya ke client side storage
            if (res.has_session) {
                if (res.station_number) {
                    localStorage.setItem("active_station", res.station_number);
                    stationNumber.value = res.station_number;
                }
                router.push("/retail/pos");
            }
        } catch (error) {
            console.error("Gagal cek session:", error);
        }
    };

    // Filter input agar murni angka bulat saja yang masuk (Sangat aman untuk format IDR)
    const handleInputModal = (e) => {
        const val = e.target.value.replace(/\D/g, "");
        modalAwal.value = val ? parseInt(val, 10) : 0;
    };

    // 🚀 Eksekusi Buka Sesi Kasir Baru
    const handleBukaKasir = async () => {
        if (modalAwal.value <= 0) {
            return Swal.fire({
                icon: "warning",
                title: "Modal Kosong?",
                text: "Masukkan nominal modal awal untuk uang kembalian di laci kasir.",
                confirmButtonColor: "#2563eb",
                customClass: { popup: 'rounded-[32px]' }
            });
        }

        loading.value = true;
        try {
            // EKSEKUSI API VIA SERVICE
            await posService.openSession({
                station_number: stationNumber.value,
                // FIX FINANSIAL: Gunakan nilai Integer murni , buang jauh-jauh parseFloat demi menghindari data rounding fraud!
                modal_awal: parseInt(modalAwal.value, 10),
            });

            // Gembok nomor stasiun kasir aktif di lokal browser kasir biar ga hilang pas di-refresh!
            localStorage.setItem("active_station", stationNumber.value);

            await Swal.fire({
                icon: "success",
                title: "SESSION INITIALIZED",
                text: `Kasir Station ${stationNumber.value} berhasil dibuka. Selamat bertugas!`,
                timer: 1500,
                showConfirmButton: false,
                customClass: { popup: "rounded-[32px]" },
            });

            router.push("/retail/pos");
        } catch (error) {
            // BYPASS KHUSUS UNTUK ERROR KUOTA PENUH (SATPAM TERMINAL MULTI-TENANT)
            if (
                error.response &&
                error.response.status === 403 &&
                error.response.data?.error_code === "QUOTA_FULL"
            ) {
                throw error; // Oper mentah-mentah ke BukaKasir.vue buat ngetrigger modal Midtrans !
            }

            const msg = error.response?.data?.error || "Gagal membuka kasir";

            // Interseptor Otorisasi Absensi Face AI / PIN Shift Kerja
            if (msg.toLowerCase().includes("absen")) {
                Swal.fire({
                    title: "Akses Ditolak",
                    text: "Sistem mendeteksi Anda belum melakukan absensi masuk shift menggunakan Face AI hari ini.",
                    icon: "error",
                    showCancelButton: true,
                    confirmButtonText: "Absen Sekarang 📸",
                    cancelButtonText: "Batal",
                    confirmButtonColor: "#2563eb",
                    customClass: { popup: "rounded-[32px]" },
                }).then((result) => {
                    if (result.isConfirmed) {
                        router.push("/retail/sdm/absensi");
                    }
                });
            } else {
                Swal.fire({
                    icon: "error",
                    title: "Gagal Inisialisasi",
                    text: msg,
                    confirmButtonColor: "#ef4444",
                    customClass: { popup: "rounded-[32px]" }
                });
            }
        } finally {
            loading.value = false;
        }
    };

    return {
        role,
        name,
        storeName,
        stationNumber,
        modalAwal,
        loading,
        checkExistingSession,
        handleInputModal,
        handleBukaKasir,
    };
}