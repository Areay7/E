import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const baseURL =
  import.meta.env.VITE_API_BASE_URL !== undefined && import.meta.env.VITE_API_BASE_URL !== ''
    ? import.meta.env.VITE_API_BASE_URL
    : ''

export const http = axios.create({
  baseURL,
  timeout: 30_000,
})

// 请求拦截器 - 添加token
http.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理错误
http.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response

      // 401 未授权，跳转到登录页
      if (status === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        router.push('/login')
        ElMessage.error(data?.error || '请先登录')
      }
      // 403 禁止访问
      else if (status === 403) {
        ElMessage.error(data?.error || '没有权限访问')
      }
      // 其他错误
      else if (status >= 500) {
        ElMessage.error('服务器错误，请稍后重试')
      }
    } else if (error.request) {
      ElMessage.error('网络错误，请检查网络连接')
    }

    return Promise.reject(error)
  }
)

