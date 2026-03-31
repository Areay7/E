import { http } from './http'

export interface PlatformStat {
  platform: string
  orderCount: number
  totalAmount: number
  pendingCount: number
  shippedCount: number
  completedCount: number
  growthRate: number
}

export interface DailySales {
  date: string
  amount: number
  count: number
}

export interface DailyOrders {
  date: string
  count: number
}

export interface ProductSales {
  productName: string
  platform: string
  quantity: number
  amount: number
}

export interface CountrySales {
  country: string
  count: number
  amount: number
}

export interface OrderStatusCount {
  status: string
  count: number
}

export interface RecentOrder {
  orderId: string
  platform: string
  amount: number
  status: string
  country: string
  orderTime: string
}

export interface SystemHealth {
  totalProducts: number
  totalInventory: number
  apiSuccessRate: number
  lastSyncTime: string
}

export interface DashboardSummary {
  // 核心指标
  todayOrders: number
  todaySales: number
  pendingShipment: number
  openOrders: number
  lowStockCount: number

  // 对比数据
  yesterdayOrders: number
  yesterdaySales: number

  // 平台统计
  platformStats: PlatformStat[]

  // 趋势数据
  salesTrend: DailySales[]
  orderTrend: DailyOrders[]

  // 排行榜
  topProducts: ProductSales[]
  topCountries: CountrySales[]

  // 状态分布
  orderStatusDist: OrderStatusCount[]

  // 实时数据
  recentOrders: RecentOrder[]

  // 系统健康
  systemHealth: SystemHealth
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
