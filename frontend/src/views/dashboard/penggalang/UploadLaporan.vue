<script setup>
import { ref, onMounted } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { apiUploadReport } from '@/services/api'

const { getToken } = useAuth()
const API_URL = import.meta.env.VITE_API_URL

const campaigns = ref([])
const form = ref({ campaign_id: '', title: '', description: '' })
const fileImage = ref(null)
const loading = ref(false)

const fetchMyCampaigns = async () => {
  const token = getToken() || localStorage.getItem('opendonasi_token')
  console.log('[UploadLaporan] Token:', token ? 'Ada ✅' : 'Kosong ❌')

  if (!token) {
    console.warn('[UploadLaporan] Tidak ada token, fetch dibatalkan.')
    return
  }

  try {
    const res = await fetch(`${API_URL}/fundraiser/campaigns`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    campaigns.value = Array.isArray(data) ? data : (data.data || [])
    console.log('[UploadLaporan] Campaigns loaded:', campaigns.value.length, campaigns.value)
  } catch (err) {
    console.error('[UploadLaporan] Fetch gagal:', err)
    campaigns.value = []
  }
}

onMounted(fetchMyCampaigns)

const handleFileChange = (e) => {
  if (e.target.files.length > 0) {
    fileImage.value = e.target.files[0]
  } else {
    fileImage.value = null
  }
}

const handleSubmit = async () => {
  if (!fileImage.value) return
  
  loading.value = true
  try {
    const formData = new FormData()
    formData.append('campaign_id', form.value.campaign_id)
    formData.append('title', form.value.title)
    formData.append('description', form.value.description)
    formData.append('image', fileImage.value)

    await apiUploadReport(formData)
    alert('Laporan berhasil diupload! Menunggu verifikasi admin.')
    
    // Reset form
    form.value = { campaign_id: '', title: '', description: '' }
    fileImage.value = null
    document.getElementById('form-upload-laporan').reset()
  } catch (err) {
    console.error('Upload error:', err)
    alert('Gagal mengupload laporan. Silakan coba lagi.')
  } finally {
    loading.value = false
  }
}
</script>
<template>
  <div>
    <h2>Upload Laporan Penyaluran</h2>
    <p class="page-desc">Unggah laporan penyaluran dana setelah menerima transfer dari admin.</p>
    <form @submit.prevent="handleSubmit" class="report-form" id="form-upload-laporan">
      <div class="form-group">
        <label>Campaign *</label>
        <select v-model="form.campaign_id" required>
          <option value="" disabled>Pilih Campaign</option>
          <option v-for="campaign in campaigns" :key="campaign.ID || campaign.id" :value="campaign.ID || campaign.id">
            {{ campaign.title }}
          </option>
        </select>
      </div>
      <div class="form-group"><label>Judul Laporan *</label><input type="text" v-model="form.title" placeholder="Laporan Penyaluran Tahap 1" required /></div>
      <div class="form-group"><label>Deskripsi Laporan *</label><textarea v-model="form.description" rows="6" placeholder="Jelaskan detail penyaluran dana..." required></textarea></div>
      <div class="form-group"><label>Foto Bukti Penyaluran *</label><input type="file" accept="image/*" @change="handleFileChange" required /></div>
      <button type="submit" class="btn-submit" :disabled="loading || !fileImage">
        <span v-if="loading" class="spinner"></span>
        <span v-else>📤 Upload Laporan</span>
      </button>
    </form>
  </div>
</template>
<style scoped>
h2 { font-family: var(--font-heading); font-size: 24px; margin-bottom: 4px; }
.page-desc { color: var(--text-muted); font-size: 14px; margin-bottom: 28px; }
.report-form { background: white; border: 1px solid var(--border-light); border-radius: var(--radius-md); padding: 28px; }
.form-group { margin-bottom: 18px; }
.form-group label { display: block; font-size: 13px; font-weight: 600; color: var(--text-heading); margin-bottom: 6px; }
.form-group input, .form-group select, .form-group textarea {
  width: 100%; padding: 12px 16px; border: 1px solid var(--border); border-radius: var(--radius-sm);
  font-size: 14px; font-family: var(--font-body); color: var(--text-heading); background: var(--bg-body); transition: all var(--transition-fast); resize: vertical;
}
.form-group input:focus, .form-group select:focus, .form-group textarea:focus { outline: none; border-color: var(--primary); box-shadow: 0 0 0 3px var(--primary-light); background: white; }
.btn-submit { width: 100%; padding: 14px; background: var(--primary); color: white; font-weight: 700; font-size: 15px; border: none; border-radius: var(--radius-sm); cursor: pointer; transition: all var(--transition-fast); display: flex; align-items: center; justify-content: center; gap: 8px; }
.btn-submit:hover:not(:disabled) { background: var(--primary-hover); box-shadow: 0 4px 16px var(--primary-glow); }
.btn-submit:disabled { opacity: 0.7; cursor: not-allowed; }
.spinner { width: 18px; height: 18px; border: 2px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 0.6s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
