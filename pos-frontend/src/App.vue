<template>
	<router-view></router-view>

	<!-- <div v-auto-animate>
		<div v-if="showInstallModal" class="fixed inset-0 bg-black/70 flex items-center justify-center z-[9999] p-4">
			<div class="relative bg-white p-6 rounded-2xl max-w-sm w-full text-center shadow-2xl border border-gray-100">
				<button @click="closeAndDeferModal" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600 transition" aria-label="Tutup">
					<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-6 h-6">
						<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>

				<div class="flex justify-center mb-4 mt-2">
					<img src="/logo-192.png" alt="Arzura POS Logo" class="w-20 h-20 rounded-2xl object-cover shadow-md" />
				</div>

				<h3 class="text-xl font-bold text-gray-900 mb-2">Install Arzura POS</h3>
				<p class="text-sm text-gray-600 leading-relaxed mb-6">Untuk performa kasir yang lebih cepat, stabil, dan optimal, disarankan menginstall aplikasi di HP Anda.</p>

				<div class="space-y-2">
					<button @click="triggerInstall" class="w-full bg-black text-white py-3 px-4 rounded-xl font-semibold hover:bg-gray-800 transition duration-200 active:scale-95">Install Sekarang</button>

					<button @click="closeAndDeferModal" class="w-full bg-transparent text-gray-500 py-2 text-sm font-medium hover:text-gray-800 transition">Nanti Saja</button>
				</div>
			</div>
		</div>
	</div> -->
</template>

<script setup>
import { onMounted, ref } from 'vue';

const deferredPrompt = ref(null);
const showInstallModal = ref(false);

const shouldShowModal = () => {
	const hideUntil = localStorage.getItem('arzura_pwa_hide_until');
	if (!hideUntil) return true;

	const now = new Date().getTime();
	return now > parseInt(hideUntil, 10);
};

onMounted(() => {
	window.addEventListener('beforeinstallprompt', (e) => {
		e.preventDefault();
		deferredPrompt.value = e;

		if (shouldShowModal()) {
			showInstallModal.value = true;
		}
	});

	window.addEventListener('appinstalled', () => {
		deferredPrompt.value = null;
		showInstallModal.value = false;
		localStorage.removeItem('arzura_pwa_hide_until');
	});
});

const triggerInstall = async () => {
	if (!deferredPrompt.value) return;

	deferredPrompt.value.prompt();
	const { outcome } = await deferredPrompt.value.userChoice;

	deferredPrompt.value = null;
	showInstallModal.value = false;
};

const closeAndDeferModal = () => {
	showInstallModal.value = false;
	const threeDaysInMs = 1 * 24 * 60 * 60 * 1000;
	const hideUntilTimestamp = new Date().getTime() + threeDaysInMs;
	localStorage.setItem('arzura_pwa_hide_until', hideUntilTimestamp.toString());
};
</script>
