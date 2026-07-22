import { Table, Button, Modal, Form, Input, InputNumber, message, Row, Col } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined, ExportOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import { getRoomList, createRoom, updateRoom, deleteRoom, batchDeleteRoom } from '@/api';
import type { Room } from '@/types';

export const RoomManage = () => {
  const [data, setData] = useState<Room[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [editId, setEditId] = useState<number | null>(null);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '描述', dataIndex: 'tag', key: 'tag' },
    { title: '基础价格', dataIndex: 'price', key: 'price', render: (val: number) => <span>¥{val}/小时</span> },
    { title: '最大人数', dataIndex: 'roomsId', key: 'maxPeople', render: () => '-' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Room) => (
        <div className="action-buttons">
          <Button type="primary" ghost size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
          <Button danger size="small" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>删除</Button>
        </div>
      ),
    },
  ];

  const fetchData = async (params: { page?: number; page_size?: number; name?: string; status?: string } = {}) => {
    setLoading(true);
    try {
      const res = await getRoomList({ page: params.page || 1, page_size: params.page_size || 10, name: params.name, status: params.status });
      setData(res.data.data);
      setTotal(res.data.total);
    } catch (error) {
      message.error('获取包间列表失败');
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

  const handleEdit = (record: Room) => {
    setEditId(record.id);
    form.setFieldsValue(record);
    setIsModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个包间类型吗？',
      onOk: async () => {
        try {
          await deleteRoom(id);
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
      message.warning('请选择要删除的包间类型');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 个包间类型吗？`,
      onOk: async () => {
        try {
          await batchDeleteRoom({ ids: selectedRowKeys.map(k => k.toString()) });
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
        await updateRoom(editId, values);
        message.success('更新成功');
      } else {
        await createRoom(values);
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
        <h2>包间类型管理</h2>
        <div className="action-buttons">
          <Button type="primary" icon={<SearchOutlined />} onClick={() => {}}>搜索</Button>
          <Button icon={<ExportOutlined />} onClick={() => {}}>导出</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>+ 添加类型</Button>
        </div>
      </div>

      <div className="table-container">
        <Table
          rowSelection={rowSelection}
          columns={columns}
          dataSource={data}
          pagination={{ total, pageSize: 10, showSizeChanger: true, showQuickJumper: true }}
          loading={loading}
          rowKey="id"
          scroll={{ x: 800 }}
        />
      </div>

      <Modal
        title={editId ? '编辑包间类型' : '添加包间类型'}
        visible={isModalVisible}
        onOk={handleSubmit}
        onCancel={() => setIsModalVisible(false)}
        className="form-modal"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入名称' }]}>
            <Input placeholder="请输入包间类型名称" />
          </Form.Item>
          <Form.Item name="tag" label="描述">
            <Input.TextArea placeholder="请输入描述信息" />
          </Form.Item>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="price" label="基础价格">
                <InputNumber placeholder="请输入价格" style={{ width: '100%' }} prefix="¥" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="roomsId" label="最大人数">
                <InputNumber placeholder="请输入最大人数" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>
    </div>
  );
};