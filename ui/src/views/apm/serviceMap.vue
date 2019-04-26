<template>
  <div style="width: 100%;max-width:100%">
     <topology style="width:calc(100vw - 180px);height:calc(100vh - 100px)"></topology>
  </div>   
</template>

<script>
import request from '@/utils/request' 
import topology from './charts/topology'
export default {
  name: 'serviceMap',
  components: {topology},
  data () {
    return {
      data : null
    }
  },
  watch: {
    "$store.state.apm.selDate"() {
        this.initServiceMap()
    },
    "$store.state.apm.appName"() {
        this.initServiceMap()
    }
  },
  computed: {
 
  },
  methods: {
    initServiceMap() {
      this.$Loading.start();
        request({
            url: '/apm/web/appServiceMap',
            method: 'GET',
            params: {
              app_name: this.$store.state.apm.appName,
              start: JSON.parse(this.$store.state.apm.selDate)[0],
              end: JSON.parse(this.$store.state.apm.selDate)[1],
            }
        }).then(res => {
          this.data = res.data.data
          this.$Loading.finish();
        }).catch(error => {
            this.$Loading.error();
          })
      }
  },
  mounted() {
    this.initServiceMap()
  }
}
</script>


<style lang="less" scoped> 
@import "../../theme/gvar.less";
</style>
