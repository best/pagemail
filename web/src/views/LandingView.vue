<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import DOMPurify from 'dompurify'
import { useSiteConfigStore } from '@/stores/siteConfig'
import PublicHeader from '@/components/common/PublicHeader.vue'
import { Monitor, Message, Timer, ArrowRight } from '@element-plus/icons-vue'

const { t } = useI18n()
const router = useRouter()
const siteConfig = useSiteConfigStore()

const features = computed(() => [
  { icon: Monitor, titleKey: 'landing.feature1Title', descKey: 'landing.feature1Desc' },
  { icon: Timer, titleKey: 'landing.feature2Title', descKey: 'landing.feature2Desc' },
  { icon: Message, titleKey: 'landing.feature3Title', descKey: 'landing.feature3Desc' }
])

const steps = computed(() => [
  { num: '01', titleKey: 'landing.step1Title', descKey: 'landing.step1Desc' },
  { num: '02', titleKey: 'landing.step2Title', descKey: 'landing.step2Desc' },
  { num: '03', titleKey: 'landing.step3Title', descKey: 'landing.step3Desc' }
])

const sanitizedSlogan = computed(() => {
  if (!siteConfig.siteSlogan) return ''
  return DOMPurify.sanitize(siteConfig.siteSlogan)
})

onMounted(() => siteConfig.fetchConfig())
</script>

<template>
  <div class="landing">
    <PublicHeader />

    <section class="hero">
      <div class="hero-bg">
        <div class="orb orb-1"></div>
        <div class="orb orb-2"></div>
        <div class="orb orb-3"></div>
        <div class="grid-overlay"></div>
      </div>
      <div class="container hero-content">
        <div class="hero-text">
          <h1 class="anim-up">{{ t('landing.heroTitle') }}<br><span class="gradient-text">{{ t('landing.heroHighlight') }}</span></h1>
          <p class="subtitle anim-up d1">{{ t('landing.heroSubtitle') }}</p>
          <div class="cta-group anim-up d2">
            <el-button type="primary" size="large" round @click="router.push('/register')">
              {{ t('landing.startFree') }} <el-icon class="el-icon--right"><ArrowRight /></el-icon>
            </el-button>
            <el-button size="large" round plain>{{ t('landing.learnMore') }}</el-button>
          </div>
        </div>
        <div class="hero-visual anim-fade d3">
          <div class="browser">
            <div class="browser-bar">
              <span class="dot red"></span>
              <span class="dot yellow"></span>
              <span class="dot green"></span>
              <span class="url">pagemail.app</span>
            </div>
            <div class="browser-body">
              <div class="mock-row">
                <div class="mock-card pulse"></div>
                <div class="mock-card pulse d1"></div>
                <div class="mock-card pulse d2"></div>
              </div>
              <div class="mock-table">
                <div class="mock-line"></div>
                <div class="mock-line"></div>
                <div class="mock-line"></div>
              </div>
            </div>
          </div>
          <div class="floating-badge badge-1">
            <el-icon><Monitor /></el-icon> {{ t('landing.captureReady') }}
          </div>
          <div class="floating-badge badge-2">
            <el-icon><Message /></el-icon> {{ t('landing.delivered') }}
          </div>
        </div>
      </div>
    </section>

    <section class="features">
      <div class="container">
        <h2 class="section-title">{{ t('landing.featuresTitle') }}</h2>
        <p class="section-subtitle">{{ t('landing.featuresSubtitle') }}</p>
        <div class="features-grid">
          <div v-for="(f, i) in features" :key="i" class="feature-card">
            <div class="feature-icon"><component :is="f.icon" /></div>
            <h3>{{ t(f.titleKey) }}</h3>
            <p>{{ t(f.descKey) }}</p>
          </div>
        </div>
      </div>
    </section>

    <section class="how-it-works">
      <div class="container">
        <h2 class="section-title">{{ t('landing.howItWorks') }}</h2>
        <div class="steps">
          <div v-for="(s, i) in steps" :key="i" class="step">
            <div class="step-num">{{ s.num }}</div>
            <h3>{{ t(s.titleKey) }}</h3>
            <p>{{ t(s.descKey) }}</p>
          </div>
        </div>
      </div>
    </section>

    <section class="cta-section">
      <div class="container cta-inner">
        <h2>{{ t('landing.ctaTitle') }}</h2>
        <p>{{ t('landing.ctaSubtitle') }}</p>
        <el-button size="large" round class="cta-btn" @click="router.push('/register')">{{ t('landing.ctaButton') }}</el-button>
      </div>
    </section>

    <footer class="footer">
      <div class="container footer-inner">
        <div class="footer-top">
          <div class="footer-brand">
            <h3>{{ siteConfig.siteName }}</h3>
          </div>
          <div class="footer-links">
            <a href="#">{{ t('landing.privacy') }}</a>
            <a href="#">{{ t('landing.terms') }}</a>
            <a href="#">{{ t('landing.contact') }}</a>
          </div>
        </div>
        <div class="footer-bottom">
          <div class="copyright">{{ siteConfig.copyright }}</div>
          <div v-if="sanitizedSlogan" class="footer-slogan" v-html="sanitizedSlogan"></div>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.landing {
  background: var(--pm-bg-page);
  color: var(--pm-text-heading);
  min-height: 100vh;
  overflow-x: hidden;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1.5rem;
}

