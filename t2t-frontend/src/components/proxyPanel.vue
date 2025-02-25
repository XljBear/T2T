<script setup lang="ts">
import { Plus, Refresh, Cherry, Setting, SwitchButton } from '@element-plus/icons-vue';
import { onMounted, ref, onUnmounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import TrafficChart from './trafficChart.vue';
import LinksDialog from './links.vue';
import SettingDialog from './setting.vue';
import IPRulesDialog from './ipRules.vue';
import { axiosInstance } from '../utils/axios';
import { useLoginStore } from '../stores/login';
import FooterInfo from './footer.vue';
const loginStore = useLoginStore();
const proxyData = ref<any>([]);
const showProxyForm = ref(false);
const proxyTrafficChartRefs = ref<any>({});
const proxyLoading = ref(false);
const proxyForm = ref({
    uuid: '',
    name: '',
    local_address: '',
    remote_address: '',
    max_link: 0,
    status: true,
});
const linksPage = ref();
const settingPage = ref();
const ipRulesPage = ref();
onMounted(() => {
    getProxyList();
});
onUnmounted(() => {
    if (refreshProxyTrafficTimer.value != 0) {
        clearInterval(refreshProxyTrafficTimer.value);
    }
});
const createProxyForm = () => {
    showProxyForm.value = true;
    proxyForm.value = {
        uuid: '',
        name: '',
        local_address: '',
        remote_address: '',
        max_link: 0,
        status: true,
    };
}
const createProxy = () => {
    if (!proxyForm.value.name || !proxyForm.value.local_address || !proxyForm.value.remote_address || proxyForm.value.max_link < 0) {
        ElMessage({
            message: '请填写完整正确的配置信息',
            type: 'warning',
        })
        return;
    }
    ElMessageBox.confirm(
        '确认创建 [' + proxyForm.value.name + '] 反代理配置?',
        '创建反代理配置',
        {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
        }
    )
        .then(() => {
            axiosInstance.post("/proxy", proxyForm.value).then(() => {
                ElMessage({
                    message: '创建配置成功',
                    type: 'success',
                });
            }).catch((resp) => {
                ElMessage({
                    message: '创建配置失败:' + resp.response.data.error,
                    type: 'error',
                });
            }).finally(() => {
                showProxyForm.value = false;
                proxyForm.value = {
                    uuid: '',
                    name: '',
                    local_address: '',
                    remote_address: '',
                    max_link: 0,
                    status: true,
                };
                getProxyList();
            });
        }).catch(() => {
        });
}
const deleteProxy = (index: number, uuid: string) => {
    ElMessageBox.confirm(
        '确认删除 [' + proxyData.value[index]['name'] + '] 反代理配置吗？',
        {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
        }
    ).then(() => {
        axiosInstance.delete("/proxy/" + uuid).then(() => {
            ElMessage({
                message: '删除配置成功',
                type: 'success',
            });
        }).catch((resp) => {
            ElMessage({
                message: '删除配置失败:' + resp.response.data.error,
                type: 'error',
            });
        }).finally(() => {
            getProxyList();
        });
    }).catch(() => {
    });
}
const refreshProxyTrafficTimer = ref<number>(0);
const getProxyList = () => {
    proxyLoading.value = true;
    if (refreshProxyTrafficTimer.value != 0) {
        clearInterval(refreshProxyTrafficTimer.value);
    }
    proxyTrafficChartRefs.value = new Object;
    axiosInstance.get("/proxy").then((response) => {
        proxyData.value = response.data;
        if (proxyData.value.length > 0) {
            refreshProxyTrafficTimer.value = setInterval(refreshProxyTrafficData, 1500);
        }
    }).catch(() => {
    }).finally(() => {
        proxyLoading.value = false;
    });
}
const editProxy = (index: number) => {
    showProxyForm.value = true;
    proxyForm.value.uuid = proxyData.value[index].uuid;
    proxyForm.value.local_address = proxyData.value[index].local_address;
    proxyForm.value.remote_address = proxyData.value[index].remote_address;
    proxyForm.value.name = proxyData.value[index].name;
    proxyForm.value.max_link = proxyData.value[index].max_link
    proxyForm.value.status = proxyData.value[index].status;
}
const updateProxy = () => {
    if (!proxyForm.value.name || !proxyForm.value.local_address || !proxyForm.value.remote_address) {
        ElMessage({
            message: '请填写完整配置信息',
            type: 'warning',
        })
        return;
    }
    ElMessageBox.confirm(
        '确认更新该反代理配置吗？',
        {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
        }
    )
        .then(() => {
            axiosInstance.put("/proxy/" + proxyForm.value.uuid, proxyForm.value).then(() => {
                ElMessage({
                    message: '更新配置成功',
                    type: 'success',
                });
            }).catch((resp) => {
                ElMessage({
                    message: '更新配置失败:' + resp.response.data.error,
                    type: 'error',
                });
            }).finally(() => {
                showProxyForm.value = false;
                proxyForm.value = {
                    uuid: '',
                    name: '',
                    local_address: '',
                    remote_address: '',
                    max_link: 0,
                    status: true,
                };
                getProxyList();
            });
        })
        .catch(() => {
        });
}
const restartService = () => {
    ElMessageBox.confirm(
        '确定要立即重启服务吗？',
        {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
        }
    ).then(() => {
        axiosInstance.post("/restart").then(() => {
            getProxyList();
            ElMessage({
                message: '服务重启成功',
                type: 'success',
            });
        }).catch(() => {
            ElMessage({
                message: '服务重启失败',
                type: 'error',
            });
        });
    }).catch(() => {
    });
}
const proxyTrafficData = ref<any>({});
const refreshProxyTrafficData = () => {
    axiosInstance.get("/traffic").then((response) => {
        proxyTrafficData.value = {};
        Object.entries(response.data).forEach(([key, data]) => {
            if (proxyTrafficChartRefs.value[key]) {
                const trafficData: any = data as object
                proxyTrafficChartRefs.value[key].pushTrafficData(trafficData.downlink_in_second, trafficData.uplink_in_second, trafficData.downlink_total, trafficData.uplink_total);
            }
            proxyTrafficData.value[key] = data;
        });
    }).catch(() => {
    });
}
const showLinks = (uuid: string, name: string, maxLink: number, localAddress: string) => {
    const localPort = localAddress.split(":")
    linksPage.value.showLinksPage(uuid, name, maxLink, localPort[localPort.length - 1]);
}
const showSetting = () => {
    settingPage.value.showSettingDialog();
}
const logout = () => {
    ElMessageBox.confirm(
        '确定要退出登录吗？',
        {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
        }
    ).then(() => {
        axiosInstance.post("/logout").then(() => {
        }).catch(() => {
        }).finally(() => {
            loginStore.logout();
        });
    }).catch(() => {
    });
}
const showIPRules = () => {
    ipRulesPage.value.showIPRulesDialog();
}
</script>

<template>
    <div class="proxyPanel">
        <h1>T2T Service Panel</h1>
        <el-card class="panel">
            <template #header>
                <div class="cardHeader">
                    <el-button :icon="Plus" type="primary" @click="createProxyForm">创建反代理配置</el-button>
                    <el-button :icon="Cherry" type="danger" @click="showIPRules">防火墙配置</el-button>
                </div>
            </template>
            <el-table v-loading="proxyLoading" :data="proxyData" stripe style="width: 100%">
                <el-table-column fixed="left" prop="name" label="名称" />
                <el-table-column prop="local_address" label="本地端口" min-width="150" />
                <el-table-column prop="remote_address" label="远程端口" min-width="150" />
                <el-table-column label="连接数" min-width="100">
                    <template #default="scope">
                        <el-link
                            @click="showLinks(scope.row.uuid, scope.row.name, scope.row.max_link, scope.row.local_address)"
                            type="warning">
                            {{ proxyTrafficData[scope.row.uuid] ? proxyTrafficData[scope.row.uuid].link_count : 0
                            }}</el-link> / {{
                                scope.row.max_link
                                    == 0 ? '无限制' : scope.row.max_link }}
                    </template>
                </el-table-column>
                <el-table-column label="网络" min-width="200">
                    <template #default="scope">
                        <TrafficChart :key="scope.row.uuid" :ref="el => proxyTrafficChartRefs[scope.row.uuid] = el"
                            class="tChart" />
                    </template>
                </el-table-column>
                <el-table-column prop="status" label="状态">
                    <template #default="scope">
                        <el-tag type="success" v-if="scope.row.status">生效</el-tag>
                        <el-tag type="warning" v-else>禁用</el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="操作" min-width="120">
                    <template #default="scope">
                        <el-button type="primary" link @click="editProxy(scope.$index)">编辑</el-button>
                        <el-button type="danger" link @click="deleteProxy(scope.$index, scope.row.uuid)">删除</el-button>
                    </template>
                </el-table-column>
                <template #empty>
                    <div style="text-align: center;">
                        无反向代理配置
                    </div>
                </template>
            </el-table>
            <template #footer>
                <div class="systemAction">
                    <el-button :icon="Setting" type="info" link @click="showSetting">系统设置</el-button>
                    <el-button :icon="SwitchButton" type="warning" link @click="logout">退出登录</el-button>
                    <el-button :icon="Refresh" type="danger" link @click="restartService">重启服务</el-button>
                </div>

                <FooterInfo />
            </template>
        </el-card>
        <el-dialog v-model="showProxyForm" :title="proxyForm.uuid != '' ? '编辑反代理配置' : '创建反代理配置'" width="100%"
            style="max-width: 500px;">
            <el-form :model="proxyForm">
                <el-form-item label="业务名称 :" :label-width="80">
                    <el-input v-model="proxyForm.name" autocomplete="off" />
                </el-form-item>
                <el-form-item label="本地端口 :" :label-width="80">
                    <el-input v-model="proxyForm.local_address" autocomplete="off" />
                </el-form-item>
                <el-form-item label="远程端口 :" :label-width="80">
                    <el-input v-model="proxyForm.remote_address" autocomplete="off" />
                </el-form-item>
                <el-form-item label="最大连接数 :" :label-width="90">
                    <el-input v-model.number="proxyForm.max_link" autocomplete="off" />
                </el-form-item>
                <el-form-item label="启用 :" :label-width="80">
                    <el-checkbox v-model="proxyForm.status" />
                </el-form-item>
            </el-form>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="showProxyForm = false">取消</el-button>
                    <el-button type="primary" @click="proxyForm.uuid != '' ? updateProxy() : createProxy()">
                        {{ proxyForm.uuid != '' ? '更新' : '创建' }}
                    </el-button>
                </div>
            </template>
        </el-dialog>
        <LinksDialog ref="linksPage" />
        <SettingDialog ref="settingPage" />
        <IPRulesDialog ref="ipRulesPage" />
    </div>
</template>

<style scoped lang="scss">
.proxyPanel {
    width: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;

    .panel {
        width: 80%;
    }

    .cardHeader {
        display: flex;
        justify-content: space-between;
    }

    .systemAction {
        margin-bottom: 5px;
    }

    .tChart {
        width: 200px;
        height: 100px;
    }
}

@media screen and (orientation: portrait) and (max-width: 768px) {
    .proxyPanel {
        h1 {
            font-size: 24px;
        }

        .panel {
            width: 95%;
        }
    }
}
</style>
