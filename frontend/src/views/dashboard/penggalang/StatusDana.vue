<script setup>
import { ref, onMounted } from 'vue'
import { apiGetFundraiserSummary, apiGetFundraiserDonations, apiRequestWithdrawal } from '@/services/api'

const summary = ref({
  total_terkumpul: 0,
  saldo_aktif: 0,
  total_pending: 0,
  has_pending: false
})
const donations = ref([])
const loading = ref(true)

// Modal State
const showModal = ref(false)
const modalLoading = ref(false)
const modalError = ref('')
const modalSuccess = ref('')
const form = ref({
  bank_name: '',
  bank_account: '',
  account_name: ''
})

const fetchSummary = async () => {
  try {
    const data = await apiGetFundraiserSummary()
    summary.value = data
  } catch (err) {
    console.error('Failed to fetch summary:', err)
  }
}

const fetchDonations = async () => {
  try {
    const data = await apiGetFundraiserDonations()
    donations.value = data
  } catch (err) {
    console.error('Failed to fetch donations:', err)
  }
}

onMounted(async () => {
  await Promise.all([fetchSummary(), fetchDonations()])
  loading.value = false
})

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(amount)
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Intl.DateTimeFormat('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'short'
  }).format(new Date(dateString))
}

const getBadgeInfo = (status) => {
  if (status === 'Menunggu Verifikasi') return { text: '🔴 Menunggu Verifikasi Admin', class: 'badge-waiting' }
  if (status === 'Berhasil') return { text: '🟢 Terverifikasi / Masuk Saldo', class: 'badge-verified' }
  if (status === 'Ditolak') return { text: '⚫ Ditolak', class: 'badge-rejected' }
  return { text: status, class: 'badge-default' }
}

const openWithdrawalModal = () => {
  showModal.value = true
  modalError.value = ''
  modalSuccess.value = ''
  form.value = { bank_name: '', bank_account: '', account_name: '' }
}

const submitWithdrawal = async () => {
  if (!form.value.bank_name || !form.value.bank_account || !form.value.account_name) {
    modalError.value = 'Semua field wajib diisi'
    return
  }

  modalLoading.value = true
  modalError.value = ''
  modalSuccess.value = ''

  try {
    const res = await apiRequestWithdrawal(form.value)
    modalSuccess.value = res.message || 'Permintaan penarikan berhasil.'
    
    // Refresh summary
    await fetchSummary()
    
    setTimeout(() => {
      showModal.value = false
    }, 2000)
  } catch (err) {
    modalError.value = err.message || 'Terjadi kesalahan saat memproses penarikan.'
  } finally {
    modalLoading.value = false
  }
}
</script>

