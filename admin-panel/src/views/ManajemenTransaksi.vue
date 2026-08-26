<!--
  === GO BACKEND INTEGRATION CHECKLIST ===
  Route: GET /api/transactions
  Auth: Must verify JWT and ensure role is admin.
  Response: JSON array of transactions:
  [
    {
      "id": 1,
      "campaign": { "title": "..." },
      "donatur": { "name": "..." },
      "amount": 50000,
      "method": "BCA Transfer",
      "status": "berhasil" // or "menunggu", "gagal"
    }
  ]
-->
<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">Manajemen Transaksi</h1>
        <p class="page-subtitle">Pantau seluruh transaksi donasi di platform.</p>
      </div>
      <div class="header-actions">
        <div class="search-box">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
          <input
            v-model="searchQuery"
            type="search"
            placeholder="Cari nama donatur..."
            class="search-input"
          />
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-toolbar">
        <div class="toolbar-left">
          <span class="result-count" v-if="!isLoading">
            Menampilkan <strong>{{ filteredTransactions.length }}</strong> transaksi
          </span>
        </div>
        <div class="toolbar-right">
          <select v-model="statusFilter" class="filter-select">
            <option value="">Semua Status</option>
            <option value="berhasil">Berhasil</option>
            <option value="menunggu">Menunggu</option>
            <option value="gagal">Gagal</option>
          </select>
        </div>
      </div>

      <!-- Error State -->
      <div v-if="errorMessage" class="empty-state">
        <div class="empty-icon-wrap" style="color: #ef4444; background: #fef2f2;">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <circle cx="12" cy="16" r="1" fill="currentColor"/>
          </svg>
        </div>
        <p class="empty-title">Gagal memuat transaksi</p>
        <p class="empty-sub">{{ errorMessage }}</p>
        <button class="btn-ghost" @click="fetchTransactions">Coba Lagi</button>
      </div>

      <!-- Loading skeleton -->
      <div v-else-if="isLoading" class="skeleton-table">
        <div v-for="i in 5" :key="i" class="skeleton-row">
          <div class="skel skel--sm"></div>
          <div class="skel skel--lg"></div>
          <div class="skel skel--md"></div>
          <div class="skel skel--sm"></div>
          <div class="skel skel--sm"></div>
        </div>
      </div>

      <!-- Table -->
      <div v-else-if="filteredTransactions.length > 0" class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th class="th-id">ID Transaksi</th>
              <th>Campaign</th>
              <th>Donatur</th>
              <th>Jumlah</th>
              <th>Metode</th>
              <th class="th-center">Bukti Transfer</th>
              <th class="th-center">Status</th>
              <th class="th-center">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="t in filteredTransactions" :key="t.id" class="table-row">
              <td class="td-id">TRX-{{ String(t.id || t.ID).padStart(5, '0') }}</td>
              <td class="td-text font-medium">{{ t.campaign?.title || t.Campaign?.Title || t.campaign_title || '—' }}</td>
              <td class="td-text">{{ t.donatur?.name || t.Donatur?.Name || t.donatur_name || 'Anonim' }}</td>
              <td class="td-text font-medium">Rp {{ formatNumber(t.amount || t.Amount || t.jumlah_donasi) }}</td>
              <td class="td-text">{{ t.method || t.PaymentMethod || t.metode_pembayaran || 'Transfer' }}</td>
              <td class="td-text th-center">
                <a v-if="t.bukti_transfer || t.BuktiTransfer" :href="`http://localhost:8080${t.bukti_transfer || t.BuktiTransfer}`" target="_blank" class="btn-bukti">Lihat Bukti</a>
                <span v-else class="text-muted">—</span>
              </td>
              <td class="td-status th-center">
                <span class="badge-status" :class="`badge-status--${(t.status || t.Status) === 'Sudah Dicairkan' ? 'dicairkan' : (t.status || t.Status)?.toLowerCase()}`">
                  {{ capitalize(t.status || t.Status) }}
                </span>
              </td>
              <td class="th-center">
                <button
                  v-if="(t.status || t.Status) === 'Berhasil' || (t.status || t.Status) === 'berhasil'"
                  class="btn-cairkan"
                  @click="cairkanDana(t.id || t.ID)"
                  :disabled="processingId === (t.id || t.ID)"
                >
                  {{ processingId === (t.id || t.ID) ? 'Memproses...' : 'Cairkan Dana' }}
                </button>
                <span v-else-if="(t.status || t.Status) === 'Sudah Dicairkan'" class="text-success-muted">✓ Selesai</span>
                <span v-else class="text-muted">—</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Empty State -->
      <div v-else class="empty-state">
        <div class="empty-icon-wrap">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <line x1="12" y1="1" x2="12" y2="23"/>
            <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
          </svg>
        </div>
        <p class="empty-title">Tidak ada transaksi</p>
        <p class="empty-sub">Coba sesuaikan filter atau kata kunci pencarian.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { adminApi } from '@/services/api'
import { useAuth } from '@/composables/useAuth'

const { getToken } = useAuth()

