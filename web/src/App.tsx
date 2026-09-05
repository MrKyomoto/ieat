import { useEffect, useState } from 'react'
import {
  Alert,
  Avatar,
  Button,
  Card,
  Col,
  Drawer,
  Empty,
  Flex,
  Form,
  Grid,
  Input,
  Layout,
  List,
  Menu,
  Row,
  Space,
  Spin,
  Statistic,
  Tag,
  Typography,
  message,
} from 'antd'
import type { MenuProps } from 'antd'
import { getCanteens, getCurrentUser, login, logout } from './api'
import type { Canteen, Role, User } from './api'

const { Header, Content, Sider } = Layout
const { Title, Text, Paragraph } = Typography

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

export default function App() {
  const [user, setUser] = useState<User | null>()

  useEffect(() => {
    getCurrentUser().then(setUser).catch(() => setUser(null))
  }, [])

  if (user === undefined) {
    return <CenteredSpin />
  }
  if (user === null) {
    return <LoginPage onLogin={setUser} />
  }
  return <RoleShell user={user} onLogout={() => setUser(null)} />
}

function CenteredSpin() {
  return (
    <main className="centered-page" aria-label="正在加载">
      <Spin size="large" />
    </main>
  )
}

function LoginPage({ onLogin }: { onLogin: (user: User) => void }) {
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')

  const submit = async (values: { email: string; password: string }) => {
    setSubmitting(true)
    setError('')
    try {
      onLogin(await login(values.email, values.password))
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : '登录失败')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="login-page">
      <Card className="login-card" bordered={false}>
        <Space direction="vertical" size={6} className="login-heading">
          <Tag color="green">USTC CAMPUS</Tag>
          <Title level={1}>今天吃什么？</Title>
          <Paragraph type="secondary">查看真实评价，也让每一次反馈都被食堂管理部门看见。</Paragraph>
        </Space>
        {error && <Alert className="login-error" type="error" message={error} showIcon />}
        <Form layout="vertical" requiredMark={false} onFinish={submit}>
          <Form.Item
            label="校园邮箱"
            name="email"
            rules={[{ required: true, type: 'email', message: '请输入正确的邮箱地址' }]}
          >
            <Input size="large" autoComplete="username" placeholder="name@mail.ustc.edu.cn" />
          </Form.Item>
          <Form.Item label="密码" name="password" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password size="large" autoComplete="current-password" />
          </Form.Item>
          <Button block size="large" type="primary" htmlType="submit" loading={submitting}>
            登录
          </Button>
        </Form>
        <Text className="dev-hint" type="secondary">开发账号与密码请查看项目 README。</Text>
      </Card>
    </main>
  )
}

function RoleShell({ user, onLogout }: { user: User; onLogout: () => void }) {
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
  if (role === 'member' && selected === 'community') return <CommunityHome />
  if (role === 'manager' && selected === 'management') return <ManagementHome />
  if (role === 'admin' && selected === 'platform') return <AdminHome />
  return <PendingPage />
}

function CommunityHome() {
  const [canteens, setCanteens] = useState<Canteen[]>()
  const [error, setError] = useState('')

  useEffect(() => {
    getCanteens().then(setCanteens).catch((reason) => {
      setError(reason instanceof Error ? reason.message : '读取食堂目录失败')
    })
  }, [])

  return (
    <Space direction="vertical" size={20} className="page-stack">
      <div>
        <Title level={2}>食堂与窗口</Title>
        <Text type="secondary">从窗口开始查看和分享就餐体验。</Text>
      </div>
      {error && <Alert type="error" message={error} showIcon />}
      {!canteens && !error && <Spin />}
      {canteens?.length === 0 && <Empty description="还没有食堂数据" />}
      {canteens?.map((canteen) => (
        <Card key={canteen.id} title={canteen.name} bordered={false}>
          {canteen.floors.map((floor) => (
            <section key={floor.id} className="floor-section">
              <Title level={4}>{floor.name}</Title>
              <List
                grid={{ gutter: 16, xs: 1, sm: 2, lg: 3 }}
                dataSource={floor.windows}
                locale={{ emptyText: '本层暂无窗口' }}
                renderItem={(window) => (
                  <List.Item>
                    <Card size="small" title={window.name} className="window-card">
                      <Paragraph type="secondary">{window.description || '暂无介绍'}</Paragraph>
                      <Text>{window.businessHours || '营业时间待补充'}</Text>
                    </Card>
                  </List.Item>
                )}
              />
            </section>
          ))}
        </Card>
      ))}
    </Space>
  )
}

function ManagementHome() {
  return (
    <Space direction="vertical" size={20} className="page-stack">
      <div><Title level={2}>经营看板</Title><Text type="secondary">仅展示所辖窗口的汇总信息。</Text></div>
      <Row gutter={[16, 16]}>
        {['今日净交易额', '成功交易笔数', '客单价', '待回复评价'].map((title) => (
          <Col xs={24} sm={12} xl={6} key={title}>
            <Card bordered={false}><Statistic title={title} value="—" /></Card>
          </Col>
        ))}
      </Row>
      <Alert message="经营数据接口将在 TX 与 MGMT 模块完成后接入。" type="info" showIcon />
    </Space>
  )
}

function AdminHome() {
  return (
    <Space direction="vertical" size={20} className="page-stack">
      <div><Title level={2}>平台管理</Title><Text type="secondary">维护基础数据并处理平台事务。</Text></div>
      <Row gutter={[16, 16]}>
        {['待处理举报', '最近导入批次', '食堂与窗口', '管理部门'].map((title) => (
          <Col xs={24} sm={12} xl={6} key={title}>
            <Card bordered={false}><Statistic title={title} value="—" /></Card>
          </Col>
        ))}
      </Row>
      <Alert message="管理功能入口已预留，具体任务见 TODO.md。" type="info" showIcon />
    </Space>
  )
}

function PendingPage() {
  return <Empty description="该模块已列入 TODO.md，等待分配开发。" />
}
