<template>
  <div class="admin-layout">

    <!-- ═══════════════════════════════════════════
         SIDEBAR
    ═══════════════════════════════════════════ -->
    <aside class="sidebar" :class="{ 'sidebar--collapsed': sidebarCollapsed }">

      <!-- Logo Area -->
      <div class="sidebar-brand">
        <div class="brand-logo" aria-hidden="true">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
            <path d="M12 21.593c-5.63-5.539-11-10.297-11-14.402 0-3.791 3.068-5.191 5.281-5.191 1.312 0 4.151.501 5.719 4.457 1.59-3.968 4.464-4.447 5.726-4.447 2.54 0 5.274 1.621 5.274 5.181 0 4.069-5.136 8.625-11 14.402z"
              fill="#60a5fa"/>
          </svg>
        </div>
        <div class="brand-text" v-show="!sidebarCollapsed">
          <span class="brand-name">OpenDonasi</span>
          <span class="brand-role">Admin Panel</span>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="sidebar-nav" aria-label="Admin navigation">
        <div class="nav-section-label" v-show="!sidebarCollapsed">MENU UTAMA</div>

        <router-link
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="nav-link"
          active-class="nav-link--active"
          :title="sidebarCollapsed ? item.label : undefined"
        >
          <span class="nav-link-icon" aria-hidden="true" v-html="item.icon"></span>
          <span class="nav-link-label" v-show="!sidebarCollapsed">{{ item.label }}</span>
          <span v-if="item.badge && !sidebarCollapsed" class="nav-badge">{{ item.badge }}</span>
        </router-link>
      </nav>

      <!-- Bottom Actions -->
      <div class="sidebar-footer">
        <div class="admin-identity" v-show="!sidebarCollapsed">
          <div class="admin-avatar" aria-hidden="true">
            {{ userInitial }}
          </div>
          <div class="admin-info">
            <span class="admin-name">{{ userName }}</span>
            <span class="admin-role-tag">Administrator</span>
          </div>
        </div>

        <button
          id="admin-logout-btn"
          class="logout-btn"
          :class="{ 'logout-btn--icon-only': sidebarCollapsed }"
          :title="sidebarCollapsed ? 'Logout' : undefined"
          @click="handleLogout"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            <polyline points="16 17 21 12 16 7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            <line x1="21" y1="12" x2="9" y2="12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
          <span v-show="!sidebarCollapsed">Logout</span>
        </button>
      </div>
    </aside>

    <!-- ═══════════════════════════════════════════
         MAIN PANEL
    ═══════════════════════════════════════════ -->
    <div class="main-panel">

      <!-- Topbar -->
      <header class="topbar">
        <!-- Left: Collapse toggle + breadcrumb -->
        <div class="topbar-left">
          <button
            class="collapse-btn"
            :aria-label="sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'"
            @click="sidebarCollapsed = !sidebarCollapsed"
          >
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <line x1="3" y1="6"  x2="21" y2="6"  stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
              <line x1="3" y1="12" x2="21" y2="12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
              <line x1="3" y1="18" x2="21" y2="18" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            </svg>
          </button>
          <div class="topbar-breadcrumb">
            <span class="breadcrumb-root">Admin</span>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <polyline points="9 18 15 12 9 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <span class="breadcrumb-current">{{ currentPageLabel }}</span>
          </div>
        </div>

        <!-- Right: Notifications + Profile -->
        <div class="topbar-right">
          <!-- Notification bell (decorative / extendable) -->
          <button class="topbar-icon-btn" title="Notifikasi" aria-label="Notifikasi">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
              <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M13.73 21a2 2 0 0 1-3.46 0" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>

          <!-- Profile chip -->
          <div class="profile-chip">
            <div class="profile-avatar" aria-hidden="true">{{ userInitial }}</div>
            <div class="profile-info">
              <span class="profile-name">{{ userName }}</span>
              <span class="profile-badge">Admin</span>
            </div>
          </div>
        </div>
      </header>

      <!-- Content Area -->
      <main class="content-area" id="admin-main-content">
        <router-view v-slot="{ Component, route }">
          <transition name="page-fade" mode="out-in">
            <component :is="Component" :key="route.path" />
          </transition>
        </router-view>
      </main>
    </div>

  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

// ─── Composables ───
const router = useRouter()
const route  = useRoute()
const { logout, user } = useAuth()

// ─── Sidebar State ───
const sidebarCollapsed = ref(false)

