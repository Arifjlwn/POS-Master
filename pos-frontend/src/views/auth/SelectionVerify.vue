<script setup>
import { useRoute, useRouter } from 'vue-router';
import Swal from 'sweetalert2';
import api from '../../api.js';

const route = useRoute();
const router = useRouter();
const email = route.query.email;
const phone = route.query.phone;
const intent = route.query.intent || 'register';

const selectMethod = async (method) => {
    // --------------------------------------------------------
    // OPSI 1: JIKA USER MEMILIH VERIFIKASI VIA EMAIL
    // --------------------------------------------------------
    if (method === 'email') {
        Swal.fire({
            title: 'Mengirim OTP...',
            text: 'Mohon tunggu, kami sedang mengirim kode ke Email Anda',
            allowOutsideClick: false,
            didOpen: () => { Swal.showLoading(); }
        });

        try {
            // Pemicu kirim OTP via Email ke Backend
            await api.post('/auth/send-otp-email', { email: email });
            Swal.close();

            Swal.fire({
                icon: 'success',
                title: 'Email Terkirim!',
                text: 'Kode OTP aman sudah meluncur ke kotak masuk/spam Email Anda.',
                confirmButtonColor: '#4f46e5'
            }).then(() => {
                router.push({ path: '/verify-otp', query: { email: email, phone: phone, method: 'email', intent: intent } });
            });
        } catch (error) {
            Swal.close();
            Swal.fire({
                icon: 'error',
                title: 'Gagal Mengirim Email',
                text: error.response?.data?.error || 'Layanan email server sedang sibuk, coba lagi nanti.',
                confirmButtonColor: '#ef4444'
            });
        }
        return;
    }
    
    // --------------------------------------------------------
    // OPSI 2: JIKA USER MEMILIH VERIFIKASI VIA WHATSAPP
    // --------------------------------------------------------
    if (method === 'whatsapp') {
        Swal.fire({
            title: 'Mengirim OTP...',
            text: 'Mohon tunggu, kami sedang mengirim kode ke WhatsApp Anda',
            allowOutsideClick: false,
            didOpen: () => { Swal.showLoading(); }
        });

        try {
            await api.post('/auth/send-otp-wa', { phone: phone });
            Swal.close();
            
            Swal.fire({
                icon: 'success',
                title: 'Berhasil Sent!',
                text: 'Kode OTP aman sudah meluncur ke WhatsApp Anda.',
                confirmButtonColor: '#4f46e5'
            }).then(() => {
                router.push({ path: '/verify-otp', query: { email: email, phone: phone, method: 'whatsapp', intent: intent } });
            });
            
        } catch (error) {
            Swal.close();
            Swal.fire({
                icon: 'error',
                title: 'Gagal Mengirim WA',
                text: error.response?.data?.error || 'Server gateway Fonnte sedang sibuk, coba lagi nanti.',
                confirmButtonColor: '#ef4444'
            });
        }
    }
};
</script>

<template>
	<div class="min-h-screen bg-[#F8FAFC] font-sans text-slate-900 selection:bg-indigo-100 selection:text-indigo-600 overflow-x-hidden antialiased flex items-center justify-center p-6 relative">
		
		<div class="absolute -top-24 -right-24 w-96 h-96 bg-blue-100/50 rounded-full blur-3xl pointer-events-none"></div>
		<div class="absolute -bottom-24 -left-24 w-96 h-96 bg-indigo-100/50 rounded-full blur-3xl pointer-events-none"></div>

		<div class="w-full max-w-md relative">
			<div class="bg-white rounded-[40px] p-10 shadow-2xl border border-white text-center">
				<h2 class="text-3xl font-black text-slate-900 tracking-tighter">VERIFIKASI AKUN</h2>
				<p class="text-slate-400 font-bold text-[10px] uppercase tracking-[0.2em] mt-2 mb-10">Pilih metode pengiriman kode OTP</p>

				<div class="space-y-4">
					<button @click="selectMethod('email')" class="w-full p-6 border-2 border-slate-50 bg-slate-50/50 rounded-3xl flex items-center gap-4 hover:border-blue-600 hover:bg-blue-50/50 transition-all group focus:outline-none">
						<div class="w-12 h-12 bg-white border border-slate-100 text-blue-600 rounded-2xl flex items-center justify-center shadow-sm group-hover:bg-blue-600 group-hover:border-blue-600 group-hover:text-white transition-all shrink-0">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>
						</div>
						<div class="text-left overflow-hidden">
							<div class="text-xs font-black text-slate-800 uppercase tracking-wide">Kirim via Email</div>
							<div class="text-[11px] text-slate-400 font-bold mt-1 truncate">{{ email }}</div>
						</div>
					</button>

					<button @click="selectMethod('whatsapp')" class="w-full p-6 border-2 border-slate-50 bg-slate-50/50 rounded-3xl flex items-center gap-4 hover:border-green-600 hover:bg-green-50/50 transition-all group focus:outline-none">
						<div class="w-12 h-12 bg-white border border-slate-100 text-green-600 rounded-2xl flex items-center justify-center shadow-sm group-hover:bg-green-600 group-hover:border-green-600 group-hover:text-white transition-all shrink-0">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
						</div>
						<div class="text-left">
							<div class="text-xs font-black text-slate-800 uppercase tracking-wide">Kirim via WhatsApp</div>
							<div class="text-[11px] text-slate-400 font-bold mt-1">+{{ phone }}</div>
						</div>
					</button>
				</div>

				<button @click="router.back()" class="mt-10 text-[10px] font-black text-slate-300 uppercase tracking-widest hover:text-red-500 transition-colors focus:outline-none">Batal & Kembali</button>
			</div>
		</div>
	</div>
</template>