import axios from 'axios'

const baseUrl = 'http://8.156.76.126:11223'
const username = 'admin'
const password = 'admin@123'

const apiService = axios.create({
  baseURL: '/api', // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 10000, // request timeout
  headers: {
    'Content-Type': 'application/json'
  }
})

apiService.interceptors.request.use(config => {
  const token = btoa(`${username}:${password}`)
  config.headers.Authorization = `Basic ${token}`
  return config
})

export default apiService
