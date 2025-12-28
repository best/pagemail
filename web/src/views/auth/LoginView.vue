<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import type { FormRules } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const form = reactive({
  email: '',
  password: ''
})

const rules: FormRules = {
  email: [
    { required: true, message: 'Please enter email', trigger: 'blur' },
    { type: 'email', message: 'Please enter valid email', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Please enter password', trigger: 'blur' }
  ]
}

async function handleSubmit() {
  loading.value = true
  try {
    await authStore.login(form.email, form.password)
    ElMessage.success('Login successful')
    router.push('/')
  } catch (error) {
    // Error handled by interceptor
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-view">
    <div class="auth-header">
      <h2>Welcome back</h2>
      <p>Sign in to continue</p>
    </div>

    <el-form :model="form" :rules="rules" label-position="top" @submit.prevent="handleSubmit">
      <el-form-item label="Email" prop="email">
        <el-input
          v-model="form.email"
          type="email"
          placeholder="Enter your email"
          size="large"
        />
      </el-form-item>

      <el-form-item label="Password" prop="password">
        <el-input
          v-model="form.password"
          type="password"
          placeholder="Enter your password"
          show-password
          size="large"
        />
      </el-form-item>

      <el-form-item class="submit-item">
        <el-button
          type="primary"
          native-type="submit"
          :loading="loading"
          size="large"
          class="submit-btn"
        >
          Sign in
        </el-button>
      </el-form-item>

      <div class="auth-footer">
        Don't have an account?
        <router-link to="/register">Create one</router-link>
      </div>
    </el-form>
  </div>
</template>

<style scoped>
.auth-view {
  width: 100%;
}

.auth-header {
  text-align: center;
  margin-bottom: 2rem;
}

.auth-header h2 {
  margin: 0 0 0.5rem;
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--pm-text-heading);
}

.auth-header p {
  margin: 0;
  color: var(--pm-text-muted);
}

:deep(.el-input__wrapper) {
  padding: 4px 15px;
  box-shadow: 0 0 0 1px var(--el-border-color) inset;
  transition: box-shadow 0.2s;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px var(--el-border-color-hover) inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px var(--pm-primary) inset !important;
}

.submit-item {
  margin-top: 1.5rem;
  margin-bottom: 1.5rem;
}

.submit-btn {
  width: 100%;
  height: 48px;
  font-size: 1rem;
  font-weight: 600;
}

.auth-footer {
  text-align: center;
  color: var(--pm-text-muted);
}

.auth-footer a {
  color: var(--pm-primary);
  text-decoration: none;
  font-weight: 500;
  margin-left: 4px;
}

.auth-footer a:hover {
  text-decoration: underline;
}
</style>
