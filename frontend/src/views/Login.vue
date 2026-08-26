<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { apiLogin } from '../services/api'

const router = useRouter()
const route = useRoute()
const { login } = useAuth()

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const loading = ref(false)
const errorMsg = ref('')

// Fallback mock users (when backend is unreachable)
const mockUsers = [
  { email: 'donatur@test.com', password: 'donatur123', name: 'Ahmad Saputra', role: 'donatur' },
  { email: 'penggalang@test.com', password: 'penggalang123', name: 'Yayasan Peduli', role: 'penggalang' },
]

const handleLogin = async () => {
  errorMsg.value = ''
  loading.value = true

  try {
    // Try real backend API first
    const data = await apiLogin(email.value, password.value)

    // Login successful via API
    login(data.user, data.token)
    loading.value = false

    // Redirect based on query or role
    const redirectPath = route.query.redirect
    if (redirectPath) {
      router.push(redirectPath)
    } else if (data.user.role === 'penggalang') {
      router.push('/penggalang')
    } else if (data.user.role === 'donatur') {
      router.push('/donatur')
    } else {
      router.push('/')
    }

  } catch (err) {
    // If it's a server auth error, show the error
    if (err.message && !err.message.includes('fetch') && !err.message.includes('NetworkError')) {
      errorMsg.value = err.message
      loading.value = false
      return
    }

    // Backend unreachable — fallback to mock users for demo
    console.warn('[Login] Backend unreachable, falling back to mock:', err.message)
    const found = mockUsers.find(u => u.email === email.value && u.password === password.value)

    if (!found) {
      errorMsg.value = 'Email atau password salah. Coba: donatur@test.com / donatur123 atau penggalang@test.com / penggalang123'
      loading.value = false
      return
    }

    login({ name: found.name, email: found.email, role: found.role }, 'mock-jwt-token-' + Date.now())
    loading.value = false

    const redirectPath = route.query.redirect
    if (redirectPath) {
      router.push(redirectPath)
    } else if (found.role === 'penggalang') {
      router.push('/penggalang')
    } else if (found.role === 'donatur') {
      router.push('/donatur')
    } else {
      router.push('/')
    }
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-container">
      <div class="auth-visual">
        <div class="auth-visual-inner">
          <h2>Selamat Datang<br />Kembali 👋</h2>
          <p>Masuk untuk mengelola campaign, memantau donasi, dan melihat laporan penyaluran dana.</p>
          <div class="visual-features">
            <div class="vf-item"><span>🔒</span> Sistem Escrow Aman</div>
            <div class="vf-item"><span>📊</span> Transparansi Dana 100%</div>
            <div class="vf-item"><span>✅</span> Verifikasi Admin</div>
          </div>
        </div>
      </div>

      <div class="auth-form-wrapper">
        <div class="auth-form-inner">
          <router-link to="/" class="auth-logo">
            <span class="logo-text">Open<span class="logo-accent">Donasi</span></span>
          </router-link>
          <h1>Masuk ke Akun</h1>
          <p class="auth-subtitle">Masukkan email dan password Anda</p>

          <div class="error-box" v-if="errorMsg">
            <p>⚠️ {{ errorMsg }}</p>
          </div>

          <form @submit.prevent="handleLogin" class="auth-form" id="login-form">
            <div class="form-group">
              <label for="email">Email</label>
              <input type="email" id="email" v-model="email" placeholder="nama@email.com" required />
            </div>

            <div class="form-group">
              <label for="password">Password</label>
              <div class="input-password">
                <input :type="showPassword ? 'text' : 'password'" id="password" v-model="password" placeholder="••••••••" required />
                <button type="button" class="toggle-pw" @click="showPassword = !showPassword">
                  {{ showPassword ? '🙈' : '👁️' }}
                </button>
              </div>
            </div>

            <button type="submit" class="btn-submit" id="btn-submit-login" :disabled="loading">
              <span v-if="loading" class="spinner"></span>
              <span v-else>Masuk</span>
            </button>
          </form>


          <p class="auth-footer-text">
            Belum punya akun? <router-link to="/register">Daftar Sekarang</router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-body);
  padding: 24px;
  padding-top: 96px;
}

