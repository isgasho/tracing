<template>
  <div class="apm-nav">
      <Row style="line-height:22px">
          <Col span="3" style="border-right: 1px solid #d2d2d2;border-bottom:3px solid #e0ebd1;height:67px" class="padding-left-20 left-nav margin-top-5">
             <div  style="margin-top:20px" class="no-border">
              <span>
                <Select v-model="$store.state.apm.appName" style="width:180px" @on-change="selAppName" filterable>
                  <Option v-for="item in appNames" :value="item" :key="item">{{ item }}</Option>
                </Select>
              </span></div>
          </Col>
          <Col span="21" class="padding-left-20" style="border-bottom:3px solid #e0ebd1;height:67px">
            <div class="color-primary font-size-16 padding-top-15 hover-cursor">
              <DatePicker :split-panels=true  confirm type="datetimerange" :options="options1" :value="getDate()"  placeholder="启止时间设定" style="width: 320px;margin-left: 10px;margin-top:5px" @on-change="changeDate"  @on-ok="confirmDate" @on-clear="clearDate" :clearable=false :editable=false></DatePicker>
            </div>
          </Col>
      </Row>

       <Row>
          <Col span="3"  class="nav-left">
            <div class="left-items">
              <div v-for="i in items" :key="i" class="item item-1"  v-if="level[i]==1" >
                  {{names[i]}}
              </div>
              <div class="item hover-cursor item-2" :class="{'selected': selItem==i}" @click="selectItem(i)" v-else>
                {{names[i]}}
              </div>
            </div>
          </Col>    
          <Col span="21" style="background-color:white">
            <router-view></router-view>
          </Col>
      </Row>
  </div>
</template>

<script>
import request from '@/utils/request' 
export default {
  name: 'apmNav',
  data () {
    return {
      items: [],
      level: {},

      path : '',
      selItem : '',

      selDate: [],

      appNames: [],

      options1: {
        disabledDate (date) {
                        return date && date.valueOf() > Date.now() + 86400000
        },
          shortcuts: [
             {
                  text: '30m',
                  value () {
                    var d = new Date()
                      return [new Date(d.getTime() - 1800 * 1000),d];
                  }
              },
             {
                  text: '1h',
                  value () {
                    var d = new Date()
                      return [new Date(d.getTime() - 3600 * 1000),d];
                  }
              },
              {
                  text: '3h',
                  value () {
                    var d = new Date()
                      return [new Date(d.getTime() - 3600 * 1000 *3) ,d];
                  }
              },
               {
                  text: '6h',
                  value () {
                    var d = new Date()
                      return [new Date(d.getTime() - 3600 * 1000 *3) ,d];
                  }
              },
              {
                  text: '1d',
                  value () {
                    var d = new Date()
                      return [new Date(d.getTime() - 3600 * 1000 * 24),d];
                  }
              },
              {
                  text: '3d',
                  value () {
                     var d = new Date()
                      return [new Date(d.getTime() - 3600 * 1000 * 24*3),d];
                  }
              },
              {
                  text: '7d',
                  value () {
                      var d = new Date()
                      return [new Date(d.getTime() - 3600 * 1000 * 24*7),d];
                  }
              }
          ]
      }
    }
  },
  watch: {
    $route() {
      this.initItem()
    }
  },
  computed: {
  },
  methods: {
    getDate() {
      var d = this.$store.state.apm.selDate
      if (d == '') {
        return [new Date((new Date()).getTime() - 3600 * 1000).toLocaleString('chinese',{hour12:false}).replace(/\//g,'-'),new Date().toLocaleString('chinese',{hour12:false}).replace(/\//g,'-')]
      } else {
        return   JSON.parse(this.$store.state.apm.selDate)
      }
    },
    clearDate() {
       this.$store.dispatch('setSelDate', '')
    },
    selAppName(appName) {
      this.$store.dispatch('setAPPName', appName)
    },
    changeDate(date) {
      this.selDate = date
    },
    confirmDate() {
      if (this.selDate != undefined) {
         this.$store.dispatch('setSelDate', JSON.stringify(this.selDate))
      }
    },
    selectItem(i) {
      this.$router.push('/apm/ui/' + i)
    },
    initItem() {
       this.appNames = [this.$store.state.apm.appName]
        this.path = window.location.pathname
        this.items = ['monitoring','dashboard','tracing','serviceMap','runtime','profiling','thread','memory','stats','database','interface','method','exception']
        this.level = {monitoring: 1,'dashboard':2, tracing:2,serviceMap:2, runtime:2, profiling:1,thread:2,memory:2,stats:1,database:2,interface:2,exception:2,method:2}
        this.names = {monitoring: '监控','dashboard': "应用总览",
            tracing: '链路跟踪',serviceMap:'应用拓扑',  
            runtime: '应用运行时', profiling: '在线诊断', thread:'线程剖析', memory: '内存剖析',
            stats: '数据统计', 
            database:'数据库', interface:'访问接口', exception:'错误异常',method:'关键事务'}
        this.selItem = this.path.split('/')[3]
        // 加载app名列表
         request({
            url: '/apm/web/appNamesWithSetting',
            method: 'GET',
            params: {
            }
        }).then(res => {   
            this.appNames = res.data.data 
        })
    }
  },
  mounted() {
      this.initItem()
  }
}

function defaultDate() {
  var d = new Date()
  var dates = [new Date(d.getTime() - 3600 * 1000),d]
  return dates
}
</script>

<style lang="less">
.ivu-date-picker {
  .ivu-picker-panel-sidebar {
    text-align: center;
    padding-top: 11px;
    .ivu-picker-panel-shortcut {
      margin-top: 5px;
    }
    .ivu-picker-panel-shortcut:hover {
      cursor: pointer
    }
  }
  .ivu-picker-panel-sidebar::before {
    content : '距离当前';
    font-size: 12px;
    color: #c5c8ce
  }
  .ivu-picker-panel-sidebar:hover {
    cursor: auto  
  }

  .ivu-picker-confirm {

  }
  input {
      border:none;
  }
}

</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";
.nav-left {
      position: relative;
      .left-items {
        margin-top:10px;
        // overflow-y: scroll;
      }
      .site {
        margin-left: 30px;
      }
       border-right: 1px solid #d2d2d2;
}

.item-2 {
    // border-radius: 4px;
    transition: background-color .3s ease-in-out;
    padding : 12px 20px;
        padding-left:30px;
    // padding-left: 45px !important
}

.item-1 {
    margin-left:20px;
    padding: 12px 0px;
    color: #888;
    font-size:15px
}

.item.selected {
    color: #555;
    font-weight: 700;
    
    border-top: 1px solid #d2d2d2;
    border-bottom: 1px solid #d2d2d2!important;
    border-left: 5px solid @primary-color;
}

    

</style>
