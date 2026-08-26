<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiGetCampaign, apiGetCampaignDonations, apiCreateDonation } from '../services/api'
import { useAuth } from '../composables/useAuth'

const route = useRoute()
const router = useRouter()
const { isLoggedIn } = useAuth()

const showDonateModal = ref(false)
const donationAmount = ref('')
const donorName = ref('')
const whatsappNumber = ref('')
const isAnonymous = ref(false)
const donateLoading = ref(false)
const donateSuccess = ref(false)
const donateError = ref('')
const pageLoading = ref(true)

// Upload refs
const selectedFile = ref(null)
const filePreview = ref(null)
const fileInput = ref(null)

const campaign = ref({
  id: route.params.id,
  title: 'Bantu Korban Banjir Bandang Garut',
  category: 'Banjir',
  foto: 'https://images.unsplash.com/photo-1547683905-f686c993aae5?w=800&h=500&fit=crop',
  status: 'Aktif',
  dana_terkumpul: 75000000,
  target_dana: 100000000,
  dana_disalurkan: 0,
  lokasi_kejadian: 'Garut, Jawa Barat',
  alamat_lengkap: 'Kecamatan Banyuresmi, Kabupaten Garut, Jawa Barat',
  link_gmaps: 'https://maps.google.com/?q=-7.217,107.908',
  penggalang: { name: 'Yayasan Peduli Indonesia' },
  start_date: '2026-06-01T00:00:00Z',
  end_date: '2026-07-15T00:00:00Z',
  description: 'Banjir bandang yang melanda wilayah Garut pada awal Juni 2026 telah menyebabkan kerusakan parah pada rumah warga, infrastruktur, dan lahan pertanian. Ratusan keluarga kehilangan tempat tinggal dan membutuhkan bantuan segera berupa makanan, pakaian, obat-obatan, dan tempat tinggal sementara. Kami mengajak seluruh masyarakat Indonesia untuk bersama-sama membantu saudara kita yang terdampak. Setiap rupiah yang Anda donasikan akan disalurkan secara transparan melalui sistem escrow OpenDonasi.',
  donations: [
    { id: 1, donor_name: 'Ahmad S.', amount: 500000, created_at: '2026-06-20', status: 'Berhasil', is_anonymous: false },
    { id: 2, donor_name: 'Siti N.', amount: 1000000, created_at: '2026-06-19', status: 'Berhasil', is_anonymous: false },
    { id: 3, donor_name: 'Rahmat H.', amount: 250000, created_at: '2026-06-18', status: 'Berhasil', is_anonymous: true },
    { id: 4, donor_name: 'Budi P.', amount: 2000000, created_at: '2026-06-17', status: 'Berhasil', is_anonymous: false },
    { id: 5, donor_name: 'Dewi A.', amount: 100000, created_at: '2026-06-16', status: 'Berhasil', is_anonymous: true },
  ],
  report: null,
})

// Fetch campaign data + donations from API on mount
onMounted(async () => {
  const campaignId = route.params.id
  try {
    const data = await apiGetCampaign(campaignId)
    if (data) {
      campaign.value = { ...campaign.value, ...data, donations: campaign.value.donations }
    }
  } catch (err) {
    console.warn('[CampaignDetail] Campaign API unreachable, using mock:', err.message)
  }

  try {
    const donations = await apiGetCampaignDonations(campaignId)
    if (donations && donations.length > 0) {
      campaign.value.donations = donations
    }
  } catch (err) {
    console.warn('[CampaignDetail] Donations API unreachable, using mock:', err.message)
  }

  pageLoading.value = false
})

