<template>
  <!-- 此页面是模版文件，用于快速搭建其他页面 -->
  <div class="app-container">
      <div>{{hello}} {{world}}</div>
      <div>
        {{world}}
        <div class="margin-left-20">{{getAddr}}</div>
        你好中国人
      </div>

      <Button>Default</Button>
    <Button type="primary">Primary</Button>
    <Button type="dashed">Dashed</Button>
    <Button type="text">Text</Button>
  </div>
</template>

<script>
import request from '@/utils/request' 
export default {
  name: 'tmpl',
  data () {
    return {
      hello: '',
      world: ''
    }
  },
  watch: {
    $route() {
      // 当路由变化时，调用此函数
    },
    world(val) {
      // world发生变化时
    }
  },
  computed: {
    getAddr() {
      return process.env.API_ADDR
    }
  },
  methods: {
  },
  mounted() {
    this.$store.dispatch('setService', 'hello')

    this.hello = this.$store.state.misc.service
    this.world = this.$t('common.test')

    setTimeout(() => {
        this.world = "new world"
    }, 3000)


    request({
        url: '/login',
        method: 'POST',
        params: {
          user: '111'
        }
    }).then(res => {
    }).catch(err => {
 
            })
  }
}
</script>

<style lang="less">
</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";


</style>
