<template>
  <span>
        <Tooltip placement="bottom" max-width="400">
        {{alert.label}}
        <div slot="content" style="padding: 15px 15px">
        <div class="alert-setting">
            <span class="label">比较方式</span>
            <span style="margin-left:10px">
            <span v-if="alert.compare==1"> > </span>
            <span v-else-if="alert.compare==1"> = </span>
            <span v-else> &lt; </span>
            </span>
        
        </div>

        <div class="alert-setting">
            <span class="label">设定阈值</span>
            <span style="margin-left:10px">{{alert.value}}{{alert.unit}}</span>
        </div>

        <div class="alert-setting">
            <span class="label">持续时间</span>
            <span style="margin-left:10px">{{alert.duration}} 分钟</span>
        </div>
        </div>
    </Tooltip>
    </span>
</template>

<script>
export default {
  props: {
    alert: {
      default: {}
    }
  },
  data() {
    return {
  
    }
  },
  mounted() {
  },
  beforeDestroy() {
  },
  methods: {
  }
}
</script>
