<template>
  <div class="callback">
  </div>
</template>

<script>
/* eslint-disable */
import request from '@/utils/request' 
import Cookies from 'js-cookie'
export default {
  name: 'callback',
  data() {
      return {}
  },
  created() {
    var subToken = this.$route.query.subToken 
    // 获取用户信息s
    request({
        url: '/apm/web/login',
        method: 'POST', 
        params: {
            subToken: subToken,
        }
    }).then(res => {
        // 存在，设置用户信息
        this.$store.dispatch('SetUserInfo', res.data.data).then(() => {  
            // 获取历史路径
            var opath = Cookies.get('lastPath')
            console.log(opath)
            if (opath != '' && opath != undefined) {
              Cookies.remove('lastPath')
               this.$router.push({ path: opath })
            } else {
              this.$router.push('/apm/ui/list')
            }
           
        })
    }).catch(error => {
      console.log(error)
      _this.$router.push({ path: '/' })
    })
  },
  destroyed() {
    // window.removeEventListener('hashchange', this.afterQRScan)
  }
}
</script>