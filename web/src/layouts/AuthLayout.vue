<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="auth-layout">
    <div class="brand-side">
      <div class="brand-content">
        <h1 class="brand-logo">Pagemail</h1>
        <p class="brand-tagline">Capture the web,<br>deliver with ease.</p>
      </div>
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
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
</template>

<style scoped>
.auth-layout {
  display: flex;
  min-height: 100vh;
  width: 100%;
}

.brand-side {
  flex: 0 0 45%;
  background: linear-gradient(135deg, var(--pm-sidebar-bg) 0%, var(--pm-sidebar-light) 100%);
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

.shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.15;
}

.shape-1 {
  width: 400px;
  height: 400px;
  background: var(--pm-primary);
  top: -15%;
  left: -10%;
}

.shape-2 {
  width: 300px;
  height: 300px;
  background: #fff;
  bottom: -10%;
  right: -5%;
}

.shape-3 {
  width: 200px;
  height: 200px;
  background: var(--pm-secondary);
  top: 50%;
  right: 20%;
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
  .brand-side {
    display: none;
  }

  .form-side {
    padding: 1.5rem;
  }
}
</style>
