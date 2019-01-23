import Vue from 'vue'
import Router from 'vue-router'

import Nav from '@/views/nav'
import ApmNav from '@/views/apm/nav'
import AlertsNav from '@/views/alerts/nav'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    { 
      path: '/', 
      component: Nav,
      redirect: '/apm/ui/list',
      children: [
        { path: '/apm/ui/list', meta: '应用监控',bg: '#348899',component: () => import('@/views/apm/list')},
        { path: '/apm/ui/admin', meta:'管理面板', bg: '#348899',component: () => import('@/views/admin')},
        { path: '/apm/ui/person', meta:'个人设置', bg: '#348899',component: () => import('@/views/personSetting')},
        { 
          path: '/apm/ui/index', 
          component: ApmNav,
          redirect: '/apm/ui/dashboard',
          meta: '应用监控',
          children: [
            { path: '/apm/ui/dashboard',     meta: '应用监控', bg: '#348899',component: () => import('@/views/apm/dashboard')},
            { path: '/apm/ui/tracing',  meta: '应用监控',bg: '#348899', component: () => import('@/views/apm/tracing')},
            { path: '/apm/ui/serviceMap',  meta: '应用监控', bg: '#348899',component: () => import('@/views/apm/serviceMap')},
            { path: '/apm/ui/runtime',  meta: '应用监控',bg: '#348899', component: () => import('@/views/apm/runtime')},
            { path: '/apm/ui/thread',  meta: '应用监控',bg: '#348899', component: () => import('@/views/apm/thread')},
            { path: '/apm/ui/memory',  meta: '应用监控', bg: '#348899',component: () => import('@/views/apm/memory')},
            { path: '/apm/ui/database',  meta: '应用监控',bg: '#348899', component: () => import('@/views/apm/database')},
            { path: '/apm/ui/interface', meta: '应用监控',bg: '#348899',  component: () => import('@/views/apm/interface')},
            { path: '/apm/ui/exception',  meta: '应用监控', bg: '#348899',component: () => import('@/views/apm/exception')}
          ]
        },
        { 
          path: '/apm/ui/alerts', 
          component: AlertsNav,
          redirect: '/apm/ui/alerts/appList',
          meta: '应用监控',
          children: [
            { path: '/apm/ui/alerts/appList',meta: '告警平台',bg:'#b286bc', component: () => import('@/views/alerts/appList')},
            { path: '/apm/ui/alerts/policy',meta: '告警平台', bg:'#b286bc',component: () => import('@/views/alerts/policy')},
            { path: '/apm/ui/alerts/group',meta: '告警平台', bg:'#b286bc',component: () => import('@/views/alerts/group')},
            { path: '/apm/ui/alerts/alertsNotify',meta: '告警平台',bg:'#b286bc', component: () => import('@/views/alerts/alertsNotify')}
          ]
        },
      ]
    },
    { path: '/apm/ui/login', component: () => import('@/views/login/index')},
    { path: '/apm/ui/callback', component: () => import('@/views/login/callback')},


    { path: '/404', component: () => import('@/views/errorPage/page404')},
    { path: '*', redirect: '/404'}
  ]
})



export const routerMap = [
  // { path: '/404', component: () => import('@/views/errorPage/404'), hidden: true },
  { path: '/', component: () => import('@/views/template'), hidden: true },
  // {
  //   path: '/infra',
  //   component: Layout,
  //   redirect: '/infra/service', 
  //   alwaysShow: true, // will always show the root menu
  //   meta: { title: 'Infra',icon: 'component'},
  //   children: [{
  //     path: 'service',
  //     component: () => import('@/views/infra/service'),
  //     name: 'service',
  //     meta: {
  //       title: 'Service'
  //     }
  //   }, 
  //   {
  //     path: 'application',
  //     component: () => import('@/views/infra/application'),
  //     name: 'application',
  //     meta: {
  //       title: 'Application'
  //     }
  //   },
  //   {
  //     path: 'server',
  //     component: () => import('@/views/infra/server'),
  //     name: 'server',
  //     meta: {
  //       title: 'Cloud Server'
  //     }
  //   }
  // ]
  // },
  { path: '*', redirect: '/404', hidden: true }
]
