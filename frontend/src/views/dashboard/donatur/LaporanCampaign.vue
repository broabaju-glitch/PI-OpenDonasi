<script setup>
import { ref } from 'vue'

const campaignReports = ref([])
</script>

<template>
  <div>
    <h2>Laporan Campaign</h2>
    <p class="page-desc">Lihat laporan penyaluran dana dari campaign yang Anda donasikan.</p>
    
    <div v-if="campaignReports.length === 0" class="empty-state">
      Belum ada laporan campaign
    </div>
    <div v-else class="report-list">
      <div v-for="report in campaignReports" :key="report.id" class="rp-card">
        <div class="rp-header">
          <h3>{{ report.campaign?.title }}</h3>
          <span class="badge badge-success">{{ report.status }}</span>
        </div>
        <div class="rp-body">
          <h4>{{ report.title }}</h4>
          <p>{{ report.description }}</p>
          <div class="rp-meta">
            <span>📅 Tanggal: {{ report.date }}</span>
            <span>💰 Dana Disalurkan: Rp {{ report.amount }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
h2 { font-family: var(--font-heading); font-size: 24px; margin-bottom: 4px; }
.page-desc { color: var(--text-muted); font-size: 14px; margin-bottom: 24px; }
.empty-state { text-align: center; padding: 40px; color: var(--text-muted); background: white; border-radius: var(--radius-md); border: 1px dashed var(--border-light); margin-bottom: 24px; }
.report-list { display: flex; flex-direction: column; gap: 16px; }
.rp-card { background: white; border: 1px solid var(--border-light); border-radius: var(--radius-md); overflow: hidden; transition: all var(--transition-fast); }
.rp-card:hover { box-shadow: var(--shadow-md); }
.rp-header { display: flex; justify-content: space-between; align-items: center; padding: 16px 20px; border-bottom: 1px solid var(--border-light); gap: 12px; }
.rp-header h3 { font-family: var(--font-heading); font-size: 16px; }
.badge { padding: 4px 10px; border-radius: var(--radius-full); font-size: 12px; font-weight: 600; white-space: nowrap; }
.badge-success { background: var(--success-light); color: var(--success); }
.badge-warning { background: var(--warning-light); color: var(--warning); }
.rp-body { padding: 20px; }
.rp-body h4 { font-family: var(--font-heading); font-size: 15px; margin-bottom: 8px; }
.rp-body p { font-size: 14px; line-height: 1.7; color: var(--text-body); }
.rp-meta { display: flex; gap: 20px; margin-top: 16px; font-size: 13px; color: var(--text-muted); }
.empty-text { color: var(--text-muted); font-style: italic; }
</style>
