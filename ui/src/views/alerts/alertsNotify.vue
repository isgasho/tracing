<template>
  <div class="app-container">
      <Row>
          <Col span="22" offset="1">
            <div class="header font-size-18">
                <Select v-model="selApp" style="width:220px" placeholder="选择你的应用" filterable>
                    <Option v-for="item in appNames" :value="item" :key="item">{{ item }}</Option>
                </Select>
            </div>
            <Table stripe :columns="appLabels" :data="appList" class="margin-top-15" @on-row-click="editApp">
                <!-- <template slot-scope="{ row }" slot="owner">
                  {{ row.owner_name + '/'+row.owner_id}}
                </template> -->
            </Table>
          </Col>
      </Row>
  </div>
</template>

<script>
import request from '@/utils/request' 
export default {
  name: 'alertsNotify',
  data () {
    return {
      appNames : [],
      selApp: '',
      appLabels: [
            {
                title: '应用名',
                key: 'app_name'
            },
            {
                title: '服务器',
                key: 'agent_id',
            },
            {
                title: '报警项',
                key: 'alert',
            },
            {
                title: 'API',
                key: 'api',
            },
            {
              title: '告警值',
                key: 'alert_value',  
            },
            {
              title: '告警通道',
                key: 'channel',  
            },
            {
              title: '告警用户',
                key: 'users',  
            },
            {
              title: '告警时间',
                key: 'alert_time',  
            }
        ],
        appList: [],
    }
  },
  watch: {
  },
  computed: {
  },
  methods: {
  },
  mounted() {
    request({
            url: '/apm/web/appNames',
            method: 'GET',
            params: {
            }
        }).then(res => {   
            this.appNames = res.data.data 
        })
  }
}
</script>

<style lang="less">
</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";

</style>
