<template>
  <div>
    <Row class="padding-20">
       <Col span="7" style="background:rgb(248, 248, 248)"  class="padding-20 padding-bottom-10">
          <div  class="font-size-18">链路过滤  <Button ghost type="primary" style="float:right;margin-top:3px;padding-left:15px;padding-right:15px;" size="small" @click="queryTraces">查询</Button></div>
          <div class="margin-top-10 no-border">
            <div class="font-size-12">请求API</div>
            <Select v-model="currentApi" style="width:100%" size="large" class="api-filter"  placeholder="默认选择全部API" filterable clearable>
              <Option v-for="api in apis" :value="api" :key="api">
                {{ api }}
              </Option>
            </Select>
          </div>
          <div class="margin-top-10 no-border">
            <div class="font-size-12">客户端地址(remote addr)</div>
            <Input class="margin-top-5" v-model="remoteAddr" placeholder="e.g.  10.0.0.1" style="width: 100%;border:none;" size="large" />
          </div>
          <Row class="margin-top-15 no-border" style="border-bottom: 1px solid #ddd;padding-bottom:25px">
            <Col span="11">
               <div class="font-size-12">最小耗时(ms)</div>
              <Input class="margin-top-5" v-model="minElapsed" placeholder="e.g.  100" style="width: 100%;border:none;" size="large" />
            </Col>
            <Col span="11" offset="2">
              <div class="font-size-12">最大耗时(ms)</div>
              <Input  class="margin-top-5" v-model="maxElapsed" placeholder="留空代表不限制" style="width: 100%" size="large" />
            </Col>
          </Row>
          <Row class="margin-top-15 no-border" style="border-bottom: 1px solid #ddd;padding-bottom:25px">
            <Col span="11">
              <div class="font-size-12">限定搜索数目</div>
              <InputNumber class="margin-top-5" :max="1000" :min="1"  :step="5" v-model="resultLimit" style="width: 100px;" size="medium"></InputNumber>
            </Col>
            <Col span="10" offset="2">
              <div class="font-size-12">只搜索错误</div>
              <i-switch v-model="searchError" class="margin-top-10 margin-left-5"/>
            </Col>
          </Row>

          <div class="margin-top-15 no-border">
            <div class="font-size-12">指定链路ID查询</div>
            <Input  class="margin-top-5" v-model="searchTraceID" placeholder="不为空时，优先于其他条件查询" style="width: 100%" size="large" />
          </div>


          <div class="margin-top-20">
           
          </div>
       </Col>
       <Col span="15" offset="1">
        <trace :graphData="tracesData" v-if="tracesData.is_suc==true" style='height: 500px' @selTraces="selectTraces"></trace>
       </Col>
    </Row>

     <div v-show="selectedTraces.length > 0" class="margin-top-10" style="padding-left:20px;padding-right:20px;">
          <span>
            <Tag style="background: rgb(18, 147, 154,.5)" size="large" @click="setTableFilter('suc')">{{selSucCount}} 成功</Tag> 
            <Tag style="background:rgba(223, 83, 83, .5)" size="large" @click="setTableFilter('error')">{{selErrCount}} 错误</Tag>
          </span>
          <!-- <span style="float:right" class="no-border">
            <Select style="width:200px"  value="1">
                <Option value="1">时间最近</Option>
                <Option value="2">时间最远</Option>
                <Option value="3">链路最短</Option>
                <Option value="3">链路最长</Option>
            </Select>
          </span> -->

        <div style="padding-left:10px;padding-right:10px" class="margin-top-20 margin-bottom-40">
          <Table :row-class-name="rowClassName"  :columns="traceLabels" :data="showTraceTable()" class="margin-top-15" @on-row-dblclick="showTrace">
              <template slot-scope="{ row }" slot="operation">
                    <Icon type="ios-eye" style="color: #777" @click="showTrace(row)"/>
              </template>
          </Table>
        </div>
      </div>

    <!-- 链路展示Modal -->
    <Modal v-model="traceVisible" :footer-hide="true" :z-index="500" fullscreen>
        <Row slot="header" style="padding-top:0px;padding-bottom:0px;border-bottom:none;height:80px;">
          <Col span="12">
             <div class="font-size-16 margin-top-30" style="font-weight:bold">{{$store.state.apm.appName}} : {{selectedTrace.api}}</div>
            <div  style="margin-top:13px;font-weight:bold;font-size:12px">
              <span class="meta-word">
                时间:
              </span>
              {{selectedTrace.showTime}}

              <span class="meta-word margin-left-10">
                耗时:
              </span>
              {{selectedTrace.elapsed}}ms

              <span class="meta-word margin-left-10">
                链路ID:
              </span>
              {{selectedTrace.traceId}}

              <span class="meta-word margin-left-10">
                服务器ID:
              </span>
              {{selectedTrace.agentId}}
              
            </div>
          </Col>
          <Col span="12" style="margin-top:-8px">
            <Row>
                <Col span="8">
                  <blueLineChart width="240px" height="130px" id="tracing-jvmcpu" :titleFontSize="12" name2="jvm" name1="system" title="cpu" :timeline="timeline" :valueList2="jvmCpuList" :valueList1="systemCpuList" :group="chartGroup" :showXAxis="false"></blueLineChart>
                </Col>
                 <Col span="8">
                   <greenLineChart width="240px" height="130px" id="tracing-jvmheap" :titleFontSize="12" title="heap" name1="heap max" name2="heap usage"  :timeline="timeline" :valueList1="heapMaxList" :valueList2="jvmHeapList"  :group="chartGroup" :showXAxis="false"></greenLineChart>
                </Col>
                <Col span="8">
                   <redLineChart width="240px" height="130px" id="tracing-fullgc" :titleFontSize="12" title="fullgc"  name2="总耗时" :timeline="timeline" :valueList2="fullgcDurationList"  :group="chartGroup" :showXAxis="false"></redLineChart>
                </Col>
            </Row>
            
            
          </Col>
        </Row>

        <div class="trace-pane" style="padding-bottom:50px;padding-top:27px">
          <Row>
            <Col span="13" class="title split-border">方法(点击具体方法名可查看详情)</Col>
            <Col span="4"  class="title">参数</Col>
            <Col span="2" class="title">耗时(ms)</Col>
            <Col span="3" class="title">服务类型</Col>
            <Col span="2" class="title">所属应用</Col>
          </Row>
          <div  class="body">
            <Row v-for="r in traceData" v-show="isShow(r)" class="hover-cursor" @click.native="nodeDetail(r)" :class="classObject(r)" >
              <Col span="13" class="item split-border" :style="{paddingLeft:r.depth * 15 +'px'}"> 
                <Icon v-if="r.show=='expand'" type="md-add" @click.stop="expand(r)" style="padding:3px 3px" />
                <Icon v-else-if="r.show=='collapse'" type="ios-remove" @click.stop="collapse(r)" style="padding:3px 3px"/>
                <!-- 这里的padding-left是为了让没有展开/收缩符号的文字跟有符号的文字左对齐 -->
                <span :style="{paddingLeft:calcTextMarginLeft(r)+'px'}">
                  <Icon v-show="r.icon=='hand'" type="ios-hand" /> 
                  <Icon v-show="r.icon=='bug'" type="ios-bug" />
                   <Icon v-show="r.icon=='info'" type="md-information-circle" />
                  {{getMethod(r.method)}}
                </span> 
              </Col>
              <Col span="4"  class="item"> {{r.params}}</Col>
              <Col span="2" class="item" style="padding-left:25px">{{showDuration(r.duration)}}</Col>
              <Col span="3" class="item"> {{r.service_type}}</Col>
              <Col span="2" class="item">{{r.app_name}}</Col>
            </Row>
          </div>
        </div>
    </Modal>

    <Drawer :title="selNode.method" :closable="false" v-model="isNodeSel" :styles="{'z-index':2000}" width=40>
        <Form :label-width="100" label-position="left">
          <FormItem label="方法名">
              {{selNode.method}}
          </FormItem>
           <FormItem label="Class">
              {{selNode.class}}
          </FormItem>
          <FormItem label="耗时(ms)">
              {{selNode.duration}}
          </FormItem>
          <FormItem label="应用名">
              {{selNode.app_name}}
          </FormItem>
          <FormItem label="服务器">
              {{selNode.agent_id}}
          </FormItem>
           <FormItem label="服务类型">
              {{selNode.service_type}}
          </FormItem>
          <FormItem label="参数">
              {{selNode.params}}
          </FormItem>
            
          <Divider orientation="center">Debug Part</Divider>

          <FormItem label="method id">
              {{selNode.method_id}}
          </FormItem>
           <FormItem label="sequence">
              {{selNode.seq}}
          </FormItem>
          <FormItem label="span depth">
              {{selNode.span_depth}}
          </FormItem>
           <FormItem label="node type">
              {{selNode.type}}
          </FormItem>
      </Form>
    </Drawer>
  </div>   
