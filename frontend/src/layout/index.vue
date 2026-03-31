<template>
  <el-container class="layout-container" direction="horizontal">
    <el-aside width="220px" class="aside">
      <div class="logo">
        <span class="logo-mark">
          <el-icon><Histogram /></el-icon>
        </span>
        <span class="logo-text">跨境电商</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        background-color="transparent"
        text-color="#94a3b8"
        active-text-color="#f8fafc"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataLine /></el-icon>
          <span>数据看板</span>
        </el-menu-item>

        <el-sub-menu index="1">
          <template #title>
            <el-icon><Shop /></el-icon>
            <span>Shopee 管理</span>
          </template>
          <el-menu-item index="/shopee/orders">订单处理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="2">
          <template #title>
            <el-icon><Goods /></el-icon>
            <span>速卖通 管理</span>
          </template>
          <el-menu-item index="/aliexpress/orders">订单处理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="3">
          <template #title>
            <el-icon><VideoCamera /></el-icon>
            <span>TK 管理</span>
          </template>
          <el-menu-item index="/tk/orders">订单处理</el-menu-item>
        </el-sub-menu>

        <el-menu-item index="/products">
          <el-icon><Box /></el-icon>
          <span>商品管理</span>
        </el-menu-item>

        <el-menu-item index="/inventory">
          <el-icon><Tickets /></el-icon>
          <span>库存管理</span>
        </el-menu-item>

        <el-menu-item index="/platform/config">
          <el-icon><Setting /></el-icon>
          <span>平台配置</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container class="main-wrap" direction="vertical">
      <el-header class="header" height="56px">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>{{ currentRouteTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <div class="user-pill">
            <el-avatar :size="32" src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" />
            <span class="user-name">管理员</span>
          </div>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activeMenu = computed(() => route.path)
const currentRouteTitle = computed(() => {
  const t = route.meta.title
  return typeof t === 'string' && t.length > 0 ? t : '工作台'
})
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.aside {
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, var(--cb-sidebar-bg) 0%, var(--cb-sidebar-bg-end) 100%);
  border-right: 1px solid var(--cb-sidebar-border);
  box-shadow: 4px 0 32px rgba(15, 23, 42, 0.14);
}

.logo {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 14px 0 16px;
  height: 64px;
  border-bottom: 1px solid var(--cb-sidebar-border);
}

.logo-mark {
  flex-shrink: 0;
  width: 38px;
  height: 38px;
  border-radius: 11px;
  background: linear-gradient(135deg, var(--cb-primary) 0%, var(--cb-primary-light) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 20px;
  box-shadow: 0 4px 14px rgba(79, 70, 229, 0.35);
}

.logo-text {
  font-size: 15px;
  font-weight: 700;
  color: var(--cb-sidebar-text-strong);
  letter-spacing: 0.03em;
  line-height: 1.25;
}

.sidebar-menu {
  flex: 1;
  border-right: none !important;
  padding: 12px 8px 16px;
  overflow-y: auto;
}

.sidebar-menu :deep(.el-menu-item),
.sidebar-menu :deep(.el-sub-menu__title) {
  border-radius: var(--cb-radius-sm);
  margin: 3px 0;
}

.sidebar-menu :deep(.el-menu-item:hover),
.sidebar-menu :deep(.el-sub-menu__title:hover) {
  background-color: var(--cb-sidebar-hover) !important;
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background: var(--cb-sidebar-active-bg) !important;
  color: var(--cb-sidebar-text-strong) !important;
}

.sidebar-menu :deep(.el-sub-menu .el-menu) {
  background: transparent !important;
}

.sidebar-menu :deep(.el-sub-menu .el-menu-item) {
  padding-left: 48px !important;
  min-width: auto;
}

.main-wrap {
  flex: 1;
  min-width: 0;
  flex-direction: column;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: var(--cb-header-bg) !important;
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  border-bottom: 1px solid var(--cb-header-border);
  box-shadow:
    0 1px 0 rgba(255, 255, 255, 0.55) inset,
    0 1px 2px rgba(15, 23, 42, 0.03);
}

.header-right {
  display: flex;
  align-items: center;
}

.user-pill {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 4px 14px 4px 5px;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.05);
  cursor: default;
  transition: background 0.2s ease;
}

.user-pill:hover {
  background: rgba(79, 70, 229, 0.09);
}

.user-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--cb-text);
}

.main-content {
  background: var(--cb-page-gradient);
  padding: 24px;
  min-height: calc(100vh - 56px);
  overflow: auto;
}
</style>
