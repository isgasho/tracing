<template>
  <div>
      <topology :graphData="JSON.parse(data)" style=' width:100%; min-height: 90vh;max-height:200vh'></topology>
  </div>   
</template>

<script>
import request from '@/utils/request' 
import topology from './charts/topology'
let $ = window.go.GraphObject.make
export default {
  name: 'serviceMap',
  components: {topology},
  data () {
    return {
      data : null
    }
  },
  watch: {
  },
  computed: {
 
  },
  methods: {
  },
  mounted() {
    this.$Loading.start();
    request({
        url: '/apm/web/serviceMap',
        method: 'GET',
        params: {
        }
    }).then(res => {
      this.data = res.data.data
      this.$Loading.finish();
    }).catch(error => {
        this.$Loading.error();
      })
  }
}
</script>


<style lang="less" scoped> 
@import "../../theme/gvar.less";
</style>
