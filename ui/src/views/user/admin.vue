<template>
  <div class="app-container" style="margin-top:30px">
      <Row>
          <Col span="22" offset="1">
            <Table stripe :columns="userLabels" :data="userList" class="margin-top-15"></Table>
          </Col>
      </Row>
      
  </div>
</template>

<script>
import request from '@/utils/request' 
export default {
  name: 'admin',
  data () {
    return {
      userLabels: [
            {
                title: 'ID',
                key: 'id'
            },
            {
                title: '姓名',
                key: 'name',
            },
            {
                title: '邮箱',
                key: 'email'
            },
            {
                title: '手机号',
                key: 'mobile'
            },
            {
                title: '上次登录时间',
                key: 'last_login_date'
            },
            {
                title: '登录次数',
                key: 'login_count'
            },
            {
                title: '权限',
                key: 'priv'
            },
            {
                        title: '操作',
                        key: 'action',
                        width: 150,
                        align: 'center',
                        render: (h, params) => {
                            var t = '取消管理'
                            if (params.row.priv == 'normal') {
                                t = '设为管理'
                            } else if (params.row.priv == 'super_admin') {
                                t = '转出超级管理'
                            }
                            return h('div', [
                                h('Button', {
                                    props: {
                                        type: 'primary',
                                        size: 'small'
                                    },
                                    style: {
                                        marginRight: '5px'
                                    },
                                    on: {
                                        click: () => {
                                            if (t == '设为管理') {
                                                 request({
                                                    url: '/web/setAdmin',
                                                    method: 'POST',
                                                    params: {
                                                        user_id: params.row.id
                                                    }
                                                }).then(res => {
                                                    params.row.priv = 'admin'
                                                    this.$Message.success({
                                                        content: '设置管理成功',
                                                        duration: 3 
                                                    })
                                                })
                                            } else if (t == '取消管理') {
                                                request({
                                                    url: '/web/cancelAdmin',
                                                    method: 'POST',
                                                    params: {
                                                        user_id: params.row.id
                                                    }
                                                }).then(res => {
                                                    params.row.priv = 'normal'
                                                    this.$Message.success({
                                                        content: '取消管理成功',
                                                        duration: 3 
                                                    })
                                                })
                                            }
                                        }
                                    }
                                }, t)
                                // h('Button', {
                                //     props: {
                                //         type: 'error',
                                //         size: 'small'
                                //     },
                                //     on: {
                                //         click: () => {
                                //             this.remove(params.index)
                                //         }
                                //     }
                                // }, 'Delete')
                            ]);
                        }
                    }
        ],
    userList: [
        ],
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
        url: '/web/manageUserList',
        method: 'GET',
        params: {
        }
    }).then(res => {
      this.userList = res.data.data
    })
  }
}
</script>

<style lang="less">
.ivu-modal-close { 
  visibility: hidden !important;
}

</style>

<style lang="less" scoped> 
@import "../../theme/gvar.less";
.margin-left-20 {
  color: @text-light-color
}


</style>