</template>

<script>
import request from '@/utils/request' 
import trace from './charts/trace'
import echarts from 'echarts'
import blueLineChart from './charts/blueLineChart'
import greenLineChart from './charts/greenLineChart'
import redLineChart from './charts/redLineChart'
export default {
  name: 'tracing',
  components: {trace,blueLineChart,greenLineChart,redLineChart},
  data () {
    return {
      tracesData: {},
      traceData: {},
       resultLimit: 50,
       minElapsed: null,
       maxElapsed: null,
      selectedTraces: [],
      searchError: false,
      searchTraceID : '',
      remoteAddr: '',
      apis : [],
      currentApi: '',
      
      chartGroup: 'tracing',

      traceLabels: [
            {
                title: '发生时间',
                key: 'showTime',
                width: 200
            },
            {
                title: 'API',
                key: 'api'
            },
             {
                title: '耗时(ms)',
                key: 'elapsed',
                width: 100
            },
            {
                title: '服务器ID',
                key: 'agentId',
                width: 150
            },
            {
                title: 'Remote Addr',
                key: 'remote_addr',
                width: 150
            },
            {
                title: '链路ID',
                key: 'traceId'
            },
            {
                title: '查看详情',
                slot: 'operation',
                width: 100,
                align: 'center'
            }
        ],

      selSucCount: 0,
      selErrCount: 0,

      traceVisible: false,

      selectedTrace: {},

       split1: 0.3,
       selNode: {},
       isNodeSel: false,
       collapseList: [],
       
       // 0: 同时显示成功、错误的点
       // 1: 只显示成功的点
       // 2: 只显示错误的点
       tableFilter: 0,


       timeline: [],
        jvmCpuList: [],
        systemCpuList: [],
        jvmHeapList: [],
        heapMaxList: [],
        fullgcDurationList: []
    }
  },
  watch: {
    "$store.state.apm.appName"() {
            this.apiList()
    }
  },
  computed: {
    
  },
  methods: {
    showTraceTable() {
      if (this.tableFilter==0) {
          return this.selectedTraces
      } else if (this.tableFilter == 1) {
        var traces = []
        for (var i=0;i<this.selectedTraces.length;i++) {
          if (this.selectedTraces[i].errCode==0) {
            traces.push(this.selectedTraces[i])
          }
        }
        return traces
      } 

      var traces = []
      for (var i=0;i<this.selectedTraces.length;i++) {
        if (this.selectedTraces[i].errCode==1) {
          traces.push(this.selectedTraces[i])
        }
      }
      return traces
    },
    showDuration(d) {
      if (d == -1) {
        return ''
      }

      return d
    },
    calcTextMarginLeft(r) {
      if (r.show != 'expand' && r.show != 'collapse') {
        return 21
      }
      return 0
    },
    expand(r) {
      r.show = 'collapse' 
      for (var i=0;i<this.collapseList.length;i++) {
        if (this.collapseList[i] == r.id) {
          this.collapseList.splice(i,1)
        }
      }
    },
    collapse(r) {
       r.show='expand'
       this.collapseList.push(r.id)
    },
    queryTraces() {
      this.$Loading.start();
      request({
          url: '/web/queryTraces',
          method: 'GET',
          params: {
            app_name: this.$store.state.apm.appName,
            api : this.currentApi,
            min_elapsed: this.minElapsed,
            max_elapsed: this.maxElapsed,
            limit: this.resultLimit,
            search_error: this.searchError,
            search_trace_id: this.searchTraceID.trim(),
            remote_addr: this.remoteAddr,
            start: JSON.parse(this.$store.state.apm.selDate)[0],
            end: JSON.parse(this.$store.state.apm.selDate)[1],
          }
      }).then(res => {
        this.tracesData = res.data.data
        this.selectedTraces = []
        console.log("query traces:",this.tracesData)
        this.$Loading.finish();
        if (!this.tracesData.is_suc) {
          this.$Message.info({
            content: '没有查询到数据',
            duration: 3 
          })
        }
      }).catch(error => {
        this.$Loading.error();
      })
    },
    isShow(r) {
      // 对于collapseList中的每个值，判断当前行是否在它的子树中，若在，则不显示，跳出循环
      // 若当前name是以collapseList的值为前缀，说明在子树中
      for (var i=0;i<this.collapseList.length;i++) {
        if (r.id == this.collapseList[i]) {
          continue
        }
        var j = r.id.indexOf(this.collapseList[i]);
        if(j == 0){
          return false
        }
      }
      return true
    },
    nodeDetail(r) {
      this.selNode = r
      this.isNodeSel = true
    },
    classObject: function (r) {
      var o = {}
      o[r.id] = true
      if (r.is_error) {
        o['error'] = true
      } 
      if (r.agent_id== this.selectedTrace.agentId && r.span_depth == 0) {
        o['current-span'] = true
      }
      return o
    },
    getMethod(s) {
      var i = s.split('(')
      return i[0]
    },
    rowClassName (row, index) {
                if (row.errCode == 1) {
                    return 'error-trace';
                } else  {
                    return 'success-trace';
                }
            },
    selectTraces(t) {
      this.selectedTraces = t
      var ec = 0,sc = 0
      for (var i=0;i<t.length;i++) {
        if (t[i].errCode==1) {
          ec+=1
        } else {
          sc+=1
        }
      }

      this.selErrCount = ec
      this.selSucCount = sc
    },
    showTrace(t) {
      this.selectedTrace = t
      this.$Loading.start();
      // 查询trace详情
      request({
          url: '/web/trace',
          method: 'GET',
          params: {
            trace_id : t.traceId
          }
      }).then(res => {
        this.traceData = res.data.data
        for (var i=0;i<this.traceData.length;i++) {
          if (this.traceData[i].span_depth != -1) {
            this.traceData[i].show = 'collapse'
          } 
        }
        this.traceVisible = true
        console.log(this.selectedTrace)
        console.log(this.traceData)
        
        this.$Loading.finish();
      }).catch(error => {
        this.$Loading.error();
      })

      this.dashboard(this.selectedTrace.agentId,this.selectedTrace.startTime)
    },
    apiList() {
       request({
          url: '/web/appApis',
          method: 'GET',
          params: {
              app_name: this.$store.state.apm.appName,
          }
      }).then(res => {
        this.apis = res.data.data
      })
    },
    // 加载JVM详细图表
    dashboard(agentID,startTime) {
      this.$Loading.start();
      // 加载当前APP的dashbord数据
      request({
          url: '/web/runtimeDashByUnixTime',
          method: 'GET',
          params: {
            app_name: this.$store.state.apm.appName,
            start: startTime/1000 - 30,
            end: startTime/1000 + 30,
            agent_id: agentID
          }
      }).then(res => {   
        console.log(res.data.data)
        this.timeline =  res.data.data.timeline
        this.jvmCpuList = res.data.data.jvm_cpu_list
        this.systemCpuList = res.data.data.sys_cpu_list        
        this.jvmHeapList = res.data.data.jvm_heap_list
        this.heapMaxList = res.data.data.heap_max_list
        this.fullgcDurationList = res.data.data.fullgc_duration_list

          this.$Loading.finish();
      }).catch(error => {
        this.$Loading.error();
      })
    },

  },
  mounted() {
    this.apiList()
    echarts.connect(this.chartGroup);
  }
}
</script>

