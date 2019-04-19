<template>
  <div>
    <Row>
      <Col span="22" offset="1" class="no-border">
        <Table stripe  :columns="sqlLabels" :data="sqlList.slice((this.currentPage-1)*10,this.currentPage*10)" class="margin-top-40" @on-sort-change="sortSql">
            <template slot-scope="{ row }" slot="sql">
              <Tooltip  max-width="400" :delay="400" placement="top-start">
                  {{row.sql}}
                  <div slot="content">                                  
                    <!-- <div><Tag style="background: #F28F20;border:none;" size="medium" class="margin-right-10">Method</Tag>{{row.method}}</div> -->
                    <div><Tag style="background: #F28F20;border:none;" size="medium"  class="margin-right-10">完整SQL</Tag> {{row.sql}}</div>
                  </div>
              </Tooltip>     
            </template>

             <template slot-scope="{ row }" slot="operation">
              <Button size="small" type="primary" ghost @click="dashboard(row)">详细图表</Button>
            </template>
        </Table>

        <Page :current="currentPage" :total="sqlList.length" size="small" class="margin-top-15" simple  @on-change="setApiPage"/>
      </Col>
    </Row>


    <Modal v-model="dashVisible" :footer-hide="true">
             <rpm width="430px" height="200px" id="apm-rpm" :dateList="dateList" :valueList="countList"></rpm>
             <error width="430px" height="200px" id="apm-error" :dateList="dateList" :valueList="errorList" class="margin-top-20"></error>
            <respTime width="430px" height="200px" id="apm-resp" :dateList="dateList" :valueList="elapsedList" class="margin-top-20"></respTime>
    </Modal>
  </div>   
</template>

<script>
import request from '@/utils/request' 
import echarts from 'echarts'
import respTime from './charts/respTime'
import rpm from './charts/rpm'
import error from './charts/error'
export default {
  name: 'database',
  components: {error,respTime,rpm},
  data () {
    return {
      sqlList: [],
      dashVisible: false,
      dateList: [],
      countList: [],
      elapsedList: [],
      errorList: [],
      sqlLabels: [
            {
                title: 'SQL',
                slot: 'sql',
                ellipsis : true
            },
            {
                title: '均耗时(ms)',
                key: 'average_elapsed',
                width:130,
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
                width: 120
            },
            {
                title: '最小耗时(ms)',
                key: 'min_elapsed',
                width: 120
            },
             {
                title: '操作',
                slot: 'operation',
                width: 170,
                align: 'center'
            }
        ],
      currentPage : 1,
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
    dashboard(r) {
      this.$Loading.start();
      // 加载当前APP的dashbord数据
      request({
          url: '/apm/web/sqlDash',
          method: 'GET',
          params: {
              app_name: this.$store.state.apm.appName,
              start: JSON.parse(this.$store.state.apm.selDate)[0],
              end: JSON.parse(this.$store.state.apm.selDate)[1],
              sql_id: r.id
          }
      }).then(res => {   
          this.dateList = res.data.data.timeline
          this.countList = res.data.data.count_list
          this.elapsedList = res.data.data.elapsed_list
          this.errorList = res.data.data.error_list

          this.dashVisible = true
          console.log(res.data.data)
          this.$Loading.finish();
      }).catch(error => {
        this.$Loading.error();
      })
    },
     sortSql(val) {
      switch (val.key) {
        case "average_elapsed": // 平均耗时排序
          if (val.order=='asc') {
            this.sqlList.sort(function(api1,api2){
                return api1.average_elapsed - api2.average_elapsed;
            });
          } else {
            this.sqlList.sort(function(api1,api2){
                return api2.average_elapsed - api1.average_elapsed;
            });
          }
          break;
        case "count":
          if (val.order=='asc') {
            this.sqlList.sort(function(api1,api2){
                return api1.count - api2.count;
            });
          } else {
            this.sqlList.sort(function(api1,api2){
                return api2.count - api1.count;
            });
          }
          break;
        case "error_count":
          if (val.order=='asc') {
            this.sqlList.sort(function(api1,api2){
                return api1.error_count - api2.error_count;
            });
          } else {
            this.sqlList.sort(function(api1,api2){
                return api2.error_count - api1.error_count;
            });
          }
          break;
        default:
          break;
      }
    },
     initStats() {
       this.$Loading.start();
       request({
            url: '/apm/web/sqlStats',
            method: 'GET',
            params: {
                app_name: this.$store.state.apm.appName,
                start: JSON.parse(this.$store.state.apm.selDate)[0],
                end: JSON.parse(this.$store.state.apm.selDate)[1],
            }
        }).then(res => {   
            this.sqlList = res.data.data
            console.log(this.sqlList)
            // 初始化时，默认对平均耗时排序
            this.sortSql({key:'average_elapsed',order:'desc'})

            this.$Loading.finish();
        }).catch(error => {
          this.$Loading.error();
        })
    }
  },
  mounted() {
    this.initStats()
     echarts.connect('group-dashboard');
  }
}
</script>


<style lang="less" scoped> 
@import "../../theme/gvar.less";
</style>
