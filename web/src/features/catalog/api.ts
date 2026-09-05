import { request } from '../../shared/api/request'
import type { Canteen } from './types'

export function getCanteens(): Promise<Canteen[]> {
  return request<Canteen[]>('/api/v1/catalog/canteens')
}
