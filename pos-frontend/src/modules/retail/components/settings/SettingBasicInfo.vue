<script setup>
import { onMounted, ref, watch } from "vue";

const props = defineProps({
  form: Object,
  logoPreview: String,
});

const emit = defineEmits(["update-file", "remove-logo"]);
const fileInput = ref(null);

const onLogoSelect = (e) => {
  const file = e.target.files[0];
  if (file) {
    if (file.size > 2 * 1024 * 1024) {
      alert("Maksimal ukuran file adalah 2MB");
      return;
    }
    emit("update-file", "logo", file, URL.createObjectURL(file));
  }
};

const removeImage = () => {
  emit("remove-logo");
  if (fileInput.value) fileInput.value.value = "";
};

const isEditingLocation = ref(false);
const toggleEditLocation = () => {
  isEditingLocation.value = !isEditingLocation.value;
};

// ==========================================
// 🚀 INTEGRASI API WILAYAH (IBNUX) - SAMA KAYA SETUP TOKO!
// ==========================================
const listProvinsi = ref([]);
const listKota = ref([]);
const listKecamatan = ref([]);
const listKelurahan = ref([]);

const regIds = ref({ provinsi: "", kota: "", kecamatan: "", kelurahan: "" });
const isLoading = ref({ reg: false, dist: false, vil: false });

const cleanStr = (str) => {
  return str ? String(str).trim().toUpperCase() : "";
};

const matchWilayah = (apiName, dbName) => {
  if (!apiName || !dbName) return false;
  // Benerin biar "Daerah Khusus Ibukota Jakarta" bisa match sama "DKI JAKARTA"
  let api = cleanStr(apiName)
    .replace(/KABUPATEN |KOTA |ADMINISTRASI /g, "")
    .replace("DAERAH KHUSUS IBUKOTA JAKARTA", "DKI JAKARTA")
    .trim();
  let db = cleanStr(dbName)
    .replace(/KABUPATEN |KOTA |ADMINISTRASI /g, "")
    .replace("DAERAH KHUSUS IBUKOTA JAKARTA", "DKI JAKARTA")
    .trim();
  return api === db;
};

// ==========================================
// LOAD DATA WILAYAH DARI IBNUX
// ==========================================
const loadKota = async (idProvinsi, targetName = null) => {
  isLoading.value.reg = true;
  try {
    const res = await fetch(
      `https://ibnux.github.io/data-indonesia/kabupaten/${idProvinsi}.json`,
    );
    const data = await res.json();
    listKota.value = data.map((item) => ({
      id: item.id,
      name: item.nama.toUpperCase(),
    }));

    if (targetName) {
      const found = listKota.value.find((r) =>
        matchWilayah(r.name, targetName),
      );
      if (found) {
        regIds.value.kota = found.id;
        props.form.kota = found.name;
        return found.id;
      }
    }
    return null;
  } catch (err) {
    console.error(err);
    return null;
  } finally {
    isLoading.value.reg = false;
  }
};

const loadKecamatan = async (idKota, targetName = null) => {
  isLoading.value.dist = true;
  try {
    const res = await fetch(
      `https://ibnux.github.io/data-indonesia/kecamatan/${idKota}.json`,
    );
    const data = await res.json();
    listKecamatan.value = data.map((item) => ({
      id: item.id,
      name: item.nama.toUpperCase(),
    }));

    if (targetName) {
      const found = listKecamatan.value.find((d) =>
        matchWilayah(d.name, targetName),
      );
      if (found) {
        regIds.value.kecamatan = found.id;
        props.form.kecamatan = found.name;
        return found.id;
      }
    }
    return null;
  } catch (err) {
    console.error(err);
    return null;
  } finally {
    isLoading.value.dist = false;
  }
};

const loadKelurahan = async (idKecamatan, targetName = null) => {
  isLoading.value.vil = true;
  try {
    const res = await fetch(
      `https://ibnux.github.io/data-indonesia/kelurahan/${idKecamatan}.json`,
    );
    const data = await res.json();
    listKelurahan.value = data.map((item) => ({
      id: item.id,
      name: item.nama.toUpperCase(),
    }));

    if (targetName) {
      const found = listKelurahan.value.find((v) =>
        matchWilayah(v.name, targetName),
      );
      if (found) {
        regIds.value.kelurahan = found.id;
        props.form.kelurahan = found.name;
        return found.id;
      }
    }
    return null;
  } catch (err) {
    console.error(err);
    return null;
  } finally {
    isLoading.value.vil = false;
  }
};

