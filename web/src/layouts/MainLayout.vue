<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUiStore } from '@/stores/ui'
import { useAuthStore } from '@/stores/auth'
import LanguageSwitcher from '@/components/common/LanguageSwitcher.vue'
import {
  House,
  Document,
  Setting,
  User,
  DataAnalysis,
  Memo,
  Moon,
  Sunny,
  Fold,
  Expand
} from '@element-plus/icons-vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const uiStore = useUiStore()
const authStore = useAuthStore()

const menuItems = computed(() => {
  const items = [
    { index: '/dashboard', titleKey: 'nav.dashboard', icon: House },
    { index: '/tasks', titleKey: 'nav.tasks', icon: Document },
    { index: '/settings', titleKey: 'nav.settings', icon: Setting }
  ]

  if (authStore.isAdmin) {
    items.push(
      { index: '/admin/users', titleKey: 'nav.users', icon: User },
      { index: '/admin/system', titleKey: 'nav.system', icon: DataAnalysis },
      { index: '/admin/audit', titleKey: 'nav.auditLogs', icon: Memo }
    )
  }

  return items
})

function handleSelect(index: string) {
  router.push(index)
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <el-container class="main-layout">
    <div class="layout-bg">
      <div class="orb orb-1"></div>
      <div class="orb orb-2"></div>
      <div class="orb orb-3"></div>
      <div class="grid-overlay"></div>
    </div>

    <el-aside :width="uiStore.sidebarCollapsed ? '64px' : '220px'" class="sidebar">
      <div class="logo">
        <span v-if="!uiStore.sidebarCollapsed">Pagemail</span>
        <span v-else>P</span>
      </div>
      <el-menu
        :default-active="route.path"
        :collapse="uiStore.sidebarCollapsed"
        :collapse-transition="false"
        @select="handleSelect"
      >
        <el-menu-item v-for="item in menuItems" :key="item.index" :index="item.index">
          <el-icon><component :is="item.icon" /></el-icon>
          <template #title>{{ t(item.titleKey) }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-button :icon="uiStore.sidebarCollapsed ? Expand : Fold" text @click="uiStore.toggleSidebar" aria-label="Toggle Sidebar" />
        </div>
        <div class="header-right">
          <LanguageSwitcher />
          <el-button :icon="uiStore.isDark ? Sunny : Moon" text @click="uiStore.toggleTheme" aria-label="Toggle Theme" />
          <el-dropdown trigger="click">
            <el-button text>
              <el-avatar :size="32">{{ authStore.user?.email?.[0]?.toUpperCase() }}</el-avatar>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item disabled>{{ authStore.user?.email }}</el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">{{ t('common.logout') }}</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main-content">
        <RouterView v-slot="{ Component }">
          <Transition name="page-fade" mode="out-in" appear>
            <component :is="Component" />
          </Transition>
        </RouterView>
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.main-layout {
  height: 100vh;
  position: relative;
  overflow: hidden;
}

.sidebar {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-right: 1px solid rgba(226, 232, 240, 0.6);
  transition: width 0.3s;
  overflow: hidden;
  z-index: 10;
}

html.dark .sidebar {
  background: rgba(15, 23, 42, 0.85);
  border-color: rgba(51, 65, 85, 0.6);
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 800;
  background: linear-gradient(135deg, var(--pm-primary) 0%, var(--pm-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  border-bottom: 1px solid rgba(226, 232, 240, 0.6);
}

html.dark .logo {
  border-color: rgba(51, 65, 85, 0.6);
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgba(226, 232, 240, 0.6);
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  padding: 0 20px;
  z-index: 10;
}

html.dark .header {
  background: rgba(15, 23, 42, 0.85);
  border-color: rgba(51, 65, 85, 0.6);
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.main-content {
  background: transparent;
  padding: 20px;
  overflow-y: auto;
  position: relative;
  z-index: 1;
}

.el-menu {
  border-right: none;
  background: transparent;
}

/* Background Effects */
.layout-bg {
  position: fixed;
  inset: 0;
  z-index: 0;
  overflow: hidden;
  background: var(--pm-bg-page);
  pointer-events: none;
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.25;
  animation: float 12s ease-in-out infinite;
  will-change: transform;
}

.orb-1 {
  width: 500px;
  height: 500px;
  background: var(--pm-primary);
  top: -8%;
  right: -8%;
}

.orb-2 {
  width: 350px;
  height: 350px;
  background: var(--pm-secondary);
  bottom: -8%;
  left: -5%;
  animation-delay: 4s;
}

.orb-3 {
  width: 280px;
  height: 280px;
  background: #22d3ee;
  top: 35%;
  left: 25%;
  animation-delay: 8s;
}

html.dark .orb {
  opacity: 0.12;
}

.grid-overlay {
  position: absolute;
  inset: 0;
  background-image: linear-gradient(rgba(79, 70, 229, 0.03) 1px, transparent 1px),
                    linear-gradient(90deg, rgba(79, 70, 229, 0.03) 1px, transparent 1px);
  background-size: 60px 60px;
}

html.dark .grid-overlay {
  background-image: linear-gradient(rgba(99, 102, 241, 0.04) 1px, transparent 1px),
                    linear-gradient(90deg, rgba(99, 102, 241, 0.04) 1px, transparent 1px);
}

@keyframes float {
  0%, 100% { transform: translateY(0) scale(1); }
  50% { transform: translateY(-25px) scale(1.03); }
}

/* Page Transitions */
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}

.page-fade-enter-from {
  opacity: 0;
  transform: translateY(12px);
}

.page-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

@media (prefers-reduced-motion: reduce) {
  .orb { display: none; }
  .page-fade-enter-active,
  .page-fade-leave-active { transition: opacity 0.15s ease; }
  .page-fade-enter-from,
  .page-fade-leave-to { transform: none; }
}
</style>
