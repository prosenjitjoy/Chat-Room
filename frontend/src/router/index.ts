import { createRouter, createWebHistory } from "vue-router";
import HomeView from "@/views/HomeView.vue";
import RegisterView from "../views/RegisterView.vue";
import LoginView from "../views/LoginView.vue";
import AppView from "@/views/AppView.vue";

import { useStore } from "@/stores";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/signin",
      name: "login",
      component: LoginView,
    },
    {
      path: "/signup",
      name: "register",
      component: RegisterView,
    },
    {
      path: "/",
      name: "home",
      component: HomeView,
    },
    {
      path: "/room/:id",
      name: "app",
      component: AppView
    },
  ],
});

router.beforeEach(async (to, from) => {
  const checkAuth = (): boolean => {
    const authFlag = sessionStorage.getItem("isAuthenticated")
    if (authFlag === "true") {
      return true
    }
    return false
  }
  useStore().isAuthenticated = checkAuth()

  if (
    // make sure the user is authenticated
    !useStore().isAuthenticated &&
    // ❗️ Avoid an infinite redirect
    to.name !== 'login' && to.name !== 'register'
  ) {
    // redirect the user to the login page
    return { name: 'login' }
  }
})


export default router;
