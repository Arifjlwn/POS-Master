import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { sessionService } from '../services/posService.js';
import Swal from 'sweetalert2';

export function useBukaKasir() {
    const router = useRouter();
    
    // Metadata User & Store dari LocalStorage
    const role = ref(localStorage.getItem('role'));
    const name = ref(localStorage.getItem('name') || 'Operator');
    const storeName = ref(localStorage.getItem('storeName') || 'POS UMKM');
    
    // Sesi Input States
    const stationNumber = ref('01');
    const modalAwal = ref(0);
    const loading = ref(false);

    // 🔍 Cek Sesi Aktif pas Halaman Pertama Kali Dibuka
    const checkExistingSession = async () => {
        try {
            const res = await sessionService.checkActiveSession();
            if (res.has_session) {
                router.push('/retail/pos');
            }
        } catch (error) {
            console.error("Gagal cek session:", error);
        }
    };

    // Filter input agar murni angka saja yang masuk
    const handleInputModal = (e) => {
        const val = e.target.value.replace(/\D/g, '');
        modalAwal.value = val ? parseInt(val, 10) : 0;
    };

    // 🚀 Eksekusi Buka Sesi Kasir Baru
    const handleBukaKasir = async () => {
        if (modalAwal.value <= 0) {
            return Swal.fire({
                icon: 'warning',
                title: 'Modal Kosong?',
                text: 'Masukkan nominal modal awal untuk uang kembalian di laci kasir.',
                confirmButtonColor: '#2563eb'
            });
        }

        loading.value = true;
        try {
            await sessionService.openSession({
                station_number: stationNumber.value,
                modal_awal: parseFloat(modalAwal.value)
            });

            await Swal.fire({
                icon: 'success',
                title: 'SESSION OPENED',
                text: `Kasir Station ${stationNumber.value} berhasil dibuka. Selamat bertugas!`,
                timer: 1500,
                showConfirmButton: false,
                customClass: { popup: 'rounded-[32px]' }
            });

            router.push('/retail/pos');

        } catch (error) {
            const msg = error.response?.data?.error || 'Gagal membuka kasir';
            
            // 📸 Interseptor Otorisasi Absensi Face AI / PIN Kak Arif
            if (msg.toLowerCase().includes('absen')) {
                Swal.fire({
                    title: 'Otorisasi Gagal',
                    text: 'Sistem mendeteksi Anda belum melakukan absensi masuk shift hari ini.',
                    icon: 'error',
                    showCancelButton: true,
                    confirmButtonText: 'Absen Sekarang 📸',
                    confirmButtonColor: '#2563eb',
                    customClass: { popup: 'rounded-[32px]' }
                }).then((result) => {
                    if (result.isConfirmed) {
                        router.push('/retail/absensi');
                    }
                });
            } else {
                Swal.fire('Error', msg, 'error');
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
        handleBukaKasir
    };
}