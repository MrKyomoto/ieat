import { Alert, Card, Col, Row, Space, Statistic, Typography } from 'antd'

const { Title, Text } = Typography

export function AdminOverviewPage() {
  return (
    <Space direction="vertical" size={20} className="page-stack">
      <div>
        <Title level={2}>平台管理</Title>
        <Text type="secondary">维护基础数据并处理平台事务。</Text>
      </div>
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
