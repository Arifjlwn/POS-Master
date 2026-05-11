import axios from "axios";

// Bikin instance Axios Khusus
const api = axios.create({
    baseURL: "http://localhost:8080/api", //Arahkan ke Golang
    headers: {
        "Content-Type": "application/json",
    },
});

// Interceptor: Satpam FE sebelum data dikirim
api.interceptors.request.use(
    (config) => {
        // Ambil tiket JWT dari localStorage
        const token = localStorage.getItem("token");

        // Kalau tiket ada, tempel di header Auth
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    },
);

export default api;