<template>
  <div class="status-dana-container">
    <div class="header-section">
      <h2 class="page-title">Status Dana & Penarikan</h2>
      <p class="page-desc">Kelola dan pantau seluruh donasi yang masuk serta ajukan penarikan dana ke rekening Anda.</p>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Memuat data...</p>
    </div>
    
    <template v-else>
      <!-- Summary Cards -->
      <div class="summary-cards">
        <div class="card glass-card card-total">
          <div class="card-icon">💰</div>
          <div class="card-content">
            <span class="card-label">Total Dana Terkumpul</span>
            <h3 class="card-value">{{ formatCurrency(summary.total_terkumpul) }}</h3>
            <p class="card-info">Semua donasi yang telah terverifikasi</p>
          </div>
        </div>

        <div class="card glass-card card-active">
          <div class="card-icon">💳</div>
          <div class="card-content">
            <span class="card-label">Saldo Siap Ditarik</span>
            <h3 class="card-value highlight">{{ formatCurrency(summary.saldo_aktif) }}</h3>
            <button 
              class="btn-tarik" 
              :disabled="summary.saldo_aktif <= 0 || summary.has_pending"
              @click="openWithdrawalModal"
            >
              {{ summary.has_pending ? 'Proses Penarikan Aktif' : 'Tarik Dana' }}
            </button>
          </div>
        </div>

        <div class="card glass-card card-pending">
          <div class="card-icon">⏳</div>
          <div class="card-content">
            <span class="card-label">Dana Dalam Proses Penarikan</span>
            <h3 class="card-value">{{ formatCurrency(summary.total_pending) }}</h3>
            <p class="card-info">Menunggu pencairan oleh Admin</p>
          </div>
        </div>
      </div>

      <!-- Transaction Table -->
      <div class="table-section">
        <h3 class="section-title">Riwayat Donasi Masuk</h3>
        <div class="table-responsive glass-panel">
          <table class="modern-table">
            <thead>
              <tr>
                <th>Donatur</th>
                <th>Campaign</th>
                <th>Tanggal</th>
                <th>Jumlah</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="donations.length === 0">
                <td colspan="5" class="empty-state">Belum ada donasi masuk.</td>
              </tr>
              <tr v-for="donasi in donations" :key="donasi.ID" class="table-row">
                <td>
                  <div class="donatur-info">
                    <span class="d-name">{{ donasi.is_anonymous ? 'Hamba Allah' : (donasi.donor_name || donasi.donatur?.name || 'Anonim') }}</span>
                  </div>
                </td>
                <td><span class="c-title">{{ donasi.campaign?.title }}</span></td>
                <td>{{ formatDate(donasi.CreatedAt) }}</td>
                <td class="t-amount">{{ formatCurrency(donasi.amount) }}</td>
                <td>
                  <span class="badge" :class="getBadgeInfo(donasi.status).class">
                    {{ getBadgeInfo(donasi.status).text }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Withdrawal Modal -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-content glass-panel">
        <button class="close-btn" @click="showModal = false">✕</button>
        <h3 class="modal-title">Tarik Dana</h3>
        <p class="modal-desc">
          Dana sebesar <strong>{{ formatCurrency(summary.saldo_aktif) }}</strong> akan ditarik ke rekening Anda.
        </p>

        <div v-if="modalError" class="alert alert-error">{{ modalError }}</div>
        <div v-if="modalSuccess" class="alert alert-success">{{ modalSuccess }}</div>

        <form @submit.prevent="submitWithdrawal" class="modal-form" v-if="!modalSuccess">
          <div class="form-group">
            <label>Nama Bank Tujuan</label>
            <input type="text" v-model="form.bank_name" placeholder="Contoh: BCA, Mandiri, BRI" required>
          </div>
          <div class="form-group">
            <label>Nomor Rekening</label>
            <input type="text" v-model="form.bank_account" placeholder="Masukkan nomor rekening" required>
          </div>
          <div class="form-group">
            <label>Atas Nama Pemilik Rekening</label>
            <input type="text" v-model="form.account_name" placeholder="Sesuai buku tabungan" required>
          </div>

          <div class="modal-actions">
            <button type="button" class="btn-cancel" @click="showModal = false">Batal</button>
            <button type="submit" class="btn-submit" :disabled="modalLoading">
              {{ modalLoading ? 'Memproses...' : 'Ajukan Penarikan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ─── Variables & Fonts ─── */
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700&display=swap');

.status-dana-container {
  font-family: 'Outfit', sans-serif;
  color: #2c3e50;
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* ─── Header ─── */
.header-section {
  margin-bottom: 2rem;
}
.page-title {
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #2563eb, #7c3aed);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 0.5rem;
}
.page-desc {
  color: #64748b;
  font-size: 1.05rem;
}

/* ─── Loading ─── */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem;
  color: #64748b;
}
.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e2e8f0;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ─── Summary Cards ─── */
.summary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2.5rem;
}

.glass-card {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.4);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.05);
  border-radius: 1.25rem;
  padding: 1.5rem;
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}
.glass-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.08);
}

.card-total { border-bottom: 4px solid #3b82f6; }
.card-active { border-bottom: 4px solid #10b981; }
.card-pending { border-bottom: 4px solid #f59e0b; }

.card-icon {
  font-size: 2.5rem;
  background: white;
  border-radius: 1rem;
  padding: 0.5rem;
  box-shadow: 0 4px 6px rgba(0,0,0,0.05);
}

.card-content {
  flex: 1;
}
.card-label {
  display: block;
  font-size: 0.9rem;
  color: #64748b;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 0.5rem;
}
.card-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 0.5rem 0;
}
.card-value.highlight {
  color: #10b981;
}
.card-info {
  font-size: 0.85rem;
  color: #94a3b8;
  margin: 0;
}

.btn-tarik {
  margin-top: 0.75rem;
  width: 100%;
  padding: 0.75rem;
  border: none;
  border-radius: 0.75rem;
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}
.btn-tarik:hover:not(:disabled) {
  transform: scale(1.02);
  box-shadow: 0 6px 16px rgba(16, 185, 129, 0.4);
}
.btn-tarik:disabled {
  background: #cbd5e1;
  color: #64748b;
  box-shadow: none;
  cursor: not-allowed;
  transform: none;
}

/* ─── Table Section ─── */
.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #1e293b;
}

.glass-panel {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.04);
  border-radius: 1.25rem;
  overflow: hidden;
}

