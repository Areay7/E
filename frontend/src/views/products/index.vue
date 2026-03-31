<template>
  <div class="products-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">商品管理</span>
          <el-button type="primary" @click="loadProducts">刷新</el-button>
        </div>
      </template>

      <div class="filter-shell">
        <el-form :inline="true" class="filters">
          <el-form-item label="平台">
            <el-select v-model="platform" placeholder="选择平台" @change="loadProducts">
              <el-option label="Shopee" value="shopee" />
              <el-option label="速卖通" value="aliexpress" />
              <el-option label="TikTok" value="tiktok" />
            </el-select>
          </el-form-item>
          <el-form-item label="关键词">
            <el-input v-model="keyword" placeholder="商品名称/SKU" clearable @clear="loadProducts" />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="status" clearable placeholder="全部" @change="loadProducts">
              <el-option label="上架" value="active" />
              <el-option label="下架" value="inactive" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadProducts">查询</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table v-loading="loading" :data="products" stripe border>
        <el-table-column prop="productId" label="商品 ID" width="150" />
        <el-table-column prop="name" label="商品名称" min-width="200" />
        <el-table-column prop="sku" label="SKU" width="120" />
        <el-table-column label="价格" width="120">
          <template #default="{ row }">{{ row.price }} {{ row.currency }}</template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="100" />
        <el-table-column prop="status" label="状态" width="100" />
        <el-table-column prop="soldCount" label="销量" width="100" />
      </el-table>

      <div class="pager">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          background
          @current-change="loadProducts"
          @size-change="onSizeChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { http } from '@/api/http'

const loading = ref(false)
const platform = ref('shopee')
const keyword = ref('')
const status = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const products = ref<any[]>([])

async function loadProducts() {
  loading.value = true
  try {
    const { data } = await http.get('/api/v1/products', {
      params: {
        platform: platform.value,
        page: page.value,
        pageSize: pageSize.value,
        keyword: keyword.value || undefined,
        status: status.value || undefined
      }
    })
    products.value = data.list
    total.value = data.total
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

function onSizeChange() {
  page.value = 1
  loadProducts()
}

onMounted(() => {
  loadProducts()
})
</script>

<style scoped>
.products-page {
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
