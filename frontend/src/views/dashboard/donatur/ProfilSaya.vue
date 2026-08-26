<script setup>
import { ref, onMounted } from 'vue'
import { apiGetProfile } from '@/services/api'
import { useAuth } from '@/composables/useAuth'

const { user } = useAuth()
const profile = ref({
  name: user.value?.name || 'Loading...',
  email: user.value?.email || '',
  role: user.value?.role || 'donatur',
  phone: '',
  address: '',
})

const loading = ref(true)
const errorMsg = ref('')

onMounted(async () => {
  try {
    const data = await apiGetProfile()
    if (data) {
      profile.value = data
    }
  } catch (err) {
    errorMsg.value = err.message
    console.warn('[Profile] Failed to fetch:', err.message)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="profil-page">
    <h2>Profil Saya</h2>
    <p class="page-desc">Kelola informasi akun Anda.</p>
    
    <div v-if="errorMsg" class="error-msg">⚠️ {{ errorMsg }}</div>

    <div class="profil-card">
      <div class="profil-header">
        <div class="profil-avatar">{{ profile.name.charAt(0).toUpperCase() }}</div>
        <div class="profil-info">
          <h3>{{ profile.name }}</h3>
          <p>{{ profile.email }}</p>
          <span class="badge badge-success">{{ profile.role === 'donatur' ? 'Donatur Aktif' : 'Penggalang Dana' }}</span>
        </div>
      </div>

      <div class="profil-stats">
        <div class="ps-item">
          <span class="ps-value">0</span>
          <span class="ps-label">Total Donasi</span>
        </div>
        <div class="ps-item">
          <span class="ps-value">Rp 0</span>
          <span class="ps-label">Total Kontribusi</span>
        </div>
        <div class="ps-item">
          <span class="ps-value">0</span>
          <span class="ps-label">Campaign Didukung</span>
        </div>
      </div>

      <form class="profil-form">
        <div class="form-row">
          <div class="form-group">
            <label>Nama Lengkap</label>
            <input type="text" v-model="profile.name" />
          </div>
          <div class="form-group">
            <label>Email</label>
            <input type="email" :value="profile.email" disabled />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>Nomor Telepon</label>
            <input type="tel" v-model="profile.phone" placeholder="08xx-xxxx-xxxx" />
          </div>
          <div class="form-group">
            <label>Alamat</label>
            <input type="text" v-model="profile.address" placeholder="Alamat lengkap Anda" />
          </div>
        </div>
        <button type="button" class="btn-save">Simpan Perubahan</button>
      </form>
    </div>
  </div>
</template>

<style scoped>
h2 { font-family: var(--font-heading); font-size: 24px; margin-bottom: 4px; }
.page-desc { color: var(--text-muted); font-size: 14px; margin-bottom: 24px; }

.profil-card {
  background: white;
  border: 1px solid var(--border-light);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.profil-header {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 28px;
  background: linear-gradient(135deg, #E0F7FF, #F0F9FF);
  border-bottom: 1px solid var(--border-light);
}

.profil-avatar {
  width: 64px;
  height: 64px;
  border-radius: var(--radius-full);
  background: var(--primary);
  color: white;
  font-family: var(--font-heading);
  font-weight: 700;
  font-size: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.profil-info h3 { font-family: var(--font-heading); font-size: 20px; margin-bottom: 2px; }
.profil-info p { font-size: 14px; color: var(--text-muted); margin-bottom: 8px; }
.badge { padding: 4px 10px; border-radius: var(--radius-full); font-size: 12px; font-weight: 600; }
.badge-success { background: var(--success-light); color: var(--success); }

.profil-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1px;
  background: var(--border-light);
  border-bottom: 1px solid var(--border-light);
}

.ps-item {
  background: white;
  padding: 20px;
  text-align: center;
}

.ps-value {
  display: block;
  font-family: var(--font-heading);
  font-weight: 800;
  font-size: 20px;
  color: var(--primary);
  margin-bottom: 4px;
}

.ps-label {
  font-size: 12px;
  color: var(--text-muted);
}

.profil-form { padding: 28px; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.form-group { margin-bottom: 18px; }
.form-group label { display: block; font-size: 13px; font-weight: 600; color: var(--text-heading); margin-bottom: 6px; }
.form-group input {
  width: 100%; padding: 12px 16px; border: 1px solid var(--border); border-radius: var(--radius-sm);
  font-size: 14px; color: var(--text-heading); background: var(--bg-body); transition: all var(--transition-fast);
}
.form-group input:focus { outline: none; border-color: var(--primary); box-shadow: 0 0 0 3px var(--primary-light); background: white; }
.form-group input:disabled { background: var(--border-light); color: var(--text-muted); cursor: not-allowed; }

.btn-save {
  padding: 12px 28px;
  background: var(--primary);
  color: white;
  font-weight: 600;
  font-size: 14px;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}
.btn-save:hover { background: var(--primary-hover); box-shadow: 0 4px 12px var(--primary-glow); }

@media (max-width: 640px) { .form-row { grid-template-columns: 1fr; } .profil-stats { grid-template-columns: 1fr; } }
</style>
