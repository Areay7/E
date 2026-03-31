import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '数据看板' }
      },
      {
        path: 'shopee/orders',
        name: 'ShopeeOrders',
        component: () => import('@/views/shopee/orders.vue'),
        meta: { title: 'Shopee 订单管理' }
      },
      {
        path: 'aliexpress/orders',
        name: 'AliExpressOrders',
        component: () => import('@/views/aliexpress/orders.vue'),
        meta: { title: '速卖通 订单管理' }
      },
      {
        path: 'tk/orders',
        name: 'TikTokOrders',
        component: () => import('@/views/tk/orders.vue'),
        meta: { title: 'TikTok 订单管理' }
      },
      {
        path: 'platform/config',
        name: 'PlatformConfig',
        component: () => import('@/views/platform/config.vue'),
        meta: { title: '平台配置管理' }
      },
      {
        path: 'products',
        name: 'Products',
        component: () => import('@/views/products/index.vue'),
        meta: { title: '商品管理' }
      },
      {
        path: 'inventory',
        name: 'Inventory',
        component: () => import('@/views/inventory/index.vue'),
        meta: { title: '库存管理' }
      }
    ]
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '页面不存在' }
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

router.afterEach((to) => {
  const t = to.meta.title
  document.title = typeof t === 'string' && t.length > 0 ? `${t} · ${appTitle}` : appTitle
})

export default router