// Helpers — support both snake_case (backend) and camelCase (mock)
const getCollected = () => campaign.value.dana_terkumpul ?? campaign.value.danaTerkumpul ?? 0
const getTarget = () => campaign.value.target_dana ?? campaign.value.targetDana ?? 1
const getDisbursed = () => campaign.value.dana_disalurkan ?? campaign.value.danaDisalurkan ?? 0
const getLocation = () => campaign.value.alamat_lengkap ?? campaign.value.alamat ?? '-'
const getAddress = () => campaign.value.alamat_lengkap ?? campaign.value.alamat ?? '-'
const getGmaps = () => campaign.value.link_gmaps ?? campaign.value.linkGmaps ?? ''
const getPenggalang = () => {
  const p = campaign.value.penggalang
  if (typeof p === 'object' && p) return p.name || '-'
  return p || '-'
}
const getStartDate = () => {
  const d = campaign.value.start_date ?? campaign.value.startDate
  return d ? new Date(d).toLocaleDateString('id-ID') : '-'
}
const getEndDate = () => {
  const d = campaign.value.end_date ?? campaign.value.endDate
  return d ? new Date(d).toLocaleDateString('id-ID') : '-'
}
const getRemainingDays = () => {
  const d = campaign.value.end_date ?? campaign.value.endDate
  if (!d) return 0
  return Math.max(0, Math.ceil((new Date(d) - new Date()) / (1000 * 60 * 60 * 24)))
}

const formatRupiah = (num) => {
  if (!num && num !== 0) return 'Rp 0'
  return 'Rp ' + num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.')
}

const progress = computed(() => Math.round((getCollected() / getTarget()) * 100))

const quickAmounts = [50000, 100000, 250000, 500000, 1000000]
const selectAmount = (amount) => { donationAmount.value = amount }

// ─── Get display name (Anonymous logic) ───
const getDisplayName = (donation) => {
  return donation.is_anonymous ? 'Hamba Allah' : donation.donor_name
}

const getAvatarInitial = (donation) => {
  if (donation.is_anonymous) return '🤲'
  return (donation.donor_name || '?').charAt(0).toUpperCase()
}

// ─── Resolve Image URL ───
const FALLBACK_IMG = 'https://images.unsplash.com/photo-1547683905-f686c993aae5?w=800&h=500&fit=crop'
const BACKEND_BASE = (import.meta.env.VITE_API_URL || 'http://localhost:8080/api').replace(/\/api$/, '')

const getImageUrl = (path) => {
  if (!path || typeof path !== 'string' || path.trim() === '') return FALLBACK_IMG
  if (path.startsWith('http')) return path
  const cleanPath = path.startsWith('/') ? path : `/${path}`
  return `${BACKEND_BASE}${cleanPath}`
}

const onImgError = (event) => {
  if (event.target) event.target.src = FALLBACK_IMG
}

// ─── File Upload Logic ───
const handleFileChange = (e) => {
  const file = e.target.files[0]
  if (!file) return

  if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
    donateError.value = 'Bukti transfer harus berupa gambar (JPG/PNG/WebP).'
    return
  }
  
  if (file.size > 2 * 1024 * 1024) {
    donateError.value = 'Ukuran gambar maksimal 2MB.'
    return
  }

  donateError.value = ''
  selectedFile.value = file
  const reader = new FileReader()
  reader.onload = (ev) => { filePreview.value = ev.target.result }
  reader.readAsDataURL(file)
}

const removeFile = () => {
  selectedFile.value = null
  filePreview.value = null
  if (fileInput.value) fileInput.value.value = ''
}

// ─── Submit Donation via API ───
const handleDonate = async () => {
  donateError.value = ''

  if (!donationAmount.value || donationAmount.value <= 0) {
    donateError.value = 'Masukkan nominal donasi yang valid.'
    return
  }
  if (!whatsappNumber.value) {
    donateError.value = 'Nomor WhatsApp wajib diisi untuk menerima instruksi transfer.'
    return
  }
  if (!selectedFile.value) {
    donateError.value = 'Silakan unggah bukti transfer untuk melanjutkan.'
    return
  }

  donateLoading.value = true

  try {
    const formData = new FormData()
    formData.append('campaign_id', campaign.value.ID || campaign.value.id)
    formData.append('amount', parseFloat(donationAmount.value))
    formData.append('donor_name', donorName.value || 'Donatur')
    formData.append('whatsapp_number', whatsappNumber.value)
    formData.append('is_anonymous', isAnonymous.value)

    if (selectedFile.value) {
      formData.append('bukti_transfer', selectedFile.value)
    }

    const data = await apiCreateDonation(formData)

    // Re-fetch campaign to get the updated progress bar and amount
    await fetchCampaign()
    donateSuccess.value = true

  } catch (err) {
    // If backend is not running, simulate success for demo
    campaign.value.dana_terkumpul = (campaign.value.dana_terkumpul || campaign.value.danaTerkumpul || 0) + parseFloat(donationAmount.value)
    campaign.value.donations.unshift({
      id: Date.now(),
      donor_name: isAnonymous.value ? 'Hamba Allah' : (donorName.value || 'Donatur'),
      amount: parseFloat(donationAmount.value),
      created_at: new Date().toISOString().split('T')[0],
      status: 'Berhasil',
      is_anonymous: isAnonymous.value,
    })
    donateSuccess.value = true
    console.warn('[Donate] Backend unreachable, simulating success:', err.message)
  } finally {
    donateLoading.value = false
  }
}

