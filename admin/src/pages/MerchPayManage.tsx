import { Table, Button, Modal, Form, Input, Select, message, Tag, Space, Spin, Row, Col, InputNumber } from 'antd';
import { SearchOutlined, DeleteOutlined, EyeOutlined, PlusOutlined, EditOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import type { MerchPay } from '@/types';
import { getMerchPayList, getMerchPayDetail, batchDeleteMerchPay, importMerchPay, createMerchPay, batchUpdateMerchPay, updateMerchPay } from '@/api';
import { CustomPagination } from '@/components/CustomPagination';
import { ExportButton } from '@/components/ExportButton';

const statusColors: Record<string, string> = {
  '未完成': 'orange',
  '已完成': 'green',
  '已关闭': 'red',
};

export const MerchPayManage = () => {
  const [data, setData] = useState<MerchPay[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [showSearch, setShowSearch] = useState(false);
  const [currentItem, setCurrentItem] = useState<MerchPay | null>(null);
  const [detailVisible, setDetailVisible] = useState(false);
  const [searchForm] = Form.useForm();
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isBatchModalVisible, setIsBatchModalVisible] = useState(false);
  const [editId, setEditId] = useState<number | null>(null);
  const [form] = Form.useForm();
  const [batchForm] = Form.useForm();

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '订单号', dataIndex: 'code', key: 'code' },
    { title: '商品名称', dataIndex: 'name', key: 'name' },
    { title: '汇付订单号', dataIndex: 'hfSeqId', key: 'hfSeqId' },
    { title: '汇付支付时间', dataIndex: 'reqDate', key: 'reqDate' },
    { title: '订单原价', dataIndex: 'originalPrice', key: 'originalPrice' },
    { title: '实际支付', dataIndex: 'price', key: 'price' },
    { title: '锁总数', dataIndex: 'locktotal', key: 'locktotal' },
    { title: '订单类型', dataIndex: 'type', key: 'type' },
    {
      title: '订单状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={statusColors[status] || 'default'}>{status}</Tag>
      ),
    },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (time: string) => formatTime(time) },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: MerchPay) => (
        <Space>
          <Button className="action-btn-detail" icon={<EyeOutlined />} onClick={() => viewDetail(record)} size="small">详情</Button>
          <Button className="action-btn-edit" icon={<EditOutlined />} onClick={() => handleEdit(record)} size="small">编辑</Button>
          <Button className="action-btn-delete" icon={<DeleteOutlined />} onClick={() => deleteItem(record.id)} size="small">删除</Button>
        </Space>
      ),
    },
  ];

  const loadData = async (params: Record<string, any> = {}) => {
    setLoading(true);
    try {
      const res = await getMerchPayList({
        ...params,
        page: params.page || currentPage,
        page_size: params.page_size || pageSize,
      });
      setData(res.data.data);
      setTotal(res.data.total);
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleSearch = () => {
    const values = searchForm.getFieldsValue();
    setCurrentPage(1);
    loadData(values);
  };

  const handleReset = () => {
    searchForm.resetFields();
    setCurrentPage(1);
    loadData();
  };

  const viewDetail = async (item: MerchPay) => {
    setDetailLoading(true);
    try {
      const res = await getMerchPayDetail(item.id);
      setCurrentItem(res.data);
      setDetailVisible(true);
    } catch (error) {
      console.error(error);
      message.error('获取详情失败');
    } finally {
      setDetailLoading(false);
    }
  };

  const deleteItem = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个商户支付记录吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteMerchPay({ ids: [id.toString()] });
          message.success('删除成功');
          loadData();
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  const handleBatchDelete = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要删除的商户支付记录');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 条记录吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteMerchPay({ ids: selectedRowKeys.map(String) });
          message.success('批量删除成功');
          setSelectedRowKeys([]);
          loadData();
        } catch (error) {
          message.error('批量删除失败');
        }
      },
    });
  };

  const handleImport = async (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    try {
      await importMerchPay(formData);
      message.success('导入成功');
      loadData();
    } catch (error) {
      message.error('导入失败');
    }
  };

  const handleAdd = () => {
    setEditId(null);
    form.resetFields();
    setIsModalVisible(true);
  };

  const handleEdit = (record: MerchPay) => {
    setEditId(record.id);
    form.setFieldsValue(record);
    setIsModalVisible(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editId) {
        await updateMerchPay(editId, values);
        message.success('更新成功');
      } else {
        await createMerchPay(values);
        message.success('创建成功');
      }
      setIsModalVisible(false);
      loadData();
    } catch (error) {
      message.error('提交失败');
    }
  };

  const handleBatchEdit = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要编辑的记录');
      return;
    }
    batchForm.resetFields();
    setIsBatchModalVisible(true);
  };

  const handleBatchSubmit = async () => {
    try {
      const values = await batchForm.validateFields();
      const updateData = {
        ids: selectedRowKeys.map(String),
        ...values,
      };
      await batchUpdateMerchPay(updateData);
      message.success('批量更新成功');
      setIsBatchModalVisible(false);
      setSelectedRowKeys([]);
      loadData();
    } catch (error) {
      message.error('批量更新失败');
    }
  };

  const exportHeaders: Record<string, string> = {
    id: 'ID',
    code: '订单号',
    name: '商品名称',
    hfSeqId: '汇付订单号',
    reqDate: '汇付支付时间',
    originalPrice: '订单原价',
    price: '实际支付',
    locktotal: '锁总数',
    type: '订单类型',
    status: '订单状态',
    remarks: '订单备注',
    createdAt: '创建时间',
  };

  const handleTableChange = (pagination: { current: number; pageSize: number }) => {
    setCurrentPage(pagination.current);
    setPageSize(pagination.pageSize);
    loadData({ page: pagination.current, page_size: pagination.pageSize });
  };

  // 格式化时间
  const formatTime = (time: string) => {
    if (!time) return '-'
    return time.replace('T', ' ').substring(0, 19)
  };

  return (
    <div className="page-container">
      <div className="page-header">
        <h2>商户支付管理</h2>
        <div className="action-buttons">
          <Button className="action-btn-search" size="small" icon={<SearchOutlined />} onClick={() => setShowSearch(!showSearch)}>{showSearch ? '收起搜索' : '搜索'}</Button>
          <ExportButton data={data} filename="商户支付列表" headers={exportHeaders} />
          <Button className="action-btn-import" size="small" onClick={() => document.getElementById('merch-pay-import')?.click()}>导入</Button>
          <input type="file" id="merch-pay-import" style={{ display: 'none' }} onChange={(e) => e.target.files?.[0] && handleImport(e.target.files[0])} />
          <Button className="action-btn-edit" size="small" icon={<EditOutlined />} onClick={handleBatchEdit} disabled={selectedRowKeys.length === 0}>编辑({selectedRowKeys.length})</Button>
          <Button className="action-btn-delete" size="small" icon={<DeleteOutlined />} onClick={handleBatchDelete} disabled={selectedRowKeys.length === 0}>删除({selectedRowKeys.length})</Button>
          <Button className="action-btn-add" size="small" icon={<PlusOutlined />} onClick={handleAdd}>添加</Button>
        </div>
      </div>

      {showSearch && (
        <Form className="search-form" form={searchForm} layout="inline">
          <Form.Item name="code">
            <Input placeholder="订单号" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="name">
            <Input placeholder="商品名称" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="hfSeqId">
            <Input placeholder="汇付订单号" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="type">
            <Select placeholder="订单类型" allowClear>
              <Select.Option value="购买">购买</Select.Option>
              <Select.Option value="续费">续费</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="status">
            <Select placeholder="订单状态" allowClear>
              <Select.Option value="未完成">未完成</Select.Option>
              <Select.Option value="已完成">已完成</Select.Option>
              <Select.Option value="已关闭">已关闭</Select.Option>
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
          columns={columns}
          dataSource={data}
          rowKey="id"
          loading={loading}
          pagination={false}
          rowSelection={{
            selectedRowKeys,
            onChange: setSelectedRowKeys,
            onSelect: (record: MerchPay, selected: boolean) => {
              if (selected) {
                setSelectedRowKeys([...selectedRowKeys, record.id]);
              } else {
                setSelectedRowKeys(selectedRowKeys.filter(k => k !== record.id));
              }
            },
          }}
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

      <Modal title="商户支付详情" open={detailVisible} onCancel={() => setDetailVisible(false)} okText="确定" cancelText="取消" footer={null} className="detail-modal">
        <Spin spinning={detailLoading}>
          {currentItem && (
            <div>
              <div className="detail-header">
                <span className="detail-title">{currentItem.code}</span>
                <Tag className="detail-tag" color={statusColors[currentItem.status] || 'default'}>
                  {currentItem.status}
                </Tag>
              </div>
              <div className="detail-content">
                <div className="detail-grid">
                  <div className="detail-item">
                    <div className="detail-item-label">商品名称</div>
                    <div className="detail-item-value">{currentItem.name}</div>
                  </div>
                  <div className="detail-item">
                    <div className="detail-item-label">订单状态</div>
                    <div className="detail-item-value">
                      <Tag color={statusColors[currentItem.status] || 'default'}>{currentItem.status}</Tag>
                    </div>
                  </div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">ID</div>
                  <div className="detail-value">{currentItem.id}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">汇付订单号</div>
                  <div className="detail-value">{currentItem.hfSeqId || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">汇付支付时间</div>
                  <div className="detail-value">{currentItem.reqDate || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">订单原价</div>
                  <div className="detail-value">{currentItem.originalPrice}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">实际支付金额</div>
                  <div className="detail-value">{currentItem.price}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">锁总数</div>
                  <div className="detail-value">{currentItem.locktotal}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">订单类型</div>
                  <div className="detail-value">{currentItem.type}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">订单备注</div>
                  <div className="detail-value">{currentItem.remarks || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">创建时间</div>
                  <div className="detail-value">{currentItem.createdAt}</div>
                </div>
              </div>
            </div>
          )}
        </Spin>
      </Modal>

      <Modal
        title={editId ? '编辑商户支付' : '添加商户支付'}
        open={isModalVisible}
        onOk={handleSubmit}
        onCancel={() => setIsModalVisible(false)}
        okText="确定"
        cancelText="取消"
        className="form-modal"
      >
        <Form form={form} layout="vertical">
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="code" label="订单号" rules={[{ required: true, message: '请输入订单号' }]}>
                <Input placeholder="请输入订单号" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="name" label="商品名称" rules={[{ required: true, message: '请输入商品名称' }]}>
                <Input placeholder="请输入商品名称" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="hfSeqId" label="汇付订单号">
                <Input placeholder="请输入汇付订单号" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="reqDate" label="汇付支付时间">
                <Input placeholder="请输入汇付支付时间" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="originalPrice" label="订单原价">
                <InputNumber placeholder="请输入订单原价" style={{ width: '100%' }} prefix="¥" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="price" label="实际支付金额">
                <InputNumber placeholder="请输入实际支付金额" style={{ width: '100%' }} prefix="¥" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="locktotal" label="锁总数">
                <InputNumber placeholder="请输入锁总数" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="type" label="订单类型">
                <Select placeholder="请选择订单类型">
                  <Select.Option value="购买">购买</Select.Option>
                  <Select.Option value="续费">续费</Select.Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="status" label="订单状态" initialValue="未完成">
                <Select placeholder="请选择订单状态">
                  <Select.Option value="未完成">未完成</Select.Option>
                  <Select.Option value="已完成">已完成</Select.Option>
                  <Select.Option value="已关闭">已关闭</Select.Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="remarks" label="订单备注">
                <Input placeholder="请输入订单备注" />
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
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="status" label="订单状态">
                <Select placeholder="请选择订单状态" allowClear>
                  <Select.Option value="未完成">未完成</Select.Option>
                  <Select.Option value="已完成">已完成</Select.Option>
                  <Select.Option value="已关闭">已关闭</Select.Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="type" label="订单类型">
                <Select placeholder="请选择订单类型" allowClear>
                  <Select.Option value="购买">购买</Select.Option>
                  <Select.Option value="续费">续费</Select.Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="price" label="实际支付金额">
                <InputNumber placeholder="请输入实际支付金额" style={{ width: '100%' }} prefix="¥" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="originalPrice" label="订单原价">
                <InputNumber placeholder="请输入订单原价" style={{ width: '100%' }} prefix="¥" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="locktotal" label="锁总数">
                <InputNumber placeholder="请输入锁总数" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="merchsId" label="商家ID">
                <InputNumber placeholder="请输入商家ID" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={24}>
              <Form.Item name="remarks" label="订单备注">
                <Input.TextArea placeholder="请输入订单备注" />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>
    </div>
  );
};