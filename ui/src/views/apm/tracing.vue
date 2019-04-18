<template>
  <div>
    <Row class="padding-20">
       <Col span="8" style="background:rgb(245, 245, 245)"  class="padding-20 padding-bottom-40">
          <div  class="font-size-18">链路过滤</div>
          <!-- <div class="margin-top-10">
            <div class="font-size-16">请求API</div>
            <Select v-model="currentApi" style="width:400px" size="large" class="api-filter"  placeholder="默认选择全部API" filterable clearable>
              <Option v-for="api in apis" :value="api" :key="api">
                {{ api }}
              </Option>
            </Select>
          </div> -->
          
          <div class="margin-top-10">
            <div class="font-size-16">最低响应时间(ms)</div>
            <Input  v-model="minElapsed" placeholder="e.g.  100，留空代表不限制" style="width: 100%" size="large" />
          </div>
          <div class="margin-top-40">
            <div class="font-size-16">最大响应时间(ms)</div>
            <Input   v-model="maxElapsed" placeholder="e.g. 3000，留空代表并不限制" style="width: 100%" size="large" />
          </div>
          <div class="margin-top-40">
            <div class="font-size-16">限定搜索数目</div>
            <InputNumber :max="100" :min="1"  :step="5" v-model="resultLimit" style="width: 100%" size="large"></InputNumber>
          </div>

          <div class="margin-top-20">
            <Button type="primary"  class="primary2-button" @click="queryTraces">查询链路</Button>
          </div>
       </Col>
       <Col span="14" offset="1">
        <trace :graphData="tracesData" v-if="tracesData.is_suc==true" style='height: 49vh' @selTraces="selectTraces"></trace>
       </Col>
    </Row>

     <div v-show="selectedTraces.length > 0" class="margin-top-10" style="padding-left:20px;padding-right:20px;">
          <span><Tag style="background: rgb(18, 147, 154,.5)" size="large">{{selSucCount}} 成功</Tag> <Tag style="background:rgba(223, 83, 83, .5)" size="large">{{selErrCount}} 错误</Tag></span>
          <!-- <span style="float:right" class="no-border">
            <Select style="width:200px"  value="1">
                <Option value="1">时间最近</Option>
                <Option value="2">时间最远</Option>
                <Option value="3">链路最短</Option>
                <Option value="3">链路最长</Option>
            </Select>
          </span> -->

        <div style="padding-left:10px;padding-right:10px" class="margin-top-20 margin-bottom-40">
          <Table :row-class-name="rowClassName"  :columns="traceLabels" :data="selectedTraces" class="margin-top-15" @on-row-click="showTrace"></Table>
        </div>
      </div>

    <!-- 链路展示Modal -->
    <Modal v-model="traceVisible" :footer-hide="true" :z-index="500" fullscreen>
        <div slot="header" style="padding-top:5px;padding-bottom:0px;border-bottom:none">
            <div class="font-size-16" style="font-weight:bold">{{$store.state.apm.appName}} : {{selectedTrace.api}}</div>
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
        </div>

        <div class="padding-top-5 trace-pane" style="padding-bottom:50px">
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
export default {
  name: 'tracing',
  components: {trace},
  data () {
    return {
      tracesData: {},
      traceData: {},
       resultLimit: 50,
       minElapsed: null,
       maxElapsed: null,
      selectedTraces: [],

      apis : [],
      currentApi: '',

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
                width: 250
            },
            {
                title: '链路ID',
                key: 'traceId'
            },

        ],

      selSucCount: 0,
      selErrCount: 0,

      traceVisible: false,

      selectedTrace: {},

       split1: 0.3,
       selNode: {},
       isNodeSel: false,
       collapseList: []
    }
  },
  watch: {
  },
  computed: {
    
  },
  methods: {
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
          url: '/apm/web/queryTraces',
          method: 'GET',
          params: {
            app_name: this.$store.state.apm.appName,
            api : this.currentApi,
            min_elapsed: this.minElapsed,
            max_eapsed: this.maxElapsed,
            limit: this.resultLimit,

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
          url: '/apm/web/trace',
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
    }
  },
  mounted() {
    request({
        url: '/apm/web/appApis',
        method: 'GET',
        params: {
           app_name: this.$store.state.apm.appName,
        }
    }).then(res => {
      this.apis = res.data.data
    })
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
</style>


