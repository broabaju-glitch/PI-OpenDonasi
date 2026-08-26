<!--
  === GO BACKEND INTEGRATION CHECKLIST ===
  Route 1: GET /api/campaigns
  Auth: Must verify JWT and ensure role is admin.
  Response: JSON array of campaigns: [{ "id": 1, "title": "...", "penggalang": { "name": "..." }, "target_dana": 1000, "status": "pending" }]

  Route 2: PATCH /api/campaigns/:id/status
  Auth: Must verify JWT and ensure role is admin.
  Body: { "status": "disetujui" } or { "status": "ditolak" }
  Response: JSON object with success message.
-->
<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">Moderasi Campaign</h1>
        <p class="page-subtitle">Tinjau dan setujui atau tolak pengajuan campaign baru.</p>
      </div>
      <div class="header-actions">
        <button class="btn-ghost" @click="fetchCampaigns">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.59-9.22l-5.36 5.36"/>
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <div class="card">
      <div class="card-toolbar">
        <div class="toolbar-left">
          <span class="result-count" v-if="!isLoading">
            Menampilkan <strong>{{ campaigns.length }}</strong> campaign
          </span>
        </div>
      </div>

      <!-- Error State -->
      <div v-if="errorMessage" class="empty-state">
        <div class="empty-icon-wrap" aria-hidden="true" style="color: #ef4444; background: #fef2f2;">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5"/>
            <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <circle cx="12" cy="16" r="1" fill="currentColor"/>
          </svg>
        </div>
        <p class="empty-title">Gagal memuat campaign</p>
        <p class="empty-sub">{{ errorMessage }}</p>
        <button class="btn-ghost" @click="fetchCampaigns">Coba Lagi</button>
      </div>

      <!-- Loading skeleton -->
      <div v-else-if="isLoading" class="skeleton-table">
        <div v-for="i in 5" :key="i" class="skeleton-row">
          <div class="skel skel--sm"></div>
          <div class="skel skel--lg"></div>
          <div class="skel skel--md"></div>
          <div class="skel skel--sm"></div>
          <div class="skel skel--md"></div>
        </div>
      </div>

      <!-- Table -->
      <div v-else-if="campaigns.length > 0" class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th class="th-id">ID</th>
              <th>Judul Campaign</th>
              <th>Penggalang</th>
              <th>Target</th>
              <th class="th-center">Status</th>
              <th class="th-center">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="c in campaigns" :key="c.id" class="table-row">
              <td class="td-id">#{{ String(c.id || c.ID).padStart(3, '0') }}</td>
              <td class="td-name font-medium">{{ c.title || c.Title }}</td>
              <td class="td-text">{{ c.penggalang?.name || c.Penggalang?.Name || 'Tidak diketahui' }}</td>
              <td class="td-text font-medium">Rp {{ formatNumber(c.target_dana || c.TargetDana) }}</td>
              <td class="th-center">
                <span class="badge" :class="statusClass(c.status || c.Status)">{{ c.status || c.Status }}</span>
              </td>
              <td class="td-actions th-center">
                <div class="action-group">
                  <button v-if="(c.status || c.Status) === 'pending'" class="btn-action btn-approve" @click="updateStatus(c.id || c.ID, 'disetujui')" title="Setujui">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                    Setujui
                  </button>
                  <button v-if="(c.status || c.Status) === 'pending'" class="btn-action btn-reject" @click="updateStatus(c.id || c.ID, 'ditolak')" title="Tolak">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <line x1="18" y1="6" x2="6" y2="18"></line>
                      <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                    Tolak
                  </button>
                  <button class="btn-action btn-delete" @click="deleteCampaign(c.id || c.ID)" title="Hapus Campaign">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <polyline points="3 6 5 6 21 6"></polyline>
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                    </svg>
                    Hapus
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Empty State -->
      <div v-else class="empty-state">
        <div class="empty-icon-wrap" aria-hidden="true">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </div>
        <p class="empty-title">Tidak ada campaign</p>
        <p class="empty-sub">Belum ada campaign yang terdaftar.</p>
      </div>
    </div>

    <!-- Toast Notification -->
    <transition name="toast">
      <div v-if="toast.show" class="toast" :class="`toast--${toast.type}`" role="alert">
        {{ toast.message }}
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted } from 'vue'
import { adminApi } from '@/services/api'
import { useAuth } from '@/composables/useAuth'

