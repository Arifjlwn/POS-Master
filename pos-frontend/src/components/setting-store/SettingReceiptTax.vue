<script setup>
const props = defineProps({ form: Object });

// Memaksa input footer struk menjadi huruf kapital sesuai standar printer thermal
const forceUppercase = (field) => {
	if (props.form[field]) {
		props.form[field] = props.form[field].toUpperCase();
	}
};

// Mengunci angka pajak agar tidak minus atau meledak lebih dari 100%
const limitTaxPercentage = () => {
	if (props.form.pajak_persen < 0) props.form.pajak_persen = 0;
	if (props.form.pajak_persen > 100) props.form.pajak_persen = 100;
};
</script>

<template>
	<div class="space-y-6 animate-fade-in-up">
		<div class="p-6 border-2 border-slate-100 rounded-3xl bg-white">
			<div class="flex items-center justify-between mb-6">
				<div>
					<h4 class="font-black text-slate-800 uppercase tracking-widest text-sm">Pajak Pelanggan (PPN/PB1)</h4>
					<p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">Aktifkan untuk membebankan pajak otomatis pada struk</p>
				</div>
				<label class="relative inline-flex items-center cursor-pointer">
					<input type="checkbox" v-model="form.is_tax_active" class="sr-only peer" />
					<div class="w-14 h-7 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-6 after:w-6 after:transition-all peer-checked:bg-emerald-500"></div>
				</label>
			</div>

			<div v-if="form.is_tax_active" class="animate-fade-in-up">
				<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Persentase Pajak (%)</label>
				<input v-model="form.pajak_persen" @blur="limitTaxPercentage" type="number" step="0.1" min="0" max="100" class="w-full md:w-1/3 p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-emerald-500 outline-none font-black text-lg text-emerald-600 transition-colors" />
				<p class="text-[9px] font-bold text-slate-400 mt-2">Batas input 0% hingga 100%</p>
			</div>
		</div>

		<div class="p-6 border-2 border-slate-100 rounded-3xl bg-white">
			<h4 class="font-black text-slate-800 uppercase tracking-widest text-sm mb-1">Pengaturan Cetak Struk</h4>
			<p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-6">Konfigurasi hardware dan identitas nota pelanggan</p>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-5 mb-6">
				<div>
					<label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-2 block">Ukuran Kertas Printer</label>
					<select v-model="form.printer_width" class="w-full p-4 bg-slate-50 border-2 border-slate-100 focus:border-blue-500 rounded-2xl outline-none font-black text-xs uppercase text-slate-800 cursor-pointer transition-all">
						<option value="58mm">Thermal Kecil (58mm)</option>
						<option value="80mm">Thermal Besar (80mm)</option>
					</select>
				</div>
				<div>
					<label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-2 block">Tipe Koneksi Printer</label>
					<select v-model="form.printer_type" class="w-full p-4 bg-slate-50 border-2 border-slate-100 focus:border-blue-500 rounded-2xl outline-none font-black text-xs uppercase text-slate-800 cursor-pointer transition-all">
						<option value="bluetooth">Bluetooth (HP / Tablet)</option>
						<option value="usb">Kabel USB (Laptop / PC)</option>
						<option value="lan">LAN / WiFi (Jaringan)</option>
					</select>
				</div>
			</div>

			<div class="pt-6 border-t border-slate-100">
				<label class="text-[9px] font-black text-slate-800 uppercase tracking-widest mb-2 block">Catatan Kaki Struk (Footer)</label>
				<textarea v-model="form.receipt_footer" @input="forceUppercase('receipt_footer')" rows="3" placeholder="TERIMA KASIH ATAS KUNJUNGAN ANDA. BARANG YANG DIBELI TIDAK DAPAT DITUKAR." class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm text-center resize-none transition-colors uppercase"></textarea>
			</div>
		</div>
	</div>
</template>
