<script setup lang="ts">
import { ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { axiosInstance } from '../utils/axios';
const dialogSettingVisible = ref(false);
const secretPassword = "~nononono$y0ucantsee.meme@";

const settingForm = ref({
    panel_password: secretPassword,
    repeat_panel_password: ""
});

const showSettingDialog = () => {
    settingForm.value.panel_password = secretPassword;
    settingForm.value.repeat_panel_password = "";
    dialogSettingVisible.value = true;
}

const updateSetting = () => {
    if (settingForm.value.panel_password == secretPassword) {
        dialogSettingVisible.value = false;
        return;
    } else if (settingForm.value.panel_password != settingForm.value.repeat_panel_password) {
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
    }).catch(() => { })

}
defineExpose({ showSettingDialog });
</script>
<template>
    <el-dialog draggable v-model="dialogSettingVisible" title="系统设置" width="100%" style="max-width:500px">
        <el-form :model="settingForm">
            <el-form-item label="面板管理密码:" :label-width="120">
                <el-input type="password" v-model="settingForm.panel_password" autocomplete="off" />
                <el-text class="mx-1" type="warning">保持默认为不修改密码</el-text>
            </el-form-item>
            <el-form-item v-show="settingForm.panel_password != secretPassword" label="重复输入密码:" :label-width="120">
                <el-input type="password" v-model="settingForm.repeat_panel_password" autocomplete="off" />
            </el-form-item>
        </el-form>
        <el-button type="primary" @click="updateSetting">保存设置</el-button>
    </el-dialog>
</template>
<style scoped lang="scss"></style>