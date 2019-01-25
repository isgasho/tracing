<template>
  <div class="app-container">
      <Row>
          <Col span="22" offset="1">
            <div class="header font-size-18 no-border">
                <span class="hover-cursor alerts-hover-primary" style="font-size:15px" @click="createAppAlert"><Icon type="ios-add-circle-outline font-size-20" style="margin-bottom:3px;margin-right:3px"/>新建应用告警</span>
                 <Tooltip placement="right" max-width="400">
                    <Icon type="ios-help" style="margin-bottom:5px;font-size:16px"  />
                  <div slot="content" style="padding: 15px 15px">
                     <div class="font-size-18 font-weight-500" style="line-height:25px">什么是应用告警</div>
                     <div>为指定的应用设定告警策略、告警通道等，一个应用只能设置一次告警</div>
                     <div class="font-size-18 font-weight-500" style="line-height:25px">我该怎么做</div>
                     <div>你应该先创建策略模版,然后再回来创建应用告警</div>
                     <div class="font-size-18 font-weight-500" style="line-height:25px">找不到应用？</div>
                     <div>联系APM管理员部署监控</div>
                     <div class="font-size-18 font-weight-500" style="line-height:25px">要设置的应用已经被设置？</div>
                     <div>联系之前设置的Owner，将该应用转交给你，或者联系APM管理员</div>
                  </div>
                </Tooltip>
                <span style="float:right">
                    <Select v-model="appSetting" style="width:220px" placeholder="" @on-change="setAppSetting">
                        <Option v-for="item in appSettingItems" :value="item.value" :key="item.value">{{ item.label }}</Option>
                    </Select>
                </span>
            </div>
            <Table stripe :columns="appLabels" :data="appList" class="margin-top-15" @on-row-click="editApp">
                <template slot-scope="{ row }" slot="owner">
                  {{ row.owner_name + '/'+row.owner_id}}
                </template>
                <template slot-scope="{ row }" slot="policy">
                  <policy :policyID="row.policy" :policyName="row.policy_name" class="hover-cursor"></policy>
                </template>
                <template slot-scope="{ row }" slot="channel">
                  <span v-if="row.channel=='mobile'">短信</span>
                  <span v-else>邮件</span>
                </template>
                <template slot-scope="{ row }" slot="users">
                   {{genUsers(row)}}
                </template>
            </Table>
          </Col>
      </Row>

      <Modal  v-model="handleAppVisible" :title="handleAppTitle" ok-text="提交" cancel-text="取消"  @on-ok="submitHandleApp" @on-cancel="cancelHandleApp" width="500">
          <Row>
              <Col span="21" offset="1">
                <Form :model="tempApp" :label-width="80">
                    <FormItem label="应用名">
                        <Select v-model="tempApp.name" v-if="handleAppType=='create'" style="width:220px" placeholder="选择你的应用" filterable>
                            <Option v-for="item in appNames" :value="item" :key="item">{{ item }}</Option>
                        </Select>
                        <span v-else>{{tempApp.name}}</span> 
                    </FormItem>
                    <FormItem label="策略模版">
                        <Select v-model="tempApp.policy" style="width:220px" placeholder="请选择.." filterable>
                            <Option v-for="policy in policyList" v-show="policy.owner_id==tempApp.owner_id" :value="policy.id" :key="policy.id">{{ policy.name }}</Option>
                        </Select>
                    </FormItem>

                    <FormItem label="告警通道">
                        <RadioGroup v-model="tempApp.channel">
                            <Radio label="mobile">短信</Radio>
                            <Radio label="email">邮件</Radio>
                        </RadioGroup>
                    </FormItem>

                    <FormItem label="告警用户">
                        <Select v-model="tempApp.users" multiple  filterable>
                            <Option v-for="u in userList" :value="u.id" :key="u.id">{{u.id}}/{{u.name}}</Option>
                        </Select>
                    </FormItem>

                    <FormItem label="删除" v-show="handleAppType=='edit'">
                        <Poptip
                            confirm
                            :title="'一旦删除不可恢复！确定删除应用告警 ' + tempApp.name + ' ？'"
                            ok-text="删" 
                            cancel-text="不,我不要删除,速速取消" 
                            @on-ok="confirmDeleteApp(tempApp.name)">
                            <Button type="warning" size="small">删除</Button>
                        </Poptip>
                    </FormItem>
                </Form>
              </Col>
          </Row>

    </Modal>
  </div>