const transactions = ref([])
const isLoading = ref(true)
const errorMessage = ref('')
const statusFilter = ref('')
const searchQuery = ref('')
const processingId = ref(null)

const filteredTransactions = computed(() => {
  let list = transactions.value
  if (statusFilter.value) {
    list = list.filter(t => (t.status || t.Status)?.toLowerCase() === statusFilter.value.toLowerCase())
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    list = list.filter(t => {
      const donatur = (t.donatur?.name || t.Donatur?.Name || t.donatur_name || '').toLowerCase()
      return donatur.includes(q)
    })
  }
  return list
})

async function fetchTransactions() {
  isLoading.value = true
  errorMessage.value = ''
  try {
    transactions.value = await adminApi.getTransactions()
  } catch (err) {
    errorMessage.value = err.message || 'Tidak dapat memuat transaksi.'
  } finally {
    isLoading.value = false
  }
}

async function cairkanDana(id) {
  if (!confirm('Apakah Anda yakin ingin mencairkan dana transaksi ini?')) {
    return
  }
  
  processingId.value = id
  try {
    await adminApi.cairkanTransaction(id)
    
    // Update state directly
    const target = transactions.value.find(t => (t.id || t.ID) === id)
    if (target) {
      if(target.Status !== undefined) {
         target.Status = 'Sudah Dicairkan'
      } else {
         target.status = 'Sudah Dicairkan'
      }
    }
  } catch (err) {
    alert(err.message || 'Gagal mencairkan dana.')
  } finally {
    processingId.value = null
  }
}

function formatNumber(num) {
  return Number(num || 0).toLocaleString('id-ID')
}

function capitalize(str) {
  return str ? str.charAt(0).toUpperCase() + str.slice(1) : '—'
}

onMounted(fetchTransactions)
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');

* { box-sizing: border-box; margin: 0; padding: 0; }
.page { font-family: 'Inter', sans-serif; display: flex; flex-direction: column; gap: 20px; }

.page-header { display: flex; justify-content: space-between; flex-wrap: wrap; gap: 16px; }
.page-title { font-size: 1.35rem; font-weight: 800; color: #0f172a; letter-spacing: -0.4px; }
.page-subtitle { font-size: 0.8rem; color: #94a3b8; margin-top: 3px; }

.search-box {
  display: flex; align-items: center; gap: 8px; border: 1.5px solid #e2e8f0;
  border-radius: 9px; padding: 8px 13px; background: #fff; color: #94a3b8;
}
.search-box:focus-within { border-color: #3b82f6; box-shadow: 0 0 0 3px rgba(59,130,246,0.1); color: #3b82f6; }
.search-input { border: none; background: transparent; outline: none; font-size: 0.825rem; color: #1e293b; width: 200px; }

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

.th-id { width: 90px; }
.td-id { font-size: 0.72rem; color: #94a3b8; font-weight: 600; font-family: monospace; }
.th-center { text-align: center !important; }
.font-medium { font-weight: 600; color: #1e293b; }
.td-text { color: #475569; }

.badge-status {
  display: inline-block; padding: 3px 10px; border-radius: 20px;
  font-size: 0.7rem; font-weight: 700;
}
.badge-status--berhasil { background: #dcfce7; color: #166534; }
.badge-status--dicairkan { background: #e0e7ff; color: #3730a3; }
.badge-status--menunggu, .badge-status--menunggu\ verifikasi { background: #fef9c3; color: #854d0e; }
.badge-status--gagal, .badge-status--ditolak { background: #fee2e2; color: #991b1b; }

.text-muted { color: #94a3b8; }
.text-success-muted { color: #10b981; font-weight: 600; font-size: 0.75rem; }

.btn-cairkan {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 5px 12px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}
.btn-cairkan:hover:not(:disabled) { background: #2563eb; }
.btn-cairkan:disabled { background: #94a3b8; cursor: not-allowed; }

.btn-bukti {
  display: inline-block;
  padding: 4px 10px;
  background: #eff6ff;
  color: #2563eb;
  border-radius: 6px;
  text-decoration: none;
  font-size: 0.7rem;
  font-weight: 600;
  transition: all 0.2s;
}
.btn-bukti:hover { background: #dbeafe; }

.empty-state { padding: 56px 24px; display: flex; flex-direction: column; align-items: center; gap: 8px; text-align: center; }
.empty-icon-wrap { width: 56px; height: 56px; background: #f1f5f9; border-radius: 14px; display: flex; align-items: center; justify-content: center; color: #94a3b8; }
.empty-title { font-size: 0.9rem; font-weight: 700; color: #475569; }
.empty-sub { font-size: 0.78rem; color: #94a3b8; }

.skeleton-table { padding: 16px 20px; display: flex; flex-direction: column; gap: 14px; }
.skeleton-row { display: flex; gap: 12px; align-items: center; }
.skel { height: 14px; border-radius: 5px; background: #f1f5f9; animation: pulse 1.4s infinite; }
.skel--sm { flex: 0.5; } .skel--md { flex: 1; } .skel--lg { flex: 2; }
@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.45; } }
</style>
