import { ref, computed } from 'vue';
import { fnbService } from '../services/fnbServices';
import Swal from 'sweetalert2';

export function useOrderFnB() {
    const cart = ref([]);
    const tipeOrder = ref('DINE_IN');
    const nomorMeja = ref('');
    const namaPemesan = ref('');
    const metodeBayar = ref('CASH');
    const uangBayarRaw = ref(0);
    const uangBayarDisplay = ref('');

    const totalBelanja = computed(() => cart.value.reduce((sum, item) => sum + (item.harga * item.qty), 0));
    const kembalian = computed(() => Math.max(0, uangBayarRaw.value - totalBelanja.value));

    const handleUangInput = (e) => {
        let val = e.target.value.replace(/\D/g, '');
        uangBayarRaw.value = Number(val);
        uangBayarDisplay.value = val ? new Intl.NumberFormat('id-ID').format(val) : '';
    };

    const addToCart = (product) => {
        const existingItem = cart.value.find(item => item.id === product.id);
        if (existingItem) existingItem.qty++;
        else cart.value.push({ ...product, qty: 1, notes: '' });
    };

    const decreaseQty = (product) => {
        const existingItem = cart.value.find(item => item.id === product.id);
        if (existingItem) {
            if (existingItem.qty > 1) existingItem.qty--;
            else cart.value = cart.value.filter(item => item.id !== product.id);
        }
    };

    const getQtyInCart = (productId) => {
        const item = cart.value.find(i => i.id === productId);
        return item ? item.qty : 0;
    };

    const resetCart = () => {
        cart.value = [];
        nomorMeja.value = '';
        namaPemesan.value = '';
        uangBayarRaw.value = 0;
        uangBayarDisplay.value = '';
    };

    return { 
        cart, tipeOrder, nomorMeja, namaPemesan, metodeBayar, uangBayarRaw, uangBayarDisplay,
        totalBelanja, kembalian, handleUangInput, addToCart, decreaseQty, getQtyInCart, resetCart
    };
}