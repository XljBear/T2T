<script setup lang="ts">
import { ref, reactive } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { axiosInstance } from '../utils/axios';
import { WarningFilled } from '@element-plus/icons-vue';
const emit = defineEmits(["onDirty"]);
const showAddDialog = ref(false);
const showAddIPRuleDialog = (ip: string, port: string) => {
    ipRuleForm.ip = ip;
    ipRuleForm.port = port;
    ipRuleForm.type = 0;
    ipRuleForm.end_time = null;
    ipRuleForm.reason = "";
    showAddDialog.value = true;
}
const ipRuleForm = reactive({
    ip: "",
    port: "",
    type: 0,
    end_time: null,
    reason: ""
})
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
        axiosInstance.post(`/ipRules/${ipRuleForm.type == 0 ? 'block' : 'allow'}`, ipRuleForm).then(() => {
            ElMessage({
                message: '添加防火墙规则成功',
                type: 'success',
            });
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
        <el-dialog v-model="showAddDialog" title="添加防火墙规则" width="100%" style="max-width: 500px;">
            <el-form :model="ipRuleForm">
                <el-form-item label="IP:" :label-width="80">
                    <el-input v-model="ipRuleForm.ip" autocomplete="off" />
                </el-form-item>
                <el-form-item label="端口:" :label-width="80">
                    <el-col :span="20">
                        <el-input v-model="ipRuleForm.port" autocomplete="off" />
                    </el-col>
                    <el-col :span="4">
                        <el-tooltip content="留空则为所有端口，多个端口可用 ',' 英文逗号分隔">
                            <el-icon size="large">
                                <WarningFilled />
                            </el-icon>
                        </el-tooltip>
                    </el-col>
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
                    <el-button @click="showAddDialog = false">取消</el-button>
                    <el-button type="primary" @click="addIPRule">添加</el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>
<style scoped lang="scss"></style>