<template>
    <Tooltip placement="bottom" max-width="800" :delay="300" @on-popper-show="beforeShow">
        {{policyName}}
        <div slot="content" style="padding: 15px 15px">
           <Table stripe :columns="policyLabels" :data="policy.alerts" class="margin-top-15">
                <template slot-scope="{ row }" slot="type">
                    <span v-if="row.type=='apm'">应用监控</span>
                    <span v-else>系统监控</span>
                </template>
                <template slot-scope="{ row }" slot="duration">
                    <span>{{row.duration}}分钟</span>
                </template>
                <template slot-scope="{ row }" slot="compare">
                    <span v-if="row.compare==1">大于</span>
                    <span v-else-if="row.compare==2">等于</span>
                    <span v-else>小于</span>
                </template>
            </Table>
        </div>
    </Tooltip>
</template>

<script>
import request from '@/utils/request' 
export default {
  props: {
    policyID: {
      default: ''
    },
    policyName: {
        default: ''
    }
  },
  data() {
    return {
        policy: {},
        policyLabels: [
            {
                title: '名称',
                key: 'label',
            },
            {
                title: '类型',
                slot: 'type',
                width: 100
            },
            {
                title: '持续时间',
                slot: 'duration',
                width: 100
            },
            {
              title: '比较方式',
                slot: 'compare',
                width: 70 
            },
            {
              title: '触发值',
                key: 'value',
                width: 100
            }
        ]
    }
  },
  methods: {
      beforeShow() {
          request({
                    url: '/apm/web/queryPolicy',
                    method: 'GET',
                    params: {
                        id: this.policyID
                    }
                }).then(res => {   
                    this.policy = res.data.data
                })    
      }
  },
  mounted() {
     
  },
  beforeDestroy() {
  }
}
</script>
