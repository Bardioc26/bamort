import { createRouter, createWebHistory } from "vue-router";
import { isLoggedIn } from "../utils/auth"; // Import the helper function
import LoginView from "../views/LoginView.vue";
import RegisterView from "../views/RegisterView.vue";
import ForgotPasswordView from "../views/ForgotPasswordView.vue";
import ResetPasswordView from "../views/ResetPasswordView.vue";
import DashboardView from "../views/DashboardView.vue";
import AusruestungView from "../views/AusruestungView.vue";
import MaintenanceView from "../views/MaintenanceView.vue";
import FileUploadPage from "../views/FileUploadPage.vue";

import CharacterDetails from "@/components/CharacterDetails.vue";
import CharacterCreation from "@/components/CharacterCreation.vue";




const routes = [
  { path: "/", name: "Login", component: LoginView },
  { path: "/register", name: "Register", component: RegisterView },
  { path: "/forgot-password", name: "ForgotPassword", component: ForgotPasswordView },
  { path: "/reset-password", name: "ResetPassword", component: ResetPasswordView },
  { path: "/dashboard", name: "Dashboard", component: DashboardView, meta: { requiresAuth: true } },
  { path: "/ausruestung/:characterId", name: "Ausruestung", component: AusruestungView, meta: { requiresAuth: true } },
  { path: "/maintenance", name: "Maintenance", component: MaintenanceView, meta: { requiresAuth: true } },
  { path: "/upload", name: "FileUpload", component: FileUploadPage },
  // Route for character details  // Pass route params as props to the component
  {    path: "/character/:id",     name: "CharacterDetails",    component: CharacterDetails,    props: true, meta: { requiresAuth: true }   },
  // Route for character creation
  {    path: "/character/create/:sessionId",     name: "CharacterCreation",    component: CharacterCreation,    props: true, meta: { requiresAuth: true }   },
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
