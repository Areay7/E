import { http } from './http'
import type { OrderPlatform, OrderRow, OrdersListResult } from '@/types/order'

export interface OrderListQuery {
  page: number
  pageSize: number
  keyword?: string
  status?: string
}

function parseOrderRow(row: unknown): OrderRow {
  if (typeof row !== 'object' || row === null) {
    throw new Error('订单行格式错误：非对象')
  }
  const r = row as Record<string, unknown>
  const rawId = r.orderId ?? r.order_id ?? r.id
  if (typeof rawId !== 'string' && typeof rawId !== 'number') {
    throw new Error('订单行缺少 orderId / id')
  }
  return {
    orderId: String(rawId),
    buyer: typeof r.buyer === 'string' ? r.buyer : undefined,
    amount: typeof r.amount === 'number' || typeof r.amount === 'string' ? r.amount : undefined,
    currency: typeof r.currency === 'string' ? r.currency : undefined,
    status: typeof r.status === 'string' ? r.status : undefined,
    createdAt:
      typeof r.createdAt === 'string'
        ? r.createdAt
        : typeof r.created_at === 'string'
          ? r.created_at
          : undefined,
  }
}

function parseOrdersPayload(data: unknown): OrdersListResult {
  if (typeof data !== 'object' || data === null) {
    throw new Error('订单列表响应格式错误：非对象')
  }
  const o = data as Record<string, unknown>
  if (!Array.isArray(o.list)) {
    throw new Error('订单列表响应格式错误：缺少 list 数组')
  }
  const total = o.total
  if (typeof total !== 'number') {
    throw new Error('订单列表响应格式错误：total 须为数字')
  }
  const list = o.list.map((row) => parseOrderRow(row))
  return { list, total }
}

/**
 * GET /api/v1/orders?platform=&page=&pageSize=&keyword=&status=
 * 后端未就绪时会得到网络错误或非 2xx，由调用方展示。
 */
export async function fetchOrders(
  platform: OrderPlatform,
  query: OrderListQuery
): Promise<OrdersListResult> {
  const { data } = await http.get<unknown>('/api/v1/orders', {
    params: {
      platform,
      page: query.page,
      pageSize: query.pageSize,
      keyword: query.keyword || undefined,
      status: query.status || undefined,
    },
  })
  return parseOrdersPayload(data)
}
