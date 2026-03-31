<template>
  <el-form :model="form" label-width="120px">
    <el-form-item label="启用状态">
      <el-switch v-model="form.enabled" />
    </el-form-item>

    <el-form-item label="店铺 ID">
      <el-input v-model="form.shopId" placeholder="请输入店铺 ID" />
    </el-form-item>

    <el-form-item label="店铺名称">
      <el-input v-model="form.shopName" placeholder="请输入店铺名称" />
    </el-form-item>

    <el-form-item label="App Key">
      <el-input v-model="form.appKey" placeholder="请输入 App Key" />
    </el-form-item>

    <el-form-item label="App Secret">
      <el-input v-model="form.appSecret" type="password" placeholder="请输入 App Secret" show-password />
    </el-form-item>

    <el-form-item label="Partner ID" v-if="platform === 'shopee'">
      <el-input v-model="form.partnerId" placeholder="请输入 Partner ID" />
    </el-form-item>

    <el-form-item label="Partner Key" v-if="platform === 'shopee'">
      <el-input v-model="form.partnerKey" type="password" placeholder="请输入 Partner Key" show-password />
    </el-form-item>

    <el-form-item label="API URL">
      <el-input v-model="form.apiUrl" placeholder="请输入 API URL" />
    </el-form-item>

    <el-form-item label="启用同步">
      <el-switch v-model="form.syncEnabled" />
    </el-form-item>

    <el-form-item>
      <el-button type="primary" @click="handleSave">保存配置</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  platform: string
  config: any
}>()

const emit = defineEmits<{
  save: [platform: string, config: any]
}>()

const form = ref({
  enabled: true,
  shopId: '',
  shopName: '',
  appKey: '',
  appSecret: '',
  partnerId: '',
  partnerKey: '',
  apiUrl: '',
  syncEnabled: true
})

watch(() => props.config, (newConfig) => {
  if (newConfig) {
    form.value = { ...newConfig }
  }
}, { immediate: true })

function handleSave() {
  emit('save', props.platform, form.value)
}
</script>
