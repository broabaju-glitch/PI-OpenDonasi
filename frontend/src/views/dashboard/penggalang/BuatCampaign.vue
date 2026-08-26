<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { getToken } = useAuth()
const API_URL = import.meta.env.VITE_API_URL

const form = ref({
  title: '',
  category: '',
  description: '',
  targetDana: '',
  alamat: '',
  linkGmaps: '',
  startDate: '',
  endDate: '',
  rekening: ''
})

// Reactive state untuk input kategori kustom
const customCategory = ref('')

const loading = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

// Reset customCategory saat dropdown berubah ke selain Lainnya
const onCategoryChange = () => {
  if (form.value.category !== 'Lainnya') {
    customCategory.value = ''
  }
}

// Resolve kategori final: jika Lainnya, gunakan customCategory
const getFinalCategory = () => {
  if (form.value.category === 'Lainnya') {
    return customCategory.value.trim()
  }
  return form.value.category
}

const handleSubmit = async () => {
  loading.value = true
  errorMsg.value = ''
  successMsg.value = ''

  try {
    // Validasi kategori kustom jika dipilih "Lainnya"
    const finalCategory = getFinalCategory()
    if (!finalCategory) {
      throw new Error('Kategori tidak boleh kosong. Silakan pilih atau isi kategori.')
    }

    const res = await fetch(`${API_URL}/fundraiser/campaigns`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${getToken()}`
      },
      body: JSON.stringify({
        title:           form.value.title,
        category:        finalCategory,
        description:     form.value.description,
        target_dana:     parseFloat(form.value.targetDana),
        alamat_lengkap:  form.value.alamat,
        link_gmaps:      form.value.linkGmaps,
        start_date:      form.value.startDate,
        end_date:        form.value.endDate,
        rekening:        form.value.rekening,
      })
    })

    const data = await res.json()

    if (!res.ok) {
      throw new Error(data.error || 'Terjadi kesalahan saat membuat campaign.')
    }

    successMsg.value = data.message || 'Campaign berhasil dibuat!'
    
    // Tunggu sebentar agar user bisa membaca pesan sukses,
    // lalu arahkan ke Campaign Saya yang akan otomatis fetch ulang dari API
    setTimeout(() => {
      router.push('/penggalang/campaign-saya')
    }, 1500)

  } catch (err) {
    errorMsg.value = err.message
    console.error('Buat campaign gagal:', err)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <h2>Buat Campaign Baru</h2>
    <p class="page-desc">Isi semua data campaign sesuai ketentuan. Campaign akan aktif setelah diverifikasi admin.</p>

    <!-- Notifikasi sukses -->
    <div v-if="successMsg" class="alert alert-success" role="alert">
      ✅ {{ successMsg }} <span class="alert-sub">Mengarahkan ke daftar campaign...</span>
    </div>

    <!-- Notifikasi error -->
    <div v-if="errorMsg" class="alert alert-error" role="alert">
      ⚠️ {{ errorMsg }}
      <button @click="errorMsg = ''" class="alert-close" aria-label="Tutup">✕</button>
    </div>

    <form @submit.prevent="handleSubmit" class="create-form" id="form-buat-campaign">
      <div class="form-row">
        <div class="form-group">
          <label>Judul Campaign *</label>
          <input type="text" v-model="form.title" placeholder="Judul campaign Anda" required :disabled="loading" />
        </div>
        <div class="form-group">
          <label>Kategori *</label>
          <select
            v-model="form.category"
            @change="onCategoryChange"
            required
            :disabled="loading"
          >
            <option value="">Pilih Kategori</option>
            <option>Banjir</option>
            <option>Gempa Bumi</option>
            <option>Gunung Meletus</option>
            <option>Tsunami</option>
            <option>Tanah Longsor</option>
            <option>Rehabilitasi</option>
            <option value="Lainnya">Lainnya (isi sendiri)</option>
          </select>
          <!-- Input kustom muncul hanya saat memilih Lainnya -->
          <input
            v-if="form.category === 'Lainnya'"
            v-model="customCategory"
            type="text"
            placeholder="Masukkan kategori baru..."
            class="custom-category-input"
            maxlength="50"
            :disabled="loading"
            :required="form.category === 'Lainnya'"
          />
        </div>
      </div>

      <div class="form-group">
        <label>Deskripsi *</label>
        <textarea v-model="form.description" rows="5" placeholder="Jelaskan detail campaign, situasi darurat, dan rencana penggunaan dana..." required :disabled="loading"></textarea>
      </div>

      <div class="form-group">
        <label>Target Dana (Rp) *</label>
        <input type="number" v-model="form.targetDana" placeholder="Contoh: 100000000" min="100000" required :disabled="loading" />
      </div>

      <div class="form-group">
        <label>Alamat Lengkap *</label>
        <input type="text" v-model="form.alamat" placeholder="Alamat lengkap lokasi kejadian" required :disabled="loading" />
      </div>

      <div class="form-group">
        <label>Link Google Maps</label>
        <input type="url" v-model="form.linkGmaps" placeholder="https://maps.google.com/?q=..." :disabled="loading" />
      </div>

      <div class="form-row">
        <div class="form-group">
          <label>Tanggal Mulai *</label>
          <input type="date" v-model="form.startDate" required :disabled="loading" />
        </div>
        <div class="form-group">
          <label>Tanggal Berakhir *</label>
          <input type="date" v-model="form.endDate" required :disabled="loading" />
        </div>
      </div>

      <div class="form-group">
        <label>Rekening Penggalang *</label>
        <input type="text" v-model="form.rekening" placeholder="BCA 1234567890 a.n. Nama Lengkap" required :disabled="loading" />
      </div>

      <p class="form-note">
        💡 Foto campaign dapat diupload setelah campaign berhasil dibuat, melalui halaman "Campaign Saya".
      </p>

      <button type="submit" class="btn-submit" :disabled="loading" id="btn-submit-campaign">
        <span v-if="loading" class="spinner"></span>
        <span v-else>📤 Buat Campaign</span>
      </button>
    </form>
  </div>
</template>

<style scoped>
h2 { font-family: var(--font-heading); font-size: 24px; margin-bottom: 4px; }
.page-desc { color: var(--text-muted); font-size: 14px; margin-bottom: 20px; }

.alert {
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 18px; border-radius: var(--radius-sm); margin-bottom: 16px; font-size: 14px; font-weight: 500;
}
.alert-success { background: var(--success-light); color: var(--success); border: 1px solid var(--success); }
.alert-error   { background: var(--danger-light);  color: var(--danger);  border: 1px solid var(--danger); }
.alert-sub { font-size: 12px; font-weight: 400; margin-left: 8px; opacity: 0.8; }
.alert-close { background: none; border: none; cursor: pointer; color: inherit; font-size: 16px; padding: 0 4px; }

.create-form { background: white; border: 1px solid var(--border-light); border-radius: var(--radius-md); padding: 28px; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.form-group { margin-bottom: 18px; }
.form-group label { display: block; font-size: 13px; font-weight: 600; color: var(--text-heading); margin-bottom: 6px; }
.form-group input, .form-group select, .form-group textarea {
  width: 100%; padding: 12px 16px; border: 1px solid var(--border); border-radius: var(--radius-sm);
  font-size: 14px; font-family: var(--font-body); color: var(--text-heading); background: var(--bg-body); transition: all var(--transition-fast); resize: vertical;
}
.form-group input:focus, .form-group select:focus, .form-group textarea:focus {
  outline: none; border-color: var(--primary); box-shadow: 0 0 0 3px var(--primary-light); background: white;
}
.form-group input:disabled, .form-group select:disabled, .form-group textarea:disabled {
  opacity: 0.65; cursor: not-allowed;
}

.form-note {
  font-size: 13px; color: var(--text-muted); background: var(--bg-body);
  border: 1px dashed var(--border); border-radius: var(--radius-sm);
  padding: 10px 14px; margin-bottom: 20px;
}

.btn-submit {
  width: 100%; padding: 14px; background: var(--primary); color: white;
  font-weight: 700; font-size: 15px; border: none; border-radius: var(--radius-sm);
  cursor: pointer; transition: all var(--transition-fast);
  display: flex; align-items: center; justify-content: center; gap: 8px;
}
.btn-submit:hover:not(:disabled) { background: var(--primary-hover); box-shadow: 0 4px 16px var(--primary-glow); }
.btn-submit:disabled { opacity: 0.7; cursor: not-allowed; }
.spinner {
  width: 18px; height: 18px; border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white; border-radius: 50%; animation: spin 0.6s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
/* Input kategori kustom (muncul di bawah select saat memilih Lainnya) */
.custom-category-input {
  margin-top: 8px;
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--primary);
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-family: var(--font-body);
  color: var(--text-heading);
  background: white;
  box-shadow: 0 0 0 3px var(--primary-light);
  transition: all var(--transition-fast);
  animation: slideDown 0.2s ease;
}
.custom-category-input:focus {
  outline: none;
  border-color: var(--primary-hover);
  box-shadow: 0 0 0 3px var(--primary-light);
}
@keyframes slideDown {
  from { opacity: 0; transform: translateY(-6px); }
  to   { opacity: 1; transform: translateY(0); }
}
@media (max-width: 640px) { .form-row { grid-template-columns: 1fr; } }
</style>
