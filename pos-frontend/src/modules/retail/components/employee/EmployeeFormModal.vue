<script setup>
import { ref, defineProps, defineEmits, watch } from 'vue';
import Swal from 'sweetalert2';

const props = defineProps({
    show: Boolean,
    isEditMode: Boolean,
    isProcessing: Boolean,
    form: Object,
    fotoProfilPreview: String,
    fotoBiometricPreview: String,
    isCameraOpen: Boolean
});

const emit = defineEmits(['close', 'submit', 'update-profile', 'update-biometric', 'update:isCameraOpen']);

const videoRef = ref(null);
const canvasRef = ref(null);
const profileInput = ref(null);

const onProfileChange = (e) => {
    const file = e.target.files[0];
    if (file) emit('update-profile', file, URL.createObjectURL(file));
};

// --- LOGIKA KAMERA LOKAL MODAL ---
const startCamera = async () => {
    emit('update:isCameraOpen', true);
    try {
        const stream = await navigator.mediaDevices.getUserMedia({ video: true });
        videoRef.value.srcObject = stream;
        videoRef.value.play();
    } catch (err) {
        Swal.fire('Gagal!', 'Akses kamera ditolak atau tidak ditemukan.', 'error');
        emit('update:isCameraOpen', false);
    }
};

const capturePhoto = () => {
    const context = canvasRef.value.getContext('2d');
    canvasRef.value.width = videoRef.value.videoWidth;
    canvasRef.value.height = videoRef.value.videoHeight;
    context.drawImage(videoRef.value, 0, 0, canvasRef.value.width, canvasRef.value.height);
    
    canvasRef.value.toBlob((blob) => {
        const file = new File([blob], "face_registration.jpg", { type: "image/jpeg" });
        emit('update-biometric', file, URL.createObjectURL(file));
        stopCamera();
    }, 'image/jpeg');
};

const stopCamera = () => {
    if (videoRef.value && videoRef.value.srcObject) {
        videoRef.value.srcObject.getTracks().forEach(track => track.stop());
    }
    emit('update:isCameraOpen', false);
};

// Matikan kamera otomatis jika modal ditutup
watch(() => props.show, (newVal) => {
    if (!newVal) stopCamera();
});
</script>

