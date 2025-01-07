import { defineStore } from 'pinia'
import { axiosInstance } from '../utils/axios';
import { ElMessage } from 'element-plus';
export const useIPRulesStore = defineStore('ipRules', {
    state: () => ({
        mode: 0,
        allow: {} as any,
        block: {} as any
    }),
    actions: {
        async refreshIPRules() {
            await axiosInstance.get('/ipRules').then(res => {
                this.mode = res.data.mode,
                    this.allow = res.data.allow,
                    this.block = res.data.block
            }).catch(() => {
                ElMessage.error('防火墙规则加载失败');
            });
        }
    },
});