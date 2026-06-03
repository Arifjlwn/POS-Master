<script setup>
  import { ref, watch, onMounted } from 'vue';

  const props = defineProps({
    form: Object,
    logoPreview: String,
  });

  const emit = defineEmits(['update-file', 'remove-logo']);
  const fileInput = ref(null);

  const onLogoSelect = (e) => {
    const file = e.target.files[0];
    if (file) {
      if (file.size > 2 * 1024 * 1024) {
        alert('Maksimal ukuran file adalah 2MB');
        return;
      }
      emit('update-file', 'logo', file, URL.createObjectURL(file));
    }
  };

  const removeImage = () => {
    emit('remove-logo');
    if (fileInput.value) fileInput.value.value = '';
  };

  // ==========================================
  // 🚀 INTEGRASI API WILAYAH INDONESIA
  // ==========================================
  const provinces = ref([]);
  const regencies = ref([]);
  const districts = ref([]);
  const villages = ref([]);

  const isLoading = ref({ reg: false, dist: false, vil: false });
  const apiWilayah = 'https://www.emsifa.com/api-wilayah-indonesia/api';

  // Fungsi sakti buat bersihin spasi dan nge-uppercase string
  const cleanStr = (str) => {
    return str ? String(str).trim().toUpperCase() : '';
  };

  const matchWilayah = (apiName, dbName) => {
    const api = cleanStr(apiName)
      .replace('KABUPATEN ', '')
      .replace('KOTA ', '');

    const db = cleanStr(dbName).replace('KABUPATEN ', '').replace('KOTA ', '');

    return api === db;
  };

  // ==========================================
  // LOAD DATA WILAYAH
  // ==========================================

  const loadRegencies = async (provinsiName) => {
    const prov = provinces.value.find(
      (p) => cleanStr(p.name) === cleanStr(provinsiName)
    );

    if (!prov) return false;

    isLoading.value.reg = true;

    try {
      const res = await fetch(`${apiWilayah}/regencies/${prov.id}.json`);

      regencies.value = await res.json();

      if (props.form.kota) {
        const found = regencies.value.find((r) =>
          matchWilayah(r.name, props.form.kota)
        );

        if (found) {
          props.form.kota = found.name;
        }
      }

      return true;
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      isLoading.value.reg = false;
    }
  };

  const loadDistricts = async (kotaName) => {
    const reg = regencies.value.find((r) => matchWilayah(r.name, kotaName));

    if (!reg) return false;

    isLoading.value.dist = true;

    try {
      const res = await fetch(`${apiWilayah}/districts/${reg.id}.json`);

      districts.value = await res.json();

      if (props.form.kecamatan) {
        const found = districts.value.find((d) =>
          matchWilayah(d.name, props.form.kecamatan)
        );

        if (found) {
          props.form.kecamatan = found.name;
        }
      }

      return true;
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      isLoading.value.dist = false;
    }
  };

  const loadVillages = async (kecamatanName) => {
    const dist = districts.value.find((d) =>
      matchWilayah(d.name, kecamatanName)
    );

    if (!dist) return false;

    isLoading.value.vil = true;

    try {
      const res = await fetch(`${apiWilayah}/villages/${dist.id}.json`);

      villages.value = await res.json();

      if (props.form.kelurahan) {
        const found = villages.value.find((v) =>
          matchWilayah(v.name, props.form.kelurahan)
        );

        if (found) {
          props.form.kelurahan = found.name;
        }
      }

      return true;
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      isLoading.value.vil = false;
    }
  };

  // ==========================================
  // INIT DARI DATABASE
  // ==========================================

  const initializeWilayah = async () => {
    if (!props.form.provinsi) return;

    const loadedReg = await loadRegencies(props.form.provinsi);

    if (!loadedReg) return;

    if (!props.form.kota) return;

    const loadedDist = await loadDistricts(props.form.kota);

    if (!loadedDist) return;

    if (!props.form.kecamatan) return;

    await loadVillages(props.form.kecamatan);
  };

  // ==========================================
  // LOAD PROVINSI
  // ==========================================

  onMounted(async () => {
    try {
      const res = await fetch(`${apiWilayah}/provinces.json`);

      provinces.value = await res.json();

      await initializeWilayah();
    } catch (e) {
      console.error('Gagal load provinsi', e);
    }
  });

  // ==========================================
  // AUTO DETEKSI DATA DB DARI SUPABASE
  // ==========================================

  watch(
    () => ({
      provinsi: props.form.provinsi,
      kota: props.form.kota,
      kecamatan: props.form.kecamatan,
      kelurahan: props.form.kelurahan,
      provincesLoaded: provinces.value.length,
    }),
    async (val) => {
      if (
        val.provinsi &&
        val.provincesLoaded > 0 &&
        regencies.value.length === 0
      ) {
        await initializeWilayah();
      }
    },
    {
      immediate: true,
      deep: true,
    }
  );

  // ==========================================
  // USER GANTI PROVINSI
  // ==========================================

  const handleProvinsiChange = async () => {
    regencies.value = [];
    districts.value = [];
    villages.value = [];

    props.form.kota = '';
    props.form.kecamatan = '';
    props.form.kelurahan = '';

    if (!props.form.provinsi) return;

    await loadRegencies(props.form.provinsi);
  };

  // ==========================================
  // USER GANTI KOTA
  // ==========================================

  const handleKotaChange = async () => {
    districts.value = [];
    villages.value = [];

    props.form.kecamatan = '';
    props.form.kelurahan = '';

    if (!props.form.kota) return;

    await loadDistricts(props.form.kota);
  };

  // ==========================================
  // USER GANTI KECAMATAN
  // ==========================================

  const handleKecamatanChange = async () => {
    villages.value = [];

    props.form.kelurahan = '';

    if (!props.form.kecamatan) return;

    await loadVillages(props.form.kecamatan);
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
              >{{ form.business_type || 'RETAIL' }}</span
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

    <div class="border-t border-slate-100 pt-6 space-y-4 w-full">
      <h4
        class="font-black text-slate-800 uppercase tracking-widest text-sm mb-2"
      >
        Lokasi Operasional
      </h4>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 w-full">
        <div class="md:col-span-2 lg:col-span-3">
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Alamat Lengkap (Jalan / Patokan)</label
          >
          <textarea
            v-model="form.alamat"
            rows="2"
            class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all resize-none"
            placeholder="Contoh: Jl. Sudirman No. 123..."
          ></textarea>
        </div>

        <div class="relative">
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Provinsi</label
          >
          <select
            v-model="form.provinsi"
            @change="handleProvinsiChange"
            class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all appearance-none cursor-pointer"
          >
            <option value="" disabled>Pilih Provinsi</option>
            <option v-for="p in provinces" :key="p.id" :value="p.name">
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

        <div class="relative">
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Kota / Kabupaten</label
          >
          <select
            v-model="form.kota"
            @change="handleKotaChange"
            :disabled="!form.provinsi || isLoading.reg"
            class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all appearance-none disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer"
          >
            <option value="" disabled>
              {{ isLoading.reg ? 'Loading...' : 'Pilih Kota/Kab' }}
            </option>
            <option v-for="r in regencies" :key="r.id" :value="r.name">
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

        <div class="relative">
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Kecamatan</label
          >
          <select
            v-model="form.kecamatan"
            @change="handleKecamatanChange"
            :disabled="!form.kota || isLoading.dist"
            class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all appearance-none disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer"
          >
            <option value="" disabled>
              {{ isLoading.dist ? 'Loading...' : 'Pilih Kecamatan' }}
            </option>
            <option v-for="d in districts" :key="d.id" :value="d.name">
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

        <div class="relative">
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Kelurahan / Desa</label
          >
          <select
            v-model="form.kelurahan"
            :disabled="!form.kecamatan || isLoading.vil"
            class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all appearance-none disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer"
          >
            <option value="" disabled>
              {{ isLoading.vil ? 'Loading...' : 'Pilih Kelurahan' }}
            </option>
            <option v-for="v in villages" :key="v.id" :value="v.name">
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

        <div>
          <label
            class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2"
            >Kode Pos</label
          >
          <input
            v-model="form.kode_pos"
            type="text"
            class="w-full p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 focus:bg-white focus:border-blue-600 focus:shadow-lg focus:shadow-blue-500/10 outline-none font-bold text-sm text-slate-800 transition-all"
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
