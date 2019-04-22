<template>
  <div class="app-container">
      <Row>
          <Col span="22" offset="1">
            <div class="header font-size-18">
                <span class="hover-cursor alerts-hover-primary" style="font-size:15px" @click="createPolicy"><Icon type="ios-add-circle-outline font-size-20" style="margin-bottom:3px;margin-right:3px"/>新建策略模版</span>
                 <Tooltip placement="right" max-width="400">
                    <Icon type="ios-help-circle-outline" style="margin-bottom:3px;font-size:16px"  class="alerts-color-primary" />
                  <div slot="content" style="padding: 15px 15px">
                     <div class="font-size-18 font-weight-500" style="line-height:20px">什么是策略模版</div>
                     <div>策略模版把多个告警条件组合成一个模版对象，用户可以在后续将该模版关联到具体的应用上，避免了大量重复操作</div>
                     <div class="font-size-18 font-weight-500" style="line-height:20px">允许设置哪些告警条件</div>
                     <div>针对三种类型监控：系统监控、应用监控和业务监控，都可以自定义设置</div>
                  </div>
                </Tooltip>
            </div>
            <Table stripe :columns="policyLabels" :data="policyList" class="margin-top-15" @on-row-click="editPolicy">
              <template slot-scope="{ row }" slot="owner">
                  {{ row.owner_name + '/'+row.owner_id}}
              </template>
              <template slot-scope="{ row }" slot="alerts">
                 <policy1 :policyID="row.id" :policyName="row.alerts.length" class="hover-cursor"></policy1>
              </template>
            </Table>
          </Col>
      </Row>

      <Modal  :mask-closable="false" v-model="handlePolicyVisible" :title="handlePolicyTitle" ok-text="提交" cancel-text="取消"  @on-ok="submitHandlePolicy" @on-cancel="cancelHandlePolicy" fullscreen>
          <Row>
              <Col span="24">
                <Form :model="tempPolicy" :label-width="120" label-position="left">
                    <FormItem label="策略模版名称">
                        <div class="margin-left-20">
                          <div class="right-meta">为你的模版选择一个准确、简洁的名称，只支持英文字母和-</div>
                           <Input style="width:400px;" class="right-body" v-model="tempPolicy.name" placeholder="e.g. tf56pay-gateway" :autofocus="handleType=='create'"></Input>
                        </div>    
                    </FormItem>

                    <FormItem label="设置监控项">
                        <div class="margin-left-20">
                           <div class="right-meta">左边进行设置，成功后会添加到右边区域；点击右边区域中的监控项，可进行修改</div>
                           <div class="right-body" >
                             <Row>
                               <Col span="9">
                                 <div style="border-right:1px solid #ddd;padding: 10px 25px">
                                    <div class="font-size-12">选择监控类型</div>
                                    <div class="margin-top-10">
                                      <RadioGroup v-model="tempAlert.type" type="button">
                                          <Radio label="apm" @click.native="selPolicyType('apm')">应用监控</Radio>
                                          <Radio label="system" @click.native="selPolicyType('system')">系统监控</Radio>
                                      </RadioGroup>
                                    </div>

                                    <div class="font-size-12 margin-top-20">定义告警阈值</div>
                                      <div class="alert-setting no-border">
                                        <span class="label">监控指标</span>
                                        <Select :value="tempAlert.name" style="width:250px" > 
                                          <!-- <OptionGroup label="关键监控指标"> -->
                                              <Option v-for="alert in alertItems[policyType]" v-show="filterOption(alert)" @click.native="selAlert(alert)"  :value="alert.name" :key="alert.name">{{ alert.label }}</Option>
                                          <!-- </OptionGroup> -->
                                          <!-- <OptionGroup label="其它监控指标">
                                              <Option v-for="item in alertItems[policyType]" @click.native="selAlert(item)" v-show="item.key==false" :value="item.name" :key="item.name">{{ item.label }}</Option>
                                          </OptionGroup> -->
                                          
                                        </Select>
                                        <Tooltip placement="left" max-width="400" v-show="tempAlert.help != undefined">
                                          <Icon type="ios-help-circle-outline" style="margin-top:-2px;font-size:16px" class="margin-left-5"  />
                                          <div slot="content" style="padding: 15px 15px">
                                              <div>{{tempAlert.help}}</div>
                                          </div>
                                      </Tooltip>
                                      </div>

                                      <div class="alert-setting" style="margin-top:15px">
                                        <span class="label">比较方式</span>
                                        <span style="margin-left:10px">
                                          <span v-if="tempAlert.compare==1"> > </span>
                                          <span v-else-if="tempAlert.compare==1"> = </span>
                                          <span v-else> &lt; </span>
                                        </span>
                                       
                                      </div>

                                      <div class="alert-setting">
                                        <span class="label">错误的HTTP CODE</span>
                                        <Input style="width:100px;margin-bottom:12px" class="right-body" v-model="tempAlert.value" placeholder=""></Input>
                                        <span class="label">{{tempAlert.unit}}</span>
                                      </div>

                                      <div class="alert-setting">
                                        <span class="label">设定阈值</span>
                                        <Input style="width:100px;margin-bottom:12px" class="right-body" v-model="tempAlert.value" placeholder=""></Input>
                                        <span class="label">{{tempAlert.unit}}</span>
                                      </div>

                                      <div class="alert-setting">
                                        <span class="label">持续时间</span>
                                         <InputNumber :max="5" :min="1" :step="2" v-model="tempAlert.duration"></InputNumber><span class="label">分钟</span>
                                      </div>
                                      

                                      <div class="font-size-12 margin-top-20">操作</div>
                                      <div class="margin-left-10 font-size-18">
                                        <Tooltip placement="bottom" max-width="400" content="添加该项">
                                          <Icon type="md-add" class="meta-color hover-cursor" @click="addAlert"/>
                                        </Tooltip>
                                         
                                        <Tooltip placement="bottom" max-width="400" content="清空该项">
                                           <Icon type="md-close" class="color-error margin-left-5 hover-cursor" @click="clearAlert"/>
                                        </Tooltip>
                                       
                                      </div>
                                  </div>
                               </Col>
                               <Col span="10" offset="1">
                                <div style="padding: 10px 25px">
                                  <div class="alert-setting" v-show="isAlertsVisible('apm')">
                                    <div class="font-size-12">应用监控</div>
                                    <div class="margin-left-10">
                                      <alert v-for="a in tempPolicy.alerts" v-show="a.type=='apm'" :key="a.value" :alert="a"  @click.native="editAlert(a)" class="margin-left-10 hover-cursor" style="background-color:#9cd9e7;padding:4px 6px;font-size:12px;border-radius:4px;"></alert>
                                    </div>
                                  </div>

                                  <div class="alert-setting" v-show="isAlertsVisible('system')">
                                    <div class="font-size-12 margin-top-10">系统监控</div>
                                     <div class="margin-left-10">
                                      <alert v-for="a in tempPolicy.alerts" v-show="a.type=='system'" :key="a.value" :alert="a"  @click.native="editAlert(a)" class="margin-left-10 hover-cursor" style="background-color:#efda83;padding:4px 6px;font-size:12px;border-radius:4px;"></alert>
                                    </div>
                                  </div>
                                </div>
                               </Col>
                             </Row>
                           </div>
                        </div>    
                    </FormItem>

                    <FormItem label="删除模版" v-show="handleType=='edit'">
                        <Poptip
                            confirm
                            :title="'一旦删除不可恢复！确定删除模版 ' + tempPolicy.name + ' ？'"
                            ok-text="删" 
                            cancel-text="不,我不要删除,速速取消" 
                            @on-ok="confirmDeletePolicy(tempPolicy.id)">
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
import alert from './components/alert' 
import policy1 from './components/policy' 
export default {
  name: 'policy',
  components: {alert,policy1},
  watch: {
  },
  computed: {
  },
  methods: {
      clearAlert() {
        this.tempAlert = {
          name: '',
          type: '',
          duration: 0,
          compare: 1,
          value: 0,
          unit: ''
        }
      },
      editAlert(a) {
        // 从tempPolicy中删除该alert
        for (var i=0;i<this.tempPolicy.alerts.length;i++) {
          if (this.tempPolicy.alerts[i].name == a.name) {
            this.tempPolicy.alerts.splice(i,1)
          }
        }
        this.tempAlert = a
      } ,  
      filterOption(alert) {
        //若当前alert已经在tempPolicy中，则不显示
        for (var i=0;i<this.tempPolicy.alerts.length;i++) {
          if (this.tempPolicy.alerts[i].name == alert.name) {
            return false
          }
        }

        return true
      },
      isAlertsVisible(type) {
        for (var i=0;i<this.tempPolicy.alerts.length;i++) {
          if (this.tempPolicy.alerts[i].type == type) {
            return true
          }
        }

        return false
      },
      addAlert() {
        this.tempPolicy.alerts.push(this.tempAlert)
        this.tempAlert =  {
          name: '',
          type: '',
          duration: 0,
          compare: 1,
          value: 0,
          unit: ''
        }
      },
      selAlert(alert) {
        // 设置tempAlert
        this.tempAlert = alert 
        this.tempAlert.type = this.policyType
      },
      selPolicyType(tp) {
        this.policyType = tp
        this.tempAlert = {}
      },
      confirmDeletePolicy(id) {
         request({
                url: '/apm/web/deletePolicy',
                method: 'POST',
                params: {
                    id: id
                }
            }).then(res => {
                this.loadPolicys()
                this.handlePolicyVisible = false
                this.$Message.success({
                    content: '删除成功',
                    duration: 3 
                })
            })
      },
      editPolicy(policy) {
          this.handleType = 'edit'
          this.handlePolicyTitle = '编辑策略模版'
          this.tempPolicy = policy
          this.handlePolicyVisible = true 
      },
      submitHandlePolicy() {
        if (this.handleType == 'create') {
          request({
              url: '/apm/web/createPolicy',
              method: 'POST',
              params: {
                policy: JSON.stringify(this.tempPolicy)
              }
          }).then(res => {
            this.tempPolicy =  {
                name: '',
                alerts: []
            }
            this.loadPolicys()
            this.$Message.success({
                content: '创建成功',
                duration: 3 
            })
          })
        } else {
          request({
              url: '/apm/web/editPolicy',
              method: 'POST',
              params: {
                policy: JSON.stringify(this.tempPolicy)
              }
          }).then(res => {
            this.tempPolicy =  {
                name: '',
                alerts: []
            }
            this.loadPolicys()
            this.$Message.success({
                content: '编辑成功',
                duration: 3 
            })
          })
        }
      },
      cancelHandlePolicy() {
          this.handlePolicyVisible = false
      },
      createPolicy() {
        this.handlePolicyTitle = '创建策略模版'
        this.handlePolicyVisible = true
        this.handleType = 'create'
        this.tempPolicy.name = ''
        this.tempPolicy.alerts =  []
        // 设置默认显示的监控项
        this.tempAlert = this.alertItems.apm[0]
        this.tempAlert.type = 'apm'
      },
      loadPolicys() {
        this.$Loading.start();
        request({
              url: '/apm/web/queryPolicies',
              method: 'GET',
              params: {
              }
          }).then(res => {
            this.policyList = res.data.data
            this.$Loading.finish();
          }).catch(error => {
            this.$Loading.error();
        })
      }
  },
  mounted() {
    this.loadPolicys()
  },
   data () {
    return {
        handlePolicyTitle: '创建策略模版',
        handlePolicyVisible: false,
        policyLabels: [
            {
                title: '模版名',
                key: 'name'
            },
            {
                title: 'Owner',
                slot: 'owner',
            },
            {
                title: '监控项',
                key: 'alerts',
                slot: 'alerts'
            }
        ],
        policyList: [],
  

        tempPolicy: {
            name: '',
            alerts: []
        },
        tempAlert : {
          name: '',
          type: '',
          duration: 0,
          compare: 1,
          value: 0,
          unit: '',
          keys: []
        },
        handleType: 'create',
        policyType: 'apm',
        alertItems: {
          apm: [
            {
              name: 'apm.apdex.count',
              label: '综合健康指数Apdex',
              compare: 3,
              unit: '',
              duration: 1,
              keys : [],
              value: 0.8
            },
            {
              name: 'apm.http_code.ratio',
              label: '错误HTTP CODE比率',
              compare: 1,
              unit: '%',
              duration: 1,
              keys : [],
              value: 10,
              help: '指定的http code占所有请求的比例'
            },
            {
              name: 'apm.http_code.count',
              label: '错误HTTP CODE次数',
              compare: 1,
              unit: '次',
              duration: 1,
              keys : [],
              value: 10,
              help: '制定的http code发生次数'
            },
            {
              name: 'apm.api_error_.atio',
              label: '接口错误率',
              compare: 1,
              unit: '%',
              duration: 1,
              keys : [],
              value: 10
            },
            {
              name: 'apm.sql_error.ratio',
              label: 'sql错误率',
              compare: 1,
              unit: '%',
              duration: 1,
              keys : [],
              value: 10
            },
            {
              name: 'apm.api.duration',
              label: '接口平均耗时',
              compare: 1,
              unit: 'ms',
              duration: 1,
              keys : [],
              value: 10000
            },
             {
              name: 'apm.sql.duration',
              label: 'sql平均耗时',
              compare: 1,
              unit: 'ms',
              duration: 1,
              keys : [],
              value: 10000
            },
            {
              name: 'apm.jvm_fullgc.count',
              label: 'JVMFullGC报警',
              compare: 1,
              unit: '次',
              duration: 1,
              keys : [],
              value: 2
            },
            {
              name: 'apm.api.count',
              label: '接口访问次数',
              compare: 1,
              unit: '次',
              duration: 1,
              keys : [],
              value: 3000
            }
          ],
          system: [
            {
              name: 'system.cpu_used.ratio',
              label: 'cpu使用率',
              compare: 1,
              unit: '%',
              duration: 1,
              keys : [],
              value: 80
            },
            {
              name: 'system.load.count',
              label: '系统Load',
              compare: 1,
              unit: '',
              duration: 1,
              keys : [],
              value: 4
            },
            {
              name: 'system.mem_used.ratio',
              label: '内存使用率',
              compare: 1,
              unit: '%',
              duration: 1,
              keys : [],
              value: 90
            },
            {
              name: 'system.disk_used.ratio',
              label: '硬盘使用率',
              compare: 1,
              unit: '%',
              duration: 1,
              keys : [],
              value: 80
            },
            {
              name: 'system.syn_recv.count',
              label: 'sync_recv数',
              compare: 1,
              unit: '个',
              duration: 1,
              keys : [],
              value: 10000
            },
             {
              name: 'system.time_wait.count',
              label: 'time_wait数',
              compare: 1,
              unit: '个',
              duration: 1,
              keys : [],
              value: 10000
            },
            {
              name: 'system.ioutil.ratio',
              label: 'diskio利用率',
              compare: 1,
              unit: '%',
              duration: 1,
              keys : [],
              value: 90
            },
            {
              name: 'system.ifstat_out.speed',
              label: '网络out速度',
              compare: 1,
              unit: 'MB/S',
              duration: 1,
              keys : [],
              value: 100
            },
            {
              name: 'system.close_wait.count',
              label: 'close_wait数',
              compare: 1,
              unit: '个',
              duration: 1,
              keys : [],
              value: 5000
            },
            {
              name: 'system.ifstat_in.speed',
              label: '网络in速度',
              compare: 1,
              unit: 'MB/S',
              duration: 1,
              keys : [],
              value: 100
            },
            {
              name: 'system.estab.count',
              label: '建立长链接数',
              compare: 1,
              unit: '个',
              duration: 1,
              keys : [],
              value: 5000
            }
          ]
        }
    }
  }
}
</script>

<style lang="less">
@import "../../theme/alerts.less";
</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";
.alert-setting {
  margin-top:8px;
  .label {
    font-size:10px;
    margin-left:10px;
    margin-right:15px;
    color: @meta-color
  }
}

</style>
