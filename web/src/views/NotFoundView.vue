<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <div class="not-found">
    <el-result icon="warning" :title="t('notFound.title')" :sub-title="t('notFound.subtitle')">
      <template #extra>
        <el-button type="primary" @click="$router.push('/')">{{ t('notFound.backHome') }}</el-button>
      </template>
    </el-result>
  </div>
</template>

<style scoped>
.not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}
</style>
