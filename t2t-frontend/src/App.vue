<script setup lang="ts">
import axios from 'axios';
import { Plus, Refresh } from '@element-plus/icons-vue';
import { onMounted, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import TrafficChart from './components/trafficChart.vue';
let devMode = true;
if (import.meta.env.MODE === 'production') {
  devMode = false;
}
const axiosInstance = axios.create({
  baseURL: devMode ? 'http://127.0.0.1:8080/api' : '../api',
  timeout: 3000,
});
const proxyData = ref<any>([]);
const showProxyForm = ref(false);
const proxyTrafficChartRefs = ref<any>({});
const proxyForm = ref({
  uuid: '',
  name: '',
  local_address: '',
  remote_address: '',
  max_link: 0,
  status: true,
});
onMounted(() => {
  getProxyList();
})
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
    })
    .catch(() => {
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
  if (refreshProxyTrafficTimer.value != 0) {
    clearInterval(refreshProxyTrafficTimer.value);
  }
  proxyTrafficChartRefs.value = new Object;
  axiosInstance.get("/proxy").then((response) => {
    proxyData.value = response.data;
    if (proxyData.value.length > 0) {
      refreshProxyTrafficTimer.value = setInterval(refreshProxyTrafficData, 1500);
    }
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
        proxyTrafficChartRefs.value[key].pushTrafficData(trafficData.downlink_in_second, trafficData.uplink_in_second,trafficData.downlink_total,trafficData.uplink_total);
      }
      proxyTrafficData.value[key] = data;
    });
  });
}
</script>

<template>
  <div class="app">
    <h1>T2T Service Panel</h1>
    <el-card style="width: 80%;">
      <template #header>
        <div class="card-header">
          <el-button :icon="Plus" type="primary" @click="createProxyForm">创建反代理配置</el-button>
          <el-button :icon="Refresh" type="danger" @click="restartService">重启服务</el-button>
        </div>
      </template>
      <el-table :data="proxyData" stripe style="width: 100%">
        <el-table-column fixed="left" prop="name" label="名称" />
        <el-table-column prop="local_address" label="本地端口" min-width="150" />
        <el-table-column prop="remote_address" label="远程端口" min-width="150" />
        <el-table-column label="连接数" min-width="100">
          <template #default="scope">
            {{ proxyTrafficData[scope.row.uuid] ? proxyTrafficData[scope.row.uuid].link_count : 0 }} / {{ scope.row.max_link
              == 0 ? '无限制' : scope.row.max_link }}
          </template>
        </el-table-column>
        <el-table-column label="网络" min-width="150">
          <template #default="scope">
            <TrafficChart :ref="el => proxyTrafficChartRefs[scope.row.uuid] = el" class="tChart" />
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.status">生效</el-tag>
            <el-tag type="warning" v-else>禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" min-width="120">
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
        <div class="card-footer">T2T Server v0.0.1 © StupidBear Studio 2024</div>
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
  </div>
</template>

<style scoped lang="scss">
.app {
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  .card-header {
    display: flex;
    justify-content: space-between;
  }

  .card-footer {
    color: #ccc;
    font-size: 12px;
  }

  .tChart {
    width: 150px;
    height: 100px;
  }
}
</style>
