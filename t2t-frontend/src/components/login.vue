<script setup lang="ts">
import { ref, reactive } from 'vue';
import { axiosInstance } from '../utils/axios';
import { ElMessage } from 'element-plus';
import { ClickDot } from 'go-captcha-vue/dist/components/click/meta/data.js';
import { useLoginStore } from '../stores/login';
const loginStore = useLoginStore();
const captchaRef = ref();
const captchaData = reactive({
  image: '',
  thumb: '',
  captcha_id: '',
});
const loginForm = reactive({
  password: '',
})
const dialogCaptchaVisible = ref(false);
const refreshCaptcha = () => {
  axiosInstance.get('/captcha').then(res => {
    captchaData.image = res.data.captcha;
    captchaData.thumb = res.data.thumb;
    captchaData.captcha_id = res.data.captcha_id;
  });
};
const captchaConfirm = (dots: Array<ClickDot>, reset: () => void) => {
  if (dots.length === 0) {
    ElMessage.warning('请点击图像进行安全验证');
    return;
  }
  axiosInstance.post('/login', {
    captcha_id: captchaData.captcha_id,
    password: loginForm.password,
    captcha_data: dots
  }).then(res => {
    loginStore.login(res.data.token);
    reset();
  }).catch(err => {
    switch (err.response.data.error) {
      case 'Invalid request data':
        ElMessage.error('参数错误，请重新尝试');
        break;
      case 'Invalid captchaID':
        ElMessage.error('验证码已失效，请重新尝试');
        break;
      case 'Captcha verification failed':
        ElMessage.error('验证码输入错误，请重新尝试');
        break;
      case 'Invalid password':
        dialogCaptchaVisible.value = false;
        loginForm.password = '';
        ElMessage.error('密码错误，请重新输入');
        return;
      default:
        ElMessage.error(err.response.data.error);
        break;
    };
    refreshCaptcha();
  }).finally(() => {
    reset();
  });
};
const showCaptcha = () => {
  dialogCaptchaVisible.value = true;
  refreshCaptcha();
};
</script>
<template>
  <div class="login">
    <el-card style="width: 80%;">
      <template #header>
        <div class="card-header">T2T 控制面板</div>
      </template>
      <div class="loginForm">
        <el-form label-width="auto" :model="loginForm" style="max-width: 600px">
          <el-form-item label="管理密码:" label-position="right">
            <el-col :span="16"><el-input type="password" v-model="loginForm.password" /></el-col>
            <el-col :span="8"><el-button type="primary" @click="showCaptcha">登录</el-button></el-col>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <div class="card-footer">T2T Server v0.0.1 © StupidBear Studio 2024</div>
      </template>
    </el-card>
    <el-dialog draggable v-model="dialogCaptchaVisible" title="安全验证" width="100%" style="max-width:400px">
      <div class="captcha">
        <gocaptcha-click :config="{}" :data="captchaData" :events="{ refresh: refreshCaptcha, confirm: captchaConfirm }"
          ref="captchaRef" />
      </div>
    </el-dialog>
  </div>
</template>
<style scoped lang="scss">
.login {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;

  .card-header {
    color: #333;
    font-size: 24px;
    font-weight: bold;
  }

  .loginForm {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .captcha {
    width: 100%;
    display: flex;
    justify-content: center;
  }

  .card-footer {
    color: #ccc;
    font-size: 12px;
  }
}
</style>