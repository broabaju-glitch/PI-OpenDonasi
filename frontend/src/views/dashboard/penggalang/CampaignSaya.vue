<script setup>
import { ref, onMounted } from 'vue'
import { useAuth } from '@/composables/useAuth'

const { getToken } = useAuth()
const API_URL = import.meta.env.VITE_API_URL

const campaigns = ref([])
const isLoading = ref(true)
const uploadingId = ref(null)    // ID campaign yang sedang diupload fotonya
const uploadError = ref('')

// ─── Fetch campaign milik penggalang dari API ───
const fetchCampaigns = async () => {
  isLoading.value = true
  try {
    const res = await fetch(`${API_URL}/fundraiser/campaigns`, {
      headers: { Authorization: `Bearer ${getToken()}` }
    })
    if (!res.ok) throw new Error('Gagal mengambil data campaign')
    campaigns.value = await res.json()
  } catch (err) {
    console.error(err)
    campaigns.value = []
  } finally {
    isLoading.value = false
  }
}

onMounted(fetchCampaigns)

// ─── Helper: Hitung sisa hari dari end_date ───
const hitungSisaHari = (endDate) => {
  if (!endDate) return 0
  const end = new Date(endDate)
  const now = new Date()
  const diff = Math.ceil((end - now) / (1000 * 60 * 60 * 24))
  return diff > 0 ? diff : 0
}

// ─── Helper format Rupiah ───
const formatRupiah = (angka) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(angka || 0)
}

// ─── Logika status otomatis (tidak bergantung pada field status DB) ───
const getCampaignStatus = (c) => {
  const terkumpul = c.dana_terkumpul || 0
  const target = c.target_dana || 1
  const sisaHari = hitungSisaHari(c.end_date)
  if (terkumpul >= target || sisaHari <= 0) return 'Selesai'
  return 'Aktif'
}

const getStatusClass = (status) => {
  if (status === 'Aktif') return 'badge-success'
  if (status === 'Selesai') return 'badge-info'
  return 'badge-warning'
}

// ─── Resolusi URL foto (lokal atau absolut) ───
const getFotoSrc = (foto) => {
  if (!foto) return null
  if (foto.startsWith('http')) return foto
  // File lokal yang di-serve backend di /uploads/...
  return `http://localhost:8080${foto}`
}

// ─── Klik overlay edit foto → trigger file input ───
const triggerFileInput = (campaignId) => {
  const input = document.getElementById(`file-input-${campaignId}`)
  if (input) input.click()
}

