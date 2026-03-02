import { createRouter, createWebHistory } from 'vue-router'
import createRouteGuard from './guard/index'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

NProgress.configure({ showSpinner: false }) // nprogress 进度条配置
const router = createRouter({
  history: createWebHistory(process.env.BASE),
  routes: [
    {
      path: '/',
      redirect: '/user/login',
    },
    {
      path: '/user/login',
      name: 'login',
      component: () => import('@/views/user/login.vue'),
    },
    {
      path: '/user/register',
      name: 'register',
      component: () => import('@/views/user/register.vue'),
    },
    {
      path: '/user/forget',
      name: 'forget',
      component: () => import('@/views/user/forget.vue'),
    },
    {
      path: '/user/entry',
      name: 'entry',
      component: () => import('@/views/user/entry.vue'),
    },
    {
      path: '/home',
      component: () => import('@/layouts/HomeLayout.vue'),
      redirect: '/home/index',
      children: [
        {
          path: 'index',
          name: 'home',
          component: () => import('@/views/home/index.vue'),
        },
        {
          path: 'team',
          name: 'team',
          component: () => import('@/views/home/team.vue'),
        },
        {
          path: 'template',
          name: 'template',
          component: () => import('@/views/home/template.vue'),
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/home/settings.vue'),
        },
        {
          path: 'feedback',
          name: 'feedback',
          component: () => import('@/views/home/feedback.vue'),
        },
        {
          path: 'editor',
          name: 'editor',
          component: () => import('@/views/home/editor.vue'),
        },
        {
          path: 'upload-demo',
          name: 'upload-demo',
          component: () => import('@/views/home/upload-demo.vue'),
        },
        {
          path: 'watermark-demo',
          name: 'watermark-demo',
          component: () => import('@/views/home/watermark-demo.vue'),
        },
      ],
    },

    {
      path: '/product/:id/:tableSchemaId?',
      name: 'product',
      component: () => import('@/views/product/index.vue'),
    },
  ],
})

createRouteGuard(router)
export default router
