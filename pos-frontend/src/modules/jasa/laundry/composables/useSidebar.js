import { ref, reactive, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Swal from 'sweetalert2';

// State global di level file biar konsisten pas di-import lintas view bray
const sidebarOpen = ref(false);
const openGroups = reactive({
    stock: true,
    admin: true
});

const user = reactive({
    name: localStorage.getItem('name') || 'Owner',
    role: localStorage.getItem('role') || 'owner',
    storeName: localStorage.getItem('storeName') || localStorage.getItem('store_name') || 'ARZU LAUNDRY',
    storeLogo: localStorage.getItem('storeLogo') || '',
    foto_url: localStorage.getItem('foto_url') || ''
});

export function useSidebar() {
    const route = useRoute();
    const router = useRouter();

    const toggleGroup = (group) => {
        openGroups[group] = !openGroups[group];
    };

    const logout = () => {
        sidebarOpen.value = false;
        Swal.fire({
            title: 'Keluar Kasir?',
            text: 'Sesi operasional mesin kasir Anda akan diakhiri bray.',
            icon: 'question',
            showCancelButton: true,
            confirmButtonColor: '#e11d48',
            cancelButtonColor: '#94a3b8',
            confirmButtonText: 'Ya, Keluar'
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