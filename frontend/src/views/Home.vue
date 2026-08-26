<script setup>
import { ref, onMounted } from 'vue'
import CampaignCard from '../components/CampaignCard.vue'
import { apiGetCampaigns, apiGetStats, apiGetStatistics } from '../services/api'

const platformStats = ref({
  activeCampaigns: 0,
  totalFundsCollected: 0,
  totalDonors: 0,
  completedCampaigns: 0
})

// Mock data fallback
const mockCampaigns = [
  { id: 1, title: 'Bantu Korban Banjir Bandang Garut', category: 'Banjir', foto: 'https://images.unsplash.com/photo-1547683905-f686c993aae5?w=600&h=400&fit=crop', status: 'Aktif', dana_terkumpul: 75000000, target_dana: 100000000, lokasi_kejadian: 'Garut, Jawa Barat', penggalang: { name: 'Yayasan Peduli' }, end_date: '2026-07-15T00:00:00Z' },
  { id: 2, title: 'Darurat Gempa Cianjur — Bangun Kembali', category: 'Gempa Bumi', foto: 'https://images.unsplash.com/photo-1599707367812-042c63e2e6f4?w=600&h=400&fit=crop', status: 'Aktif', dana_terkumpul: 250000000, target_dana: 500000000, lokasi_kejadian: 'Cianjur, Jawa Barat', penggalang: { name: 'Tim Relawan Cianjur' }, end_date: '2026-07-21T00:00:00Z' },
  { id: 3, title: 'Bantuan Logistik Erupsi Gunung Semeru', category: 'Gunung Meletus', foto: 'https://images.unsplash.com/photo-1580974852861-c381510bc98a?w=600&h=400&fit=crop', status: 'Target Tercapai', dana_terkumpul: 150000000, target_dana: 150000000, lokasi_kejadian: 'Lumajang, Jawa Timur', penggalang: { name: 'BNPB Lumajang' }, end_date: '2026-06-30T00:00:00Z' },
  { id: 4, title: 'Rehabilitasi Pasca Tsunami Palu', category: 'Tsunami', foto: 'https://images.unsplash.com/photo-1542393545-10f5cde2c810?w=600&h=400&fit=crop', status: 'Aktif', dana_terkumpul: 42000000, target_dana: 200000000, lokasi_kejadian: 'Palu, Sulawesi Tengah', penggalang: { name: 'Komunitas Palu Bangkit' }, end_date: '2026-07-30T00:00:00Z' },
  { id: 5, title: 'Tanggap Darurat Longsor Sumedang', category: 'Tanah Longsor', foto: 'https://images.unsplash.com/photo-1613842127673-3f3e5b45e4c7?w=600&h=400&fit=crop', status: 'Aktif', dana_terkumpul: 18000000, target_dana: 80000000, lokasi_kejadian: 'Sumedang, Jawa Barat', penggalang: { name: 'Relawan Sumedang' }, end_date: '2026-07-07T00:00:00Z' },
  { id: 6, title: 'Bangun Sekolah Darurat Korban Bencana', category: 'Rehabilitasi', foto: 'https://images.unsplash.com/photo-1497633762265-9d179a990aa6?w=600&h=400&fit=crop', status: 'Aktif', dana_terkumpul: 320000000, target_dana: 400000000, lokasi_kejadian: 'NTT, Indonesia', penggalang: { name: 'Yayasan Cerdas Bangsa' }, end_date: '2026-08-15T00:00:00Z' },
]

const campaigns = ref(mockCampaigns)
const loading = ref(true)
const apiError = ref('')

const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(value || 0);
}

onMounted(async () => {
  // Fetch real data from backend
  try {
    const data = await apiGetCampaigns()
    if (data && data.length > 0) {
      campaigns.value = data.slice(0, 6) // show first 6
    }
  } catch (err) {
    apiError.value = 'Menggunakan data demo — backend belum terhubung'
    console.warn('[Home] API unreachable, using mock data:', err.message)
  }

  // Fetch stats dari endpoint /statistics (real data dari DB)
  try {
    const s = await apiGetStatistics()
    platformStats.value = {
      activeCampaigns:     s.active_campaigns    || 0,
      totalFundsCollected: s.total_funds         || s.total_collected || 0,
      totalDonors:         s.total_donors        || s.total_donatur   || 0,
      completedCampaigns:  s.completed_campaigns || 0
    }
    console.log('[Home] Stats loaded:', platformStats.value)
  } catch (err) {
    console.warn('[Home] Stats API unreachable:', err.message)
  }

  loading.value = false
})

