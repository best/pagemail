<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useSiteConfigStore } from '@/stores/siteConfig'
import PublicHeader from '@/components/common/PublicHeader.vue'

const { t } = useI18n()
const siteConfig = useSiteConfigStore()

onMounted(() => siteConfig.fetchConfig())
</script>

<template>
  <div class="auth-layout">
    <PublicHeader solid />
    <div class="auth-content">
      <div class="brand-side">
        <div class="brand-content">
          <h1 class="brand-logo">{{ siteConfig.siteName }}</h1>
          <p class="brand-tagline" v-html="t('auth.brandTagline').replace('\n', '<br>')"></p>
        </div>
        <div class="orb orb-1"></div>
        <div class="orb orb-2"></div>
        <div class="orb orb-3"></div>
        <div class="grid-overlay"></div>
      </div>
      <div class="form-side">
        <div class="form-container">
          <RouterView v-slot="{ Component }">
            <Transition name="fade-slide" appear>
              <component :is="Component" />
            </Transition>
          </RouterView>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  width: 100%;
}

.auth-content {
  display: flex;
  flex: 1;
  margin-top: 73px;
}

.brand-side {
  flex: 0 0 45%;
  background-color: var(--pm-sidebar-bg);
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.brand-content {
  position: relative;
  z-index: 2;
  text-align: center;
  padding: 2rem;
}

.brand-logo {
  font-size: 3rem;
  font-weight: 700;
  color: #fff;
  margin: 0 0 1rem;
  letter-spacing: -0.02em;
}

.brand-tagline {
  font-size: 1.25rem;
  color: rgba(255, 255, 255, 0.85);
  line-height: 1.6;
  margin: 0;
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: float 8s ease-in-out infinite;
  will-change: transform;
}

.orb-1 {
  width: 600px;
  height: 600px;
  background: var(--pm-primary);
  top: -20%;
  right: -20%;
  animation-delay: 0s;
}

.orb-2 {
  width: 400px;
  height: 400px;
  background: var(--pm-secondary);
  bottom: 0;
  left: -10%;
  animation-delay: 2s;
}

.orb-3 {
  width: 300px;
  height: 300px;
  background: #22d3ee;
  top: 40%;
  left: 20%;
  animation-delay: 4s;
}

html.dark .orb { opacity: 0.2; }

.grid-overlay {
  position: absolute;
  inset: 0;
  background-image: linear-gradient(rgba(79,70,229,0.03) 1px, transparent 1px),
                    linear-gradient(90deg, rgba(79,70,229,0.03) 1px, transparent 1px);
  background-size: 60px 60px;
}

html.dark .grid-overlay {
  background-image: linear-gradient(rgba(99,102,241,0.05) 1px, transparent 1px),
                    linear-gradient(90deg, rgba(99,102,241,0.05) 1px, transparent 1px);
}

@keyframes float {
  0%, 100% { transform: translateY(0) scale(1); }
  50% { transform: translateY(-30px) scale(1.05); }
}

.form-side {
  flex: 1;
  background: var(--pm-bg-card);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.form-container {
  width: 100%;
  max-width: 400px;
}

.fade-slide-enter-active {
  transition: opacity 0.5s ease, transform 0.5s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

@media (max-width: 768px) {
  .auth-content {
    margin-top: 56px;
  }

  .brand-side {
    display: none;
  }

  .form-side {
    padding: 1.5rem;
  }
}
</style>
