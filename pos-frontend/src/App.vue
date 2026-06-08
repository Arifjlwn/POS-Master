<script setup>
import { onMounted, ref } from 'vue';

const deferredPrompt = ref(null);
const showInstallModal = ref(false);

onMounted(() => {
	// Tangkap event install bawaan browser jika belum ter-install di HP
	window.addEventListener('beforeinstallprompt', (e) => {
		e.preventDefault();
		deferredPrompt.value = e;
		showInstallModal.value = true;
	});

	// Jika user akhirnya berhasil menginstall, langsung tutup modalnya
	window.addEventListener('appinstalled', () => {
		deferredPrompt.value = null;
		showInstallModal.value = false;
		console.log('Arzura POS berhasil dipasang di perangkat!');
	});
});

const triggerInstall = async () => {
	if (!deferredPrompt.value) return;

	// Munculkan pop-up instalasi bawaan browser (Chrome/Edge/dll)
	deferredPrompt.value.prompt();

	// Cek apakah user klik "Install" atau "Cancel"
	const { outcome } = await deferredPrompt.value.userChoice;
	console.log(`Pilihan user: ${outcome}`);

	// Reset prompt dan tutup modal kustom kita
	deferredPrompt.value = null;
	showInstallModal.value = false;
};
</script>

<template>
	<!-- Router view utama Arzura POS -->
	<router-view></router-view>

	<!-- Pop-Up Install (Menggunakan Tailwind & AutoAnimate) -->
	<div v-if="showInstallModal" v-auto-animate class="fixed inset-0 bg-black/70 flex items-center justify-center z-[9999] p-4">
		<!-- Tambahkan class 'relative' agar bisa memosisikan tombol silang -->
		<div class="relative bg-white p-6 rounded-2xl max-w-sm w-full text-center shadow-2xl border border-gray-100">
			<!-- PILIHAN 1: Tombol Silang (X) di Pojok Kanan Atas -->
			<button @click="showInstallModal = false" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600 transition" aria-label="Tutup">
				<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-6 h-6">
					<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
				</svg>
			</button>

			<!-- Wrapper Logo -->
			<div class="flex justify-center mb-4 mt-2">
				<img src="/logo-192.png" alt="Arzura POS Logo" class="w-20 h-20 rounded-2xl object-cover shadow-md" />
			</div>

			<!-- Teks Informasi -->
			<h3 class="text-xl font-bold text-gray-900 mb-2">Install Arzura POS</h3>
			<p class="text-sm text-gray-600 leading-relaxed mb-6">Untuk performa kasir yang lebih cepat, stabil, dan optimal, disarankan menginstall aplikasi di HP Anda.</p>

			<!-- Tombol Aksi -->
			<div class="space-y-2">
				<button @click="triggerInstall" class="w-full bg-black text-white py-3 px-4 rounded-xl font-semibold hover:bg-gray-800 transition duration-200 active:scale-95">Install Sekarang</button>

				<!-- PILIHAN 2: Tombol Teks "Nanti Saja" di bawah (Biar makin user-friendly) -->
				<button @click="showInstallModal = false" class="w-full bg-transparent text-gray-500 py-2 text-sm font-medium hover:text-gray-800 transition">Nanti Saja</button>
			</div>
		</div>
	</div>
</template>
