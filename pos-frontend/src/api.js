import axios from 'axios';

const api = axios.create({
    baseURL: 'http://localhost:8080/api',
});

// Otomatis bawa tiket JWT setiap kali Kurir pergi
api.interceptors.request.use(config => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export default api;