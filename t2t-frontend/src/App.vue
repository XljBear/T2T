<script setup lang="ts">
import { onMounted, ref, reactive, watch } from 'vue'
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
  await panelStore.refreshInfo();
  inInit.value = false;
});
const font = reactive({
  color: 'rgba(0, 0, 0, .15)',
})
watch(
  () => panelStore.darkMode,
  () => {
    font.color = panelStore.darkMode
      ? 'rgba(255, 255, 255, .15)'
      : 'rgba(0, 0, 0, .15)'
  },
  {
    immediate: true,
  }
)
</script>

<template>
  <el-watermark class="watermark" content="T2T Server" :font="font">
    <div v-if="!inInit" class="app">
      <proxyPanel v-if="loginStore.isLogin" />
      <login v-else />
    </div>
  </el-watermark>
</template>

<style scoped lang="scss">
.app {
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.watermark {
  min-height: 100vh;
  display: flex;
  justify-items: center;
}
</style>
