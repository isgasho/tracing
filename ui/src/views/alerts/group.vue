<template>
  <div class="app-container">
      <Row>
          <Col span="22" offset="1">
            <div class="header font-size-18">
                <span class="hover-cursor alerts-hover-primary" style="font-size:15px"  @click="createGroup"><Icon type="ios-add-circle-outline font-size-20" style="margin-bottom:3px;margin-right:3px"/>新建用户组</span>
                <Tooltip placement="right" max-width="400">
                    <Icon type="ios-help-circle-outline" style="margin-bottom:3px;font-size:16px" class="alerts-color-primary" />
                  <div slot="content" style="padding: 15px 15px">
                     <div class="font-size-18 font-weight-500" style="line-height:20px">什么是用户组</div>
                     <div>用户组可以关联到应用，当应用发生告警时，会通过用户组设置的告警通道通知组内成员(包括Owner)</div>
                  </div>
                </Tooltip>
            </div>
            <Table stripe :columns="groupLabels" :data="groupList" class="margin-top-15" @on-row-click="editGroup" on-row-dblclick="deleteGroup">
                <template slot-scope="{ row }" slot="owner">
                  {{ row.owner_name + '/'+row.owner_id}}
                </template>
                <template slot-scope="{ row }" slot="channel">
                  <span v-if="row.channel=='mobile'">短信</span>
                  <span v-else>邮件</span>
                </template>
                <template slot-scope="{ row }" slot="users">
                  {{ row.users.length }}
                </template>
            </Table>
          </Col>
      </Row>

      <Modal  :mask-closable="false" v-model="handleGroupVisible" :title="handleGroupTitle" ok-text="提交" cancel-text="取消"  @on-ok="submitHandleGroup" @on-cancel="cancelHandleGroup" width="500">
          <Row>
              <Col span="21" offset="1">
                <Form :model="tempGroup" :label-width="80">
                    <FormItem label="组名">
                        <Input v-model="tempGroup.name" placeholder="为你的用户组起一个简洁、直观的名称.." ></Input>
                    </FormItem>
                    <FormItem label="告警通道">
                        <RadioGroup v-model="tempGroup.channel">
                            <Radio label="mobile">短信</Radio>
                            <Radio label="email">邮件</Radio>
                        </RadioGroup>
                    </FormItem>
                    <FormItem label="组员">
                        <Select v-model="tempGroup.users" multiple  filterable>
                            <Option v-show="u.id!=$store.state.user.id" v-for="u in userList" :value="u.id" :key="u.id">{{u.id}}/{{u.name}}</Option>
                        </Select>
                    </FormItem>
                    <FormItem label="删除组" v-show="handleType=='edit'">
                        <Poptip
                            confirm
                            :title="'一旦删除不可恢复！确定删除组 ' + tempGroup.name + ' ？'"
                            ok-text="删" cancel-text="不,我不要删除,速速取消" 
                            @on-ok="confirmDeleteGroup(tempGroup.id)">
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
export default {
  name: 'group',
  data () {
    return {
        handleGroupTitle: '创建组',
        handleGroupVisible: false,
        groupLabels: [
            {
                title: '组名',
                key: 'name'
            },
            {
                title: 'Owner',
                slot: 'owner'
            },
            {
                title: '告警通道',
                slot: 'channel'
            },
             {
                title: '组员数',
                slot: 'users'
            },
        ],
        groupList: [],
        userList : [],

        tempGroup: {
            id :'',
            name: '',
            channel: 'mobile',
            users: []
        },
        handleType: 'create'
    }
  },
  watch: {
  },
  computed: {
  },
  methods: {
      confirmDeleteGroup(id) {
          request({
                url: '/apm/web/deleteGroup',
                method: 'POST',
                params: {
                    id: id
                }
            }).then(res => {
                 this.loadGroups()
                this.handleGroupVisible = false
                this.$Message.success({
                    content: '删除成功',
                    duration: 3 
                })
            })
      },
      editGroup(group) {
          this.handleType = 'edit'
          this.handleGroupTitle = '编辑组'
          this.tempGroup = group
          this.handleGroupVisible = true 
      },
      submitHandleGroup() {
          if (this.handleType == 'create') {
            request({
                url: '/apm/web/createGroup',
                method: 'POST',
                params: {
                    name: this.tempGroup.name,
                    channel: this.tempGroup.channel,
                    users: JSON.stringify(this.tempGroup.users)
                }
            }).then(res => {
                this.tempGroup = {
                    id:'',
                    name: '',
                    channel: 'mobile',
                    users: []
                }
                this.loadGroups()
                this.handleGroupVisible = false
                this.$Message.success({
                    content: '创建成功',
                    duration: 3 
                })
            })
          } else {
            request({
                url: '/apm/web/editGroup',
                method: 'POST',
                params: {
                    id : this.tempGroup.id,
                    name: this.tempGroup.name,
                    channel: this.tempGroup.channel,
                    users: JSON.stringify(this.tempGroup.users)
                }
            }).then(res => {
                this.handleGroupVisible = false
                this.loadGroups()
                this.$Message.success({
                    content: '更新成功 : ' + this.tempGroup.name,
                    duration: 3 
                })
                this.tempGroup = {
                    name: '',
                    channel: 'mobile',
                    users: []
                }
            })
          }
      },
      cancelHandleGroup() {
          this.handleGroupVisible = false
          this.tempGroup = {
                    name: '',
                    channel: 'mobile',
                    users: []
                }
      },
      createGroup() {
        this.handleGroupTitle = '创建组'
        this.handleGroupVisible = true
        this.handleType = 'create'
      },
      loadGroups() {
        request({
            url: '/apm/web/queryGroups',
            method: 'GET',
            params: {
            }
        }).then(res => {
        this.groupList = res.data.data
        for (var i=0;i<this.groupList.length;i++) {
            this.groupList[i].user_count = this.groupList[i].users.length
        }
        })
      }
  },
  mounted() {
    request({
        url: '/apm/web/userList',
        method: 'GET',
        params: {
        }
    }).then(res => {
      this.userList = res.data.data
    })
    this.loadGroups()
  }
}
</script>

<style lang="less">
@import "../../theme/alerts.less";
</style>

<style lang="less" scoped> 



</style>
