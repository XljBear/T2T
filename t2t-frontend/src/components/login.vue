<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { axiosInstance } from '../utils/axios';
import { ElMessage } from 'element-plus';
import { useLoginStore } from '../stores/login';
const loginStore = useLoginStore();
const captchaRef = ref();
const captchaData = reactive({
  image: '',
  thumb: '',
  thumbWidth: 0,
  thumbHeight: 0,
  thumbX: 0,
  thumbY: 0,
  captcha_id: '',

});
const captchaType = ref(0);
const loginForm = reactive({
  password: '',
})
onMounted(() => {
  axiosInstance.get('/config').then(res => {
    captchaType.value = res.data.captcha_type;
  }).catch(() => {
    ElMessage.error('获取配置参数失败');
  });
})
const dialogCaptchaVisible = ref(false);
const refreshCaptcha = () => {
  axiosInstance.get('/captcha').then(res => {
    captchaData.image = res.data.captcha;
    captchaData.thumb = res.data.thumb;
    captchaData.thumbWidth = res.data.thumb_width;
    captchaData.thumbHeight = res.data.thumb_height;
    captchaData.thumbX = res.data.thumb_x;
    captchaData.thumbY = res.data.thumb_y;
    captchaData.captcha_id = res.data.captcha_id;
  });
};
const captchaConfirm = (captcha: any, reset: () => void) => {
  if (captcha.length === 0) {
    ElMessage.warning('请点击图像进行安全验证');
    return;
  }
  axiosInstance.post('/login', {
    captcha_id: captchaData.captcha_id,
    password: loginForm.password,
    click_captcha_data: captchaType.value == 1 || captchaType.value == 2 ? captcha : null,
    slide_captcha_data: captchaType.value == 3 || captchaType.value == 4 ? captcha : null,
    rotate_captcha_data: captchaType.value == 5 ? captcha : null,
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
  if (captchaType.value === 0) {
    axiosInstance.post('/login', {
      captcha_id: captchaData.captcha_id,
      password: loginForm.password
    }).then(res => {
      loginStore.login(res.data.token);
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
    });
  } else {
    dialogCaptchaVisible.value = true;
    refreshCaptcha();
  }
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
            <el-col :span="16"><el-input type="password" placeholder="请输入管理密码" v-model="loginForm.password" clearable
                show-password /></el-col>
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
        <gocaptcha-click v-if="captchaType == 1 || captchaType == 2" :config="{}" :data="captchaData"
          :events="{ refresh: refreshCaptcha, confirm: captchaConfirm }" ref="captchaRef" />
        <gocaptcha-slide v-if="captchaType == 3" :config="{}" :data="captchaData"
          :events="{ refresh: refreshCaptcha, confirm: captchaConfirm }" ref="captchaRef" />
        <gocaptcha-slide-region v-if="captchaType == 4" :config="{}" :data="captchaData"
          :events="{ refresh: refreshCaptcha, confirm: captchaConfirm }" ref="captchaRef" />
        <gocaptcha-rotate v-if="captchaType == 5" :config="{}" :data="captchaData"
          :events="{ refresh: refreshCaptcha, confirm: captchaConfirm }" ref="captchaRef" />
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