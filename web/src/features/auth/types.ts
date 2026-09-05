export type Role = 'member' | 'manager' | 'admin'

export interface User {
  id: string
  email: string
  nickname: string
  role: Role
}