const initializeWilayah = async () => {
  if (!props.form.provinsi) return;
  const prov = listProvinsi.value.find((p) =>
    matchWilayah(p.name, props.form.provinsi),
  );
  if (!prov) return;

  regIds.value.provinsi = prov.id;
  props.form.provinsi = prov.name;

  if (!props.form.kota) return;
  const idKota = await loadKota(prov.id, props.form.kota);
  if (!idKota) return;

  if (!props.form.kecamatan) return;
  const idKecamatan = await loadKecamatan(idKota, props.form.kecamatan);
  if (!idKecamatan) return;

  if (props.form.kelurahan) {
    await loadKelurahan(idKecamatan, props.form.kelurahan);
  }
};

onMounted(async () => {
  try {
    const res = await fetch(
      `https://ibnux.github.io/data-indonesia/provinsi.json`,
    );
    const data = await res.json();
    listProvinsi.value = data.map((item) => ({
      id: item.id,
      name: item.nama.toUpperCase(),
    }));
  } catch (e) {
    console.error("Gagal load provinsi dari IBNUX:", e);
  }
});

watch(
  () => ({
    provinsi: props.form.provinsi,
    provincesLoaded: listProvinsi.value.length,
  }),
  async (val) => {
    if (val.provinsi && val.provincesLoaded > 0 && !regIds.value.provinsi) {
      await initializeWilayah();
    }
  },
  { immediate: true, deep: true },
);

// ==========================================
// HANDLE SELECT CHANGE OLEH USER
// ==========================================
const handleProvinsiChange = async () => {
  regIds.value.kota = "";
  regIds.value.kecamatan = "";
  regIds.value.kelurahan = "";
  listKota.value = [];
  listKecamatan.value = [];
  listKelurahan.value = [];
  props.form.kota = "";
  props.form.kecamatan = "";
  props.form.kelurahan = "";

  const prov = listProvinsi.value.find((p) => p.id === regIds.value.provinsi);
  if (prov) props.form.provinsi = prov.name;

  if (regIds.value.provinsi) {
    await loadKota(regIds.value.provinsi);
  }
};

const handleKotaChange = async () => {
  regIds.value.kecamatan = "";
  regIds.value.kelurahan = "";
  listKecamatan.value = [];
  listKelurahan.value = [];
  props.form.kecamatan = "";
  props.form.kelurahan = "";

  const kota = listKota.value.find((k) => k.id === regIds.value.kota);
  if (kota) props.form.kota = kota.name;

  if (regIds.value.kota) {
    await loadKecamatan(regIds.value.kota);
  }
};

const handleKecamatanChange = async () => {
  regIds.value.kelurahan = "";
  listKelurahan.value = [];
  props.form.kelurahan = "";

  const kec = listKecamatan.value.find((k) => k.id === regIds.value.kecamatan);
  if (kec) props.form.kecamatan = kec.name;

  if (regIds.value.kecamatan) {
    await loadKelurahan(regIds.value.kecamatan);
  }
};

const handleKelurahanChange = () => {
  const kel = listKelurahan.value.find((k) => k.id === regIds.value.kelurahan);
  if (kel) props.form.kelurahan = kel.name;
};

const forceUppercase = (field) => {
  if (props.form[field]) {
    props.form[field] = props.form[field].toUpperCase();
  }
};
</script>

