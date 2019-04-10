<template>
  <div class="login-container">
    <video  poster="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEWvAP9YXMgKAAAAAXRSTlMAQObYZgAAAApJREFUeJxjYgAAAAYAAzY3fKgAAAAASUVORK5CYII=" id="datadog-video" playsinline="" autoplay="" muted="" loop="" class="bgvid" preload="auto" data-play="true" style="background-image: url(https://datadog-prod.imgix.net/img/dd_homepage_still_opt.png?w=1700&amp;h=956&amp;fit=crop);height:900px"><source src="https://don08600y3gfm.cloudfront.net/video/dd_homepage_3.mp4"></video>
    <Form class="login-form" autoComplete="on" label-position="left">
      <div class="title-container">
        <h3 class="title">传化监控平台</h3>
      </div>


      <Button type="primary" style="width:100%;margin-bottom:10px;margin-top:80px;background:#2d8cf0;border-color:#2d8cf0"  @click.native.prevent="handleLogin">进入监控</Button>

      <div class="tips">
        <span>当前使用SSO登录</span>
      </div>
    </Form> 

    <div class="index-bg-footer"><img src="../../assets/login-bg2.png" style="width:100%"></div>
  </div>
</template>

<script>
/* eslint-disable */
import request from '@/utils/request' 
const config = require('@/../config')
export default {
  name: 'login',
  data() {
    return {
    }
  },
  methods: {
    handleLogin() {
      let url = "http://sitetest.tf56.com/openssoWeb/opensso/login?clientNo=1" + "&redirectUrl=http://" + window.location.origin+ "/apm/ui/callback"
      window.location.href = url
    }
  },
  created() {
 
  }
}
</script>

<style  lang="less">
/* reset element-ui css */
</style>

<style lang="less" scoped>
.login-container {
  width: 100%;
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  background-color: #fff;
  background-image: url(../../assets/login-bg1.jpg);
  background-repeat: no-repeat;
  background-size: cover;
  background-position: 50%;
  overflow: hidden;
  .login-form {
    position: absolute;
    left: 50px;
    width: 460px;
    padding: 35px 35px 15px 35px;
    margin-top: 150px;
    z-index: 99;
  }
  .tips {
    font-size: 13px;
    color: rgb(180,180,180);
    margin-bottom: 10px;
    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }
  .svg-container {
    padding: 6px 5px 6px 15px;
    color: #889aa4;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
    &_login {
      font-size: 20px;
    }
  }
  .title-container {
    position: relative;
    .title {
      font-size: 36px;
      font-weight: bold;
    }
    .title::after {
        content: " ";
        border-bottom-style: solid;
        border-bottom-width: 3px;
        display: block;
        width: 35px;
        color: #2d8cf0;
        padding-top: 15px;
        margin-left: 6px
    }
    .set-language {
      color: #fff;
      position: absolute;
      top: 5px;
      right: 0px;
    }
  }
  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: #889aa4;
    cursor: pointer;
    user-select: none;
  }
  .thirdparty-button {
    position: absolute;
    right: 35px;
    bottom: 28px;
  }
}
.index-bg-footer {
  width: 260px;
  height: 90px;
  position: absolute;
  bottom: 0;
  left: 0;
  img {
    z-index: 1
  }
}
 @media only screen and (max-width: 800px) {
     .login-form { 
        margin-top: 30px !important;
     }
     .index-bg-footer {
       display: none !important
     }
}
</style>