.auth-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  max-width: 960px;
  width: 100%;
  background: white;
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: var(--shadow-xl);
  min-height: 560px;
}

.auth-visual {
  background: linear-gradient(135deg, var(--primary), #0077B6);
  padding: 48px;
  display: flex;
  align-items: center;
}

.auth-visual-inner h2 { font-family: var(--font-heading); font-size: 32px; font-weight: 800; color: white; line-height: 1.2; margin-bottom: 16px; }
.auth-visual-inner p { color: rgba(255,255,255,0.8); font-size: 15px; line-height: 1.7; margin-bottom: 32px; }
.visual-features { display: flex; flex-direction: column; gap: 12px; }
.vf-item { display: flex; align-items: center; gap: 10px; font-size: 14px; color: rgba(255,255,255,0.9); font-weight: 500; }
.vf-item span { font-size: 18px; }

.auth-form-wrapper { padding: 48px; display: flex; align-items: center; }
.auth-form-inner { width: 100%; }
.auth-logo { display: inline-block; text-decoration: none; margin-bottom: 32px; }
.logo-text { font-family: var(--font-heading); font-weight: 800; font-size: 22px; color: var(--text-heading); }
.logo-accent { color: var(--primary); }
.auth-form-inner h1 { font-family: var(--font-heading); font-size: 24px; font-weight: 700; color: var(--text-heading); margin-bottom: 6px; }
.auth-subtitle { color: var(--text-muted); font-size: 14px; margin-bottom: 24px; }

.error-box {
  background: var(--danger-light);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: var(--radius-sm);
  padding: 12px 16px;
  margin-bottom: 16px;
}
.error-box p { font-size: 13px; color: var(--danger); margin: 0; }

.form-group { margin-bottom: 20px; }
.form-group label { display: block; font-size: 13px; font-weight: 600; color: var(--text-heading); margin-bottom: 6px; }
.form-group input {
  width: 100%; padding: 12px 16px; border: 1px solid var(--border); border-radius: var(--radius-sm);
  font-size: 14px; color: var(--text-heading); background: var(--bg-body); transition: all var(--transition-fast);
}
.form-group input:focus { outline: none; border-color: var(--primary); box-shadow: 0 0 0 3px var(--primary-light); background: white; }
.form-group input::placeholder { color: var(--text-muted); }

.input-password { position: relative; }
.input-password input { padding-right: 48px; }
.toggle-pw { position: absolute; right: 12px; top: 50%; transform: translateY(-50%); background: none; border: none; cursor: pointer; font-size: 18px; padding: 4px; }

.btn-submit {
  width: 100%; padding: 14px 24px; background: var(--primary); color: white;
  font-family: var(--font-body); font-weight: 600; font-size: 15px;
  border: none; border-radius: var(--radius-sm); cursor: pointer;
  transition: all var(--transition-fast); display: flex; align-items: center; justify-content: center; gap: 8px; margin-top: 8px;
}
.btn-submit:hover:not(:disabled) { background: var(--primary-hover); box-shadow: 0 4px 16px var(--primary-glow); }
.btn-submit:disabled { opacity: 0.7; cursor: not-allowed; }
.spinner { width: 18px; height: 18px; border: 2px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 0.6s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.auth-footer-text { text-align: center; font-size: 14px; color: var(--text-muted); margin-top: 20px; }
.auth-footer-text a { color: var(--primary); font-weight: 600; }

@media (max-width: 768px) {
  .auth-container { grid-template-columns: 1fr; }
  .auth-visual { display: none; }
  .auth-form-wrapper { padding: 32px 24px; }
}
</style>
