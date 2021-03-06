/* eslint-disable */
// 该模块用来请求API网关
import axios from 'axios'
import { Message } from 'iview'
import { getToken } from '@/utils/auth'
import Cookies from 'js-cookie'


// create an axios instance
const service = axios.create({
  baseURL: process.env.WEB_ADDR, // api的base_url
  timeout: 30000 // request timeout
})

// request interceptor
service.interceptors.request.use(
  config => {
    // 设置token
    config.headers['X-Token'] = getToken()
    return config
  }, 
  error => {
    // Do something with request error
    Promise.reject(error)
})

// respone interceptor
service.interceptors.response.use(
  response => {
    // 1054表示需要重新登陆
    if (response.data.err_code == 1054) {
      Message.error({
        content: response.data.message,
        duration: 3
      })
      // 记录现在的路径，登录后恢复
      Cookies.set("lastPath", window.location.pathname)

      setTimeout(function() {
        window.location.href = "/"
      },600)
      return response
    }

    // 错误码不为0，代表发生了错误1
    if (response.data.err_code != 0) {
      Message.error({
        content: response.data.message+' : '+ response.data.err_code,
        duration: 3 
      })
      return Promise.reject(response.data.message+' : '+ response.data.err_code)
    }
    return response
  },
  error => {
    Message.error({
      content: error.message,
      duration: 3
    })
    return Promise.reject(error)
  })

export default service
