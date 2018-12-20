<template>
  <div>
      <topology :graphData="JSON.parse(testData)" style=' width:100%; min-height: 90vh;max-height:200vh'></topology>
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
      testData : `{"value":{"agentNodeModels":[{"key":"10","text":"sms","source":"TOMCAT","instanceCount":7,"agentSources":{}},{"key":"1","text":"ehuodiApi","source":"TOMCAT","instanceCount":3,"agentSources":{}},{"key":"2","text":"passport","source":"TOMCAT","instanceCount":4,"agentSources":{}},{"key":"3","text":"partyGetway","source":"TOMCAT","instanceCount":6,"agentSources":{}},{"key":"4","text":"user","source":"USER","instanceCount":1,"agentSources":{}},{"key":"5","text":"invoiceApi","source":"TOMCAT","instanceCount":3,"agentSources":{}},{"key":"6","text":"partyOpenApi","source":"TOMCAT","instanceCount":3,"agentSources":{}},{"key":"7","text":"pushServiceApi","source":"TOMCAT","instanceCount":3,"agentSources":{}},{"key":"8","text":"site","source":"TOMCAT","instanceCount":3,"agentSources":{}},{"key":"9","text":"myportal","source":"TOMCAT","instanceCount":3,"agentSources":{}},{"key":"11","text":"120.55.197.77:1210(1) 36.26.68.39:7862(21) 116.62.244.37(1) ","source":"UNKNOWN_GROUP","instanceCount":1,"agentSources":{"116.62.244.37":{"badCount":0,"errorCount":0,"goodCount":1,"warnCount":0},"120.55.197.77:1210":{"badCount":0,"errorCount":0,"goodCount":1,"warnCount":0},"36.26.68.39:7862":{"badCount":0,"errorCount":0,"goodCount":21,"warnCount":0}}},{"key":"14","text":"REDIS","source":"REDIS","instanceCount":1,"agentSources":{}},{"key":"15","text":"sms","source":"MYSQL","instanceCount":1,"agentSources":{}},{"key":"16","text":"voice.3tong.net(72) 10.33.41.150:5381(62) www.dh3t.com(2) ","source":"UNKNOWN_GROUP","instanceCount":1,"agentSources":{"10.33.41.150:5381":{"badCount":0,"errorCount":0,"goodCount":62,"warnCount":0},"voice.3tong.net":{"badCount":0,"errorCount":0,"goodCount":72,"warnCount":0},"www.dh3t.com":{"badCount":0,"errorCount":0,"goodCount":2,"warnCount":0}}}],"agentLinkModels":[{"from":"1","to":"10","totalCount":101},{"from":"2","to":"10","totalCount":12},{"from":"3","to":"10","totalCount":1},{"from":"4","to":"10","totalCount":159},{"from":"5","to":"10","totalCount":6},{"from":"6","to":"10","totalCount":40},{"from":"7","to":"10","totalCount":8},{"from":"8","to":"10","totalCount":39},{"from":"9","to":"10","totalCount":2},{"from":"10","to":"11","totalCount":23},{"from":"10","to":"14","totalCount":6818},{"from":"10","to":"15","totalCount":885},{"from":"10","to":"16","totalCount":136}],"nowNodeKey":"10"},"success":true,"message":"","code":0}`
    }
  },
  watch: {
  },
  computed: {
 
  },
  methods: {
  },
  mounted() {
  }
}
</script>


<style lang="less" scoped> 
@import "../../theme/gvar.less";
</style>
