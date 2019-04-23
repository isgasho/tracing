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

     <Modal v-model="dashVisible" :footer-hide="true">
        <div class="padding-left-20">
          <blueLineChart width="430px" height="200px" id="runtime-jvmcpu" title="JVM CPU Usage" :timeline="timeline" :valueList="jvmCpuList"></blueLineChart>
          <blueLineChart width="430px" height="200px" id="runtime-jvmheap" title="JVM Heap Usage" unit="(MB)" :timeline="timeline" :valueList="jvmHeapList" class="margin-top-20"></blueLineChart>
        </div>
    </Modal>
  </div>   
</template>

<script>
import request from '@/utils/request' 
import blueLineChart from './charts/blueLineChart'
export default {
  name: 'runtime',
  components: {blueLineChart},
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
        jvmHeapList: []
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
          url: '/apm/web/runtimeDash',
          method: 'GET',
          params: {
            app_name: this.$store.state.apm.appName,
            start: JSON.parse(this.$store.state.apm.selDate)[0],
            end: JSON.parse(this.$store.state.apm.selDate)[1],
            agent_id: r.agent_id
          }
      }).then(res => {   
        console.log('runtime dash',res.data.data)
        this.timeline =  res.data.data.timeline
        this.jvmCpuList = res.data.data.jvm_cpu_list        
        this.jvmHeapList = res.data.data.jvm_heap_list

       
        this.dashVisible = true
          this.$Loading.finish();
      }).catch(error => {
        this.$Loading.error();
      })
    },
    initAgents() {
      request({
            url: '/apm/web/agentList',
            method: 'GET',
            params: {
                app_name: this.$store.state.apm.appName
            }
        }).then(res => {   
            this.agents = res.data.data
        })
    }
  },
  mounted() {
    this.initAgents()
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
