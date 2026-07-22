import { Table, Button, Modal, Form, Input, Select, message, Row, Col } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import { getGroupList, createGroup, updateGroup, deleteGroup, batchDeleteGroup } from '@/api';
import type { Group } from '@/types';

const { Option } = Select;

export const GroupManage = () => {
  const [data, setData] = useState<Group[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [editId, setEditId] = useState<number | null>(null);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '分组名称', dataIndex: 'name', key: 'name' },
    { title: '电话', dataIndex: 'phone', key: 'phone' },
    { title: '房间数量', dataIndex: 'roomCount', key: 'roomCount' },
    { title: '类型', dataIndex: 'type', key: 'type' },
    { title: '位置', dataIndex: 'location', key: 'location' },
    { title: '地址', dataIndex: 'address', key: 'address' },
    { title: '规则名称', dataIndex: 'ruleName', key: 'ruleName' },
    { title: '绑定编号', dataIndex: 'bindNumber', key: 'bindNumber' },
    { title: '消费推送', dataIndex: 'consumePush', key: 'consumePush' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Group) => (
        <div className="action-buttons">
          <Button type="primary" ghost size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
          <Button danger size="small" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>删除</Button>
        </div>
      ),
    },
  ];

  const fetchData = async (params: { page?: number; page_size?: number; name?: string; type?: string } = {}) => {
    setLoading(true);
    try {
      const res = await getGroupList({ page: params.page || 1, page_size: params.page_size || 10, name: params.name, type: params.type });
      setData(res.data.data);
      setTotal(res.data.total);
    } catch (error) {
      message.error('获取分组列表失败');
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

  const handleEdit = (record: Group) => {
    setEditId(record.id);
    form.setFieldsValue(record);
    setIsModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个分组吗？',
      onOk: async () => {
        try {
          await deleteGroup(id);
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
      message.warning('请选择要删除的分组');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 个分组吗？`,
      onOk: async () => {
        try {
          await batchDeleteGroup({ ids: selectedRowKeys.map(k => k.toString()) });
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
        await updateGroup(editId, values);
        message.success('更新成功');
      } else {
        await createGroup(values);
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
        <h2>分组管理</h2>
        <div className="action-buttons">
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>新增分组</Button>
          <Button danger icon={<DeleteOutlined />} onClick={handleBatchDelete}>批量删除</Button>
        </div>
      </div>

      <Form className="search-form" form={searchForm} layout="inline">
        <Form.Item name="name">
          <Input placeholder="分组名称" prefix={<SearchOutlined />} />
        </Form.Item>
        <Form.Item name="type">
          <Select placeholder="分组类型" allowClear>
            <Option value="hotel">酒店</Option>
            <Option value="office">办公室</Option>
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
          scroll={{ x: 1300 }}
        />
      </div>

      <Modal
        title={editId ? '编辑分组' : '新增分组'}
        visible={isModalVisible}
        onOk={handleSubmit}
        onCancel={() => setIsModalVisible(false)}
        className="form-modal"
      >
        <Form form={form} layout="vertical">
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="name" label="分组名称" rules={[{ required: true, message: '请输入分组名称' }]}>
                <Input placeholder="请输入分组名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="phone" label="电话">
                <Input placeholder="请输入电话" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="roomCount" label="房间数量">
                <Input type="number" placeholder="请输入房间数量" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="type" label="类型">
                <Input placeholder="请输入类型" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="location" label="位置">
                <Input placeholder="请输入位置" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="address" label="地址">
                <Input placeholder="请输入地址" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="ruleName" label="规则名称">
                <Input placeholder="请输入规则名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="bindNumber" label="绑定编号">
                <Input placeholder="请输入绑定编号" />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="consumePush" label="消费推送">
            <Select defaultValue="1">
              <Option value="1">开启</Option>
              <Option value="0">关闭</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};