<template>
  <div class="app-containe padding-top-10" style="overflow-y:scroll;overflow-x:none;background-color:white">
     <Row :gutter="20">
         <Col span="12">
            <respTime width="100%" height="400px" id="apm-resp" :dateList="dateList" :valueList="elapsedList"></respTime>
         </Col>
         <Col span="11">
            <apdex width="100%" height="410px" id="apm-apdex" :dateList="dateList" :valueList="apdexList"></apdex>
         </Col>
     </Row>
     <Row :gutter="20">
        <Col span="9">
        <rpm width="100%" height="300px" id="apm-rpm" :dateList="dateList" :valueList="countList"></rpm>
              <error width="100%" height="300px" id="apm-error" :dateList="dateList" :valueList="errorList"></error>
             
         </Col>
          <Col span="6" offset="1" style="padding:8px 10px;padding-left:20px">
                <div class="font-size-18">应用动态</div>
                <div class="margin-top-10 card-tab">
                    <Button type="primary" ghost>告警通知</Button>
                    <Button >事件日志</Button>
                </div>
               <div>
                   <Icon type="ios-happy-outline" class="margin-top-20 color-primary2 margin-left-20" style="font-size:60px" />
               </div>
               <div class="margin-top-10 font-size-18 margin-left-5">
                   恭喜，当前没有任何告警
               </div>
          </Col>
          <Col span="8" style="padding:8px 10px;padding-left:20px">
                <div class="font-size-18">节点最新状态</div>
                <Table :columns="trLabels" :data="agentList"  class="margin-top-10"></Table>
          </Col>
     </Row>
  </div>
</template>

<script>
import echarts from 'echarts'
import request from '@/utils/request'
import apdex from './charts/apdex'
import error from './charts/error'
import respTime from './charts/respTime'
import rpm from './charts/rpm'
export default {
  name: 'apmDashboard',
  components: {apdex,error,respTime,rpm},
  data () {
    return {
        trLabels: [
            {
                title: '节点名',
                key: 'host_name'
            },
            {
                title: '是否容器',
                key: 'is_container',
                width: 100,
            },
            {
                title: '运行状态',
                key: 'is_live',
                width:100,
                render: (h, params) => {
                    if (params.row.is_live) {
                          return h('div', [
                                h('Tag', {
                                    props: {
                                        type: 'dot',
                                        size:"small",
                                        color: "success"
                                    }
                                })
                            ]);
                    }
                       return h('div', [
                                h('Tag', {
                                    props: {
                                        type: 'dot',
                                        size:"small",
                                        color: "warning"
                                    }
                                })
                        ]);     
                 }
            }
        ],
        dateList: [],
        countList: [],
        elapsedList: [],
        errorList: [],
        apdexList: [],

        agentList: []
    }
  },
  watch: {
    "$store.state.apm.selDate"() {
            this.initDash()
    },
    "$store.state.apm.appName"() {
            this.initDash()
    }
  },
  computed: {
  },
  methods: {
      initDash() {
        this.$Loading.start();
        // 加载当前APP的dashbord数据
        request({
            url: '/apm/web/appDash',
            method: 'GET',
            params: {
                app_name: this.$store.state.apm.appName,
                start: JSON.parse(this.$store.state.apm.selDate)[0],
                end: JSON.parse(this.$store.state.apm.selDate)[1],
            }
        }).then(res => {   
            this.dateList = res.data.data.timeline
            this.countList = res.data.data.count_list
            this.elapsedList = res.data.data.elapsed_list
            this.errorList = res.data.data.error_list
            this.apdexList = res.data.data.apdex_list

            this.$Loading.finish();
        }).catch(error => {
          this.$Loading.error();
        })

        request({
            url: '/apm/web/agentList',
            method: 'GET',
            params: {
                app_name: this.$store.state.apm.appName
            }
        }).then(res => {   
            this.agentList = res.data.data
            console.log(this.agentList)
        })
      }
  },
  mounted() {
    this.initDash()
    echarts.connect('group-dashboard');
  }
}
</script>

<style lang="less">
.ivu-tag.ivu-tag-dot {
    border:none !important;
    background: transparent !important
}
.ivu-tag.ivu-tag-dot:hover {
    background:transparent !important;
}
</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";


</style>
