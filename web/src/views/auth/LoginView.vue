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
  <el-form :model="form" :rules="rules" label-position="top" @submit.prevent="handleSubmit">
    <el-form-item label="Email" prop="email">
      <el-input v-model="form.email" type="email" placeholder="Enter your email" />
    </el-form-item>

    <el-form-item label="Password" prop="password">
      <el-input v-model="form.password" type="password" placeholder="Enter your password" show-password />
    </el-form-item>

    <el-form-item>
      <el-button type="primary" native-type="submit" :loading="loading" style="width: 100%">
        Login
      </el-button>
    </el-form-item>

    <div class="auth-footer">
      Don't have an account?
      <router-link to="/register">Register</router-link>
    </div>
  </el-form>
</template>

<style scoped>
.auth-footer {
  text-align: center;
  color: var(--el-text-color-secondary);
}

.auth-footer a {
  color: var(--el-color-primary);
  text-decoration: none;
}
</style>
