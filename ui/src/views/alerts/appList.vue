<template>
  <div class="app-container">
      <Row>
          <Col span="22" offset="1">
            <div class="header font-size-18 no-border">
                <span class="hover-cursor alerts-hover-primary" style="font-size:15px" @click="createAppAlert"><Icon type="ios-add-circle-outline font-size-20" style="margin-bottom:3px;margin-right:3px"/>新建应用告警</span>
                 <Tooltip placement="right" max-width="400">
                    <Icon type="ios-help-circle-outline" style="margin-bottom:3px;font-size:16px" class="alerts-color-primary" />
                    <div slot="content" style="padding: 15px 15px">
                        <div class="font-size-18 font-weight-500" style="line-height:20px">什么是应用告警</div>
                        <div>为指定的应用设定告警策略、告警通道等，一个应用只能设置一次告警</div>
                        <div class="font-size-18 font-weight-500" style="line-height:20px">我该怎么做</div>
                        <div>你应该先创建策略模版,然后再回来创建应用告警</div>
                        <div class="font-size-18 font-weight-500" style="line-height:20px">找不到应用？</div>
                        <div>联系APM管理员部署监控</div>
                        <div class="font-size-18 font-weight-500" style="line-height:20px">要设置的应用已经被设置？</div>
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

      <Modal  v-model="handleAppVisible" :title="handleAppTitle" ok-text="提交" cancel-text="取消"  @on-ok="submitHandleApp" @on-cancel="cancelHandleApp" width="800">
          <Row>
              <Col span="21" offset="1">
                <Form :model="tempApp" :label-width="100">
                    <FormItem label="应用名">
                        <Select v-model="tempApp.name" v-if="handleAppType=='create'" style="width:280px" placeholder="选择你的应用" filterable @on-change="selApp">
                            <Option v-for="item in appNames" :value="item" :key="item">{{ item }}</Option>
                        </Select>
                        <span v-else>{{tempApp.name}}</span> 
                    </FormItem>
                    <FormItem label="策略模版">
                        <Select v-model="tempApp.policy" style="width:280px" placeholder="请选择.." filterable> 
                            <Option v-for="policy in policyList" :value="policy.id" :key="policy.id">{{ policy.name }}</Option>
                        </Select>
                         <policy :policyID="tempApp.policy" policyName="view" class="hover-cursor font-size-12 margin-left-5 alerts-color-primary"></policy>
                    </FormItem>

                    <FormItem label="告警通道">
                        <RadioGroup v-model="tempApp.channel">
                            <Radio label="mobile">短信</Radio>
                            <Radio label="email">邮件</Radio>
                        </RadioGroup>
                    </FormItem>

                    <FormItem label="告警用户">
                        <Select v-model="tempApp.users" multiple  filterable style="width:280px">
                            <Option v-for="u in userList" :value="u.id" :key="u.id">{{u.id}}/{{u.name}}</Option>
                        </Select>
                    </FormItem>
                    

                    <!-- 特殊告警规则设定 -->
                    <FormItem label="特殊告警" v-show="tempApp.policy != '' && tempApp.name != undefined">
                        <div class="right-meta">为指定接口设定告警规则，此规则将覆盖策略模版</div>
                        <div class="right-body" style="border: 1px solid #eee; padding:10px" >
                            <div class="font-size-12">设定API</div>
                            <div class="margin-top-10">
                                <Select v-model="tempApi"   filterable placeholder="设定应用API" style="width:280px">
                                    <Option v-for="api in apis" :value="api" :key="api">{{api}}</Option>
                                </Select>
                            </div>
                            
                            <div v-show="tempApi != undefined && tempApi != ''">
                                <div class="font-size-12 margin-top-20">添加告警项</div>
                                <div class="margin-top-10">
                                        <span class="no-border">
                                            <Select v-model="tempAlert.key"   filterable placeholder="选择监控项" style="width:130px">
                                                <Option v-for="alert in alertItems" :value="alert.name" :key="alert.name">{{alert.label}}</Option>
                                            </Select>
                                        </span>
                                    <InputNumber  style="width:100px"  :min="0" :max="alertNumberMax(tempAlert.key)" v-model="tempAlert.value"  placeholder="告警值.."></InputNumber>
                                    <Icon type="md-add"  class="margin-left-10 alerts-color-primary hover-cursor" @click="addAlert"/>
                                </div>
                            </div>

                            <div>
                                <div class="font-size-12 margin-top-20">已设定告警项</div>
                                <div class="margin-top-10">
                                    <Tag v-for="alert in apiAlerts()" @click.native="updateAlert(alert)">{{getAlertLabel(alert.key) + '/' +alert.value}}</Tag>
                                </div>
                            </div>

                             <div>
                                <div class="font-size-12 margin-top-20">已设定API</div>
                                <div class="margin-top-10">
                                    <Tag v-for="api in setAlertApis()" closable @on-close="delApiAlerts(api)" @click.native="selApi(api)">{{api}}</Tag>
                                </div>
                            </div>
                        </div> 
                    </FormItem>
                   
                    <FormItem label="危险区域" v-show="handleAppType=='edit'">
                        <Poptip
                            confirm
                            :title="'一旦删除不可恢复！确定删除应用告警 ' + tempApp.name + ' ？'"
                            ok-text="删" 
                            cancel-text="不,我不要删除,速速取消" 
                            @on-ok="confirmDeleteApp(tempApp.name)">
                            <Button type="warning" size="small">删除应用</Button>
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
import alertItems from './data.js';
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
            },
            {
              title: '更新时间',
                key: 'update_date',  
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
        ],
        apis: [],
        tempApi: '',
        tempAlert: {
            value: 0
        },
        alertItems: []
    }
  },
  watch: {
  },
  computed: {
  },
  methods: {
      alertNumberMax(key) {
          switch (key) {
              case 'apm.api_error.ratio':
                  return 100
                  break;
              case 'apm.api.duration':
                  return 300000
                  break;
              case 'apm.api.count':
                  return 10000000
                  break;
              default:
                  return 0
                  break;
          }
      },
      selApi(api) {
          this.tempApi = api
      },
      getAlertLabel(key) {
          for (var i=0;i<this.alertItems.length;i++) {
              if (this.alertItems[i].name == key) {
                  return this.alertItems[i].label
              }
          }

          return key
      },
      delApiAlerts(api) {
          if (this.tempApp.api_alerts == undefined) {
              return []
          }

           for (var i=0;i<this.tempApp.api_alerts.length;i++) {
              if (this.tempApp.api_alerts[i].api == api) {
                  this.tempApp.api_alerts.splice(i,1)
              }
          }
      },
      setAlertApis() {
          if (this.tempApp.api_alerts == undefined) {
              return []
          }
          var apis = []
          for (var i=0;i<this.tempApp.api_alerts.length;i++) {
              apis.push(this.tempApp.api_alerts[i].api)
          }

          return apis
      },
      updateAlert(alert) {
        this.tempAlert = {
            key: alert.key,
            value: alert.value
        }
        for (var i=0;i<this.tempApp.api_alerts.length;i++) {
            if (this.tempApp.api_alerts[i].api == this.tempApi) {
                for (var j=0;j<this.tempApp.api_alerts[i].alerts.length;j++) {
                    if (this.tempApp.api_alerts[i].alerts[j].key == alert.key) {
                        this.tempApp.api_alerts[i].alerts.splice(j,1)
                    }
                }
            }
        }
      },
      apiAlerts() {
          if (this.tempApp.api_alerts == undefined) {
              return []
          }
          for (var i=0;i<this.tempApp.api_alerts.length;i++) {
              if (this.tempApp.api_alerts[i].api == this.tempApi) {
                  return this.tempApp.api_alerts[i].alerts
              }
          }

          return []
      },
      selApp(app) {
        if (app == undefined) {
            return 
        }
          this.apiList(app)
      },
      addAlert() {
        for (var i=0;i<this.tempApp.api_alerts.length;i++) {
            if (this.tempApp.api_alerts[i].api == this.tempApi) {
                for (var j=0;j<this.tempApp.api_alerts[i].alerts.length;j++) {
                    if (this.tempApp.api_alerts[i].alerts[j].key == this.tempAlert.key) {
                        this.$Message.warning('告警项已存在')
                        return 
                    }
                }
                var alert = {
                    key: this.tempAlert.key,
                    value: this.tempAlert.value
                }
                this.tempApp.api_alerts[i].alerts.push(alert)
                this.tempAlert = {
                    value: 0
                }

                console.log(this.tempApp.api_alerts)
                return 
            }
        }

        var alert = {
            key: this.tempAlert.key,
            value: this.tempAlert.value
        }
        this.tempAlert = {
            value:0
        }
        this.tempApp.api_alerts.push({
            api: this.tempApi,
            alerts: [alert]
        })

        console.log(this.tempApp.api_alerts)
      },
      apiList(app) {
        request({
            url: '/apm/web/appApis',
            method: 'GET',
            params: {
                app_name: app,
            }
        }).then(res => {
            this.apis = res.data.data
        })
      },
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
                        users:  JSON.stringify(this.tempApp.users),
                        api_alerts: JSON.stringify(this.tempApp.api_alerts)
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
                        users:  JSON.stringify(this.tempApp.users),
                        api_alerts: JSON.stringify(this.tempApp.api_alerts)
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
          this.apiList(app.name)
      },
      createAppAlert() {
          this.handleAppTitle = '新建应用告警'
          this.handleAppType = 'create'
          this.handleAppVisible = true
          this.apis = []
          this.tempApp = {
              policy :this.policyList[0].id,
              channel: 'mobile',
              users: [this.$store.state.user.id],
              api_alerts: []
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
            for (var i=0;i<this.appList.length;i++) {
                this.appList[i].api_alerts = JSON.parse(this.appList[i].api_alerts)
            }
             this.$Loading.finish();
        }).catch(error => {
            this.$Loading.error();
        })
      }
  },
  mounted() {
      // 设置特殊告警项
      for (var i=0;i<alertItems.apm.length;i++) {
          if (alertItems.apm[i].name== 'apm.api_error.ratio' || alertItems.apm[i].name == 'apm.api.count' || alertItems.apm[i].name == 'apm.api.duration') {
              this.alertItems.push(alertItems.apm[i])
          }
      }
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
 input {
                border-top:none !important;
                border-left:none !important;
                border-right:none !important;
                border-radius: 0 !important
              }
  
            .ivu-input:focus {
                border-color: @primary-color;
                outline: 0;
                box-shadow: none;
            }
</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";


</style>
