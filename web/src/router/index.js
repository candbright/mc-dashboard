import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'main',
      component: () => import('../views/main/Main.vue'),
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('../views/home/Home.vue')
        },
        {
          path: 'server_list',
          name: 'server_list',
          component: () => import('../views/server/list/ServerList.vue')
        },
        {
          path: 'server_list/:id/info',
          name: 'server_info',
          component: () => import('../views/server/list/ServerListItem.vue'),
          props: (route) => ({
            id: route.params.id,
            row: route.query.row
          })
        },
        {
          path: 'save_list',
          name: 'save_list',
          component: () => import('../views/save/list/SaveList.vue')
        },
        {
          path: 'impl',
          name: 'impl',
          component: () => import('../views/server/Impl.vue')
        }
      ]
    }
  ]
})

export default router
