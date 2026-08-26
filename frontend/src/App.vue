<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from './composables/useAuth'

const router = useRouter()
const { isLoggedIn, user, role, logout } = useAuth()

const mobileMenuOpen = ref(false)
const scrolled = ref(false)

const handleScroll = () => {
  scrolled.value = window.scrollY > 60
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

const handleLogout = () => {
  logout()
  mobileMenuOpen.value = false
  router.push('/')
}

const goToDashboard = () => {
  mobileMenuOpen.value = false
  if (role.value === 'penggalang') router.push('/penggalang')
  else if (role.value === 'donatur') router.push('/donatur')
}
</script>

<template>
  <div class="app-shell">
    
    <header class="navbar" :class="{ 'navbar-scrolled': scrolled }" id="main-navbar">
      <div class="container navbar-inner">
        <router-link to="/" class="logo" id="logo-link">
          <span class="logo-icon">💙</span>
          <span class="logo-text">Open<span class="logo-accent">Donasi</span></span>
        </router-link>

        <nav class="nav-links" :class="{ 'nav-open': mobileMenuOpen }">
          <template v-if="!isLoggedIn || role !== 'penggalang'">
            <router-link to="/" @click="mobileMenuOpen = false">Beranda</router-link>
            <router-link to="/campaigns" @click="mobileMenuOpen = false">Campaign</router-link>
          </template>
        </nav>

        <div class="nav-actions">
          
          <template v-if="!isLoggedIn">
            <router-link to="/login" class="btn-nav-login" id="btn-login">Masuk</router-link>
            <router-link to="/register" class="btn-nav-register" id="btn-register">Daftar</router-link>
          </template>

          
          <template v-else>
            <button class="btn-nav-dashboard" @click="goToDashboard" id="btn-dashboard">
              👤 {{ user?.name || 'User' }}
            </button>
            <button class="btn-nav-logout" @click="handleLogout" id="btn-logout">Keluar</button>
          </template>

          <button class="mobile-toggle" @click="mobileMenuOpen = !mobileMenuOpen" id="btn-mobile-toggle" aria-label="Toggle menu">
            <span :class="{ 'open': mobileMenuOpen }"></span>
          </button>
        </div>
      </div>
    </header>


    <main class="main-content">
      <router-view />
    </main>


    <footer class="footer" id="main-footer">
      <div class="container footer-inner">
        <div class="footer-brand">
          <span class="footer-logo-text">Open<span class="logo-accent">Donasi</span></span>
          <p>Platform crowdfunding transparan untuk membantu korban bencana alam di Indonesia.</p>
        </div>
        <div class="footer-links">
          <h4>Navigasi</h4>
          <router-link to="/">Beranda</router-link>
          <router-link to="/campaigns">Campaign</router-link>
          <router-link to="/login">Masuk</router-link>
        </div>
        <div class="footer-links">
          <h4>Informasi</h4>
          <a href="#">Tentang Kami</a>
          <a href="#">Kebijakan Privasi</a>
          <a href="#">Syarat & Ketentuan</a>
        </div>
      </div>
      <div class="footer-bottom">
        <div class="container">
          <p>&copy; 2026 OpenDonasi — Sistem Open Donasi untuk Indonesia.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>

.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  transition: all 0.35s ease;
  background: transparent;
}

.navbar-scrolled {
  background: rgba(255, 255, 255, 0.97);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border);
  box-shadow: var(--shadow-xs);
}

.navbar-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 72px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
}

.logo-icon { font-size: 24px; }

.logo-text {
  font-family: var(--font-heading);
  font-weight: 800;
  font-size: 22px;
  color: var(--text-heading);
  transition: color 0.35s ease;
}

.navbar-scrolled .logo-text {
  color: var(--text-heading);
}

.logo-accent { color: var(--primary); }

.nav-links {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nav-links a {
  font-family: var(--font-body);
  font-weight: 500;
  font-size: 15px;
  color: var(--text-heading);
  text-decoration: none;
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  transition: all 0.25s ease;
}

.navbar-scrolled .nav-links a {
  color: var(--text-heading);
}

.nav-links a:hover {
  color: var(--primary);
  background: var(--primary-light);
}

.navbar-scrolled .nav-links a:hover {
  color: var(--primary);
  background: var(--primary-light);
}

.nav-links a.router-link-exact-active {
  color: var(--primary);
  background: var(--primary-light);
  font-weight: 600;
}

.navbar-scrolled .nav-links a.router-link-exact-active {
  color: var(--primary);
  background: var(--primary-light);
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.btn-nav-login {
  font-weight: 600;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.9);
  padding: 8px 20px;
  border-radius: var(--radius-sm);
  text-decoration: none;
  transition: all 0.25s ease;
}

.navbar-scrolled .btn-nav-login {
  color: var(--primary);
}

.btn-nav-login:hover {
  background: rgba(255, 255, 255, 0.12);
}

.navbar-scrolled .btn-nav-login:hover {
  background: var(--primary-light);
}

.btn-nav-register {
  font-weight: 600;
  font-size: 14px;
  color: var(--text-heading);
  background: white;
  padding: 8px 20px;
  border-radius: var(--radius-sm);
  text-decoration: none;
  transition: all 0.25s ease;
}

.navbar-scrolled .btn-nav-register {
  color: white;
  background: var(--primary);
}

.btn-nav-register:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  color: var(--text-heading);
}