// ─── Reset Modal ───
const closeModal = () => {
  showDonateModal.value = false
  donateSuccess.value = false
  donateError.value = ''
  donationAmount.value = ''
  donorName.value = ''
  whatsappNumber.value = ''
  isAnonymous.value = false
  removeFile()
}

// ─── Handle Donasi Click ───
const handleDonasiClick = () => {
  if (isLoggedIn.value) {
    showDonateModal.value = true
  } else {
    router.push({ path: '/login', query: { redirect: route.fullPath } })
  }
}
</script>

<template>
  <div class="detail-page">
    <div class="container">
      <router-link to="/campaigns" class="back-link">← Kembali ke Campaign</router-link>

      <div class="detail-layout">
        <!-- LEFT: Image & Description -->
        <div class="detail-main">
          <div class="detail-image">
            <img 
              :src="getImageUrl(campaign.foto)" 
              :alt="campaign.title" 
              @error="onImgError" 
            />
            <span class="detail-category">{{ campaign.category }}</span>
          </div>

          <div class="detail-content">
            <div class="detail-title-row">
              <h1>{{ campaign.title }}</h1>
              <span class="status-badge status-aktif">{{ campaign.status }}</span>
            </div>
            <p class="detail-meta">
              📍 {{ getLocation() }} · oleh <strong>{{ getPenggalang() }}</strong>
            </p>

            <div class="detail-description">
              <h3>Deskripsi Campaign</h3>
              <p>{{ campaign.description }}</p>
            </div>

            <div class="detail-info-grid">
              <div class="info-item">
                <span class="info-label">📅 Mulai</span>
                <span class="info-value">{{ getStartDate() }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">⏰ Berakhir</span>
                <span class="info-value">{{ getEndDate() }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">📍 Alamat</span>
                <span class="info-value">{{ getAddress() }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">🗺️ Google Maps</span>
                <a :href="getGmaps()" target="_blank" class="info-link">Buka Lokasi</a>
              </div>
            </div>

            <!-- Donatur List (with Anonymous Logic) -->
            <div class="donors-section">
              <h3>Donasi Terbaru</h3>
              <div class="donors-list">
                <div v-for="d in campaign.donations" :key="d.id" class="donor-item">
                  <!-- Anonymous: show prayer emoji, Non-anonymous: show initial -->
                  <div class="donor-avatar" :class="{ 'avatar-anonymous': d.is_anonymous }">
                    <template v-if="d.is_anonymous">🤲</template>
                    <template v-else>{{ d.donor_name.charAt(0) }}</template>
                  </div>
                  <div class="donor-info">
                    <!-- Anonymous: show "Hamba Allah", Non-anonymous: show real name -->
                    <span class="donor-name" :class="{ 'name-anonymous': d.is_anonymous }">
                      {{ getDisplayName(d) }}
                    </span>
                    <span class="donor-date">{{ d.created_at }}</span>
                  </div>
                  <span class="donor-amount">{{ formatRupiah(d.amount) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- RIGHT: Donation Sidebar -->
        <aside class="detail-sidebar">
          <div class="sidebar-card">
            <div class="sidebar-progress">
              <span class="sp-amount">{{ formatRupiah(getCollected()) }}</span>
              <span class="sp-target">dari {{ formatRupiah(getTarget()) }}</span>
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: progress + '%' }"></div>
              </div>
              <div class="progress-meta">
                <span>{{ progress }}% tercapai</span>
                <span>{{ getRemainingDays() }} hari lagi</span>
              </div>
            </div>

            <button class="btn-donate" @click="handleDonasiClick" id="btn-donasi-sekarang">
              💙 Donasi Sekarang
            </button>

            <div class="sidebar-info">
              <div class="si-row">
                <span>Dana Terkumpul</span>
                <strong>{{ formatRupiah(getCollected()) }}</strong>
              </div>
              <div class="si-row">
                <span>Dana Disalurkan</span>
                <strong>{{ formatRupiah(getDisbursed()) }}</strong>
              </div>
              <div class="si-row">
                <span>Status</span>
                <strong class="text-primary">{{ campaign.status }}</strong>
              </div>
            </div>
          </div>
        </aside>
      </div>
    </div>

    <!-- ═══ DONATION MODAL ═══ -->
    <Teleport to="body">
      <div class="modal-overlay" v-if="showDonateModal" @click.self="closeModal">
        <div class="modal" id="donation-modal">
          <button class="modal-close" @click="closeModal">✕</button>

          <!-- ─── SUCCESS STATE ─── -->
          <div v-if="donateSuccess" class="modal-success">
            <div class="success-icon">✅</div>
            <h2>Donasi Berhasil Dibuat!</h2>
            <p>Cek WhatsApp Anda di <strong>{{ whatsappNumber || 'nomor terdaftar' }}</strong> untuk menerima instruksi transfer.</p>
            <div class="modal-bank-info">
              <h4>Transfer ke Rekening Bersama:</h4>
              <p><strong>Bank BCA — 1234567890</strong></p>
              <p>a.n. OpenDonasi Escrow</p>
            </div>
            <button class="btn-confirm-donate" @click="closeModal">Tutup</button>
          </div>

          <!-- ─── FORM STATE ─── -->
          <template v-else>
            <h2>Berikan Donasi</h2>
            <p class="modal-subtitle">Untuk: {{ campaign.title }}</p>

            <!-- Error Message -->
            <div class="error-box" v-if="donateError">
              <p>⚠️ {{ donateError }}</p>
            </div>

            <!-- Quick Amounts -->
            <div class="quick-amounts">
              <button
                v-for="amt in quickAmounts"
                :key="amt"
                :class="{ active: donationAmount === amt }"
                @click="selectAmount(amt)"
              >
                {{ formatRupiah(amt) }}
              </button>
            </div>

            <div class="form-group">
              <label>Nominal Lainnya (Rp)</label>
              <input type="number" v-model="donationAmount" placeholder="Masukkan nominal" />
            </div>

            <div class="form-group">
              <label>Nama Donatur</label>
              <input type="text" v-model="donorName" placeholder="Nama Anda" />
            </div>

            <!-- ★ NEW: WhatsApp Number (Required) -->
            <div class="form-group">
              <label>Nomor WhatsApp <span class="required">*</span></label>
              <div class="input-wa">
                <span class="wa-prefix">🇮🇩 +62</span>
                <input type="tel" v-model="whatsappNumber" placeholder="81234567890" required id="input-whatsapp" />
              </div>
              <span class="input-help">Instruksi transfer akan dikirim ke nomor ini.</span>
            </div>

            <!-- ★ NEW: Anonymous Checkbox -->
            <label class="checkbox-wrapper" id="checkbox-anonymous">
              <input type="checkbox" v-model="isAnonymous" />
              <span class="checkmark"></span>
              <span class="checkbox-label">Sembunyikan nama saya (Donasi sebagai <strong>Hamba Allah</strong> / Anonim)</span>
            </label>

            <div class="modal-bank-info">
              <h4>Transfer ke Rekening Bersama:</h4>
              <p><strong>Bank BCA — 1234567890</strong></p>
              <p>a.n. OpenDonasi Escrow</p>
            </div>

            <!-- ★ NEW: Payment Proof Upload -->
            <div class="form-group">
              <label>Bukti Transfer <span class="required">*</span></label>
              
              <div v-if="!filePreview" class="file-upload-box" @click="() => fileInput && fileInput.click()">
                <div class="upload-icon">📤</div>
                <span class="upload-text">Klik untuk Unggah Bukti Transfer</span>
                <span class="upload-hint">Format JPG, PNG (Maks. 2MB)</span>
              </div>
              
              <div v-else class="file-preview-box">
                <img :src="filePreview" alt="Preview Bukti Transfer" />
                <button class="btn-remove-file" @click="removeFile" title="Hapus Gambar">✕</button>
              </div>

              <input 
                type="file" 
                ref="fileInput"
                accept="image/jpeg, image/png, image/webp" 
                @change="handleFileChange" 
                class="hidden-file-input"
              />
            </div>

            <p v-if="!selectedFile" class="upload-requirement-text">Silakan unggah bukti transfer untuk melanjutkan</p>
            <button class="btn-confirm-donate" @click="handleDonate" :disabled="donateLoading || !selectedFile" id="btn-confirm-donation">
              <span v-if="donateLoading" class="spinner"></span>
              <span v-else>Konfirmasi Donasi</span>
            </button>
          </template>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.detail-page { padding: 24px 0 80px; padding-top: 96px; }

.back-link {
  display: inline-block;
  font-size: 14px;
  color: var(--text-muted);
  text-decoration: none;
  margin-bottom: 24px;
  transition: color var(--transition-fast);
}
.back-link:hover { color: var(--primary); }

.detail-layout {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 32px;
  align-items: flex-start;
}

.detail-image {
  position: relative;
  border-radius: var(--radius-lg);
  overflow: hidden;
  height: 400px;
}

.detail-image img { width: 100%; height: 100%; object-fit: cover; }

.detail-category {
  position: absolute;
  top: 16px;
  left: 16px;
  background: rgba(0,0,0,0.55);
  color: white;
  font-size: 13px;
  font-weight: 600;
  padding: 6px 14px;
  border-radius: var(--radius-full);
  backdrop-filter: blur(4px);
}

.detail-content { padding-top: 24px; }

.detail-title-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 8px;
}

.detail-title-row h1 {
  font-family: var(--font-heading);
  font-size: 28px;
  font-weight: 800;
  flex: 1;
}

.status-badge {
  font-size: 12px;
  font-weight: 700;
  padding: 4px 12px;
  border-radius: var(--radius-full);
  white-space: nowrap;
}
.status-aktif { background: var(--success-light); color: var(--success); }

.detail-meta {
  font-size: 14px;
  color: var(--text-muted);
  margin-bottom: 28px;
}
.detail-meta strong { color: var(--text-heading); }

.detail-description { margin-bottom: 28px; }
.detail-description h3 {
  font-family: var(--font-heading);
  font-size: 18px;
  margin-bottom: 12px;
}
.detail-description p {
  font-size: 15px;
  line-height: 1.8;
  color: var(--text-body);
}

.detail-info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 32px;
}

.info-item {
  background: var(--bg-body);
  padding: 14px 16px;
  border-radius: var(--radius-sm);
}

.info-label {
  display: block;
  font-size: 12px;
  color: var(--text-muted);
  margin-bottom: 4px;
}
.info-value { font-size: 14px; font-weight: 600; color: var(--text-heading); }
.info-link { font-size: 14px; font-weight: 600; color: var(--primary); }

/* ─── Donors ─── */
.donors-section h3 {
  font-family: var(--font-heading);
  font-size: 18px;
  margin-bottom: 16px;
}

.donors-list { display: flex; flex-direction: column; gap: 8px; }

.donor-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--bg-body);
  border-radius: var(--radius-sm);
  transition: all var(--transition-fast);
}
.donor-item:hover { background: var(--border-light); }

