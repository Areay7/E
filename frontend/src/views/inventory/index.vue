<template>
  <div class="inventory-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">库存管理</span>
          <el-button type="primary" @click="loadInventory">刷新</el-button>
        </div>
      </template>

      <div class="filter-shell">
        <el-form :inline="true" class="filters">
          <el-form-item label="平台">
            <el-select v-model="platform" placeholder="选择平台" @change="loadInventory">
              <el-option label="全部" value="" />
              <el-option label="Shopee" value="shopee" />
              <el-option label="速卖通" value="aliexpress" />
              <el-option label="TikTok" value="tiktok" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadInventory">查询</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table v-loading="loading" :data="inventory" stripe border>
        <el-table-column prop="platform" label="平台" width="120" />
        <el-table-column prop="sku" label="SKU" width="150" />
        <el-table-column prop="productId" label="商品 ID" width="150" />
        <el-table-column prop="stock" label="总库存" width="100" />
        <el-table-column prop="reservedStock" label="预留库存" width="100" />
        <el-table-column prop="availableStock" label="可用库存" width="100" />
        <el-table-column prop="warehouseCode" label="仓库" width="120" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleUpdate(row)">更新库存</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pager">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          background
          @current-change="loadInventory"
          @size-change="onSizeChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" title="更新库存" width="500px">
      <el-form :model="updateForm" label-width="100px">
        <el-form-item label="SKU">
          <el-input v-model="updateForm.sku" disabled />
        </el-form-item>
        <el-form-item label="平台">
          <el-input v-model="updateForm.platform" disabled />
        </el-form-item>
        <el-form-item label="新库存数量">
          <el-input-number v-model="updateForm.quantity" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitUpdate">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { http } from '@/api/http'

const loading = ref(false)
const platform = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const inventory = ref<any[]>([])
const dialogVisible = ref(false)
const updateForm = ref({
  platform: '',
  sku: '',
  quantity: 0
})

async function loadInventory() {
  loading.value = true
  try {
    const { data } = await http.get('/api/v1/inventory', {
      params: {
        platform: platform.value || undefined,
        page: page.value,
        pageSize: pageSize.value
      }
    })
    inventory.value = data.list
    total.value = data.total
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

function onSizeChange() {
  page.value = 1
  loadInventory()
}

function handleUpdate(row: any) {
  updateForm.value = {
    platform: row.platform,
    sku: row.sku,
    quantity: row.stock
  }
  dialogVisible.value = true
}

async function submitUpdate() {
  try {
    await http.put('/api/v1/inventory', updateForm.value)
    ElMessage.success('更新成功')
    dialogVisible.value = false
    loadInventory()
  } catch (e) {
    ElMessage.error('更新失败')
  }
}

onMounted(() => {
  loadInventory()
})
</script>

<style scoped>
.inventory-page {
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

.filter-shell {
  padding: 16px 18px;
  margin-bottom: 16px;
  background: linear-gradient(180deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 8px;
  border: 1px solid rgba(15, 23, 42, 0.06);
}

.filters {
  margin-bottom: 0;
}

.pager {
  margin-top: 18px;
  display: flex;
  justify-content: flex-end;
}
</style>