const howItWorks = ref([
  { step: '01', icon: '📝', title: 'Buat Campaign', desc: 'Penggalang membuat campaign bantuan dan menunggu verifikasi admin.' },
  { step: '02', icon: '💳', title: 'Donasi Masuk', desc: 'Donatur memilih campaign dan melakukan transfer ke rekening bersama (escrow).' },
  { step: '03', icon: '✅', title: 'Verifikasi Admin', desc: 'Admin memverifikasi bukti transfer dan memvalidasi donasi yang masuk.' },
  { step: '04', icon: '📊', title: 'Penyaluran Dana', desc: 'Setelah campaign selesai, dana disalurkan ke penggalang disertai laporan.' },
])
</script>

<template>
  <div class="home-page">
    <!-- ═══════════════════════════════════════════ -->
    <!-- HERO SECTION — FULL SCREEN                 -->
    <!-- ═══════════════════════════════════════════ -->
    <section class="hero" id="hero-section">
      <div class="hero-overlay"></div>
      <div class="container hero-content">
        <span class="hero-badge">🇮🇩 Platform Donasi Transparan Indonesia</span>
        <h1 class="hero-title">
          Mari Ringankan<br />
          Beban Korban<br />
          <span class="hero-highlight">Bencana</span>
        </h1>
        <p class="hero-desc">
          Setiap rupiah yang Anda donasikan terverifikasi melalui sistem escrow. Transparansi penuh dari donasi hingga penyaluran dana.
        </p>
        <div class="hero-actions">
          <router-link to="/campaigns" class="btn btn-hero-primary" id="hero-cta-explore">
            Jelajahi Campaign
          </router-link>
          <router-link to="/register" class="btn btn-hero-outline" id="hero-cta-register">
            Mulai Galang Dana →
          </router-link>
        </div>
        <div class="hero-stats-row">
          <div class="hs-item">
            <span class="hs-value">{{ platformStats.activeCampaigns }}</span>
            <span class="hs-label">Campaign Aktif</span>
          </div>
          <div class="hs-item">
            <span class="hs-value">{{ formatCurrency(platformStats.totalFundsCollected) }}</span>
            <span class="hs-label">Dana Terkumpul</span>
          </div>
          <div class="hs-item">
            <span class="hs-value">{{ platformStats.totalDonors }}</span>
            <span class="hs-label">Donatur</span>
          </div>
          <div class="hs-item">
            <span class="hs-value">{{ platformStats.completedCampaigns }}</span>
            <span class="hs-label">Campaign Selesai</span>
          </div>
        </div>
      </div>
      <div class="hero-scroll-hint">
        <span>Scroll ke bawah</span>
        <div class="scroll-arrow"></div>
      </div>
    </section>

    <!-- ═══ CAMPAIGN TERBARU ═══ -->
    <section class="campaigns-section" id="campaigns-section">
      <div class="container">
        <div class="section-header">
          <div>
            <h2>Campaign <span class="text-gradient">Terbaru</span></h2>
            <p>Bantu mereka yang membutuhkan dengan donasi Anda</p>
          </div>
          <router-link to="/campaigns" class="btn-see-all" id="btn-see-all-campaigns">
            Lihat Semua →
          </router-link>
        </div>
        <div class="campaigns-grid">
          <CampaignCard v-for="campaign in campaigns" :key="campaign.id" :campaign="campaign" />
        </div>
      </div>
    </section>

    <!-- ═══ HOW IT WORKS ═══ -->
    <section class="how-section" id="how-section">
      <div class="container">
        <div class="section-header center">
          <h2>Bagaimana <span class="text-gradient">OpenDonasi</span> Bekerja?</h2>
          <p>Alur yang transparan dari donasi hingga penyaluran</p>
        </div>
        <div class="steps-grid">
          <div v-for="item in howItWorks" :key="item.step" class="step-card">
            <span class="step-number">{{ item.step }}</span>
            <span class="step-icon">{{ item.icon }}</span>
            <h3>{{ item.title }}</h3>
            <p>{{ item.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ═══ CTA ═══ -->
    <section class="cta-section" id="cta-section">
      <div class="container cta-inner">
        <h2>Siap Membantu Sesama?</h2>
        <p>Bergabunglah bersama ribuan donatur dan penggalang dana di seluruh Indonesia.</p>
        <div class="cta-actions">
          <router-link to="/register" class="btn btn-white-lg">Daftar Sekarang — Gratis</router-link>
          <router-link to="/campaigns" class="btn btn-ghost-lg">Jelajahi Campaign</router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
/* ═══════════════════════════════════════════════ */
/* HERO SECTION — FULL SCREEN                     */
/* ═══════════════════════════════════════════════ */
.hero {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  overflow: hidden;
  background-image: url('https://images.unsplash.com/photo-1469571486292-0ba58a3f068b?q=80&w=1920&auto=format&fit=crop');
  background-size: cover;
  background-position: center;
  background-attachment: fixed;
}

.hero-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    rgba(0, 0, 0, 0.75) 0%,
    rgba(0, 0, 0, 0.55) 40%,
    rgba(0, 77, 128, 0.45) 100%
  );
  z-index: 1;
}