const { getToken } = useAuth()

const campaigns = ref([])
const isLoading = ref(true)
const errorMessage = ref('')

async function fetchCampaigns() {
  isLoading.value = true
  errorMessage.value = ''
  try {
    campaigns.value = await adminApi.getCampaigns()
  } catch (err) {
    errorMessage.value = err.message || 'Tidak dapat memuat campaign.'
  } finally {
    isLoading.value = false
  }
}

async function deleteCampaign(id) {
  if (!confirm('Apakah Anda yakin ingin menghapus campaign ini?')) return

  try {
    await adminApi.deleteCampaign(id)
    campaigns.value = campaigns.value.filter(c => (c.id || c.ID) !== id)
    showToast('Campaign berhasil dihapus.', 'success')
  } catch (err) {
    showToast(`Gagal menghapus: ${err.message}`, 'error')
  }
}

async function updateStatus(id, newStatus) {
  try {
    await adminApi.updateCampaignStatus(id, newStatus)
    const idx = campaigns.value.findIndex(c => (c.id || c.ID) === id)
    if (idx !== -1) {
      if (campaigns.value[idx].Status !== undefined) {
        campaigns.value[idx].Status = newStatus === 'disetujui' ? 'aktif' : newStatus
      } else {
        campaigns.value[idx].status = newStatus === 'disetujui' ? 'aktif' : newStatus
      }
    }
    showToast(`Status campaign berhasil diubah menjadi ${newStatus}.`, 'success')
  } catch (err) {
    showToast(`Gagal mengubah status: ${err.message}`, 'error')
  }
}

// ─── Helpers ───
const toast = reactive({ show: false, message: '', type: 'success' })
let toastTimer = null
function showToast(message, type = 'success') {
  clearTimeout(toastTimer)
  toast.message = message
  toast.type = type
  toast.show = true
  toastTimer = setTimeout(() => { toast.show = false }, 3500)
}

function formatNumber(num) {
  return Number(num || 0).toLocaleString('id-ID')
}

function capitalize(str) {
  return str ? str.charAt(0).toUpperCase() + str.slice(1) : '—'
}

function statusClass(status) {
  if (!status) return 'badge--grey'
  const s = status.toLowerCase()
  if (s === 'aktif' || s === 'disetujui') return 'badge--green'
  if (s === 'pending') return 'badge--yellow'
  if (s === 'ditolak') return 'badge--red'
  if (s === 'selesai') return 'badge--blue'
  return 'badge--grey'
}

onMounted(fetchCampaigns)
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');

* { box-sizing: border-box; margin: 0; padding: 0; }
.page { font-family: 'Inter', sans-serif; display: flex; flex-direction: column; gap: 20px; }

