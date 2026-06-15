import { ref } from 'vue';

// State global tunggal terkunci aman bray
const isGlobalLoading = ref(false);
let loadingTimeout = null;

export function useLoading() {
	const startLoading = () => {
		if (loadingTimeout) clearTimeout(loadingTimeout);
		isGlobalLoading.value = true;
	};

	const stopLoading = () => {
		// Kasih delay mini 200ms biar animasinya keliatan smooth kaga kaku jepret hilang bray
		loadingTimeout = setTimeout(() => {
			isGlobalLoading.value = false;
		}, 200);
	};

	return {
		isGlobalLoading,
		startLoading,
		stopLoading,
	};
}
