<script setup>
import { ref, computed, onMounted } from 'vue'
import CampaignCard from '../components/CampaignCard.vue'
import { apiGetCampaigns } from '../services/api'

const searchQuery = ref('')
const selectedCategory = ref('Semua')
const loading = ref(true)
const apiError = ref('')

const categories = ['Semua', 'Banjir', 'Gempa Bumi', 'Gunung Meletus', 'Tsunami', 'Tanah Longsor', 'Rehabilitasi']

const mockCampaigns = [
  { id: 1, title: 'Bantu Korban Banjir Bandang Garut', category: 'Banjir', foto: 'https://images.unsplash.com/photo-1547683905-f686c993aae5?w=600&h=400&fit=crop', status: 'Aktif', dana_terkumpul: 75000000, target_dana: 100000000, lokasi_kejadian: 'Garut, Jawa Barat', penggalang: { name: 'Yayasan Peduli' }, end_date: '2026-07-15T00:00:00Z' },
  { id: 2, title: 'Darurat Gempa Cianjur — Bangun Kembali', category: 'Gempa Bumi', foto: 'https://images.unsplash.com/photo-1599707367812-042c63e2e6f4?w=600&h=400&fit=crop', status: 'Aktif', dana_terkumpul: 250000000, target_dana: 500000000, lokasi_kejadian: 'Cianjur, Jawa Barat', penggalang: { name: 'Tim Relawan Cianjur' }, end_date: '2026-07-21T00:00:00Z' },
  { id: 3, title: 'Bantuan Logistik Erupsi Gunung Semeru', category: 'Gunung Meletus', foto: 'https://images.unsplash.com/photo-1580974852861-c381510bc98a?w=600&h=400&fit=crop', status: 'Target Tercapai', dana_terkumpul: 150000000, target_dana: 150000000, lokasi_kejadian: 'Lumajang, Jawa Timur', penggalang: { name: 'BNPB Lumajang' }, end_date: '2026-06-30T00:00:00Z' },
  { id: 4, title: 'Rehabilitasi Pasca Tsunami Palu', category: 'Tsunami', foto: 'https://images.unsplash.com/photo-1542393545-10f5cde2c810?w=600&h=400&fit=crop', status: 'Aktif', dana_terkumpul: 42000000, target_dana: 200000000, lokasi_kejadian: 'Palu, Sulawesi Tengah', penggalang: { name: 'Komunitas Palu Bangkit' }, end_date: '2026-07-30T00:00:00Z' },
  { id: 5, title: 'Tanggap Darurat Longsor Sumedang', category: 'Tanah Longsor', foto: 'https://images.unsplash.com/photo-1613842127673-3f3e5b45e4c7?w=600&h=400&fit=crop', status: 'Aktif', dana_terkumpul: 18000000, target_dana: 80000000, lokasi_kejadian: 'Sumedang, Jawa Barat', penggalang: { name: 'Relawan Sumedang' }, end_date: '2026-07-07T00:00:00Z' },
  { id: 6, title: 'Bangun Sekolah Darurat Korban Bencana', category: 'Rehabilitasi', foto: 'https://images.unsplash.com/photo-1497633762265-9d179a990aa6?w=600&h=400&fit=crop', status: 'Aktif', dana_terkumpul: 320000000, target_dana: 400000000, lokasi_kejadian: 'NTT, Indonesia', penggalang: { name: 'Yayasan Cerdas Bangsa' }, end_date: '2026-08-15T00:00:00Z' },
  { id: 7, title: 'Evakuasi Warga Banjir Kalimantan', category: 'Banjir', foto: 'https://images.unsplash.com/photo-1559128010-7c1ad6e1b6a5?w=600&h=400&fit=crop', status: 'Aktif', dana_terkumpul: 60000000, target_dana: 120000000, lokasi_kejadian: 'Banjarmasin, Kalimantan Selatan', penggalang: { name: 'Relawan Borneo' }, end_date: '2026-07-18T00:00:00Z' },
  { id: 8, title: 'Pangan untuk Korban Gempa Lombok', category: 'Gempa Bumi', foto: 'https://images.unsplash.com/photo-1488521787991-ed7bbaae773c?w=600&h=400&fit=crop', status: 'Selesai', dana_terkumpul: 90000000, target_dana: 90000000, lokasi_kejadian: 'Lombok, NTB', penggalang: { name: 'ACT Lombok' }, end_date: '2026-06-01T00:00:00Z' },
]

