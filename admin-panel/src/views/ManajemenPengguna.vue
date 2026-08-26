<!--
  === GO BACKEND INTEGRATION CHECKLIST ===
  Route: GET /api/users must be implemented.
  Auth: The endpoint must verify the JWT and ensure the role is admin.
  Response Format: The backend must return a JSON array:
  [
    {
      "id": 1,
      "name": "...",
      "email": "...",
      "role": "...",
      "status": "aktif", // or "nonaktif"
      "joined_at": "2024-01-10T00:00:00Z"
    }
  ]
-->
<template>
  <div class="page">

    <!-- ── Page Header ── -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Manajemen Pengguna</h1>
        <p class="page-subtitle">Kelola semua pengguna yang terdaftar di platform OpenDonasi.</p>
      </div>
      <div class="header-actions">
        <div class="search-box" :class="{ 'search-box--focused': searchFocused }">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <circle cx="11" cy="11" r="8" stroke="currentColor" stroke-width="1.5"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
          <input
            v-model="searchQuery"
            type="search"
            id="user-search"
            placeholder="Cari nama atau email..."
            class="search-input"
            @focus="searchFocused = true"
            @blur="searchFocused = false"
            aria-label="Cari pengguna"
          />
        </div>
        <button class="btn-primary" @click="openAddModal">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <line x1="12" y1="5" x2="12" y2="19" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="5" y1="12" x2="19" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          Tambah Pengguna
        </button>
      </div>
    </div>

    <!-- ── Stats Row ── -->
    <div class="stats-row">
      <div class="mini-stat" v-for="s in miniStats" :key="s.label" :style="{ '--c': s.color, '--cs': s.colorSoft }">
        <span class="mini-stat-dot"></span>
        <span class="mini-stat-value">{{ s.value }}</span>
        <span class="mini-stat-label">{{ s.label }}</span>
      </div>
    </div>

    <!-- ── Table Card ── -->
    <div class="card">

      <!-- Toolbar -->
      <div class="card-toolbar">
        <div class="toolbar-left">
          <span class="result-count" v-if="!isLoading">
            Menampilkan <strong>{{ filteredUsers.length }}</strong> dari {{ users.length }} pengguna
          </span>
        </div>
        <div class="toolbar-right">
          <select v-model="roleFilter" class="filter-select" aria-label="Filter role">
            <option value="">Semua Role</option>
            <option value="admin">Admin</option>
            <option value="penggalang">Penggalang</option>
            <option value="donatur">Donatur</option>
          </select>
          <select v-model="statusFilter" class="filter-select" aria-label="Filter status">
            <option value="">Semua Status</option>
            <option value="aktif">Aktif</option>
            <option value="nonaktif">Nonaktif</option>
          </select>
        </div>
      </div>

      <!-- Loading skeleton -->
      <div v-if="isLoading" class="skeleton-table">
        <div v-for="i in 6" :key="i" class="skeleton-row">
          <div class="skel skel--sm"></div>
          <div class="skel skel--md"></div>
          <div class="skel skel--lg"></div>
          <div class="skel skel--sm"></div>
          <div class="skel skel--sm"></div>
          <div class="skel skel--md"></div>
        </div>
      </div>

      <!-- Table -->
      <div v-else-if="filteredUsers.length > 0" class="table-wrapper">
        <table class="data-table" aria-label="Tabel manajemen pengguna">
          <thead>
            <tr>
              <th class="th-id">ID</th>
              <th>Nama</th>
              <th>Email</th>
              <th>Role</th>
              <th class="th-center">Status</th>
              <th class="th-center">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="user in filteredUsers"
              :key="user.id"
              class="table-row"
              :class="{ 'row--inactive': user.status === 'nonaktif' }"
            >
              <td class="td-id">#{{ String(user.id || user.ID).padStart(3, '0') }}</td>
              <td class="td-name">
                <div class="user-chip">
                  <div class="avatar" :style="{ background: avatarColor(user.name || user.Name) }" aria-hidden="true">
                    {{ (user.name || user.Name).charAt(0).toUpperCase() }}
                  </div>
                  <div>
                    <span class="user-name">{{ user.name || user.Name }}</span>
                    <span class="user-joined">Bergabung {{ formatDate(user.joined_at || user.CreatedAt) }}</span>
                  </div>
                </div>
              </td>
              <td class="td-email">{{ user.email || user.Email }}</td>
              <td class="td-role">
                <span class="badge-role" :class="`badge-role--${user.role || user.Role}`">{{ capitalize(user.role || user.Role) }}</span>
              </td>
              <td class="td-status th-center">
                <button
                  class="status-toggle"
                  :class="(user.status || 'aktif') === 'aktif' ? 'status-toggle--on' : 'status-toggle--off'"
                  :title="(user.status || 'aktif') === 'aktif' ? 'Klik untuk nonaktifkan' : 'Klik untuk aktifkan'"
                  @click="toggleStatus(user)"
                  :aria-label="`Toggle status ${user.name || user.Name}`"
                >
                  <span class="toggle-knob"></span>
                </button>
              </td>
              <td class="td-actions th-center">
                <div class="action-group">
                  <button class="btn-action btn-edit" @click="openEditModal(user)" title="Edit pengguna">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                    Edit
                  </button>
                  <button class="btn-action btn-delete" @click="confirmDelete(user)" title="Hapus pengguna">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                      <polyline points="3 6 5 6 21 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                      <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                      <path d="M10 11v6M14 11v6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                      <path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                    Hapus
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Error State -->
      <div v-else-if="errorMessage" class="empty-state">
        <div class="empty-icon-wrap" aria-hidden="true" style="color: #ef4444; background: #fef2f2;">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5"/>
            <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <circle cx="12" cy="16" r="1" fill="currentColor"/>
          </svg>
        </div>
        <p class="empty-title">Gagal memuat pengguna</p>
        <p class="empty-sub">{{ errorMessage }}</p>
        <button class="btn-ghost" @click="fetchUsers">Coba Lagi</button>
      </div>

      <!-- Empty State -->
      <div v-else class="empty-state">
        <div class="empty-icon-wrap" aria-hidden="true">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            <circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="1.5"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </div>
        <p class="empty-title">Tidak ada pengguna ditemukan</p>
        <p class="empty-sub">Coba ubah filter atau kata kunci pencarian.</p>
        <button class="btn-ghost" @click="resetFilters">Reset Filter</button>
      </div>

    </div>

    <!-- ══════════════════════════════════
         DELETE CONFIRMATION MODAL
    ══════════════════════════════════ -->
    <transition name="modal">
      <div v-if="showDeleteModal" class="modal-overlay" role="dialog" aria-modal="true" aria-labelledby="delete-modal-title" @click.self="showDeleteModal = false">
        <div class="modal-card">
          <div class="modal-icon modal-icon--danger" aria-hidden="true">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
              <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <line x1="12" y1="9" x2="12" y2="13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
              <circle cx="12" cy="17" r="0.5" fill="currentColor" stroke="currentColor"/>
            </svg>
          </div>
          <h2 id="delete-modal-title" class="modal-title">Hapus Pengguna?</h2>
          <p class="modal-body">
            Kamu akan menghapus pengguna
            <strong>{{ userToDelete?.name || userToDelete?.Name }}</strong>
            (<em>{{ userToDelete?.email || userToDelete?.Email }}</em>).
            Tindakan ini tidak dapat dibatalkan.
          </p>
          <div class="modal-actions">
            <button class="btn-ghost" @click="showDeleteModal = false" :disabled="isDeleting">Batal</button>
            <button class="btn-danger" @click="deleteUser" :disabled="isDeleting">
              <span v-if="!isDeleting">Hapus Sekarang</span>
              <span v-else class="btn-spinner-wrap">
                <span class="btn-spinner" aria-hidden="true"></span>
                Menghapus...
              </span>
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- ══════════════════════════════════
         EDIT / ADD USER MODAL
    ══════════════════════════════════ -->
    <transition name="modal">
      <div v-if="showEditModal" class="modal-overlay" role="dialog" aria-modal="true" :aria-labelledby="editModalMode === 'add' ? 'add-modal-title' : 'edit-modal-title'" @click.self="showEditModal = false">
        <div class="modal-card modal-card--wide">
          <h2 :id="editModalMode === 'add' ? 'add-modal-title' : 'edit-modal-title'" class="modal-title modal-title--left">
            {{ editModalMode === 'add' ? 'Tambah Pengguna' : 'Edit Pengguna' }}
          </h2>

          <form class="edit-form" @submit.prevent="saveUser">
            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="edit-name">Nama Lengkap</label>
                <input id="edit-name" v-model="editForm.name" class="form-input" type="text" placeholder="Masukkan nama" required />
              </div>
              <div class="form-group">
                <label class="form-label" for="edit-email">Email</label>
                <input id="edit-email" v-model="editForm.email" class="form-input" type="email" placeholder="email@contoh.com" required />
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="edit-role">Role</label>
                <select id="edit-role" v-model="editForm.role" class="form-input form-select" required>
                  <option value="donatur">Donatur</option>
                  <option value="penggalang">Penggalang</option>
                  <option value="admin">Admin</option>
                </select>
              </div>
              <div class="form-group">
                <label class="form-label" for="edit-status">Status</label>
                <select id="edit-status" v-model="editForm.status" class="form-input form-select" required>
                  <option value="aktif">Aktif</option>
                  <option value="nonaktif">Nonaktif</option>
                </select>
              </div>
            </div>
            <div class="modal-actions">
              <button type="button" class="btn-ghost" @click="showEditModal = false">Batal</button>
              <button type="submit" class="btn-primary">
                {{ editModalMode === 'add' ? 'Tambah Pengguna' : 'Simpan Perubahan' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- ── Toast Notification ── -->
    <transition name="toast">
      <div v-if="toast.show" class="toast" :class="`toast--${toast.type}`" role="alert" aria-live="polite">
        <svg v-if="toast.type === 'success'" width="16" height="16" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          <polyline points="22 4 12 14.01 9 11.01" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5"/>
          <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          <circle cx="12" cy="16" r="0.5" fill="currentColor" stroke="currentColor"/>
        </svg>
        {{ toast.message }}
      </div>
    </transition>

  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted } from 'vue'
import { adminApi } from '@/services/api'

// ─── State ───────────────────────────────────────────────────────
const users        = ref([])
const isLoading    = ref(true)
const errorMessage = ref('')
const searchQuery  = ref('')
const searchFocused = ref(false)
const roleFilter    = ref('')
const statusFilter  = ref('')

// ─── Computed: Filtered Users ─────────────────────────────────────
const filteredUsers = computed(() => {
  let list = users.value
  if (roleFilter.value) {
    list = list.filter(u => (u.role || u.Role) === roleFilter.value)
  }
  if (statusFilter.value) {
    list = list.filter(u => (u.status || 'aktif') === statusFilter.value)
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    list = list.filter(u =>
      (u.name || u.Name).toLowerCase().includes(q) ||
      (u.email || u.Email).toLowerCase().includes(q)
    )
  }
  return list
})

// ─── Computed: Mini Stats ─────────────────────────────────────────
const miniStats = computed(() => [
  { label: 'Total Pengguna', value: users.value.length,                                         color: '#3b82f6', colorSoft: 'rgba(59,130,246,0.12)' },
  { label: 'Donatur',        value: users.value.filter(u => (u.role || u.Role) === 'donatur').length,        color: '#10b981', colorSoft: 'rgba(16,185,129,0.12)' },
  { label: 'Penggalang',     value: users.value.filter(u => (u.role || u.Role) === 'penggalang').length,     color: '#8b5cf6', colorSoft: 'rgba(139,92,246,0.12)'  },
  { label: 'Admin',          value: users.value.filter(u => (u.role || u.Role) === 'admin').length,          color: '#f59e0b', colorSoft: 'rgba(245,158,11,0.12)'  },
  { label: 'Akun Nonaktif',  value: users.value.filter(u => (u.status || 'aktif') === 'nonaktif').length,     color: '#ef4444', colorSoft: 'rgba(239,68,68,0.12)'   },
])

// ─── Delete Modal ─────────────────────────────────────────────────
const showDeleteModal = ref(false)
const userToDelete    = ref(null)
const isDeleting      = ref(false)

function confirmDelete(user) {
  userToDelete.value  = user
  showDeleteModal.value = true
}

async function deleteUser() {
  if (!userToDelete.value) return
  isDeleting.value = true
  const uid = userToDelete.value.id || userToDelete.value.ID
  const uname = userToDelete.value.name || userToDelete.value.Name
  try {
    await adminApi.deleteUser(uid)
    users.value = users.value.filter(u => (u.id || u.ID) !== uid)
    showDeleteModal.value = false
    showToast(`Pengguna "${uname}" berhasil dihapus.`, 'success')
    userToDelete.value = null
  } catch (err) {
    showToast('Gagal menghapus pengguna. ' + (err.message || 'Coba lagi.'), 'error')
    console.error('[ManajemenPengguna] deleteUser error:', err)
  } finally {
    isDeleting.value = false
  }
}

// ─── Toggle Status ────────────────────────────────────────────────
function toggleStatus(user) {
  const currentStatus = user.status || 'aktif'
  const newStatus = currentStatus === 'aktif' ? 'nonaktif' : 'aktif'
  user.status = newStatus
  const uname = user.name || user.Name
  const msg = newStatus === 'aktif'
    ? `${uname} diaktifkan kembali.`
    : `${uname} dinonaktifkan.`
  showToast(msg, 'success')
}

// ─── Edit / Add Modal ─────────────────────────────────────────────
const showEditModal = ref(false)
const editModalMode = ref('edit') // 'edit' | 'add'
const editingUserId = ref(null)
const editForm      = reactive({ name: '', email: '', role: 'donatur', status: 'aktif' })

function openEditModal(user) {
  editModalMode.value  = 'edit'
  editingUserId.value  = user.id || user.ID
  editForm.name        = user.name || user.Name
  editForm.email       = user.email || user.Email
  editForm.role        = user.role || user.Role
  editForm.status      = user.status || 'aktif'
  showEditModal.value  = true
}

function openAddModal() {
  editModalMode.value  = 'add'
  editingUserId.value  = null
  editForm.name        = ''
  editForm.email       = ''
  editForm.role        = 'donatur'
  editForm.status      = 'aktif'
  showEditModal.value  = true
}

function saveUser() {
  if (editModalMode.value === 'edit') {
    const idx = users.value.findIndex(u => (u.id || u.ID) === editingUserId.value)
    if (idx !== -1) {
      // Just visually update
      if(users.value[idx].Name !== undefined) {
          users.value[idx].Name = editForm.name
          users.value[idx].Email = editForm.email
          users.value[idx].Role = editForm.role
          users.value[idx].status = editForm.status
      } else {
          users.value[idx] = { ...users.value[idx], ...editForm }
      }
    }
    showToast(`Data ${editForm.name} berhasil diperbarui (visual saja).`, 'success')
  } else {
    const maxId = Math.max(...users.value.map(u => u.id || u.ID), 0)
    const newId = maxId > 0 ? maxId + 1 : 1
    users.value.unshift({
      ID: newId,
      Name: editForm.name,
      Email: editForm.email,
      Role: editForm.role,
      status: editForm.status,
      CreatedAt: new Date().toISOString(),
    })
    showToast(`Pengguna "${editForm.name}" berhasil ditambahkan (visual saja).`, 'success')
  }
  showEditModal.value = false
}

// ─── Reset Filters ────────────────────────────────────────────────
function resetFilters() {
  searchQuery.value  = ''
  roleFilter.value   = ''
  statusFilter.value = ''
}

// ─── Toast ───────────────────────────────────────────────────────
const toast = reactive({ show: false, message: '', type: 'success' })
let toastTimer = null

function showToast(message, type = 'success') {
  clearTimeout(toastTimer)
  toast.message = message
  toast.type    = type
  toast.show    = true
  toastTimer = setTimeout(() => { toast.show = false }, 3500)
}

// ─── Helpers ─────────────────────────────────────────────────────
function formatDate(iso) {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

function capitalize(str) {
  return str ? str.charAt(0).toUpperCase() + str.slice(1) : '—'
}

const AVATAR_COLORS = [
  'linear-gradient(135deg, #3b82f6, #6366f1)',
  'linear-gradient(135deg, #10b981, #059669)',
  'linear-gradient(135deg, #f59e0b, #d97706)',
  'linear-gradient(135deg, #8b5cf6, #7c3aed)',
  'linear-gradient(135deg, #ef4444, #dc2626)',
  'linear-gradient(135deg, #06b6d4, #0284c7)',
]
function avatarColor(name) {
  const idx = name.charCodeAt(0) % AVATAR_COLORS.length
  return AVATAR_COLORS[idx]
}

// ─── API Integration ──────────────────────────────────────────────
async function fetchUsers() {
  isLoading.value = true
  errorMessage.value = ''
  try {
    users.value = await adminApi.getUsers()
  } catch (err) {
    if (err.message?.toLowerCase().includes('forbidden') || err.message?.toLowerCase().includes('admin')) {
      errorMessage.value = 'Akses ditolak. Hanya admin yang dapat melihat data ini.'
    } else {
      errorMessage.value = err.message || 'Tidak dapat memuat data pengguna.'
    }
    console.error('[ManajemenPengguna] Fetch error:', err)
  } finally {
    isLoading.value = false
  }
}

// ─── Lifecycle ───────────────────────────────────────────────────
onMounted(fetchUsers)
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');

* { box-sizing: border-box; margin: 0; padding: 0; }

.page {
  font-family: 'Inter', system-ui, sans-serif;
  display: flex;
  flex-direction: column;
  gap: 20px;
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
}

.page-subtitle {
  font-size: 0.8rem;
  color: #94a3b8;
  margin-top: 3px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

/* ── Search Box ── */
.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  border: 1.5px solid #e2e8f0;
  border-radius: 9px;
  padding: 8px 13px;
  background: #fff;
  color: #94a3b8;
  transition: border-color 0.18s, box-shadow 0.18s;
}

.search-box--focused {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.search-input {
  border: none;
  background: transparent;
  outline: none;
  font-family: inherit;
  font-size: 0.825rem;
  color: #1e293b;
  width: 210px;
}

.search-input::placeholder { color: #cbd5e1; }

/* ── Buttons ── */
.btn-primary {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 9px 16px;
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: #fff;
  border: none;
  border-radius: 9px;
  font-family: inherit;
  font-size: 0.825rem;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(59,130,246,0.3);
  transition: box-shadow 0.18s, transform 0.18s, background 0.18s;
  white-space: nowrap;
}

.btn-primary:hover {
  background: linear-gradient(135deg, #2563eb, #1d4ed8);
  box-shadow: 0 4px 14px rgba(59,130,246,0.38);
  transform: translateY(-1px);
}

.btn-ghost {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 9px 16px;
  background: #fff;
  color: #64748b;
  border: 1.5px solid #e2e8f0;
  border-radius: 9px;
  font-family: inherit;
  font-size: 0.825rem;
  font-weight: 600;
  cursor: pointer;
  transition: border-color 0.15s, background 0.15s, color 0.15s;
}

.btn-ghost:hover:not(:disabled) {
  border-color: #cbd5e1;
  background: #f8fafc;
  color: #1e293b;
}

.btn-danger {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 9px 20px;
  background: #ef4444;
  color: #fff;
  border: none;
  border-radius: 9px;
  font-family: inherit;
  font-size: 0.825rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s, transform 0.15s;
}

.btn-danger:hover:not(:disabled) { background: #dc2626; transform: translateY(-1px); }
.btn-danger:disabled, .btn-ghost:disabled { opacity: 0.55; cursor: not-allowed; }

/* ── Mini Stats ── */
.stats-row {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.mini-stat {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fff;
  border: 1px solid #f1f5f9;
  border-radius: 10px;
  padding: 10px 14px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
  flex: 1;
  min-width: 100px;
}

.mini-stat-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--c);
  flex-shrink: 0;
}

.mini-stat-value {
  font-size: 1.05rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1;
}

.mini-stat-label {
  font-size: 0.72rem;
  color: #94a3b8;
  white-space: nowrap;
}

/* ── Card ── */
.card {
  background: #fff;
  border-radius: 16px;
  border: 1px solid #f1f5f9;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.05);
  overflow: hidden;
}

/* ── Toolbar ── */
.card-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 20px;
  border-bottom: 1px solid #f1f5f9;
  flex-wrap: wrap;
}

.result-count {
  font-size: 0.8rem;
  color: #64748b;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-select {
  border: 1.5px solid #e2e8f0;
  border-radius: 8px;
  padding: 6px 10px;
  font-family: inherit;
  font-size: 0.78rem;
  color: #475569;
  background: #f8fafc;
  cursor: pointer;
  outline: none;
  transition: border-color 0.15s;
}

.filter-select:focus { border-color: #3b82f6; }

/* ── Skeleton ── */
.skeleton-table {
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.skeleton-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.skel {
  height: 14px;
  border-radius: 5px;
  background: #f1f5f9;
  animation: pulse 1.4s ease-in-out infinite;
}

.skel--sm  { flex: 0.6; }
.skel--md  { flex: 1.2; }
.skel--lg  { flex: 2;   }

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50%       { opacity: 0.45; }
}

/* ── Table ── */
.table-wrapper { overflow-x: auto; }

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

.th-id     { width: 60px; }
.th-center { text-align: center !important; }

.table-row {
  border-bottom: 1px solid #f8fafc;
  transition: background 0.12s;
}

.table-row:last-child { border-bottom: none; }
.table-row:hover      { background: #fafbfc; }
.row--inactive        { opacity: 0.62; }

.data-table td { padding: 13px 16px; vertical-align: middle; color: #334155; }

.td-id { font-size: 0.72rem; color: #94a3b8; font-weight: 600; font-variant-numeric: tabular-nums; }

/* ── User Chip ── */
.user-chip {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.user-name {
  display: block;
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
}

.user-joined {
  display: block;
  font-size: 0.7rem;
  color: #94a3b8;
  margin-top: 1px;
}

.td-email {
  color: #475569;
  font-size: 0.8rem;
}

/* ── Role Badge ── */
.badge-role {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 20px;
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 0.03em;
  white-space: nowrap;
}

.badge-role--admin     { background: #fef9c3; color: #854d0e; }
.badge-role--penggalang { background: #ede9fe; color: #6d28d9; }
.badge-role--donatur    { background: #dbeafe; color: #1d4ed8; }

/* ── Status Toggle ── */
.status-toggle {
  position: relative;
  display: inline-flex;
  align-items: center;
  width: 38px;
  height: 20px;
  border-radius: 20px;
  border: none;
  cursor: pointer;
  transition: background 0.2s;
  padding: 0;
  flex-shrink: 0;
}

.status-toggle--on  { background: #22c55e; }
.status-toggle--off { background: #cbd5e1; }

.toggle-knob {
  position: absolute;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: #fff;
  transition: left 0.2s;
  box-shadow: 0 1px 3px rgba(0,0,0,0.2);
}

.status-toggle--on  .toggle-knob { left: 21px; }
.status-toggle--off .toggle-knob { left: 3px; }

/* ── Action Buttons ── */
.action-group {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.btn-action {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 10px;
  border: none;
  border-radius: 7px;
  font-family: inherit;
  font-size: 0.72rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s, transform 0.15s;
  white-space: nowrap;
}

.btn-action:hover { transform: translateY(-1px); }

.btn-edit   { background: #eff6ff; color: #2563eb; }
.btn-delete { background: #fef2f2; color: #dc2626; }
.btn-edit:hover   { background: #dbeafe; }
.btn-delete:hover { background: #fee2e2; }

/* ── Empty State ── */
.empty-state {
  padding: 56px 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  text-align: center;
}

.empty-icon-wrap {
  width: 56px;
  height: 56px;
  background: #f1f5f9;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  margin-bottom: 4px;
}

.empty-title { font-size: 0.9rem; font-weight: 700; color: #475569; }
.empty-sub   { font-size: 0.78rem; color: #94a3b8; }

/* ── Modal ── */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  backdrop-filter: blur(3px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9000;
  padding: 24px;
}

.modal-card {
  background: #fff;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.18);
  padding: 36px 36px 32px;
  width: 100%;
  max-width: 420px;
  text-align: center;
}

.modal-card--wide { max-width: 580px; text-align: left; }

.modal-icon {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.modal-icon--danger { background: #fef2f2; color: #ef4444; }

.modal-title {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 10px;
}

.modal-title--left { text-align: left; margin-bottom: 20px; }

.modal-body {
  font-size: 0.85rem;
  color: #64748b;
  line-height: 1.6;
  margin-bottom: 24px;
}

.modal-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 24px;
}

/* ── Edit Form ── */
.edit-form { display: flex; flex-direction: column; gap: 14px; }

.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }

.form-group { display: flex; flex-direction: column; gap: 5px; }

.form-label { font-size: 0.78rem; font-weight: 600; color: #374151; }

.form-input {
  border: 1.5px solid #e2e8f0;
  border-radius: 9px;
  padding: 10px 12px;
  font-family: inherit;
  font-size: 0.825rem;
  color: #0f172a;
  background: #f8fafc;
  outline: none;
  transition: border-color 0.18s, box-shadow 0.18s, background 0.18s;
}

.form-input:focus {
  border-color: #3b82f6;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(59,130,246,0.1);
}

.form-select { cursor: pointer; }

/* ── Toast ── */
.toast {
  position: fixed;
  bottom: 28px;
  right: 28px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 13px 18px;
  border-radius: 12px;
  font-family: inherit;
  font-size: 0.825rem;
  font-weight: 600;
  box-shadow: 0 8px 24px rgba(0,0,0,0.14);
  z-index: 9999;
  pointer-events: none;
}

.toast--success { background: #0f172a; color: #fff; }
.toast--error   { background: #ef4444; color: #fff; }

/* Spinner in button */
.btn-spinner-wrap { display: flex; align-items: center; gap: 8px; }
.btn-spinner {
  display: inline-block;
  width: 14px; height: 14px;
  border: 2px solid rgba(255,255,255,0.35);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.65s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* ── Transitions ── */
.modal-enter-active, .modal-leave-active { transition: opacity 0.22s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }

.modal-enter-active .modal-card, .modal-leave-active .modal-card { transition: transform 0.22s ease; }
.modal-enter-from .modal-card, .modal-leave-to .modal-card { transform: scale(0.95) translateY(8px); }

.toast-enter-active, .toast-leave-active { transition: opacity 0.25s ease, transform 0.25s ease; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(10px); }

/* ── Responsive ── */
@media (max-width: 700px) {
  .page-header   { flex-direction: column; align-items: flex-start; }
  .header-actions { width: 100%; }
  .search-input  { width: 100%; }
  .form-row      { grid-template-columns: 1fr; }
  .stats-row     { flex-direction: column; }
}
</style>
