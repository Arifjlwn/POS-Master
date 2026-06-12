<script setup>
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Swal from 'sweetalert2';

const props = defineProps({
    isOpen: Boolean
});

const emit = defineEmits(['close']);
const route = useRoute();
const router = useRouter();

// Konfigurasi data navigasi Mission Control
const menuItems = [
    { name: 'Mission Control', path: '/admin/dashboard', icon: '<path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" />' },
    { name: 'Analytics', path: '/admin/analytics', icon: '<path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />' },
    { name: 'Tenant & Billing Hub', path: '/admin/tenant-hub', icon: '<path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />' },
    { name: 'Users', path: '/admin/users', icon: '<path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />' },
    { name: 'Monitoring', path: '/admin/monitoring', icon: '<path stroke-linecap="round" stroke-linejoin="round" d="M9 17.25v1.007a3 3 0 01-.879 2.122L7.5 21h9l-.621-.621A3 3 0 0115 18.257V17.25m6-12V15a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 15V5.25m18 0A2.25 2.25 0 0018.75 3H5.25A2.25 2.25 0 003 5.25m18 0V12a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 12V5.25" />' },
    { name: 'Audit Logs', path: '/admin/audit', icon: '<path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />' },
    { name: 'System Settings', path: '/admin/settings', icon: '<path stroke-linecap="round" stroke-linejoin="round" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.324.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.24-.438.613-.431.992a6.759 6.759 0 010 .255c-.007.378.138.75.43.99l1.005.828c.424.35.534.954.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.57 6.57 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.28c-.09.543-.56.941-1.11.941h-2.594c-.55 0-1.02-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.992a6.932 6.932 0 010-.255c.007-.378-.138-.75-.43-.99l-1.004-.828a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.087.22-.128.332-.183.582-.495.644-.869l.214-1.281z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />' }
];

const checkActive = (path) => route.path.startsWith(path);

// 🚀 FUNGSI LOGOUT TERPUSAT: Membersihkan session dan mengarahkan kembali ke gerbang login admin
const handleLogout = () => {
    Swal.fire({
        title: 'Keluar Sistem?',
        text: 'Sesi otentikasi Mission Control Anda akan diakhiri.',
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'Ya, Keluar',
        cancelButtonText: 'Batalkan',
        confirmButtonColor: '#ef4444',
        cancelButtonColor: '#334155',
        customClass: { popup: 'rounded-[24px]' }
    }).then((result) => {
        if (result.isConfirmed) {
            // Bersihkan data otentikasi admin dari penyimpanan lokal
            localStorage.removeItem('token');
            localStorage.removeItem('role');
            localStorage.removeItem('name');
            
            // Redirect murni menuju gerbang login admin pusat
            router.push('/admin/login');
        }
    });
};
</script>

<template>
    <div v-if="isOpen" @click="emit('close')" class="fixed inset-0 bg-black/60 z-40 lg:hidden backdrop-blur-sm transition-opacity"></div>

    <aside 
        :class="isOpen ? 'translate-x-0 ml-0' : '-translate-x-full lg:translate-x-0 lg:-ml-72'"
        class="fixed lg:static inset-y-0 left-0 z-50 w-72 bg-[#0B0F19] border-r border-slate-800 transform transition-all duration-300 ease-in-out flex flex-col shrink-0"
    >
        <div class="h-20 flex items-center px-8 border-b border-slate-800">
            <div class="flex items-center gap-3">
                <div class="w-8 h-8 bg-indigo-600 rounded-lg flex items-center justify-center shadow-lg shadow-indigo-600/20">
                    <svg class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                    </svg>
                </div>
                <div>
                    <h1 class="text-white font-black tracking-widest text-sm uppercase whitespace-nowrap">Arzura POS</h1>
                    <p class="text-[10px] text-slate-500 font-bold uppercase tracking-widest whitespace-nowrap">Root Admin</p>
                </div>
            </div>
        </div>

        <div class="flex-1 overflow-y-auto py-6 px-4 space-y-1 scrollbar-hide">
            <router-link 
                v-for="item in menuItems" 
                :key="item.name" 
                :to="item.path"
                @click="emit('close')"
                class="flex items-center gap-3 px-4 py-3.5 rounded-xl font-bold text-sm transition-all duration-200 group relative whitespace-nowrap"
                :class="checkActive(item.path) ? 'bg-indigo-600/10 text-indigo-400' : 'text-slate-400 hover:bg-slate-800/50 hover:text-white'"
            >
                <div v-if="checkActive(item.path)" class="absolute left-0 top-1/2 -translate-y-1/2 w-1.5 h-8 bg-indigo-500 rounded-r-full"></div>
                
                <svg v-html="item.icon" class="w-5 h-5 shrink-0" :class="checkActive(item.path) ? 'text-indigo-400' : 'text-slate-500 group-hover:text-slate-300'" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor"></svg>
                <span>{{ item.name }}</span>
            </router-link>
        </div>

        <div class="p-4 border-t border-slate-800">
            <button @click="handleLogout" class="w-full flex items-center justify-center gap-2 py-3 px-4 rounded-xl text-xs font-black text-red-400 bg-red-500/10 hover:bg-red-500/20 transition-all uppercase tracking-widest whitespace-nowrap">
                <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
                </svg>
                <span>Keluar Sistem</span>
            </button>
        </div>
    </aside>
</template>