.hero-content {
  position: relative;
  z-index: 2;
  padding-top: 120px;
  padding-bottom: 80px;
  max-width: 720px;
}

.hero-badge {
  display: inline-block;
  background: rgba(255, 255, 255, 0.12);
  color: rgba(255, 255, 255, 0.9);
  font-weight: 600;
  font-size: 13px;
  padding: 6px 18px;
  border-radius: var(--radius-full);
  margin-bottom: 28px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(6px);
  animation: fadeInUp 0.8s ease both;
}

.hero-title {
  font-family: var(--font-heading);
  font-size: 60px;
  font-weight: 800;
  line-height: 1.08;
  color: white;
  margin-bottom: 24px;
  letter-spacing: -1px;
  animation: fadeInUp 0.8s ease 0.15s both;
}

.hero-highlight {
  color: var(--primary);
  position: relative;
}

.hero-highlight::after {
  content: '';
  position: absolute;
  bottom: 4px;
  left: 0;
  right: 0;
  height: 6px;
  background: var(--primary);
  opacity: 0.3;
  border-radius: 4px;
}

.hero-desc {
  font-size: 18px;
  line-height: 1.7;
  color: rgba(255, 255, 255, 0.75);
  max-width: 520px;
  margin-bottom: 36px;
  animation: fadeInUp 0.8s ease 0.3s both;
}

.hero-actions {
  display: flex;
  gap: 14px;
  margin-bottom: 60px;
  animation: fadeInUp 0.8s ease 0.45s both;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-body);
  font-weight: 600;
  border-radius: var(--radius-md);
  text-decoration: none;
  cursor: pointer;
  transition: all 0.25s ease;
  border: none;
}

.btn-hero-primary {
  padding: 16px 36px;
  font-size: 16px;
  background: var(--primary);
  color: white;
}

.btn-hero-primary:hover {
  background: var(--primary-hover);
  box-shadow: 0 8px 32px rgba(0, 191, 255, 0.35);
  transform: translateY(-2px);
  color: white;
}

.btn-hero-outline {
  padding: 16px 36px;
  font-size: 16px;
  background: transparent;
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.btn-hero-outline:hover {
  border-color: white;
  background: rgba(255, 255, 255, 0.08);
  color: white;
}

/* Hero Stats Row */
.hero-stats-row {
  display: flex;
  gap: 40px;
  animation: fadeInUp 0.8s ease 0.6s both;
}

.hs-item {
  display: flex;
  flex-direction: column;
  position: relative;
}

.hs-value {
  font-family: var(--font-heading);
  font-weight: 800;
  font-size: 28px;
  color: white;
}

.hs-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  margin-top: 2px;
}

/* Scroll Hint */
.hero-scroll-hint {
  position: absolute;
  bottom: 32px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 2;
  text-align: center;
  color: rgba(255, 255, 255, 0.4);
  font-size: 12px;
  animation: bounce 2s ease infinite;
}