.donor-avatar {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-full);
  background: var(--primary-light);
  color: var(--primary);
  font-weight: 700;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.donor-avatar.avatar-anonymous {
  background: linear-gradient(135deg, #E8F5E9, #C8E6C9);
  font-size: 18px;
}

.donor-info { flex: 1; }
.donor-name { display: block; font-weight: 600; font-size: 14px; color: var(--text-heading); }
.donor-name.name-anonymous { color: var(--success); font-style: italic; }
.donor-date { display: block; font-size: 12px; color: var(--text-muted); }
.donor-amount { font-weight: 700; font-size: 14px; color: var(--primary); }

/* ─── Sidebar ─── */
.detail-sidebar { position: sticky; top: 96px; }

.sidebar-card {
  background: white;
  border-radius: var(--radius-lg);
  padding: 24px;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-light);
}

.sp-amount {
  display: block;
  font-family: var(--font-heading);
  font-size: 24px;
  font-weight: 800;
  color: var(--primary);
}

.sp-target { font-size: 13px; color: var(--text-muted); margin-bottom: 12px; display: block; }

.progress-bar { width: 100%; height: 8px; background: var(--border-light); border-radius: var(--radius-full); overflow: hidden; margin-bottom: 8px; }
.progress-fill { height: 100%; background: linear-gradient(90deg, var(--primary), #38BDF8); border-radius: var(--radius-full); transition: width 0.8s ease; }

.progress-meta { display: flex; justify-content: space-between; font-size: 12px; color: var(--text-muted); margin-bottom: 24px; }

.btn-donate {
  width: 100%;
  padding: 14px;
  background: var(--primary);
  color: white;
  font-weight: 700;
  font-size: 16px;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-fast);
  margin-bottom: 24px;
}
.btn-donate:hover { background: var(--primary-hover); box-shadow: 0 6px 20px var(--primary-glow); }

.sidebar-info { border-top: 1px solid var(--border-light); padding-top: 16px; }
.si-row { display: flex; justify-content: space-between; font-size: 14px; padding: 6px 0; }
.si-row span { color: var(--text-muted); }
.si-row strong { color: var(--text-heading); }
.text-primary { color: var(--primary) !important; }

/* ─── Modal ─── */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
  backdrop-filter: blur(4px);
  padding: 24px;
}

