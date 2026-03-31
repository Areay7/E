import { http } from './http'

export interface LoginRequest {
  username: string
  password: string
  captcha?: string
  captchaId?: string
}

export interface LoginResponse {
  token: string
  user: {
    id: number
    username: string
    nickname: string
    email: string
    avatar: string
    role: string
    status: number
  }
}

export interface CaptchaResponse {
  imageUrl: string
  captchaId: string
}

/** POST /api/v1/auth/login */
export async function login(data: LoginRequest): Promise<LoginResponse> {
  const { data: result } = await http.post<LoginResponse>('/api/v1/auth/login', data)
  return result
}

/** GET /api/v1/auth/captcha */
export async function getCaptcha(): Promise<CaptchaResponse> {
  const response = await http.get('/api/v1/auth/captcha', {
    responseType: 'blob',
  })

  const captchaId = response.headers['x-captcha-id']
  const imageUrl = URL.createObjectURL(response.data)

  return {
    imageUrl,
    captchaId,
  }
}

/** GET /api/v1/user/current */
export async function getCurrentUser() {
  const { data } = await http.get('/api/v1/user/current')
  return data
}

/** POST /api/v1/auth/logout */
export async function logout() {
  const { data } = await http.post('/api/v1/auth/logout')
  return data
}
