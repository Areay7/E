<template>
  <div class="order-list-view">
    <el-card shadow="never" class="page-card">
      <template #header>
        <div class="card-header">
          <div class="title-block">
            <span class="title-accent" aria-hidden="true" />
            <span class="title">{{ pageTitle }}</span>
          </div>
          <el-button type="primary" :loading="loading" @click="load">刷新</el-button>
        </div>
      </template>

      <div class="filter-shell">
        <el-form :inline="true" class="filters" @submit.prevent>
          <el-form-item label="订单号 / 关键词">
            <el-input
              v-model="keyword"
              clearable
              placeholder="支持后端实现后筛选"
              class="filter-input"
              @clear="load"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select
              v-model="status"
              clearable
              placeholder="全部"
              class="filter-select"
              @change="load"
            >
              <el-option v-for="opt in statusOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="loading" @click="load">查询</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-alert
        v-if="lastError"
        type="error"
        :title="lastError"
        show-icon
        class="alert"
        :closable="false"
      />

      <div class="table-wrap">
        <el-table
          v-loading="loading"
          :data="tableData"
          stripe
          border
          empty-text="暂无订单数据"
          class="data-table"
        >
          <el-table-column prop="orderId" label="订单号" min-width="140" />
          <el-table-column label="买家" min-width="120">
            <template #default="{ row }">{{ row.buyer || '—' }}</template>
          </el-table-column>
          <el-table-column label="金额" min-width="110">
            <template #default="{ row }">{{ formatAmount(row) }}</template>
          </el-table-column>
          <el-table-column label="状态" min-width="100">
            <template #default="{ row }">{{ row.status || '—' }}</template>
          </el-table-column>
          <el-table-column label="下单时间" min-width="170">
            <template #default="{ row }">{{ row.createdAt || '—' }}</template>
          </el-table-column>
        </el-table>
      </div>

      <div class="pager">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          background
          @current-change="load"
          @size-change="onSizeChange"
        />
      </div>

      <p class="hint">
        接口约定：<code>GET /api/v1/orders</code>，返回 <code>{ list, total }</code>；开发环境可在
        <code>vite.config.ts</code> 中配置 proxy 指向后端。
      </p>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { fetchOrders } from '@/api/orders'
import type { OrderPlatform, OrderRow } from '@/types/order'
import { messageFromError } from '@/utils/error'

const props = defineProps<{
  platform: OrderPlatform
  pageTitle: string
}>()

const statusOptions = [
  { label: '全部', value: '' },
  { label: '待发货', value: 'pending_shipment' },
  { label: '已发货', value: 'shipped' },
  { label: '已完成', value: 'completed' },
]

const loading = ref(false)
const tableData = ref<OrderRow[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')
const status = ref('')
const lastError = ref<string | null>(null)

function formatAmount(row: OrderRow): string {
  if (row.amount === undefined || row.amount === '') return '—'
  const cur = row.currency ? ` ${row.currency}` : ''
  return `${row.amount}${cur}`
}

async function load() {
  loading.value = true
  lastError.value = null
  try {
    const res = await fetchOrders(props.platform, {
      page: page.value,
      pageSize: pageSize.value,
      keyword: keyword.value.trim() || undefined,
      status: status.value || undefined,
    })
    tableData.value = res.list
    total.value = res.total
  } catch (e) {
    lastError.value = messageFromError(e)
    tableData.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function onSizeChange() {
  page.value = 1
  void load()
}

onMounted(() => {
  void load()
})

watch(
  () => props.platform,
  () => {
    page.value = 1
    void load()
  }
)
</script>

<style scoped>
.order-list-view {
  max-width: 1280px;
}

.page-card :deep(.el-card__header) {
  padding: 18px 22px;
  border-bottom: 1px solid rgba(15, 23, 42, 0.06);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.title-block {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.title-accent {
  width: 4px;
  height: 22px;
  border-radius: 4px;
  background: linear-gradient(180deg, var(--cb-primary-light), var(--cb-primary));
  flex-shrink: 0;
}

.title {
  font-size: 17px;
  font-weight: 700;
  color: var(--cb-text);
  letter-spacing: -0.02em;
}

.filter-shell {
  padding: 16px 18px;
  margin-bottom: 16px;
  background: linear-gradient(180deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: var(--cb-radius-sm);
  border: 1px solid rgba(15, 23, 42, 0.06);
}

.filters {
  margin-bottom: 0;
}

.filter-input {
  width: 220px;
}

.filter-select {
  width: 160px;
}

.alert {
  margin-bottom: 12px;
}

.table-wrap {
  border-radius: var(--cb-radius-sm);
  overflow: hidden;
}

.data-table {
  width: 100%;
}

.table-wrap :deep(.el-table) {
  --el-table-border: 1px solid rgba(15, 23, 42, 0.06);
}

.pager {
  margin-top: 18px;
  display: flex;
  justify-content: flex-end;
}

.hint {
  margin: 18px 0 0;
  font-size: 12px;
  color: var(--cb-text-muted);
  line-height: 1.55;
}

.hint code {
  font-size: 12px;
  padding: 2px 6px;
  background: #f1f5f9;
  border-radius: 6px;
  color: var(--cb-text-secondary);
}
</style>
