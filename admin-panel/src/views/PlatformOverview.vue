<template>
  <div class="overview-page">

    <!-- ── Page Header ── -->
    <div class="page-header">
      <div class="page-header-text">
        <h1 class="page-title">Dashboard Overview</h1>
        <p class="page-subtitle">Ringkasan aktivitas platform OpenDonasi hari ini.</p>
      </div>
      <button class="refresh-btn" :class="{ 'refresh-btn--spinning': isRefreshing }" @click="fetchAll" title="Refresh data">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <path d="M23 4v6h-6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M1 20v-6h6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        Refresh
      </button>
    </div>

    <!-- ── Stats Cards ── -->
    <section class="stats-grid" aria-label="Statistik platform">
      <div
        v-for="card in statCards"
        :key="card.key"
        class="stat-card"
        :style="{ '--accent': card.color, '--accent-soft': card.colorSoft }"
      >
        <div class="stat-card-icon" aria-hidden="true" v-html="card.icon"></div>
        <div class="stat-card-body">
          <p class="stat-label">{{ card.label }}</p>
          <p class="stat-value">
            <span v-if="statsLoading" class="skeleton-text">—</span>
            <template v-else>{{ card.formatted }}</template>
          </p>
          <p class="stat-sub">{{ card.sub }}</p>
        </div>
        <div class="stat-card-glow" aria-hidden="true"></div>
      </div>
    </section>

    <!-- ── Recent Donations Table ── -->
    <section class="table-section">
      <div class="section-header">
        <div class="section-header-left">
          <h2 class="section-title">Transaksi Donasi Terbaru</h2>
          <span class="section-count" v-if="!donationsLoading">{{ recentDonations.length }} entri</span>
        </div>
        <div class="table-filters">
          <div class="search-box">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <circle cx="11" cy="11" r="8" stroke="currentColor" stroke-width="1.5"/>
              <line x1="21" y1="21" x2="16.65" y2="16.65" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            </svg>
            <input
              v-model="searchQuery"
              type="search"
              placeholder="Cari donatur atau campaign..."
              class="search-input"
              aria-label="Cari transaksi"
            />
          </div>
          <select v-model="statusFilter" class="filter-select" aria-label="Filter status">
            <option value="">Semua Status</option>
            <option value="Berhasil">Berhasil</option>
            <option value="Menunggu Verifikasi">Menunggu</option>
            <option value="Ditolak">Ditolak</option>
          </select>
        </div>
      </div>

      <!-- Loading Skeleton -->
      <div v-if="donationsLoading" class="table-skeleton">
        <div v-for="i in 5" :key="i" class="skeleton-row">
          <div class="skeleton-cell skeleton-cell--wide"></div>
          <div class="skeleton-cell skeleton-cell--mid"></div>
          <div class="skeleton-cell skeleton-cell--mid"></div>
          <div class="skeleton-cell skeleton-cell--narrow"></div>
          <div class="skeleton-cell skeleton-cell--narrow"></div>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="donationsError" class="empty-state">
        <svg width="40" height="40" viewBox="0 0 24 24" fill="none" class="empty-icon" aria-hidden="true">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5"/>
          <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <circle cx="12" cy="16" r="1" fill="currentColor"/>
        </svg>
        <p class="empty-title">Gagal memuat data</p>
        <p class="empty-sub">{{ donationsError }}</p>
        <button class="retry-btn" @click="fetchDonations">Coba Lagi</button>
      </div>

      <!-- Table -->
      <div v-else-if="filteredDonations.length > 0" class="table-wrapper">
        <table class="data-table" aria-label="Tabel transaksi donasi terbaru">
          <thead>
            <tr>
              <th>Tanggal</th>
              <th>Campaign</th>
              <th>Donatur</th>
              <th class="text-right">Total Donasi</th>
              <th class="text-center">Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="tx in filteredDonations" :key="tx.id || tx.ID" class="table-row">
              <td class="td-date">
                <span class="date-main">{{ formatDate(tx.created_at || tx.CreatedAt) }}</span>
                <span class="date-time">{{ formatTime(tx.created_at || tx.CreatedAt) }}</span>
              </td>
              <td class="td-campaign">
                <span class="campaign-name">{{ tx.campaign?.title || tx.Campaign?.Title || '—' }}</span>
                <span class="campaign-cat">{{ tx.campaign?.category || tx.Campaign?.Category || '' }}</span>
              </td>
              <td class="td-donor">
                <div class="donor-chip">
                  <div class="donor-avatar" aria-hidden="true">{{ donorInitial(tx) }}</div>
                  <span>{{ (tx.is_anonymous || tx.IsAnonymous) ? 'Hamba Allah' : (tx.donor_name || tx.Donatur?.Name || '—') }}</span>
                </div>
              </td>
              <td class="td-amount text-right">{{ formatCurrency(tx.amount || tx.Amount) }}</td>
              <td class="td-status text-center">
                <span class="badge" :class="statusClass(tx.status || tx.Status)">{{ tx.status || tx.Status }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Empty State -->
      <div v-else class="empty-state">
        <svg width="40" height="40" viewBox="0 0 24 24" fill="none" class="empty-icon" aria-hidden="true">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <p class="empty-title">Tidak ada transaksi</p>
        <p class="empty-sub">{{ searchQuery || statusFilter ? 'Tidak ada hasil yang cocok.' : 'Belum ada donasi yang masuk.' }}</p>
      </div>
    </section>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { adminApi } from '@/services/api'
// ─── State ───────────────────────────────────────
const stats          = ref(null)
const donations      = ref([])
const statsLoading   = ref(true)
const donationsLoading = ref(true)
const donationsError = ref('')
const isRefreshing   = ref(false)
const searchQuery    = ref('')
const statusFilter   = ref('')

// ─── Computed: Stat Cards ──────────────────────
const statCards = computed(() => [
  {
    key: 'donatur',
    label: 'Total Donatur',
    formatted: formatNumber(stats.value?.total_donatur ?? 0),
    sub: 'Pengguna terdaftar sebagai donatur',
    color: '#3b82f6',
    colorSoft: 'rgba(59,130,246,0.1)',
    icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none">
      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
      <circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="1.5"/>
    </svg>`,
  },
  {
    key: 'campaign_aktif',
    label: 'Campaign Aktif',
    formatted: formatNumber(stats.value?.active_campaigns ?? 0),
    sub: 'Campaign yang sedang berjalan',
    color: '#10b981',
    colorSoft: 'rgba(16,185,129,0.1)',
    icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none">
      <path d="M12 21.593c-5.63-5.539-11-10.297-11-14.402 0-3.791 3.068-5.191 5.281-5.191 1.312 0 4.151.501 5.719 4.457 1.59-3.968 4.464-4.447 5.726-4.447 2.54 0 5.274 1.621 5.274 5.181 0 4.069-5.136 8.625-11 14.402z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>
    </svg>`,
  },
  {
    key: 'campaign_selesai',
    label: 'Campaign Selesai',
    formatted: formatNumber(stats.value?.completed_campaigns ?? 0),
    sub: 'Campaign yang telah tercapai targetnya',
    color: '#8b5cf6',
    colorSoft: 'rgba(139,92,246,0.1)',
    icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none">
      <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
      <polyline points="22 4 12 14.01 9 11.01" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>`,
  },
  {
    key: 'total_dana',
    label: 'Total Dana Terkumpul',
    formatted: formatCurrency(stats.value?.total_collected ?? 0),
    sub: 'Akumulasi seluruh donasi masuk',
    color: '#f59e0b',
    colorSoft: 'rgba(245,158,11,0.1)',
    icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none">
      <line x1="12" y1="1" x2="12" y2="23" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>`,
  },
])

// ─── Computed: Filtered Donations ─────────────
const recentDonations = computed(() =>
  [...donations.value]
    .sort((a, b) => new Date(b.created_at || b.CreatedAt) - new Date(a.created_at || a.CreatedAt))
    .slice(0, 50)
)

const filteredDonations = computed(() => {
  let list = recentDonations.value
  if (statusFilter.value) {
    list = list.filter(d => (d.status || d.Status)?.toLowerCase() === statusFilter.value.toLowerCase())
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    list = list.filter(d =>
      (d.donor_name || d.Donatur?.Name || '').toLowerCase().includes(q) ||
      (d.campaign?.title || d.Campaign?.Title || '').toLowerCase().includes(q)
    )
  }
  return list
})

// ─── API Calls ─────────────────────────────────
async function fetchStats() {
  statsLoading.value = true
  try {
    stats.value = await adminApi.getStatistics()
  } catch (e) {
    console.error('[PlatformOverview] Stats error:', e)
    // Fallback to zeros — UI already handles null gracefully
  } finally {
    statsLoading.value = false
  }
}

async function fetchDonations() {
  donationsLoading.value = true
  donationsError.value = ''
  try {
    donations.value = await adminApi.getDonations()
  } catch (e) {
    // adminApi throws Error with server message already
    if (e.message?.toLowerCase().includes('forbidden') || e.message?.toLowerCase().includes('admin')) {
      donationsError.value = 'Akses ditolak. Hanya admin yang dapat melihat data ini.'
    } else {
      donationsError.value = e.message || 'Tidak dapat memuat data transaksi.'
    }
    console.error('[PlatformOverview] Donations error:', e)
  } finally {
    donationsLoading.value = false
  }
}

async function fetchAll() {
  isRefreshing.value = true
  await Promise.all([fetchStats(), fetchDonations()])
  // Keep spin animation visible briefly
  setTimeout(() => { isRefreshing.value = false }, 600)
}

// ─── Helpers ───────────────────────────────────
function formatCurrency(val) {
  if (val === null || val === undefined) return 'Rp 0'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(val)
}

function formatNumber(val) {
  if (val === null || val === undefined) return '0'
  return new Intl.NumberFormat('id-ID').format(val)
}

function formatDate(iso) {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

function formatTime(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

function donorInitial(tx) {
  if (tx.is_anonymous || tx.IsAnonymous) return '?'
  return (tx.donor_name || tx.Donatur?.Name || 'D').charAt(0).toUpperCase()
}

function statusClass(status) {
  if (!status) return 'badge--grey'
  const s = status.toLowerCase()
  if (s.includes('berhasil'))          return 'badge--green'
  if (s.includes('menunggu'))          return 'badge--yellow'
  if (s.includes('ditolak'))           return 'badge--red'
  return 'badge--grey'
}

// ─── Lifecycle ─────────────────────────────────
onMounted(fetchAll)
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');

* { box-sizing: border-box; margin: 0; padding: 0; }

/* ═══════════════════════════════════
   PAGE
═══════════════════════════════════ */
.overview-page {
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ── Page Header ── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.page-title {
  font-size: 1.35rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.4px;
  line-height: 1.2;
}

.page-subtitle {
  font-size: 0.8rem;
  color: #94a3b8;
  margin-top: 3px;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.8rem;
  font-weight: 600;
  color: #475569;
  cursor: pointer;
  transition: background 0.15s, color 0.15s, border-color 0.15s;
  flex-shrink: 0;
}

.refresh-btn:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
  color: #1e293b;
}

.refresh-btn svg {
  transition: transform 0.6s ease;
}

.refresh-btn--spinning svg {
  animation: spin 0.7s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* ═══════════════════════════════════
   STATS GRID
═══════════════════════════════════ */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  position: relative;
  background: #fff;
  border-radius: 14px;
  border: 1px solid #f1f5f9;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04), 0 4px 12px rgba(0,0,0,0.04);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0,0,0,0.08);
}

.stat-card-glow {
  position: absolute;
  top: -30px;
  right: -20px;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: var(--accent-soft);
  pointer-events: none;
}

.stat-card-icon {
  width: 40px;
  height: 40px;
  background: var(--accent-soft);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent);
  flex-shrink: 0;
}

.stat-card-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.5px;
  line-height: 1.2;
}

.stat-sub {
  font-size: 0.72rem;
  color: #94a3b8;
  margin-top: 2px;
}

.skeleton-text {
  color: #e2e8f0;
  animation: pulse 1.4s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50%       { opacity: 0.4; }
}

/* ═══════════════════════════════════
   TABLE SECTION
═══════════════════════════════════ */
.table-section {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #f1f5f9;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04), 0 4px 12px rgba(0,0,0,0.04);
  overflow: hidden;
}

/* ── Section Header ── */
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 18px 20px;
  border-bottom: 1px solid #f1f5f9;
  flex-wrap: wrap;
}

