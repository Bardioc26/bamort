import { createRouter, createWebHistory } from "vue-router";
import { isLoggedIn } from "../utils/auth";
import { useUserStore } from '../stores/userStore';

// Static imports for landing/auth pages (needed immediately)
import LandingView from "../views/LandingView.vue";
import LoginView from "../views/LoginView.vue";
import RegisterView from "../views/RegisterView.vue";
import ForgotPasswordView from "../views/ForgotPasswordView.vue";
import ResetPasswordView from "../views/ResetPasswordView.vue";

// Lazy-loaded views (code-split into separate chunks)
const DashboardView = () => import("../views/DashboardView.vue");
const AusruestungView = () => import("../views/AusruestungView.vue");
const MaintenanceView = () => import("../views/MaintenanceView.vue");
const FileUploadPage = () => import("../views/FileUploadPage.vue");
const UserProfileView = () => import("../views/UserProfileView.vue");
const UserManagementView = () => import("../views/UserManagementView.vue");
const SponsorsView = () => import("../views/SponsorsView.vue");
const HelpView = () => import("../views/HelpView.vue");
const SystemInfoView = () => import("../views/SystemInfoView.vue");
const CharacterDetails = () => import("@/components/CharacterDetails.vue");
const CharacterCreation = () => import("@/components/CharacterCreation.vue");




const routes = [
  { path: "/", name: "Landing", component: LandingView },
  { path: "/login", name: "Login", component: LoginView },
  { path: "/register", name: "Register", component: RegisterView },
  { path: "/forgot-password", name: "ForgotPassword", component: ForgotPasswordView },
  { path: "/reset-password", name: "ResetPassword", component: ResetPasswordView },
  { path: "/dashboard", name: "Dashboard", component: DashboardView, meta: { requiresAuth: true } },
  { path: "/profile", name: "UserProfile", component: UserProfileView, meta: { requiresAuth: true } },
  { path: "/users", name: "UserManagement", component: UserManagementView, meta: { requiresAuth: true, requiresAdmin: true } },
  { path: "/ausruestung/:characterId", name: "Ausruestung", component: AusruestungView, meta: { requiresAuth: true } },
  { path: "/maintenance", name: "Maintenance", component: MaintenanceView, meta: { requiresAuth: true } },
  { path: "/upload", name: "FileUpload", component: FileUploadPage, meta: { requiresAuth: true } },
  { path: "/sponsors", name: "Sponsors", component: SponsorsView },
  { path: "/help", name: "Help", component: HelpView },
  { path: "/system-info", name: "SystemInfo", component: SystemInfoView },
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
router.beforeEach(async (to, from, next) => {
  if (to.meta.requiresAuth && !isLoggedIn()) {
    // Redirect to login if not authenticated
    next({ name: "Login" });
  } else if (to.meta.requiresAdmin) {
    // Check if route requires admin role
    const userStore = useUserStore()
    
    // Fetch user if not already loaded
    if (!userStore.currentUser) {
      await userStore.fetchCurrentUser()
    }
    
    if (!userStore.isAdmin) {
      // Redirect to dashboard if not admin
      next({ name: "Dashboard" });
    } else {
      next();
    }
  } else {
    next(); // Allow navigation
  }
});

export default router;
