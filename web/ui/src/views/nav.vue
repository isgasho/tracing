<template>
  <div>
      <Row class="nav" >
        <ul class="product-switcher" :class="{'open':isOpen}">
            <li class="notched"  @click="switchProduct">
                <a href="#"  title="Switch to another New Relic product">
                    <div style="color:white;font-weight:bold;margin-left:35px;">OpenAPM</div>
                </a>
            </li>
            <li class="product ng-scope insights active">
            </li>
            <li class="product ng-scope apm not-active">
                <a href="/apm/ui/list">应用监控</a>
            </li>
            <li class="product ng-scope browser not-active">
                业务监控
            </li>
            <li class="product ng-scope system not-active">
                系统监控
            </li>
            <li class="product ng-scope alerts not-active">
                <a href="/apm/ui/alerts">告警平台</a>
            </li>
            <li class="product ng-scope infrastructure not-active">
                使用文档
            </li>
        </ul>
        <span class="hover-cursor">
            <span style="color:white;float:right;margin-top:6px;margin-right:29px;font-size:15px">
               
             <Dropdown>
                {{$store.state.user.name}}/{{$store.state.user.id}}
                <DropdownMenu slot="list">
                    <DropdownItem>个人设置</DropdownItem>
                    <DropdownItem v-show="$store.state.user.priv!='normal'" @click.native="goAdmin">管理面板</DropdownItem>
                    <DropdownItem @click.native="logout">退出登录</DropdownItem>
                </DropdownMenu>
            </Dropdown>    
            </span>
            <span style="color:white;float:right;margin-top:6px;margin-right:9px;font-size:15px">
                <img alt="Icon-unknown_user" class="gravatar" height="30" src="https://gravatar.com/avatar/13eb272664a43d2a962cb7c1e778ddab?s=34&amp;d=https://rpm.newrelic.com/images/v5/template/icon-unknown_user.png?a848d37" width="30">
            </span>
        </span>
        <span  :style="{'background':getBG()}" style="float:right;color:white;margin-top:.4px;margin-right:40px;padding:6px 10px">
            {{getMeta()}}
        </span>
      </Row>
      <router-view></router-view>
  </div>
</template> 
 
<script>
export default {
  name: 'Nav',
  data () {
    return {
        isOpen:false
    }
  }, 
  watch: {
  }, 
  computed: {
  },
  mounted() {
  },
  methods: {
    getBG() {
        var p = window.location.pathname
        var routes = this.$router.options.routes 
        for (var i=0;i<routes.length;i++) {
            if (routes[i].path == p) {
                return routes[i].bg
            }
            for (var j=0;j<routes[i].children.length;j++) {
                if (routes[i].children[j].path==p) {
                    return routes[i].children[j].bg
                }
                if (routes[i].children[j].children != undefined) {
                    for (var k=0;k<routes[i].children[j].children.length;k++) {
                        if (routes[i].children[j].children[k].path==p) {
                            return routes[i].children[j].children[k].bg
                        }
                    }
                }
            }
        }
    },
    getMeta() {
        var p = window.location.pathname
        var routes = this.$router.options.routes 
        for (var i=0;i<routes.length;i++) {
            if (routes[i].path == p) {
                return routes[i].meta
            }
            for (var j=0;j<routes[i].children.length;j++) {
                if (routes[i].children[j].path==p) {
                    return routes[i].children[j].meta
                }
                if (routes[i].children[j].children != undefined) {
                    for (var k=0;k<routes[i].children[j].children.length;k++) {
                        if (routes[i].children[j].children[k].path==p) {
                            return routes[i].children[j].children[k].meta
                        }
                    }
                }
            }
        }
    },
    goAdmin() {
        console.log(this.$router)
        this.$router.push('/apm/ui/admin')
        this.$store.dispatch('setPage', '管理面板')
    },
    switchProduct() {
        this.isOpen = !this.isOpen
    },
    logout() {
      this.$store.dispatch('Logout').then(() => {
        this.$router.push('/apm/ui/login') // In order to re-instantiate the vue-router object to avoid bugs
      }).catch(error => {
        // 登出错误，登陆数据已经清除，返回登陆页面
        this.$router.push('/apm/ui/login')
      })
    },
  },
}
</script>

<style lang="less">
</style>

<style lang="less" scoped> 
@import "../theme/gvar.less";
// .nav {
//     top: 0;
//     width: 100%;
//     z-index: 999;
//     // position: fixed;
    
//     .display-in-small {
//       display : none;
//     }

//      @media only screen and (max-width: 600px) {
//       .display-in-large {
//         display: none !important
//       }
//       .display-in-small {
//         display: inherit !important;
//       }
//     }
  
