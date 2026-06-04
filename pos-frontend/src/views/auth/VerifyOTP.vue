<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Swal from 'sweetalert2';
import api from '../../api.js';

const route = useRoute();
const router = useRouter();

const email = route.query.email || '';
const phone = route.query.phone || '';
const method = route.query.method || 'email';
const intent = route.query.intent || 'register';
const otp = ref('');
const isLoading = ref(false);

// --- LOGIKA TIMER & RESEND OTP ---
const timeLeft = ref(0);
let timerInterval = null;
const storageKey = computed(() => `nexa_otp_expiry_${email || phone}`);

const startTimer = (isResend = false) => {
	clearInterval(timerInterval);
	let expiryTime = localStorage.getItem(storageKey.value);

	if (!expiryTime || isResend) {
		expiryTime = Date.now() + (180 * 1000);
		localStorage.setItem(storageKey.value, expiryTime);
	}

	const updateTimer = () => {
		const now = Date.now();
		const diff = Math.floor((expiryTime - now) / 1000);

		if (diff > 0) {
			timeLeft.value = diff;
		} else {
			timeLeft.value = 0;
			clearInterval(timerInterval);
			localStorage.removeItem(storageKey.value);
		}
	};

	updateTimer();
	timerInterval = setInterval(updateTimer, 1000);
};

const formattedTime = computed(() => {
	const minutes = Math.floor(timeLeft.value / 60);
	const seconds = timeLeft.value % 60;
	return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
});

const isTimerActive = computed(() => timeLeft.value > 0);

// --- HIT API VERIFIKASI ---
const handleVerify = async () => {
	if (otp.value.length < 6) return;
	isLoading.value = true;

	try {
		await api.post('/verify-otp', {
			email: email,
			phone: phone,
			otp: otp.value,
			intent: intent
		});
		
		if (intent === 'reset-password') {
			await Swal.fire({
				icon: 'success',
				title: 'Kode OTP Valid!',
				text: 'Silakan buat password baru Anda bray.',
				confirmButtonColor: '#4f46e5'
			});
			localStorage.removeItem(storageKey.value);
			router.push({ path: '/reset-password', query: { email: email, token: otp.value } });
		} else {
			await Swal.fire({
				icon: 'success',
				title: 'Verifikasi Sukses!',
				text: 'Akun Owner Anda sudah aktif, mari login dan mulai setup toko.',
				confirmButtonColor: '#4f46e5'
			});
			localStorage.removeItem(storageKey.value);
			router.push('/login');
		}

	} catch (error) {
		Swal.fire({
			icon: 'error',
			title: 'Verifikasi Gagal',
			text: error.response?.data?.error || 'Kode OTP salah atau sudah kadaluarsa bray.',
			confirmButtonColor: '#ef4444'
		});
	} finally {
		isLoading.value = false;
	}
};

// --- HIT API RESEND OTP ---
const handleResendOTP = async () => {
	if (isTimerActive.value) return;
	isLoading.value = true;

	Swal.fire({
		title: 'Mengirim Ulang OTP...',
		text: 'Mohon tunggu sejenak bray...',
		allowOutsideClick: false,
		didOpen: () => { Swal.showLoading(); }
	});

	try {
		// 🚀 FIX: Menggunakan local state variable secara konsisten & aman untuk payload resend bray
		if (intent === 'reset-password') {
			await api.post('/auth/send-otp-wa', { phone: phone });
		} else {
			// Menggunakan utility resend OTP terpadu ke gateway WhatsApp bray
			await api.post('/auth/send-otp-wa', { phone: phone });
		}

		Swal.close();
		Swal.fire({
			icon: 'success',
			title: 'Kode OTP Baru Dikirim!',
			text: 'Silakan cek kembali kotak masuk WhatsApp Anda bray.',
			confirmButtonColor: '#4f46e5'
		});

		otp.value = ''; 
		startTimer(true);
	} catch (error) {
		Swal.close();
		Swal.fire({
			icon: 'error',
			title: 'Gagal Kirim Ulang',
			text: error.response?.data?.error || 'Sistem mendeteksi limitasi kuota kirim, coba lagi nanti bray.',
			confirmButtonColor: '#ef4444'
		});
	} finally {
		isLoading.value = false;
	}
};

