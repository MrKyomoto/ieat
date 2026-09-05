import { useEffect, useState } from 'react'
import { Spin } from 'antd'
import { getCurrentUser } from '../features/auth/api'
import { LoginPage } from '../features/auth/LoginPage'
import type { User } from '../features/auth/types'
import { AppShell } from './AppShell'

export default function App() {
  const [user, setUser] = useState<User | null>()

  useEffect(() => {
    getCurrentUser().then(setUser).catch(() => setUser(null))
  }, [])

  if (user === undefined) {
    return (
      <main className="centered-page" aria-label="正在加载">
        <Spin size="large" />
      </main>
    )
  }
  if (user === null) {
    return <LoginPage onLogin={setUser} />
  }
  return <AppShell user={user} onLogout={() => setUser(null)} />
}
