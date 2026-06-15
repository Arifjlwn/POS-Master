// 🚀 THE ULTIMATE MULTI-TENANT CATEGORY SVG REPOSITORY (VERSION 2.0 - FULL DATA)
// Murni SVG Heroicons kasta tertinggi, steril dari emoji alay bray bray!

export const getCategoryIcon = (categoryName) => {
	const cat = categoryName ? categoryName.toLowerCase().trim() : '';

	// 1. MATERIAL BANGUNAN (Arzu Baja Core)
	if (cat.includes('bangunan') || cat.includes('semen') || cat.includes('besi') || cat.includes('pipa') || cat.includes('paku') || cat.includes('kayu') || cat.includes('baja') || cat.includes('cat') || cat.includes('keramik') || cat.includes('baut') || cat.includes('tools') || cat.includes('perkakas') || cat.includes('material')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M11 4a2 2 0 114 0v1a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-1a2 2 0 100 4h1a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-1a2 2 0 10-4 0v1a1 1 0 01-1 1H7a1 1 0 01-1-1v-3a1 1 0 00-1-1H4a2 2 0 110-4h1a1 1 0 001-1V7a1 1 0 011-1h3a1 1 0 001-1V4z" />
        </svg>`;
	}

	// 2. OBAT & KESEHATAN / PHARMACY
	if (cat.includes('obat') || cat.includes('apotek') || cat.includes('farmasi') || cat.includes('herbal') || cat.includes('vitamin') || cat.includes('medis') || cat.includes('alkes') || cat.includes('salep') || cat.includes('kapsul') || cat.includes('tablet') || cat.includes('sirup') || cat.includes('suplemen') || cat.includes('p3k') || cat.includes('health') || cat.includes('pharmacy')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
        </svg>`;
	}

	// 3. MAKANAN / GROCERY (FMCG Makanan, Bumbu, Beras, Gula, Tepung, Cemilan)
	if (cat.includes('makanan') || cat.includes('snack') || cat.includes('biskuit') || cat.includes('roti') || cat.includes('permen') || cat.includes('coklat') || cat.includes('sereal') || cat.includes('kaleng') || cat.includes('grocery') || cat.includes('beras') || cat.includes('gula') || cat.includes('tepung') || cat.includes('bumbu') || cat.includes('cemilan') || cat.includes('bakery') || cat.includes('pastry') || cat.includes('kue')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707m12.728 0l-.707-.707M6.343 6.343l-.707-.707M14 12a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>`;
	}

	// 4. MINUMAN (Air Mineral, Teh, Kopi, Jus, Soda, Susu Kotak)
	if (cat.includes('minuman') || cat.includes('mineral') || cat.includes('teh') || cat.includes('kopi') || cat.includes('susu') || cat.includes('energi') || cat.includes('jus') || cat.includes('sirup') || cat.includes('soda') || cat.includes('drink') || cat.includes('beverage')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 5V3m0 2a3 3 0 100 6h0a3 3 0 100-6m0 6v10m6-5a6 6 0 11-12 0H6" />
        </svg>`;
	}

	// 5. BAYI / BABY CARE (Pampers, Susu Bayi, Tisu Basah)
	if (cat.includes('bayi') || cat.includes('baby') || cat.includes('pampers') || cat.includes('tisu basah') || cat.includes('popok') || cat.includes('bedak bayi')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>`;
	}

	// 6. PERAWATAN TUBUH & KOSMETIK (Personal Care, Sabun, Skincare, Parfum)
	if (cat.includes('tubuh') || cat.includes('personal') || cat.includes('care') || cat.includes('sabun') || cat.includes('shampoo') || cat.includes('shampo') || cat.includes('gigi') || cat.includes('deodoran') || cat.includes('skincare') || cat.includes('kosmetik') || cat.includes('bedak') || cat.includes('lipstik') || cat.includes('foundation') || cat.includes('maskara') || cat.includes('parfum')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
        </svg>`;
	}

	// 7. KEBERSHIHAN RUMAH & RUMAH TANGGA (Deterjen, Sapu, Pel, Alat Dapur)
	if (cat.includes('kebersihan') || cat.includes('household') || cat.includes('deterjen') || cat.includes('pewangi') || cat.includes('lantai') || cat.includes('piring') || cat.includes('tangga') || cat.includes('sapu') || cat.includes('pel') || cat.includes('ember') || cat.includes('spons') || cat.includes('dapur')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
        </svg>`;
	}

	// 8. ROKOK (Kretek, Filter, Cerutu)
	if (cat.includes('rokok') || cat.includes('kretek') || cat.includes('filter') || cat.includes('cerutu')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>`;
	}

	// 9. ATK / STATIONERY / BUKU
	if (cat.includes('atk') || cat.includes('stationery') || cat.includes('buku') || cat.includes('kertas') || cat.includes('pulpen') || cat.includes('pensil') || cat.includes('spidol') || cat.includes('kantor')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
        </svg>`;
	}

	// 10. ELEKTRONIK / LISTRIK (Lampu, Kabel, Baterai, Charger)
	if (cat.includes('elektronik') || cat.includes('listrik') || cat.includes('lampu') || cat.includes('kabel') || cat.includes('baterai') || cat.includes('charger') || cat.includes('gadget') || cat.includes('headset') || cat.includes('power bank') || cat.includes('case') || cat.includes('hp')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>`;
	}

	// 11. FASHION & SEPATU (Kaos, Kemeja, Sandal, Sepatu, Tas, Dompet)
	if (cat.includes('fashion') || cat.includes('kaos') || cat.includes('kemeja') || cat.includes('celana') || cat.includes('jaket') || cat.includes('pakaian') || cat.includes('alas kaki') || cat.includes('sandal') || cat.includes('sepatu') || cat.includes('tas') || cat.includes('dompet') || cat.includes('butik')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
        </svg>`;
	}

	// 12. HEWAN PELIHARAAN / PET SHOP
	if (cat.includes('hewan') || cat.includes('pet') || cat.includes('kucing') || cat.includes('anjing') || cat.includes('pakan') || cat.includes('whiskas')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
        </svg>`;
	}

	// 13. OTOMOTIF / SPAREPART / BENGKEL (Oli, Aki, Radiator)
	if (cat.includes('otomotif') || cat.includes('oli') || cat.includes('aki') || cat.includes('radiator') || cat.includes('sparepart') || cat.includes('bengkel') || cat.includes('motor') || cat.includes('mobil')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <circle cx="12" cy="12" r="3" stroke-linecap="round" stroke-linejoin="round" />
        </svg>`;
	}

	// 14. FROZEN FOOD (Nugget, Sosis, Bakso)
	if (cat.includes('frozen') || cat.includes('nugget') || cat.includes('sosis') || cat.includes('bakso')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m3.343-5.657l.707.707m2.122 2.122l.707.707m3.536 3.536l.707.707m-9.193 2.122l.707-.707m7.071-7.072l.707-.707" />
        </svg>`;
	}

	// 15. FRESH FOOD / BUAH & SAYUR / DAGING & SEAFOOD
	if (cat.includes('buah') || cat.includes('sayur') || cat.includes('fresh') || cat.includes('daging') || cat.includes('seafood') || cat.includes('ayam') || cat.includes('ikan') || cat.includes('telur')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
        </svg>`;
	}

	// 16. DIGITAL PRODUCT & PULSA (Paket Data, Token PLN, Voucher Game)
	if (cat.includes('pulsa') || cat.includes('digital') || cat.includes('data') || cat.includes('token') || cat.includes('pln') || cat.includes('voucher') || cat.includes('game')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
        </svg>`;
	}

	// 17. PERTANIAN (Pupuk, Benih, Pestisida)
	if (cat.includes('tani') || cat.includes('pertanian') || cat.includes('pupuk') || cat.includes('benih') || cat.includes('pestisida')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
        </svg>`;
	}

	// 18. MAINAN & SOUVENIR (Mainan Anak, Edukasi, Boneka, Souvenir, Parcel, Gift)
	if (cat.includes('mainan') || cat.includes('boneka') || cat.includes('souvenir') || cat.includes('gift') || cat.includes('parcel') || cat.includes('bingkisan') || cat.includes('hadiah')) {
		return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5a2 2 0 10-2 2h2zm0 0h4m-4 0H8m13 4v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6m18 0H3m18 0h-5M3 13h5m0 0v8m8-8v8" />
        </svg>`;
	}

	// ⚙️ DEFAULT FALLBACK: Folder Minimalis Clean untuk Kategori "Lainnya" / Umum
	return `<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
    </svg>`;
};