onMounted(() => { startTimer() });
onUnmounted(() => { clearInterval(timerInterval) });
</script>

<template>
	<div class="min-h-screen bg-[#F8FAFC] font-sans text-slate-900 selection:bg-indigo-100 selection:text-indigo-600 overflow-x-hidden antialiased flex items-center justify-center p-6 relative">
		
		<div class="absolute -top-24 -right-24 w-96 h-96 bg-blue-100/40 rounded-full blur-3xl pointer-events-none"></div>
		<div class="absolute -bottom-24 -left-24 w-96 h-96 bg-indigo-100/40 rounded-full blur-3xl pointer-events-none"></div>

		<div class="w-full max-w-md relative z-10">
			<div class="bg-white rounded-[40px] p-10 shadow-2xl border border-white text-center">
				
				<div class="w-20 h-20 bg-indigo-50 border border-indigo-100/30 text-indigo-600 rounded-[28px] flex items-center justify-center mx-auto mb-6 shadow-sm">
					<svg v-if="method === 'email'" xmlns="http://www.w3.org/2000/svg" class="w-9 h-9" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>
					<svg v-else xmlns="http://www.w3.org/2000/svg" class="w-9 h-9" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
				</div>

				<h2 class="text-3xl font-black text-slate-900 tracking-tighter uppercase">INPUT OTP</h2>
				<p class="text-slate-400 font-bold text-[10px] uppercase tracking-widest mt-2 mb-8 leading-relaxed">
					Kode rahasia 6-digit telah dikirim ke <br/>
					<span class="text-indigo-600 font-black text-xs normal-case tracking-normal">
						{{ method === 'email' ? email : (phone || '').replace(/^62/, '0') }}
					</span>
				</p>

				<form @submit.prevent="handleVerify" class="space-y-6">
					<div class="relative">
						<input v-model="otp" type="text" maxlength="6" pattern="[0-9]*" inputmode="numeric" class="w-full text-center text-4xl font-black tracking-[0.4em] pl-[0.4em] py-5 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:bg-white focus:border-indigo-600 outline-none transition-all placeholder:text-slate-200" placeholder="000000" required>
					</div>

					<div class="flex items-center justify-center gap-2 px-4 py-3 bg-slate-50 border border-slate-100 rounded-xl text-xs">
						<span class="font-bold text-slate-400 uppercase tracking-wider">Sesi Berlaku:</span>
						<span :class="['font-black tracking-wider', timeLeft < 30 ? 'text-red-500 animate-pulse' : 'text-indigo-600']">
							{{ formattedTime }}
						</span>
					</div>

					<button type="submit" :disabled="otp.length < 6 || isLoading || !isTimerActive" class="w-full py-5 rounded-2xl bg-indigo-600 text-white font-black text-xs uppercase tracking-widest shadow-xl shadow-indigo-200 hover:bg-slate-900 transition-all active:scale-95 disabled:opacity-30 disabled:pointer-events-none">
						<span v-if="isLoading">MEMPROSES...</span>
						<span v-else>{{ intent === 'reset-password' ? 'Validasi & Buat Password Baru' : 'Aktifkan Akun Owner' }}</span>
					</button>
				</form>

				<div class="mt-8 pt-6 border-t border-slate-100">
					<p class="text-xs font-bold text-slate-400 uppercase tracking-tight">
						Tidak menerima kode OTP? <br class="sm:hidden"/>
						<button @click="handleResendOTP" :disabled="isTimerActive || isLoading" :class="['ml-1 font-black transition-colors outline-none focus:underline', isTimerActive ? 'text-slate-300 cursor-not-allowed' : 'text-indigo-600 hover:text-slate-900']">
							Kirim Ulang Kode
						</button>
					</p>
				</div>

			</div>
		</div>
	</div>
</template>

<style scoped>
input::-webkit-outer-spin-button,
input::-webkit-inner-spin-button {
	-webkit-appearance: none;
	margin: 0;
}
input[type=number] { -moz-appearance: textfield; }
</style>