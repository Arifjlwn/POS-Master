import Swal from 'sweetalert2';
import { onMounted, onUnmounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../../../api.js';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const getImageUrl = (path) => {
    // Pengetatan filter validasi karakter string bawaan localstorage
    if (!path || path === 'null' || path === 'undefined' || path === '') return '';
    if (path.startsWith('http://') || path.startsWith('https://')) return path;
    
    const cleanPath = path.startsWith('/') ? path : `/${path}`;
    return `${API_BASE_URL}${cleanPath}`;
};

export function useSidebar() {
    const router = useRouter();
    const route = useRoute();
    const sidebarOpen = ref(false);

    const openGroups = ref({
        stock: route.path.includes('master-produk') || route.path.includes('penerimaan-barang') || route.path.includes('stock-opname') || route.path.includes('retur-barang'),
        admin: route.path.startsWith('/retail/dashboard') || route.path.startsWith('/retail/karyawan') || route.path.includes('report'),
    });

    const toggleGroup = (group) => {
        openGroups.value[group] = !openGroups.value[group];
    };

    // Handler fungsi bernama agar bisa di-remove tanpa menyebabkan Memory Leak
    const syncUserData = () => {
        user.value.name = localStorage.getItem('name') || 'User';
        user.value.role = localStorage.getItem('role') || 'staff';
        user.value.foto_url = getImageUrl(localStorage.getItem('foto_url'));
    };

    const syncStoreData = () => {
        user.value.storeName = localStorage.getItem('storeName') || 'POS UMKM';
        user.value.storeLogo = getImageUrl(localStorage.getItem('storeLogo'));
    };

    // REVISI UTAMA: Bentuk objek state user dengan fungsi pembersih URL langsung sejak awal inisialisasi
    const user = ref({
        name: localStorage.getItem('name') || 'User',
        role: localStorage.getItem('role') || 'staff',
        storeName: localStorage.getItem('storeName') || 'POS UMKM',
        storeLogo: getImageUrl(localStorage.getItem('storeLogo')),
        foto_url: getImageUrl(localStorage.getItem('foto_url')),
    });

    onMounted(async () => {
        // Jalankan sinkronisasi instan saat komponen diakses
        syncUserData();
        syncStoreData();

        // REVISI KEAMANAN DATA SINKRON: Tarik data profil paling valid (/me) langsung dari server database 
        // untuk memastikan foto profil staff tetap muncul terupdate meskipun localstorage sempat di-clear saat logout.
        try {
            const resMe = await api.get('/me');
            const freshProfile = resMe.data;

            if (freshProfile.name) localStorage.setItem('name', freshProfile.name);
            if (freshProfile.role) localStorage.setItem('role', freshProfile.role);
            if (freshProfile.foto_url) {
                localStorage.setItem('foto_url', freshProfile.foto_url);
            } else {
                localStorage.removeItem('foto_url');
            }
            
            // Perbarui state secara reaktif di internal UI
            syncUserData();
        } catch (err) {
            console.error("Gagal sinkronisasi data profil dari server:", err.message);
        }

        window.addEventListener('storage', syncUserData);
        window.addEventListener('profile-updated', syncUserData);
        window.addEventListener('store-updated', syncStoreData);
    });

    onUnmounted(() => {
        window.removeEventListener('storage', syncUserData);
        window.removeEventListener('profile-updated', syncUserData);
        window.removeEventListener('store-updated', syncStoreData);
    });

    const logout = () => {
        const isOwner = user.value.role === 'owner';

        Swal.fire({
            title: 'Sesi Operasional',
            text: isOwner ? 'Pilih tindakan untuk mengakhiri sesi Anda saat ini.' : 'Apakah Anda yakin ingin keluar dari sistem?',
            icon: 'question',
            showCancelButton: true,
            showDenyButton: isOwner,
            confirmButtonText: isOwner ? 'Ganti Cabang' : 'Keluar Total',
            denyButtonText: isOwner ? 'Keluar Total' : '',
            cancelButtonText: 'Batal',
            buttonsStyling: false,
            customClass: {
                popup: 'rounded-[32px] p-6 border border-slate-100 shadow-2xl',
                title: 'text-2xl font-black text-slate-800 tracking-tight',
                htmlContainer: 'text-sm font-bold text-slate-500 mt-2 mb-6',
                actions: 'flex flex-col sm:flex-row gap-3 w-full',
                confirmButton: isOwner 
                    ? 'w-full sm:flex-1 bg-indigo-600 hover:bg-indigo-700 text-white font-black text-[10px] md:text-xs uppercase tracking-widest py-4 rounded-2xl transition-all shadow-xl shadow-indigo-200 active:scale-95 order-1 sm:order-1'
                    : 'w-full sm:flex-1 bg-rose-500 hover:bg-rose-600 text-white font-black text-[10px] md:text-xs uppercase tracking-widest py-4 rounded-2xl transition-all shadow-xl shadow-rose-200 active:scale-95 order-1 sm:order-1',
                denyButton: isOwner 
                    ? 'w-full sm:flex-1 bg-rose-500 hover:bg-rose-600 text-white font-black text-[10px] md:text-xs uppercase tracking-widest py-4 rounded-2xl transition-all shadow-xl shadow-rose-200 active:scale-95 order-2 sm:order-2' 
                    : '',
                cancelButton: 'w-full sm:flex-1 bg-slate-100 hover:bg-slate-200 text-slate-500 hover:text-slate-700 font-black text-[10px] md:text-xs uppercase tracking-widest py-4 rounded-2xl transition-all active:scale-95 order-3 sm:order-3',
            },
        }).then(async (result) => {
            if (isOwner && result.isConfirmed) {
                try {
                    const resSetting = await api.get('/retail/store/settings');
                    const freshData = resSetting.data.data;
                    let tempStores = localStorage.getItem('temp_stores');

                    if (tempStores) {
                        let storesArr = JSON.parse(tempStores);
                        storesArr = storesArr.map((store) => {
                            if (store.id === freshData.id) {
                                return { ...store, nama_toko: freshData.nama_toko, logo_url: freshData.logo_url, kota: freshData.kota, subscription_plan: freshData.subscription_plan };
                            }
                            return store;
                        });
                        localStorage.setItem('temp_stores', JSON.stringify(storesArr));
                    }
                } catch (e) {
                    console.error('Sinkronisasi list cabang ditunda:', e.message);
                }

                localStorage.removeItem('store_id');
                localStorage.removeItem('storeName');
                localStorage.removeItem('storeLogo');
                localStorage.removeItem('subscriptionPlan');

                router.push('/select-store');

            } else if ((isOwner && result.isDenied) || (!isOwner && result.isConfirmed)) {
                localStorage.clear();
                window.location.href = '/login'; 
            }
        });
    };

    return { route, sidebarOpen, openGroups, user, toggleGroup, logout };
}