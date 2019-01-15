/* eslint-disable */
import Cookies from 'js-cookie'
const apm = {
  state: {
    appid:  Cookies.get('apm-appid') || '',
    appName:  Cookies.get('apm-appName') || '',
    selDate: Cookies.get('sel-date') || JSON.stringify([new Date((new Date()).getTime() - 3600 * 1000).toLocaleString('chinese',{hour12:false}).replace(/\//g,'-'),new Date().toLocaleString('chinese',{hour12:false}).replace(/\//g,'-')])
  },

  mutations: {
    SET_APPID: (state, appid) => {
      state.appid = appid
      Cookies.set('apm-appid', appid)
    },
    SET_APPName: (state, name) => {
        state.appName = name
        Cookies.set('apm-appName', name)
      },
      SET_SEL_DATE: (state, date) => {
        state.selDate = date
        Cookies.set('sel-date', date)
      }
  },

  actions: {
    setAPPID({ commit }, appid) {
        commit('SET_APPID', appid)
    },
    setAPPName({ commit }, name) {
        commit('SET_APPName', name)
    },
    setSelDate({ commit }, date) {
      commit('SET_SEL_DATE', date)
    }  
  }
}

export default apm