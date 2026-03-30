import axios from 'axios'

const baseURL =
  import.meta.env.VITE_API_BASE_URL !== undefined && import.meta.env.VITE_API_BASE_URL !== ''
    ? import.meta.env.VITE_API_BASE_URL
    : ''

export const http = axios.create({
  baseURL,
  timeout: 30_000,
})
