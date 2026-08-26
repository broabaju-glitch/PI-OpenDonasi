import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/composables/useAuth'


const AdminLogin          = () => import('@/views/AdminLogin.vue')
const AdminDashboard      = () => import('@/views/AdminDashboard.vue')
const PlatformOverview    = () => import('@/views/PlatformOverview.vue')
const ManajemenPengguna   = () => import('@/views/ManajemenPengguna.vue')
const ModerasiCampaign    = () => import('@/views/ModerasiCampaign.vue')
const ManajemenTransaksi  = () => import('@/views/ManajemenTransaksi.vue')

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior: () => ({ top: 0 }),
  routes: [
    
    {
      path: '/',
      redirect: '/admin/platform-overview',
    },

    
    {
      path: '/admin/login',
      name: 'admin-login',
      component: AdminLogin,
      meta: { guestOnly: true },
    },

    
    {
      path: '/admin',
      name: 'admin-dashboard',
      component: AdminDashboard,
      meta: { requiresAuth: true, role: 'admin' },
      redirect: '/admin/platform-overview',
      children: [
        {
          path: 'platform-overview',
          name: 'admin-platform-overview',
          component: PlatformOverview,
        },
        {
          path: 'manajemen-pengguna',
          name: 'admin-manajemen-pengguna',
          component: ManajemenPengguna,
        },
        {
          path: 'moderasi-campaign',
          name: 'admin-moderasi-campaign',
          component: ModerasiCampaign,
        },
        {
          path: 'manajemen-transaksi',
          name: 'admin-manajemen-transaksi',
          component: ManajemenTransaksi,
        },
      ],
    },

    
    {
      path: '/:pathMatch(.*)*',
      redirect: '/admin/platform-overview',
    },
  ],
})


router.beforeEach((to, _from, next) => {
  const { isLoggedIn, role } = useAuth()

  
  if (to.meta.guestOnly && isLoggedIn.value && role.value === 'admin') {
    return next('/admin/platform-overview')
  }

  
  if (to.meta.requiresAuth && !isLoggedIn.value) {
    return next('/admin/login')
  }

  
  if (to.meta.requiresAuth && isLoggedIn.value && role.value !== 'admin') {
    return next('/admin/login')
  }

  next()
})

export default router
