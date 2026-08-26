import { reactive, computed } from 'vue'

// Module-level singleton state — persists across hot-reloads
const state = reactive({
  user:  JSON.parse(localStorage.getItem('opendonasi_admin_user')  || 'null'),
  token: localStorage.getItem('opendonasi_admin_token') || null,
})

export function useAuth() {
  const isLoggedIn = computed(() => !!state.user)
  const user       = computed(() => state.user)
  const role       = computed(() => state.user?.role || null)

  function login(userData, token) {
    state.user  = userData
    state.token = token
    localStorage.setItem('opendonasi_admin_user',  JSON.stringify(userData))
    localStorage.setItem('opendonasi_admin_token', token)
  }

  function logout() {
    state.user  = null
    state.token = null
    localStorage.removeItem('opendonasi_admin_user')
    localStorage.removeItem('opendonasi_admin_token')
  }

  function getToken() {
    return state.token
  }

  return { isLoggedIn, user, role, login, logout, getToken }
}
