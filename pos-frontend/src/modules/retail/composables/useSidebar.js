import Swal from 'sweetalert2';
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

// 🚀 FUNGSI PINTER ANTI-GAMBAR-RUSAK
const getImageUrl = (path) => {
    if (!path) return '';
    if (path.startsWith('http')) return path;
    return `http://localhost:8080${path}`;
};

export function useSidebar() {
    const router = useRouter();
    const route = useRoute();
    const sidebarOpen = ref(false);

    // State untuk kontrol buka-tutup grup menu accordion
    const openGroups = ref({
        stock:
            route.path.includes('master-produk') ||
            route.path.includes('penerimaan-barang') ||
            route.path.includes('stock-opname') ||
            route.path.includes('retur-barang'),
        admin: route.path.startsWith('/retail/dashboard') || route.path.startsWith('/retail/karyawan') || route.path.includes('report')
    });

    const toggleGroup = (group) => {
        openGroups.value[group] = !openGroups.value[group];
    };

    // 🚀 BUNGKUS PAKAI getImageUrl DI SINI
    const user = ref({
        name: localStorage.getItem('name') || 'User',
        role: localStorage.getItem('role') || 'staff',
        storeName: localStorage.getItem('storeName') || 'POS UMKM',
        storeLogo: getImageUrl(localStorage.getItem('storeLogo')),
        foto_url: getImageUrl(localStorage.getItem('foto_url'))
    });

    onMounted(() => {
        // Tarik data saat halaman pertama dimuat
        user.value.name = localStorage.getItem('name') || 'User';
        user.value.role = localStorage.getItem('role') || 'staff';
        user.value.storeName = localStorage.getItem('storeName') || 'POS UMKM';

        // 🚀 BUNGKUS PAKAI getImageUrl JUGA DI SINI (Buat nge-refresh ulang pas pindah rute)
        user.value.storeLogo = getImageUrl(localStorage.getItem('storeLogo'));
        user.value.foto_url = getImageUrl(localStorage.getItem('foto_url'));

        // 🚀 EVENT LISTENER SAKTI
        // Biar kalau lu update data di halaman Pengaturan/Akun, Sidebar auto ganti tanpa refresh!
        window.addEventListener('storage', () => {
            user.value.name = localStorage.getItem('name') || 'User';
            user.value.storeLogo = getImageUrl(localStorage.getItem('storeLogo'));
            user.value.foto_url = getImageUrl(localStorage.getItem('foto_url'));
        });
        window.addEventListener('profile-updated', () => {
            user.value.name = localStorage.getItem('name') || 'User';
            user.value.foto_url = getImageUrl(localStorage.getItem('foto_url'));
        });
        window.addEventListener('store-updated', () => {
            user.value.storeName = localStorage.getItem('storeName') || 'POS UMKM';
            user.value.storeLogo = getImageUrl(localStorage.getItem('storeLogo'));
        });
    });

    const logout = () => {
        Swal.fire({
            title: 'Sesi Operasional',
            text: 'Pilih tindakan untuk mengakhiri sesi Anda saat ini.',
            icon: 'question',
            showCancelButton: true,
            showDenyButton: true,
            confirmButtonText: 'Ganti Cabang',
            denyButtonText: 'Keluar Total',
            cancelButtonText: 'Batal',

            // 🚀 INI KUNCINYA: Matikan CSS bawaan biar Tailwind bisa ambil alih 100%
            buttonsStyling: false,

            customClass: {
                popup: 'rounded-[32px] p-6 border border-slate-100 shadow-2xl',
                title: 'text-2xl font-black text-slate-800 tracking-tight',
                htmlContainer: 'text-sm font-bold text-slate-500 mt-2 mb-6',

                // Bikin tombolnya sejajar di desktop, numpuk ke bawah di mobile
                actions: 'flex flex-col sm:flex-row gap-3 w-full',

                // TOMBOL 1: GANTI CABANG (Indigo)
                confirmButton:
                    'w-full sm:flex-1 bg-indigo-600 hover:bg-indigo-700 text-white font-black text-[10px] md:text-xs uppercase tracking-widest py-4 rounded-2xl transition-all shadow-xl shadow-indigo-200 active:scale-95 order-1 sm:order-1',

                // TOMBOL 2: KELUAR TOTAL (Merah / Rose)
                denyButton:
                    'w-full sm:flex-1 bg-rose-500 hover:bg-rose-600 text-white font-black text-[10px] md:text-xs uppercase tracking-widest py-4 rounded-2xl transition-all shadow-xl shadow-rose-200 active:scale-95 order-2 sm:order-2',

                // TOMBOL 3: BATAL (Abu-abu)
                cancelButton:
                    'w-full sm:flex-1 bg-slate-100 hover:bg-slate-200 text-slate-500 hover:text-slate-700 font-black text-[10px] md:text-xs uppercase tracking-widest py-4 rounded-2xl transition-all active:scale-95 order-3 sm:order-3'
            }
        }).then((result) => {
            if (result.isConfirmed) {
                // 🚀 SKENARIO GANTI CABANG (JANGAN HAPUS TOKEN!)
                // Biarin token tetep ada biar Vue Router ngijinin lu masuk ke Select Store
                router.push('/select-store');
            } else if (result.isDenied) {
                // 🚀 SKENARIO KELUAR TOTAL (Bersihin semua memory)
                localStorage.clear();
                router.push('/login');
            }
        });
    };

    return {
        route,
        sidebarOpen,
        openGroups,
        user,
        toggleGroup,
        logout
    };
}
