import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import Timeline from '../views/Timeline.vue'
import Reports from '../views/Reports.vue'
import Categories from '../views/Categories.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: Dashboard
    },
    {
      path: '/timeline',
      name: 'timeline',
      component: Timeline
    },
    {
      path: '/reports',
      name: 'reports',
      component: Reports
    },
    {
      path: '/categories',
      name: 'categories',
      component: Categories
    }
  ]
})

export default router
