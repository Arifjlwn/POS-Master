<script setup>
import { useAbsensi } from '../../composables/useAbsensi.js';
import Sidebar from '../../components/Sidebar.vue';

// Import Komponen Modular
import AbsensiHeader from '../../components/attendance/AbsensiHeader.vue';
import AbsensiList from '../../components/attendance/AbsensiList.vue';
import AttendanceLogTable from '../../components/attendance/AttendanceLogTable.vue';
import CameraModal from '../../components/attendance/CameraModal.vue';

const {
    karyawan, riwayat, urutanTanggalTerbaru, isLoading, isAiLoading, showCameraModal,
    tanggalDipilih, bulanDipilih, filterMode, currentTime, currentUser,
    absenTarget, me, isSubmitting, setVideoRef, jepretDanKirim, mulaiAbsen, 
    stopCamera, toggleSortTanggal, lihatFoto, downloadLaporan
} = useAbsensi();
</script>

<template>
    <Sidebar>
        <div class="p-4 md:p-8 max-w-7xl mx-auto font-sans bg-[#f8fafc] min-h-screen">
            
            <AbsensiHeader 
                :me="me"
            />

            <AbsensiList 
                :karyawan="karyawan"
                :currentUser="currentUser"
                @mulai-absen="mulaiAbsen"
            />

            <AttendanceLogTable 
                :riwayat="riwayat"
                :isLoading="isLoading"
                v-model:filterMode="filterMode"
                v-model:tanggalDipilih="tanggalDipilih"
                v-model:bulanDipilih="bulanDipilih"
                :urutanTanggalTerbaru="urutanTanggalTerbaru"
                :currentUser="currentUser"
                @toggle-sort="toggleSortTanggal"
                @download-laporan="downloadLaporan"
                @lihat-foto="lihatFoto"
            />

            <CameraModal 
                :show="showCameraModal"
                :target="absenTarget"
                :currentTime="currentTime"
                :isAiLoading="isAiLoading"
                :isSubmitting="isSubmitting"
                :setVideoRef="setVideoRef"
                @close="stopCamera"
                @capture="jepretDanKirim"
            />

        </div>
    </Sidebar>
</template>