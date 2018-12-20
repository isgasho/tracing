<template>
  <div class="apm-nav">
      <Row style="line-height:22px">
          <Col span="3" style="border-right: 1px solid #d2d2d2;border-bottom:3px solid #e0ebd1" class="padding-left-20 left-nav margin-top-5">
             <div class="color-primary font-size-16 padding-top-15 hover-cursor">应用设定</div>
             <div class="padding-bottom-5">{{$store.state.apm.appName}}</div>
          </Col>
          <Col span="21" class="padding-left-20" style="border-bottom:3px solid #e0ebd1;">
            <div class="color-primary font-size-16 padding-top-15 hover-cursor">时间设定</div>
            <div class="padding-bottom-5">2018-12-05 00:00 - 2018-12-06 00:00</div>
             <!-- <DatePicker type="datetimerange" format="yyyy-MM-dd HH:mm" placeholder="Select date and time(Excluding seconds)" style="width: 280px"></DatePicker> -->
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
      selItem : ''
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
    selectItem(i) {
      this.$router.push('/apm/' + i)
    },
    initItem() {
        this.path = window.location.pathname
        this.items = ['monitoring','dashboard','tracing','serviceMap','runtime','profiling','thread','memory','stats','database','interface','exception']
        this.level = {monitoring: 1,'dashboard':2, tracing:2,serviceMap:2, runtime:2, profiling:1,thread:2,memory:2,stats:1,database:2,interface:2,exception:2}
        this.names = {monitoring: '监控','dashboard': "应用总览",
            tracing: '链路跟踪',serviceMap:'应用拓扑',  
            runtime: '应用运行时', profiling: '在线诊断', thread:'线程剖析', memory: '内存剖析',
            stats: '数据统计', 
            database:'数据库', interface:'访问接口', exception:'错误异常'}
        this.selItem = this.path.split('/')[2]
    }
  },
  mounted() {
      this.initItem()
  }
}
</script>

<style lang="less">

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
