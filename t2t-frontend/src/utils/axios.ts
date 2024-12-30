import axios from 'axios';
import { useLoginStore } from '../stores/login';
import { ElMessage } from 'element-plus';
let devMode = true;
if (import.meta.env.MODE === 'production') {
    devMode = false;
}
export const axiosInstance = axios.create({
    baseURL: devMode ? './api' : '../api',
    timeout: 3000,
});
axiosInstance.interceptors.response.use(function (response) {
    return response;
}, function (error) {
    const loginStore = useLoginStore();
    if (error.response.status === 401 && loginStore.isLogin) {
        loginStore.logout();
        ElMessage.error('登录已过期，请重新登录');
    }
    return Promise.reject(error);
});