</template>

<script>
import request from '@/utils/request' 
import policy from './components/policy' 
export default {
  name: 'appList',
  components: {policy},
  data () {
    return {
        appLabels: [
            {
                title: '应用名',
                key: 'name',
                width:250
            },
            {
                title: 'Owner',
                slot: 'owner',
                width:200
            },
            {
                title: '策略模版',
                slot: 'policy',
                width: 200
            },
            {
                title: '告警通道',
                slot: 'channel',
                width: 150
            },
            {
              title: '告警用户',
                slot: 'users',  
            }
        ],
        appList: [],
        handleAppType: '',
        handleAppVisible: false,
        handleAppTitle: '',

        appNames: [],
        userList: [],
        policyList: [],

        tempApp: {},

        appSetting: 1,
        appSettingItems :[
            {
                value: 1,
                label: '查看全部应用'
            },
            {
                value: 2,
                label: '我创建的应用'
            },
            {
                value: 3,
                label: '我设定的应用'
            }
        ]
    }
  },
  watch: {
  },
  computed: {
  },
  methods: {
      genUsers(app) {
          var us = []
          for (var i=0;i<app.users.length;i++) {
              us.push(app.user_names[i]+ '/' + app.users[i])
          }
          return us
      },
      setAppSetting() {
          this.loadApps()
      },
      confirmDeleteApp() {
          request({
            url: '/apm/web/deleteAppAlert',
            method: 'POST',
            params: {
                name: this.tempApp.name
            }
        }).then(res => {   
            this.loadApps()
            this.handleAppVisible = false
            this.$Message.success({
                content: '删除成功',
                duration: 3 
            })
        })
      },
      submitHandleApp() {
          if (this.handleAppType =='create') {
              request({
                    url: '/apm/web/createAppAlert',
                    method: 'POST',
                    params: {
                        app_name: this.tempApp.name,
                        policy: this.tempApp.policy,
                        channel: this.tempApp.channel,
                        users:  JSON.stringify(this.tempApp.users)
                    }
                }).then(res => {   
                    this.loadApps()
                    this.$Message.success({
                        content: '创建成功',
                        duration: 3 
                    })
                })
          } else {
              request({
                    url: '/apm/web/editAppAlert',
                    method: 'POST',
                    params: {
                        app_name: this.tempApp.name,
                        policy: this.tempApp.policy,
                        channel: this.tempApp.channel,
                        users:  JSON.stringify(this.tempApp.users)
                    }
                }).then(res => {   
                    this.loadApps()
                    this.$Message.success({
                        content: '编辑成功',
                        duration: 3 
                    })
                })
          }
      },
      cancelHandleApp() {
          this.handleAppVisible = false
      },
      editApp(app) {
          this.handleAppVisible = true
          this.handleAppTitle = '编辑应用告警'
          this.tempApp = app
          this.handleAppType = 'edit'
          console.log(this.tempApp)
      },
      createAppAlert() {
          this.handleAppTitle = '新建应用告警'
          this.handleAppType = 'create'
          this.handleAppVisible = true
          this.tempApp = {
              policy :this.policyList[0].id,
              channel: 'mobile',
              users: [this.$store.state.user.id]
          }
      },
      loadApps() {
        this.$Loading.start();
        request({
            url: '/apm/web/alertsAppList',
            method: 'GET',
            params: {
                type: this.appSetting,
            }
        }).then(res => {
            this.appList = res.data.data
             this.$Loading.finish();
        }).catch(error => {
            this.$Loading.error();
        })
      }
  },
  mounted() {
      this.loadApps()
      // 加载app名列表
         request({
            url: '/apm/web/appNames',
            method: 'GET',
            params: {
            }
        }).then(res => {   
            this.appNames = res.data.data 
        })

        request({
            url: '/apm/web/userList',
            method: 'GET',
            params: {
            }
        }).then(res => {
            this.userList = res.data.data
        })

        request({
              url: '/apm/web/queryPolicies',
              method: 'GET',
              params: {
              }
          }).then(res => {
            this.policyList = res.data.data
          })
  }
}
</script>

<style lang="less">
@import "../../theme/alerts.less";
</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";


</style>
