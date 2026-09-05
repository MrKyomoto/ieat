import { Alert, Card, Col, Row, Space, Statistic, Typography } from 'antd'

const { Title, Text } = Typography

export function ManagementDashboardPage() {
  return (
    <Space direction="vertical" size={20} className="page-stack">
      <div>
        <Title level={2}>经营看板</Title>
        <Text type="secondary">仅展示所辖窗口的汇总信息。</Text>
      </div>
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
