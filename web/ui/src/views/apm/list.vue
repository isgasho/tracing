<template>
  <div>
    <Row class="subnav" style="background-color:#595959;border-bottom: 1px solid #474747;color:#eaeaea;padding: 8px 10px;padding-bottom:6px;vertical-align:middle;font-size:13px;margin-top:6px">
          <span class="item">总应用数：<span class="count bg-second">300</span></span>
          <span class="item">不健康应用数： <span class="count bg-second">0</span></span>
          <span class="item">最近3小时告警数： <span class="count bg-second">0</span></span>
      </Row>
      <div class="app-container">
       <Row style="padding:0 10px;" class="split-border-bottom no-border">
          <Col span="17" class="split-border-right">
            <span class="padding-bottom-5 font-size-18">应用列表</span>
            <Tag style="margin-top: -3px">最近5分钟</Tag>
             <Select v-model="selApps" filterable multiple style="width:300px;border:none;float:right;margin-right:20px" placeholder="过滤应用">
                <Option v-for="item in appNames" :value="item.value" :key="item.value">{{ item.label }}</Option>
            </Select>
          </Col>
           <Col span="6">
            <span class="padding-bottom-5 font-size-18 margin-left-10" >应用动态</span>
          </Col>
      </Row>
      <Row style="padding:0 10px">
          <Col span="17" class="split-border-right no-border" style="padding:8px 10px;">
             <Table stripe :columns="appLabels" :data="appList" class="margin-top-15" @on-row-click="gotoApp"></Table>
             <Page :current="1" :total="totalApps" size="small" class="margin-top-15" simple />
          </Col>
           <Col span="6"  style="padding:8px 10px;padding-left:20px">
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
      </Row>

      </div>   
  </div>
</template>

<script>
import request from '@/utils/request' 
export default {
  name: 'appList',
  data () {
    return {
        appNames: [
        ],
        selApps: [],

         appLabels: [
            {
                title: '应用名',
                key: 'name'
            },
            {
                title: '请求总数',
                key: 'count',
            },
            {
                title: 'Apdex',
                key: 'apdex'
            },
            {
                title: '响应时间(ms)',
                key: 'average_elapsed'
            },
            {
                title: '错误率',
                key: 'error_percent'
            },
        ],
        appList: [
        ],

        totalApps: 3
    }
  },
  watch: {
  },
  computed: {

  },
  methods: {
      gotoApp(app) {
          this.$store.dispatch('setAPPID', app.id)
          this.$store.dispatch('setAPPName', app.name)
          this.$router.push("/apm/ui/dashboard")
      }
  },
  mounted() {
      // 加载APPS
       request({
        url: '/apm/web/appListWithSetting',
        method: 'GET',
        params: {
          
        }
    }).then(res => {
      this.appList = res.data.data
      console.log(res.data)
    })
  }
}
</script>

<style lang="less">
</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";
.subnav {
    .item {
        margin-left: 20px;
        .count {
            display: inline-block;
            padding: 1px 7px;
            border-radius: 4px;
            text-shadow: 0 1px 2px rgba(0,0,0,0.2);
        }
    }
}
</style>