// ─── Handle upload foto ke backend ───
const handlePhotoUpload = async (event, campaignId) => {
  const file = event.target.files?.[0]
  if (!file) return

  uploadingId.value = campaignId
  uploadError.value = ''

  const formData = new FormData()
  formData.append('photo', file)

  try {
    const res = await fetch(`${API_URL}/campaigns/${campaignId}/photo`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${getToken()}` },
      body: formData
    })

    if (!res.ok) {
      const errData = await res.json()
      throw new Error(errData.error || 'Gagal mengupload foto')
    }

    // Refresh daftar campaign agar foto terbaru ter-render secara reaktif
    await fetchCampaigns()
  } catch (err) {
    uploadError.value = err.message
    console.error('Upload gagal:', err)
  } finally {
    uploadingId.value = null
    // Reset file input agar bisa pilih file yang sama lagi
    event.target.value = ''
  }
}
</script>

<template>
  <div>
    <div class="page-top">
      <div>
        <h2>Campaign Saya</h2>
        <p class="page-desc">Kelola semua campaign yang Anda buat.</p>
      </div>
      <router-link to="/penggalang/buat-campaign" class="btn-create">➕ Buat Campaign</router-link>
    </div>

    <!-- Error upload -->
    <div v-if="uploadError" class="alert-error" role="alert">
      ⚠️ {{ uploadError }}
      <button @click="uploadError = ''" class="alert-close" aria-label="Tutup">✕</button>
    </div>

    <!-- Loading state -->
    <div v-if="isLoading" class="loading-state">
      <span class="spinner-ring"></span>
      <p>Memuat campaign Anda...</p>
    </div>

    <!-- Empty state -->
    <div v-else-if="!campaigns.length" class="empty-state">
      <span class="empty-icon">📂</span>
      <h3>Belum ada campaign</h3>
      <p>Buat campaign pertama Anda untuk mulai menggalang dana.</p>
      <router-link to="/penggalang/buat-campaign" class="btn-create" style="display: inline-flex;">➕ Buat Campaign Sekarang</router-link>
    </div>

    <!-- Daftar kartu campaign -->
    <div v-else class="cards-grid">
      <div v-for="c in campaigns" :key="c.ID" class="mc-card">

        <!-- ─── THUMBNAIL KIRI dengan Overlay Edit Foto ─── -->
        <div class="thumb-wrapper">
          <img
            :src="getFotoSrc(c.foto) || `https://ui-avatars.com/api/?name=${encodeURIComponent(c.title)}&size=400&background=e2e8f0&color=475569&bold=true`"
            :alt="`Foto campaign: ${c.title}`"
            class="mc-thumb"
          />

          <!-- Overlay tombol edit (muncul saat hover) -->
          <div
            class="thumb-overlay"
            :class="{ 'thumb-uploading': uploadingId === c.ID }"
            @click="triggerFileInput(c.ID)"
            :title="uploadingId === c.ID ? 'Sedang mengupload...' : 'Klik untuk mengganti foto'"
            role="button"
            :aria-label="`Ganti foto campaign ${c.title}`"
          >
            <span v-if="uploadingId === c.ID" class="spinner-ring spinner-white"></span>
            <template v-else>
              <span class="edit-icon">✏️</span>
              <span class="edit-label">Ganti Foto</span>
            </template>
          </div>

          <!-- Input file tersembunyi, satu per kartu -->
          <input
            :id="`file-input-${c.ID}`"
            type="file"
            accept="image/jpeg,image/png,image/webp"
            class="file-input-hidden"
            @change="handlePhotoUpload($event, c.ID)"
          />
        </div>

        <!-- ─── KONTEN KANAN ─── -->
        <div class="mc-content">
          <div class="mc-header">
            <h3>{{ c.title }}</h3>
            <span class="badge" :class="getStatusClass(getCampaignStatus(c))">
              {{ getCampaignStatus(c) }}
            </span>
          </div>

          <div class="mc-body">
            <div class="mc-stat">
              <span class="mcs-label">Terkumpul</span>
              <span class="mcs-val primary">{{ formatRupiah(c.dana_terkumpul) }}</span>
            </div>
            <div class="mc-stat">
              <span class="mcs-label">Target</span>
              <span class="mcs-val">{{ formatRupiah(c.target_dana) }}</span>
            </div>
            <div class="mc-stat">
              <span class="mcs-label">Sisa Hari</span>
              <span class="mcs-val">{{ hitungSisaHari(c.end_date) }} Hari</span>
            </div>
          </div>

          <!-- Progress bar -->
          <div class="progress-wrap">
            <div
              class="progress-bar"
              :style="{ width: Math.min(100, ((c.dana_terkumpul || 0) / (c.target_dana || 1)) * 100) + '%' }"
            ></div>
          </div>
          <p class="progress-label">
            {{ Math.min(100, Math.round(((c.dana_terkumpul || 0) / (c.target_dana || 1)) * 100)) }}% tercapai
          </p>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
.page-top { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 24px; }
h2 { font-family: var(--font-heading); font-size: 24px; margin-bottom: 4px; }
.page-desc { color: var(--text-muted); font-size: 14px; }
.btn-create { padding: 10px 20px; background: var(--primary); color: white; font-weight: 600; font-size: 14px; border-radius: var(--radius-sm); text-decoration: none; transition: all var(--transition-fast); white-space: nowrap; }
.btn-create:hover { background: var(--primary-hover); color: white; }

/* ─── ALERT ERROR ─── */
.alert-error {
  display: flex; align-items: center; justify-content: space-between;
  background: var(--danger-light); color: var(--danger); border: 1px solid var(--danger);
  border-radius: var(--radius-sm); padding: 10px 16px; margin-bottom: 16px; font-size: 14px;
}
.alert-close { background: none; border: none; cursor: pointer; color: var(--danger); font-size: 16px; padding: 0 4px; }

/* ─── LOADING & EMPTY ─── */
.loading-state, .empty-state {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  padding: 60px 24px; background: white; border: 1px solid var(--border-light);
  border-radius: var(--radius-md); text-align: center; gap: 12px;
}
.empty-icon { font-size: 48px; }
.empty-state h3 { font-family: var(--font-heading); font-size: 18px; color: var(--text-heading); }
.empty-state p { font-size: 14px; color: var(--text-muted); }

/* ─── SPINNER ─── */
.spinner-ring {
  display: inline-block; width: 28px; height: 28px;
  border: 3px solid rgba(0,0,0,0.1); border-top-color: var(--primary);
  border-radius: 50%; animation: spin 0.7s linear infinite;
}
.spinner-white { border-color: rgba(255,255,255,0.3); border-top-color: white; width: 22px; height: 22px; }
@keyframes spin { to { transform: rotate(360deg); } }

/* ─── KARTU ─── */
.cards-grid { display: flex; flex-direction: column; gap: 16px; }
.mc-card {
  background: white; border: 1px solid var(--border-light); border-radius: var(--radius-md);
  padding: 16px; transition: all var(--transition-fast);
  display: flex; flex-direction: row; align-items: stretch; gap: 20px;
}
.mc-card:hover { box-shadow: var(--shadow-md); }

/* ─── THUMBNAIL DENGAN OVERLAY ─── */
.thumb-wrapper {
  position: relative; flex-shrink: 0;
  width: 220px; height: 150px; border-radius: var(--radius-sm); overflow: hidden;
}
.mc-thumb {
  width: 100%; height: 100%;
  object-fit: cover; display: block;
  background-color: var(--border-light);
  transition: filter 0.25s ease;
}

/* Overlay edit foto */
.thumb-overlay {
  position: absolute; inset: 0;
  background: rgba(15, 23, 42, 0);
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  gap: 4px; cursor: pointer; opacity: 0;
  transition: all 0.25s ease;
  border-radius: var(--radius-sm);
}
.thumb-wrapper:hover .mc-thumb { filter: brightness(0.6); }
.thumb-wrapper:hover .thumb-overlay { opacity: 1; background: rgba(15, 23, 42, 0.4); }
.thumb-uploading {
  opacity: 1 !important;
  background: rgba(15, 23, 42, 0.55) !important;
  cursor: wait !important;
}
.thumb-uploading ~ .mc-thumb { filter: brightness(0.55); }

.edit-icon { font-size: 22px; }
.edit-label {
  font-size: 12px; font-weight: 700; color: white;
  background: rgba(0,0,0,0.35); padding: 2px 8px; border-radius: var(--radius-full);
}

/* Hidden file input */
.file-input-hidden {
  position: absolute; inset: 0; opacity: 0; cursor: pointer; width: 100%; height: 100%;
  /* pointer-events harus none agar klik dari overlay yang menangani */
  pointer-events: none;
}

/* ─── KONTEN ─── */
.mc-content { flex-grow: 1; display: flex; flex-direction: column; justify-content: space-between; }
.mc-header { display: flex; justify-content: space-between; align-items: flex-start; gap: 12px; margin-bottom: 12px; }
.mc-header h3 { font-family: var(--font-heading); font-size: 18px; margin: 0; color: var(--text-heading); }
.badge { padding: 4px 10px; border-radius: var(--radius-full); font-size: 12px; font-weight: 600; white-space: nowrap; }
.badge-success { background: var(--success-light); color: var(--success); }
.badge-info { background: var(--info-light); color: var(--info); }
.badge-warning { background: var(--warning-light); color: var(--warning); }

.mc-body { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; margin-bottom: 12px; }
.mc-stat { background: var(--bg-body); padding: 10px 14px; border-radius: var(--radius-sm); border: 1px solid var(--border-light); }
.mcs-label { display: block; font-size: 12px; color: var(--text-muted); margin-bottom: 2px; }
.mcs-val { display: block; font-weight: 700; font-size: 15px; color: var(--text-heading); }
.mcs-val.primary { color: var(--primary); }

/* ─── PROGRESS BAR ─── */
.progress-wrap {
  width: 100%; height: 6px; background: var(--border-light);
  border-radius: var(--radius-full); overflow: hidden; margin-bottom: 4px;
}
.progress-bar {
  height: 100%; background: var(--primary);
  border-radius: var(--radius-full); transition: width 0.4s ease;
  min-width: 2px;
}
.progress-label { font-size: 12px; color: var(--text-muted); text-align: right; }

/* ─── RESPONSIF ─── */
@media (max-width: 768px) {
  .mc-card { flex-direction: column; padding: 16px; gap: 16px; }
  .thumb-wrapper { width: 100%; height: 180px; }
  .mc-body { grid-template-columns: 1fr 1fr; }
}
</style>
