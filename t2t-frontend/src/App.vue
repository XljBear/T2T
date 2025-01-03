<script setup lang="ts">
import { onMounted, ref } from 'vue'
import proxyPanel from './components/proxyPanel.vue';
import login from './components/login.vue';
import { useLoginStore } from './stores/login';
import { usePanelStore } from './stores/panel';
const loginStore = useLoginStore();
const panelStore = usePanelStore();
const inInit = ref(true);
onMounted(async () => {
  const cookies = document.cookie;
  if (cookies.includes('token=')) {
    loginStore.isLogin = true;
  };
  await panelStore.refreshConfig();
  inInit.value = false;
});
</script>

<template>
  <div v-if="!inInit" class="app">
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
