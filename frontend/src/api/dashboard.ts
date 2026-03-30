import { http } from './http'
import type { DashboardSummary } from '@/types/order'

function parseSummary(data: unknown): DashboardSummary {
  if (typeof data !== 'object' || data === null) {
    throw new Error('看板汇总响应格式错误：非对象')
  }
  const o = data as Record<string, unknown>
  const todaySalesUsd = o.todaySalesUsd ?? o.today_sales_usd
  const pendingShipment = o.pendingShipment ?? o.pending_shipment
  const openOrders = o.openOrders ?? o.open_orders
  if (typeof todaySalesUsd !== 'number' || typeof pendingShipment !== 'number' || typeof openOrders !== 'number') {
    throw new Error('看板汇总字段类型不符合约定（todaySalesUsd / pendingShipment / openOrders 须为数字）')
  }
  return { todaySalesUsd, pendingShipment, openOrders }
}

/** GET /api/v1/dashboard/summary */
export async function fetchDashboardSummary(): Promise<DashboardSummary> {
  const { data } = await http.get<unknown>('/api/v1/dashboard/summary')
  return parseSummary(data)
}
