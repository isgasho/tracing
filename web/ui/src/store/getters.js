/* eslint-disable */
const getters = {
  //misc
  service: state => state.misc.service,
  currentPage: state => state.misc.currentPage,

  selDate: state => state.misc.selDate,
  appid: state => state.apm.appid,

  //user
  token: state => state.user.token,
  userid : state => state.user.id,
  avatar: state => state.user.avatar,
  name: state => state.user.name,
  priv: state => state.user.priv
}
export default getters
 