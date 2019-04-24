<template>
  <div class="mf-doc-page">
      <!-- <div class="font-size-30 font-weight-500 page-header margin-left-20">{{$t('pageName.install')}}</div> -->
     
  </div>
</template>

<script>
export default {
  name: 'install',
  data () {
    return {
    }
  },

  computed: {
 
  },
  methods: {
  },
  mounted() {
      this.$store.state.misc.content = this.$t('pages.install')
  }
}
</script>

<style lang="less">


</style>

<style lang="less" scoped> 


</style>
