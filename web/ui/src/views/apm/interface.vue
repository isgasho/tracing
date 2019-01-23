<template>
  <div>
interface
      </div>   
</template>

<script>
import request from '@/utils/request' 
export default {
  name: 'interface',
  data () {
    return {
      apiStats: []
    }
  },
  watch: {
  },
  computed: {

  },
  methods: {
  },
  mounted() {
    request({
            url: '/apm/web/apiStats',
            method: 'GET',
            params: {
                app_name: this.$store.state.apm.appName,
                start: JSON.parse(this.$store.state.apm.selDate)[0],
                end: JSON.parse(this.$store.state.apm.selDate)[1],
            }
        }).then(res => {   
            this.apiStats = res.data.data
            console.log(this.apiStats)
        })
  }
}
</script>


<style lang="less" scoped> 
@import "../../theme/gvar.less";
</style>
