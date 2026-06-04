import { ref, computed, onMounted, nextTick, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { posService } from '../services/posService.js';
import api from '../../../api.js';
import Swal from 'sweetalert2';
import { Html5Qrcode } from 'html5-qrcode';

export function usePos() {
  const router = useRouter();

  // --- SETUP USER & ROLE ---
  const getUserInfo = () => {
    const token = localStorage.getItem('token'), role = localStorage.getItem('role') || 'kasir';
    let name = localStorage.getItem('name');
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        if (!name || name === 'undefined' || name === '') name = payload.name || payload.username || 'Kasir Toko';
        return { userId: payload.user_id, role, name };
      } catch (e) { return { userId: 0, role, name: 'Kasir Toko' }; }
    }
    return { userId: 0, role, name: 'Kasir Toko' };
  };

  const currentUser = ref(getUserInfo()); const currentSession = ref(null); const currentTime = ref(''); let timer;
  const products = ref([]), isLoadingProducts = ref(true), cart = ref([]), heldOrders = ref([]), showHeldModal = ref(false);
  const payAmount = ref(0), paymentMethod = ref('Cash'), showReceipt = ref(false), showQrisModal = ref(false);
  const lastTransaction = ref(null), showReceiptClosing = ref(false), lastClosingData = ref(null), noHpPelanggan = ref('');
  const isMobileCartOpen = ref(false), searchQuery = ref(''), searchInput = ref(null), showScanner = ref(false);
  let html5QrCode = null;

  const storeSettings = ref({ payment_type: 'qris_static', qris_image: '', qris_name: '', is_tax_active: false, pajak_persen: 0 });

  // --- 🚀 UTILS: PAGAR PENGAMAN URL GAMBAR SUPABASE CLOUD ---
  const getImageUrl = (path) => {
    if (!path) return null;
    // Kalau path gambar sudah link internet utuh (https://), langsung balikin mentah bray!
    if (path.startsWith('http://') || path.startsWith('https://')) return path;
    return `${import.meta.env.VITE_API_BASE_URL}${path}`;
  };

  // --- LOGIKA KAMERA SCANNER KASIR ---
  const startScanner = async () => {
    showScanner.value = true;
    setTimeout(async () => {
      try {
        html5QrCode = new Html5Qrcode('reader-kasir');
        await html5QrCode.start(
          { facingMode: 'environment' }, { fps: 15, qrbox: { width: 250, height: 100 } },
          (decodedText) => {
            searchQuery.value = decodedText; stopScanner(); handleBarcodeScan();
            const audio = new Audio('https://www.soundjay.com/buttons/beep-07a.mp3'); audio.play().catch(() => {});
          }, () => {}
        );
      } catch (err) {
        console.error(err); Swal.fire('Oops!', 'Gagal menyalakan kamera. Pastikan izin kamera aktif.', 'error'); stopScanner();
      }
    }, 300);
  };

  const stopScanner = () => {
    if (html5QrCode) {
      html5QrCode.stop().then(() => { html5QrCode.clear(); showScanner.value = false; }).catch(() => { showScanner.value = false; });
    } else { showScanner.value = false; }
  };

  // --- LOGIKA DATA PRODUK & PENCARIAN ---
  const fetchProducts = async () => {
    try {
      const response = await posService.getProducts();
      products.value = response.data.map((p) => ({
        id: p.id, sku: p.sku || `SKU-${p.id}`, name: p.nama_produk, price: p.harga_jual, stock: p.stok, image: p.gambar,
        satuan_dasar: p.satuan_dasar || 'PCS', satuan_besar: p.satuan_besar || null, isi_per_besar: p.isi_per_besar || 0, harga_jual_besar: p.harga_jual_besar || 0,
        is_nested_uom: p.is_nested_uom || false, satuan_tengah: p.satuan_tengah || null, isi_tengah_ke_dasar: p.isi_tengah_ke_dasar || 0,
        harga_jual_tengah: p.harga_jual_tengah || p.harga_jual * (p.isi_tengah_ke_dasar || 1),
      }));
    } catch (error) { console.error('Gagal narik produk:', error); } 
    finally { isLoadingProducts.value = false; }
  };

  const filteredProducts = computed(() => {
    if (!searchQuery.value) return products.value;
    const query = searchQuery.value.toLocaleLowerCase();
    return products.value.filter((p) => p.name.toLowerCase().includes(query) || (p.sku && p.sku.toLowerCase().includes(query)));
  });

  const handleBarcodeScan = () => {
    if (!searchQuery.value) return;
    const query = String(searchQuery.value).trim().toLowerCase();
    const exactMatch = products.value.find((p) => p.sku && String(p.sku).toLowerCase() === query);

    if (exactMatch) { addToCart(exactMatch); searchQuery.value = ''; } 
    else if (filteredProducts.value.length === 1) { addToCart(filteredProducts.value[0]); searchQuery.value = ''; }
    nextTick(() => { if (searchInput.value) searchInput.value.focus(); });
  };

  // --- LOGIKA KERANJANG BELANJA ---
  const addToCart = (product) => {
    if (product.stock <= 0) return Swal.fire({ icon: 'error', title: 'Stok Habis!', text: `Stok ${product.name} kosong` });

    const defaultUom = product.satuan_dasar;
    const existingItem = cart.value.find((item) => item.id === product.id && item.selected_uom === defaultUom);

    if (existingItem) {
      if ((existingItem.qty + 1) * existingItem.uom_multiplier <= product.stock) existingItem.qty++;
      else return Swal.fire({ icon: 'warning', title: 'Stok Terbatas', text: 'Kuantitas melebihi stok!' });
    } else {
      cart.value.unshift({
        id: product.id, name: product.name, price: product.price, qty: 1, selected_uom: defaultUom, uom_multiplier: 1, has_grosir: !!product.satuan_besar,
        satuan_dasar: product.satuan_dasar, harga_dasar: product.price, is_nested: product.is_nested_uom, satuan_tengah: product.satuan_tengah,
        isi_tengah: product.isi_tengah_ke_dasar, harga_tengah: product.harga_jual_tengah, satuan_besar: product.satuan_besar, isi_besar: product.isi_per_besar, harga_besar: product.harga_jual_besar,
      });
    }

    if (window.innerWidth < 1024 && !isMobileCartOpen.value) {
      Swal.fire({ toast: true, position: 'top', icon: 'success', title: `${product.name} Masuk Keranjang`, showConfirmButton: false, timer: 800, timerProgressBar: true });
    }
  };

  const toggleUom = (item) => {
    if (!item.has_grosir) return;
    const prodMaster = products.value.find((p) => p.id === item.id);

    if (item.is_nested) {
      if (item.selected_uom === item.satuan_dasar) {
        if (item.qty * item.isi_tengah > prodMaster.stock) return Swal.fire({ icon: 'error', title: 'Stok Kurang', text: `Stok tak cukup untuk ${item.satuan_tengah}` });
        item.selected_uom = item.satuan_tengah; item.uom_multiplier = item.isi_tengah; item.price = item.harga_tengah;
      } else if (item.selected_uom === item.satuan_tengah) {
        if (item.qty * item.isi_besar > prodMaster.stock) {
          Swal.fire({ icon: 'error', title: 'Stok Kurang', text: `Stok tak cukup untuk ${item.satuan_besar}` });
          item.selected_uom = item.satuan_dasar; item.uom_multiplier = 1; item.price = item.harga_dasar; return;
        }
        item.selected_uom = item.satuan_besar; item.uom_multiplier = item.isi_besar; item.price = item.harga_besar;
      } else {
        item.selected_uom = item.satuan_dasar; item.uom_multiplier = 1; item.price = item.harga_dasar;
      }
    } else {
      if (item.selected_uom === item.satuan_dasar) {
        if (item.qty * item.isi_besar > prodMaster.stock) return Swal.fire({ icon: 'error', title: 'Stok Kurang', text: `Stok tak cukup untuk ${item.satuan_besar}` });
        item.selected_uom = item.satuan_besar; item.uom_multiplier = item.isi_besar; item.price = item.harga_besar;
      } else {
        item.selected_uom = item.satuan_dasar; item.uom_multiplier = 1; item.price = item.harga_dasar;
      }
    }
  };

  const decreaseQty = (item) => {
    const existingItem = cart.value.find((i) => i.id === item.id && i.selected_uom === item.selected_uom);
    if (existingItem) {
      if (existingItem.qty > 1) existingItem.qty--;
      else {
        cart.value = cart.value.filter((i) => !(i.id === item.id && i.selected_uom === item.selected_uom));
        if (cart.value.length === 0) isMobileCartOpen.value = false;
      }
    }
  };

  const increaseQty = (item) => {
    const existingItem = cart.value.find((i) => i.id === item.id && i.selected_uom === item.selected_uom), prodMaster = products.value.find((p) => p.id === item.id);
    if (existingItem && prodMaster) {
      if ((existingItem.qty + 1) * existingItem.uom_multiplier <= prodMaster.stock) existingItem.qty++;
      else Swal.fire({ icon: 'warning', title: 'Stok Terbatas', text: 'Kuantitas melebihi stok!' });
    }
  };

  const validateQty = (item) => {
    const existingItem = cart.value.find((i) => i.id === item.id && i.selected_uom === item.selected_uom);
    if (existingItem && existingItem.qty < 1) existingItem.qty = 1;
    const prodMaster = products.value.find((p) => p.id === item.id);
    if (existingItem && prodMaster && existingItem.qty * existingItem.uom_multiplier > prodMaster.stock) {
      Swal.fire({ icon: 'warning', title: 'Overstock!', text: 'Kuantitas dikembalikan ke batas stok' });
      existingItem.qty = Math.floor(prodMaster.stock / existingItem.uom_multiplier);
    }
  };

  const clearCart = () => {
    if (cart.value.length === 0) return;
    Swal.fire({ title: 'Batalkan Transaksi?', text: 'Semua barang di keranjang akan dihapus!', icon: 'warning', showCancelButton: true, confirmButtonColor: '#ef4444', confirmButtonText: 'Ya, Hapus Semua!', cancelButtonText: 'Batal' }).then((result) => {
      if (result.isConfirmed) { cart.value = []; payAmount.value = 0; setPaymentMethod('Cash'); isMobileCartOpen.value = false; }
    });
  };

  const holdTransaction = () => {
    if (cart.value.length === 0) return;
    heldOrders.value.push({ id: Date.now(), customer: `Pelanggan ${heldOrders.value.length + 1}`, items: [...cart.value], time: new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }), total: totalBelanja.value });
    cart.value = []; payAmount.value = 0; setPaymentMethod('Cash'); isMobileCartOpen.value = false;
    Swal.fire({ toast: true, position: 'top-end', icon: 'info', title: 'Pesanan ditunda!', showConfirmButton: false, timer: 1500 });
  };

  const resumeOrder = (order) => {
    if (cart.value.length > 0) {
      Swal.fire({ title: 'Timpa Keranjang?', text: 'Ada barang di keranjang saat ini. Lanjutkan memuat pesanan tertunda?', icon: 'warning', showCancelButton: true, confirmButtonText: 'Ya, Timpa!' }).then((res) => { if (res.isConfirmed) processResume(order); });
    } else { processResume(order); }
  };

  const processResume = (order) => { cart.value = [...order.items]; heldOrders.value = heldOrders.value.filter((h) => h.id !== order.id); showHeldModal.value = false; if (window.innerWidth < 1024) isMobileCartOpen.value = true; };

  // --- KALKULASI UANG & PAJAK ---
  const subTotalBelanja = computed(() => cart.value.reduce((total, item) => total + item.price * item.qty, 0));
  const nilaiPajak = computed(() => (!storeSettings.value || !storeSettings.value.is_tax_active) ? 0 : (storeSettings.value.pajak_persen / 100) * subTotalBelanja.value);
  const totalBelanja = computed(() => Math.round((subTotalBelanja.value + nilaiPajak.value) / 100) * 100);
  const kembalian = computed(() => payAmount.value - totalBelanja.value);

  const setPaymentMethod = (method) => { paymentMethod.value = method; payAmount.value = method !== 'Cash' ? totalBelanja.value : 0; };
  const setNominal = (amount) => { payAmount.value = amount; };

  const formatInputRupiah = (event) => {
    let rawValue = event.target.value.replace(/\D/g, ''); payAmount.value = rawValue ? parseInt(rawValue, 10) : 0;
    event.target.value = payAmount.value === 0 ? '' : payAmount.value.toLocaleString('id-ID');
  };

  // --- 🚀 PROSES CHECKOUT FINAL GOLANG DB SAVE ---
  const isProcessingCheckout = ref(false);

  const executeCheckout = async () => {
    if (isProcessingCheckout.value) return; isProcessingCheckout.value = true;
    const payloadItems = cart.value.map((item) => ({ product_id: item.id, kuantitas: item.qty * item.uom_multiplier, uom_label: `${item.qty} ${item.selected_uom}`, harga_uom: item.price }));

    try {
      const response = await posService.checkout({ session_id: currentSession.value.id, items: payloadItems, nominal_bayar: Number(payAmount.value), metode_bayar: paymentMethod.value, no_hp_pelanggan: noHpPelanggan.value ? String(noHpPelanggan.value) : '' });
      lastTransaction.value = {
        invoice: response.invoice, cart: [...cart.value], total: response.tagihan, pay: payAmount.value, return: response.counter || response.kembali, method: paymentMethod.value, subtotal: subTotalBelanja.value, pajak: nilaiPajak.value,
        date: new Date().toLocaleString('id-ID', { year: '2-digit', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }).replace(/\//g, '.'),
      };
      showQrisModal.value = false; isMobileCartOpen.value = false; showReceipt.value = true; cart.value = []; payAmount.value = 0; setPaymentMethod('Cash'); fetchProducts();
      nextTick(() => { if (searchInput.value) searchInput.value.focus(); });
    } catch (error) { Swal.fire('Gagal!', error.response?.data?.error || 'Koneksi terputus', 'error'); } 
    finally { isProcessingCheckout.value = false; }
  };

  const processCheckout = async () => {
    if (paymentMethod.value === 'Cash' && payAmount.value < totalBelanja.value) return Swal.fire({ icon: 'error', title: 'Uang Kurang!', text: `Kurang Rp ${(totalBelanja.value - payAmount.value).toLocaleString('id-ID')}` });

    if (paymentMethod.value === 'QRIS') {
      const tipePayment = storeSettings.value?.payment_type;
      if (tipePayment === 'midtrans') {
        if (typeof window.snap === 'undefined') return Swal.fire({ icon: 'error', title: 'Script Midtrans Gagal!', text: `Client Key kosong atau script belum ke-load. Pastikan Client Key terisi di Pengaturan Toko!` });
        isProcessingCheckout.value = true;

        try {
          const payRes = await api.post('/retail/pos/midtrans-order', { total: totalBelanja.value });
          if (!payRes.data || !payRes.data.token) { Swal.fire('Error Backend', 'Golang tidak mengembalikan Token Snap!', 'error'); isProcessingCheckout.value = false; return; }

          window.snap.pay(payRes.data.token, {
            onSuccess: () => { Swal.fire('Berhasil', 'Pembayaran QRIS Diterima!', 'success'); executeCheckout(); },
            onPending: () => { Swal.fire('Menunggu', 'Pelanggan belum bayar.', 'info'); isProcessingCheckout.value = false; },
            onError: () => { Swal.fire('Gagal', 'Pembayaran ditolak bank.', 'error'); isProcessingCheckout.value = false; },
            onClose: () => { isProcessingCheckout.value = false; },
          });
        } catch (error) {
          Swal.fire({ icon: 'error', title: 'API Golang Error!', text: error.response?.data?.error || 'Endpoint /retail/pos/midtrans-order gagal dipanggil.' });
          isProcessingCheckout.value = false;
        }
      } else {
        Swal.fire({ icon: 'warning', title: 'Bukan Midtrans!', text: `Tipe di Database saat ini adalah: "${tipePayment}". Makanya lari ke QRIS Statis!` });
        showQrisModal.value = true;
      }
    } else { executeCheckout(); }
  };

  const openSession = async (stationNumber, modalAwalValue) => {
    try { return await posService.openSession({ station_number: stationNumber, modal_awal: parseFloat(modalAwalValue) }); } 
    catch (error) { throw error; }
  };

  const showClosingModal = ref(false);
  const pecahan = ref({ p100k: 0, p50k: 0, p20k: 0, p10k: 0, p5k: 0, p2k: 0, p1k: 0, p500: 0, p200: 0, p100: 0, p50: 0, p25: 0 });
  const totalUangFisik = computed(() => {
    return (pecahan.value.p100k * 100000 + pecahan.value.p50k * 50000 + pecahan.value.p20k * 20000 + pecahan.value.p10k * 10000 + pecahan.value.p5k * 5000 + pecahan.value.p2k * 2000 + pecahan.value.p1k * 1000 + pecahan.value.p500 * 500 + pecahan.value.p200 * 200 + pecahan.value.p100 * 100 + pecahan.value.p50 * 50 + pecahan.value.p25 * 25);
  });

  const handleClosing = async () => {
    try {
      const res = await posService.closeSession(currentSession.value.id, { total_aktual: totalUangFisik.value, pecahan: pecahan.value });
      Swal.fire('Closing Berhasil!', 'Struk closing akan dicetak.', 'success'); lastClosingData.value = res; showClosingModal.value = false; showReceiptClosing.value = true;
    } catch (error) { Swal.fire('Gagal Closing', error.response?.data?.error, 'error'); }
  };

  const logout = () => {
    Swal.fire({ title: 'Akhiri Shift?', text: 'Hitung uang laci (Cash Count) sebelum tutup shift.', icon: 'question', showCancelButton: true, confirmButtonColor: '#2563eb', confirmButtonText: 'Ya, Tutup Shift' }).then((result) => {
      if (result.isConfirmed) { Object.keys(pecahan.value).forEach((k) => (pecahan.value[k] = 0)); showClosingModal.value = true; }
    });
  };

  // --- LIFECYCLE HOOKS ---
  onMounted(async () => {
    const token = localStorage.getItem('token'); if (!token) return router.push('/login');

    try {
      const setRes = await api.get('/retail/store/settings'); storeSettings.value = setRes.data.data;

      if (storeSettings.value.payment_type === 'midtrans' && storeSettings.value.midtrans_client_key) {
        if (!document.getElementById('midtrans-script')) {
          const script = document.createElement('script'); script.id = 'midtrans-script'; script.src = 'https://app.sandbox.midtrans.com/snap/snap.js';
          script.setAttribute('data-client-key', storeSettings.value.midtrans_client_key); document.head.appendChild(script);
        }
      }

      const res = await posService.checkSession(token);
      if (!res.has_session) return router.push('/retail/pos/buka-kasir');

      currentSession.value = res.session; await fetchProducts();
      if (searchInput.value) searchInput.value.focus();
    } catch (error) { if (error.response?.status === 401) router.push('/login'); }

    timer = setInterval(() => {
      currentTime.value = new Date().toLocaleString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit', second: '2-digit' }).replace(/\//g, '.');
    }, 1000);
  });

  onUnmounted(() => { if (showScanner.value) stopScanner(); clearInterval(timer); });

  return {
    currentUser, currentSession, currentTime, products, isLoadingProducts, cart, heldOrders, showHeldModal, payAmount, paymentMethod, showReceipt, showQrisModal, lastTransaction, showReceiptClosing, lastClosingData, isMobileCartOpen, searchQuery, searchInput, showScanner, pecahan, totalUangFisik, filteredProducts, subTotalBelanja, nilaiPajak, totalBelanja, kembalian, isProcessingCheckout, showClosingModal, noHpPelanggan, storeSettings,
    getImageUrl, startScanner, stopScanner, handleBarcodeScan, addToCart, decreaseQty, increaseQty, validateQty, clearCart, holdTransaction, resumeOrder, setPaymentMethod, executeCheckout, formatInputRupiah, processCheckout, handleClosing, logout, toggleUom, setNominal,
  };
}