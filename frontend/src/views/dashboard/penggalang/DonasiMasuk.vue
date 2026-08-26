<script setup>
import { ref, onMounted } from 'vue'
import { apiGetFundraiserDonations } from '@/services/api'

const incomingDonations = ref([])

onMounted(async () => {
  try {
    const data = await apiGetFundraiserDonations()
    incomingDonations.value = data.map(d => ({
      id: d.ID,
      donorName: d.is_anonymous ? 'Hamba Allah' : (d.donor_name || d.donatur?.name || 'Hamba Allah'),
      campaignTitle: d.campaign?.title || 'Unknown Campaign',
      amount: new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(d.amount),
      date: new Date(d.CreatedAt).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' }),
      status: d.status
    }))
  } catch (error) {
    console.error("Gagal memuat donasi masuk:", error)
  }
})
</script>

<template>
  <div>
    <h2>Donasi Masuk</h2>
    <p class="page-desc">Daftar donasi yang masuk untuk campaign Anda.</p>
    <div class="table-wrapper">
      <table class="data-table">
        <thead>
          <tr>
            <th>Donatur</th>
            <th>Campaign</th>
            <th>Jumlah</th>
            <th>Tanggal</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="incomingDonations.length === 0">
            <td colspan="5" class="empty-state">Belum ada data saat ini</td>
          </tr>
          <tr v-else v-for="donation in incomingDonations" :key="donation.id">
            <td><strong>{{ donation.donorName }}</strong></td>
            <td>{{ donation.campaignTitle }}</td>
            <td>{{ donation.amount }}</td>
            <td>{{ donation.date }}</td>
            <td><span class="badge badge-success">{{ donation.status }}</span></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
h2 { font-family: var(--font-heading); font-size: 24px; margin-bottom: 4px; }
.page-desc { color: var(--text-muted); font-size: 14px; margin-bottom: 24px; }
.table-wrapper { background: white; border-radius: var(--radius-md); border: 1px solid var(--border-light); overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th { background: var(--bg-body); padding: 12px 16px; text-align: left; font-size: 13px; font-weight: 600; color: var(--text-muted); text-transform: uppercase; }
.data-table td { padding: 14px 16px; border-top: 1px solid var(--border-light); font-size: 14px; color: var(--text-body); }
.data-table td strong { color: var(--text-heading); }
.badge { padding: 4px 10px; border-radius: var(--radius-full); font-size: 12px; font-weight: 600; }
.badge-success { background: var(--success-light); color: var(--success); }
.badge-warning { background: var(--warning-light); color: var(--warning); }
.empty-state { text-align: center; color: var(--text-muted); padding: 32px !important; font-style: italic; }
</style>
