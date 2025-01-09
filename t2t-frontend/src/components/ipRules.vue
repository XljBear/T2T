<script setup lang="ts">
import { ref } from 'vue';
import { useIPRulesStore } from '../stores/ipRules';
import { ElMessage, ElMessageBox } from 'element-plus';
import { axiosInstance } from '../utils/axios';
import IPRulesList from './ipRulesList.vue';
import { Plus, Refresh } from '@element-plus/icons-vue';
import AddIPRule from './addIPRule.vue';
const ipRulesStore = useIPRulesStore();
const dialogIPRulesVisible = ref(false);
const addIPRuleDialog = ref();
const ipRulesForm = ref({
    mode: 0,
    allow: [],
    block: []
});
const dataInLoading = ref(false);
const refreshIPRules = async () => {
    dataInLoading.value = true;
    await ipRulesStore.refreshIPRules();
    ipRulesForm.value.mode = ipRulesStore.mode;
    ipRulesForm.value.allow = ipRulesStore.allow.allow_ips || [];
    ipRulesForm.value.block = ipRulesStore.block.block_ips || [];
    dataInLoading.value = false;
}
const showIPRulesDialog = async () => {
    await refreshIPRules();
    dialogIPRulesVisible.value = true;
};
const modeChange = () => {
    ElMessageBox.confirm(
        '确认要修改防火墙运行模式吗?',
        '防火墙运行模式',
        {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
        }
    ).then(() => {
        axiosInstance.put("/ipRules/mode", { mode: ipRulesForm.value.mode }).then(() => {
            ElMessage({
                message: '防火墙运行模式修改成功',
                type: 'success',
            });
        }).catch((resp) => {
            ElMessage({
                message: '防火墙运行模式修改失败:' + resp.response.data.error,
                type: 'error',
            });
            ipRulesForm.value.mode = ipRulesStore.mode;
        });
    }).catch(() => {
        ipRulesForm.value.mode = ipRulesStore.mode;
    });
};
const onDirty = () => {
    refreshIPRules();
}
const reloadIPRules = () => {
    dataInLoading.value = true;
    axiosInstance.post('/ipRules/reload', {}).then(() => {
        ElMessage({
            message: '防火墙配置重载成功',
            type: 'success',
        });
    }).catch(() => {
        ElMessage({
            message: '防火墙配置重载失败',
            type: 'error',
        });
    }).finally(() => {
        refreshIPRules();
    });
};
const showAddIPRuleDialog = () => {
    addIPRuleDialog.value.showAddIPRuleDialog();
}
defineExpose({ showIPRulesDialog });
</script>
<template>
    <el-dialog v-model="dialogIPRulesVisible" title="防火墙配置" width="100%" style="max-width:800px">
        <el-form class="actionForm" label-position="top" :model="ipRulesForm">
            <el-form-item label="运行模式:">
                <el-radio-group @change="modeChange" v-model="ipRulesForm.mode" size="small">
                    <el-radio-button label="关闭" :value="0" />
                    <el-radio-button label="黑名单" :value="1" />
                    <el-radio-button label="白名单" :value="2" />
                </el-radio-group>
            </el-form-item>
            <el-form-item label="操作:">
                <el-button type="warning" :icon="Plus" @click="showAddIPRuleDialog" size="small">添加记录</el-button>
                <el-button type="primary" :loading="dataInLoading" :icon="Refresh" @click="reloadIPRules"
                    size="small">重载配置</el-button>
            </el-form-item>
        </el-form>
        <IPRulesList name="黑名单" :in-loading="dataInLoading" type="block" @on-dirty="onDirty" :data="ipRulesForm.block"
            :now-work="ipRulesForm.mode == 1" />
        <IPRulesList name="白名单" :in-loading="dataInLoading" type="allow" @on-dirty="onDirty" :data="ipRulesForm.allow"
            :now-work="ipRulesForm.mode == 2" />
    </el-dialog>
    <AddIPRule ref="addIPRuleDialog" @on-dirty="refreshIPRules" />
</template>
<style scoped lang="scss">
.actionForm {
    align-items: center;
    text-align: center;
    display: flex;
    justify-content: space-around;
}
</style>