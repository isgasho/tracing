/* eslint-disable */
import Cookies from 'js-cookie'
const apm = {
  state: {
    appid:  Cookies.get('apm-appid') || '',
    appName:  Cookies.get('apm-appName') || ''
  },

  mutations: {
    SET_APPID: (state, appid) => {
      state.appid = appid
      Cookies.set('apm-appid', appid)
    },
    SET_APPName: (state, name) => {
        state.appName = name
        Cookies.set('apm-appName', name)
      }
  },

  actions: {
    setAPPID({ commit }, appid) {
        commit('SET_APPID', appid)
    },
    setAPPName({ commit }, name) {
        commit('SET_APPName', name)
      }  
  }
}

export default apm