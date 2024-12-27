<script setup lang="ts">
import * as echarts from 'echarts';
import { ref, onMounted, nextTick } from 'vue';
import { Download, Upload, Box } from '@element-plus/icons-vue';
const chart = ref(null);
const chartInstance = ref<echarts.ECharts>();
const downlinkData = ref<number[]>([]);
const uplinkData = ref<number[]>([]);
const cacheDataLimit = 10;
const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
}
nextTick(() => {
    redrawChart();
});
const redrawChart = () => {
    if (chartInstance.value) {
        chartInstance.value.setOption({
            color: ['#03c2df', '#FF0087'],
            xAxis: {
                type: 'category',
                show: false
            },
            yAxis: {
                type: 'value',
                show: false
            },
            series: [
                {
                    name: '上行',
                    data: uplinkData.value,
                    type: 'line',
                    smooth: true,

                },
                {
                    name: '下行',
                    data: downlinkData.value,
                    type: 'line',
                    smooth: true,
                }
            ]
        });
    }
}
onMounted(() => {
    chartInstance.value = echarts.init(chart.value);
    clearTrafficData();
})
const downlinkVal = ref(0);
const uplinkVal = ref(0);
const downlinkTotalVal = ref(0);
const uplinkValTotalVal = ref(0);
const pushTrafficData = (downlink: number, uplink: number, downlinkTotal: number, uplinkTotal: number) => {
    downlinkVal.value = downlink;
    uplinkVal.value = uplink;
    downlinkTotalVal.value = downlinkTotal;
    uplinkValTotalVal.value = uplinkTotal;
    if (downlinkData.value.length >= cacheDataLimit) {
        downlinkData.value.shift();
        uplinkData.value.shift();
    }
    downlinkData.value.push(downlink);
    uplinkData.value.push(uplink);
    redrawChart();
}
const clearTrafficData = () => {
    downlinkData.value = [];
    uplinkData.value = [];
    for (let i = 0; i < cacheDataLimit; i++) {
        pushTrafficData(0, 0, 0, 0);
    }
}
defineExpose({ pushTrafficData, clearTrafficData });
</script>
<template>
    <div class="traffic">
        <div ref="chart" class="chart"></div>
        <div class="info now">
            <div class="val up"><el-icon>
                    <Upload />
                </el-icon> {{ formatBytes(uplinkVal) }}/s</div>
            <div class="val down"><el-icon>
                    <Download />
                </el-icon> {{ formatBytes(downlinkVal) }}/s</div>
        </div>
        <div class="info">
            <div class="val up"><el-icon>
                    <Box />
                </el-icon> {{ formatBytes(uplinkValTotalVal) }}</div>
            <div class="val down"><el-icon>
                    <Box />
                </el-icon> {{ formatBytes(downlinkTotalVal) }}</div>
        </div>
    </div>

</template>
<style scoped lang="scss">
.traffic {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    justify-items: center;

    .info {
        display: flex;
        justify-content: center;
        height: 15px;

        .val {
            margin-right: 10px;
            font-size: 10px;
            font-weight: bold;
            color: #999999;
        }

        &.now {
            .val {
                font-size: 10px;
                font-weight: normal;

                &.up {
                    color: #03c2df;
                }

                &.down {
                    color: #FF0087;
                }
            }
        }
    }

    .chart {
        width: 100%;
        height: calc(100% - 20px);
        margin-bottom: -18px;
    }
}
</style>