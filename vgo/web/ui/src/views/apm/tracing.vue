<template>
  <div>
    <Row class="padding-20">
       <Col span="8" style="background:rgb(245, 245, 245)"  class="padding-20 padding-bottom-40">
          <div  class="font-size-18">链路过滤</div>
          <div class="margin-top-10">
            <div class="font-size-16">请求API</div>
            <Select style="width:300px" size="large">
                <Option value="beijing">New York</Option>
                <Option value="shanghai">London</Option>
                <Option value="shenzhen">Sydney</Option>
            </Select>
          </div>
          
          <div class="margin-top-40">
            <div class="font-size-16">最低响应时间</div>
            <Input  placeholder="e.g. 1.2s 100ms 500us" style="width: 300px" size="large" />
          </div>
           <div class="margin-top-40">
            <div class="font-size-16">最大响应时间</div>
            <Input  placeholder="e.g. 3s" style="width: 300px" size="large" />
          </div>
          <div class="margin-top-40">
            <div class="font-size-16">限定搜索数目</div>
            <InputNumber :max="100" :min="1"  :step="5" v-model="resultLimit" style="width: 300px" size="large"></InputNumber>
          </div>
       </Col>
       <Col span="14" offset="1">
        <trace :graphData="JSON.parse(tracesData)" style='height: 49vh' @selTraces="selectTraces"></trace>
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

        <div style="padding-left:10px;padding-right:10px" class="margin-top-20">
          <Table :row-class-name="rowClassName"  :columns="traceLabels" :data="selectedTraces" class="margin-top-15" @on-row-click="showTrace"></Table>
        </div>
        </div>
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
      tracesData: null,
      traceData: null,
       resultLimit: 20,
      selectedTraces: [],

      traceLabels: [
            {
                title: '发生时间',
                key: 'showTime',
                width: 200
            },
            {
                title: 'API',
                key: 'url'
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
      selErrCount: 0
    }
  },
  watch: {
  },
  computed: {

  },
  methods: {
    rowClassName (row, index) {
                if (row.errCode == 1) {
                    return 'error-trace';
                } else  {
                    return 'success-trace';
                }
            },
    selectTraces(t) {
      this.selectedTraces = t
      for (var i=0;i<t.length;i++) {
        if (t[i].errCode==1) {
          this.selErrCount++
        } else {
          this.selSucCount++
        }
      }
    },
    showTrace(t) {
      // 查询trace详情
      request({
        url: '/apm/query/trace',
        method: 'GET',
        params: {
          traceId : t.traceId
        }
    }).then(res => {
      this.traceData = res.data.data
      console.log(JSON.parse(this.traceData))
    })
    }
  },
  mounted() {
    request({
        url: '/apm/query/traces',
        method: 'GET',
        params: {
        }
    }).then(res => {
      this.tracesData = res.data.data
    })
  }
}
</script>

</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";
</style>

<style>
    .ivu-table .error-trace td{
        background-color: rgba(223, 83, 83, .5);
        color: #333;
    }
    .ivu-table td:hover {
      cursor: pointer;
    }
</style>