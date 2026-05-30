import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import Swal from 'sweetalert2';

export function useSidebar() {
    const router = useRouter();
    const route = useRoute();
    const sidebarOpen = ref(false);

    // State untuk kontrol buka-tutup grup menu accordion
    const openGroups = ref({
        stock: route.path.includes('master-produk') || route.path.includes('penerimaan-barang') || route.path.includes('stock-opname') || route.path.includes('retur-barang'),
        admin: route.path.startsWith('/retail/dashboard') || route.path.startsWith('/retail/karyawan') || route.path.includes('report')
    });

    const toggleGroup = (group) => {
        openGroups.value[group] = !openGroups.value[group];
    };

    const user = ref({
        name: localStorage.getItem('name') || 'User',
        role: localStorage.getItem('role') || 'staff',
        storeName: localStorage.getItem('storeName') || 'POS UMKM',
        storeLogo: localStorage.getItem('storeLogo') || '',
        foto_url: localStorage.getItem('foto_url') || '', // 🚀 TAMBAHIN INI
    });

    onMounted(() => {
        // Tarik data saat halaman pertama dimuat
        user.value.name = localStorage.getItem('name') || 'User';
        user.value.role = localStorage.getItem('role') || 'staff';
        user.value.storeName = localStorage.getItem('storeName') || 'POS UMKM';
        user.value.storeLogo = localStorage.getItem('storeLogo') || ''; // 🚀 TAMBAHIN INI
        user.value.foto_url = localStorage.getItem('foto_url') || '';   // 🚀 TAMBAHIN INI

        // 🚀 EVENT LISTENER SAKTI
        // Biar kalau lu update data di halaman Pengaturan/Akun, Sidebar auto ganti tanpa refresh!
        window.addEventListener('storage', () => {
            user.value.name = localStorage.getItem('name') || 'User';
            user.value.storeLogo = localStorage.getItem('storeLogo') || '';
            user.value.foto_url = localStorage.getItem('foto_url') || '';
        });
    });

    const logout = () => {
        Swal.fire({
            title: 'Mau keluar, Bos?',
            text: "Pastikan semua kerjaan hari ini sudah tersimpan ya!",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#4f46e5',
            cancelButtonColor: '#94a3b8',
            confirmButtonText: 'Ya, Logout',
            cancelButtonText: 'Batal',
            customClass: {
                popup: 'rounded-[32px]',
                confirmButton: 'rounded-[16px] font-black px-6 py-3',
                cancelButton: 'rounded-[16px] font-black px-6 py-3'
            }
        }).then((result) => {
            if (result.isConfirmed) {
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