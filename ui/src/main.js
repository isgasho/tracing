// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import iView from 'iview';

// 全局范围加载通用样式，每个vue page里无需重复引入
import '!style-loader!css-loader!less-loader!./theme/gcss.less'
import '!style-loader!css-loader!less-loader!./theme/mavon-md.css'

// 引入markdown editor和代码高亮
import mavonEditor from 'mavon-editor'
import 'mavon-editor/dist/css/index.css'

import hljs from 'highlight.js'
import 'highlight.js/styles/googlecode.css' //样式文件
Vue.directive('highlight', function(el) {
  let blocks = el.querySelectorAll('pre code');
  blocks.forEach((block) => {
      hljs.highlightBlock(block)
  })
})
Vue.use(mavonEditor)

Vue.config.productionTip = false

import i18n from './lang' // Internationalization

Vue.use(iView);

import store from './store'

 
router.beforeEach((to, _, next) => {
    next()
})

router.afterEach(() => {
})



/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  i18n,
  components: { App },
  template: '<App/>'
})