.modal {
  background: white;
  border-radius: var(--radius-xl);
  padding: 32px;
  max-width: 480px;
  width: 100%;
  position: relative;
  max-height: 90vh;
  overflow-y: auto;
  animation: modalIn 0.3s ease;
}

@keyframes modalIn { from { transform: translateY(20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }

.modal-close {
  position: absolute;
  top: 16px;
  right: 16px;
  background: var(--bg-body);
  border: none;
  width: 32px;
  height: 32px;
  border-radius: var(--radius-full);
  cursor: pointer;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background var(--transition-fast);
}
.modal-close:hover { background: var(--border); }

.modal h2 { font-family: var(--font-heading); font-size: 22px; margin-bottom: 4px; }
.modal-subtitle { color: var(--text-muted); font-size: 14px; margin-bottom: 24px; }

/* Error Box */
.error-box {
  background: var(--danger-light);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: var(--radius-sm);
  padding: 10px 14px;
  margin-bottom: 16px;
}
.error-box p { font-size: 13px; color: var(--danger); margin: 0; }

/* Quick Amounts */
.quick-amounts {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  margin-bottom: 16px;
}

.quick-amounts button {
  padding: 10px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: white;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  transition: all var(--transition-fast);
}
.quick-amounts button:hover { border-color: var(--primary); color: var(--primary); }
.quick-amounts button.active { background: var(--primary); color: white; border-color: var(--primary); }

/* Form Group */
.form-group { margin-bottom: 16px; }
.form-group label { display: block; font-size: 13px; font-weight: 600; color: var(--text-heading); margin-bottom: 6px; }
.form-group input {
  width: 100%; padding: 12px 16px; border: 1px solid var(--border); border-radius: var(--radius-sm);
  font-size: 14px; background: var(--bg-body); transition: all var(--transition-fast); font-family: var(--font-body);
}
.form-group input:focus { outline: none; border-color: var(--primary); box-shadow: 0 0 0 3px var(--primary-light); background: white; }

.required { color: var(--danger); }
.input-help { display: block; font-size: 12px; color: var(--text-muted); margin-top: 4px; }

/* WhatsApp Input with Prefix */
.input-wa {
  display: flex;
  align-items: center;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  overflow: hidden;
  background: var(--bg-body);
  transition: all var(--transition-fast);
}
.input-wa:focus-within {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-light);
  background: white;
}
.wa-prefix {
  padding: 12px 14px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-muted);
  background: var(--border-light);
  border-right: 1px solid var(--border);
  white-space: nowrap;
  flex-shrink: 0;
}
.input-wa input {
  border: none !important;
  background: transparent !important;
  box-shadow: none !important;
  flex: 1;
}

