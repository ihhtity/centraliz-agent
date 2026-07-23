import { Table, Button, Form, Input, Select, message, Modal, Tag, Space, Spin, Row, Col } from 'antd';
import { SearchOutlined, EditOutlined, DeleteOutlined, PlusOutlined, EyeOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import { getOrderList, getOrderDetail, batchDeleteOrder, batchUpdateOrder, createOrder, updateOrder } from '@/api';
import { CustomPagination } from '@/components/CustomPagination';
import { ExportButton } from '@/components/ExportButton';
import type { Order } from '@/types';

const { Option } = Select;

const statusColors: Record<string, string> = {
  '未完成': 'orange',
  '进行中': 'blue',
  '已完成': 'green',
  '申请退款': 'yellow',
  '已退款': 'red',
  '拒绝退款': 'gray',
};

export const OrderManage = () => {
  const [data, setData] = useState<Order[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [showSearch, setShowSearch] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [isBatchModalVisible, setIsBatchModalVisible] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [batchForm] = Form.useForm();
  const [form] = Form.useForm();
  const [editId, setEditId] = useState<number | null>(null);
  const [detailVisible, setDetailVisible] = useState(false);
  const [currentOrder, setCurrentOrder] = useState<Order | null>(null);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '订单编号', dataIndex: 'orderNo', key: 'orderNo' },
    { title: '订单码', dataIndex: 'code', key: 'code' },
    { title: '支付码', dataIndex: 'payCode', key: 'payCode' },
    { title: '类型', dataIndex: 'type', key: 'type' },
    { title: '模式', dataIndex: 'mode', key: 'mode' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={statusColors[status] || 'default'}>{status}</Tag>
      ),
    },
    { title: '金额', dataIndex: 'amount', key: 'amount', render: (val: number) => <span>¥{val}</span> },
    { title: '时长', dataIndex: 'duration', key: 'duration' },
    { title: '价格', dataIndex: 'price', key: 'price', render: (val: number) => <span>¥{val}</span> },
    { title: '押金', dataIndex: 'deposit', key: 'deposit', render: (val: number) => <span>¥{val}</span> },
    { title: '支付金额', dataIndex: 'payPrice', key: 'payPrice', render: (val: number) => <span>¥{val}</span> },
    { title: '支付方式', dataIndex: 'payType', key: 'payType' },
    { title: '用户电话', dataIndex: 'userPhone', key: 'userPhone' },
    { title: '开始时间', dataIndex: 'startTime', key: 'startTime', render: (time: string) => formatTime(time) },
    { title: '结束时间', dataIndex: 'endTime', key: 'endTime', render: (time: string) => formatTime(time) },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (time: string) => formatTime(time) },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: Order) => (
        <Space>
          <Button className="action-btn-detail" icon={<EyeOutlined />} onClick={() => viewDetail(record)} size="small">详情</Button>
          <Button className="action-btn-edit" icon={<EditOutlined />} onClick={() => handleEdit(record)} size="small">编辑</Button>
          <Button className="action-btn-delete" icon={<DeleteOutlined />} onClick={() => deleteOrderItem(record.id)} size="small">删除</Button>
        </Space>
      ),
    },
  ];

  const fetchData = async (params: { page?: number; page_size?: number; status?: string; code?: string; orderNo?: string; userPhone?: string; payType?: string } = {}) => {
    setLoading(true);
    try {
      const res = await getOrderList({ 
        page: params.page || currentPage, 
        page_size: params.page_size || pageSize, 
        status: params.status, 
        code: params.code,
        order_no: params.orderNo,
        user_phone: params.userPhone,
        pay_type: params.payType,
      });
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

  const viewDetail = async (order: Order) => {
    setDetailLoading(true);
    try {
      const res = await getOrderDetail(order.id);
      setCurrentOrder(res.data);
      setDetailVisible(true);
    } catch (error) {
      console.error(error);
      message.error('获取详情失败');
    } finally {
      setDetailLoading(false);
    }
  };

  const deleteOrderItem = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个订单吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteOrder({ ids: [id.toString()] });
          message.success('删除成功');
          fetchData();
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  const handleBatchDelete = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要删除的订单');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 个订单吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteOrder({ ids: selectedRowKeys.map(k => k.toString()) });
          message.success('批量删除成功');
          setSelectedRowKeys([]);
          fetchData();
        } catch (error) {
          message.error('批量删除失败');
        }
      },
    });
  };

  const handleAdd = () => {
    setEditId(null);
    form.resetFields();
    setIsModalVisible(true);
  };

  const handleEdit = (record: Order) => {
    setEditId(record.id);
    form.setFieldsValue(record);
    setIsModalVisible(true);
  };

  const handleBatchEdit = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要编辑的订单');
      return;
    }
    batchForm.resetFields();
    setIsBatchModalVisible(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editId) {
        await updateOrder(editId, values);
        message.success('更新成功');
      } else {
        await createOrder(values);
        message.success('创建成功');
      }
      setIsModalVisible(false);
      fetchData();
    } catch (error) {
      message.error('提交失败');
    }
  };

  const handleBatchSubmit = async () => {
    try {
      const values = await batchForm.validateFields();
      await batchUpdateOrder({ ids: selectedRowKeys.map(k => k.toString()), data: values });
      message.success('批量更新成功');
      setIsBatchModalVisible(false);
      setSelectedRowKeys([]);
      fetchData();
    } catch (error) {
      message.error('批量更新失败');
    }
  };

  const handleImport = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.xlsx,.xls';
    input.onchange = async (e: Event) => {
      const target = e.target as HTMLInputElement;
      const file = target.files?.[0];
      if (!file) return;
      const formData = new FormData();
      formData.append('file', file);
      try {
        const response = await fetch('/admin/order/import', {
          method: 'POST',
          headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` },
          body: formData,
        });
        const result = await response.json();
        if (result.code === 200) {
          message.success('导入成功');
          fetchData();
        } else {
          message.error(result.message || '导入失败');
        }
      } catch (error) {
        message.error('导入失败');
      }
    };
    input.click();
  };

  const exportHeaders: Record<string, string> = {
    id: 'ID',
    orderNo: '订单编号',
    code: '订单码',
    payCode: '支付码',
    type: '类型',
    mode: '模式',
    status: '状态',
    tag: '订单标签',
    amount: '数量',
    duration: '时长',
    price: '价格',
    deposit: '押金',
    payPrice: '支付金额',
    payType: '支付方式',
    refundPrice: '退款金额',
    remark: '备注',
    userPhone: '用户电话',
    merchPhone: '商家电话',
    startTime: '开始时间',
    endTime: '结束时间',
    payTime: '支付时间',
    createdAt: '创建时间',
    updatedAt: '更新时间',
  };

  const handleTableChange = (pagination: { current: number; pageSize: number }) => {
    setCurrentPage(pagination.current);
    setPageSize(pagination.pageSize);
    fetchData({ page: pagination.current, page_size: pagination.pageSize });
  };

  const [searchForm] = Form.useForm();

  const handleSearch = () => {
    const values = searchForm.getFieldsValue();
    setCurrentPage(1);
    fetchData({ ...values, page: 1 });
  };

  const handleReset = () => {
    searchForm.resetFields();
    setCurrentPage(1);
    fetchData();
  };

  const rowSelection = {
    selectedRowKeys,
    onChange: (keys: React.Key[]) => setSelectedRowKeys(keys),
    onSelect: (record: Order, selected: boolean) => {
      if (selected) {
        setSelectedRowKeys([...selectedRowKeys, record.id]);
      } else {
        setSelectedRowKeys(selectedRowKeys.filter(k => k !== record.id));
      }
    },
  };

  // 格式化时间
  const formatTime = (time: string) => {
    if (!time) return '-'
    return time.replace('T', ' ').substring(0, 19)
  };

  return (
    <div className="page-container">
      <div className="page-header">
        <h2>订单管理</h2>
        <div className="action-buttons">
          <Button className="action-btn-search" size="small" icon={<SearchOutlined />} onClick={() => setShowSearch(!showSearch)}>{showSearch ? '收起搜索' : '搜索'}</Button>
          <ExportButton data={data} filename="订单列表" headers={exportHeaders} />
          <Button className="action-btn-import" size="small" icon={<SearchOutlined />} onClick={handleImport}>导入</Button>
          <Button className="action-btn-edit" size="small" icon={<EditOutlined />} onClick={handleBatchEdit} disabled={selectedRowKeys.length === 0}>编辑({selectedRowKeys.length})</Button>
          <Button className="action-btn-delete" size="small" icon={<DeleteOutlined />} onClick={handleBatchDelete} disabled={selectedRowKeys.length === 0}>删除({selectedRowKeys.length})</Button>
          <Button className="action-btn-add" size="small" icon={<PlusOutlined />} onClick={handleAdd}>添加</Button>
        </div>
      </div>

      {showSearch && (
        <Form className="search-form" form={searchForm} layout="inline">
          <Form.Item name="orderNo">
            <Input placeholder="订单编号" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="code">
            <Input placeholder="订单码" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="payCode">
            <Input placeholder="支付码" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="userPhone">
            <Input placeholder="用户电话" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="type">
            <Select placeholder="订单类型" allowClear>
              <Option value="normal">普通</Option>
              <Option value="vip">VIP</Option>
            </Select>
          </Form.Item>
          <Form.Item name="mode">
            <Select placeholder="模式" allowClear>
              <Option value="hour">按时长</Option>
              <Option value="day">按天</Option>
              <Option value="month">按月</Option>
            </Select>
          </Form.Item>
          <Form.Item name="payType">
            <Select placeholder="支付方式" allowClear>
              <Option value="alipay">支付宝</Option>
              <Option value="wechat">微信</Option>
            </Select>
          </Form.Item>
          <Form.Item name="status">
            <Select placeholder="状态" allowClear>
              <Option value="未完成">未完成</Option>
              <Option value="进行中">进行中</Option>
              <Option value="已完成">已完成</Option>
              <Option value="申请退款">申请退款</Option>
              <Option value="已退款">已退款</Option>
              <Option value="拒绝退款">拒绝退款</Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" onClick={handleSearch}>搜索</Button>
            <Button onClick={handleReset} style={{ marginLeft: 8 }}>重置</Button>
          </Form.Item>
        </Form>
      )}

      <div className="table-container">
        <Table
          rowSelection={rowSelection}
          columns={columns}
          dataSource={data}
          pagination={false}
          loading={loading}
          rowKey="id"
          scroll={{ x: 1600 }}
          onRow={(record) => ({
            onClick: () => {
              if (selectedRowKeys.includes(record.id)) {
                setSelectedRowKeys(selectedRowKeys.filter(k => k !== record.id));
              } else {
                setSelectedRowKeys([...selectedRowKeys, record.id]);
              }
            },
          })}
        />
        <CustomPagination
          total={total}
          current={currentPage}
          pageSize={pageSize}
          onChange={(page, pageSize) => handleTableChange({ current: page, pageSize })}
        />
      </div>

      <Modal
        title={editId ? '编辑订单' : '添加订单'}
        open={isModalVisible}
        onOk={handleSubmit}
        onCancel={() => setIsModalVisible(false)}
        okText="确定"
        cancelText="取消"
        className="form-modal"
        width={600}
      >
        <Form form={form} layout="vertical">
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="orderNo" label="订单编号" rules={[{ required: true, message: '请输入订单编号' }]}>
                <Input placeholder="请输入订单编号" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="code" label="订单码">
                <Input placeholder="请输入订单码" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="payCode" label="支付码">
                <Input placeholder="请输入支付码" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="type" label="类型">
                <Select placeholder="请选择类型">
                  <Option value="normal">普通</Option>
                  <Option value="vip">VIP</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="mode" label="模式">
                <Select placeholder="请选择模式">
                  <Option value="hour">按时长</Option>
                  <Option value="day">按天</Option>
                  <Option value="month">按月</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="status" label="状态">
                <Select placeholder="请选择状态">
                  <Option value="未完成">未完成</Option>
                  <Option value="进行中">进行中</Option>
                  <Option value="已完成">已完成</Option>
                  <Option value="申请退款">申请退款</Option>
                  <Option value="已退款">已退款</Option>
                  <Option value="拒绝退款">拒绝退款</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="tag" label="订单标签">
                <Input placeholder="请输入订单标签" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="payType" label="支付方式">
                <Select placeholder="请选择支付方式">
                  <Option value="alipay">支付宝</Option>
                  <Option value="wechat">微信</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="amount" label="数量">
                <Input type="number" placeholder="请输入数量" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="duration" label="时长">
                <Input type="number" placeholder="请输入时长" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="price" label="价格">
                <Input type="number" placeholder="请输入价格" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="deposit" label="押金">
                <Input type="number" placeholder="请输入押金" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="payPrice" label="支付金额">
                <Input type="number" placeholder="请输入支付金额" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="refundPrice" label="退款金额">
                <Input type="number" placeholder="请输入退款金额" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="userPhone" label="用户电话">
                <Input placeholder="请输入用户电话" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="merchPhone" label="商家电话">
                <Input placeholder="请输入商家电话" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="remark" label="备注">
                <Input.TextArea placeholder="请输入备注" />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>

      <Modal
        title={`批量编辑 (${selectedRowKeys.length} 条)`}
        open={isBatchModalVisible}
        onOk={handleBatchSubmit}
        onCancel={() => setIsBatchModalVisible(false)}
        okText="确定"
        cancelText="取消"
        className="form-modal"
      >
        <Form form={batchForm} layout="vertical">
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="status" label="状态">
                <Select placeholder="请选择状态" allowClear>
                  <Option value="未完成">未完成</Option>
                  <Option value="进行中">进行中</Option>
                  <Option value="已完成">已完成</Option>
                  <Option value="申请退款">申请退款</Option>
                  <Option value="已退款">已退款</Option>
                  <Option value="拒绝退款">拒绝退款</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="payType" label="支付方式">
                <Select placeholder="请选择支付方式" allowClear>
                  <Option value="alipay">支付宝</Option>
                  <Option value="wechat">微信</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>

      <Modal title="订单详情" open={detailVisible} onCancel={() => setDetailVisible(false)} okText="确定" cancelText="取消" footer={null} width={600} className="detail-modal">
        <Spin spinning={detailLoading}>
          {currentOrder && (
          <div>
            <div className="detail-header">
              <span className="detail-title">{currentOrder.orderNo}</span>
              <Tag className="detail-tag" color={statusColors[currentOrder.status] || 'default'}>
                {currentOrder.status}
              </Tag>
            </div>
            <div className="detail-content">
              <div className="detail-grid">
              <div className="detail-item">
                <div className="detail-item-label">订单编号</div>
                <div className="detail-item-value">{currentOrder.orderNo}</div>
              </div>
              <div className="detail-item">
                <div className="detail-item-label">金额</div>
                <div className="detail-item-value">¥{currentOrder.amount}</div>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">ID</div>
              <div className="detail-value">{currentOrder.id}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">订单码</div>
              <div className="detail-value">{currentOrder.code}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">支付码</div>
              <div className="detail-value">{currentOrder.payCode || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">类型</div>
              <div className="detail-value">{currentOrder.type || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">模式</div>
              <div className="detail-value">{currentOrder.mode || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">状态</div>
              <div className="detail-value">
                <Tag color={statusColors[currentOrder.status] || 'default'}>
                  {currentOrder.status}
                </Tag>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">时长</div>
              <div className="detail-value">{currentOrder.duration} 分钟</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">价格</div>
              <div className="detail-value">¥{currentOrder.price}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">押金</div>
              <div className="detail-value">¥{currentOrder.deposit}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">支付金额</div>
              <div className="detail-value">¥{currentOrder.payPrice}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">支付方式</div>
              <div className="detail-value">{currentOrder.payType === 'alipay' ? '支付宝' : '微信'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">用户电话</div>
              <div className="detail-value">{currentOrder.userPhone || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">开始时间</div>
              <div className="detail-value">{currentOrder.startTime || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">结束时间</div>
              <div className="detail-value">{currentOrder.endTime || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">创建时间</div>
              <div className="detail-value">{currentOrder.createdAt}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">更新时间</div>
              <div className="detail-value">{currentOrder.updatedAt}</div>
            </div>
            </div>
          </div>
        )}
        </Spin>
      </Modal>
    </div>
  );
};