<template>
    <div v-if="show" class="fixed inset-0 bg-slate-900/80 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
        <div class="bg-white rounded-[32px] md:rounded-[40px] shadow-2xl w-full max-w-4xl flex flex-col max-h-[95vh] overflow-hidden border-[8px] md:border-[10px] border-slate-900/5">
            
            <div class="p-6 md:p-8 border-b border-slate-50 flex justify-between items-center bg-slate-50/80 shrink-0">
                <h3 class="font-black text-xl md:text-2xl text-slate-800 tracking-tighter uppercase italic">{{ isEditMode ? 'Update Profile' : 'New Registration' }}</h3>
                <button @click="emit('close')" class="text-slate-400 hover:text-red-500 transition-colors bg-white p-2 md:p-2.5 rounded-xl shadow-sm border border-slate-200">✕</button>
            </div>
            
            <div class="overflow-y-auto custom-scrollbar p-6 md:p-8">
                <form @submit.prevent="emit('submit')" class="flex flex-col gap-6 md:gap-8">
                    
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-5 md:gap-6">
                        <div>
                            <label class="block text-[9px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 md:mb-2">Nama Lengkap</label>
                            <input v-model="form.name" type="text" required class="w-full p-3.5 md:p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm text-slate-800">
                        </div>
                        <div>
                            <label class="block text-[9px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 md:mb-2">Jabatan / Role</label>
                            <select v-model="form.role" required class="w-full p-3.5 md:p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-black text-xs text-slate-700 uppercase tracking-widest">
                                <option value="manager">Manager</option>
                                <option value="supervisor">Supervisor</option>
                                <option value="staff">Staff / Cashier</option>
                            </select>
                        </div>
                        <div class="grid grid-cols-2 gap-3 md:gap-4">
                            <div>
                                <label class="block text-[9px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 md:mb-2">Tempat Lahir</label>
                                <input v-model="form.tempat_lahir" type="text" required class="w-full p-3.5 md:p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm text-slate-800">
                            </div>
                            <div>
                                <label class="block text-[9px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 md:mb-2">Tgl Lahir</label>
                                <input v-model="form.tanggal_lahir" type="date" required class="w-full p-3.5 md:p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-xs uppercase tracking-widest text-slate-800">
                            </div>
                        </div>
                        <div>
                            <label class="block text-[9px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 md:mb-2">No. Handphone (WA)</label>
                            <input v-model="form.no_hp" type="text" required class="w-full p-3.5 md:p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm text-slate-800">
                        </div>
                        <div class="md:col-span-2">
                            <label class="block text-[9px] md:text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1.5 md:mb-2">Password <span v-if="isEditMode" class="text-amber-500 italic lowercase tracking-normal">(Kosongi jika tidak diubah)</span></label>
                            <input v-model="form.password" type="password" :required="!isEditMode" placeholder="••••••••" class="w-full p-3.5 md:p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:border-blue-600 outline-none font-bold text-sm text-slate-800">
                        </div>
                    </div>

                    <div class="grid grid-cols-1 md:grid-cols-2 gap-5 md:gap-6 pt-6 md:pt-8 border-t border-slate-100">
                        
                        <div class="flex flex-col items-center p-6 bg-slate-50 rounded-[32px] border-2 border-slate-100">
                            <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest mb-4">1. Foto Profil (Aesthetic)</label>
                            <div @click="$refs.profileInput.click()" class="w-36 h-36 md:w-40 md:h-40 bg-white rounded-[24px] shadow-sm flex items-center justify-center cursor-pointer overflow-hidden border-2 border-slate-200 group relative">
                                <img v-if="fotoProfilPreview" :src="fotoProfilPreview" class="w-full h-full object-cover">
                                <div v-else class="text-slate-300">📷</div>
                            </div>
                            <input type="file" ref="profileInput" @change="onProfileChange" class="hidden" accept="image/*">
                        </div>

                        <div class="flex flex-col items-center p-6 bg-slate-900 rounded-[32px] relative overflow-hidden shadow-2xl border-4 border-slate-800">
                            <label class="text-[10px] font-black text-indigo-400 uppercase tracking-widest mb-4">2. Foto Biometrik (Face AI)</label>
                            
                            <div class="w-full flex flex-col items-center gap-4 relative z-10">
                                <div class="w-36 h-36 md:w-40 md:h-40 bg-slate-950 rounded-[24px] overflow-hidden relative border-2 border-slate-700 shadow-inner">
                                    <video ref="videoRef" v-show="isCameraOpen" autoplay muted playsinline class="w-full h-full object-cover scale-x-[-1]"></video>
                                    <img v-if="fotoBiometricPreview && !isCameraOpen" :src="fotoBiometricPreview" class="w-full h-full object-cover">
                                    <canvas ref="canvasRef" class="hidden"></canvas>
                                    
                                    <div v-if="isCameraOpen" class="absolute inset-0 border-2 border-indigo-500/50 rounded-[24px] pointer-events-none">
                                        <div class="w-full h-0.5 bg-indigo-500 absolute top-0 animate-[scan_2s_infinite] shadow-[0_0_8px_#6366f1]"></div>
                                    </div>
                                </div>

                                <div class="flex gap-2 w-full justify-center">
                                    <button v-if="!isCameraOpen" @click.prevent="startCamera" class="bg-indigo-600 hover:bg-indigo-500 text-white px-6 py-3 md:py-3.5 rounded-xl font-black text-[9px] md:text-[10px] uppercase tracking-widest w-3/4">
                                        Nyalakan Kamera
                                    </button>
                                    <button v-else @click.prevent="capturePhoto" class="bg-emerald-500 hover:bg-emerald-400 text-white px-6 py-3 md:py-3.5 rounded-xl font-black text-[9px] md:text-[10px] uppercase tracking-widest animate-pulse w-3/4">
                                        Jepret Wajah!
                                    </button>
                                </div>
                            </div>
                        </div>

                    </div>

                    <div class="pt-6 border-t border-slate-100 mt-2">
                        <button type="submit" :disabled="isProcessing" class="w-full bg-blue-600 hover:bg-slate-900 text-white py-4 md:py-5 rounded-[20px] md:rounded-[24px] font-black text-xs md:text-sm shadow-xl active:scale-95 uppercase tracking-[0.2em] flex items-center justify-center gap-3">
                            {{ isProcessing ? 'MENYIMPAN DATA...' : (isEditMode ? 'SIMPAN PERUBAHAN' : 'DAFTARKAN KARYAWAN') }}
                        </button>
                    </div>
                    
                </form>
            </div>
        </div>
    </div>
</template>

<style scoped>
@keyframes scan { 0% { top: 0%; opacity: 0; } 10% { opacity: 1; } 90% { opacity: 1; } 100% { top: 100%; opacity: 0; } }
</style>