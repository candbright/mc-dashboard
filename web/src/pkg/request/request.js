import axios from 'axios'

const apiService = axios.create({
  baseURL: '/api', // url = base url + request url
  timeout: 10000, // request timeout
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
apiService.interceptors.request.use(config => {
  // 从localStorage获取token
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器
apiService.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          // token过期或无效，清除token并跳转到登录页
          localStorage.removeItem('token')
          window.location.href = '/login'
          break
        case 403:
          // 权限不足
          console.error('没有权限访问该资源')
          break
        case 500:
          // 服务器错误
          console.error('服务器错误')
          break
        default:
          console.error('请求错误:', error.response.data)
      }
    }
    return Promise.reject(error)
  }
)

export default apiService
