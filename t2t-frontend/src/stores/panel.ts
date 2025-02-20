import { defineStore } from 'pinia'
import { useDark } from "@vueuse/core";
import { axiosInstance } from '../utils/axios';
import { ElMessage } from 'element-plus';
const isDark = useDark();
export const usePanelStore = defineStore('panel', {
    state: () => ({
        darkMode: false,
        captchaType: 0,
        commitId: '',
        version: '',
        buildTime: '',
        os: '',
        arch: ''
    }),
    actions: {
        async refreshInfo() {
            await axiosInstance.get('/info').then(res => {
                this.captchaType = res.data.captcha_type;
                this.darkMode = res.data.dark_mode;
                this.commitId = res.data.commit_id;
                this.version = res.data.version;
                this.buildTime = res.data.build_time;
                this.os = res.data.os;
                this.arch = res.data.arch;
                isDark.value = this.darkMode;
            }).catch(() => {
                ElMessage.error('获取配置参数失败');
            });
        }
    },
});