import { ref, onMounted, computed } from 'vue';
import api from '../../../api.js';
import Swal from 'sweetalert2';

export function useReturReport() {
    const returns = ref([]);
    const isLoading = ref(true);
    const searchQuery = ref('');

    // State Modal Detail
    const isModalOpen = ref(false);
    const selectedDocument = ref(null);

    const user = ref({
        storeName: localStorage.getItem('storeName') || 'POS UMKM',
        name: localStorage.getItem('name') || 'Admin'
    });

    const fetchReturns = async () => {
        isLoading.value = true;
        try {
            const response = await api.get('/retail/returns');
            returns.value = response.data.data || [];
        } catch (error) {
            Swal.fire('Error', 'Gagal memuat data riwayat retur.', 'error');
        } finally {
            isLoading.value = false;
        }
    };

    onMounted(() => fetchReturns());

    // 🚀 SIHIR GHAIB: Grouping data
    const groupedReturns = computed(() => {
        const groups = {};
        returns.value.forEach(ret => {
            const docNo = ret.return_no || `RET-OLD-${ret.id}`; 
            
            if (!groups[docNo]) {
                groups[docNo] = {
                    return_no: docNo,
                    created_at: ret.created_at,
                    user: ret.user || { name: 'Sistem' },
                    items: [],
                    total_qty: 0
                };
            }
            groups[docNo].items.push(ret);
            groups[docNo].total_qty += ret.qty;
        });
        return Object.values(groups).sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    });

    // 🚀 PENCARIAN
    const filteredDocuments = computed(() => {
        if (!searchQuery.value) return groupedReturns.value;
        const query = searchQuery.value.toLowerCase();
        
        return groupedReturns.value.filter(doc => 
            doc.return_no.toLowerCase().includes(query) ||
            doc.user.name.toLowerCase().includes(query) ||
            doc.items.some(item => item.product?.nama_produk?.toLowerCase().includes(query))
        );
    });

    const openDetail = (doc) => {
        selectedDocument.value = doc;
        isModalOpen.value = true;
    };

    const closeDetail = () => {
        isModalOpen.value = false;
        setTimeout(() => { selectedDocument.value = null; }, 300);
    };

    const printDocument = () => {
        window.print();
    };

    return {
        returns, isLoading, searchQuery, isModalOpen, selectedDocument, user,
        filteredDocuments, openDetail, closeDetail, printDocument
    };
}