import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: to => {
        // 如果已登录，重定向到服务器列表
        return localStorage.getItem('token') ? '/server/list' : '/home'
      }
    },
    {
      path: '/',
      name: 'container',
      component: () => import('../views/Container.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: 'home',
          name: 'home',
          component: () => import('../views/Home.vue')
        },
        {
          path: '/',
          name: 'dashboard',
          component: () => import('../views/Dashboard.vue'),
          meta: { requiresAuth: true },
          children: [
            {
              path: 'server/list',
              name: 'server_list',
              component: () => import('../views/server/ServerList.vue')
            },
            {
              path: 'server/:id/detail',
              name: 'server_detail',
              component: () => import('../views/server/ServerDetail.vue'),
              props: (route) => ({
                id: route.params.id
              })
            },
            {
              path: 'save/list',
              name: 'save_list',
              component: () => import('../views/save/SaveList.vue')
            },
            {
              path: 'impl',
              name: 'impl',
              component: () => import('../views/server/Impl.vue')
            }
          ]
        }
      ]
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/auth/Login.vue'),
      meta: {
        transition: 'fade'
      }
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('token')
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  if (requiresAuth && !isAuthenticated) {
    // 如果需要登录但未登录，重定向到登录页面
    next({ name: 'login', query: { redirect: to.fullPath }})
  } else if (to.name === 'login' && isAuthenticated) {
    // 如果已登录但访问登录页面，重定向到服务器列表
    next({ name: 'server_list' })
  } else if (to.name === 'register' && isAuthenticated) {
    // 如果已登录但访问注册页面，重定向到服务器列表
    next({ name: 'server_list' })
  } else {
    next()
  }
})

export default router