// ─── User Info ───
const userName    = computed(() => user.value?.name  || 'Admin Utama')
const userInitial = computed(() => (user.value?.name || 'A').charAt(0).toUpperCase())

// ─── Navigation Items ───
const navItems = [
  {
    to: '/admin/platform-overview',
    label: 'Dashboard',
    icon: `<svg width="18" height="18" viewBox="0 0 24 24" fill="none">
      <rect x="3" y="3" width="7" height="7" rx="1" stroke="currentColor" stroke-width="1.5"/>
      <rect x="14" y="3" width="7" height="7" rx="1" stroke="currentColor" stroke-width="1.5"/>
      <rect x="3" y="14" width="7" height="7" rx="1" stroke="currentColor" stroke-width="1.5"/>
      <rect x="14" y="14" width="7" height="7" rx="1" stroke="currentColor" stroke-width="1.5"/>
    </svg>`,
  },
  {
    to: '/admin/manajemen-pengguna',
    label: 'Manajemen Pengguna',
    icon: `<svg width="18" height="18" viewBox="0 0 24 24" fill="none">
      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
      <circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="1.5"/>
      <path d="M23 21v-2a4 4 0 0 0-3-3.87" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
      <path d="M16 3.13a4 4 0 0 1 0 7.75" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>`,
  },
  {
    to: '/admin/moderasi-campaign',
    label: 'Moderasi Campaign',
    icon: `<svg width="18" height="18" viewBox="0 0 24 24" fill="none">
      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
      <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
      <line x1="9" y1="13" x2="15" y2="13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      <line x1="9" y1="17" x2="12" y2="17" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
    </svg>`,
  },
  {
    to: '/admin/manajemen-transaksi',
    label: 'Manajemen Transaksi',
    icon: `<svg width="18" height="18" viewBox="0 0 24 24" fill="none">
      <line x1="12" y1="1" x2="12" y2="23" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
      <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>`,
  },
]

// ─── Breadcrumb label from active route ───
const currentPageLabel = computed(() => {
  const matched = navItems.find(item => route.path.startsWith(item.to))
  return matched?.label ?? 'Panel Admin'
})

// ─── Logout Handler ───
async function handleLogout() {
  logout()
  await router.push('/admin/login')
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');

/* ─── Reset & Base ─── */
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

/* ═══════════════════════════════════════
   ROOT LAYOUT
═══════════════════════════════════════ */
.admin-layout {
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  display: flex;
  height: 100vh;
  overflow: hidden;
  background: #f8fafc;
}

/* ═══════════════════════════════════════
   SIDEBAR
═══════════════════════════════════════ */
.sidebar {
  width: 248px;
  min-width: 248px;
  background: #1e293b;
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  transition: width 0.25s ease, min-width 0.25s ease;
  z-index: 100;
  flex-shrink: 0;
}

.sidebar--collapsed {
  width: 68px;
  min-width: 68px;
}

/* ── Brand / Logo ── */
.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 16px 18px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  flex-shrink: 0;
}

.brand-logo {
  width: 36px;
  height: 36px;
  background: rgba(96, 165, 250, 0.12);
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid rgba(96, 165, 250, 0.2);
}

.brand-text {
  display: flex;
  flex-direction: column;
  gap: 1px;
  overflow: hidden;
}

.brand-name {
  font-size: 0.9rem;
  font-weight: 700;
  color: #f1f5f9;
  white-space: nowrap;
  letter-spacing: -0.2px;
}

.brand-role {
  font-size: 0.68rem;
  font-weight: 500;
  color: #64748b;
  white-space: nowrap;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

/* ── Navigation ── */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 12px 10px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  scrollbar-width: thin;
  scrollbar-color: #334155 transparent;
}

.sidebar-nav::-webkit-scrollbar {
  width: 4px;
}
.sidebar-nav::-webkit-scrollbar-track { background: transparent; }
.sidebar-nav::-webkit-scrollbar-thumb { background: #334155; border-radius: 4px; }

.nav-section-label {
  font-size: 0.62rem;
  font-weight: 700;
  color: #475569;
  letter-spacing: 0.1em;
  padding: 8px 8px 4px;
  white-space: nowrap;
  overflow: hidden;
}

/* ── Nav Link ── */
.nav-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 10px;
  border-radius: 8px;
  color: #94a3b8;
  text-decoration: none;
  font-size: 0.825rem;
  font-weight: 500;
  transition: background 0.15s ease, color 0.15s ease, border-color 0.15s ease;
  border-left: 3px solid transparent;
  white-space: nowrap;
  overflow: hidden;
  position: relative;
}

