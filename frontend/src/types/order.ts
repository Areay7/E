/** 与后端约定的平台标识（查询参数 platform） */
export type OrderPlatform = 'shopee' | 'aliexpress' | 'tiktok'

export interface OrderRow {
  orderId: string
  buyer?: string
  amount?: string | number
  currency?: string
  status?: string
  createdAt?: string
}

export interface OrdersListResult {
  list: OrderRow[]
  total: number
}

export interface DashboardSummary {
  todaySalesUsd: number
  pendingShipment: number
  openOrders: number
}
