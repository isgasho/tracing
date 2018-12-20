<template>
  <div class="app-containe padding-top-10" style="overflow-y:scroll;overflow-x:none;max-height:100vh;background-color:white">
     <Row :gutter="20">
         <Col span="12">
            <respTime width="100%" height="400px" id="apm-resp"></respTime>
         </Col>
         <Col span="11">
            <apdex width="100%" height="410px" id="apm-apdex"></apdex>
         </Col>
     </Row>
     <Row :gutter="20">
        <Col span="10">
              <error width="100%" height="300px" id="apm-error"></error>
             <rpm width="100%" height="300px" id="apm-rpm"></rpm>
         </Col>
          <Col span="5" offset="1" style="padding:8px 10px;padding-left:20px">
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
          <Col span="7" style="padding:8px 10px;padding-left:20px">
                <div class="font-size-18">慢事务列表</div>
                <Table :columns="trLabels" class="margin-top-10" @on-row-click="gotoApp"></Table>
          </Col>
     </Row>
  </div>
</template>

<script>
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
                title: '接口名',
                key: 'name'
            },
            {
                title: '响应时间(ms)',
                key: 'resp_time'
            }
        ],
    }
  },
  watch: {
  },
  computed: {
  },
  methods: {
  },
  mounted() {
  }
}
</script>

<style lang="less">

</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";


</style>