.nav-link:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #e2e8f0;
}

.nav-link--active {
  background: #334155;
  color: #f1f5f9;
  border-left-color: #3b82f6;
}

.nav-link-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 20px;
  height: 20px;
}

.nav-link-label {
  overflow: hidden;
  text-overflow: ellipsis;
}

.nav-badge {
  margin-left: auto;
  background: #ef4444;
  color: #fff;
  font-size: 0.65rem;
  font-weight: 700;
  padding: 1px 6px;
  border-radius: 20px;
  flex-shrink: 0;
}

/* ── Sidebar Footer ── */
.sidebar-footer {
  padding: 12px 10px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex-shrink: 0;
}

.admin-identity {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 4px;
  overflow: hidden;
}

.admin-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #3b82f6, #6366f1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.admin-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
  overflow: hidden;
}

.admin-name {
  font-size: 0.8rem;
  font-weight: 600;
  color: #e2e8f0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.admin-role-tag {
  font-size: 0.65rem;
  color: #64748b;
  white-space: nowrap;
}

/* ── Logout Button ── */
.logout-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 9px 12px;
  background: rgba(239, 68, 68, 0.1);
  color: #f87171;
  border: 1px solid rgba(239, 68, 68, 0.15);
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.825rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease, border-color 0.15s ease;
  white-space: nowrap;
  overflow: hidden;
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.18);
  color: #fca5a5;
  border-color: rgba(239, 68, 68, 0.28);
}

.logout-btn--icon-only {
  padding: 9px;
}

/* ═══════════════════════════════════════
   MAIN PANEL
═══════════════════════════════════════ */
.main-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

/* ── Topbar ── */
.topbar {
  height: 60px;
  background: #ffffff;
  border-bottom: 1px solid #e9edf2;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  gap: 16px;
  flex-shrink: 0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  z-index: 50;
}

.topbar-left {
  display: flex;
  align-items: center;
  gap: 14px;
}

/* ── Collapse / Hamburger Button ── */
.collapse-btn {
  width: 36px;
  height: 36px;
  background: none;
  border: none;
  border-radius: 8px;
  color: #64748b;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s ease, color 0.15s ease;
  flex-shrink: 0;
}

.collapse-btn:hover {
  background: #f1f5f9;
  color: #334155;
}

/* ── Breadcrumb ── */
.topbar-breadcrumb {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.825rem;
}

.breadcrumb-root {
  color: #94a3b8;
  font-weight: 500;
}

.breadcrumb-current {
  color: #1e293b;
  font-weight: 600;
}

/* ── Topbar Right ── */
.topbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

.topbar-icon-btn {
  width: 36px;
  height: 36px;
  background: none;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  color: #64748b;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s ease, color 0.15s ease;
}

.topbar-icon-btn:hover {
  background: #f1f5f9;
  color: #334155;
}

/* ── Profile Chip ── */
.profile-chip {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 5px 12px 5px 6px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 40px;
  cursor: default;
}

.profile-avatar {
  width: 30px;
  height: 30px;
  background: linear-gradient(135deg, #3b82f6, #6366f1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.profile-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.profile-name {
  font-size: 0.78rem;
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
  line-height: 1;
}

.profile-badge {
  font-size: 0.62rem;
  font-weight: 600;
  color: #3b82f6;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  line-height: 1;
}

/* ═══════════════════════════════════════
   CONTENT AREA
═══════════════════════════════════════ */
.content-area {
  flex: 1;
  overflow-y: auto;
  background: #f8fafc;
  padding: 24px;
  scrollbar-width: thin;
  scrollbar-color: #cbd5e1 transparent;
}

.content-area::-webkit-scrollbar       { width: 6px; }
.content-area::-webkit-scrollbar-track { background: transparent; }
.content-area::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 6px; }

/* ── Page Transition ── */
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.18s ease, transform 0.18s ease;
}

.page-fade-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.page-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* ═══════════════════════════════════════
   RESPONSIVE (Mobile)
═══════════════════════════════════════ */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    z-index: 200;
    transform: translateX(0);
  }

  .sidebar--collapsed {
    transform: translateX(-100%);
    width: 248px;
    min-width: 248px;
  }

  .main-panel {
    width: 100%;
  }

  .topbar {
    padding: 0 16px;
  }

  .content-area {
    padding: 16px;
  }

  .profile-info {
    display: none;
  }
}
</style>
