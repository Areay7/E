import { http } from './http'

export interface PlatformStat {
  platform: string
  orderCount: number
  totalAmount: number
}

export interface DailySales {
  date: string
  amount: number
  count: number
}

export interface ProductSales {
  productName: string
  quantity: number
  amount: number
}

export interface OrderStatusCount {
  status: string
  count: number
}

export interface DashboardSummary {
  todayOrders: number
  todaySales: number
  pendingShipment: number
  openOrders: number
  lowStockCount: number
  platformStats: PlatformStat[]
  salesTrend: DailySales[]
  topProducts: ProductSales[]
  orderStatusDist: OrderStatusCount[]
}

export interface SalesReport {
  totalOrders: number
  totalAmount: number
  avgAmount: number
  paidOrders: number
  shippedCount: number
}

/** GET /api/v1/dashboard/summary */
export async function fetchDashboardSummary(): Promise<DashboardSummary> {
  const { data } = await http.get<DashboardSummary>('/api/v1/dashboard/summary')
  return data
}

/** GET /api/v1/dashboard/sales-report */
export async function fetchSalesReport(params: {
  startDate: string
  endDate: string
  platform?: string
}): Promise<SalesReport> {
  const { data } = await http.get<SalesReport>('/api/v1/dashboard/sales-report', { params })
  return data
}
