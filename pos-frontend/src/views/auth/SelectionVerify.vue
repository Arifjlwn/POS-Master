<script setup>
import { useRoute, useRouter } from 'vue-router';
import Swal from 'sweetalert2';

const route = useRoute();
const router = useRouter();
const email = route.query.email;
const phone = route.query.phone;

const selectMethod = (method) => {
    if (method === 'whatsapp') {
        // Nanti Mas bisa integrasi API WA Blast di sini
        Swal.fire('Info', 'Fitur WhatsApp sedang dalam pengembangan, gunakan Email dulu ya!', 'info');
        return;
    }
    
    // Jika pilih email, langsung lanjut ke halaman input OTP
    router.push({ path: '/verify', query: { email: email } });
};
</script>

<template>
    <div class="min-h-screen flex items-center justify-center bg-[#F8FAFC] p-6">
        <div class="w-full max-w-md bg-white rounded-[40px] p-10 shadow-2xl border border-white text-center">
            <h2 class="text-3xl font-black text-slate-900 tracking-tighter">VERIFIKASI AKUN</h2>
            <p class="text-slate-400 font-bold text-xs uppercase tracking-widest mt-2 mb-10">Pilih metode pengiriman kode OTP</p>

            <div class="space-y-4">
                <button @click="selectMethod('email')" class="w-full p-6 border-2 border-slate-50 bg-slate-50 rounded-3xl flex items-center gap-4 hover:border-blue-600 hover:bg-blue-50 transition-all group">
                    <div class="w-12 h-12 bg-white rounded-2xl flex items-center justify-center text-blue-600 shadow-sm group-hover:bg-blue-600 group-hover:text-white transition-all">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>
                    </div>
                    <div class="text-left">
                        <div class="text-sm font-black text-slate-800 uppercase">Kirim via Email</div>
                        <div class="text-[10px] text-slate-400 font-bold">{{ email }}</div>
                    </div>
                </button>

                <button @click="selectMethod('whatsapp')" class="w-full p-6 border-2 border-slate-50 bg-slate-50 rounded-3xl flex items-center gap-4 hover:border-green-600 hover:bg-green-50 transition-all group">
                    <div class="w-12 h-12 bg-white rounded-2xl flex items-center justify-center text-green-600 shadow-sm group-hover:bg-green-600 group-hover:text-white transition-all">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
                    </div>
                    <div class="text-left">
                        <div class="text-sm font-black text-slate-800 uppercase">Kirim via WhatsApp</div>
                        <div class="text-[10px] text-slate-400 font-bold">+62 {{ phone }}</div>
                    </div>
                </button>
            </div>

            <button @click="router.back()" class="mt-10 text-[10px] font-black text-slate-300 uppercase tracking-widest hover:text-red-500 transition-colors">Batal & Kembali</button>
        </div>
    </div>
</template>