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
            <Input  v-model="minElapsed" placeholder="e.g.  100，留空代表不限制" style="width: 400px" size="large" />
          </div>
           <div class="margin-top-40">
            <div class="font-size-16">最大响应时间(ms)</div>
            <Input   v-model="maxElapsed" placeholder="e.g. 3000，留空代表并不限制" style="width: 400px" size="large" />
          </div>
          <div class="margin-top-40">
            <div class="font-size-16">限定搜索数目</div>
            <InputNumber :max="100" :min="1"  :step="5" v-model="resultLimit" style="width: 400px" size="large"></InputNumber>
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
          <span style="float:right" class="no-border">
            <Select style="width:200px"  value="1">
                <Option value="1">时间最近</Option>
                <Option value="2">时间最远</Option>
                <Option value="3">链路最短</Option>
                <Option value="3">链路最长</Option>
            </Select>
          </span>

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

        <div class="padding-top-5 trace-pane">
          <Row>
            <Col span="10" class="title split-border">方法(点击具体方法名可查看详情)</Col>
            <Col span="5"  class="title split-border">参数</Col>
            <Col span="3" class="title split-border">耗时(ms)</Col>
            <Col span="3" class="title split-border">API</Col>
            <Col span="3" class="title">所属应用</Col>
          </Row>
          <div  class="body">
            <Row v-for="r in traceData.value" v-show="isShow(r)" class="hover-cursor" @click.native="spanDetail(r)" :class="classObject(r)" >
              <Col span="10" class="item split-border" :style="{paddingLeft:r.paddingLeft+'px'}"> {{getMethod(r.method)}}</Col>
              <Col span="5"  class="item split-border"> {{r.params}}</Col>
              <Col span="3" class="item split-border">{{r.self}}</Col>
              <Col span="3" class="item split-border"> {{r.api}}</Col>
              <Col span="3" class="item">{{r.agentName}}</Col>
            </Row>
          </div>
        </div>
    </Modal>

    <Drawer :title="selItem.method" :closable="false" v-model="isItemSel" :styles="{'z-index':2000}" width=40>
        <Form :label-width="80">
          <FormItem label="发生时间">
              {{selItem.startTimeStr}}
          </FormItem>
          <FormItem label="耗时(ms)">
              {{selItem.self}}
          </FormItem>
          <FormItem label="应用名">
              {{selItem.agentName}}
          </FormItem>
          <FormItem label="服务器">
              {{selItem.agentId}}
          </FormItem>
           <FormItem label="Class">
              {{selItem.class}}
          </FormItem>
           <FormItem label="Api">
              {{selItem.api}}
          </FormItem>
          <FormItem label="参数">
              {{selItem.params}}
          </FormItem>
          <FormItem label="层级(Debug)">
              {{selItem.name}}
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
       selItem: {},
       isItemSel: false,
       collapseList: []
    }
  },
  watch: {
  },
  computed: {
    
  },
  methods: {
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
        console.log(this.tracesData)
        this.$Loading.finish();
      }).catch(error => {
        this.$Loading.error();
      })
    },
    isShow(r) {
      // 对于collapseList中的每个值，判断当前行是否在它的子树中，若在，则不显示，跳出循环
      // 若当前name是以collapseList的值为前缀，说明在子树中
      for (var i=0;i<this.collapseList.length;i++) {
        var j = r.name.indexOf(this.collapseList[i]);
        if(j == 0){
          return false
        }
      }
      return true
    },
    spanDetail(r) {
      this.selItem = r
      this.isItemSel = true
    },
    classObject: function (r) {
      var o = {}
      o[r.name] = true
      if (r.classStyle == 'error') {
        o['error'] = true
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
            traceId : t.traceId
          }
      }).then(res => {
        this.traceData = JSON.parse(res.data.data)
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
    padding-bottom:6px
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
</style>


