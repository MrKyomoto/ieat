import { useState } from 'react'
import { Alert, Button, Card, Form, Input, Space, Tag, Typography } from 'antd'
import { login } from './api'
import type { User } from './types'
import './styles.css'

const { Title, Text, Paragraph } = Typography

export function LoginPage({ onLogin }: { onLogin: (user: User) => void }) {
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
