import { ref } from 'vue';
import { fnbService } from '../services/fnbServices'; // Pastikan path ini benar ke service kamu
import Swal from 'sweetalert2';

export function useMenuFnB() {
    const products = ref([]);
    const isLoading = ref(false);
    const isSubmitting = ref(false);

    // READ
    const fetchProducts = async () => {
        isLoading.value = true;
        try {
            const data = await fnbService.getProducts();
            products.value = data.map(p => ({
                ...p,
                nama: p.nama_produk,
                harga: p.harga_jual,
                kategori: p.kategori,
                gambar: p.gambar
            }));
        } catch (error) {
            Swal.fire('Gagal!', 'Tidak dapat memuat data menu.', 'error');
        } finally {
            isLoading.value = false;
        }
    };

    // CREATE & UPDATE
    const saveMenu = async (form) => {
        isSubmitting.value = true;
        try {
            const payload = {
                nama: form.nama,
                harga: Number(form.harga),
                kategori: form.kategori,
                gambar: form.gambar
            };

            if (form.id) {
                await fnbService.updateProduct(form.id, payload);
                Swal.fire('Berhasil!', 'Menu diperbarui.', 'success');
            } else {
                await fnbService.createProduct(payload);
                Swal.fire('Berhasil!', 'Menu ditambahkan.', 'success');
            }
            await fetchProducts();
            return true;
        } catch (error) {
            Swal.fire('Error!', 'Gagal menyimpan data.', 'error');
            return false;
        } finally {
            isSubmitting.value = false;
        }
    };

    // DELETE
    const deleteMenu = async (id) => {
        const confirm = await Swal.fire({
            title: 'Yakin hapus?',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#ef4444',
            confirmButtonText: 'Ya, Hapus'
        });

        if (confirm.isConfirmed) {
            try {
                await fnbService.deleteProduct(id);
                Swal.fire('Terhapus!', 'Menu telah dihapus.', 'success');
                await fetchProducts();
            } catch (error) {
                Swal.fire('Error!', 'Gagal menghapus.', 'error');
            }
        }
    };

    // TOGGLE
    const toggleStatus = async (prod) => {
        try {
            const res = await fnbService.toggleStatus(prod.id);
            Swal.fire({ toast: true, position: 'top-end', icon: 'success', title: res.message, showConfirmButton: false, timer: 1500 });
            await fetchProducts();
        } catch (e) {
            Swal.fire('Gagal!', 'Tidak dapat ubah status.', 'error');
        }
    };

    return { products, isLoading, isSubmitting, fetchProducts, saveMenu, deleteMenu, toggleStatus };
}