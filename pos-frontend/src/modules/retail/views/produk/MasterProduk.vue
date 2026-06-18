<script setup>
import { onMounted, ref } from 'vue';
import api from '../../../../api.js';
import Sidebar from '../../components/Sidebar.vue';
import { useMasterProduk } from '../../composables/useMasterProduk.js';

// Import Kepingan Komponen Baru
import ProductFormModal from '../../components/inventory/ProductFormModal.vue';
import ProductHeader from '../../components/inventory/ProductHeader.vue';
import ProductTable from '../../components/inventory/ProductTable.vue';

const { products, categories, isLoading, isSubmitting, stok_dalam_karton, eceran_tambahan, searchQuery, selectedCategory, currentPage, totalPages, showFormModal, isEditing, imagePreview, form, showScanner, visiblePages, getImageUrl, changePage, onFileChange, openAddModal, editProduct, submitProduct, deleteProduct, triggerImport, handleImport, exportCSV, startScanner, stopScanner } = useMasterProduk();

// 🚀 NEGARA REPUBLIK STATE: Wadah penampung profile konfigurasi bisnis tenant
const storeData = ref({});

// 🚀 AMBIL LOKASI IDENTITAS BISNIS SEBELUM MODAL DI-RENDER BRAY
onMounted(async () => {
	try {
		const res = await api.get('/store/settings');
		storeData.value = res.data.data;
	} catch (e) {
		console.error('Gagal memuat tipe bisnis merchant di katalog produk:', e);
	}
});
</script>

<template>
	<Sidebar>
		<div class="p-4 md:p-8 lg:p-10 w-full max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
			<ProductHeader v-model:searchQuery="searchQuery" v-model:selectedCategory="selectedCategory" :categories="categories" @export="exportCSV" @trigger-import="triggerImport" @handle-import="handleImport" @add-new="openAddModal" />

			<ProductTable :products="products" :isLoading="isLoading" :currentPage="currentPage" :totalPages="totalPages" :visiblePages="visiblePages" :getImageUrl="getImageUrl" @edit="editProduct" @delete="deleteProduct" @change-page="changePage" />

			<ProductFormModal :show="showFormModal" :isEditing="isEditing" :isSubmitting="isSubmitting" :form="form" :categories="categories" :imagePreview="imagePreview" v-model:stokDalamKarton="stok_dalam_karton" v-model:eceranTambahan="eceran_tambahan" :businessSubType="storeData?.business_type" @close="showFormModal = false" @submit="submitProduct" @start-scanner="startScanner" @file-change="onFileChange" />

			<div v-if="showScanner" class="fixed inset-0 bg-slate-900/90 flex items-center justify-center z-[150] p-4 backdrop-blur-sm">
				<div class="bg-white rounded-[32px] shadow-2xl w-full max-w-md overflow-hidden flex flex-col">
					<div class="p-6 border-b border-slate-50 flex justify-between items-center bg-slate-50/50">
						<h2 class="text-lg font-black text-slate-800 uppercase tracking-widest flex items-center gap-2">Scan Barcode Produk</h2>
						<button @click="stopScanner" class="p-2 rounded-xl bg-slate-100 text-slate-400 hover:text-red-500 transition-all">✕</button>
					</div>
					<div class="p-6 bg-black relative">
						<div id="reader" class="w-full rounded-2xl overflow-hidden border-2 border-slate-700"></div>
					</div>
				</div>
			</div>
		</div>
	</Sidebar>
</template>