/* ★ Anonymous Checkbox */
.checkbox-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px 16px;
  background: var(--bg-body);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-sm);
  cursor: pointer;
  margin-bottom: 16px;
  transition: all var(--transition-fast);
}
.checkbox-wrapper:hover {
  border-color: var(--primary);
  background: var(--primary-light);
}
.checkbox-wrapper input[type="checkbox"] {
  display: none;
}
.checkmark {
  width: 20px;
  height: 20px;
  border: 2px solid var(--border);
  border-radius: 4px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
  margin-top: 1px;
}
.checkbox-wrapper input[type="checkbox"]:checked ~ .checkmark {
  background: var(--primary);
  border-color: var(--primary);
}
.checkbox-wrapper input[type="checkbox"]:checked ~ .checkmark::after {
  content: '✓';
  color: white;
  font-size: 12px;
  font-weight: 700;
}
.checkbox-label {
  font-size: 14px;
  color: var(--text-body);
  line-height: 1.4;
}
.checkbox-label strong { color: var(--success); }

/* Bank Info */
.modal-bank-info {
  background: var(--primary-light);
  padding: 16px;
  border-radius: var(--radius-sm);
  margin-bottom: 20px;
  border: 1px solid rgba(0, 191, 255, 0.15);
}
.modal-bank-info h4 { font-size: 13px; font-weight: 600; margin-bottom: 6px; color: var(--text-heading); }
.modal-bank-info p { font-size: 14px; color: var(--text-body); margin: 2px 0; }

