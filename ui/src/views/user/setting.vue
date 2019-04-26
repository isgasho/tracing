<template>
  <div class="app-container" style="margin-top:30px">
      <Row>
          <Col span="14" offset="5">
            <Form :model="formItem" :label-width="120">
                <FormItem label="应用设定">
                    <RadioGroup v-model="formItem.app_show">
                        <Radio :label=1>显示全部应用</Radio>
                        <Radio :label=2>自定义应用列表</Radio>
                    </RadioGroup>
                    <Tooltip placement="right" max-width="400">
                        <Icon type="ios-help" style="margin-bottom:5px;font-size:16px"  />
                    <div slot="content" style="padding: 15px 15px">
                        <div class="font-size-18 font-weight-500" style="line-height:25px">什么是应用设定</div>
                        <div>APM默认展示所有应用的监控状态，因此应用列表往往太长，不便于迅速找到自己名下的应用列表。用户可以在此处进行自行设置，忽略非设定的应用</div>
                        <div class="font-size-18 font-weight-500" style="line-height:25px">还是不明白？</div>
                        <div>请联系APM管理员</div>
                    </div>
                    </Tooltip>
                    <div class="margin-top-10">
                         <Select v-model="formItem.select" v-show="formItem.app_show==2" style="width:300px" multiple  filterable>
                            <Option v-for="item in appNames" :value="item" :key="item">{{ item }}</Option>
                        </Select>
                    </div>
                   
                </FormItem>
                
                <FormItem >
                    <Button type="primary" @click="setPerson">提交</Button>
                </FormItem>
            </Form>
          </Col>
      </Row>
  </div>
</template>

<script>
import request from '@/utils/request' 
export default {
  name: 'personSetting',
  data () {
    return {
        formItem: {
                    select: [],
                    app_show: 1
                },
        appNames: []
    }
  },
  watch: {
  },
  computed: {
  },
  methods: {
      setPerson() {
          request({
            url: '/web/setPerson',
            method: 'POST',
            params: {
               app_names : JSON.stringify(this.formItem.select),
               app_show: this.formItem.app_show
            }
        }).then(res => {   
            this.appNames = res.data.data 
            this.$Message.success({
                content: '设置成功',
                duration: 3 
            })
        })
      }
  },
  mounted() {
      request({
            url: '/web/appNames',
            method: 'GET',
            params: {
            }
        }).then(res => {   
            this.appNames = res.data.data 
        })
        request({
            url: '/web/getAppSetting',
            method: 'GET',
            params: {
            }
        }).then(res => {   
            console.log(res.data.data)
            this.formItem.select = JSON.parse(res.data.data.app_names)
            this.formItem.app_show = res.data.data.app_show
        })
        
  }
}
</script>

<style lang="less">

</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";


</style>
