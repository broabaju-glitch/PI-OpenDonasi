import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior() {
    return { top: 0 }
  },
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/Home.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue'),
      meta: { guestOnly: true }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/Register.vue'),
      meta: { guestOnly: true }
    },
    {
      path: '/campaigns',
      name: 'campaigns',
      component: () => import('../views/CampaignList.vue')
    },
    {
      path: '/campaign/:id',
      name: 'campaign-detail',
      component: () => import('../views/CampaignDetail.vue')
    },

    {
      path: '/penggalang',
      name: 'penggalang-dashboard',
      component: () => import('../views/dashboard/penggalang/PenggalangDashboard.vue'),
      meta: { requiresAuth: true, role: 'penggalang' },
      redirect: '/penggalang/campaign-saya',
      children: [
        { path: 'campaign-saya', name: 'penggalang-campaigns', component: () => import('../views/dashboard/penggalang/CampaignSaya.vue') },
        { path: 'buat-campaign', name: 'penggalang-buat-campaign', component: () => import('../views/dashboard/penggalang/BuatCampaign.vue') },
        { path: 'donasi-masuk', name: 'penggalang-donasi-masuk', component: () => import('../views/dashboard/penggalang/DonasiMasuk.vue') },
        { path: 'status-dana', name: 'penggalang-status-dana', component: () => import('../views/dashboard/penggalang/StatusDana.vue') },
        { path: 'upload-laporan', name: 'penggalang-upload-laporan', component: () => import('../views/dashboard/penggalang/UploadLaporan.vue') },
      ]
    },
  
    {
      path: '/donatur',
      name: 'donatur-dashboard',
      component: () => import('../views/dashboard/donatur/DonaturDashboard.vue'),
      meta: { requiresAuth: true, role: 'donatur' },
      redirect: '/donatur/profil',
      children: [
        { path: 'profil', name: 'donatur-profil', component: () => import('../views/dashboard/donatur/ProfilSaya.vue') },
        { path: 'riwayat', name: 'donatur-riwayat', component: () => import('../views/dashboard/donatur/RiwayatDonasi.vue') },
        { path: 'laporan', name: 'donatur-laporan', component: () => import('../views/dashboard/donatur/LaporanCampaign.vue') },
      ]
    },

  ]
})


router.beforeEach((to, from, next) => {
  const { isLoggedIn, role } = useAuth()

  
  if (to.meta.guestOnly && isLoggedIn.value) {
    if (role.value === 'penggalang') return next('/penggalang')
    if (role.value === 'donatur') return next('/donatur')
    return next('/')
  }

  
  if (to.meta.requiresAuth && !isLoggedIn.value) {
    return next('/login')
  }

  
  if (isLoggedIn.value && to.meta.role && role.value !== to.meta.role) {
    return next('/')
  }

  next()
})

export default router
