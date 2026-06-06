<script setup>
import { onMounted } from 'vue';
import Sidebar from '../../components/Sidebar.vue';
import SettingBasicInfo from '../../components/settings/SettingBasicInfo.vue';
import SettingPayment from '../../components/settings/SettingPayment.vue';
import SettingReceiptTax from '../../components/settings/SettingReceiptTax.vue';
import { useStoreSetting } from '../../composables/useStoreSetting.js';

const { isLoading, isSaving, activeTab, form, logoPreview, qrisPreview, handleFileChange, removeLogo, removeQris, fetchSettings, saveSettings } = useStoreSetting();

const tabs = [
	{ id: 'basic', label: 'Info Dasar', icon: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z' },
	{ id: 'payment', label: 'Pembayaran', icon: 'M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z' },
	{ id: 'tax', label: 'Pajak & Struk', icon: 'M9 14l6-6m-5.5.5h.01m4.99 5h.01M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16l3.5-2 3.5 2 3.5-2 3.5 2zM10 8.5a.5.5 0 11-1 0 .5.5 0 011 0zm5 5a.5.5 0 11-1 0 .5.5 0 011 0z' },
];

// Eksekusi tarikan data backend saat halaman pertama kali dibuka 
onMounted(() => {
	fetchSettings();
});
</script>

<template>
	<Sidebar>
		<div class="p-4 md:p-8 lg:p-10 max-w-5xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
			<div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
				<div>
					<h1 class="text-2xl md:text-3xl font-black text-slate-800 tracking-tighter uppercase">Pengaturan Toko</h1>
					<p class="text-[10px] md:text-xs font-black text-slate-400 uppercase tracking-widest mt-1">Konfigurasi Identitas, Pembayaran & Struk Tenant</p>
				</div>
				<button @click="saveSettings" :disabled="isSaving || isLoading" class="px-8 py-3.5 bg-blue-600 hover:bg-slate-900 text-white rounded-2xl font-black text-[10px] md:text-xs uppercase tracking-widest shadow-xl transition-all disabled:opacity-50 flex items-center justify-center gap-2">
					<span v-if="isSaving" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
					<svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
					{{ isSaving ? 'Menyimpan...' : 'Simpan Perubahan' }}
				</button>
			</div>

			<div v-if="isLoading" class="py-20 flex flex-col items-center justify-center bg-white rounded-[32px] border border-slate-100 shadow-sm">
				<div class="w-10 h-10 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin mb-4"></div>
				<p class="text-slate-400 font-black text-[10px] uppercase tracking-widest animate-pulse">Menarik Data Konfigurasi...</p>
			</div>

			<div v-else class="bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden flex flex-col md:flex-row">
				<div class="w-full md:w-64 bg-slate-50/50 border-r border-slate-100 p-6 flex flex-row md:flex-col gap-2 overflow-x-auto custom-scrollbar">
					<button v-for="tab in tabs" :key="tab.id" @click="activeTab = tab.id" :class="['flex items-center gap-3 p-3.5 rounded-2xl font-black text-[10px] uppercase tracking-widest transition-all whitespace-nowrap', activeTab === tab.id ? 'bg-white text-blue-600 shadow-sm border border-slate-200' : 'text-slate-500 hover:bg-slate-100']">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
							<path stroke-linecap="round" stroke-linejoin="round" :d="tab.icon" />
						</svg>
						{{ tab.label }}
					</button>
				</div>

				<div class="flex-1 p-6 md:p-8 overflow-hidden">
					<SettingBasicInfo v-show="activeTab === 'basic'" :form="form" :logoPreview="logoPreview" @update-file="handleFileChange" @remove-logo="removeLogo" />
					<SettingPayment v-show="activeTab === 'payment'" :form="form" :qrisPreview="qrisPreview" @update-file="handleFileChange" @remove-qris="removeQris" />
					<SettingReceiptTax v-show="activeTab === 'tax'" :form="form" />
				</div>
			</div>
		</div>
	</Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
	height: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
	background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
	background: #cbd5e1;
	border-radius: 10px;
}
</style>
