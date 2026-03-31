<template>
  <div class="dashboard">
    <div class="toolbar-shell">
      <div class="toolbar">
        <el-button type="primary" size="default" :loading="summaryLoading" @click="refreshSummary">
          刷新数据
        </el-button>
        <span class="toolbar-hint">实时数据看板</span>
      </div>
    </div>

    <el-alert
      v-if="summaryError"
      type="error"
      :title="summaryError"
      show-icon
      class="alert"
      closable
      @close="summaryError = null"
    />

    <el-row :gutter="20" class="stats">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover" class="stat-uplift stat-card stat-card--orders">
          <div class="stat-inner">
            <div class="stat-icon-wrap stat-icon-wrap--orders">
              <el-icon :size="22"><ShoppingCart /></el-icon>
            </div>
            <div class="stat-body">
              <div class="stat-label">今日订单</div>
              <div class="stat-value">{{ formatInt(summary?.todayOrders) }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover" class="stat-uplift stat-card stat-card--sales">
          <div class="stat-inner">
            <div class="stat-icon-wrap stat-icon-wrap--sales">
              <el-icon :size="22"><TrendCharts /></el-icon>
            </div>
            <div class="stat-body">
              <div class="stat-label">今日销售额</div>
              <div class="stat-value">{{ formatUsd(summary?.todaySales) }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover" class="stat-uplift stat-card stat-card--ship">
          <div class="stat-inner">
            <div class="stat-icon-wrap stat-icon-wrap--ship">
              <el-icon :size="22"><Box /></el-icon>
            </div>
            <div class="stat-body">
              <div class="stat-label">待发货</div>
              <div class="stat-value">{{ formatInt(summary?.pendingShipment) }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover" class="stat-uplift stat-card stat-card--stock">
          <div class="stat-inner">
            <div class="stat-icon-wrap stat-icon-wrap--stock">
              <el-icon :size="22"><Warning /></el-icon>
            </div>
            <div class="stat-body">
              <div class="stat-label">低库存商品</div>
              <div class="stat-value">{{ formatInt(summary?.lowStockCount) }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon><DataLine /></el-icon>
              <span>近7天销售趋势</span>
            </div>
          </template>
          <div class="chart-container">
            <div v-if="summary?.salesTrend && summary.salesTrend.length > 0" class="trend-list">
              <div v-for="item in summary.salesTrend" :key="item.date" class="trend-item">
                <span class="trend-date">{{ item.date }}</span>
                <span class="trend-count">{{ item.count }}单</span>
                <span class="trend-amount">{{ formatUsd(item.amount) }}</span>
              </div>
            </div>
            <el-empty v-else description="暂无数据" :image-size="80" />
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon><Platform /></el-icon>
              <span>平台销售统计</span>
            </div>
          </template>
          <div class="chart-container">
            <div v-if="summary?.platformStats && summary.platformStats.length > 0" class="platform-list">
              <div v-for="item in summary.platformStats" :key="item.platform" class="platform-item">
                <span class="platform-name">{{ item.platform }}</span>
                <span class="platform-count">{{ item.orderCount }}单</span>
                <span class="platform-amount">{{ formatUsd(item.totalAmount) }}</span>
              </div>
            </div>
            <el-empty v-else description="暂无数据" :image-size="80" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon><Goods /></el-icon>
              <span>热销商品TOP10</span>
            </div>
          </template>
          <div class="chart-container">
            <div v-if="summary?.topProducts && summary.topProducts.length > 0" class="product-list">
              <div v-for="(item, index) in summary.topProducts" :key="index" class="product-item">
                <span class="product-rank">{{ index + 1 }}</span>
                <span class="product-name">{{ item.productName }}</span>
                <span class="product-quantity">{{ item.quantity }}件</span>
              </div>
            </div>
            <el-empty v-else description="暂无数据" :image-size="80" />
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon><PieChart /></el-icon>
              <span>订单状态分布</span>
            </div>
          </template>
          <div class="chart-container">
            <div v-if="summary?.orderStatusDist && summary.orderStatusDist.length > 0" class="status-list">
              <div v-for="item in summary.orderStatusDist" :key="item.status" class="status-item">
                <span class="status-label">{{ getStatusLabel(item.status) }}</span>
                <span class="status-count">{{ item.count }}</span>
              </div>
            </div>
            <el-empty v-else description="暂无数据" :image-size="80" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="welcome" shadow="never">
      <template #header>
        <div class="welcome-head">
          <el-icon class="welcome-head-icon"><Right /></el-icon>
          <span>快捷入口</span>
        </div>
      </template>
      <div class="quick-links">
        <router-link v-slot="{ navigate }" to="/shopee/orders" custom>
          <el-button class="quick-btn" round @click="navigate">Shopee 订单</el-button>
        </router-link>
        <router-link v-slot="{ navigate }" to="/aliexpress/orders" custom>
          <el-button class="quick-btn" round @click="navigate">速卖通订单</el-button>
        </router-link>
        <router-link v-slot="{ navigate }" to="/tk/orders" custom>
          <el-button class="quick-btn" round @click="navigate">TikTok 订单</el-button>
        </router-link>
        <router-link v-slot="{ navigate }" to="/products" custom>
          <el-button class="quick-btn" round @click="navigate">商品管理</el-button>
        </router-link>
        <router-link v-slot="{ navigate }" to="/inventory" custom>
          <el-button class="quick-btn" round @click="navigate">库存管理</el-button>
        </router-link>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchDashboardSummary, type DashboardSummary } from '@/api/dashboard'
import { messageFromError } from '@/utils/error'

const summary = ref<DashboardSummary | null>(null)
const summaryLoading = ref(false)
const summaryError = ref<string | null>(null)

function formatUsd(v: number | undefined): string {
  if (v === undefined) return '—'
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(v)
}

function formatInt(v: number | undefined): string {
  if (v === undefined) return '—'
  return String(v)
}

function getStatusLabel(status: string): string {
  const statusMap: Record<string, string> = {
    pending: '待处理',
    processing: '处理中',
    paid: '已支付',
    shipped: '已发货',
    completed: '已完成',
    cancelled: '已取消'
  }
  return statusMap[status] || status
}

async function refreshSummary() {
  summaryLoading.value = true
  summaryError.value = null
  try {
    summary.value = await fetchDashboardSummary()
  } catch (e) {
    summary.value = null
    summaryError.value = messageFromError(e)
  } finally {
    summaryLoading.value = false
  }
}

onMounted(() => {
  refreshSummary()
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 1400px;
}

.toolbar-shell {
  padding: 18px 20px;
  background: var(--cb-surface);
  border-radius: var(--cb-radius);
  border: 1px solid rgba(15, 23, 42, 0.06);
  box-shadow: var(--cb-shadow);
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 14px;
}

.toolbar-hint {
  font-size: 13px;
  color: var(--cb-text-muted);
  font-weight: 500;
}

.alert {
  margin: 0;
}

.stats {
  margin: 0 !important;
}

.stat-uplift {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.stat-uplift:hover {
  transform: translateY(-2px);
  box-shadow: var(--cb-shadow-hover) !important;
}

.stat-card {
  margin-bottom: 4px;
  border-top: 3px solid transparent;
}

.stat-card--orders {
  border-top-color: #059669;
}

.stat-card--sales {
  border-top-color: #4f46e5;
}

.stat-card--ship {
  border-top-color: #d97706;
}

.stat-card--stock {
  border-top-color: #dc2626;
}

.stat-inner {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.stat-icon-wrap {
  flex-shrink: 0;
  width: 48px;
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.stat-icon-wrap--orders {
  background: linear-gradient(135deg, #059669, #34d399);
  box-shadow: 0 6px 16px rgba(5, 150, 105, 0.22);
}

.stat-icon-wrap--sales {
  background: linear-gradient(135deg, #4f46e5, #818cf8);
  box-shadow: 0 6px 16px rgba(79, 70, 229, 0.25);
}

.stat-icon-wrap--ship {
  background: linear-gradient(135deg, #d97706, #fbbf24);
  box-shadow: 0 6px 16px rgba(217, 119, 6, 0.25);
}

.stat-icon-wrap--stock {
  background: linear-gradient(135deg, #dc2626, #f87171);
  box-shadow: 0 6px 16px rgba(220, 38, 38, 0.25);
}

.stat-body {
  min-width: 0;
}

.stat-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--cb-text-muted);
  margin-bottom: 6px;
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: var(--cb-text);
  letter-spacing: -0.02em;
  line-height: 1.15;
}

.chart-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--cb-text);
}

.chart-container {
  min-height: 200px;
}

.trend-list, .platform-list, .product-list, .status-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.trend-item, .platform-item, .product-item, .status-item {
  display: flex;
  align-items: center;
  padding: 12px;
  background: #f8fafc;
  border-radius: 8px;
  gap: 12px;
}

.trend-date, .platform-name, .product-name, .status-label {
  flex: 1;
  font-weight: 500;
  color: var(--cb-text);
}

.trend-count, .platform-count, .product-quantity, .status-count {
  color: var(--cb-text-secondary);
  font-size: 14px;
}

.trend-amount, .platform-amount {
  font-weight: 600;
  color: #4f46e5;
}

.product-rank {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #4f46e5;
  color: white;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 600;
}

.welcome-head {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--cb-text);
}

.welcome-head-icon {
  color: var(--cb-primary);
}

.welcome :deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid rgba(15, 23, 42, 0.06);
}

.quick-links {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.quick-btn {
  border: 1px solid rgba(79, 70, 229, 0.35);
  color: var(--cb-primary);
  background: rgba(79, 70, 229, 0.06);
}

.quick-btn:hover {
  background: rgba(79, 70, 229, 0.12);
  border-color: var(--cb-primary);
  color: var(--cb-primary-dark);
}
</style>
