import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import "go-captcha-vue/dist/style.css"
import GoCaptcha from "go-captcha-vue"
import { createPinia } from 'pinia'
import 'element-plus/theme-chalk/dark/css-vars.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
createApp(App).use(createPinia()).use(GoCaptcha).use(ElementPlus, { locale: zhCn }).mount('#app')
