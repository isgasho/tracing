/* eslint-disable */
const getters = {
  //misc
  service: state => state.misc.service,

  selDate: state => state.apm.selDate,
  dashSelDate: state => state.apm.dashSelDate,
  dashNav: state => state.apm.dashNav,


  appid: state => state.apm.appid,

  //user
  token: state => state.user.token,
  userid : state => state.user.id,
  avatar: state => state.user.avatar,
  name: state => state.user.name,
  priv: state => state.user.priv
}
export default getters
 