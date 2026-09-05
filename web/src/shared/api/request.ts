interface APIError {
  error?: string
}

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
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