const allCampaigns = ref(mockCampaigns)

onMounted(async () => {
  try {
    const data = await apiGetCampaigns()
    if (data && data.length > 0) {
      allCampaigns.value = data
    }
  } catch (err) {
    apiError.value = 'Menggunakan data demo'
    console.warn('[CampaignList] API unreachable, using mock data:', err.message)
  }
  loading.value = false
})

const filteredCampaigns = computed(() => {
  return allCampaigns.value.filter(c => {
    const matchCategory = selectedCategory.value === 'Semua' || c.category === selectedCategory.value
    const loc = c.lokasi_kejadian || c.lokasi || ''
    const matchSearch = c.title.toLowerCase().includes(searchQuery.value.toLowerCase()) || loc.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchCategory && matchSearch
  })
})
</script>

<template>
  <div class="campaign-list-page">
    <section class="page-header">
      <div class="container">
        <h1>Daftar <span class="text-gradient">Campaign</span></h1>
        <p>Temukan campaign yang membutuhkan bantuan Anda</p>
      </div>
    </section>

    <section class="filters-section">
      <div class="container">
        <div class="filters">
          <div class="search-box" id="search-campaigns">
            <span class="search-icon">🔍</span>
            <input type="text" v-model="searchQuery" placeholder="Cari campaign atau lokasi..." />
          </div>
          <div class="category-tabs">
            <button
              v-for="cat in categories"
              :key="cat"
              class="cat-tab"
              :class="{ active: selectedCategory === cat }"
              @click="selectedCategory = cat"
            >
              {{ cat }}
            </button>
          </div>
        </div>
      </div>
    </section>

    <section class="results-section">
      <div class="container">
        <p class="results-count">Menampilkan {{ filteredCampaigns.length }} campaign</p>
        <div class="campaigns-grid" v-if="filteredCampaigns.length">
          <CampaignCard v-for="c in filteredCampaigns" :key="c.id" :campaign="c" />
        </div>
        <div class="empty-state" v-else>
          <span class="empty-icon">📭</span>
          <h3>Tidak Ada Campaign Ditemukan</h3>
          <p>Coba ubah kata kunci atau filter kategori Anda.</p>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.page-header {
  padding: 48px 0 32px;
  text-align: center;
  background: linear-gradient(135deg, #E0F7FF, #F0F9FF);
}

.page-header h1 {
  font-family: var(--font-heading);
  font-size: 32px;
  font-weight: 800;
  margin-bottom: 8px;
}

.text-gradient {
  background: linear-gradient(135deg, var(--primary), #0077B6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-header p { color: var(--text-muted); font-size: 16px; }

.filters-section {
  padding: 24px 0;
  background: white;
  border-bottom: 1px solid var(--border-light);
  position: sticky;
  top: 72px;
  z-index: 10;
}

.filters {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.search-box {
  position: relative;
  max-width: 480px;
}

.search-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 16px;
}

.search-box input {
  width: 100%;
  padding: 12px 16px 12px 44px;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 14px;
  background: var(--bg-body);
  transition: all var(--transition-fast);
}

.search-box input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-light);
  background: white;
}

.category-tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.cat-tab {
  padding: 8px 18px;
  border: 1px solid var(--border);
  border-radius: var(--radius-full);
  font-size: 13px;
  font-weight: 500;
  color: var(--text-body);
  background: white;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.cat-tab:hover { border-color: var(--primary); color: var(--primary); }
.cat-tab.active { background: var(--primary); color: white; border-color: var(--primary); }

.results-section { padding: 32px 0 80px; }
.results-count { font-size: 14px; color: var(--text-muted); margin-bottom: 24px; }

.campaigns-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}

.empty-state {
  text-align: center;
  padding: 64px 24px;
}

.empty-icon { font-size: 48px; display: block; margin-bottom: 16px; }
.empty-state h3 { font-size: 20px; margin-bottom: 8px; }
.empty-state p { color: var(--text-muted); }

/* Removed conflicting media queries for grid */
</style>
