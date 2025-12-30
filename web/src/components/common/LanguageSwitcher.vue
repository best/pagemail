<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUiStore } from '@/stores/ui'

const { t } = useI18n()
const uiStore = useUiStore()

const label = computed(() => uiStore.language === 'zh' ? '中' : 'En')

const handleCommand = (command: string) => {
  uiStore.setLanguage(command as 'en' | 'zh')
}
</script>

<template>
  <el-dropdown trigger="click" @command="handleCommand">
    <el-button class="lang-btn" text circle :aria-label="t('language.switch')">
      {{ label }}
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item command="en" :disabled="uiStore.language === 'en'">
          {{ t('language.en') }}
        </el-dropdown-item>
        <el-dropdown-item command="zh" :disabled="uiStore.language === 'zh'">
          {{ t('language.zh') }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<style scoped>
.lang-btn {
  font-weight: 600;
  width: 32px;
  height: 32px;
  transition: background-color 0.2s ease;
}

.lang-btn:hover {
  background-color: var(--el-fill-color-light);
}

.lang-btn:focus-visible {
  outline: 2px solid var(--el-color-primary);
  outline-offset: 2px;
}
</style>
