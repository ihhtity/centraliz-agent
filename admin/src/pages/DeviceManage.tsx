import { Table, Button, Modal, Form, Input, Select, message, Row, Col } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import { getDeviceList, createDevice, updateDevice, deleteDevice, batchDeleteDevice } from '@/api';
import type { Device } from '@/types';

const { Option } = Select;

export const DeviceManage = () => {
  const [data, setData] = useState<Device[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [editId, setEditId] = useState<number | null>(null);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '设备名称', dataIndex: 'name', key: 'name' },
    { title: '设备编码', dataIndex: 'code', key: 'code' },
    { title: '板号', dataIndex: 'boardNo', key: 'boardNo' },
    { title: '状态', dataIndex: 'status', key: 'status', render: (status: string) => (status === '1' ? '在线' : '离线') },
    { title: '类型', dataIndex: 'type', key: 'type' },
    { title: '开锁次数', dataIndex: 'lockCount', key: 'lockCount' },
    { title: '充值', dataIndex: 'recharge', key: 'recharge' },
    { title: '版本', dataIndex: 'version', key: 'version' },
    { title: '信号', dataIndex: 'signal', key: 'signal' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Device) => (
        <div className="action-buttons">
          <Button type="primary" ghost size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
          <Button danger size="small" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>删除</Button>
        </div>
      ),
    },
  ];

  const fetchData = async (params: { page?: number; page_size?: number; name?: string; status?: string; type?: string } = {}) => {
    setLoading(true);
    try {
      const res = await getDeviceList({ page: params.page || 1, page_size: params.page_size || 10, name: params.name, status: params.status, type: params.type });
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

  const rowSelection = {
    selectedRowKeys,
    onChange: (keys: React.Key[]) => setSelectedRowKeys(keys),
  };

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
        <h2>设备管理</h2>
        <div className="action-buttons">
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>新增设备</Button>
          <Button danger icon={<DeleteOutlined />} onClick={handleBatchDelete}>批量删除</Button>
        </div>
      </div>

      <Form className="search-form" form={searchForm} layout="inline">
        <Form.Item name="name">
          <Input placeholder="设备名称" prefix={<SearchOutlined />} />
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
        <Form.Item>
          <Button type="primary" onClick={handleSearch}>搜索</Button>
          <Button onClick={handleReset} style={{ marginLeft: 8 }}>重置</Button>
        </Form.Item>
      </Form>

      <div className="table-container">
        <Table
          rowSelection={rowSelection}
          columns={columns}
          dataSource={data}
          pagination={{ total, pageSize: 10, showSizeChanger: true, showQuickJumper: true }}
          loading={loading}
          rowKey="id"
          scroll={{ x: 1200 }}
        />
      </div>

      <Modal
        title={editId ? '编辑设备' : '新增设备'}
        visible={isModalVisible}
        onOk={handleSubmit}
        onCancel={() => setIsModalVisible(false)}
        className="form-modal"
      >
        <Form form={form} layout="vertical">
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="name" label="设备名称" rules={[{ required: true, message: '请输入设备名称' }]}>
                <Input placeholder="请输入设备名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="code" label="设备编码">
                <Input placeholder="请输入设备编码" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="boardNo" label="板号">
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
    </div>
  );
};