.navbar-scrolled .btn-nav-register:hover {
  background: var(--primary-hover);
  box-shadow: 0 4px 12px var(--primary-glow);
  color: white;
}

.btn-nav-dashboard,
.btn-nav-logout {
  font-family: var(--font-body);
  font-size: 14px;
  font-weight: 600;
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  border: none;
  cursor: pointer;
  transition: all 0.25s ease;
}

.btn-nav-dashboard {
  background: var(--bg-surface);
  color: var(--text-heading);
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 6px;
  border: 1px solid var(--border);
}

.navbar-scrolled .btn-nav-dashboard {
  background: var(--bg-surface-hover);
  color: var(--text-heading);
}

.btn-nav-dashboard:hover {
  background: var(--border-light);
  color: var(--primary);
  border-color: var(--primary);
}

.navbar-scrolled .btn-nav-dashboard:hover {
  background: var(--border-light);
  color: var(--primary);
}

.btn-nav-logout {
  background: transparent;
  color: var(--text-heading);
  border: 1px solid var(--border);
  display: flex;
  flex-direction: row;
  align-items: center;
}

.navbar-scrolled .btn-nav-logout {
  color: var(--text-heading);
  border-color: var(--border);
}

.btn-nav-logout:hover {
  background: var(--danger-light);
  color: var(--danger);
  border-color: var(--danger);
}

.navbar-scrolled .btn-nav-logout:hover {
  background: var(--danger-light);
  color: var(--danger);
  border-color: var(--danger);
}

.mobile-toggle {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  width: 32px;
  height: 32px;
  position: relative;
}

.mobile-toggle span,
.mobile-toggle span::before,
.mobile-toggle span::after {
  display: block;
  width: 22px;
  height: 2px;
  background: white;
  border-radius: 2px;
  position: absolute;
  left: 5px;
  transition: all 0.3s ease;
}

.navbar-scrolled .mobile-toggle span,
.navbar-scrolled .mobile-toggle span::before,
.navbar-scrolled .mobile-toggle span::after {
  background: var(--text-heading);
}

.mobile-toggle span { top: 15px; }
.mobile-toggle span::before { content: ''; top: -7px; }
.mobile-toggle span::after { content: ''; top: 7px; }
.mobile-toggle span.open { background: transparent; }
.mobile-toggle span.open::before { transform: rotate(45deg); top: 0; }
.mobile-toggle span.open::after { transform: rotate(-45deg); top: 0; }

.main-content {
  min-height: 100vh;
  padding-top: 80px; 
}

.footer {
  background: var(--text-heading);
  color: rgba(255, 255, 255, 0.7);
  padding-top: 60px;
}

.footer-inner {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  gap: 40px;
  padding-bottom: 40px;
}

.footer-logo-text {
  font-family: var(--font-heading);
  font-weight: 800;
  font-size: 20px;
  color: white;
  display: block;
  margin-bottom: 12px;
}

.footer-brand p {
  font-size: 14px;
  line-height: 1.7;
  max-width: 300px;
}

.footer-links h4 {
  font-family: var(--font-heading);
  color: white;
  font-size: 14px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 16px;
}

.footer-links a {
  display: block;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  padding: 4px 0;
  transition: color 0.2s ease;
}

.footer-links a:hover { color: var(--primary); }

.footer-bottom {
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px 0;
}

.footer-bottom p {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.4);
  text-align: center;
}


@media (max-width: 768px) {
  .nav-links {
    position: fixed;
    top: 72px;
    left: 0;
    right: 0;
    background: white;
    flex-direction: column;
    padding: 16px 24px;
    border-bottom: 1px solid var(--border);
    transform: translateY(-120%);
    opacity: 0;
    transition: all 0.3s ease;
    z-index: 99;
  }

  .nav-links.nav-open {
    transform: translateY(0);
    opacity: 1;
  }

  .nav-links a {
    width: 100%;
    padding: 12px 16px;
    color: var(--text-body) !important;
  }

  .mobile-toggle { display: block; }
  .btn-nav-login, .btn-nav-register { display: none; }
  .btn-nav-dashboard, .btn-nav-logout { font-size: 12px; padding: 6px 12px; }

  .footer-inner {
    grid-template-columns: 1fr;
    gap: 24px;
  }
}
</style>
