<script setup>
import Swal from 'sweetalert2';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../../api.js';

const router = useRouter();
const email = ref('');
const password = ref('');
const isLoading = ref(false);

const handleAdminLogin = async () => {
	if (!email.value || !password.value) return;

	isLoading.value = true;
	try {
		// 🚀 FIX: Ubah 'identifier' menjadi 'email' agar sinkron dengan struct GIN Go!
		const res = await api.post('/admin/login', {
			email: email.value,
			password: password.value,
		});

		// Simpan token kedaulatan kasta tertinggi
		localStorage.setItem('token', res.data.token);
		localStorage.setItem('role', res.data.user.role);
		localStorage.setItem('name', res.data.user.name);

		Swal.fire({
			icon: 'success',
			title: 'Akses Diberikan!',
			text: 'Membuka sistem telemetri Mission Control...',
			timer: 1500,
			showConfirmButton: false,
			customClass: { popup: 'rounded-[24px]' },
		}).then(() => {
			router.push('/admin/dashboard');
		});
	} catch (err) {
		Swal.fire({
			icon: 'error',
			title: 'Akses Ditolak!',
			text: err.response?.data?.error || 'Sistem mendeteksi anomali kredensial.',
			confirmButtonColor: '#ef4444',
			customClass: { popup: 'rounded-[24px]' },
		});
	} finally {
		isLoading.value = false;
	}
};
</script>

<template>
	<div class="min-h-screen bg-[#0B0F19] flex items-center justify-center p-4 font-sans select-none relative overflow-hidden">
		<div class="absolute w-[500px] h-[500px] bg-indigo-900/20 rounded-full blur-[120px] -top-40 -left-40"></div>
		<div class="absolute w-[400px] h-[400px] bg-purple-900/10 rounded-full blur-[100px] -bottom-20 -right-20"></div>

		<div class="w-full max-w-md bg-[#131B2E] border border-slate-800 rounded-[32px] p-8 md:p-10 shadow-2xl relative z-10 transition-all duration-300">
			<div class="text-center mb-8">
				<span class="px-3 py-1 bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl text-[10px] font-black uppercase tracking-widest">Pusat Otoritas Tertinggi</span>
				<h1 class="text-2xl font-black text-white tracking-tight mt-3">ARZURA POS FOUNDER</h1>
				<p class="text-slate-400 text-xs font-medium mt-1">Silakan verifikasi kunci enkripsi admin Anda.</p>
			</div>

			<form @submit.prevent="handleAdminLogin" class="space-y-5">
				<div>
					<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Email Root</label>
					<input v-model="email" type="email" required placeholder="founder@arzura.com" class="w-full px-5 py-4 bg-[#0B0F19] border border-slate-800 rounded-2xl text-white text-sm font-bold focus:outline-none focus:border-indigo-500 transition-all placeholder:text-slate-600" />
				</div>

				<div>
					<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Password Sandi</label>
					<input v-model="password" type="password" required placeholder="••••••••" class="w-full px-5 py-4 bg-[#0B0F19] border border-slate-800 rounded-2xl text-white text-sm font-bold focus:outline-none focus:border-indigo-500 transition-all placeholder:text-slate-600" />
				</div>

				<button type="submit" :disabled="isLoading" class="w-full py-4 mt-2 bg-indigo-600 hover:bg-indigo-700 active:scale-[0.98] disabled:opacity-50 text-white font-black text-xs uppercase tracking-widest rounded-2xl shadow-lg shadow-indigo-600/20 transition-all flex items-center justify-center gap-2 relative overflow-hidden">
					<span :class="{ 'opacity-0 invisible': isLoading }" class="transition-all">Otorisasi Masuk</span>
					<div v-if="isLoading" class="absolute inset-0 flex items-center justify-center bg-indigo-700">
						<svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
					</div>
				</button>
			</form>
		</div>
	</div>
</template>
