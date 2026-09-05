import { useState } from 'react'
import { Avatar, Button, Drawer, Empty, Flex, Grid, Layout, Menu, Space, Tag, Typography, message } from 'antd'
import type { MenuProps } from 'antd'
import { AdminOverviewPage } from '../features/admin/AdminOverviewPage'
import { logout } from '../features/auth/api'
import type { Role, User } from '../features/auth/types'
import { CatalogPage } from '../features/catalog/CatalogPage'
import { ManagementDashboardPage } from '../features/management/ManagementDashboardPage'

const { Header, Content, Sider } = Layout
const { Text } = Typography

const roleNames: Record<Role, string> = {
  member: '普通用户',
  manager: '管理部门人员',
  admin: '平台管理员',
}

const homeKeys: Record<Role, string> = {
  member: 'community',
  manager: 'management',
  admin: 'platform',
}

const menuItems: Record<Role, MenuProps['items']> = {
  member: [
    { key: 'community', label: '食堂与窗口' },
    { key: 'notifications', label: '站内通知' },
  ],
  manager: [
    { key: 'management', label: '经营看板' },
    { key: 'feedback', label: '窗口反馈' },
  ],
  admin: [
    { key: 'platform', label: '平台概览' },
    { key: 'moderation', label: '举报处理' },
    { key: 'imports', label: '流水导入' },
  ],
}

export function AppShell({ user, onLogout }: { user: User; onLogout: () => void }) {
  const screens = Grid.useBreakpoint()
  const compact = !screens.lg
  const [selected, setSelected] = useState(homeKeys[user.role])
  const [navigationOpen, setNavigationOpen] = useState(false)
  const [messageAPI, contextHolder] = message.useMessage()

  const signOut = async () => {
    try {
      await logout()
      onLogout()
    } catch (reason) {
      messageAPI.error(reason instanceof Error ? reason.message : '退出失败')
    }
  }

  return (
    <Layout className="app-layout">
      {contextHolder}
      <Header className="app-header">
        <Flex align="center" justify="space-between" gap={16}>
          <Space size={10}>
            <Avatar className="brand-avatar">I</Avatar>
            <Text className="brand-name">IEat</Text>
          </Space>
          <Space>
            {!compact && <Text>{user.nickname}</Text>}
            <Tag>{roleNames[user.role]}</Tag>
            {compact && <Button onClick={() => setNavigationOpen(true)}>菜单</Button>}
            <Button type="text" onClick={signOut}>退出</Button>
          </Space>
        </Flex>
      </Header>
      <Drawer
        title="功能导航"
        placement="left"
        width={280}
        open={navigationOpen}
        onClose={() => setNavigationOpen(false)}
      >
        <Menu
          mode="inline"
          selectedKeys={[selected]}
          items={menuItems[user.role]}
          onClick={({ key }) => {
            setSelected(key)
            setNavigationOpen(false)
          }}
        />
      </Drawer>
      <Layout>
        <Sider className="app-sider" width={220} breakpoint="lg" collapsedWidth={compact ? 0 : 72}>
          <Menu
            mode="inline"
            selectedKeys={[selected]}
            items={menuItems[user.role]}
            onClick={({ key }) => setSelected(key)}
          />
        </Sider>
        <Content className="app-content">
          <RolePage role={user.role} selected={selected} />
        </Content>
      </Layout>
    </Layout>
  )
}

function RolePage({ role, selected }: { role: Role; selected: string }) {
  if (role === 'member' && selected === 'community') return <CatalogPage />
  if (role === 'manager' && selected === 'management') return <ManagementDashboardPage />
  if (role === 'admin' && selected === 'platform') return <AdminOverviewPage />
  return <Empty description="该模块已列入 TODO.md，等待分配开发。" />
}
