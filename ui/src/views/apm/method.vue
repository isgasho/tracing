<template>
  <div>
    <Row>
      <Col span="22" offset="1" class="no-border">
        <Table stripe :columns="methodLabels" :data="methods.slice((this.currentMethodPage-1)*10,this.currentMethodPage*10)" class="margin-top-40"  @on-sort-change="sortMethod">
        </Table>
        <Page :current="currentMethodPage" :total="methods.length" size="small" class="margin-top-15" simple  @on-change="setMethodPage"/>
      </Col>
    </Row>
  </div>   
</template>

<script>
import request from '@/utils/request' 
export default {
  name: 'method',
  data () {
    return {
        methodLabels: [
            {
                title: 'Method',
                key: 'method'
            },
            {
                title: '均耗时(ms)',
                key: 'average_elapsed',
                width:140,
                sortable: 'custom'
            },
                        {
                title: '请求数',
                key: 'count',
                width: 120,
                sortable: 'custom'
            },
            {
                title: '错误数',
                key: 'error_count',
                width: 120,
                sortable: 'custom'
            },
            {
                title: '最大耗时(ms)',
                key: 'max_elapsed',
                width: 200
            },
                        {
                title: '服务类型',
                key: 'service_type',
                width: 120
            },
        ],

      methods: [],

      currentMethodPage: 1,
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
      sortMethod(val) {
      switch (val.key) {
        case "average_elapsed": // 平均耗时排序
          if (val.order=='asc') {
            this.methods.sort(function(api1,api2){
                return api1.average_elapsed - api2.average_elapsed;
            });
          } else {
            this.methods.sort(function(api1,api2){
                return api2.average_elapsed - api1.average_elapsed;
            });
          }

          break;
        case "count":
          if (val.order=='asc') {
            this.methods.sort(function(api1,api2){
                return api1.count - api2.count;
            });
          } else {
            this.methods.sort(function(api1,api2){
                return api2.count - api1.count;
            });
          }
          break;
        case "error_count":
          if (val.order=='asc') {
            this.methods.sort(function(api1,api2){
                return api1.error_count - api2.error_count;
            });
          } else {
            this.methods.sort(function(api1,api2){
                return api2.error_count - api1.error_count;
            });
          }
          break;
        case "ratio_elapsed": 
           if (val.order=='asc') {
            this.detailApi.sort(function(api1,api2){
                return api1.ratio_elapsed - api2.ratio_elapsed;
            });
          } else {
            this.detailApi.sort(function(api1,api2){
                return api2.ratio_elapsed - api1.ratio_elapsed;
            });
          }
        default:
          break;
      }
    },
    setMethodPage(page) {
      this.currentMethodPage = page
    },
     initStats() {
       this.$Loading.start();
       request({
            url: '/apm/web/appMethods',
            method: 'GET',
            params: {
                app_name: this.$store.state.apm.appName,
                start: JSON.parse(this.$store.state.apm.selDate)[0],
                end: JSON.parse(this.$store.state.apm.selDate)[1],
            }
        }).then(res => {   
            this.methods = res.data.data
            // 初始化时，默认对平均耗时排序
            this.sortMethod({key:'average_elapsed',order:'desc'})

            this.$Loading.finish();
        }).catch(error => {
          this.$Loading.error();
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