<template>
  <div class="space-y-6 animate-fade-in-up w-full">
    <div class="flex flex-col md:flex-row gap-6">
      <div
        class="w-full md:w-1/3 flex flex-col items-center p-6 bg-slate-50 rounded-3xl border-2 border-slate-100"
      >
        <label
          class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-4"
          >Logo Toko</label
        >

        <div
          class="relative w-32 h-32 bg-white rounded-[24px] shadow-sm flex items-center justify-center overflow-hidden border-2 border-slate-200 group hover:border-blue-400 transition-colors"
        >
          <div v-if="logoPreview" class="w-full h-full relative cursor-default">
            <img :src="logoPreview" class="w-full h-full object-contain p-2" />
            <button
              @click.prevent="removeImage"
              class="absolute inset-0 bg-rose-500/80 text-white flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity backdrop-blur-sm cursor-pointer z-10"
            >
              <svg
                class="w-8 h-8"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="3"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>

          <div
            v-else
            @click="$refs.fileInput.click()"
            class="w-full h-full flex items-center justify-center cursor-pointer"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="w-8 h-8 text-slate-300 group-hover:text-blue-500 transition-colors"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
              />
            </svg>
          </div>
        </div>

        <input
          type="file"
          ref="fileInput"
          @change="onLogoSelect"
          class="hidden"
          accept="image/*"
        />
        <p class="text-[9px] text-slate-400 font-bold mt-4 text-center">
          Format: JPG, PNG (Max 2MB)
        </p>
      </div>

      <div class="w-full md:w-2/3 grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="md:col-span-2">
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Nama Toko / Outlet</label
          >
          <input
            v-model="form.nama_toko"
            @input="forceUppercase('nama_toko')"
            type="text"
            class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all"
          />
        </div>
        <div>
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >No. WA Bisnis</label
          >
          <input
            v-model="form.telepon"
            type="tel"
            class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all"
            placeholder="08xxxxxxxx"
          />
        </div>

        <div>
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Kategori Bisnis</label
          >
          <div
            class="w-full p-4 bg-slate-100/80 rounded-2xl border-2 border-slate-100 flex items-center justify-between cursor-not-allowed"
          >
            <span
              class="font-black text-xs text-slate-400 uppercase tracking-widest"
              >{{ form.business_type || "RETAIL" }}</span
            >
            <svg
              class="w-4 h-4 text-slate-300"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2.5"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
              />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- BARIS 2: LOKASI OPERASIONAL -->
    <div class="border-t border-slate-100 pt-6 space-y-4 w-full">
      <div class="flex items-center justify-between mb-2">
        <h4 class="font-black text-slate-800 uppercase tracking-widest text-sm">
          Lokasi Operasional
        </h4>
        <button
          @click.prevent="toggleEditLocation"
          class="text-[10px] font-black uppercase tracking-widest flex items-center gap-1 transition-colors"
          :class="
            isEditingLocation
              ? 'text-rose-500 hover:text-rose-600'
              : 'text-blue-500 hover:text-blue-600'
          "
        >
          <svg
            v-if="!isEditingLocation"
            xmlns="http://www.w3.org/2000/svg"
            class="h-4 w-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
            />
          </svg>
          <svg
            v-else
            xmlns="http://www.w3.org/2000/svg"
            class="h-4 w-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
          {{ isEditingLocation ? "Batal Ubah" : "Ubah Lokasi" }}
        </button>
      </div>

      <div
        v-if="isEditingLocation"
        class="bg-amber-50 border border-amber-200 text-amber-700 p-3 rounded-xl text-xs font-bold flex items-start gap-2 mb-4"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 shrink-0"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
        </svg>
        <p>
          Mengubah lokasi akan mempengaruhi sistem ongkir dan data pajak daerah.
          Pastikan data yang dimasukkan akurat.
        </p>
      </div>

      <!-- 🚀 PEMBARUAN GRID RESPONSIVE BERDASARKAN TOTAL 12 KOLOM -->
      <div
        class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-12 gap-4 w-full"
        :class="{ 'opacity-60 pointer-events-none': !isEditingLocation }"
      >
        <!-- Alamat Lengkap makan 12 kolom penuh -->
        <div class="col-span-1 sm:col-span-2 md:col-span-12">
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Alamat Lengkap (Jalan / Patokan)</label
          >
          <textarea
            v-model="form.alamat"
            @input="forceUppercase('alamat')"
            :disabled="!isEditingLocation"
            rows="2"
            class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all resize-none disabled:bg-slate-100/50"
            placeholder="Contoh: Jl. Sudirman No. 123..."
          ></textarea>
        </div>

        <!-- Provinsi makan 5 kolom karena teksnya panjang (ex: DAERAH KHUSUS IBUKOTA JAKARTA) -->
        <div
          class="relative w-full min-w-0 col-span-1 sm:col-span-1 md:col-span-5"
        >
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Provinsi</label
          >
          <select
            v-model="regIds.provinsi"
            @change="handleProvinsiChange"
            :disabled="!isEditingLocation"
            class="w-full p-4 pr-10 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all appearance-none cursor-pointer disabled:bg-slate-100/50"
          >
            <option value="" disabled>Pilih Provinsi</option>
            <option v-for="p in listProvinsi" :key="p.id" :value="p.id">
              {{ p.name }}
            </option>
          </select>
          <div
            class="absolute inset-y-0 right-4 flex items-center top-[28px] pointer-events-none text-slate-400"
          >
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              stroke-width="3"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M19 9l-7 7-7-7"
              ></path>
            </svg>
          </div>
        </div>

        <!-- Kota/Kabupaten makan 7 kolom sisanya di baris pertama -->
        <div
          class="relative w-full min-w-0 col-span-1 sm:col-span-1 md:col-span-7"
        >
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Kota / Kabupaten</label
          >
          <select
            v-model="regIds.kota"
            @change="handleKotaChange"
            :disabled="!regIds.provinsi || isLoading.reg || !isEditingLocation"
            class="w-full p-4 pr-10 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all appearance-none disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer disabled:bg-slate-100/50"
          >
            <option value="" disabled>
              {{ isLoading.reg ? "Loading..." : "Pilih Kota/Kab" }}
            </option>
            <option v-for="r in listKota" :key="r.id" :value="r.id">
              {{ r.name }}
            </option>
          </select>
          <div
            class="absolute inset-y-0 right-4 flex items-center top-[28px] pointer-events-none text-slate-400"
          >
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              stroke-width="3"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M19 9l-7 7-7-7"
              ></path>
            </svg>
          </div>
        </div>

        <!-- Kecamatan makan 4 kolom di baris kedua -->
        <div
          class="relative w-full min-w-0 col-span-1 sm:col-span-1 md:col-span-4"
        >
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Kecamatan</label
          >
          <select
            v-model="regIds.kecamatan"
            @change="handleKecamatanChange"
            :disabled="!regIds.kota || isLoading.dist || !isEditingLocation"
            class="w-full p-4 pr-10 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all appearance-none disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer disabled:bg-slate-100/50"
          >
            <option value="" disabled>
              {{ isLoading.dist ? "Loading..." : "Pilih Kecamatan" }}
            </option>
            <option v-for="d in listKecamatan" :key="d.id" :value="d.id">
              {{ d.name }}
            </option>
          </select>
          <div
            class="absolute inset-y-0 right-4 flex items-center top-[28px] pointer-events-none text-slate-400"
          >
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              stroke-width="3"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M19 9l-7 7-7-7"
              ></path>
            </svg>
          </div>
        </div>

        <!-- Kelurahan makan 5 kolom di baris kedua -->
        <div
          class="relative w-full min-w-0 col-span-1 sm:col-span-1 md:col-span-5"
        >
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Kelurahan / Desa</label
          >
          <select
            v-model="regIds.kelurahan"
            @change="handleKelurahanChange"
            :disabled="!regIds.kecamatan || isLoading.vil || !isEditingLocation"
            class="w-full p-4 pr-10 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all appearance-none disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer disabled:bg-slate-100/50"
          >
            <option value="" disabled>
              {{ isLoading.vil ? "Loading..." : "Pilih Kelurahan" }}
            </option>
            <option v-for="v in listKelurahan" :key="v.id" :value="v.id">
              {{ v.name }}
            </option>
          </select>
          <div
            class="absolute inset-y-0 right-4 flex items-center top-[28px] pointer-events-none text-slate-400"
          >
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              stroke-width="3"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M19 9l-7 7-7-7"
              ></path>
            </svg>
          </div>
        </div>

        <!-- Kode Pos makan 3 kolom sisa terakhir -->
        <div class="w-full min-w-0 col-span-1 sm:col-span-2 md:col-span-3">
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Kode Pos</label
          >
          <input
            v-model="form.kode_pos"
            :disabled="!isEditingLocation"
            type="text"
            class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all disabled:bg-slate-100/50"
            placeholder="Masukkan kode pos"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.4s ease-out;
}
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
