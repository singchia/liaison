import { createRouter, createWebHistory } from 'vue-router'
import ProductManagement from '../views/ProductManagement.vue'

const routes = [
  {
    path: '/',
    redirect: '/product'
  },
  {
    path: '/product',
    name: 'ProductManagement',
    component: ProductManagement
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router 