<template>
  <div>
    <Row>
      <Col span="22" offset="1" class="no-border">
        <Table stripe :columns="apiLabels" :data="apiStats" class="margin-top-20" @on-row-click="apiDetail" >
        </Table>
      </Col>
    </Row>
  </div>   
</template>

<script>
import request from '@/utils/request' 
export default {
  name: 'interface',
  data () {
    return {
      apiStats: [],
      apiLabels: [
            {
                title: 'API',
                key: 'url',
                width: 400
            },
            {
                title: '平均耗时(ms)',
                key: 'average_elapsed',
                width:170
            },
            {
                title: '请求次数',
                key: 'count',
                width: 170
            },
            {
                title: '错误次数',
                key: 'error_count',
                width: 170
            },
             {
                title: '最大耗时(ms)',
                key: 'max_elapsed',
                width: 170
            },
            {
                title: '最小耗时(ms)',
                key: 'min_elapsed',
                width: 170
            },
        ],
      detailApi: {}
    }
  },
   watch: {
    "$store.state.apm.selDate"() {
            this.initStats()
    },
    "$store.state.apm.appName"() {
            this.initStats()
    }
  },
  computed: {

  },
  methods: {
    apiDetail(api) {
      request({
            url: '/apm/web/apiDetail',
            method: 'GET',
            params: {
                app_name: this.$store.state.apm.appName,
                url: api.url,
                start: JSON.parse(this.$store.state.apm.selDate)[0],
                end: JSON.parse(this.$store.state.apm.selDate)[1]
            }
        }).then(res => {   
            this.detailApi = res.data.data
            console.log(this.detailApi)
        })
    },
    initStats() {
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
  },
  mounted() {
    this.initStats()
  }
}
</script>


<style lang="less" scoped> 
@import "../../theme/gvar.less";
</style>
