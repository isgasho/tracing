import Vue from 'vue'
import Router from 'vue-router'

import Nav from '@/views/nav'
import ApmNav from '@/views/apm/nav'
import AlertsNav from '@/views/alerts/nav'

import DocNav from '@/views/docs/nav'
Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    { 
      path: '/ui', 
      component: Nav,
      redirect: '/ui/dashboard',
      children: [
        // 首页大盘
        { path: '/ui/dashboard', meta: '监控大盘',bg: '#00B6D8',component: () => import('@/views/dashboard/index')},
        // 管理员页面
        { path: '/ui/admin', meta:'管理面板', bg: '#39c',component: () => import('@/views/user/admin')},
        // 用户个人设置
        { path: '/ui/user/setting', meta:'个人设置', bg: '#39c',component: () => import('@/views/user/setting')},
        // 应用监控
        { 
          path: '/ui/apm', 
          component: ApmNav,
          redirect: '/ui/apm/dashboard',
          meta: '应用监控',
          children: [
            { path: '/ui/apm/dashboard',     meta: '应用监控', bg: '#39c',component: () => import('@/views/apm/dashboard')},
            { path: '/ui/apm/tracing',  meta: '应用监控',bg: '#39c', component: () => import('@/views/apm/tracing')},
            { path: '/ui/apm/runtime',  meta: '应用监控',bg: '#39c', component: () => import('@/views/apm/runtime')},
            { path: '/ui/apm/system',  meta: '应用监控',bg: '#39c', component: () => import('@/views/apm/system')},
            { path: '/ui/apm/thread',  meta: '应用监控',bg: '#39c', component: () => import('@/views/apm/thread')},
            { path: '/ui/apm/memory',  meta: '应用监控', bg: '#39c',component: () => import('@/views/apm/memory')},
            { path: '/ui/apm/database',  meta: '应用监控',bg: '#39c', component: () => import('@/views/apm/database')},
            { path: '/ui/apm/api', meta: '应用监控',bg: '#39c',  component: () => import('@/views/apm/api')},
            { path: '/ui/apm/exception',  meta: '应用监控', bg: '#39c',component: () => import('@/views/apm/exception')},
            { path: '/ui/apm/method',  meta: '应用监控', bg: '#39c',component: () => import('@/views/apm/method')},
            { path: '/ui/apm/serviceMap',  meta: '应用监控', bg: '#39c',component: () => import('@/views/apm/serviceMap')}
          ]
        },
        // 告警平台
        { 
          path: '/ui/alerts', 
          component: AlertsNav,
          redirect: '/ui/alerts/appList',
          meta: '应用监控',
          children: [
            { path: '/ui/alerts/appList',meta: '告警平台',bg:'#b286bc', component: () => import('@/views/alerts/appList')},
            { path: '/ui/alerts/policy',meta: '告警平台', bg:'#b286bc',component: () => import('@/views/alerts/policy')},
            { path: '/ui/alerts/group',meta: '告警平台', bg:'#b286bc',component: () => import('@/views/alerts/group')},
            { path: '/ui/alerts/alertsNotify',meta: '告警平台',bg:'#b286bc', component: () => import('@/views/alerts/alertsNotify')}
          ]
        },
      ]
    },
    { path: '/', component: () => import('@/views/index')},
    { path: '/ui/callback', component: () => import('@/views/login/callback')},

    //帮助文档
    { 
      path: '/ui/docs',
      component: DocNav,
      redirect: '/ui/docs/about',
      children: [
        { path: 'about', component: () => import('@/views/docs/pages/about')},

        { path: 'install', component: () => import('@/views/docs/pages/install')}
      ]
    },
    { path: '/404', component: () => import('@/views/misc/page404')},
    { path: '*', redirect: '/404'}
  ]
})
