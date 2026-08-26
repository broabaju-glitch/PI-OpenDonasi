<template>
  <div class="login-page">
    <!-- Background decorative blobs -->
    <div class="bg-blob bg-blob--1" aria-hidden="true"></div>
    <div class="bg-blob bg-blob--2" aria-hidden="true"></div>

    <div class="login-card" role="main">
      <!-- Brand Header -->
      <div class="brand-header">
        <div class="brand-icon" aria-hidden="true">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 21.593c-5.63-5.539-11-10.297-11-14.402 0-3.791 3.068-5.191 5.281-5.191 1.312 0 4.151.501 5.719 4.457 1.59-3.968 4.464-4.447 5.726-4.447 2.54 0 5.274 1.621 5.274 5.181 0 4.069-5.136 8.625-11 14.402z"
              fill="#3B82F6"/>
          </svg>
        </div>
        <h1 class="brand-name">Open<span class="brand-accent">Donasi</span></h1>
      </div>

      <!-- Card Body -->
      <div class="card">
        <div class="card-top">
          <div class="shield-badge" aria-hidden="true">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z" fill="#3B82F6" opacity="0.15"/>
              <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z" stroke="#3B82F6" stroke-width="1.5" stroke-linejoin="round"/>
              <path d="M9 12l2 2 4-4" stroke="#3B82F6" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <div class="card-title-group">
            <h2 class="card-title">Admin Portal</h2>
            <p class="card-subtitle">Gunakan kredensial admin untuk melanjutkan.</p>
          </div>
        </div>

        <!-- Error Alert -->
        <transition name="fade-slide">
          <div v-if="errorMessage" class="error-alert" role="alert" aria-live="polite">
            <svg class="error-icon" width="16" height="16" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5"/>
              <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
              <circle cx="12" cy="16" r="0.5" fill="currentColor" stroke="currentColor" stroke-width="1"/>
            </svg>
            <span>{{ errorMessage }}</span>
          </div>
        </transition>

        <!-- Login Form -->
        <form id="admin-login-form" class="login-form" @submit.prevent="handleLogin" novalidate>

          <!-- Email Field -->
          <div class="form-group">
            <label for="admin-email" class="form-label">Alamat Email</label>
            <div class="input-wrapper" :class="{ 'input-wrapper--focused': emailFocused }">
              <span class="input-icon" aria-hidden="true">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                  <rect x="2" y="4" width="20" height="16" rx="3" stroke="currentColor" stroke-width="1.5"/>
                  <path d="m2 7 8.586 5.586a2 2 0 0 0 2.828 0L22 7" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>
                </svg>
              </span>
              <input
                id="admin-email"
                v-model="email"
                type="email"
                class="form-input"
                placeholder="admin@opendonasi.com"
                autocomplete="email"
                required
                :disabled="isLoading"
                @focus="emailFocused = true"
                @blur="emailFocused = false"
              />
            </div>
          </div>

          <!-- Password Field -->
          <div class="form-group">
            <label for="admin-password" class="form-label">Password</label>
            <div class="input-wrapper" :class="{ 'input-wrapper--focused': passwordFocused }">
              <span class="input-icon" aria-hidden="true">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                  <rect x="3" y="11" width="18" height="11" rx="2" stroke="currentColor" stroke-width="1.5"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
              </span>
              <input
                id="admin-password"
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                class="form-input"
                placeholder="••••••••"
                autocomplete="current-password"
                required
                :disabled="isLoading"
                @focus="passwordFocused = true"
                @blur="passwordFocused = false"
              />
              <button
                type="button"
                class="toggle-password"
                :aria-label="showPassword ? 'Sembunyikan password' : 'Tampilkan password'"
                @click="showPassword = !showPassword"
              >
                <svg v-if="!showPassword" width="16" height="16" viewBox="0 0 24 24" fill="none">
                  <path d="M2 12s3.636-7 10-7 10 7 10 7-3.636 7-10 7-10-7-10-7z" stroke="currentColor" stroke-width="1.5"/>
                  <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="1.5"/>
                </svg>
                <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-6.364 0-10-7-10-7a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c6.364 0 10 7 10 7a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                  <line x1="1" y1="1" x2="23" y2="23" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
              </button>
            </div>
          </div>

          <!-- Submit Button -->
          <button
            id="admin-login-submit"
            type="submit"
            class="submit-btn"
            :class="{ 'submit-btn--loading': isLoading }"
            :disabled="isLoading || !email || !password"
          >
            <span v-if="!isLoading" class="btn-content">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>
              </svg>
              Masuk ke Admin Panel
            </span>
            <span v-else class="btn-content btn-loading-content">
              <span class="spinner" aria-hidden="true"></span>
              Memproses...
            </span>
          </button>
        </form>
      </div>

      <!-- Footer -->
      <p class="login-footer">&copy; {{ currentYear }} OpenDonasi &mdash; Sistem Administrasi Internal</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { authApi } from '@/services/api'

// ─── Router & Composables ───
const router = useRouter()
const { login } = useAuth()

// ─── Reactive State ───
const email        = ref('')
const password     = ref('')
const errorMessage = ref('')
const isLoading    = ref(false)
const showPassword = ref(false)
const emailFocused    = ref(false)
const passwordFocused = ref(false)

// ─── Computed ───
const currentYear = computed(() => new Date().getFullYear())

