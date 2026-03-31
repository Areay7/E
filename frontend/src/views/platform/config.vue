<template>
  <div class="platform-config">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">平台配置管理</span>
          <el-button type="primary" @click="loadConfigs">刷新</el-button>
        </div>
      </template>

      <el-tabs v-model="activePlatform">
        <el-tab-pane label="Shopee" name="shopee">
          <PlatformConfigForm platform="shopee" :config="configs.shopee" @save="handleSave" />
        </el-tab-pane>
        <el-tab-pane label="速卖通" name="aliexpress">
          <PlatformConfigForm platform="aliexpress" :config="configs.aliexpress" @save="handleSave" />
        </el-tab-pane>
        <el-tab-pane label="TikTok" name="tiktok">
          <PlatformConfigForm platform="tiktok" :config="configs.tiktok" @save="handleSave" />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-card shadow="never" style="margin-top: 20px">
      <template #header>
        <span class="title">数据同步</span>
      </template>
      <div class="sync-actions">
        <el-button type="primary" @click="syncOrders(activePlatform)">同步订单</el-button>
        <el-button type="success" @click="syncProducts(activePlatform)">同步商品</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import PlatformConfigForm from '@/components/platform/PlatformConfigForm.vue'
import { http } from '@/api/http'

const activePlatform = ref('shopee')
const configs = ref<any>({
  shopee: null,
  aliexpress: null,
  tiktok: null
})

async function loadConfigs() {
  try {
    const { data } = await http.get('/api/v1/platforms')
    data.forEach((config: any) => {
      configs.value[config.platform] = config
    })
  } catch (e) {
    ElMessage.error('加载配置失败')
  }
}

async function handleSave(platform: string, config: any) {
  try {
    await http.put(`/api/v1/platforms/${platform}/config`, config)
    ElMessage.success('保存成功')
    loadConfigs()
  } catch (e) {
    ElMessage.error('保存失败')
  }
}

async function syncOrders(platform: string) {
  try {
    await http.post(`/api/v1/platforms/${platform}/sync/orders`)
    ElMessage.success('同步任务已启动')
  } catch (e) {
    ElMessage.error('启动失败')
  }
}

async function syncProducts(platform: string) {
  try {
    await http.post(`/api/v1/platforms/${platform}/sync/products`)
    ElMessage.success('同步任务已启动')
  } catch (e) {
    ElMessage.error('启动失败')
  }
}

onMounted(() => {
  loadConfigs()
})
</script>

<style scoped>
.platform-config {
  max-width: 1280px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 17px;
  font-weight: 700;
}

.sync-actions {
  display: flex;
  gap: 12px;
}
</style>
