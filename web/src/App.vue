<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useUiStore } from '@/stores/ui'

const uiStore = useUiStore()
uiStore.initTheme()
</script>

<template>
  <RouterView />
</template>

<style>
html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
}
</style>
