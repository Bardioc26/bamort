import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import AusruestungView from '../views/AusruestungView.vue'

const routes = [
  { path: '/', name: 'Login', component: LoginView },
  { path: '/dashboard', name: 'Dashboard', component: DashboardView },
  { path: '/ausruestung/:characterId', name: 'Ausruestung', component: AusruestungView },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
