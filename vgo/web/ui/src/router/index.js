import Vue from 'vue'
import Router from 'vue-router'

import Nav from '@/views/nav'
import ApmNav from '@/views/apm/nav'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    { 
      path: '/', 
      component: Nav,
      redirect: '/apm/list',
      children: [
        { path: '/apm/list', component: () => import('@/views/apm/list')},
        { 
          path: '/apm/index', 
          component: ApmNav,
          redirect: '/apm/dashboard',
          children: [
            { path: '/apm/dashboard', component: () => import('@/views/apm/dashboard')},
            { path: '/apm/tracing', component: () => import('@/views/apm/tracing')},
            { path: '/apm/serviceMap', component: () => import('@/views/apm/serviceMap')},
            { path: '/apm/runtime', component: () => import('@/views/apm/runtime')},
            { path: '/apm/thread', component: () => import('@/views/apm/thread')},
            { path: '/apm/memory', component: () => import('@/views/apm/memory')},
            { path: '/apm/database', component: () => import('@/views/apm/database')},
            { path: '/apm/interface', component: () => import('@/views/apm/interface')},
            { path: '/apm/exception', component: () => import('@/views/apm/exception')}
          ]
        },
      ]
    },

    

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
