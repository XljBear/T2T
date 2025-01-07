<script setup lang="ts">
import moment from 'moment';
import { Delete, WarningFilled } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { axiosInstance } from '../utils/axios';
const emit = defineEmits(["onDirty"]);
const props = defineProps({
    name: {
        type: String,
        default: ""
    },
    data: {
        type: Array,
        default: []
    },
    nowWork: {
        type: Boolean,
        default: false
    },
    inLoading: {
        type: Boolean,
        default: false
    },
    type: {
        type: String,
        default: ""
    }
});
const deleteRule = (uuid: string) => {
    ElMessageBox.confirm(
        '确认要删除该规则吗?',
        '删除防火墙规则',
        {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
        }
    ).then(() => {
        axiosInstance.delete(`/ipRules/${props.type == 'allow' ? 'allow' : 'block'}/${uuid}`).then(() => {
            ElMessage({
                message: '防火墙规则删除成功',
                type: 'success',
            });
        }).catch((resp) => {
            ElMessage({
                message: '防火墙规则删除失败:' + resp.response.data.error,
                type: 'error',
            });
        }).finally(() => {
            emit("onDirty");
        });
    })
};
</script>
<template>
    <div :class="['list', props.nowWork ? 'work' : '']">
        <el-divider border-style="dashed">
            {{ props.name }} ({{ props.data.length }})
        </el-divider>
        <el-table v-loading="props.inLoading" :data="props.data" width="100%">
            <el-table-column fixed="left" prop="ip" label="IP地址" min-width="160" />
            <el-table-column label="端口" prop="port" min-width="100">
                <template #default="scope">
                    <span v-if="scope.row.port == ''">所有</span>
                </template>
            </el-table-column>
            <el-table-column label="添加时间" min-width="160">
                <template #default="scope">
                    {{ moment(scope.row.start_time).format("YYYY-MM-DD HH:mm:ss") }}
                </template>
            </el-table-column>
            <el-table-column label="失效时间" min-width="160">
                <template #default="scope">
                    {{ scope.row.end_time ? moment(scope.row.end_time).format("YYYY-MM-DD HH:mm:ss") : '无期限' }}
                </template>
            </el-table-column>
            <el-table-column label="原由">
                <template #default="scope">
                    <span v-if="scope.row.reason == ''">无</span>
                    <el-tooltip v-else :content="scope.row.reason">
                        <el-icon size="large">
                            <WarningFilled />
                        </el-icon>
                    </el-tooltip>
                </template>
            </el-table-column>
            <el-table-column fixed="right" label="操作" min-width="100">
                <template #default="scope">
                    <el-button :icon="Delete" @click="deleteRule(scope.row.uuid)" />
                </template>
            </el-table-column>
            <template #empty>
                无记录
            </template>
        </el-table>
    </div>
</template>
<style scoped lang="scss">
.list {
    transition: all linear .2s;
    margin-top: 30px;
    opacity: 0.5;

    &.work {
        opacity: 1;
    }
}
</style>