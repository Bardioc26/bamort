import { createRouter, createWebHistory } from "vue-router";
import { isLoggedIn } from "../utils/auth"; // Import the helper function
import LoginView from "../views/LoginView.vue";
import RegisterView from "../views/RegisterView.vue";
import DashboardView from "../views/DashboardView.vue";
import AusruestungView from "../views/AusruestungView.vue";

const routes = [
  { path: "/", name: "Login", component: LoginView },
  { path: "/register", name: "Register", component: RegisterView },
  { path: "/dashboard", name: "Dashboard", component: DashboardView, meta: { requiresAuth: true } },
  { path: "/ausruestung/:characterId", name: "Ausruestung", component: AusruestungView, meta: { requiresAuth: true } },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  //linkActiveClass: "current-page", // Custom active class
});

// Navigation guard
router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !isLoggedIn()) {
    // Redirect to login if not authenticated
    next({ name: "Login" });
  } else {
    next(); // Allow navigation
  }
});

export default router;
