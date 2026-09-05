export type Role = 'member' | 'manager' | 'admin'

export interface User {
  id: string
  email: string
  nickname: string
  role: Role
}

export interface FoodWindow {
  id: string
  externalCode: string
  name: string
  description: string
  businessHours: string
}

export interface Floor {
  id: string
  name: string
  windows: FoodWindow[]
}

export interface Canteen {
  id: string
  name: string
  floors: Floor[]
}

interface APIError {
  error?: string
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
	...init,
    credentials: 'same-origin',
    headers: init?.body ? { 'Content-Type': 'application/json', ...init.headers } : init?.headers,
  })
  if (!response.ok) {
    let body: APIError = {}
    try {
      body = (await response.json()) as APIError
    } catch {
      // Keep the server status as the fallback error.
    }
    throw new Error(body.error || `请求失败（${response.status}）`)
  }
  if (response.status === 204) {
    return undefined as T
  }
  return (await response.json()) as T
}

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

export function getCanteens(): Promise<Canteen[]> {
  return request<Canteen[]>('/api/v1/catalog/canteens')
}
