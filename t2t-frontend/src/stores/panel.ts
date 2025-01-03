import { defineStore } from 'pinia'
import { useDark } from "@vueuse/core";
import { axiosInstance } from '../utils/axios';
import { ElMessage } from 'element-plus';
const isDark = useDark();
export const usePanelStore = defineStore('panel', {
    state: () => ({
        darkMode: false,
        captchaType: 0,
    }),
    actions: {
        async refreshConfig() {
            await axiosInstance.get('/config').then(res => {
                this.captchaType = res.data.captcha_type;
                this.darkMode = res.data.dark_mode;
                isDark.value = this.darkMode;
            }).catch(() => {
                ElMessage.error('获取配置参数失败');
            });
        }
    },
});