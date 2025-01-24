import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('@/views/AboutView.vue'),
    },
    {
      path: '/signup',
      name: 'signup',
      component: () => import('@/views/SignUpView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
    },
    {
      path: '/logout',
      name: 'logout',
      component: () => import('@/components/LogoutComponent.vue'),
    },
    {
      path: '/login_sms',
      name: 'login_sms',
      component: () => import('@/views/LoginSMSView.vue'),
    },
    {
      path: '/login_gitea',
      name: 'login_gitea',
      component: () => import('@/views/LoginGiteaView.vue'),
    },
    {
      path: '/user/profile',
      name: 'user profile',
      component: () => import('@/views/users/ProfileView.vue'),
    },
    {
      path: '/oauth2success',
      name: 'oauth2 callback',
      component: () => import('@/components/OAuth2Success.vue'),
    },
    {
      path: '/articles/list',
      name: 'articles list',
      component: () => import('@/views/articles/ListView.vue'),
    },
    {
      path: '/articles/edit',
      name: 'articles edit',
      component: () => import('@/views/articles/EditView.vue'),
    },
    {
      path: '/articles/view',
      name: 'articles view',
      component: () => import('@/views/articles/ReadView.vue'),
    },
  ],
})

export default router
