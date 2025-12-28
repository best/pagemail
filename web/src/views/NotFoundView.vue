<script setup lang="ts">
</script>

<template>
  <div class="not-found">
    <el-result icon="warning" title="404" sub-title="Page not found">
      <template #extra>
        <el-button type="primary" @click="$router.push('/')">Back to Home</el-button>
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
