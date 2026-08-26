<script setup>
import { ref, computed, onMounted } from 'vue'
import { apiGetMyDonations } from '@/services/api'

const donationHistory = ref([])

const fetchDonationHistory = async () => {
  try {
    const response = await apiGetMyDonations()
    console.log("Fetched donations:", response)
    donationHistory.value = response?.data?.data || response?.data || response || []
  } catch (error) {
    console.error("Failed to fetch donation history:", error)
  }
}

onMounted(() => {
  fetchDonationHistory()
})

const totalDonasi = computed(() => {
  if (donationHistory.value.length === 0) return 0
  return donationHistory.value.reduce((acc, curr) => acc + curr.amount, 0)
})

const formatCurrency = (val) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(val)
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
}
</script>

<template>
  <div>
    <h2>Riwayat Donasi</h2>
    <p class="page-desc">Semua donasi yang pernah Anda berikan.</p>
    
    <div v-if="donationHistory.length === 0" class="empty-state">
      Belum ada riwayat donasi saat ini
    </div>
    <div v-else class="donation-list">
      <div v-for="donation in donationHistory" :key="donation.ID || donation.id" class="dl-card">
        <div class="dl-left">
          <h3>{{ donation.campaign?.title }}</h3>
          <p class="dl-meta">{{ formatDate(donation.CreatedAt || donation.created_at) }} · oleh {{ donation.campaign?.penggalang?.name || donation.campaign?.penggalang || '-' }}</p>
        </div>
        <div class="dl-right">
          <span class="dl-amount">{{ formatCurrency(donation.amount) }}</span>
          <span class="badge badge-success">{{ donation.status }}</span>
        </div>
      </div>
    </div>
    
    <div class="total-box">
      <span>Total Donasi Anda</span>
      <strong>{{ formatCurrency(totalDonasi) }}</strong>
    </div>
  </div>
</template>

<style scoped>
h2 { font-family: var(--font-heading); font-size: 24px; margin-bottom: 4px; }
.page-desc { color: var(--text-muted); font-size: 14px; margin-bottom: 24px; }
.empty-state { text-align: center; padding: 40px; color: var(--text-muted); background: white; border-radius: var(--radius-md); border: 1px dashed var(--border-light); margin-bottom: 24px; }
.donation-list { display: flex; flex-direction: column; gap: 12px; margin-bottom: 24px; }
.dl-card { display: flex; justify-content: space-between; align-items: center; background: white; border: 1px solid var(--border-light); border-radius: var(--radius-md); padding: 20px; transition: all var(--transition-fast); }
.dl-card:hover { box-shadow: var(--shadow-md); }
.dl-left h3 { font-family: var(--font-heading); font-size: 16px; margin-bottom: 4px; }
.dl-meta { font-size: 13px; color: var(--text-muted); }
.dl-right { text-align: right; display: flex; flex-direction: column; align-items: flex-end; gap: 6px; }
.dl-amount { font-family: var(--font-heading); font-weight: 700; font-size: 18px; color: var(--primary); }
.badge { padding: 4px 10px; border-radius: var(--radius-full); font-size: 12px; font-weight: 600; }
.badge-success { background: var(--success-light); color: var(--success); }
.badge-warning { background: var(--warning-light); color: var(--warning); }
.total-box { background: var(--primary-light); border: 1px solid rgba(0,191,255,0.15); border-radius: var(--radius-md); padding: 20px; display: flex; justify-content: space-between; align-items: center; }
.total-box span { font-size: 15px; color: var(--text-body); }
.total-box strong { font-family: var(--font-heading); font-size: 22px; font-weight: 800; color: var(--primary); }
@media (max-width: 640px) { .dl-card { flex-direction: column; align-items: flex-start; gap: 12px; } .dl-right { align-items: flex-start; } }
</style>
