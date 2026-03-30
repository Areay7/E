<template>
  <div class="dashboard">
    <div class="toolbar-shell">
      <div class="toolbar">
        <el-button type="primary" size="default" :loading="summaryLoading" @click="refreshSummary">
          从服务器刷新汇总
        </el-button>
        <span class="toolbar-hint">对接 <code>GET /api/v1/dashboard/summary</code> 后显示数字</span>
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
      <el-col :xs="24" :sm="12" :md="8">
        <el-card shadow="hover" class="stat-uplift stat-card stat-card--sales">
          <div class="stat-inner">
            <div class="stat-icon-wrap stat-icon-wrap--sales">
              <el-icon :size="22"><TrendCharts /></el-icon>
            </div>
            <div class="stat-body">
              <div class="stat-label">今日销售额 (USD)</div>
              <div class="stat-value">{{ formatUsd(summary?.todaySalesUsd) }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="8">
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
      <el-col :xs="24" :sm="12" :md="8">
        <el-card shadow="hover" class="stat-uplift stat-card stat-card--orders">
          <div class="stat-inner">
            <div class="stat-icon-wrap stat-icon-wrap--orders">
              <el-icon :size="22"><List /></el-icon>
            </div>
            <div class="stat-body">
              <div class="stat-label">进行中订单</div>
              <div class="stat-value">{{ formatInt(summary?.openOrders) }}</div>
            </div>
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
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { fetchDashboardSummary } from '@/api/dashboard'
import type { DashboardSummary } from '@/types/order'
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
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 1280px;
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
  font-size: 12px;
  color: var(--cb-text-muted);
}

.toolbar-hint code {
  font-size: 12px;
  padding: 2px 6px;
  background: #f1f5f9;
  border-radius: 6px;
  color: var(--cb-text-secondary);
}

.alert {
  margin: 0;
}

.stats {
  margin: 0 !important;
}

.stat-uplift {
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.stat-uplift:hover {
  transform: translateY(-2px);
  box-shadow: var(--cb-shadow-hover) !important;
}

.stat-card {
  margin-bottom: 4px;
  border-top: 3px solid transparent;
}

.stat-card--sales {
  border-top-color: #4f46e5;
}

.stat-card--ship {
  border-top-color: #d97706;
}

.stat-card--orders {
  border-top-color: #059669;
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

.stat-icon-wrap--sales {
  background: linear-gradient(135deg, #4f46e5, #818cf8);
  box-shadow: 0 6px 16px rgba(79, 70, 229, 0.25);
}

.stat-icon-wrap--ship {
  background: linear-gradient(135deg, #d97706, #fbbf24);
  box-shadow: 0 6px 16px rgba(217, 119, 6, 0.25);
}

.stat-icon-wrap--orders {
  background: linear-gradient(135deg, #059669, #34d399);
  box-shadow: 0 6px 16px rgba(5, 150, 105, 0.22);
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