</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";
.meta-word {
  margin-right:4px;
}

.trace-pane {
  .title {
    font-size: 15px;
    background: #f3f3f3;
    padding-left:5px;
    padding-top:6px;
    padding-bottom:6px;
  }
  .body {
    .item {
      padding-left:7px;
      padding-top: 8px;
      padding-bottom: 8px;
      overflow: hidden;
      text-overflow:ellipsis;
      white-space: nowrap;
      font-size:12px;
      // line-height: 20px;
      // height: 30px;
    }
    .error {
      color: #d62727;
    }
    .hover-cursor:hover {
      background-color: #ebf7ff !important
    }
    .current-span {
      background: #dff0d8;
      color: #1469eb; 
    }
  }
}
</style>

<style lang="less">
    .ivu-table .error-trace td{
        background-color: rgba(223, 83, 83, .5);
        color: #333;
    }
    .ivu-table td:hover {
      cursor: pointer;
    }

    .custom-trigger {
      cursor:col-resize;
    }

    .api-filter {
      .ivu-select-dropdown-list {
        max-width: 395px;
        // overflow: hidden;
        // text-overflow:ellipsis;
        // white-space: nowrap;
      }
    }
  .ivu-drawer {
    .ivu-drawer-header {
      display: none;
    }

    .ivu-drawer-body {
      padding-left: 30px;
      padding-top:30px;
    }
  }

  .ivu-modal-header {
    padding-top: 0px !important;
  }
  .ivu-modal-body {
    padding: 0;
    margin-top: 16px;
  }
</style>


