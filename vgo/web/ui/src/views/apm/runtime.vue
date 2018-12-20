<template>
  <div>
            <Tabs value="name1" :animated=false>
                <TabPane label="Open incidents" name="name1">标签一的内容</TabPane>
                <TabPane label="All incidents" name="name2">标签二的内容</TabPane>
            </Tabs>
      </div>   
</template>

<script>
import request from '@/utils/request' 
export default {
  name: 'runtime',
  data () {
    return {
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


<style lang="less" scoped> 
@import "../../theme/gvar.less";
</style>
