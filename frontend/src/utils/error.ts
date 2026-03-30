import axios from 'axios'

export function messageFromError(err: unknown): string {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data
    if (data && typeof data === 'object' && 'message' in data) {
      const m = (data as { message: unknown }).message
      if (typeof m === 'string') return m
    }
    return err.message
  }
  if (err instanceof Error) return err.message
  return String(err)
}
