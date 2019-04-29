/* eslint-disable */
import Cookies from 'js-cookie'
const apm = {
  state: {
    appid:  Cookies.get('apm-appid') || '',
    appName:  Cookies.get('apm-appName') || '', // 当前选择的APP name
    selDate:  getDate(), // APM内的日历
    dashSelDate: Cookies.get('apm-dash-selDate') ||'10', // 首页过去X分钟选项
    dashNav:  Cookies.get('apm-dash-nav') || '1', // 首页显示应用地图还是应用列表
    errorFilterNav: Cookies.get('apm-error-filter-nav') || ''// 首页应用地图的错误过滤 
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
    },
    SET_DASH_SELDATE: (state, date) => {
      state.dashSelDate = date
      Cookies.set('apm-dash-selDate', date)
    },
    SET_DASH_NAV: (state, val) => {
      state.dashNav = val
      Cookies.set('apm-dash-nav', val)
    },
    SET_ERROR_FILTER_NAV: (state, val) => {
      state.errorFilterNav = val
      Cookies.set('apm-error-filter-nav', val)
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
    } ,
    setDashSelDate({ commit }, date) {
      commit('SET_DASH_SELDATE', date)
    } ,
    setDashNav({ commit }, val) {
      commit('SET_DASH_NAV', val)
    } ,
    setErrorFilterNav({ commit }, val) {
      commit('SET_ERROR_FILTER_NAV', val)
    }     
  }
}

function getDate() {
  var d = Cookies.get('sel-date')
  if (d == '' || d == '[]' || d == undefined) {
    return JSON.stringify([new Date((new Date()).getTime() - 3600 * 1000).toLocaleString('chinese',{hour12:false}).replace(/\//g,'-'),new Date().toLocaleString('chinese',{hour12:false}).replace(/\//g,'-')])
  }

  return d
}
export default apm