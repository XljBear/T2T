import axios from 'axios';
let devMode = true;
if (import.meta.env.MODE === 'production') {
    devMode = false;
}
export const axiosInstance = axios.create({
    baseURL: devMode ? 'http://127.0.0.1:8080/api' : '../api',
    timeout: 3000,
});