<script setup lang="ts">
import { onMounted } from 'vue'
import proxyPanel from './components/proxyPanel.vue';
import login from './components/login.vue';
import { useLoginStore } from './stores/login';
const loginStore = useLoginStore();
onMounted(() => {
  const cookies = document.cookie;
  if (cookies.includes('token=')) {
    loginStore.isLogin = true;
  }
});
</script>

<template>
  <div class="app">
    <proxyPanel v-if="loginStore.isLogin" />
    <login v-else />
  </div>
</template>

<style scoped lang="scss">
.app {
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
</style>
