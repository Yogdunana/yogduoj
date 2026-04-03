import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: AppLayout,
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
      },
      {
        path: 'login',
        name: 'Login',
        component: () => import('@/views/Login.vue'),
        meta: { guest: true },
      },
      {
        path: 'register',
        name: 'Register',
        component: () => import('@/views/Register.vue'),
        meta: { guest: true },
      },
      {
        path: 'problems',
        name: 'ProblemList',
        component: () => import('@/views/ProblemList.vue'),
      },
      {
        path: 'problems/:id',
        name: 'ProblemDetail',
        component: () => import('@/views/ProblemDetail.vue'),
        props: true,
      },
      {
        path: 'submissions',
        name: 'SubmissionList',
        component: () => import('@/views/SubmissionList.vue'),
        meta: { requiresAuth: true },
      },
      {
        path: 'submissions/:id',
        name: 'SubmissionDetail',
        component: () => import('@/views/SubmissionDetail.vue'),
        meta: { requiresAuth: true },
        props: true,
      },
      {
        path: 'contests',
        name: 'ContestList',
        component: () => import('@/views/ContestList.vue'),
      },
      {
        path: 'contests/:id',
        name: 'ContestDetail',
        component: () => import('@/views/ContestDetail.vue'),
        props: true,
      },
      {
        path: 'contests/:id/ranking',
        name: 'ContestRanking',
        component: () => import('@/views/ContestRanking.vue'),
        props: true,
      },
      {
        path: 'profile',
        name: 'UserProfile',
        component: () => import('@/views/UserProfile.vue'),
        meta: { requiresAuth: true },
      },
      {
        path: 'users/:id',
        name: 'UserPublicProfile',
        component: () => import('@/views/UserPublicProfile.vue'),
        props: true,
      },
      {
        path: 'teams',
        name: 'TeamList',
        component: () => import('@/views/TeamList.vue'),
      },
      {
        path: 'teams/create',
        name: 'TeamCreate',
        component: () => import('@/views/TeamCreate.vue'),
        meta: { requiresAuth: true },
      },
      {
        path: 'teams/:id',
        name: 'TeamDetail',
        component: () => import('@/views/TeamDetail.vue'),
        props: true,
      },
      {
        path: 'ctf',
        name: 'CTFPractice',
        component: () => import('@/views/CTFPractice.vue'),
      },
      {
        path: 'ctf/:category',
        name: 'CTFCategoryView',
        component: () => import('@/views/CTFCategoryView.vue'),
        props: true,
      },
      {
        path: 'announcements',
        name: 'AnnouncementList',
        component: () => import('@/views/AnnouncementList.vue'),
      },
      {
        path: 'announcements/:id',
        name: 'AnnouncementDetail',
        component: () => import('@/views/AnnouncementDetail.vue'),
        props: true,
      },
      {
        path: 'help',
        name: 'HelpCenter',
        component: () => import('@/views/HelpCenter.vue'),
      },
    ],
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue'),
      },
      {
        path: 'users',
        name: 'AdminUserManage',
        component: () => import('@/views/admin/UserManage.vue'),
      },
      {
        path: 'problems',
        name: 'AdminProblemManage',
        component: () => import('@/views/admin/ProblemManage.vue'),
      },
      {
        path: 'problems/create',
        name: 'AdminProblemCreate',
        component: () => import('@/views/admin/ProblemCreate.vue'),
      },
      {
        path: 'problems/:id/edit',
        name: 'AdminProblemEdit',
        component: () => import('@/views/admin/ProblemEdit.vue'),
        props: true,
      },
      {
        path: 'contests',
        name: 'AdminContestManage',
        component: () => import('@/views/admin/ContestManage.vue'),
      },
      {
        path: 'contests/create',
        name: 'AdminContestCreate',
        component: () => import('@/views/admin/ContestCreate.vue'),
      },
      {
        path: 'contests/:id/edit',
        name: 'AdminContestEdit',
        component: () => import('@/views/admin/ContestEdit.vue'),
        props: true,
      },
      {
        path: 'submissions',
        name: 'AdminSubmissionManage',
        component: () => import('@/views/admin/SubmissionManage.vue'),
      },
      {
        path: 'judge',
        name: 'AdminJudgeMonitor',
        component: () => import('@/views/admin/JudgeMonitor.vue'),
      },
      {
        path: 'announcements',
        name: 'AdminAnnouncementManage',
        component: () => import('@/views/admin/AnnouncementManage.vue'),
      },
      {
        path: 'ai-problems',
        name: 'AdminAIProblemManage',
        component: () => import('@/views/admin/AIProblemManage.vue'),
      },
      {
        path: 'import',
        name: 'AdminImportManage',
        component: () => import('@/views/admin/ImportManage.vue'),
      },
      {
        path: 'cheat',
        name: 'AdminCheatManage',
        component: () => import('@/views/admin/CheatManage.vue'),
      },
      {
        path: 'config',
        name: 'AdminSystemConfig',
        component: () => import('@/views/admin/SystemConfig.vue'),
      },
      {
        path: 'diy-templates',
        name: 'AdminDIYTemplates',
        component: () => import('@/views/admin/DIYTemplates.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  },
})

// Navigation guards
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  const isLoggedIn = !!token

  if (to.meta.requiresAuth && !isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.meta.guest && isLoggedIn) {
    next({ name: 'Home' })
  } else if (to.meta.requiresAdmin && isLoggedIn) {
    // Check admin status - will be refined once auth store is available
    const userStr = localStorage.getItem('user')
    if (userStr) {
      try {
        const user = JSON.parse(userStr)
        if (user.role !== 'admin') {
          next({ name: 'Home' })
          return
        }
      } catch {
        next({ name: 'Login' })
        return
      }
    }
    next()
  } else {
    next()
  }
})

export default router