.page-header { display: flex; justify-content: space-between; flex-wrap: wrap; gap: 16px; }
.page-title { font-size: 1.35rem; font-weight: 800; color: #0f172a; letter-spacing: -0.4px; }
.page-subtitle { font-size: 0.8rem; color: #94a3b8; margin-top: 3px; }

.btn-ghost {
  display: inline-flex; align-items: center; gap: 6px; padding: 8px 14px;
  background: #fff; color: #64748b; border: 1.5px solid #e2e8f0; border-radius: 8px;
  font-size: 0.8rem; font-weight: 600; cursor: pointer; transition: 0.15s;
}
.btn-ghost:hover { background: #f8fafc; color: #1e293b; }

.card {
  background: #fff; border-radius: 16px; border: 1px solid #f1f5f9;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.05); overflow: hidden;
}

.card-toolbar {
  display: flex; justify-content: space-between; padding: 14px 20px;
  border-bottom: 1px solid #f1f5f9; flex-wrap: wrap; gap: 12px;
}
.result-count { font-size: 0.8rem; color: #64748b; }
.filter-select {
  border: 1.5px solid #e2e8f0; border-radius: 8px; padding: 6px 10px;
  font-size: 0.78rem; color: #475569; background: #f8fafc; outline: none;
}
.filter-select:focus { border-color: #3b82f6; }

.table-wrapper { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; font-size: 0.825rem; }
.data-table th {
  padding: 11px 16px; text-align: left; font-size: 0.7rem; font-weight: 700;
  color: #64748b; text-transform: uppercase; background: #f8fafc; border-bottom: 1px solid #f1f5f9;
}
.data-table td { padding: 13px 16px; vertical-align: middle; color: #334155; border-bottom: 1px solid #f8fafc; }
.table-row:last-child td { border-bottom: none; }
.table-row:hover { background: #fafbfc; }

.th-id { width: 60px; }
.td-id { font-size: 0.72rem; color: #94a3b8; font-weight: 600; }
.th-center { text-align: center !important; }
.font-medium { font-weight: 600; color: #1e293b; }
.td-text { color: #475569; }
.text-muted { color: #94a3b8; }

.btn-action {
  display: inline-flex; align-items: center; gap: 5px; padding: 5px 10px;
  border: none; border-radius: 7px; font-size: 0.72rem; font-weight: 600;
  cursor: pointer; transition: 0.15s;
}
.btn-action:hover { transform: translateY(-1px); }
.btn-delete { background: #fee2e2; color: #dc2626; margin-left: 8px; }
.btn-delete:hover { background: #fca5a5; }

.btn-approve { background: #dcfce7; color: #16a34a; }
.btn-approve:hover { background: #bbf7d0; }

.btn-reject { background: #fef3c7; color: #d97706; }
.btn-reject:hover { background: #fde68a; }

.badge { display: inline-block; padding: 3px 8px; border-radius: 12px; font-size: 0.7rem; font-weight: 700; text-transform: uppercase; }
.badge--green { background: #dcfce7; color: #16a34a; }
.badge--yellow { background: #fef9c3; color: #ca8a04; }
.badge--red { background: #fee2e2; color: #dc2626; }
.badge--blue { background: #dbeafe; color: #2563eb; }
.badge--grey { background: #f1f5f9; color: #64748b; }

.empty-state { padding: 56px 24px; display: flex; flex-direction: column; align-items: center; gap: 8px; text-align: center; }
.empty-icon-wrap { width: 56px; height: 56px; background: #f1f5f9; border-radius: 14px; display: flex; align-items: center; justify-content: center; color: #94a3b8; }
.empty-title { font-size: 0.9rem; font-weight: 700; color: #475569; }
.empty-sub { font-size: 0.78rem; color: #94a3b8; }

.skeleton-table { padding: 16px 20px; display: flex; flex-direction: column; gap: 14px; }
.skeleton-row { display: flex; gap: 12px; align-items: center; }
.skel { height: 14px; border-radius: 5px; background: #f1f5f9; animation: pulse 1.4s ease-in-out infinite; }
.skel--sm { flex: 0.5; } .skel--md { flex: 1; } .skel--lg { flex: 2; }
@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.45; } }

.toast {
  position: fixed; bottom: 28px; right: 28px; padding: 13px 18px;
  border-radius: 12px; font-size: 0.825rem; font-weight: 600;
  box-shadow: 0 8px 24px rgba(0,0,0,0.14); z-index: 9999;
}
.toast--success { background: #0f172a; color: #fff; }
.toast--error { background: #ef4444; color: #fff; }
.toast-enter-active, .toast-leave-active { transition: opacity 0.25s, transform 0.25s; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(10px); }
</style>
