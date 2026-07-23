import { Card, Row, Col, Statistic, Table } from 'antd';
import { FileTextOutlined, HomeOutlined, UserOutlined, DollarOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import ReactECharts from 'echarts-for-react';
import { getDashboardStats, getOrderList, getRoomList, getTrendStats } from '@/api';
import type { Order, Room } from '@/types';

interface StatItem {
  title: string;
  value: number;
  icon: React.ReactNode;
  color: string;
  bgColor: string;
  subText: string;
}

interface TrendData {
  date: string;
  value: number;
}

export const Dashboard = () => {
  const [stats, setStats] = useState<StatItem[]>([
    { title: '今日订单', value: 0, icon: <FileTextOutlined />, color: '#5470c6', bgColor: '#ecf5ff', subText: '共 0 单' },
    { title: '使用中包间', value: 0, icon: <HomeOutlined />, color: '#91cc75', bgColor: '#f0f9eb', subText: '共 0 间空闲' },
    { title: '会员总数', value: 0, icon: <UserOutlined />, color: '#fac858', bgColor: '#fffbe6', subText: '共 0 人' },
    { title: '今日营收', value: 0, icon: <DollarOutlined />, color: '#ee6666', bgColor: '#fff1f0', subText: '累计 ¥0.00' },
  ]);
  const [orders, setOrders] = useState<Order[]>([]);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [orderTrend, setOrderTrend] = useState<TrendData[]>([]);
  const [revenueTrend, setRevenueTrend] = useState<TrendData[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [statsRes, orderRes, roomRes, orderTrendRes, revenueTrendRes] = await Promise.all([
        getDashboardStats(),
        getOrderList({ page: 1, page_size: 8 }),
        getRoomList({ page: 1, page_size: 10 }),
        getTrendStats({ type: 'order', days: 7 }),
        getTrendStats({ type: 'revenue', days: 7 }),
      ]);

      const statData = statsRes.data as any || {};
      setStats([
        { title: '今日订单', value: statData.todayOrders || 0, icon: <FileTextOutlined />, color: '#5470c6', bgColor: '#ecf5ff', subText: `共 ${orderRes.data.total} 单` },
        { title: '使用中包间', value: roomRes.data.data.filter(r => r.status === '1').length, icon: <HomeOutlined />, color: '#91cc75', bgColor: '#f0f9eb', subText: `共 ${roomRes.data.total} 间` },
        { title: '会员总数', value: statData.totalUsers || 0, icon: <UserOutlined />, color: '#fac858', bgColor: '#fffbe6', subText: '共 0 人' },
        { title: '今日营收', value: statData.todayRevenue || 0, icon: <DollarOutlined />, color: '#ee6666', bgColor: '#fff1f0', subText: `累计 ¥${(statData.totalRevenue || 0).toFixed(2)}` },
      ]);

      setOrders(orderRes.data.data);
      setRooms(roomRes.data.data);
      setOrderTrend(orderTrendRes.data.data || []);
      setRevenueTrend(revenueTrendRes.data.data || []);
    } catch (error) {
      console.error('获取数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const orderColumns = [
    { title: '订单号', dataIndex: 'orderNo', key: 'orderNo' },
    { title: '用户', dataIndex: 'userPhone', key: 'userPhone' },
    { title: '包间', dataIndex: 'code', key: 'code' },
    { title: '金额', dataIndex: 'amount', key: 'amount', render: (val: number) => <span style={{ color: '#ee6666', fontWeight: 500 }}>¥{val}</span> },
    { title: '状态', dataIndex: 'status', key: 'status', render: (status: string) => {
      const statusMap: Record<string, { text: string; color: string; bgColor: string }> = {
        'pending': { text: '待支付', color: '#fac858', bgColor: '#fffbe6' },
        'paid': { text: '已支付', color: '#5470c6', bgColor: '#ecf5ff' },
        'completed': { text: '已完成', color: '#91cc75', bgColor: '#f0f9eb' },
        'refunded': { text: '已退款', color: '#ee6666', bgColor: '#fff1f0' },
      };
      const s = statusMap[status] || { text: status, color: '#666', bgColor: '#f5f5f5' };
      return <span style={{ padding: '2px 8px', borderRadius: '4px', backgroundColor: s.bgColor, color: s.color, fontSize: '12px' }}>{s.text}</span>;
    }},
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
  ];

  const roomColumns = [
    { title: '包间名称', dataIndex: 'name', key: 'name' },
    { title: '类型', dataIndex: 'tag', key: 'tag' },
    { title: '楼层', dataIndex: 'tag', key: 'floor' },
    { title: '状态', dataIndex: 'status', key: 'status', render: (status: string) => {
      const isActive = status === '1';
      return <span style={{ padding: '2px 8px', borderRadius: '4px', backgroundColor: isActive ? '#f0f9eb' : '#f5f5f5', color: isActive ? '#91cc75' : '#999', fontSize: '12px' }}>{isActive ? '空闲' : '使用中'}</span>;
    }},
  ];

  const orderTrendOption = {
    tooltip: { trigger: 'axis' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: orderTrend.map(d => d.date), axisLine: { lineStyle: { color: '#ddd' } } },
    yAxis: { type: 'value', axisLine: { lineStyle: { color: '#ddd' } } },
    series: [{
      name: '订单数',
      type: 'line',
      smooth: true,
      data: orderTrend.map(d => d.value),
      lineStyle: { color: '#5470c6', width: 3 },
      areaStyle: {
        color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [
          { offset: 0, color: 'rgba(84, 112, 198, 0.3)' },
          { offset: 1, color: 'rgba(84, 112, 198, 0.05)' }
        ]}
      },
      symbol: 'circle',
      symbolSize: 8,
      itemStyle: { color: '#5470c6' }
    }]
  };

  const revenueTrendOption = {
    tooltip: { trigger: 'axis' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: revenueTrend.map(d => d.date), axisLine: { lineStyle: { color: '#ddd' } } },
    yAxis: { type: 'value', axisLine: { lineStyle: { color: '#ddd' } } },
    series: [{
      name: '营收(¥)',
      type: 'bar',
      data: revenueTrend.map(d => d.value),
      itemStyle: {
        color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [
          { offset: 0, color: '#ee6666' },
          { offset: 1, color: '#ff9999' }
        ]},
        borderRadius: [4, 4, 0, 0]
      }
    }]
  };

  const roomStatusOption = {
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { orient: 'horizontal', bottom: 0 },
    series: [{
      name: '包间状态',
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      label: { show: false },
      emphasis: { label: { show: true, fontSize: 14, fontWeight: 'bold' } },
      labelLine: { show: false },
      data: [
        { value: rooms.filter(r => r.status === '1').length, name: '空闲', itemStyle: { color: '#91cc75' } },
        { value: rooms.filter(r => r.status === '0').length, name: '使用中', itemStyle: { color: '#fac858' } },
      ]
    }]
  };

  return (
    <div className="dashboard">
      <div className="dashboard-header">
        <h2>仪表盘</h2>
      </div>

      <Row gutter={[16, 16]} className="stats-row">
        {stats.map((stat) => (
          <Col xs={24} sm={12} lg={6} key={stat.title}>
            <Card className="stat-card" hoverable loading={loading}>
              <div className="stat-content">
                <div className="stat-icon" style={{ backgroundColor: stat.bgColor }}>
                  <span style={{ color: stat.color }}>{stat.icon}</span>
                </div>
                <div className="stat-info">
                  <Statistic
                    title={stat.title}
                    value={stat.title === '今日营收' ? stat.value : stat.value}
                    prefix={stat.title === '今日营收' ? '¥' : ''}
                    valueStyle={{ color: stat.color, fontSize: '28px', fontWeight: 600 }}
                  />
                  <p className="stat-subtext" style={{ color: '#999', fontSize: '12px', marginTop: '4px' }}>{stat.subText}</p>
                </div>
              </div>
            </Card>
          </Col>
        ))}
      </Row>

      <Row gutter={[16, 16]} className="charts-row">
        <Col xs={24} lg={12}>
          <Card title="订单趋势" className="chart-card">
            <ReactECharts option={orderTrendOption} style={{ height: '300px' }} />
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card title="营收趋势" className="chart-card">
            <ReactECharts option={revenueTrendOption} style={{ height: '300px' }} />
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card title="包间状态分布" className="chart-card">
            <ReactECharts option={roomStatusOption} style={{ height: '250px' }} />
          </Card>
        </Col>
        <Col xs={24} lg={16}>
          <Card title="最近订单" className="chart-card">
            <Table
              columns={orderColumns}
              dataSource={orders}
              pagination={false}
              loading={loading}
              rowKey="id"
              size="small"
            />
          </Card>
        </Col>
        <Col xs={24} lg={24}>
          <Card title="包间状态" className="chart-card">
            <Table
              columns={roomColumns}
              dataSource={rooms}
              pagination={false}
              loading={loading}
              rowKey="id"
              size="small"
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};