// }
.nav {
    // width: 224px;
    height: 41px;
    -webkit-transition: all 250ms ease-in-out;
    transition: all 250ms ease-in-out;
    // position: fixed;
    top: 6px;
    left: 0;
    z-index: 100;
    background-color: #474747;
    // box-shadow: 0 1px 1px rgba(0, 0, 0, 0.5);
}
.product-switcher ul {
    list-style: none;
    margin-bottom: 0;
}

.product-switcher .notched, .product-switcher .product {
    width: 185px;
    // height: 41px;
    position: absolute;
    top: 0;
    left: 0;
    box-shadow: 2px 0 3px rgba(0, 0, 0, 0.4);
    clip: rect(auto, 493.33333px, auto, -10px);
    padding: 4.5px 0;
}

.product-switcher .notched {
    -webkit-transition: left 250ms ease-in-out, width 250ms ease-in-out, opacity 250ms ease-in-out;
    transition: left 250ms ease-in-out, width 250ms ease-in-out, opacity 250ms ease-in-out;
    width: 164px;
    z-index: 1027;
    background-color: #474747;
    opacity: 1;
}

.product-switcher .notched a, .product-switcher .product a {
    height: 41px;
    color:white;
}

.product-switcher .notched a {
    right: -60px;
}

.product-switcher .product-logo {
    margin-left: 20px;
    height: 27px;
    background-size: contain;
    background-repeat: no-repeat;
    background-position: left center;
}

.product-switcher .notched .product-logo {
    background: url("../assets/logo.png") no-repeat top left;
    background-size: contain;
    width: 140px;
}

.product-switcher .notched::after {
    width: 24px;
    height: 41px;
    content: '';
    position: absolute;
    top: 0;
    right: -24px;
    z-index: 1027;
    background: transparent url(../assets/logo2.png) no-repeat top right;
    background-size: contain;
}

.product-switcher .insights {
    -webkit-transition: left 140ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    transition: left 140ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    background-color: #F2BC56;
    z-index: 1026;
    left: -35px;
    width: 224px;
}
.product-switcher .apm {
    -webkit-transition: left 180ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    transition: left 180ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    background-color: #348899;
    z-index: 1025;
    left: 9px;
}
.product-switcher .browser {
    -webkit-transition: left 220ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    transition: left 220ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    background-color: #F28F20;
    z-index: 1024;
    left: 14px;
}
.product-switcher .alerts {
    -webkit-transition: left 300ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    transition: left 300ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    background-color: #7a5aa6;
    z-index: 1022;
    left: 24px;
}
.product-switcher .plugins {
    -webkit-transition: left 380ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    transition: left 380ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    background-color: #8CC641;
    z-index: 1020;
    left: 34px;
}
.product-switcher .infrastructure {
    -webkit-transition: left 340ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    transition: left 340ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    background-color: #226582;
    z-index: 1021;
    left: 29px;
}
.product-switcher .system {
    -webkit-transition: left 260ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    transition: left 260ms cubic-bezier(0.455, 0.03, 0.515, 0.955);
    background-color: #00B6D8;
    z-index: 1023;
    left: 19px;
}


.product-switcher.open {
    width: 100%;
    z-index: 1100;
    text-align:center;
    font-weight:bold;
    color:white;
}
// .product-switcher.open .notched {
//     opacity: 0;
//     z-index: 1093;
// }
.product-switcher.open .apm {
    z-index: 1107;
    box-shadow: none; 
    left: 188px;
}
.product-switcher.open .browser {
    z-index: 1106;
    box-shadow: none;
    left: 373px;
}
.product-switcher.open .system {
    z-index: 1105;
    box-shadow: none;
    left: 558px;
}
.product-switcher.open .alerts {
    z-index: 1104;
    box-shadow: none;
    left: 743px;
}
.product-switcher.open .infrastructure {
    z-index: 1103;
    box-shadow: none;
    left: 928px;
}
.product-switcher.open .plugins {
    z-index: 1102;
    box-shadow: none;
    left: 1100px;
}




@media only screen and (max-width: 1280px) {
    .product-switcher.open {
        height: 123px;
    }
    .product-switcher.open .apm {
        width: 33.33333%;
        top: 0px;
        left: 33.33333%;
    }
    .product-switcher.open .browser {
        width: 33.33333%;
        top: 0px;
        left: 66.66667%;
    }
    .product-switcher.open .alerts {
        width: 33.33333%;
        top: 41px;
        left: 33.33333%;
    }
    .product-switcher.open .plugins {
        width: 33.33333%;
        top: 82px;
        left: 0%;
    }
    .product-switcher.open .infrastructure {
        width: 33.33333%;
        top: 41px;
        left: 66.66667%;
    }
    .product-switcher.open .system {
        width: 33.33333%;
        top: 41px;
        left: 0%;
    }
}





</style>
