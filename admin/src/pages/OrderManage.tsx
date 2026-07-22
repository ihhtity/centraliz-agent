import { Table, Button, Form, Input, Select, message, Row, Col } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import { getOrderList } from '@/api';
import type { Order } from '@/types';

const { Option } = Select;

export const OrderManage = () => {
  const [data, setData] = useState<Order[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '订单编号', dataIndex: 'orderNo', key: 'orderNo' },
    { title: '订单码', dataIndex: 'code', key: 'code' },
    { title: '支付码', dataIndex: 'payCode', key: 'payCode' },
    { title: '类型', dataIndex: 'type', key: 'type' },
    { title: '模式', dataIndex: 'mode', key: 'mode' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '金额', dataIndex: 'amount', key: 'amount' },
    { title: '时长', dataIndex: 'duration', key: 'duration' },
    { title: '价格', dataIndex: 'price', key: 'price' },
    { title: '押金', dataIndex: 'deposit', key: 'deposit' },
    { title: '支付金额', dataIndex: 'payPrice', key: 'payPrice' },
    { title: '支付方式', dataIndex: 'payType', key: 'payType' },
    { title: '用户电话', dataIndex: 'userPhone', key: 'userPhone' },
    { title: '开始时间', dataIndex: 'startTime', key: 'startTime' },
    { title: '结束时间', dataIndex: 'endTime', key: 'endTime' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
  ];

  const fetchData = async (params: { page?: number; page_size?: number; status?: string; code?: string } = {}) => {
    setLoading(true);
    try {
      const res = await getOrderList({ page: params.page || 1, page_size: params.page_size || 10, status: params.status, code: params.code });
      setData(res.data.data);
      setTotal(res.data.total);
    } catch (error) {
      message.error('获取订单列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const [searchForm] = Form.useForm();

  const handleSearch = () => {
    const values = searchForm.getFieldsValue();
    fetchData({ ...values, page: 1 });
  };

  const handleReset = () => {
    searchForm.resetFields();
    fetchData();
  };

  return (
    <div className="page-container">
      <div className="page-header">
        <h2>订单管理</h2>
      </div>

      <Form className="search-form" form={searchForm} layout="inline">
        <Form.Item name="code">
          <Input placeholder="订单编号" prefix={<SearchOutlined />} />
        </Form.Item>
        <Form.Item name="status">
          <Select placeholder="状态" allowClear>
            <Option value="pending">待支付</Option>
            <Option value="paid">已支付</Option>
            <Option value="completed">已完成</Option>
            <Option value="refunded">已退款</Option>
          </Select>
        </Form.Item>
        <Form.Item>
          <Button type="primary" onClick={handleSearch}>搜索</Button>
          <Button onClick={handleReset} style={{ marginLeft: 8 }}>重置</Button>
        </Form.Item>
      </Form>

      <div className="table-container">
        <Table
          columns={columns}
          dataSource={data}
          pagination={{ total, pageSize: 10, showSizeChanger: true, showQuickJumper: true }}
          loading={loading}
          rowKey="id"
          scroll={{ x: 1600 }}
        />
      </div>
    </div>
  );
};