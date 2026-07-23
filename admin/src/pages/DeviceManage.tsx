import { Table, Button, Modal, Form, Input, message, Row, Col, Select, Tag, Space, Spin } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined, EyeOutlined, UploadOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import { getDeviceList, getDeviceDetail, createDevice, updateDevice, deleteDevice, batchDeleteDevice, batchUpdateDevice } from '@/api';
import { CustomPagination } from '@/components/CustomPagination';
import { ExportButton } from '@/components/ExportButton';
import type { Device } from '@/types';

const statusColors: Record<string, string> = {
  '1': 'green',
  '0': 'red',
};

const { Option } = Select;

export const DeviceManage = () => {
  const [data, setData] = useState<Device[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isBatchModalVisible, setIsBatchModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [batchForm] = Form.useForm();
  const [editId, setEditId] = useState<number | null>(null);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [showSearch, setShowSearch] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [detailVisible, setDetailVisible] = useState(false);
  const [currentDevice, setCurrentDevice] = useState<Device | null>(null);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '设备名称', dataIndex: 'name', key: 'name' },
    { title: '设备编码', dataIndex: 'code', key: 'code' },
    { title: '板号', dataIndex: 'boardNo', key: 'boardNo' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={statusColors[status] || 'default'}>{status === '1' ? '在线' : '离线'}</Tag>
      ),
    },
    { title: '类型', dataIndex: 'type', key: 'type' },
    { title: '开锁次数', dataIndex: 'lockCount', key: 'lockCount' },
    { title: '充值', dataIndex: 'recharge', key: 'recharge' },
    { title: '版本', dataIndex: 'version', key: 'version' },
    { title: '信号', dataIndex: 'signal', key: 'signal' },
    { title: '温度', dataIndex: 'heat', key: 'heat' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (time: string) => formatTime(time) },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Device) => (
        <Space>
          <Button className="action-btn-detail" icon={<EyeOutlined />} onClick={() => viewDetail(record)} size="small">详情</Button>
          <Button className="action-btn-edit" icon={<EditOutlined />} onClick={() => handleEdit(record)} size="small">编辑</Button>
          <Button className="action-btn-delete" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)} size="small">删除</Button>
        </Space>
      ),
    },
  ];

  const fetchData = async (params: { page?: number; page_size?: number; name?: string; status?: string; type?: string } = {}) => {
    setLoading(true);
    try {
      const res = await getDeviceList({ page: params.page || currentPage, page_size: params.page_size || pageSize, name: params.name, status: params.status, type: params.type });
      setData(res.data.data);
      setTotal(res.data.total);
    } catch (error) {
      message.error('获取设备列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const viewDetail = async (device: Device) => {
    setDetailLoading(true);
    try {
      const res = await getDeviceDetail(device.id);
      setCurrentDevice(res.data);
      setDetailVisible(true);
    } catch (error) {
      console.error(error);
      message.error('获取详情失败');
    } finally {
      setDetailLoading(false);
    }
  };

  const handleAdd = () => {
    setEditId(null);
    form.resetFields();
    setIsModalVisible(true);
  };

  const handleEdit = (record: Device) => {
    setEditId(record.id);
    form.setFieldsValue(record);
    setIsModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个设备吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await deleteDevice(id);
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
      message.warning('请选择要删除的设备');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 个设备吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteDevice({ ids: selectedRowKeys.map(k => k.toString()) });
          message.success('批量删除成功');
          setSelectedRowKeys([]);
          fetchData();
        } catch (error) {
          message.error('批量删除失败');
        }
      },
    });
  };

  const handleBatchEdit = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要编辑的设备');
      return;
    }
    batchForm.resetFields();
    setIsBatchModalVisible(true);
  };

  const handleBatchSubmit = async () => {
    try {
      const values = await batchForm.validateFields();
      await batchUpdateDevice({ ids: selectedRowKeys.map(k => k.toString()), data: values });
      message.success('批量更新成功');
      setIsBatchModalVisible(false);
      setSelectedRowKeys([]);
      fetchData();
    } catch (error) {
      message.error('批量更新失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editId) {
        await updateDevice(editId, values);
        message.success('更新成功');
      } else {
        await createDevice(values);
        message.success('创建成功');
      }
      setIsModalVisible(false);
      fetchData();
    } catch (error) {
      message.error('提交失败');
    }
  };

  const exportHeaders: Record<string, string> = {
    id: 'ID',
    name: '设备名称',
    code: '设备编码',
    boardNo: '板号',
    status: '状态',
    type: '类型',
    lockCount: '开锁次数',
    recharge: '充值',
    version: '版本',
    signal: '信号',
    heat: '温度',
    createdAt: '创建时间',
    updatedAt: '更新时间',
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
        const response = await fetch('/admin/device/import', {
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

  const handleTableChange = (pagination: { current: number; pageSize: number }) => {
    setCurrentPage(pagination.current);
    setPageSize(pagination.pageSize);
    fetchData({ page: pagination.current, page_size: pagination.pageSize });
  };

  const rowSelection = {
    selectedRowKeys,
    onChange: (keys: React.Key[]) => setSelectedRowKeys(keys),
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

    // 格式化时间
  const formatTime = (time: string) => {
    if (!time) return '-'
    return time.replace('T', ' ').substring(0, 19)
  };

  return (
    <div className="page-container">
      <div className="page-header">
        <h2>设备管理</h2>
        <div className="action-buttons">
          <Button className="action-btn-search" size="small" icon={<SearchOutlined />} onClick={() => setShowSearch(!showSearch)}>{showSearch ? '收起搜索' : '搜索'}</Button>
          <ExportButton data={data} filename="设备列表" headers={exportHeaders} />
          <Button className="action-btn-import" size="small" icon={<UploadOutlined />} onClick={handleImport}>导入</Button>
          <Button className="action-btn-edit" size="small" icon={<EditOutlined />} onClick={handleBatchEdit} disabled={selectedRowKeys.length === 0}>编辑({selectedRowKeys.length})</Button>
          <Button className="action-btn-delete" size="small" icon={<DeleteOutlined />} onClick={handleBatchDelete} disabled={selectedRowKeys.length === 0}>删除({selectedRowKeys.length})</Button>
          <Button className="action-btn-add" size="small" icon={<PlusOutlined />} onClick={handleAdd}>添加</Button>
        </div>
      </div>

      {showSearch && (
        <Form className="search-form" form={searchForm} layout="inline">
          <Form.Item name="name">
            <Input placeholder="设备名称" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="code">
            <Input placeholder="设备编码" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="boardNo">
            <Input placeholder="板号" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="status">
            <Select placeholder="状态" allowClear>
              <Option value="1">在线</Option>
              <Option value="0">离线</Option>
            </Select>
          </Form.Item>
          <Form.Item name="type">
            <Select placeholder="设备类型" allowClear>
              <Option value="lock">门锁</Option>
              <Option value="device">设备</Option>
            </Select>
          </Form.Item>
          <Form.Item name="version">
            <Input placeholder="版本号" prefix={<SearchOutlined />} />
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
          scroll={{ x: 1200 }}
        />
        <CustomPagination
          total={total}
          current={currentPage}
          pageSize={pageSize}
          onChange={(page, pageSize) => handleTableChange({ current: page, pageSize })}
        />
      </div>

      <Modal
        title={editId ? '编辑设备' : '新增设备'}
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
              <Form.Item name="name" label="设备名称" rules={[{ required: true, message: '请输入设备名称' }, { max: 100, message: '名称长度不超过100' }]}>
                <Input placeholder="请输入设备名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="code" label="设备编码" rules={[{ required: true, message: '请输入设备编码' }, { max: 50, message: '编码长度不超过50' }]}>
                <Input placeholder="请输入设备编码" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="boardNo" label="板号" rules={[{ max: 50, message: '板号长度不超过50' }]}>
                <Input placeholder="请输入板号" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="status" label="状态">
                <Select defaultValue="1">
                  <Option value="1">在线</Option>
                  <Option value="0">离线</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="type" label="类型">
                <Input placeholder="请输入类型" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="lockCount" label="开锁次数">
                <Input type="number" placeholder="请输入开锁次数" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="recharge" label="充值">
                <Input placeholder="请输入充值信息" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="version" label="版本">
                <Input placeholder="请输入版本号" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="signal" label="信号">
                <Input placeholder="请输入信号强度" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="heat" label="温度">
                <Input placeholder="请输入温度" />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="protectHeat" label="保护温度">
            <Input placeholder="请输入保护温度" />
          </Form.Item>
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
          <Form.Item name="status" label="状态">
            <Select placeholder="请选择状态" allowClear>
              <Option value="1">在线</Option>
              <Option value="0">离线</Option>
            </Select>
          </Form.Item>
          <Form.Item name="type" label="类型">
            <Input placeholder="请输入类型" />
          </Form.Item>
          <Form.Item name="version" label="版本">
            <Input placeholder="请输入版本号" />
          </Form.Item>
          <Form.Item name="signal" label="信号">
            <Input placeholder="请输入信号强度" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal title="设备详情" open={detailVisible} onCancel={() => setDetailVisible(false)}
        okText="确定" cancelText="取消" footer={null} className="detail-modal">
        <Spin spinning={detailLoading}>
          {currentDevice && (
          <div>
            <div className="detail-header">
              <span className="detail-title">{currentDevice.name}</span>
              <Tag className="detail-tag" color={currentDevice.status === '1' ? 'green' : 'red'}>
                {currentDevice.status === '1' ? '在线' : '离线'}
              </Tag>
            </div>
            <div className="detail-content">
              <div className="detail-grid">
              <div className="detail-item">
                <div className="detail-item-label">设备名称</div>
                <div className="detail-item-value">{currentDevice.name}</div>
              </div>
              <div className="detail-item">
                <div className="detail-item-label">设备编码</div>
                <div className="detail-item-value">{currentDevice.code}</div>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">ID</div>
              <div className="detail-value">{currentDevice.id}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">板号</div>
              <div className="detail-value">{currentDevice.boardNo}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">类型</div>
              <div className="detail-value">{currentDevice.type}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">状态</div>
              <div className="detail-value">
                <Tag color={currentDevice.status === '1' ? 'green' : 'red'}>
                  {currentDevice.status === '1' ? '在线' : '离线'}
                </Tag>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">开锁次数</div>
              <div className="detail-value">{currentDevice.lockCount}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">充值</div>
              <div className="detail-value">{currentDevice.recharge || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">版本</div>
              <div className="detail-value">{currentDevice.version || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">信号</div>
              <div className="detail-value">{currentDevice.signal || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">温度</div>
              <div className="detail-value">{currentDevice.heat || '-'}°C</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">保护温度</div>
              <div className="detail-value">{currentDevice.protectHeat || '-'}°C</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">创建时间</div>
              <div className="detail-value">{currentDevice.createdAt}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">更新时间</div>
              <div className="detail-value">{currentDevice.updatedAt}</div>
            </div>
            </div>
          </div>
        )}
        </Spin>
      </Modal>
    </div>
  );
};