<script setup lang="ts">
import { ref, reactive } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { axiosInstance } from '../utils/axios';
const emit = defineEmits(["onDirty"]);
const showAddDialog = ref(false);
const showAddIPRuleDialog = (ip: string, port: string) => {
    ipRuleForm.ip = ip;
    if (port) {
        ipRuleForm.port = [port];
    }
    ipRuleForm.type = 0;
    ipRuleForm.end_time = null;
    ipRuleForm.reason = "";
    showAddDialog.value = true;
}
const ipRuleForm = reactive({
    ip: "",
    port: [] as string[],
    type: 0,
    end_time: null,
    reason: ""
})
const cleanRuleForm = () => {
    ipRuleForm.ip = "";
    ipRuleForm.port = [];
    ipRuleForm.type = 0;
    ipRuleForm.end_time = null;
    ipRuleForm.reason = "";
}
const closeDialog = () => {
    cleanRuleForm();
    showAddDialog.value = false;
}
const addIPRule = () => {
    if (ipRuleForm.ip == "") {
        ElMessage({
            message: '请先输入IP地址',
            type: 'warning',
        });
        return;
    }
    ElMessageBox.confirm(
        '确认要添加该规则吗?',
        '添加防火墙规则',
        {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
        }
    ).then(() => {
        axiosInstance.post('/ipRules', ipRuleForm).then(() => {
            ElMessage({
                message: '添加防火墙规则成功',
                type: 'success',
            });
            cleanRuleForm();
            showAddDialog.value = false;
            emit("onDirty");
        }).catch((resp) => {
            ElMessage({
                message: '添加防火墙规则失败:' + resp.response.data.error,
                type: 'error',
            });
        });
    });
};
defineExpose({ showAddIPRuleDialog });
</script>
<template>
    <div>
        <el-dialog @close="closeDialog" v-model="showAddDialog" title="添加防火墙规则" width="100%" style="max-width: 500px;">
            <el-form :model="ipRuleForm">
                <el-form-item label="IP:" :label-width="80">
                    <el-input v-model="ipRuleForm.ip" autocomplete="off" />
                </el-form-item>
                <el-form-item label="端口:" :label-width="80">
                    <el-input-tag v-model="ipRuleForm.port" placeholder="留空则为对所有端口有效" />
                </el-form-item>
                <el-form-item label="规则类型:" :label-width="80">
                    <el-select v-model="ipRuleForm.type">
                        <el-option label="黑名单" :value="0" />
                        <el-option label="白名单" :value="1" />
                    </el-select>
                </el-form-item>
                <el-form-item label="失效时间:" :label-width="90">
                    <el-date-picker v-model="ipRuleForm.end_time" type="datetime" placeholder="可使该规则自动失效" />
                </el-form-item>
                <el-form-item label="原由:" :label-width="80">
                    <el-input v-model="ipRuleForm.reason" autocomplete="off" />
                </el-form-item>
            </el-form>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="closeDialog">取消</el-button>
                    <el-button type="primary" @click="addIPRule">添加</el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>
<style scoped lang="scss"></style>