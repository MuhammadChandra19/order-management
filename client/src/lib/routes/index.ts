import { createRouter, createWebHistory } from 'vue-router'
import Main from '@/pages/Main/index.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'main',
      component: Main,
    },
    {
      path: '/order',
      name: 'order',
      component: () => import('@/pages/Order/index.vue')
    }
  ],
})
