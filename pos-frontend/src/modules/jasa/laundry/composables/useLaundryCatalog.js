import { ref, computed } from 'vue';
import { laundryCatalogService } from '../services/laundryCatalogService.js';
import Swal from 'sweetalert2';

export function useLaundryCatalog() {
    const activeTab = ref('jasa');
    const services = ref([]);
    const perfumes = ref([]);
    const isLoading = ref(false);
    const searchQuery = ref('');
    const isEditing = ref(false);
    const editId = ref(null);
    const showForm = ref(false);

    const formJasa = ref({ nama_produk: '', harga_jual: '', satuan_dasar: 'KG', estimasi: '1 Hari' });
    const formParfum = ref({ nama: '', harga: 0 });

    const formatRupiah = (angka) => {
        return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(angka || 0);
    };

    const filteredItems = computed(() => {
        const query = searchQuery.value.trim().toLowerCase();
        if (activeTab.value === 'jasa') {
            if (!query) return services.value;
            return services.value.filter(s => s.nama_produk.toLowerCase().includes(query));
        } else {
            if (!query) return perfumes.value;
            return perfumes.value.filter(p => p.nama.toLowerCase().includes(query));
        }
    });

    const loadAllData = async () => {
        isLoading.value = true;
        try {
            const [srvData, prfData] = await Promise.all([
                laundryCatalogService.getServices(),
                laundryCatalogService.getPerfumes()
            ]);
            services.value = srvData;
            perfumes.value = prfData;
        } catch (err) {
            console.error('Katalog Delay:', err);
        } finally {
            isLoading.value = false;
        }
    };

    const switchTab = (tab) => {
        activeTab.value = tab;
        cancelForm();
    };

    const triggerEdit = (item) => {
        isEditing.value = true;
        editId.value = item.id;
        formJasa.value = {
            nama_produk: item.nama_produk,
            harga_jual: item.harga_jual,
            satuan_dasar: item.satuan_dasar,
            estimasi: item.estimasi || '1 Hari'
        };
        showForm.value = true;
        window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    const cancelForm = () => {
        showForm.value = false;
        setTimeout(() => {
            isEditing.value = false;
            editId.value = null;
            formJasa.value = { nama_produk: '', harga_jual: '', satuan_dasar: 'KG', estimasi: '1 Hari' };
            formParfum.value = { nama: '', harga: 0 };
        }, 200);
    };

    const handleSave = async () => {
        if (activeTab.value === 'jasa') {
            if (!formJasa.value.nama_produk || !formJasa.value.harga_jual) {
                return Swal.fire('Data Kurang', 'Nama dan Harga paket cuci wajib diisi bray!', 'warning');
            }
            try {
                const payload = {
                    nama_produk: formJasa.value.nama_produk,
                    harga_jual: parseFloat(formJasa.value.harga_jual),
                    satuan_dasar: formJasa.value.satuan_dasar,
                    estimasi: formJasa.value.estimasi
                };
                if (isEditing.value) {
                    await laundryCatalogService.updateService(editId.value, payload);
                    Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Paket Jasa Berhasil Diubah!', showConfirmButton: false, timer: 1500 });
                } else {
                    await laundryCatalogService.createService(payload);
                    Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Paket Jasa Berhasil Disimpan!', showConfirmButton: false, timer: 1500 });
                }
                cancelForm();
                services.value = await laundryCatalogService.getServices();
            } catch (e) { Swal.fire('Gagal!', 'Gagal memproses konfigurasi paket.', 'error'); }
        } else {
            if (!formParfum.value.nama) {
                return Swal.fire('Data Kurang', 'Nama varian aroma parfum wajib diisi bray!', 'warning');
            }
            try {
                const payload = {
                    nama: formParfum.value.nama,
                    harga: parseFloat(formParfum.value.harga || 0)
                };
                await laundryCatalogService.createPerfume(payload);
                Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Aroma Parfum Ditambahkan!', showConfirmButton: false, timer: 1500 });
                cancelForm();
                perfumes.value = await laundryCatalogService.getPerfumes();
            } catch (e) { Swal.fire('Gagal!', 'Gagal meregistrasikan aroma.', 'error'); }
        }
    };

    const handleConfirmDelete = (id, nama) => {
        Swal.fire({
            title: `Hapus ${nama}?`,
            text: 'Tindakan ini menghapus data dari katalog selamanya!',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#e11d48',
            cancelButtonColor: '#94a3b8',
            confirmButtonText: 'Ya, Eksekusi Hapus'
        }).then(async (result) => {
            if (result.isConfirmed) {
                try {
                    if (activeTab.value === 'jasa') {
                        await laundryCatalogService.deleteService(id);
                        services.value = await laundryCatalogService.getServices();
                    } else {
                        await laundryCatalogService.deletePerfume(id);
                        perfumes.value = await laundryCatalogService.getPerfumes();
                    }
                    Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: 'Item Dihapus!', showConfirmButton: false, timer: 1500 });
                } catch (error) {
                    Swal.fire('Gagal!', 'Gagal menghapus data dari cluster.', 'error');
                }
            }
        });
    };

    return {
        activeTab, isLoading, searchQuery, isEditing, showForm,
        formJasa, formParfum, filteredItems, formatRupiah,
        loadAllData, switchTab, triggerEdit, cancelForm, handleSave, handleConfirmDelete
    };
}