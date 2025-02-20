<script setup lang="ts">
import { usePanelStore } from '../stores/panel';
const panelStore = usePanelStore();
</script>
<template>
    <div class="versionInfo">T2T Server {{ panelStore.version }} © StupidBear Studio {{ new Date().getFullYear() }}
    </div>
    <div class="systemInfo">
        hash [<a v-if="panelStore.commitId != '' && panelStore.commitId != 'unknown'" target="_blank"
            :href="`https://github.com/XljBear/T2T/commit/${panelStore.commitId}`">{{ panelStore.commitId
            }}</a><span v-else>unknown</span>] run on {{ panelStore.arch }} {{ panelStore.os }}
    </div>
</template>
<style scoped lang="scss">
a {
    color: #999;
    font-size: 10px;
    text-decoration: none;

    &:hover {
        color: #ff9012 !important;
    }

    &:visited,
    &:active {
        color: #999;
    }
}

.versionInfo {
    color: #ccc;
    font-size: 12px;
}

.systemInfo {
    color: #999;
    font-size: 10px;
}
</style>