// ─── Login Handler ───
async function handleLogin() {
  // Clear previous error
  errorMessage.value = ''

  // Basic client-side validation
  if (!email.value || !password.value) {
    errorMessage.value = 'Email dan password wajib diisi.'
    return
  }

  isLoading.value = true

  try {
    // Uses authApi → proxied through Vite to http://localhost:8080/api/login
    const data = await authApi.login(email.value.trim(), password.value)

    // Validate that the authenticated user has the admin role
    if (data.user?.role !== 'admin') {
      errorMessage.value = 'Akses ditolak. Hanya akun admin yang dapat masuk ke panel ini.'
      return
    }

    // Persist auth state via composable (stores user + JWT in localStorage)
    login(data.user, data.token)

    // Redirect to admin dashboard
    await router.push('/admin')

  } catch (err) {
    // err.message is set by api.js from the server's JSON error field
    if (err.message && !err.message.startsWith('Server error')) {
      errorMessage.value = err.message
    } else {
      errorMessage.value = 'Tidak dapat terhubung ke server. Pastikan backend berjalan.'
    }
    console.error('[AdminLogin] Login error:', err)
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* ─── Google Font Import ─── */
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');

/* ─── CSS Variables ─── */
* {
  box-sizing: border-box;
}

/* ─── Page Layout ─── */
.login-page {
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  min-height: 100vh;
  background-color: #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 16px;
  position: relative;
  overflow: hidden;
}

/* ─── Decorative Background Blobs ─── */
.bg-blob {
  position: fixed;
  border-radius: 50%;
  filter: blur(80px);
  pointer-events: none;
  z-index: 0;
}

.bg-blob--1 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.12) 0%, transparent 70%);
  top: -100px;
  right: -80px;
}

.bg-blob--2 {
  width: 350px;
  height: 350px;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.08) 0%, transparent 70%);
  bottom: -80px;
  left: -60px;
}

/* ─── Login Card Container ─── */
.login-card {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 420px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

/* ─── Brand Header ─── */
.brand-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-icon {
  width: 40px;
  height: 40px;
  background: #fff;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 1px 3px rgba(0,0,0,0.08), 0 4px 12px rgba(59,130,246,0.12);
}

.brand-name {
  font-size: 1.5rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.5px;
  margin: 0;
}

.brand-accent {
  color: #3b82f6;
}

/* ─── Main Card ─── */
.card {
  width: 100%;
  background: #ffffff;
  border-radius: 20px;
  box-shadow:
    0 0 0 1px rgba(0, 0, 0, 0.04),
    0 4px 6px -1px rgba(0, 0, 0, 0.06),
    0 16px 32px -4px rgba(0, 0, 0, 0.08);
  padding: 36px 36px 32px;
}

/* ─── Card Top Row ─── */
.card-top {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  margin-bottom: 28px;
}

.shield-badge {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  background: #eff6ff;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 2px;
}

.card-title-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.card-title {
  font-size: 1.2rem;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
  line-height: 1.3;
}

.card-subtitle {
  font-size: 0.8rem;
  color: #94a3b8;
  margin: 0;
  line-height: 1.4;
}

/* ─── Error Alert ─── */
.error-alert {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 10px;
  padding: 10px 14px;
  margin-bottom: 20px;
  color: #dc2626;
  font-size: 0.8rem;
  font-weight: 500;
  line-height: 1.4;
}

.error-icon {
  flex-shrink: 0;
  color: #ef4444;
}

/* ─── Transition ─── */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}

.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

/* ─── Login Form ─── */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

/* ─── Form Group ─── */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: #374151;
  letter-spacing: 0.01em;
}

/* ─── Input Wrapper ─── */
.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  border: 1.5px solid #e2e8f0;
  border-radius: 10px;
  background: #f8fafc;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background 0.2s ease;
}

.input-wrapper--focused {
  border-color: #3b82f6;
  background: #ffffff;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.12);
}

.input-icon {
  display: flex;
  align-items: center;
  padding-left: 14px;
  color: #94a3b8;
  pointer-events: none;
  flex-shrink: 0;
  transition: color 0.2s ease;
}

.input-wrapper--focused .input-icon {
  color: #3b82f6;
}

.form-input {
  flex: 1;
  border: none;
  background: transparent;
  padding: 12px 14px;
  font-size: 0.875rem;
  font-family: inherit;
  color: #0f172a;
  outline: none;
  width: 100%;
}

.form-input::placeholder {
  color: #cbd5e1;
}

.form-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* ─── Toggle Password Button ─── */
.toggle-password {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0 14px;
  color: #94a3b8;
  display: flex;
  align-items: center;
  flex-shrink: 0;
  transition: color 0.2s ease;
}

.toggle-password:hover {
  color: #475569;
}

/* ─── Submit Button ─── */
.submit-btn {
  width: 100%;
  padding: 13px 20px;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: #ffffff;
  font-family: inherit;
  font-size: 0.875rem;
  font-weight: 600;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.35);
  margin-top: 4px;
}

.submit-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
  box-shadow: 0 4px 14px rgba(59, 130, 246, 0.45);
  transform: translateY(-1px);
}

.submit-btn:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 2px 6px rgba(59, 130, 246, 0.3);
}

.submit-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
  transform: none;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.2);
}

.submit-btn--loading {
  pointer-events: none;
}

/* ─── Button Content ─── */
.btn-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn-loading-content {
  gap: 10px;
}

/* ─── Loading Spinner ─── */
.spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.35);
  border-top-color: #ffffff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ─── Footer ─── */
.login-footer {
  font-size: 0.72rem;
  color: #94a3b8;
  text-align: center;
  margin: 0;
  letter-spacing: 0.01em;
}

/* ─── Responsive ─── */
@media (max-width: 480px) {
  .card {
    padding: 28px 22px 24px;
    border-radius: 16px;
  }

  .brand-name {
    font-size: 1.3rem;
  }
}
</style>