.scroll-arrow {
  width: 20px;
  height: 20px;
  border-right: 2px solid rgba(255, 255, 255, 0.4);
  border-bottom: 2px solid rgba(255, 255, 255, 0.4);
  transform: rotate(45deg);
  margin: 8px auto 0;
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes bounce {
  0%, 100% { transform: translateX(-50%) translateY(0); }
  50% { transform: translateX(-50%) translateY(8px); }
}

/* ═══ TEXT GRADIENT ═══ */
.text-gradient {
  background: linear-gradient(135deg, var(--primary), #0077B6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* ═══ CAMPAIGNS SECTION ═══ */
.campaigns-section {
  padding: 80px 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 40px;
}

.section-header.center {
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.section-header h2 {
  font-family: var(--font-heading);
  font-size: 32px;
  font-weight: 800;
  margin-bottom: 8px;
}

.section-header p {
  font-size: 16px;
  color: var(--text-muted);
}

.btn-see-all {
  font-weight: 600;
  font-size: 14px;
  color: var(--primary);
  text-decoration: none;
  white-space: nowrap;
  transition: all 0.2s ease;
}

.btn-see-all:hover { color: var(--primary-hover); }

.campaigns-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}

/* ═══ HOW IT WORKS ═══ */
.how-section {
  padding: 80px 0;
  background: linear-gradient(180deg, var(--bg-body), #E8F4FD);
}

.steps-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  margin-top: 48px;
}

.step-card {
  background: white;
  padding: 32px 24px;
  border-radius: var(--radius-lg);
  text-align: center;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  position: relative;
  transition: all 0.25s ease;
}

.step-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

.step-number {
  position: absolute;
  top: -14px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--primary);
  color: white;
  font-family: var(--font-heading);
  font-weight: 800;
  font-size: 12px;
  width: 28px;
  height: 28px;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
}

.step-icon { display: block; font-size: 40px; margin-bottom: 16px; margin-top: 8px; }
.step-card h3 { font-family: var(--font-heading); font-size: 16px; font-weight: 700; margin-bottom: 8px; }
.step-card p { font-size: 14px; color: var(--text-muted); line-height: 1.6; }

/* ═══ CTA ═══ */
.cta-section { padding: 80px 0; }

.cta-inner {
  background: linear-gradient(135deg, var(--primary), #0077B6);
  border-radius: var(--radius-xl);
  padding: 64px 48px;
  text-align: center;
  position: relative;
  overflow: hidden;
}

.cta-inner::before {
  content: '';
  position: absolute;
  width: 400px;
  height: 400px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
  top: -200px;
  right: -100px;
}

.cta-inner h2 { font-family: var(--font-heading); font-size: 36px; font-weight: 800; color: white; margin-bottom: 12px; }
.cta-inner p { font-size: 17px; color: rgba(255, 255, 255, 0.85); margin-bottom: 32px; }

.cta-actions { display: flex; gap: 14px; justify-content: center; }

.btn-white-lg {
  padding: 14px 32px;
  font-size: 15px;
  background: white;
  color: var(--primary);
  font-weight: 700;
}
.btn-white-lg:hover { background: #F0F9FF; transform: translateY(-2px); box-shadow: 0 8px 24px rgba(0,0,0,0.15); color: var(--primary); }

.btn-ghost-lg {
  padding: 14px 32px;
  font-size: 15px;
  background: transparent;
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.35);
}
.btn-ghost-lg:hover { border-color: white; background: rgba(255,255,255,0.1); color: white; }

/* ═══ RESPONSIVE ═══ */
@media (max-width: 1024px) {
  .steps-grid { grid-template-columns: repeat(2, 1fr); }
  .hero-title { font-size: 44px; }
}

@media (max-width: 768px) {
  .hero-title { font-size: 36px; }
  .hero-desc { font-size: 16px; }
  .hero-actions { flex-direction: column; }
  .hero-stats-row { flex-wrap: wrap; gap: 20px; }
  .steps-grid { grid-template-columns: 1fr; }
  .section-header { flex-direction: column; align-items: flex-start; gap: 12px; }
  .cta-inner { padding: 40px 24px; }
  .cta-inner h2 { font-size: 28px; }
  .cta-actions { flex-direction: column; }
}
</style>
