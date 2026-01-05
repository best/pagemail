<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useSiteConfigStore } from '@/stores/siteConfig'
import { useUiStore } from '@/stores/ui'
import LanguageSwitcher from '@/components/common/LanguageSwitcher.vue'
import { Moon, Sunny } from '@element-plus/icons-vue'
import apiClient from '@/api/client'

const props = withDefaults(defineProps<{
  solid?: boolean
}>(), {
  solid: false
})

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()
const siteConfig = useSiteConfigStore()
const uiStore = useUiStore()

const isScrolled = ref(false)
const avatarUrl = ref<string>('')

const handleScroll = () => {
  isScrolled.value = window.scrollY > 50
}

const fetchAvatar = async () => {
  if (!authStore.user?.avatar_url) {
    avatarUrl.value = ''
    return
  }
  try {
    const response = await apiClient.get(authStore.user.avatar_url, { responseType: 'blob' })
    if (avatarUrl.value) URL.revokeObjectURL(avatarUrl.value)
    avatarUrl.value = URL.createObjectURL(response.data)
  } catch {
    avatarUrl.value = ''
  }
}

const hasBackground = computed(() => props.solid || isScrolled.value)
const isShrunk = computed(() => isScrolled.value)

const userInitial = computed(() => authStore.user?.email?.[0]?.toUpperCase() || '?')

const handleUserCommand = (command: string) => {
  if (command === 'dashboard') {
    router.push({ name: 'dashboard' })
  } else if (command === 'logout') {
    authStore.logout()
    router.push({ name: 'login' })
  }
}

onMounted(() => {
  handleScroll()
  window.addEventListener('scroll', handleScroll, { passive: true })
  fetchAvatar()
})

watch(() => authStore.user?.avatar_url, () => fetchAvatar())

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
  if (avatarUrl.value) URL.revokeObjectURL(avatarUrl.value)
})
</script>

<template>
  <nav :class="['pub-nav', { 'has-bg': hasBackground, 'is-shrunk': isShrunk }]">
    <div class="container nav-inner">
      <router-link to="/" class="logo">{{ siteConfig.siteName }}</router-link>
      <div class="nav-actions">
        <LanguageSwitcher />
        <el-button :icon="uiStore.isDark ? Sunny : Moon" text circle aria-label="Toggle theme" @click="uiStore.toggleTheme" />
        <template v-if="authStore.isAuthenticated">
          <el-dropdown trigger="click" @command="handleUserCommand">
            <el-button text class="avatar-btn" aria-label="User menu">
              <el-avatar :size="32" :src="avatarUrl">{{ userInitial }}</el-avatar>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item v-if="authStore.user?.email" disabled>
                  {{ authStore.user.email }}
                </el-dropdown-item>
                <el-dropdown-item command="dashboard">{{ t('nav.dashboard') }}</el-dropdown-item>
                <el-dropdown-item divided command="logout">{{ t('common.logout') }}</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <template v-else>
          <router-link to="/login" class="nav-link">{{ t('landing.signIn') }}</router-link>
          <el-button type="primary" round @click="router.push('/register')">
            {{ t('landing.getStarted') }}
          </el-button>
        </template>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1.5rem;
}

.pub-nav {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  padding: 1.25rem 0;
  transition: all 0.3s;
}

.pub-nav.has-bg {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(12px);
  box-shadow: var(--pm-shadow-sm);
}

.pub-nav.is-shrunk {
  padding: 0.75rem 0;
}

html.dark .pub-nav.has-bg {
  background: rgba(15, 23, 42, 0.9);
}

.nav-inner {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  font-weight: 800;
  font-size: 1.5rem;
  color: var(--pm-primary);
  letter-spacing: -0.02em;
  text-decoration: none;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.nav-link {
  color: var(--pm-text-body);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;
}

.nav-link:hover {
  color: var(--pm-primary);
}

.avatar-btn {
  padding: 0;
}

@media (max-width: 768px) {
  .nav-actions {
    gap: 0.5rem;
  }

  .nav-link {
    display: none;
  }

  .logo {
    font-size: 1.25rem;
  }
}
</style>
