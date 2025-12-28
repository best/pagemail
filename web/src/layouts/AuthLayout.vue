<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="auth-layout">
    <div class="auth-card">
      <div class="auth-header">
        <h1>Pagemail</h1>
        <p>Web Page Capture & Delivery</p>
      </div>
      <RouterView />
    </div>
  </div>
</template>

<style scoped>
.auth-layout {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, var(--el-color-primary-light-9) 0%, var(--el-bg-color-page) 100%);
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 2rem;
  background: var(--el-bg-color);
  border-radius: 12px;
  box-shadow: var(--el-box-shadow-light);
}

.auth-header {
  text-align: center;
  margin-bottom: 2rem;
}

.auth-header h1 {
  margin: 0;
  color: var(--el-color-primary);
  font-size: 2rem;
}

.auth-header p {
  margin: 0.5rem 0 0;
  color: var(--el-text-color-secondary);
}
</style>
