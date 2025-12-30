<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import type { FormRules } from 'element-plus'

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const form = reactive({
  email: '',
  password: '',
  confirmPassword: ''
})

const rules = computed<FormRules>(() => ({
  email: [
    { required: true, message: t('validation.emailRequired'), trigger: 'blur' },
    { type: 'email', message: t('validation.emailInvalid'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('validation.passwordRequired'), trigger: 'blur' },
    { min: 8, message: t('validation.passwordMin'), trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: t('validation.confirmRequired'), trigger: 'blur' },
    {
      validator: (_: unknown, value: string, callback: (error?: Error) => void) => {
        if (value !== form.password) {
          callback(new Error(t('validation.passwordMismatch')))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}))

async function handleSubmit() {
  loading.value = true
  try {
    await authStore.register(form.email, form.password)
    ElMessage.success(t('auth.registerSuccess'))
    router.push('/login')
  } catch {
    // Error handled by interceptor
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-view">
    <div class="auth-header">
      <h2>{{ t('auth.createAccount') }}</h2>
      <p>{{ t('auth.getStartedWith') }}</p>
    </div>

    <el-form :model="form" :rules="rules" label-position="top" @submit.prevent="handleSubmit">
      <el-form-item :label="t('auth.email')" prop="email">
        <el-input
          v-model="form.email"
          type="email"
          :placeholder="t('auth.email')"
          size="large"
        />
      </el-form-item>

      <el-form-item :label="t('auth.password')" prop="password">
        <el-input
          v-model="form.password"
          type="password"
          :placeholder="t('auth.password')"
          show-password
          size="large"
        />
      </el-form-item>

      <el-form-item :label="t('auth.confirmPassword')" prop="confirmPassword">
        <el-input
          v-model="form.confirmPassword"
          type="password"
          :placeholder="t('auth.confirmPassword')"
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
          {{ t('auth.createAccountBtn') }}
        </el-button>
      </el-form-item>

      <div class="auth-footer">
        {{ t('auth.haveAccount') }}
        <router-link to="/login">{{ t('auth.signInLink') }}</router-link>
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
