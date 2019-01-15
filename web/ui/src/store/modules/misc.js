import Cookies from 'js-cookie'

const misc = {
  state: {
    service: Cookies.get('sel-service') || 'empty',
    currentPage: Cookies.get('current-page') || '应用监控'
  },
  mutations: {
    SET_SERVICE: (state, service) => {
      state.service = service
      Cookies.set('sel-service', service)
    },
    SET_PAGE: (state, page) => {
      state.currentPage = page
      Cookies.set('current-page', page)
    }
  },
  actions: {
    setService({ commit }, service) {
      commit('SET_SERVICE', service)
    },
    setPage({ commit }, page) {
      commit('SET_PAGE', page)
    }
  }
}

export default misc
