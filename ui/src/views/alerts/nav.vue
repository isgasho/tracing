<template>
  <div class="alerts-nav">
      <div class="nav">
          <span v-for="i in items" :key="i" :class="{'selected': selItem==i}" class="hover-cursor item" @click="selectItem(i)">{{names[i]}}</span>
      </div>
            <router-view></router-view>
  </div>
</template>

<script>
import request from '@/utils/request' 
export default {
  name: 'alertsNav',
  data () {
    return {
      items: [],
      level: {},

      path : '',
      selItem : '',


      appNames: []
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
      this.$router.push('/ui/alerts/' + i)
    },
    initItem() {
        this.appNames = [this.$store.state.apm.appName]
        this.path = window.location.pathname
        this.items = ['appList','policy','alertsNotify']
        this.level = {'appList':2, alertsNotify:2,policy:2, group:2}
        this.names = {appList: '应用告警','alertsNotify': "告警消息查询",
            policy: '策略模版',group:'用户组管理'}
        this.selItem = this.path.split('/')[4]
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
.nav {
    // border-bottom: 1px solid black;
    background-color: #595959;
    color:white;
    padding-top:8px;
    padding-bottom:6px;
    font-size:13px;
    .item {
        margin-left: 15px;
        padding-left: 10px;
        padding-right: 10px;
    }
}

.item.selected {
    font-weight: 700;
    background-color: #6b6b6b;
    padding-top:11px;
    padding-bottom:10px
}

    

</style>
