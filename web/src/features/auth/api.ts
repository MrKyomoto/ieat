import { request } from '../../shared/api/request'
import type { User } from './types'

export function login(email: string, password: string): Promise<User> {
  return request<User>('/api/v1/auth/session', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export function logout(): Promise<void> {
  return request<void>('/api/v1/auth/session', { method: 'DELETE' })
}

export function getCurrentUser(): Promise<User> {
  return request<User>('/api/v1/auth/me')
}
