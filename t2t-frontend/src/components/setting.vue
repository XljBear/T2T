<script setup lang="ts">
import { ref, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { axiosInstance } from '../utils/axios';
import { Moon, Sunny } from '@element-plus/icons-vue';
import { usePanelStore } from '../stores/panel';
import { useDark } from "@vueuse/core";
const isDark = useDark();
const panelStore = usePanelStore();
const dialogSettingVisible = ref(false);
const secretPassword = "~nononono$y0ucantsee.meme@";

const settingForm = ref({
    panel_password: secretPassword,
    repeat_panel_password: "",
    captcha_type: 0,
    dark_mode: false,
});

const showSettingDialog = () => {
    panelStore.refreshConfig();
    settingForm.value.panel_password = secretPassword;
    settingForm.value.repeat_panel_password = "";
    settingForm.value.dark_mode = panelStore.darkMode;
    settingForm.value.captcha_type = panelStore.captchaType;
    dialogSettingVisible.value = true;
};

const updateSetting = () => {
    if (settingForm.value.panel_password != secretPassword && settingForm.value.panel_password != settingForm.value.repeat_panel_password) {
        ElMessage({
            message: '重复密码输入不正确',
            type: 'warning',
        });
        return;
    }
    ElMessageBox.confirm(
        '确认要更新设置吗？',
        '更新系统设置',
        {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
        }
    ).then(() => {
        axiosInstance.post("/setting", settingForm.value).then(() => {
            ElMessage({
                message: '系统设置更新成功',
                type: 'success',
            });
            dialogSettingVisible.value = false;
        }).catch((resp) => {
            ElMessage({
                message: '更新系统设置失败:' + resp.response.data.error,
                type: 'error',
            });
        })
    }).catch(() => { });
}
watch(() => dialogSettingVisible.value, (val: boolean) => {
    if (!val) {
        panelStore.refreshConfig();
    }
})
defineExpose({ showSettingDialog });
</script>
<template>
    <el-dialog draggable v-model="dialogSettingVisible" title="系统设置" width="100%" style="max-width:500px">
        <el-form :model="settingForm">
            <el-form-item label="面板管理密码:" :label-width="120">
                <el-input placeholder="请输入新的管理密码，可留空" type="password" v-model="settingForm.panel_password"
                    autocomplete="off" clearable show-password />
                <el-text class="mx-1" type="warning">保持默认为不修改密码</el-text>
            </el-form-item>
            <el-form-item v-show="settingForm.panel_password != secretPassword" label="重复输入密码:" :label-width="120">
                <el-input placeholder="请再次输入相同密码" type="password" v-model="settingForm.repeat_panel_password"
                    autocomplete="off" clearable show-password />
            </el-form-item>
            <el-form-item label="面板登陆行为验证:" :label-width="140">
                <el-radio-group v-model="settingForm.captcha_type">
                    <el-radio :value="0">关闭</el-radio>
                    <el-radio :value="1">点选式(文字)</el-radio>
                    <el-radio :value="2">点选式(图形)</el-radio>
                    <el-radio :value="3">滑动式</el-radio>
                    <el-radio :value="4">拖拽式</el-radio>
                    <el-radio :value="5">旋转拼图</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item label="面板主题:" :label-width="140">

                <el-switch v-model="settingForm.dark_mode" :active-icon="Moon" :inactive-icon="Sunny" inline-prompt
                    @change="isDark = settingForm.dark_mode" />
            </el-form-item>
        </el-form>
        <el-button type="primary" @click="updateSetting">保存设置</el-button>
    </el-dialog>
</template>
<style scoped lang="scss"></style>