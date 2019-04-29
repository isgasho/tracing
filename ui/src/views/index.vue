<template>
  <div class="index">
    <Row class="nav">
            <Col class="navbar-brand" span="2" offset="3">
                <span class="dog"><img src="https://datadog-prod.imgix.net/img/dd_logo_70x75.png?fm=png&amp;auto=format&amp;lossless=1" style="height: 50px;margin-top:10px"></span>
            </Col>
            <Col class="navbar-collapse hover-cursor"  style="color:#777" span="14" >
                <span>产品介绍</span>
                <span>文档说明</span>
                <span>联系我们</span>
            </Col>
              <span class="login-button" @click="gotoApmAdmin">进入监控后台</span>
    </Row>

    <div class="header">
        <div id="datadog-video-container" class="section-video section-video-play text-center" style="top: -20px;"><video poster="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEWvAP9YXMgKAAAAAXRSTlMAQObYZgAAAApJREFUeJxjYgAAAAYAAzY3fKgAAAAASUVORK5CYII=" id="datadog-video" playsinline="" autoplay="" muted="" loop="" class="bgvid" preload="auto" data-play="true" style="background-image: url(../assets/dd_homepage_still_opt.png);"><source src="../assets/dd_homepage_3.mp4"></video></div>
          <div class="container">
            <Row>
                <Col span="12" offset="6">
                    <div style="font-size:48px;line-height:48px;font-weight:300;color:black;text-align:center;">现代化的监控&分析平台</div> 
                    <div style="margin-top:15px;font-size:24px;font-weight:300;color:black;text-align:center;">查看你的系统、应用、业务监控，all in one!</div>  
                </Col>
            </Row>
        </div>
    </div>

    <Row style="margin-top:800px;height:533px;">
        <Col span="15">
            <img width="100%" srcset="https://datadog-prod.imgix.net/img/blog/monitoring-mean-stack-applications-with-datadog/monitoring-mean-stack.png?auto=format&amp;w=767&amp;dpr=1 1x,
                                    https://datadog-prod.imgix.net/img/blog/monitoring-mean-stack-applications-with-datadog/monitoring-mean-stack.png?auto=format&amp;w=767&amp;dpr=2 2x" alt="" base-img="https://datadog-prod.imgix.net/img/blog/monitoring-mean-stack-applications-with-datadog/monitoring-mean-stack.png">
        </Col>
        <Col span="9" style="background-color: rgb(238,235,12);height:100%;padding-top:60px;padding-left:30px">
            <div style="font-size:21px;font-weight:600;color:black;">完全可视化你的应用</div>
            <p style="font-size:17px;color:black">- 微服务端到端可视化</p>
            <p style="font-size:17px;color:black">- 应用拓扑可视化</p>
            <p style="font-size:17px;color:black">- 系统指标可视化</p>
            <p style="font-size:17px;color:black">- 日志可视化</p>
            <p style="font-size:17px;color:black">- 所有数据无缝关联</p>
        </Col>
    </Row>

    <Row style="margin-top:60px">
        <Col span="5" offset="4">
            <div style="font-size:22px;font-weight: 600;line-height: 40px;color:black">自定义应用埋点</div>
            <div style="font-size:17px;font-weight: 400;line-height: 30px;color:black;margin-top:10px;">支持通过日志或者SDK进行自定义埋点</div>
  
            <p style="font-size:17px;color:black;margin-top:15px">- 使用客户端SDK进行自定义埋点</p>
            <p style="font-size:17px;color:black">- 完全支持Opentracing标准</p>
            <p style="font-size:17px;color:black">- 在日志中输出监控链路ID</p>
            <p style="font-size:17px;color:black">- 埋点、日志跟APM自动采集的数据进行关联</p>
            <p style="font-size:17px;color:black">- 实现业务链路监控</p>
        </Col>
        <Col span="6" offset="1">
                 <img  alt="" height="500" class="img-responsive   lazyloaded" data-src="https://datadog-prod.imgix.net/img/Homepage-Screencaps-API-v3.png?ch=Width&amp;fit=max&amp;fm=png&amp;auto=format&amp;lossless=1" src="https://datadog-prod.imgix.net/img/Homepage-Screencaps-API-v3.png?ch=Width&amp;fit=max&amp;fm=png&amp;auto=format&amp;lossless=1">
        </Col>
    </Row>
  </div>
</template>

<script>
import request from '@/utils/request' 
const config = require('@/../config')
export default {
  name: 'index',
  data () {
    return {

    }
  },
  watch: {
    $route() {
      // 当路由变化时，调用此函数
    },
    world(val) {
      // world发生变化时
    }
  },
  methods: {
      gotoApmAdmin() {
        // let host = config.dev.host + ":" + config.dev.port
        let url = "http://sitetest.tf56.com/openssoWeb/opensso/login?clientNo=1" + "&redirectUrl=" + window.location.origin + "/ui/callback"
        window.location.href = url
      }
  },
  mounted() {
  }
}
</script>

<style lang="less">
</style>

<style lang="less" scoped> 
@import "../theme/gvar.less";
.nav {
    border-bottom: 1px solid transparent;
    box-shadow: 0 0 6px 6px rgba(0,0,0,.07);
    height:65px;
    display: flex;
    align-items: center;
    align-content: center;
    position: fixed;
    top: 0;
    left:0;
    right:0;
    z-index: 1030;
    background:white;
}

.navbar-collapse {
    span {
        border: 1px solid #fff;
        padding: 8.5px 9px 9.5px;
        font-size: 17px;
        font-weight: 600;
        letter-spacing: .4px;
        border-radius: 5px;
        margin-left:10px
    }
    span:hover {
            border-color: #c8c8c8;
            color : #444;
                transition: all .2s;
    }
}

.login-button {
    border: 1px solid #7749a3;
    background: #fff;
    color: #7749a3;
    margin: 0;
    border-radius: 4px;
    transition: all .5s;
    padding: 11px 20px 10px;
    font-weight: 600;
}
.login-button:hover {
    background: linear-gradient(45deg,#3023ae,#c024b6);
    color: white;
}

.header {
    margin-top: 150px;
    .section-video {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 900.25px;
        max-width: 1700px;
        margin: 0 auto;
        right: 0;
    }
    .container {
        margin-right: auto;
        margin-left: auto;
        padding-left: 15px;
        padding-right: 15px;
        z-index: 100;
    }
}

    .footer {
        color: hsla(0,0%,100%,.5);
        background: #110617;
    }
</style>