/* Hero */
.hero {
  padding: 10rem 0 6rem;
  position: relative;
  overflow: hidden;
}

.hero-bg {
  position: absolute;
  inset: 0;
  z-index: 0;
  overflow: hidden;
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
  right: -10%;
  animation-delay: 0s;
}

.orb-2 {
  width: 400px;
  height: 400px;
  background: var(--pm-secondary);
  bottom: 0;
  left: -5%;
  animation-delay: 2s;
}

.orb-3 {
  width: 300px;
  height: 300px;
  background: #22d3ee;
  top: 40%;
  left: 30%;
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

.hero-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 4rem;
  align-items: center;
  position: relative;
  z-index: 1;
}

.hero-text h1 {
  font-size: 3.5rem;
  line-height: 1.1;
  font-weight: 800;
  margin-bottom: 1.5rem;
}

.gradient-text {
  background: linear-gradient(135deg, var(--pm-primary) 0%, var(--pm-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  font-size: 1.125rem;
  color: var(--pm-text-body);
  margin-bottom: 2rem;
  line-height: 1.7;
}

.cta-group {
  display: flex;
  gap: 1rem;
}

/* Browser Mockup */
.hero-visual {
  position: relative;
}

.browser {
  background: var(--pm-bg-card);
  border-radius: 12px;
  box-shadow: var(--pm-shadow-lg), 0 0 60px rgba(79,70,229,0.15);
  border: 1px solid var(--pm-border-color);
  overflow: hidden;
  transform: perspective(1000px) rotateY(-8deg) rotateX(5deg);
  transition: transform 0.5s;
}

.browser:hover {
  transform: perspective(1000px) rotateY(-2deg) rotateX(2deg);
}

.browser-bar {
  padding: 0.75rem 1rem;
  background: var(--pm-bg-elevated);
  border-bottom: 1px solid var(--pm-border-color);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.dot.red { background: #FF5F57; }
.dot.yellow { background: #FFBD2E; }
.dot.green { background: #28C840; }

.url {
  margin-left: 1rem;
  background: var(--pm-bg-page);
  padding: 0.25rem 1rem;
  border-radius: 4px;
  font-size: 0.8rem;
  color: var(--pm-text-muted);
  flex: 1;
}

.browser-body {
  padding: 1.5rem;
  min-height: 280px;
}

.mock-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.mock-card {
  height: 70px;
  background: linear-gradient(135deg, var(--pm-border-color) 0%, transparent 100%);
  border-radius: 8px;
}

.mock-table {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.mock-line {
  height: 16px;
  background: var(--pm-border-color);
  border-radius: 4px;
  opacity: 0.5;
}

.mock-line:nth-child(2) { width: 80%; }
.mock-line:nth-child(3) { width: 60%; }

.pulse {
  animation: pulse 2s ease-in-out infinite;
}

.pulse.d1 { animation-delay: 0.3s; }
.pulse.d2 { animation-delay: 0.6s; }

@keyframes pulse {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 1; }
}

/* Floating Badges */
.floating-badge {
  position: absolute;
  background: var(--pm-bg-card);
  padding: 0.75rem 1.25rem;
  border-radius: 50px;
  box-shadow: var(--pm-shadow-md);
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  font-weight: 600;
  border: 1px solid var(--pm-border-color);
  animation: badge-float 3s ease-in-out infinite;
}

.badge-1 {
  top: 10%;
  right: -10%;
  color: var(--pm-primary);
}

.badge-2 {
  bottom: 15%;
  left: -5%;
  color: var(--el-color-success);
  animation-delay: 1.5s;
}

@keyframes badge-float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

/* Features */
.features {
  padding: 6rem 0;
  background: var(--pm-bg-elevated);
}

.section-title {
  font-size: 2.25rem;
  font-weight: 700;
  text-align: center;
  margin-bottom: 0.5rem;
}

.section-subtitle {
  text-align: center;
  color: var(--pm-text-body);
  font-size: 1.125rem;
  margin-bottom: 3rem;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 2rem;
}

.feature-card {
  background: var(--pm-bg-card);
  padding: 2rem;
  border-radius: 16px;
  border: 1px solid var(--pm-border-color);
  transition: transform 0.3s, box-shadow 0.3s;
}

.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--pm-shadow-md);
}

.feature-icon {
  width: 48px;
  height: 48px;
  background: var(--el-color-primary-light-9);
  color: var(--pm-primary);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1.25rem;
  font-size: 1.5rem;
}

html.dark .feature-icon {
  background: rgba(79,70,229,0.15);
}

.feature-card h3 {
  font-size: 1.125rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.feature-card p {
  color: var(--pm-text-body);
  line-height: 1.6;
  font-size: 0.9375rem;
}

/* How It Works */
.how-it-works {
  padding: 6rem 0;
}

.steps {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2rem;
  margin-top: 3rem;
}

.step {
  position: relative;
  padding: 1.5rem;
}

.step-num {
  font-size: 4rem;
  font-weight: 900;
  color: var(--el-color-primary-light-7);
  line-height: 1;
  margin-bottom: 1rem;
}

html.dark .step-num {
  color: rgba(79,70,229,0.15);
}

.step h3 {
  font-size: 1.25rem;
  margin-bottom: 0.5rem;
}

.step p {
  color: var(--pm-text-body);
}

/* CTA Section */
.cta-section {
  padding: 5rem 0;
  background: linear-gradient(135deg, var(--pm-primary-dark) 0%, var(--pm-primary) 50%, var(--pm-primary-light) 100%);
  text-align: center;
}

.cta-inner h2 {
  font-size: 2.25rem;
  color: white;
  margin-bottom: 0.75rem;
}

.cta-inner p {
  font-size: 1.125rem;
  color: rgba(255,255,255,0.85);
  margin-bottom: 2rem;
}

.cta-btn {
  background: white !important;
  color: var(--pm-primary) !important;
  border: none !important;
  font-weight: 600;
  padding: 1.25rem 2.5rem !important;
}

.cta-btn:hover {
  background: var(--pm-bg-page) !important;
}

/* Footer */
.footer {
  padding: 3rem 0 2rem;
  border-top: 1px solid var(--pm-border-color);
}

.footer-inner {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.footer-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-brand h3 {
  font-size: 1.25rem;
  font-weight: 700;
  margin: 0;
}

.footer-links {
  display: flex;
  gap: 2rem;
}

.footer-links a {
  color: var(--pm-text-body);
  text-decoration: none;
  font-size: 0.875rem;
  transition: color 0.2s;
}

.footer-links a:hover {
  color: var(--pm-primary);
}

.footer-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 1.5rem;
  border-top: 1px solid var(--pm-border-color);
  font-size: 0.8125rem;
  color: var(--pm-text-muted);
}

.footer-slogan :deep(a) {
  color: var(--pm-text-muted);
  text-decoration: none;
  transition: color 0.2s;
}

.footer-slogan :deep(a:hover) {
  color: var(--pm-primary);
}

/* Animations */
.anim-up {
  animation: slideUp 0.8s ease-out forwards;
  opacity: 0;
}

.anim-fade {
  animation: fadeIn 1s ease-out forwards;
  opacity: 0;
}

.d1 { animation-delay: 0.15s; }
.d2 { animation-delay: 0.3s; }
.d3 { animation-delay: 0.45s; }

@keyframes slideUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes fadeIn {
  to { opacity: 1; }
}

/* Responsive */
@media (max-width: 960px) {
  .hero {
    padding: 7rem 0 4rem;
  }

  .hero-content {
    grid-template-columns: 1fr;
    text-align: center;
  }

  .hero-text h1 {
    font-size: 2.5rem;
  }

  .cta-group {
    justify-content: center;
  }

  .hero-visual {
    display: none;
  }

  .steps {
    grid-template-columns: 1fr;
  }

  .footer-top {
    flex-direction: column;
    gap: 1.5rem;
    text-align: center;
  }

  .footer-bottom {
    flex-direction: column;
    gap: 0.75rem;
    text-align: center;
    justify-content: center;
  }
}
</style>
