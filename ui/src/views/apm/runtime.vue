<template>
  <div>
     <Row>
      <Col span="22" offset="1" class="no-border">
        <Table stripe  :columns="tableLabels" :data="agents" class="margin-top-40">
            <template slot-scope="{ row }" slot="host_ip">
              <div style="font-weight:bold">{{row.host_name}}</div>
              <div class="font-size-12">{{row.ip}}</div>
            </template>
             <template slot-scope="{ row }" slot="version">
              <div><span style="font-weight:bold">agent:</span> {{row.info.agentVersion}}</div>
             <div><span style="font-weight:bold">jvm:</span> {{row.info.vmVersion}}</div>
            </template>
            <template slot-scope="{ row }" slot="type">
               {{row.info.serverMetaData.serverInfo}}
            </template>
            <template slot-scope="{ row }" slot="is_live">
               <Tag color="success" type="dot" v-if="row.is_live"></Tag>
               <Tag color="warning" type="dot" v-else></Tag>
            </template>
            <template slot-scope="{ row }" slot="jvm_args">
               <Tooltip  max-width="600" placement="bottom-end" :delay="300">
                  <Icon type="ios-eye" class="color-table-detail-icon "/>
                  <div slot="content">     
                      <p v-for="arg in row.info.serverMetaData.vmArgs">{{arg}}</p>                             
                  </div>
              </Tooltip>    
             
            </template>

            <template slot-scope="{ row }" slot="operation">
              <Button size="small" type="primary" ghost @click="dashboard(row)">详细图表</Button>
            </template>
        </Table>
      </Col>
    </Row>  

     <Modal v-model="dashVisible" :footer-hide="true" width="1000px"  :styles="{top: '20px'}">
        <Row style="padding:10px">
          <Col span="12">
             <blueLineChart width="430px" height="200px" id="runtime-jvmcpu"  name2="jvm" name1="system" title="CPU Usage" :timeline="timeline" :valueList2="jvmCpuList" :valueList1="systemCpuList" :group="chartGroup"></blueLineChart>
             <div>
               <greenLineChart width="430px" height="200px" id="runtime-jvmheap" title="JVM Heap" name1="heap max" name2="heap usage" unit="(MB)" :timeline="timeline" :valueList1="heapMaxList" :valueList2="jvmHeapList" class="margin-top-20" :group="chartGroup"></greenLineChart>
              </div>
          </Col>
          <Col span="12">
            <yInverseChart width="430px" height="420px" id="runtime-fullgc" title="Full GC" name1="次数" name2="耗时" unit1="发生次数" unit2="累计耗时(ms)" :timeline="timeline" :valueList1="fullgcCountList" :valueList2="fullgcDurationList" class="margin-top-20" :group="chartGroup"></yInverseChart>
          </Col>

         
          
          
        </Row>
    </Modal>
  </div>   
</template>

<script>
import echarts from 'echarts'
import request from '@/utils/request' 
import blueLineChart from './charts/blueLineChart'
import greenLineChart from './charts/greenLineChart'
import yInverseChart from './charts/yInverseChart'
export default {
  name: 'runtime',
  components: {blueLineChart,greenLineChart,yInverseChart},
  data () {
    return {
      agents: [],
      dashVisible: false,
      tableLabels: [
            {
                title: 'Host/Ip',
                slot: 'host_ip'
            },
             {
                title: '运行状态',
                slot: 'is_live',
                width: 150
            },
            {
                title: 'Version',
                slot: 'version'
            },
            {
                title: '启动时间',
                key: 'start_time'
            },
            {
                title: '类型',
                slot: 'type'
            },
            {
                title: 'JVM参数',
                slot: 'jvm_args',
                align: 'center'
            },
             {
                title: '操作',
                slot: 'operation',
                width: 170,
                align: 'center'
            }
        ],

        timeline: [],
        jvmCpuList: [],
        systemCpuList: [],
        jvmHeapList: [],
        heapMaxList: [],
        fullgcCountList :[],
        fullgcDurationList: [],
        chartGroup: 'runtimeGroup'
    }
  },
  watch: {
    "$store.state.apm.selDate"() {
           
    },
    "$store.state.apm.appName"() {
            this.initAgents()
    }
  },
  computed: {

  },
  methods: { 
    // 加载JVM详细图表
    dashboard(r) {
      this.$Loading.start();
      // 加载当前APP的dashbord数据
      request({
          url: '/web/runtimeDash',
          method: 'GET',
          params: {
            app_name: this.$store.state.apm.appName,
            start: JSON.parse(this.$store.state.apm.selDate)[0],
            end: JSON.parse(this.$store.state.apm.selDate)[1],
            agent_id: r.agent_id
          }
      }).then(res => {   
        this.timeline =  res.data.data.timeline
        this.jvmCpuList = res.data.data.jvm_cpu_list
        this.systemCpuList = res.data.data.sys_cpu_list        
        this.jvmHeapList = res.data.data.jvm_heap_list
        this.heapMaxList = res.data.data.heap_max_list
        this.fullgcCountList = res.data.data.fullgc_count_list
        this.fullgcDurationList = res.data.data.fullgc_duration_list
         
        this.dashVisible = true
          this.$Loading.finish();
      }).catch(error => {
        this.$Loading.error();
      })
    },
    initAgents() {
      var appName = this.$store.state.apm.appName
      if (appName == '') {
        return 
      }
      request({
            url: '/web/agentList',
            method: 'GET',
            params: {
                app_name: appName
            }
        }).then(res => {   
            this.agents = res.data.data
        })
    }
  },
  mounted() {
    this.initAgents()
    echarts.connect(this.chartGroup);
  }
}
</script>


<style lang="less" scoped> 
@import "../../theme/gvar.less";
</style>

<style lang="less">
.ivu-tag {
  border: none !important;
  background: transparent !important
} 
</style>
