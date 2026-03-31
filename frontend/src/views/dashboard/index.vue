<template>
  <div class="dashboard-container">
    <!-- 顶部工具栏 -->
    <div class="dashboard-header">
      <div class="header-left">
        <h1 class="page-title">数据看板</h1>
        <span class="last-update">最后更新: {{ lastUpdateTime }}</span>
      </div>
      <div class="header-right">
        <el-button type="primary" :icon="Refresh" :loading="loading" @click="refreshData">
          刷新数据
        </el-button>
      </div>
    </div>

    <el-alert
      v-if="error"
      type="error"
      :title="error"
      show-icon
      closable
      @close="error = null"
      style="margin-bottom: 20px"
    />

    <!-- 核心指标卡片 -->
    <el-row :gutter="20" class="metrics-row">
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="metric-card metric-card-primary">
          <div class="metric-icon">
            <el-icon :size="32"><ShoppingCart /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">今日订单</div>
            <div class="metric-value">{{ summary?.todayOrders || 0 }}</div>
            <div class="metric-compare" :class="getCompareClass(summary?.todayOrders, summary?.yesterdayOrders)">
              <el-icon><component :is="getCompareIcon(summary?.todayOrders, summary?.yesterdayOrders)" /></el-icon>
              {{ getCompareText(summary?.todayOrders, summary?.yesterdayOrders) }}
            </div>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :sm="12" :lg="6">
        <div class="metric-card metric-card-success">
          <div class="metric-icon">
            <el-icon :size="32"><Money /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">今日销售额</div>
            <div class="metric-value">{{ formatCurrency(summary?.todaySales) }}</div>
            <div class="metric-compare" :class="getCompareClass(summary?.todaySales, summary?.yesterdaySales)">
              <el-icon><component :is="getCompareIcon(summary?.todaySales, summary?.yesterdaySales)" /></el-icon>
              {{ getCompareText(summary?.todaySales, summary?.yesterdaySales) }}
            </div>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :sm="12" :lg="6">
        <div class="metric-card metric-card-warning">
          <div class="metric-icon">
            <el-icon :size="32"><Box /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">待发货</div>
            <div class="metric-value">{{ summary?.pendingShipment || 0 }}</div>
            <div class="metric-sub">进行中: {{ summary?.openOrders || 0 }}</div>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :sm="12" :lg="6">
        <div class="metric-card metric-card-danger">
          <div class="metric-icon">
            <el-icon :size="32"><Warning /></el-icon>
          </div>
          <div class="metric-content">
            <div class="metric-label">低库存预警</div>
            <div class="metric-value">{{ summary?.lowStockCount || 0 }}</div>
            <div class="metric-sub">需要补货</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 平台统计 -->
    <el-row :gutter="20" class="section-row">
      <el-col :xs="24" :lg="16">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><TrendCharts /></el-icon>
                <span>销售趋势（近7天）</span>
              </div>
            </div>
          </template>
          <div class="chart-wrapper">
            <div v-if="summary?.salesTrend && summary.salesTrend.length > 0" class="trend-chart">
              <div v-for="(item, index) in summary.salesTrend" :key="index" class="trend-bar">
                <div class="bar-wrapper">
                  <div
                    class="bar-fill"
                    :style="{ height: getBarHeight(item.amount, maxSalesAmount) + '%' }"
                  >
                    <span class="bar-value">{{ formatCurrency(item.amount) }}</span>
                  </div>
                </div>
                <div class="bar-label">{{ formatDate(item.date) }}</div>
                <div class="bar-count">{{ item.count }}单</div>
              </div>
            </div>
            <el-empty v-else description="暂无数据" :image-size="100" />
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="8">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><Platform /></el-icon>
                <span>平台分布</span>
              </div>
            </div>
          </template>
          <div class="platform-list">
            <div v-for="item in summary?.platformStats" :key="item.platform" class="platform-item">
              <div class="platform-header">
                <span class="platform-name">{{ getPlatformName(item.platform) }}</span>
                <span class="platform-amount">{{ formatCurrency(item.totalAmount) }}</span>
              </div>
              <div class="platform-stats">
                <span class="stat-item">订单: {{ item.orderCount }}</span>
                <span class="stat-item">待处理: {{ item.pendingCount }}</span>
                <span class="stat-item">已完成: {{ item.completedCount }}</span>
              </div>
              <div class="platform-progress">
                <el-progress
                  :percentage="getPlatformPercentage(item.totalAmount)"
                  :color="getPlatformColor(item.platform)"
                  :show-text="false"
                />
              </div>
            </div>
            <el-empty v-if="!summary?.platformStats || summary.platformStats.length === 0" description="暂无数据" :image-size="80" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 排行榜 -->
    <el-row :gutter="20" class="section-row">
      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><Goods /></el-icon>
                <span>热销商品 TOP10</span>
              </div>
            </div>
          </template>
          <div class="ranking-list">
            <div v-for="(item, index) in summary?.topProducts" :key="index" class="ranking-item">
              <div class="rank-badge" :class="'rank-' + (index + 1)">{{ index + 1 }}</div>
              <div class="rank-content">
                <div class="rank-name">{{ item.productName }}</div>
                <div class="rank-platform">{{ getPlatformName(item.platform) }}</div>
              </div>
              <div class="rank-value">
                <div class="rank-quantity">{{ item.quantity }}件</div>
                <div class="rank-amount">{{ formatCurrency(item.amount) }}</div>
              </div>
            </div>
            <el-empty v-if="!summary?.topProducts || summary.topProducts.length === 0" description="暂无数据" :image-size="80" />
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><Location /></el-icon>
                <span>热门国家 TOP10</span>
              </div>
            </div>
          </template>
          <div class="ranking-list">
            <div v-for="(item, index) in summary?.topCountries" :key="index" class="ranking-item">
              <div class="rank-badge" :class="'rank-' + (index + 1)">{{ index + 1 }}</div>
              <div class="rank-content">
                <div class="rank-name">{{ item.country }}</div>
              </div>
              <div class="rank-value">
                <div class="rank-quantity">{{ item.count }}单</div>
                <div class="rank-amount">{{ formatCurrency(item.amount) }}</div>
              </div>
            </div>
            <el-empty v-if="!summary?.topCountries || summary.topCountries.length === 0" description="暂无数据" :image-size="80" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 实时订单和系统状态 -->
    <el-row :gutter="20" class="section-row">
      <el-col :xs="24" :lg="16">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><Clock /></el-icon>
                <span>最近订单</span>
              </div>
            </div>
          </template>
          <div class="recent-orders">
            <el-table :data="summary?.recentOrders" style="width: 100%" :show-header="true">
              <el-table-column prop="orderId" label="订单号" width="150" />
              <el-table-column prop="platform" label="平台" width="100">
                <template #default="{ row }">
                  <el-tag :type="getPlatformTagType(row.platform)" size="small">
                    {{ getPlatformName(row.platform) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="amount" label="金额" width="100">
                <template #default="{ row }">
                  {{ formatCurrency(row.amount) }}
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="getStatusTagType(row.status)" size="small">
                    {{ getStatusLabel(row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="country" label="国家" width="120" />
              <el-table-column prop="orderTime" label="下单时间">
                <template #default="{ row }">
                  {{ formatDateTime(row.orderTime) }}
                </template>
              </el-table-column>
            </el-table>
            <el-empty v-if="!summary?.recentOrders || summary.recentOrders.length === 0" description="暂无订单" :image-size="80" />
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="8">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <div class="card-title">
                <el-icon><Monitor /></el-icon>
                <span>系统状态</span>
              </div>
            </div>
          </template>
          <div class="system-health">
            <div class="health-item">
              <div class="health-label">商品总数</div>
              <div class="health-value">{{ summary?.systemHealth?.totalProducts || 0 }}</div>
            </div>
            <div class="health-item">
              <div class="health-label">库存总量</div>
              <div class="health-value">{{ summary?.systemHealth?.totalInventory || 0 }}</div>
            </div>
            <div class="health-item">
              <div class="health-label">API成功率</div>
              <div class="health-value">
                <el-progress
                  :percentage="summary?.systemHealth?.apiSuccessRate || 0"
                  :color="getHealthColor(summary?.systemHealth?.apiSuccessRate)"
                />
              </div>
            </div>
            <div class="health-item">
              <div class="health-label">最后同步</div>
              <div class="health-value health-time">{{ summary?.systemHealth?.lastSyncTime || '-' }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { fetchDashboardSummary, type DashboardSummary } from '@/api/dashboard'
import { messageFromError } from '@/utils/error'
import {
  Refresh, ShoppingCart, Money, Box, Warning, TrendCharts,
  Platform, Goods, Location, Clock, Monitor, ArrowUp, ArrowDown
} from '@element-plus/icons-vue'

const summary = ref<DashboardSummary | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const lastUpdateTime = ref('')

const maxSalesAmount = computed(() => {
  if (!summary.value?.salesTrend || summary.value.salesTrend.length === 0) return 0
  return Math.max(...summary.value.salesTrend.map(item => item.amount))
})

const totalPlatformAmount = computed(() => {
  if (!summary.value?.platformStats) return 0
  return summary.value.platformStats.reduce((sum, item) => sum + item.totalAmount, 0)
})

function formatCurrency(value: number | undefined): string {
  if (value === undefined || value === null) return '$0.00'
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(value)
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

function formatDateTime(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getCompareClass(current: number | undefined, previous: number | undefined): string {
  if (!current || !previous) return ''
  return current >= previous ? 'compare-up' : 'compare-down'
}

function getCompareIcon(current: number | undefined, previous: number | undefined) {
  if (!current || !previous) return null
  return current >= previous ? ArrowUp : ArrowDown
}

function getCompareText(current: number | undefined, previous: number | undefined): string {
  if (!current || !previous) return '较昨日 -'
  const diff = current - previous
  const percent = previous > 0 ? Math.abs((diff / previous) * 100).toFixed(1) : '0.0'
  return `较昨日 ${diff >= 0 ? '+' : ''}${percent}%`
}

function getBarHeight(value: number, max: number): number {
  if (max === 0) return 0
  return (value / max) * 100
}

function getPlatformName(platform: string): string {
  const names: Record<string, string> = {
    shopee: 'Shopee',
    aliexpress: '速卖通',
    tiktok: 'TikTok'
  }
  return names[platform] || platform
}

function getPlatformColor(platform: string): string {
  const colors: Record<string, string> = {
    shopee: '#ee4d2d',
    aliexpress: '#ff6a00',
    tiktok: '#000000'
  }
  return colors[platform] || '#409eff'
}

function getPlatformTagType(platform: string): '' | 'success' | 'warning' | 'danger' | 'info' {
  const types: Record<string, '' | 'success' | 'warning' | 'danger' | 'info'> = {
    shopee: 'danger',
    aliexpress: 'warning',
    tiktok: ''
  }
  return types[platform] || 'info'
}

function getPlatformPercentage(amount: number): number {
  if (totalPlatformAmount.value === 0) return 0
  return (amount / totalPlatformAmount.value) * 100
}

function getStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    pending: '待处理',
    processing: '处理中',
    paid: '已支付',
    shipped: '已发货',
    completed: '已完成',
    cancelled: '已取消'
  }
  return labels[status] || status
}

function getStatusTagType(status: string): '' | 'success' | 'warning' | 'danger' | 'info' {
  const types: Record<string, '' | 'success' | 'warning' | 'danger' | 'info'> = {
    pending: 'info',
    processing: 'warning',
    paid: '',
    shipped: '',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

function getHealthColor(rate: number | undefined): string {
  if (!rate) return '#f56c6c'
  if (rate >= 95) return '#67c23a'
  if (rate >= 80) return '#e6a23c'
  return '#f56c6c'
}

async function refreshData() {
  loading.value = true
  error.value = null
  try {
    summary.value = await fetchDashboardSummary()
    lastUpdateTime.value = new Date().toLocaleTimeString('zh-CN')
  } catch (e) {
    error.value = messageFromError(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refreshData()
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
  background: #f5f7fa;
  min-height: 100vh;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.last-update {
  font-size: 13px;
  color: #909399;
}

.metrics-row {
  margin-bottom: 20px;
}

.metric-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
  margin-bottom: 20px;
}

.metric-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
}

.metric-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  color: white;
}

.metric-card-primary .metric-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.metric-card-success .metric-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.metric-card-warning .metric-icon {
  background: linear-gradient(135deg, #ffa751 0%, #ffe259 100%);
}

.metric-card-danger .metric-icon {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
}

.metric-content {
  flex: 1;
}

.metric-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.metric-value {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  margin-bottom: 4px;
}

.metric-sub {
  font-size: 13px;
  color: #606266;
}

.metric-compare {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 500;
}

.compare-up {
  color: #67c23a;
}

.compare-down {
  color: #f56c6c;
}

.section-row {
  margin-bottom: 20px;
}

.chart-card {
  border-radius: 12px;
  margin-bottom: 20px;
}

.chart-card :deep(.el-card__header) {
  border-bottom: 1px solid #f0f0f0;
  padding: 16px 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.chart-wrapper {
  min-height: 300px;
  padding: 20px 0;
}

.trend-chart {
  display: flex;
  align-items: flex-end;
  justify-content: space-around;
  height: 280px;
  padding: 0 10px;
}

.trend-bar {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.bar-wrapper {
  width: 100%;
  height: 220px;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.bar-fill {
  width: 60%;
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px 8px 0 0;
  position: relative;
  min-height: 20px;
  transition: all 0.3s;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 8px;
}

.bar-fill:hover {
  opacity: 0.8;
}

.bar-value {
  font-size: 11px;
  color: white;
  font-weight: 600;
}

.bar-label {
  font-size: 12px;
  color: #606266;
  font-weight: 500;
}

.bar-count {
  font-size: 11px;
  color: #909399;
}

.platform-list {
  padding: 10px 0;
}

.platform-item {
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.platform-item:last-child {
  border-bottom: none;
}

.platform-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.platform-name {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.platform-amount {
  font-size: 16px;
  font-weight: 700;
  color: #409eff;
}

.platform-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.stat-item {
  font-size: 12px;
  color: #909399;
}

.platform-progress {
  margin-top: 8px;
}

.ranking-list {
  padding: 10px 0;
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.ranking-item:last-child {
  border-bottom: none;
}

.rank-badge {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 700;
  color: white;
  background: #909399;
}

.rank-1 {
  background: linear-gradient(135deg, #ffd700 0%, #ffed4e 100%);
  color: #333;
}

.rank-2 {
  background: linear-gradient(135deg, #c0c0c0 0%, #e8e8e8 100%);
  color: #333;
}

.rank-3 {
  background: linear-gradient(135deg, #cd7f32 0%, #e8a87c 100%);
  color: white;
}

.rank-content {
  flex: 1;
  min-width: 0;
}

.rank-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rank-platform {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.rank-value {
  text-align: right;
}

.rank-quantity {
  font-size: 14px;
  font-weight: 600;
  color: #409eff;
}

.rank-amount {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.recent-orders {
  min-height: 300px;
}

.system-health {
  padding: 10px 0;
}

.health-item {
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.health-item:last-child {
  border-bottom: none;
}

.health-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 12px;
}

.health-value {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
}

.health-time {
  font-size: 14px;
  font-weight: 400;
  color: #606266;
}
</style>