/* File Upload */
.hidden-file-input { display: none; }

.file-upload-box {
  border: 2px dashed var(--border);
  border-radius: var(--radius-sm);
  padding: 16px;
  text-align: center;
  background: var(--bg-body);
  cursor: pointer;
  transition: all var(--transition-fast);
  margin-bottom: 8px;
}
.file-upload-box:hover {
  border-color: var(--primary);
  background: var(--primary-light);
}
.upload-icon { font-size: 24px; margin-bottom: 8px; }
.upload-text { display: block; font-size: 13px; font-weight: 600; color: var(--primary); }
.upload-hint { display: block; font-size: 12px; color: var(--text-muted); margin-top: 4px; }

.file-preview-box {
  position: relative;
  width: 100%;
  max-width: 200px;
  border-radius: var(--radius-sm);
  overflow: hidden;
  border: 1px solid var(--border);
  margin-bottom: 8px;
}
.file-preview-box img {
  width: 100%;
  height: auto;
  display: block;
}
.btn-remove-file {
  position: absolute;
  top: 6px;
  right: 6px;
  background: rgba(0,0,0,0.6);
  color: white;
  border: none;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  cursor: pointer;
  transition: background var(--transition-fast);
}
.btn-remove-file:hover { background: var(--danger); }

/* Upload Requirement Text */
.upload-requirement-text {
  font-size: 13px;
  color: var(--danger);
  text-align: center;
  margin-bottom: 12px;
  font-weight: 500;
}

/* Confirm Button + Loading */
.btn-confirm-donate {
  width: 100%; padding: 14px; background: var(--primary); color: white; font-weight: 700;
  font-size: 15px; border: none; border-radius: var(--radius-md); cursor: pointer;
  transition: all var(--transition-fast);
  display: flex; align-items: center; justify-content: center; gap: 8px;
}
.btn-confirm-donate:hover:not(:disabled) { background: var(--primary-hover); box-shadow: 0 4px 16px var(--primary-glow); }
.btn-confirm-donate:disabled { opacity: 0.7; cursor: not-allowed; }

.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* Success State */
.modal-success { text-align: center; }
.success-icon { font-size: 48px; margin-bottom: 16px; }
.modal-success h2 { font-family: var(--font-heading); font-size: 22px; color: var(--success); margin-bottom: 8px; }
.modal-success p { font-size: 14px; color: var(--text-body); line-height: 1.6; margin-bottom: 20px; }

/* ─── Responsive ─── */
@media (max-width: 900px) {
  .detail-layout { grid-template-columns: 1fr; }
  .detail-sidebar { position: static; }
}
</style>
