<script setup lang="ts">
import { ref, onMounted, watch, onUnmounted } from 'vue';
import moment from 'moment';
import TrafficChart from './trafficChart.vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { axiosInstance } from '../utils/axios';
const proxyTrafficChartRefs = ref<any>({});
onMounted(() => {
    moment.defineLocale('zh-cn', {
        relativeTime: {
            future: '%s内',
            past: '%s前',
            s: '几秒',
            m: '1分钟',
            mm: '%d分钟',
            h: '1小时',
            hh: '%d小时',
            d: '1天',
            dd: '%d天',
            M: '1个月',
            MM: '%d个月',
            y: '1年',
            yy: '%d年'
        },
    });
});
onUnmounted(() => {
    if (linksDataRefreshTimer.value != 0) {
        clearInterval(linksDataRefreshTimer.value);
    }
});
const dialogTableVisible = ref(false);
const proxyUUID = ref('');
const proxyName = ref('');
const proxyMaxLink = ref(0);
const linksDataRefreshTimer = ref<number>(0);
const linksLoading = ref(false);
const showLinksPage = (uuid: string, name: string, maxLink: number) => {
    proxyLinksData.value = [];
    proxyMaxLink.value = maxLink;
    proxyUUID.value = uuid;
    proxyName.value = name;
    dialogTableVisible.value = true;
    linksLoading.value = true;
    refreshLinks();
    linksDataRefreshTimer.value = setInterval(refreshLinks, 1500);
}
watch(() => dialogTableVisible.value, (show: boolean) => {
    if (!show) {
        clearInterval(linksDataRefreshTimer.value);
    }
});
const proxyLinksData = ref<any>([]);
const refreshLinks = () => {
    if (proxyUUID.value == "") {
        return;
    }
    axiosInstance.get(`/proxy/${proxyUUID.value}/links`).then((response) => {
        proxyLinksData.value = [];
        proxyLinksData.value = response.data;
        proxyLinksData.value.forEach((data: any) => {
            if (proxyTrafficChartRefs.value[data.uuid]) {
                const trafficData: any = data.traffic;
                proxyTrafficChartRefs.value[data.uuid].pushTrafficData(trafficData.downlink_in_second, trafficData.uplink_in_second, trafficData.downlink_total, trafficData.uplink_total);
            }
        });
    }).catch(() => {
    }).finally(() => {
        if (linksLoading.value) {
            linksLoading.value = false;
        };
    });
}
const kickLink = (uuid: string) => {
    ElMessageBox.confirm(
        '确定要立即断开此连接吗？',
        {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
        }
    ).then(() => {
        axiosInstance.delete(`/proxy/${proxyUUID.value}/links/${uuid}`).then(() => {
            ElMessage({
                message: '连接已成功断开',
                type: 'success',
            });
        }).catch((resp) => {
            ElMessage({
                message: '断开连接失败:' + resp.response.data.error,
                type: 'error',
            });
        });
    }).catch(() => {
    });
}
defineExpose({ showLinksPage })
</script>
<template>
    <el-dialog draggable v-model="dialogTableVisible"
        :title="`${proxyName} 连接池 (${proxyLinksData.length}/${proxyMaxLink == 0 ? '无限制' : proxyMaxLink})`" width="100%"
        style="max-width:800px">
        <el-table v-loading="linksLoading" :data="proxyLinksData" width="100%">
            <el-table-column fixed="left" prop="ip" label="IP地址" />
            <el-table-column label="连接时长" min-width="120">
                <template #default="scope">
                    <div>
                        {{ moment(scope.row.link_time).fromNow() }}
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="网络" min-width="200">
                <template #default="scope">
                    <TrafficChart :key="scope.row.uuid" :ref="el => proxyTrafficChartRefs[scope.row.uuid] = el"
                        class="tChart" />
                </template>
            </el-table-column>
            <el-table-column fixed="right" label="操作" min-width="60">
                <template #default="scope">
                    <el-button type="danger" link @click="kickLink(scope.row.uuid)">断开</el-button>
                </template>
            </el-table-column>
            <template #empty>
                <div style="text-align: center;">
                    当前暂无连接
                </div>
            </template>
        </el-table>
    </el-dialog>
</template>
<style scoped lang="scss">
.tChart {
    width: 200px;
    height: 100px;
}
</style>