<script setup>
// ─── Backend base URL (strip /api suffix) ───
const BACKEND_BASE = (import.meta.env.VITE_API_URL || 'http://localhost:8080/api')
  .replace(/\/api$/, '')

const FALLBACK_IMG = 'https://images.unsplash.com/photo-1547683905-f686c993aae5?w=600&h=400&fit=crop'

defineProps({
  campaign: {
    type: Object,
    required: true
  }
})

// ─── Resolve foto URL: handle null, undefined, relative path, dan absolute URL ───
const getImageUrl = (foto) => {
  if (!foto || typeof foto !== 'string' || foto.trim() === '') return FALLBACK_IMG
  if (foto.startsWith('http://') || foto.startsWith('https://')) return foto
  // Path relatif dari backend: /uploads/campaigns/file.jpg
  return `${BACKEND_BASE}${foto}`
}

// Handler @error fallback (plain JS, TANPA TypeScript cast)
const onImgError = (event) => {
  if (event.target) event.target.src = FALLBACK_IMG
}

// ─── Format angka ke Rupiah ───
const formatRupiah = (num) => {
  const n = Number(num)
  if (isNaN(n)) return 'Rp 0'
  return 'Rp ' + n.toLocaleString('id-ID')
}

// ─── Ambil nilai (aman dari null/undefined, support snake_case & camelCase) ───
const getCollected = (c) => Number(c?.dana_terkumpul ?? c?.danaTerkumpul ?? 0) || 0
const getTarget    = (c) => Number(c?.target_dana    ?? c?.targetDana    ?? 0) || 1  // min 1 hindari div/0
const getLocation  = (c) => c?.alamat_lengkap || c?.alamat || '-'
const getPenggalang = (c) => {
  if (!c) return '-'
  const p = c.penggalang
  if (p && typeof p === 'object') return p.name || '-'
  return p || '-'
}

// ─── Hitung sisa hari dari end_date ───
const getRemainingDays = (c) => {
  if (!c) return 0
  if (c.sisaHari !== undefined) return c.sisaHari      // mock data compat
  if (!c.end_date) return 0
  const end = new Date(c.end_date)
  if (isNaN(end.getTime())) return 0                    // invalid date guard
  const now = new Date()
  const diff = Math.ceil((end.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  return Math.max(0, diff)
}

// ─── Status DINAMIS — aturan bisnis identik dengan dashboard ───
// Selesai = (dana_terkumpul >= target_dana) ATAU (sisa_hari <= 0)
// Aktif   = (dana_terkumpul <  target_dana) DAN  (sisa_hari >  0)
const getCampaignStatus = (c) => {
  if (!c) return 'Aktif'
  const terkumpul = getCollected(c)
  const target    = getTarget(c)
  const sisaHari  = getRemainingDays(c)
  if (terkumpul >= target || sisaHari <= 0) return 'Selesai'
  return 'Aktif'
}

// ─── CSS class badge berdasarkan status dinamis ───
const getStatusClass = (status) => {
  if (status === 'Aktif')   return 'status-aktif'
  if (status === 'Selesai') return 'status-selesai'
  return 'status-menunggu'
}

// ─── Persen progres (aman dari NaN & Infinity) ───
const getProgressPct = (c) => {
  const pct = (getCollected(c) / getTarget(c)) * 100
  return Math.min(isNaN(pct) ? 0 : pct, 100)
}
</script>

<template>
  <div class="campaign-card" :id="'campaign-card-' + (campaign?.ID || campaign?.id)">
    <div class="card-image">
      <img
        :src="getImageUrl(campaign?.foto)"
        :alt="campaign?.title || 'Campaign'"
        @error="onImgError"
      />
      <span class="card-category">{{ campaign?.category || '' }}</span>
      <span class="card-status" :class="getStatusClass(getCampaignStatus(campaign))">
        {{ getCampaignStatus(campaign) }}
      </span>
    </div>
    <div class="card-body">
      <h3 class="card-title">{{ campaign?.title }}</h3>
      <p class="card-location">📍 {{ getLocation(campaign) }}</p>

      <div class="progress-section">
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: getProgressPct(campaign) + '%' }"></div>
        </div>
        <div class="progress-info">
          <span class="progress-amount">{{ formatRupiah(getCollected(campaign)) }}</span>
          <span class="progress-percent">{{ Math.round(getProgressPct(campaign)) }}%</span>
        </div>
        <p class="progress-target">dari {{ formatRupiah(getTarget(campaign)) }}</p>
      </div>

      <div class="card-meta">
        <span class="meta-item">👤 {{ getPenggalang(campaign) }}</span>
        <span class="meta-item">⏳ {{ getRemainingDays(campaign) }} hari lagi</span>
      </div>

      <router-link
        :to="'/campaign/' + (campaign?.ID || campaign?.id)"
        class="btn-detail"
        :id="'btn-detail-' + (campaign?.ID || campaign?.id)"
      >
        Lihat Detail
      </router-link>
    </div>
  </div>
</template>

<style scoped>
.campaign-card {
  background: var(--bg-surface);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  transition: all var(--transition-base);
  display: flex;
  flex-direction: column;
  height: 100%;
}

.campaign-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
  border-color: var(--primary-glow);
}

.card-image {
  position: relative;
  height: 200px;
  overflow: hidden;
}

.card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-slow);
}

.campaign-card:hover .card-image img {
  transform: scale(1.05);
}

.card-category {
  position: absolute;
  top: 12px;
  left: 12px;
  background: rgba(0, 0, 0, 0.55);
  color: white;
  font-size: 12px;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: var(--radius-full);
  backdrop-filter: blur(4px);
}

.card-status {
  position: absolute;
  top: 12px;
  right: 12px;
  font-size: 11px;
  font-weight: 700;
  padding: 4px 10px;
  border-radius: var(--radius-full);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.status-aktif    { background: var(--success-light); color: var(--success); }
.status-menunggu { background: var(--warning-light); color: var(--warning); }
.status-tercapai { background: var(--info-light);    color: var(--info); }
.status-berakhir { background: var(--danger-light);  color: var(--danger); }
.status-disalurkan { background: var(--info-light);  color: var(--info); }
.status-laporan  { background: var(--warning-light); color: var(--warning); }
/* Selesai = Biru/Abu (sama dengan info) — konsisten dengan dashboard */
.status-selesai  { background: var(--info-light);    color: var(--info); }

.card-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}

.card-title {
  font-family: var(--font-heading);
  font-size: 17px;
  font-weight: 700;
  color: var(--text-heading);
  margin-bottom: 6px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-location {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 16px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.progress-section {
  margin-bottom: 16px;
}

.progress-bar {
  width: 100%;
  height: 6px;
  background: var(--border-light);
  border-radius: var(--radius-full);
  overflow: hidden;
  margin-bottom: 8px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--primary), #38BDF8);
  border-radius: var(--radius-full);
  transition: width 0.8s ease;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.progress-amount {
  font-weight: 700;
  font-size: 15px;
  color: var(--primary);
}

.progress-percent {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
}

.progress-target {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 2px;
}

.card-meta {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 16px;
  padding-top: 12px;
  border-top: 1px solid var(--border-light);
}

.btn-detail {
  display: block;
  text-align: center;
  background: var(--primary);
  color: white;
  font-weight: 600;
  font-size: 14px;
  padding: 10px 20px;
  border-radius: var(--radius-sm);
  text-decoration: none;
  transition: all var(--transition-fast);
  margin-top: auto;
}

.btn-detail:hover {
  background: var(--primary-hover);
  box-shadow: 0 4px 12px var(--primary-glow);
  color: white;
}
</style>
