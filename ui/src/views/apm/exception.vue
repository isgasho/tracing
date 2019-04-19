<template>
  <div>
    <Row>
      <Col span="22" offset="1" class="no-border">
        <Table stripe :columns="tableLabels" :data="tableData.slice((this.currentPage-1)*10,this.currentPage*10)" class="margin-top-40"  @on-sort-change="sortTable">
             <template slot-scope="{ row }" slot="method">
              <Tooltip  max-width="400" placement="top-start" :delay="400">
                  <div>{{row.method}}</div>
                  <div slot="content">                                  
                    <!-- <div><Tag style="background: #F28F20;border:none;" size="medium" class="margin-right-10">Method</Tag>{{row.method}}</div> -->
                    <div><Tag style="background: #F28F20;border:none;" size="medium"  class="margin-right-10">Class</Tag> {{row.class}}</div>
                  </div>
              </Tooltip>     
            </template>

             <template slot-scope="{ row }" slot="operation">
              <Button size="small" type="primary" ghost @click="dashboard(row)">详细图表</Button>
            </template>
        </Table>
        <Page :current="currentPage" :total="tableData.length" size="small" class="margin-top-15" simple  @on-change="setTablePage"/>
      </Col>
    </Row>

     <Modal v-model="dashVisible" :footer-hide="true">
             <rpm width="430px" height="200px" id="apm-rpm" :dateList="dateList" :valueList="countList"></rpm>
            <respTime width="430px" height="200px" id="apm-resp" :dateList="dateList" :valueList="elapsedList" class="margin-top-20"></respTime>
    </Modal>
  </div>   
</template>

<script>
import request from '@/utils/request'
import echarts from 'echarts'
import respTime from './charts/respTime'
import rpm from './charts/rpm'
export default {
  name: 'exception',
  components: {respTime,rpm},
  data () {
    return {
        tableLabels: [
             {
                title: 'Exception',
                key: 'exception'
            },
            {
                title: 'Method',
                slot: 'method'
            },
            {
                title: '均耗时(ms)',
                key: 'average_elapsed',
                width:130,
                sortable: 'custom'
            },
            {
                title: '次数',
                key: 'count',
                width: 100,
                sortable: 'custom'
            },
            {
                title: '最大耗时(ms)',
                key: 'max_elapsed',
                width: 140
            },
                        {
                title: '服务类型',
                key: 'service_type',
                width: 180
            },
            {
                title: '操作',
                slot: 'operation',
                width: 170,
                align: 'center'
            }
        ],

      tableData: [],

      currentPage: 1,

      dateList: [],
      countList: [],
      elapsedList: [],
       dashVisible: false
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
     // 加载api的详细图表
     dashboard(r) {
      this.$Loading.start();
      // 加载当前APP的dashbord数据
      request({
          url: '/apm/web/exceptionDash',
          method: 'GET',
          params: {
            app_name: this.$store.state.apm.appName,
            start: JSON.parse(this.$store.state.apm.selDate)[0],
            end: JSON.parse(this.$store.state.apm.selDate)[1],
            exception_id: r.id
          }
      }).then(res => {   
          this.dateList = res.data.data.timeline
          this.countList = res.data.data.count_list
          this.elapsedList = res.data.data.elapsed_list


          this.dashVisible = true
          console.log(res.data.data)
          this.$Loading.finish();
      }).catch(error => {
        this.$Loading.error();
      })
    },
      sortTable(val) {
      switch (val.key) {
        case "average_elapsed": // 平均耗时排序
          if (val.order=='asc') {
            this.tableData.sort(function(d1,d2){
                return d1.average_elapsed - d2.average_elapsed;
            });
          } else {
            this.tableData.sort(function(d1,d2){
                return d1.average_elapsed - d2.average_elapsed;
            });
          }

          break;
        case "count":
          if (val.order=='asc') {
            this.tableData.sort(function(d1,d2){
                return d1.count - d2.count;
            });
          } else {
            this.tableData.sort(function(d1,d2){
                return d1.count - d2.count;
            });
          }
          break;
        case "error_count":
          if (val.order=='asc') {
            this.tableData.sort(function(d1,d2){
                return d1.error_count - d2.error_count;
            });
          } else {
            this.tableData.sort(function(d1,d2){
                return d1.error_count - d2.error_count;
            });
          }
          break;
        default:
          break;
      }
    },
    setTablePage(page) {
      this.currentPage = page
    },
     initStats() {
       this.$Loading.start();
       request({
            url: '/apm/web/appException',
            method: 'GET',
            params: {
                app_name: this.$store.state.apm.appName,
                start: JSON.parse(this.$store.state.apm.selDate)[0],
                end: JSON.parse(this.$store.state.apm.selDate)[1],
            }
        }).then(res => {   
            this.tableData = res.data.data
            // 初始化时，默认对平均耗时排序
            this.sortTable({key:'count',order:'desc'})

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
