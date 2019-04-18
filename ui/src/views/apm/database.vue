<template>
  <div>
    <Row>
      <Col span="22" offset="1" class="no-border">
        <Table stripe  :columns="sqlLabels" :data="sqlList.slice((this.currentPage-1)*10,this.currentPage*10)" class="margin-top-40" @on-sort-change="sortSql">
            <template slot-scope="{ row }" slot="sql">
              <Tooltip :content="row.sql" max-width="400">
                  {{row.sql}}
              </Tooltip>     
            </template>
        </Table>

        <Page :current="currentPage" :total="sqlList.length" size="small" class="margin-top-15" simple  @on-change="setApiPage"/>
      </Col>
    </Row>
  </div>   
</template>

<script>
import request from '@/utils/request' 
export default {
  name: 'database',
  data () {
    return {
      sqlList: [],
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
  }
}
</script>


<style lang="less" scoped> 
@import "../../theme/gvar.less";
</style>
