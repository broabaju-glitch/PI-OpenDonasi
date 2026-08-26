<script setup>
import { useRoute } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const route = useRoute()
const { user } = useAuth()

const menuItems = [
  { label: 'Campaign Saya', icon: '📋', to: '/penggalang/campaign-saya' },
  { label: 'Buat Campaign', icon: '➕', to: '/penggalang/buat-campaign' },
  { label: 'Donasi Masuk', icon: '💳', to: '/penggalang/donasi-masuk' },
  { label: 'Status Dana', icon: '💰', to: '/penggalang/status-dana' },
  { label: 'Upload Laporan', icon: '📄', to: '/penggalang/upload-laporan' },
]
</script>

<template>
  <div class="dashboard-page">
    <div class="container dashboard-layout">
      <aside class="sidebar">
        <div class="sidebar-header">
          <div class="sidebar-avatar">{{ user?.name ? user.name.charAt(0).toUpperCase() : '📢' }}</div>
          <div>
            <h4>Penggalang Dana</h4>
            <span class="sidebar-role">{{ user?.name || 'Loading...' }}</span>
          </div>
        </div>
        <nav class="sidebar-nav">
          <router-link v-for="item in menuItems" :key="item.to" :to="item.to" class="sidebar-link" :class="{ active: route.path === item.to }">
            <span class="sl-icon">{{ item.icon }}</span>
            {{ item.label }}
          </router-link>
        </nav>
      </aside>
      <main class="dashboard-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<style scoped>
.dashboard-page { padding: 32px 0 80px; padding-top: 104px; }
.dashboard-layout { display: grid; grid-template-columns: 260px 1fr; gap: 28px; align-items: flex-start; }
.sidebar { background: white; border-radius: var(--radius-lg); box-shadow: var(--shadow-sm); border: 1px solid var(--border-light); padding: 20px; position: sticky; top: 96px; }
.sidebar-header { display: flex; align-items: center; gap: 12px; padding-bottom: 16px; border-bottom: 1px solid var(--border-light); margin-bottom: 12px; }
.sidebar-avatar { width: 40px; height: 40px; border-radius: var(--radius-full); background: var(--primary-light); display: flex; align-items: center; justify-content: center; font-size: 20px; }
.sidebar-header h4 { font-family: var(--font-heading); font-size: 15px; font-weight: 700; }
.sidebar-role { font-size: 12px; color: var(--text-muted); }
.sidebar-nav { display: flex; flex-direction: column; gap: 4px; }
.sidebar-link { display: flex; align-items: center; gap: 10px; padding: 10px 14px; border-radius: var(--radius-sm); font-size: 14px; font-weight: 500; color: var(--text-body); text-decoration: none; transition: all var(--transition-fast); }
.sidebar-link:hover { background: var(--bg-body); color: var(--text-heading); }
.sidebar-link.active { background: var(--primary-light); color: var(--primary); font-weight: 600; }
.sl-icon { font-size: 18px; }
@media (max-width: 768px) { .dashboard-layout { grid-template-columns: 1fr; } .sidebar { position: static; } }
</style>
