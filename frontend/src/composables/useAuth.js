import { reactive, computed } from 'vue'

const state = reactive({
  user: JSON.parse(localStorage.getItem('opendonasi_user') || 'null'),
  token: localStorage.getItem('opendonasi_token') || null,
})

export function useAuth() {
  const isLoggedIn = computed(() => !!state.user)
  const user = computed(() => state.user)
  const role = computed(() => state.user?.role || null)

  function login(userData, token) {
    state.user = userData
    state.token = token
    localStorage.setItem('opendonasi_user', JSON.stringify(userData))
    localStorage.setItem('opendonasi_token', token)
  }

  function logout() {
    state.user = null
    state.token = null
    localStorage.removeItem('opendonasi_user')
    localStorage.removeItem('opendonasi_token')
  }

  function getToken() {
    return state.token
  }

  return { isLoggedIn, user, role, login, logout, getToken }
}
