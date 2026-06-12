<script setup>
import { computed, watch } from 'vue';

const props = defineProps({
	formData: {
		type: Object,
		required: true,
	},
});

const kategoriOptions = [
	{ id: 'Retail', label: 'Retail & Barang', icon: 'M16 11V7a4 4 0 0 0-8 0v4M5 9h14l1 12H4L5 9z' },
	{ id: 'F&B', label: 'Food & Beverage', icon: 'M18 8h1a4 4 0 0 1 0 8h-1M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z M6 1v3 M10 1v3 M14 1v3' },
	{ id: 'Jasa', label: 'Layanan & Jasa', icon: 'M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z M3.27 6.96L12 12.01l8.73-5.05 M12 22.08V12' },
	{ id: 'Lainnya', label: 'Bisnis Lainnya', icon: 'M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z M22 6l-10 7L2 6' },
];

const detailOptions = {
	Retail: ['Minimarket / Toko Kelontong', 'Toko Pakaian / Butik', 'Elektronik / Gadget', 'Apotek / Farmasi', 'Toko Bangunan', 'Pet Shop', 'Lainnya'],
	'F&B': ['Restoran / Rumah Makan', 'Cafe / Coffee Shop', 'Bakery / Toko Roti', 'Food Court / Kaki Lima', 'Katering', 'Lainnya'],
	Jasa: ['Laundry', 'Barbershop / Salon', 'Cuci Mobil / Motor', 'Bengkel Otomotif', 'Klinik / Praktek', 'Lainnya'],
	Lainnya: ['Bisnis Umum / Lainnya'],
};

const activeKategori = computed(() => {
	return kategoriOptions.find((k) => k.id === props.formData.kategori_bisnis) || kategoriOptions[0];
});

watch(
	() => props.formData.kategori_bisnis,
	(newVal) => {
		props.formData.detail_bisnis = detailOptions[newVal][0];
	},
	{ immediate: true }
);

const formatNoHp = () => {
	let val = String(props.formData.telepon);
	if (val.startsWith('0')) val = val.substring(1);
	if (val.startsWith('62')) val = val.substring(2);
	props.formData.telepon = val;
};
</script>

<template>
	<div class="flex flex-col gap-6">
		<div class="flex items-center gap-3 border-b border-slate-100 pb-3">
			<div class="w-8 h-8 rounded-full bg-indigo-50 flex items-center justify-center text-indigo-600 font-black text-xs border border-indigo-100 shadow-sm">1</div>
			<h3 class="text-lg font-black text-slate-800 uppercase tracking-tight">Identitas Bisnis</h3>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-5 md:gap-6">
			<div class="md:col-span-2">
				<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Nama Brand / Toko</label>
				<input v-model="formData.nama_toko" type="text" required class="input-modern text-lg" placeholder="Contoh: Indomaret, Laundry Bersih, dsb..." />
			</div>

			<div class="md:col-span-2">
				<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Kategori Industri</label>
				<div class="p-4 rounded-[20px] border-2 border-indigo-200 bg-indigo-50/50 flex items-center gap-4 select-none opacity-90">
					<div class="w-12 h-12 bg-white rounded-xl shadow-sm border border-indigo-100 flex items-center justify-center">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
							<path :d="activeKategori.icon" />
						</svg>
					</div>
					<div class="flex flex-col">
						<span class="font-black text-slate-800 uppercase tracking-widest">{{ activeKategori.label }}</span>
						<span class="text-[10px] font-bold text-indigo-500 uppercase tracking-widest">Pilihan Dari Landing Page</span>
					</div>
					<div class="ml-auto text-slate-400" title="Kategori telah dikunci">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
							<rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
							<path d="M7 11V7a5 5 0 0 1 10 0v4" />
						</svg>
					</div>
				</div>
			</div>

			<div class="md:col-span-2 animate-[fadeIn_0.3s_ease-out]">
				<label class="text-[10px] font-black text-indigo-500 uppercase tracking-widest ml-1 mb-2 flex items-center gap-2">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
						<polyline points="9 18 15 12 9 6" />
					</svg>
					Spesifikasi {{ formData.kategori_bisnis }}
				</label>
				<div class="relative">
					<select v-model="formData.detail_bisnis" required class="input-modern bg-slate-50/50 cursor-pointer appearance-none text-indigo-900 border-indigo-100">
						<option v-for="opt in detailOptions[formData.kategori_bisnis]" :key="opt" :value="opt">{{ opt }}</option>
					</select>
					<div class="absolute inset-y-0 right-5 flex items-center pointer-events-none">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
							<path d="m6 9 6 6 6-6" />
						</svg>
					</div>
				</div>
			</div>

			<div class="md:col-span-2">
				<label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">No. WhatsApp Bisnis</label>
				<div class="flex items-center bg-white border-2 border-slate-200 rounded-2xl focus-within:border-indigo-500 focus-within:ring-4 focus-within:ring-indigo-500/10 transition-all shadow-sm overflow-hidden">
					<div class="pl-5 pr-4 py-4 bg-slate-50 border-r border-slate-200 flex items-center justify-center select-none">
						<span class="text-slate-500 font-black text-sm">+62</span>
					</div>
					<input v-model="formData.telepon" @input="formatNoHp" type="number" required class="w-full px-4 py-4 bg-transparent outline-none font-black text-slate-800 placeholder:text-slate-300 placeholder:font-bold" placeholder="81234567890" />
				</div>
			</div>
		</div>
	</div>
</template>
