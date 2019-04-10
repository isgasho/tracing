import Vue from 'vue'
import Vuex from 'vuex'
import misc from './modules/misc'
import user from './modules/user'
import apm from './modules/apm'
import getters from './getters'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    misc,
    user,
    apm
  },
  getters
})

export default store
 