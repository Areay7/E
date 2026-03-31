import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录', requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '数据看板', requiresAuth: true }
      },
      {
        path: 'shopee/orders',
        name: 'ShopeeOrders',
        component: () => import('@/views/shopee/orders.vue'),
        meta: { title: 'Shopee 订单管理', requiresAuth: true }
      },
      {
        path: 'aliexpress/orders',
        name: 'AliExpressOrders',
        component: () => import('@/views/aliexpress/orders.vue'),
        meta: { title: '速卖通 订单管理', requiresAuth: true }
      },
      {
        path: 'tk/orders',
        name: 'TikTokOrders',
        component: () => import('@/views/tk/orders.vue'),
        meta: { title: 'TikTok 订单管理', requiresAuth: true }
      },
      {
        path: 'platform/config',
        name: 'PlatformConfig',
        component: () => import('@/views/platform/config.vue'),
        meta: { title: '平台配置管理', requiresAuth: true }
      },
      {
        path: 'products',
        name: 'Products',
        component: () => import('@/views/products/index.vue'),
        meta: { title: '商品管理', requiresAuth: true }
      },
      {
        path: 'inventory',
        name: 'Inventory',
        component: () => import('@/views/inventory/index.vue'),
        meta: { title: '库存管理', requiresAuth: true }
      }
    ]
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '页面不存在', requiresAuth: false }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const appTitle = '跨境电商'

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false)

  // 需要认证但没有token，跳转到登录页
  if (requiresAuth && !token) {
    next('/login')
  }
  // 已登录访问登录页，跳转到首页
  else if (to.path === '/login' && token) {
    next('/')
  }
  // 正常访问
  else {
    next()
  }
})

router.afterEach((to) => {
  const t = to.meta.title
  document.title = typeof t === 'string' && t.length > 0 ? `${t} · ${appTitle}` : appTitle
})

export default router