import { useEffect, useState } from 'react'
import { Alert, Card, Empty, List, Space, Spin, Typography } from 'antd'
import { getCanteens } from './api'
import type { Canteen } from './types'
import './styles.css'

const { Title, Text, Paragraph } = Typography

export function CatalogPage() {
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
                renderItem={(foodWindow) => (
                  <List.Item>
                    <Card size="small" title={foodWindow.name} className="window-card">
                      <Paragraph type="secondary">{foodWindow.description || '暂无介绍'}</Paragraph>
                      <Text>{foodWindow.businessHours || '营业时间待补充'}</Text>
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
