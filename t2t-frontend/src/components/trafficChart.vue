<script setup lang="ts">
import * as echarts from 'echarts';
import { ref, onMounted, nextTick } from 'vue';
const chart = ref(null);
const chartInstance = ref<echarts.ECharts>();
const downlinkData = ref<number[]>([]);
const uplinkData = ref<number[]>([]);
const cacheDataLimit = 10;
const formatBytes = (params: any): string => {
    const bytes = params.value
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${params.seriesName} ${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
}
nextTick(() => {
    redrawChart();
});
const redrawChart = () => {
    if (chartInstance.value) {
        chartInstance.value.setOption({
            color: ['#00DDFF', '#FF0087'],
            xAxis: {
                type: 'category',
                show: false
            },
            yAxis: {
                type: 'value',
                show: false
            },
            tooltip: {
                trigger: 'item',
                axisPointer: {
                    type: 'shadow'
                },
                confine: true,
                formatter: formatBytes
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
const pushTrafficData = (downlink: number, uplink: number) => {
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
        pushTrafficData(0, 0);
    }
}
defineExpose({ pushTrafficData, clearTrafficData });
</script>
<template>
    <div ref="chart" class="chart"></div>
</template>
<style scoped lang="scss">
.chart {
    width: 100%;
    height: 100%;
}
</style>