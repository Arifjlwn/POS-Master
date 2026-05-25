<script setup>
import { useMasterProduk } from '../../composables/useMasterProduk.js';
import Sidebar from '../../components/Sidebar.vue';

// 🚀 Destructuring seluruh kepingan reaktif & fungsi dari Composable
const {
    products, categories, isLoading, isSubmitting, stok_dalam_karton, eceran_tambahan,
    searchQuery, selectedCategory, currentPage, totalPages, showFormModal, isEditing,
    fileInput, imagePreview, form, showScanner, visiblePages,
    getImageUrl, changePage, onFileChange, openAddModal, editProduct, submitProduct,
    deleteProduct, triggerImport, handleImport, exportCSV, startScanner, stopScanner
} = useMasterProduk();
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 lg:p-10 w-full max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <div class="bg-gradient-to-br from-slate-900 via-slate-800 to-indigo-900 rounded-[32px] p-6 md:p-10 mb-8 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden">
                <div class="z-10 text-center md:text-left">
                    <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-1 uppercase italic">Master <span class="text-blue-400">Inventory</span></h1>
                    <p class="text-slate-300 font-bold text-[10px] uppercase tracking-[0.2em]">Katalog Produk & Harga Gudang</p>
                </div>

                <div class="z-10 mt-6 md:mt-0 flex flex-wrap justify-center gap-3">
                    <button @click="exportCSV" class="bg-emerald-500/20 hover:bg-emerald-500 text-emerald-400 hover:text-white px-5 py-3 rounded-[16px] font-black transition-all flex items-center gap-2 text-[10px] uppercase tracking-widest border border-emerald-500/50">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
                        Ekspor
                    </button>
                    <button @click="triggerImport" class="bg-amber-500/20 hover:bg-amber-500 text-amber-400 hover:text-white px-5 py-3 rounded-[16px] font-black transition-all flex items-center gap-2 text-[10px] uppercase tracking-widest border border-amber-500/50">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>
                        Impor
                    </button>
                    <input type="file" ref="fileInput" class="hidden" accept=".csv" @change="handleImport">
                    <button @click="openAddModal" class="bg-blue-600 hover:bg-blue-500 text-white px-6 py-3 rounded-[16px] font-black transition-all shadow-lg flex items-center gap-2 text-[10px] uppercase tracking-widest active:scale-95 border border-blue-400">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" /></svg>
                        Tambah
                    </button>
                </div>
            </div>

            <div class="flex flex-col sm:flex-row gap-4 mb-6">
                <div class="relative flex-1 group">
                    <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-blue-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                    </div>
                    <input v-model="searchQuery" type="text" placeholder="Cari nama barang atau barcode..." class="block w-full pl-14 pr-4 py-4 bg-white rounded-2xl border-2 border-slate-100 shadow-sm focus:border-blue-600 outline-none font-bold text-sm transition-all text-slate-700">
                </div>
                <div class="w-full sm:w-64 shrink-0 relative">
                    <select v-model="selectedCategory" class="block w-full pl-4 pr-10 py-4 bg-white rounded-2xl border-2 border-slate-100 shadow-sm focus:border-blue-600 text-sm font-bold text-slate-700 cursor-pointer outline-none appearance-none transition-all uppercase">
                        <option value="">SEMUA KATEGORI</option>
                        <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
                    </select>
                    <div class="absolute inset-y-0 right-0 pr-4 flex items-center pointer-events-none">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
                    </div>
                </div>
            </div>

            <div class="w-full bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden flex flex-col">
                <div class="overflow-x-auto custom-scrollbar">
                    <table class="w-full text-left border-collapse whitespace-nowrap min-w-[800px]">
                        <thead>
                            <tr class="bg-slate-50 border-b border-slate-100">
                                <th class="p-5 font-black uppercase tracking-widest text-[10px] text-slate-400">Produk & Detail</th>
                                <th class="p-5 font-black uppercase tracking-widest text-[10px] text-slate-400">Kategori</th>
                                <th class="p-5 font-black uppercase tracking-widest text-[10px] text-right text-slate-400">Modal Dasar</th>
                                <th class="p-5 font-black uppercase tracking-widest text-[10px] text-right text-slate-400">Harga Jual</th>
                                <th class="p-5 font-black uppercase tracking-widest text-[10px] text-center text-slate-400">Stok</th>
                                <th class="p-5 font-black uppercase tracking-widest text-[10px] text-center text-slate-400">Aksi</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-50">
                            <tr v-if="isLoading">
                                <td colspan="6" class="p-16 text-center text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse">Memuat Katalog Produk...</td>
                            </tr>
                            <tr v-else-if="products.length === 0">
                                <td colspan="6" class="p-16 text-center text-slate-400 font-black text-xs uppercase tracking-widest opacity-50">Produk Tidak Ditemukan</td>
                            </tr>
                            <tr v-else v-for="product in products" :key="product.id" class="hover:bg-blue-50/30 transition-colors group">
                                <td class="p-5 flex items-center gap-4">
                                    <div class="w-12 h-12 md:w-14 md:h-14 rounded-[16px] border-2 border-white shadow-sm bg-slate-50 flex items-center justify-center text-slate-300 overflow-hidden shrink-0 group-hover:border-blue-200 transition-colors">
                                        <img v-if="product.gambar" :src="getImageUrl(product.gambar)" class="w-full h-full object-cover">
                                        <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                                    </div>
                                    <div>
                                        <div class="font-black text-slate-800 text-xs md:text-sm uppercase">{{ product.nama_produk }}</div>
                                        <div class="flex items-center gap-2 mt-1.5">
                                            <span class="text-[8px] md:text-[9px] font-black bg-slate-100 text-slate-500 px-2 py-0.5 rounded border border-slate-200 uppercase tracking-widest">{{ product.sku || 'NO-SKU' }}</span>
                                            <span v-if="product.satuan_besar" class="text-[8px] md:text-[9px] font-black bg-purple-50 text-purple-600 px-2 py-0.5 rounded border border-purple-100 uppercase tracking-widest">{{ product.satuan_besar }} / {{ product.isi_per_besar }}</span>
                                        </div>
                                    </div>
                                </td>
                                <td class="p-5">
                                    <span class="bg-blue-50 text-blue-600 border border-blue-100 px-3 py-1.5 rounded-lg font-black text-[9px] uppercase tracking-widest">{{ product.kategori || 'General' }}</span>
                                </td>
                                <td class="p-5 text-right">
                                    <div class="text-slate-500 font-black text-xs">Rp {{ (product.harga_modal || 0).toLocaleString('id-ID') }}</div>
                                    <div class="text-[8px] font-bold text-slate-400 uppercase mt-0.5 tracking-widest">/ {{ product.satuan_dasar || 'PCS' }}</div>
                                </td>
                                <td class="p-5 text-right">
                                    <div class="font-black text-slate-800 text-sm">Rp {{ (product.harga_jual || 0).toLocaleString('id-ID') }}</div>
                                    <div class="mt-1 flex justify-end">
                                        <span class="text-[9px] font-black px-2 py-0.5 rounded tracking-widest uppercase"
                                            :class="{
                                                'bg-red-100 text-red-600': ((product.harga_jual - product.harga_modal) / product.harga_jual * 100) < 5,
                                                'bg-amber-100 text-amber-600': ((product.harga_jual - product.harga_modal) / product.harga_jual * 100) >= 5 && ((product.harga_jual - product.harga_modal) / product.harga_jual * 100) <= 15,
                                                'bg-emerald-100 text-emerald-600': ((product.harga_jual - product.harga_modal) / product.harga_jual * 100) > 15
                                            }">
                                            +Rp {{ (product.harga_jual - product.harga_modal).toLocaleString('id-ID') }} 
                                            ({{ Math.round(((product.harga_jual - product.harga_modal) / product.harga_jual) * 100) || 0 }}%)
                                        </span>
                                    </div>
                                </td>
                                <td class="p-5 text-center">
                                    <span class="px-3 py-1.5 text-[10px] rounded-lg font-black tracking-widest shadow-sm inline-flex items-center gap-1" :class="(product.stok || 0) > 10 ? 'bg-emerald-50 text-emerald-600 border border-emerald-100' : 'bg-red-50 text-red-600 border border-red-100'">
                                        {{ product.stok || 0 }} <span class="uppercase text-[8px] opacity-70">{{ product.satuan_dasar || 'PCS' }}</span>
                                    </span>
                                </td>
                                <td class="p-5 text-center">
                                    <div class="flex justify-center gap-2">
                                        <button @click="editProduct(product)" class="p-2.5 bg-slate-50 border border-slate-200 text-slate-400 rounded-xl hover:bg-blue-600 hover:text-white hover:border-blue-600 transition-colors" title="Edit">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2-2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                                        </button>
                                        <button @click="deleteProduct(product.id)" class="p-2.5 bg-slate-50 border border-slate-200 text-slate-400 rounded-xl hover:bg-red-500 hover:text-white hover:border-red-500 transition-colors" title="Hapus">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h14" /></svg>
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div class="p-4 bg-slate-50 border-t border-slate-100 flex flex-col md:flex-row justify-between items-center gap-4 shrink-0">
                    <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest text-center md:text-left">
                        Halaman <span class="text-blue-600">{{ currentPage }}</span> dari {{ totalPages }} 
                        <span v-if="products.length > 0" class="hidden sm:inline">| Menampilkan {{ products.length }} Item</span>
                    </span>
                    
                    <div class="flex flex-wrap justify-center gap-1.5 md:gap-2">
                        <button @click="changePage(1)" :disabled="currentPage === 1 || isLoading" class="w-8 h-8 md:w-9 md:h-9 flex items-center justify-center bg-white border border-slate-200 rounded-lg text-slate-400 hover:text-blue-600 hover:border-blue-300 hover:bg-blue-50 disabled:opacity-40 shadow-sm transition-all" title="Halaman Pertama">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" /></svg>
                        </button>
                        <button @click="changePage(currentPage - 1)" :disabled="currentPage === 1 || isLoading" class="w-8 h-8 md:w-9 md:h-9 flex items-center justify-center bg-white border border-slate-200 rounded-lg text-slate-400 hover:text-blue-600 hover:border-blue-300 hover:bg-blue-50 disabled:opacity-40 shadow-sm transition-all" title="Sebelumnya">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" /></svg>
                        </button>
                        
                        <button v-for="page in visiblePages" :key="page" @click="changePage(page)" :disabled="isLoading" 
                            :class="currentPage === page ? 'bg-blue-600 text-white border-blue-600 shadow-md shadow-blue-200' : 'bg-white text-slate-600 border-slate-200 hover:text-blue-600 hover:border-blue-300 hover:bg-blue-50'"
                            class="w-8 h-8 md:w-9 md:h-9 flex items-center justify-center rounded-lg text-[10px] md:text-xs font-black transition-all shadow-sm">
                            {{ page }}
                        </button>
                        
                        <button @click="changePage(currentPage + 1)" :disabled="currentPage === totalPages || isLoading" class="w-8 h-8 md:w-9 md:h-9 flex items-center justify-center bg-white border border-slate-200 rounded-lg text-slate-400 hover:text-blue-600 hover:border-blue-300 hover:bg-blue-50 disabled:opacity-40 shadow-sm transition-all" title="Selanjutnya">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" /></svg>
                        </button>
                        <button @click="changePage(totalPages)" :disabled="currentPage === totalPages || isLoading" class="w-8 h-8 md:w-9 md:h-9 flex items-center justify-center bg-white border border-slate-200 rounded-lg text-slate-400 hover:text-blue-600 hover:border-blue-300 hover:bg-blue-50 disabled:opacity-40 shadow-sm transition-all" title="Halaman Terakhir">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M13 5l7 7-7 7M5 5l7 7-7 7" /></svg>
                        </button>
                    </div>
                </div>
            </div>

            <div v-if="showFormModal" class="fixed inset-0 bg-slate-900/80 flex items-center justify-center z-[100] p-4 backdrop-blur-sm">
                <div class="bg-white rounded-[32px] shadow-2xl w-full max-w-xl flex flex-col max-h-[90vh] overflow-hidden border border-slate-100">
                    
                    <div class="p-6 border-b border-slate-50 bg-slate-50/50 flex justify-between items-center shrink-0">
                        <h2 class="text-xl font-black text-slate-800 uppercase italic">{{ isEditing ? 'Edit Data Produk' : 'Registrasi Produk Baru' }}</h2>
                        <button @click="showFormModal = false" class="p-2 rounded-xl bg-white text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all border border-slate-100 shadow-sm">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                        </button>
                    </div>

                    <div class="p-6 md:p-8 overflow-y-auto custom-scrollbar">
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
                            
                            <div class="md:col-span-2 flex items-center gap-5 p-4 rounded-[24px] bg-slate-50 border border-slate-100 mb-2">
                                <div @click="fileInput.click()" class="w-16 h-16 rounded-[18px] border-2 border-dashed border-slate-300 flex items-center justify-center bg-white cursor-pointer overflow-hidden shadow-inner shrink-0 hover:border-blue-400 transition-colors">
                                    <img v-if="imagePreview" :src="imagePreview" class="w-full h-full object-cover">
                                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"/><circle cx="12" cy="13" r="3"/></svg>
                                </div>
                                <div class="flex-1">
                                    <input v-model="form.name" type="text" placeholder="NAMA BARANG..." class="w-full px-4 py-3 rounded-xl border border-slate-200 focus:border-blue-600 outline-none font-black text-sm uppercase mb-2 text-slate-800 transition-colors">
                                    <div class="flex gap-2">
                                        <input v-model="form.sku" type="text" placeholder="BARCODE / SKU..." class="w-full px-4 py-2 rounded-xl border border-slate-200 focus:border-blue-600 outline-none font-bold text-xs uppercase text-slate-800 transition-colors">
                                        <button @click.prevent="startScanner" class="px-4 bg-blue-50 text-blue-600 hover:bg-blue-600 hover:text-white rounded-xl border border-blue-100 transition-colors shadow-sm">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                                        </button>
                                    </div>
                                </div>
                                <input type="file" ref="fileInput" @change="onFileChange" accept="image/*" class="hidden">
                            </div>

                            <div class="md:col-span-2">
                                <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5 block">Kategori</label>
                                <input list="kategori-list" v-model="form.category" placeholder="Pilih / Ketik Kategori Baru..." class="w-full px-4 py-3.5 rounded-xl border border-slate-200 focus:border-blue-600 outline-none font-bold text-sm bg-white uppercase transition-all text-slate-800">
                                <datalist id="kategori-list"><option v-for="cat in categories" :key="cat" :value="cat"></option></datalist>
                            </div>

                            <div class="md:col-span-2双 p-5 bg-slate-900 rounded-[28px] text-white shadow-xl mt-2 mb-2">
                                <div class="flex items-center gap-2 mb-4">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-blue-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m7.5 4.27 9 5.15"/><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="M12 22V12"/></svg>
                                    <h4 class="font-black text-[10px] uppercase tracking-[0.2em]">Konversi & Satuan Jual</h4>
                                </div>

                                <div class="grid grid-cols-2 gap-4">
                                    <div>
                                        <label class="text-[8px] font-black text-slate-500 uppercase block mb-1">Satuan Dasar Jual</label>
                                        <select v-model="form.satuan_dasar" class="w-full p-3.5 bg-slate-800 border border-slate-700 rounded-xl outline-none font-black text-xs uppercase cursor-pointer text-white">
                                            <option value="PCS">PCS</option>
                                            <option value="KG">KG</option>
                                            <option value="GRAM">GRAM</option>
                                            <option value="PACK">PACK</option>
                                            <option value="BOX">BOX</option>
                                            <option value="LITER">LITER</option>
                                            <option value="BOTOL">BOTOL</option>
                                        </select>
                                    </div>
                                    
                                    <div>
                                        <label class="text-[8px] font-black text-slate-500 uppercase block mb-1">Beli Dalam Karton/Kemasan Besar?</label>
                                        <div @click="form.has_satuan_besar = !form.has_satuan_besar" class="w-full p-3.5 bg-slate-800 border border-slate-700 rounded-xl font-black text-[10px] uppercase cursor-pointer flex items-center justify-between transition-colors hover:border-blue-500">
                                            {{ form.has_satuan_besar ? 'YA (AKTIF)' : 'TIDAK (HANYA PCS)' }}
                                            <div :class="form.has_satuan_besar ? 'bg-blue-500' : 'bg-slate-600'" class="w-2 h-2 rounded-full shadow-[0_0_8px_rgba(59,130,246,0.5)] transition-colors"></div>
                                        </div>
                                    </div>

                                    <div v-if="form.has_satuan_besar" class="col-span-2 grid grid-cols-3 gap-3 pt-4 border-t border-slate-800 mt-2">
                                        <div>
                                            <label class="text-[8px] font-black text-blue-400 uppercase block mb-1">Sebutannya Apa?</label>
                                            <input v-model="form.satuan_besar" type="text" placeholder="KARTON / BOX" class="w-full p-3 bg-slate-800 border border-blue-900 focus:border-blue-500 rounded-xl outline-none font-black text-xs uppercase text-white transition-all">
                                        </div>
                                        <div>
                                            <label class="text-[8px] font-black text-blue-400 uppercase block mb-1">1 {{ form.satuan_besar || 'KEMASAN' }} Isi Berapa {{ form.satuan_dasar }}?</label>
                                            <input v-model.number="form.isi_per_besar" type="number" class="w-full p-3 bg-slate-800 border border-blue-900 focus:border-blue-500 rounded-xl outline-none font-black text-xs text-white transition-all">
                                        </div>
                                        <div>
                                            <label class="text-[8px] font-black text-amber-400 uppercase block mb-1">Harga Beli 1 {{ form.satuan_besar || 'KEMASAN' }}</label>
                                            <input v-model.number="form.harga_beli_besar" type="number" placeholder="Rp" class="w-full p-3 bg-amber-900/20 border border-amber-900 focus:border-amber-500 rounded-xl outline-none font-black text-xs text-amber-400 transition-all">
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div class="p-5 border rounded-[28px] transition-all duration-300" :class="form.has_satuan_besar ? 'bg-slate-100/80 border-transparent' : 'bg-white border-slate-200'">
                                <label class="text-[9px] font-black uppercase tracking-widest mb-3 block" :class="form.has_satuan_besar ? 'text-indigo-500' : 'text-slate-400'">Harga Modal Dasar ({{ form.satuan_dasar }})</label>
                                <div class="relative">
                                    <span class="absolute inset-y-0 left-0 pl-4 flex items-center text-xs font-black" :class="form.has_satuan_besar ? 'text-indigo-400' : 'text-slate-400'">Rp</span>
                                    <input v-model.number="form.cost_price" type="number" :disabled="form.has_satuan_besar" :class="form.has_satuan_besar ? 'text-indigo-600 font-black bg-slate-200/40 cursor-not-allowed border-transparent shadow-none' : 'text-slate-800 font-black bg-white border-slate-200 focus:border-blue-600 outline-none shadow-inner'" class="w-full pl-12 pr-4 py-4 rounded-2xl text-lg border-2">
                                </div>
                                <div v-if="form.has_satuan_besar" class="mt-2.5 flex items-center justify-center gap-1.5">
                                    <span class="text-[8px] font-black text-indigo-500 uppercase tracking-widest italic">* Terkunci dari kalkulator grosir</span>
                                </div>
                            </div>

                            <div class="p-5 bg-blue-50 border border-blue-100 rounded-[28px] shadow-sm">
                                <label class="text-[9px] font-black text-blue-400 uppercase tracking-widest mb-3 block">Harga Jual Eceran</label>
                                <div class="relative">
                                    <span class="absolute inset-y-0 left-0 pl-4 flex items-center text-xs font-black text-blue-400">Rp</span>
                                    <input v-model.number="form.price" type="number" class="w-full pl-12 pr-4 py-4 rounded-2xl bg-white border border-blue-200 focus:border-blue-600 outline-none font-black text-lg text-blue-700 shadow-inner transition-all">
                                </div>
                                <p class="text-[9px] font-black text-blue-500 mt-2 uppercase tracking-widest text-center">Profit: Rp {{ form.price - form.cost_price }} / {{ form.satuan_dasar }}</p>
                            </div>

                            <div class="md:col-span-2 grid grid-cols-1 gap-4 bg-slate-50 p-5 rounded-[24px] border border-slate-200">
                                <div v-if="form.has_satuan_besar" class="grid grid-cols-1 sm:grid-cols-3 gap-4 items-center">
                                    <div>
                                        <label class="text-[9px] font-black text-indigo-600 uppercase tracking-widest mb-2 block">Jumlah {{ form.satuan_besar || 'KARTON' }}</label>
                                        <input v-model.number="stok_dalam_karton" type="number" min="0" placeholder="0" class="w-full px-4 py-3.5 rounded-xl bg-white border border-slate-200 focus:border-indigo-600 outline-none font-black text-lg text-indigo-600 transition-all shadow-sm">
                                    </div>
                                    <div>
                                        <label class="text-[9px] font-black text-amber-600 uppercase tracking-widest mb-2 block">+ Lebih Eceran ({{ form.satuan_dasar }})</label>
                                        <input v-model.number="eceran_tambahan" type="number" min="0" placeholder="0" class="w-full px-4 py-3.5 rounded-xl bg-white border border-slate-200 focus:border-amber-600 outline-none font-black text-lg text-amber-600 transition-all shadow-sm">
                                    </div>
                                    <div>
                                        <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-2 block">Total Stok Akhir</label>
                                        <input :value="form.stock" type="number" disabled class="w-full px-4 py-3.5 rounded-xl bg-slate-100 border border-transparent font-black text-xl text-slate-500 text-center">
                                    </div>
                                </div>
                                <div v-else>
                                    <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-2 block text-center">Stok Awal Fisik (dalam {{ form.satuan_dasar }})</label>
                                    <input v-model.number="form.stock" type="number" min="0" placeholder="0" class="w-full px-4 py-4 rounded-2xl bg-white border-2 border-slate-200 focus:border-blue-600 outline-none font-black text-center text-2xl text-slate-800 transition-all">
                                </div>
                            </div>

                        </div>
                    </div>

                    <div class="p-6 bg-slate-50 border-t border-slate-100 shrink-0">
                        <button @click="submitProduct" :disabled="isSubmitting" class="w-full py-5 font-black text-xs uppercase tracking-[0.2em] bg-blue-600 text-white rounded-[24px] shadow-xl shadow-blue-200 hover:bg-slate-900 transition-all active:scale-95 flex items-center justify-center gap-3 disabled:opacity-50">
                            <template v-if="isSubmitting">
                                <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                                Menyimpan Data...
                            </template>
                            <template v-else>
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
                                Simpan Perubahan Produk
                            </template>
                        </button>
                    </div>
                </div>
            </div>
            
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

<style scoped>
.custom-scrollbar::-webkit-scrollbar { height: 6px; width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { 
  -webkit-appearance: none; 
  margin: 0; 
}
input[type=number] { -moz-appearance: textfield; }
</style>