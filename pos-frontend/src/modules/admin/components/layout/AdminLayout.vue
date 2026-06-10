<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import Sidebar from './Sidebar.vue';
import TopNavbar from './TopNavbar.vue';

// 🚀 KONFIGURASI AWAL: Sidebar diatur tertutup (false) secara default untuk semua resolusi layar
const isSidebarOpen = ref(false);

// Fungsi memantau perubahan ukuran layar guna menjaga konsistensi UI
const handleResize = () => {
    if (window.innerWidth < 1024) {
        isSidebarOpen.value = false;
    }
};

onMounted(() => window.addEventListener('resize', handleResize));
onUnmounted(() => window.removeEventListener('resize', handleResize));

// Fungsi menutup sidebar secara aman pada perangkat mobile
const handleCloseSidebar = () => {
    if (window.innerWidth < 1024) {
        isSidebarOpen.value = false;
    }
};

// Fungsi mengubah status buka/tutup sidebar dari TopNavbar
const toggleSidebar = () => {
    isSidebarOpen.value = !isSidebarOpen.value;
};
</script>

<template>
    <div class="min-h-screen bg-[#0B0F19] text-slate-300 font-sans flex overflow-hidden selection:bg-indigo-500/30">
        
        <Sidebar :isOpen="isSidebarOpen" @close="handleCloseSidebar" />

        <div class="flex-1 flex flex-col h-screen overflow-hidden relative">
            
            <TopNavbar @toggle-sidebar="toggleSidebar" />

            <main class="flex-1 overflow-x-hidden overflow-y-auto bg-[#0B0F19] p-4 lg:p-8">
                <router-view v-slot="{ Component }">
                    <transition name="fade" mode="out-in">
                        <component :is="Component" />
                    </transition>
                </router-view>
            </main>

        </div>
    </div>
</template>

<style>
/* Animasi transisi perpindahan halaman */
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.2s ease, transform 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
    opacity: 0;
    transform: translateY(10px);
}
</style>