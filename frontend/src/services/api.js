
const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'


async function request(endpoint, options = {}) {
  const url = `${API_BASE}${endpoint}`

  const config = {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  }

  
  if (options.body instanceof FormData) {
    delete config.headers['Content-Type']
  }

 
  const token = localStorage.getItem('opendonasi_token')
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(url, config)
  const data = await response.json()

  if (!response.ok) {
    throw new Error(data.error || `Request failed with status ${response.status}`)
  }

  return data
}



export async function apiLogin(email, password) {
  return request('/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export async function apiRegister(name, email, password, role) {
  return request('/register', {
    method: 'POST',
    body: JSON.stringify({ name, email, password, role }),
  })
}

export async function apiGetProfile() {
  return request('/profile', {
    method: 'GET'
  })
}



export async function apiGetCampaigns() {
  return request('/campaigns')
}

export async function apiGetCampaign(id) {
  return request(`/campaigns/${id}`)
}



export async function apiCreateDonation(payload) {
  const token = localStorage.getItem('opendonasi_token')
  const isFormData = payload instanceof FormData

  return request('/donations', {
    method: 'POST',
    headers: {
      'Authorization': token ? `Bearer ${token}` : ''
    },
    body: isFormData ? payload : JSON.stringify(payload),
  })
}

export async function apiGetCampaignDonations(campaignId) {
  return request(`/campaigns/${campaignId}/donations`)
}

export async function apiGetMyDonations() {
  return request('/donations/me')
}




export async function apiGetStats() {
  return request('/stats')
}

export async function apiGetStatistics() {
  return request('/statistics')
}



export async function apiGetFundraiserSummary() {
  return request('/fundraiser/funds/summary')
}

export async function apiGetFundraiserDonations() {
  return request('/fundraiser/donations')
}

export async function apiRequestWithdrawal(payload) {
  return request('/fundraiser/withdrawals', {
    method: 'POST',
    body: JSON.stringify(payload)
  })
}

export async function apiUploadReport(formData) {
  const token = localStorage.getItem('opendonasi_token')
  return request('/fundraiser/reports', {
    method: 'POST',
    headers: {
      'Authorization': token ? `Bearer ${token}` : ''
    },
    body: formData
  })
}

export default {
  apiLogin,
  apiRegister,
  apiGetProfile,
  apiGetCampaigns,
  apiGetCampaign,
  apiCreateDonation,
  apiGetCampaignDonations,
  apiGetStats,
  apiGetStatistics,
  apiGetFundraiserSummary,
  apiGetFundraiserDonations,
  apiRequestWithdrawal,
  apiUploadReport,
}
