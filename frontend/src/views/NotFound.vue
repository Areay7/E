<template>
  <div class="not-found">
    <div class="not-found-glow" aria-hidden="true" />
    <el-card class="not-found-card" shadow="never">
      <el-result
        icon="warning"
        title="页面不存在"
        sub-title="请检查链接是否正确，或返回数据看板继续操作。"
      >
        <template #extra>
          <el-button type="primary" size="large" round @click="goHome">回到数据看板</el-button>
        </template>
      </el-result>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

function goHome() {
  void router.replace({ path: '/dashboard' })
}
</script>

<style scoped>
.not-found {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  position: relative;
  background: var(--cb-page-gradient);
  overflow: hidden;
}

.not-found-glow {
  position: absolute;
  width: 480px;
  height: 480px;
  top: -120px;
  right: -80px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(129, 140, 248, 0.35) 0%, transparent 70%);
  pointer-events: none;
}

.not-found-glow::after {
  content: '';
  position: absolute;
  width: 360px;
  height: 360px;
  bottom: -200px;
  left: -100px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(79, 70, 229, 0.15) 0%, transparent 70%);
}

.not-found-card {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 520px;
  border-radius: var(--cb-radius) !important;
  box-shadow: var(--cb-shadow-hover) !important;
}

.not-found-card :deep(.el-card__body) {
  padding: 40px 32px 36px;
}

.not-found-card :deep(.el-result__title) {
  font-weight: 700;
  color: var(--cb-text);
}

.not-found-card :deep(.el-result__subtitle) {
  color: var(--cb-text-secondary);
  line-height: 1.6;
}
</style>