.table-responsive {
  width: 100%;
  overflow-x: auto;
}

.modern-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}
.modern-table th {
  background: rgba(248, 250, 252, 0.8);
  padding: 1.25rem 1.5rem;
  font-weight: 600;
  font-size: 0.9rem;
  color: #475569;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 2px solid #e2e8f0;
}
.modern-table td {
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}
.table-row {
  transition: background 0.2s ease;
}
.table-row:hover {
  background: rgba(241, 245, 249, 0.5);
}
.d-name {
  font-weight: 600;
  color: #1e293b;
}
.c-title {
  color: #64748b;
  font-size: 0.95rem;
}
.t-amount {
  font-weight: 600;
  color: #10b981;
}

.badge {
  display: inline-block;
  padding: 0.4rem 0.8rem;
  border-radius: 2rem;
  font-size: 0.85rem;
  font-weight: 500;
}
.badge-waiting { background: #fef3c7; color: #b45309; }
.badge-verified { background: #dcfce7; color: #15803d; }
.badge-rejected { background: #fee2e2; color: #b91c1c; }
.badge-default { background: #f1f5f9; color: #475569; }

.empty-state {
  text-align: center;
  padding: 3rem !important;
  color: #94a3b8;
  font-style: italic;
}

/* ─── Modal ─── */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease-out;
}

.modal-content {
  background: white;
  width: 100%;
  max-width: 450px;
  padding: 2.5rem;
  border-radius: 1.5rem;
  position: relative;
  box-shadow: 0 25px 50px rgba(0,0,0,0.15);
  transform: translateY(0);
  animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.close-btn {
  position: absolute;
  top: 1.25rem; right: 1.25rem;
  background: none;
  border: none;
  font-size: 1.25rem;
  color: #94a3b8;
  cursor: pointer;
  transition: color 0.2s;
}
.close-btn:hover { color: #0f172a; }

.modal-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 0 0.5rem 0;
  color: #1e293b;
}
.modal-desc {
  color: #64748b;
  margin-bottom: 1.5rem;
  line-height: 1.5;
}

.alert {
  padding: 1rem;
  border-radius: 0.75rem;
  margin-bottom: 1.5rem;
  font-size: 0.95rem;
}
.alert-error { background: #fee2e2; color: #b91c1c; border: 1px solid #fca5a5; }
.alert-success { background: #dcfce7; color: #15803d; border: 1px solid #86efac; }

.modal-form { display: flex; flex-direction: column; gap: 1.25rem; }
.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: 0.9rem; font-weight: 500; color: #475569; }
.form-group input {
  padding: 0.875rem 1rem;
  border: 1px solid #cbd5e1;
  border-radius: 0.75rem;
  font-family: inherit;
  font-size: 1rem;
  transition: all 0.2s;
  background: #f8fafc;
}
.form-group input:focus {
  outline: none;
  border-color: #3b82f6;
  background: white;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.modal-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}
.btn-cancel {
  flex: 1;
  padding: 0.875rem;
  border: 1px solid #cbd5e1;
  background: white;
  color: #475569;
  border-radius: 0.75rem;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.2s;
}
.btn-cancel:hover { background: #f1f5f9; }
.btn-submit {
  flex: 2;
  padding: 0.875rem;
  border: none;
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: white;
  border-radius: 0.75rem;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}
.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(59, 130, 246, 0.4);
}
.btn-submit:disabled {
  background: #94a3b8;
  box-shadow: none;
  cursor: not-allowed;
  transform: none;
}

@media (max-width: 768px) {
  .modal-content { padding: 1.5rem; margin: 1rem; }
  .modern-table th, .modern-table td { padding: 1rem; }
}
</style>
