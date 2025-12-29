<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, useRouter, useRoute } from 'vue-router'
import { useUiStore } from '@/stores/ui'
import { useAuthStore } from '@/stores/auth'
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

const router = useRouter()
const route = useRoute()
const uiStore = useUiStore()
const authStore = useAuthStore()

const menuItems = computed(() => {
  const items = [
    { index: '/dashboard', title: 'Dashboard', icon: House },
    { index: '/tasks', title: 'Tasks', icon: Document },
    { index: '/settings', title: 'Settings', icon: Setting }
  ]

  if (authStore.isAdmin) {
    items.push(
      { index: '/admin/users', title: 'Users', icon: User },
      { index: '/admin/system', title: 'System', icon: DataAnalysis },
      { index: '/admin/audit', title: 'Audit Logs', icon: Memo }
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
          <template #title>{{ item.title }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-button :icon="uiStore.sidebarCollapsed ? Expand : Fold" text @click="uiStore.toggleSidebar" />
        </div>
        <div class="header-right">
          <el-button :icon="uiStore.isDark ? Sunny : Moon" text @click="uiStore.toggleTheme" />
          <el-dropdown trigger="click">
            <el-button text>
              <el-avatar :size="32">{{ authStore.user?.email?.[0]?.toUpperCase() }}</el-avatar>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item disabled>{{ authStore.user?.email }}</el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">Logout</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main-content">
        <RouterView />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.main-layout {
  height: 100vh;
}

.sidebar {
  background: var(--el-bg-color);
  border-right: 1px solid var(--el-border-color-light);
  transition: width 0.3s;
  overflow: hidden;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--el-color-primary);
  border-bottom: 1px solid var(--el-border-color-light);
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--el-border-color-light);
  background: var(--el-bg-color);
  padding: 0 20px;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.main-content {
  background: var(--el-bg-color-page);
  padding: 20px;
  overflow-y: auto;
}

.el-menu {
  border-right: none;
}
</style>
