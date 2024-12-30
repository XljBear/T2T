import { defineStore } from 'pinia'
export const useLoginStore = defineStore('login', {
    state: () => ({
        isLogin: false,
    }),
    actions: {
        login(token: string) {
            const date = new Date();
            date.setDate(date.getDate() + 7);
            document.cookie = `token=${token}; path=/; expires=${date.toUTCString()}`;
            this.isLogin = true;
        },
        logout() {
            document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC;';
            this.isLogin = false
        }
    }
});