.section-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.section-title {
  font-size: 0.95rem;
  font-weight: 700;
  color: #0f172a;
}

.section-count {
  font-size: 0.72rem;
  font-weight: 600;
  color: #64748b;
  background: #f1f5f9;
  border-radius: 20px;
  padding: 2px 9px;
}

.table-filters {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  border: 1.5px solid #e2e8f0;
  border-radius: 8px;
  padding: 7px 12px;
  background: #f8fafc;
  color: #94a3b8;
  transition: border-color 0.15s;
}

.search-box:focus-within {
  border-color: #3b82f6;
  background: #fff;
}

.search-input {
  border: none;
  background: transparent;
  outline: none;
  font-family: inherit;
  font-size: 0.8rem;
  color: #1e293b;
  width: 200px;
}

.search-input::placeholder { color: #cbd5e1; }

.filter-select {
  border: 1.5px solid #e2e8f0;
  border-radius: 8px;
  padding: 7px 12px;
  font-family: inherit;
  font-size: 0.8rem;
  color: #475569;
  background: #f8fafc;
  cursor: pointer;
  outline: none;
  transition: border-color 0.15s;
}

.filter-select:focus { border-color: #3b82f6; }

/* ── Skeleton ── */
.table-skeleton {
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.skeleton-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.skeleton-cell {
  height: 14px;
  background: #f1f5f9;
  border-radius: 4px;
  animation: pulse 1.4s ease-in-out infinite;
}

.skeleton-cell--wide   { flex: 2; }
.skeleton-cell--mid    { flex: 1.5; }
.skeleton-cell--narrow { flex: 1; }

/* ── Table ── */
.table-wrapper {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.825rem;
}

.data-table thead th {
  padding: 11px 16px;
  text-align: left;
  font-size: 0.7rem;
  font-weight: 700;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  background: #f8fafc;
  border-bottom: 1px solid #f1f5f9;
  white-space: nowrap;
}

.text-right { text-align: right !important; }
.text-center { text-align: center !important; }

.table-row {
  border-bottom: 1px solid #f8fafc;
  transition: background 0.12s;
}

.table-row:last-child { border-bottom: none; }
.table-row:hover { background: #fafbfc; }

.data-table td {
  padding: 13px 16px;
  vertical-align: middle;
  color: #334155;
}

/* ── Column Styles ── */
.td-date {
  white-space: nowrap;
}

.date-main {
  display: block;
  font-weight: 600;
  color: #1e293b;
  font-size: 0.8rem;
}

.date-time {
  display: block;
  font-size: 0.7rem;
  color: #94a3b8;
  margin-top: 1px;
}

.td-campaign {
  max-width: 220px;
}

.campaign-name {
  display: block;
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.campaign-cat {
  display: block;
  font-size: 0.7rem;
  color: #94a3b8;
  margin-top: 1px;
}

.td-donor {}

.donor-chip {
  display: flex;
  align-items: center;
  gap: 8px;
}

.donor-avatar {
  width: 28px;
  height: 28px;
  background: linear-gradient(135deg, #3b82f6, #6366f1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.td-amount {
  font-weight: 700;
  color: #0f172a;
  white-space: nowrap;
}

/* ── Status Badges ── */
.badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 20px;
  font-size: 0.7rem;
  font-weight: 700;
  white-space: nowrap;
  letter-spacing: 0.02em;
}

.badge--green {
  background: #dcfce7;
  color: #16a34a;
}

.badge--yellow {
  background: #fef9c3;
  color: #ca8a04;
}

.badge--red {
  background: #fee2e2;
  color: #dc2626;
}

.badge--grey {
  background: #f1f5f9;
  color: #64748b;
}

/* ── Empty / Error State ── */
.empty-state {
  padding: 52px 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  text-align: center;
}

.empty-icon {
  color: #cbd5e1;
  margin-bottom: 4px;
}

.empty-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: #475569;
}

.empty-sub {
  font-size: 0.78rem;
  color: #94a3b8;
  max-width: 320px;
}

.retry-btn {
  margin-top: 8px;
  padding: 8px 18px;
  background: #3b82f6;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s;
}

.retry-btn:hover { background: #2563eb; }

/* ═══════════════════════════════════
   RESPONSIVE
═══════════════════════════════════ */
@media (max-width: 1100px) {
  .stats-grid { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 600px) {
  .stats-grid { grid-template-columns: 1fr; }

  .page-header { flex-direction: column; align-items: flex-start; }

  .section-header { flex-direction: column; align-items: flex-start; }

  .table-filters { flex-direction: column; align-items: stretch; }

  .search-input { width: 100%; }

  .stat-value { font-size: 1.2rem; }
}
</style>
