/**
 * api.js — Centralized API service for OpenDonasi Admin Panel
 *
 * All requests use the Vite dev proxy (/api → http://localhost:8080/api)
 * so there are zero CORS issues in development.
 * In production, point BASE_URL to the real backend.
 */

const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'

// ─── Core request helper ──────────────────────────────────────────
async function request(method, path, { body } = {}) {
  const token = localStorage.getItem('opendonasi_admin_token')
  const headers = { 'Content-Type': 'application/json' }
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch(`${BASE_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  })

  // Always try to parse JSON — even for error responses
  let data
  try {
    data = await res.json()
  } catch {
    data = null
  }

  if (!res.ok) {
    // Throw an Error with the server's message so callers can catch it
    const msg = data?.error || `Server error ${res.status}`
    throw new Error(msg)
  }

  return data
}

// ─── Convenience methods ─────────────────────────────────────────
const api = {
  get:    (path, opts)       => request('GET',    path, opts),
  post:   (path, body, opts) => request('POST',   path, { body, ...opts }),
  put:    (path, body, opts) => request('PUT',    path, { body, ...opts }),
  patch:  (path, body, opts) => request('PATCH',  path, { body, ...opts }),
  delete: (path, opts)       => request('DELETE', path, opts),
}

export default api

// ─── Named Auth API ──────────────────────────────────────────────
export const authApi = {
  /**
   * POST /api/login
   * Returns { user: { id, name, email, role }, token }
   */
  login: (email, password) =>
    api.post('/login', { email, password }),
}

// ─── Named Admin API ─────────────────────────────────────────────
export const adminApi = {
  /**
   * GET /api/statistics
   * Returns { active_campaigns, total_collected, total_donatur, completed_campaigns }
   */
  getStatistics: () =>
    api.get('/statistics'),

  /**
   * GET /api/admin/donations  (JWT required)
   * Returns array of all donation records
   */
  getDonations: () =>
    api.get('/admin/donations'),

  /**
   * GET /api/users  (JWT required, role: admin)
   * Returns array: [{ id, name, email, role, status, created_at }]
   */
  getUsers: () =>
    api.get('/users'),

  /**
   * DELETE /api/users/:id (JWT required, role: admin)
   * Deletes a user
   */
  deleteUser: (id) =>
    api.delete(`/users/${id}`),

  /**
   * GET /api/admin/campaigns
   * Returns array of all campaigns
   */
  getCampaigns: () =>
    api.get('/admin/campaigns'),

  /**
   * PATCH /api/admin/campaigns/:id/status
   * Body: { status: 'disetujui' | 'ditolak' | 'selesai' }
   */
  updateCampaignStatus: (id, status) =>
    api.patch(`/admin/campaigns/${id}/status`, { status }),

  /**
   * DELETE /api/admin/campaigns/:id
   * Soft deletes a campaign
   */
  deleteCampaign: (id) =>
    api.delete(`/admin/campaigns/${id}`),

  /**
   * GET /api/transactions (alias for getDonations or standalone)
   * Returns array of all transactions/donations
   */
  getTransactions: () =>
    api.get('/transactions'),

  /**
   * PATCH /api/admin/transactions/:id/cairkan
   * Updates a transaction status to 'Dicairkan'
   */
  cairkanTransaction: (id) =>
    api.patch(`/admin/transactions/${id}/cairkan`),
}
