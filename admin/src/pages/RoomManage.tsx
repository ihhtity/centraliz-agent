import { Table, Button, Modal, Form, Input, InputNumber, message, Row, Col, Select, Tag, Space, Spin } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined, EyeOutlined, UploadOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import { getRoomList, getRoomDetail, createRoom, updateRoom, deleteRoom, batchDeleteRoom, batchUpdateRoom } from '@/api';
import { CustomPagination } from '@/components/CustomPagination';
import { ExportButton } from '@/components/ExportButton';
import type { Room } from '@/types';

const statusColors: Record<string, string> = {
  '1': 'green',
  '0': 'red',
};

export const RoomManage = () => {
  const [data, setData] = useState<Room[]>([]);
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
  const [currentRoom, setCurrentRoom] = useState<Room | null>(null);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '描述', dataIndex: 'tag', key: 'tag' },
    { title: '基础价格', dataIndex: 'price', key: 'price', render: (val: number) => <span>¥{val}/小时</span> },
    { title: '房号', dataIndex: 'boardNo', key: 'boardNo' },
    { title: '锁号', dataIndex: 'lockNo', key: 'lockNo' },
    { title: '免费时长', dataIndex: 'freeTime', key: 'freeTime' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={statusColors[status] || 'default'}>{status === '1' ? '启用' : '禁用'}</Tag>
      ),
    },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (time: string) => formatTime(time) },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Room) => (
        <Space>
          <Button className="action-btn-detail" icon={<EyeOutlined />} onClick={() => viewDetail(record)} size="small">详情</Button>
          <Button className="action-btn-edit" icon={<EditOutlined />} onClick={() => handleEdit(record)} size="small">编辑</Button>
          <Button className="action-btn-delete" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)} size="small">删除</Button>
        </Space>
      ),
    },
  ];

  const fetchData = async (params: { page?: number; page_size?: number; name?: string; status?: string; boardNo?: string } = {}) => {
    setLoading(true);
    try {
      const res = await getRoomList({ page: params.page || currentPage, page_size: params.page_size || pageSize, name: params.name, status: params.status, board_no: params.boardNo });
      setData(res.data.data);
      setTotal(res.data.total);
    } catch (error) {
      message.error('获取房间列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const viewDetail = async (room: Room) => {
    setDetailLoading(true);
    try {
      const res = await getRoomDetail(room.id);
      setCurrentRoom(res.data);
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

  const handleEdit = (record: Room) => {
    setEditId(record.id);
    form.setFieldsValue(record);
    setIsModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个房间吗？',
      okText: '确定',
      cancelText: '取消',
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
      message.warning('请选择要删除的房间');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 个房间吗？`,
      okText: '确定',
      cancelText: '取消',
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

  const handleBatchEdit = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要编辑的房间');
      return;
    }
    batchForm.resetFields();
    setIsBatchModalVisible(true);
  };

  const handleBatchSubmit = async () => {
    try {
      const values = await batchForm.validateFields();
      await batchUpdateRoom({ ids: selectedRowKeys.map(k => k.toString()), data: values });
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

  const exportHeaders: Record<string, string> = {
    id: 'ID',
    name: '名称',
    tag: '描述',
    price: '基础价格',
    boardNo: '房号',
    lockNo: '锁号',
    freeTime: '免费时长',
    status: '状态',
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
        const response = await fetch('/admin/room/import', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
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
        <h2>房间管理</h2>
        <div className="action-buttons">
          <Button className="action-btn-search" size="small" icon={<SearchOutlined />} onClick={() => setShowSearch(!showSearch)}>{showSearch ? '收起搜索' : '搜索'}</Button>
          <ExportButton data={data} filename="房间列表" headers={exportHeaders} />
          <Button className="action-btn-import" size="small" icon={<UploadOutlined />} onClick={handleImport}>导入</Button>
          <Button className="action-btn-edit" size="small" icon={<EditOutlined />} onClick={handleBatchEdit} disabled={selectedRowKeys.length === 0}>编辑({selectedRowKeys.length})</Button>
          <Button className="action-btn-delete" size="small" icon={<DeleteOutlined />} onClick={handleBatchDelete} disabled={selectedRowKeys.length === 0}>删除({selectedRowKeys.length})</Button>
          <Button className="action-btn-add" size="small" icon={<PlusOutlined />} onClick={handleAdd}>添加</Button>
        </div>
      </div>

      {showSearch && (
        <Form className="search-form" form={searchForm} layout="inline">
          <Form.Item name="name">
            <Input placeholder="请输入房间名称" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="boardNo">
            <Input placeholder="房号" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="status">
            <Select placeholder="请选择状态" allowClear>
              <Select.Option value="1">启用</Select.Option>
              <Select.Option value="0">禁用</Select.Option>
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
          scroll={{ x: 800 }}
        />
        <CustomPagination
          total={total}
          current={currentPage}
          pageSize={pageSize}
          onChange={(page, pageSize) => handleTableChange({ current: page, pageSize })}
        />
      </div>

      {/* 添加/编辑房间弹窗 */}
      <Modal
        title={editId ? '编辑房间' : '添加房间'}
        open={isModalVisible}
        onOk={handleSubmit}
        onCancel={() => setIsModalVisible(false)}
        okText="确定"
        cancelText="取消"
        className="form-modal"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入名称' }, { max: 100, message: '名称长度不超过100' }]}>
            <Input placeholder="请输入房间名称" />
          </Form.Item>
          <Form.Item name="tag" label="描述" rules={[{ max: 500, message: '描述长度不超过500' }]}>
            <Input.TextArea placeholder="请输入描述信息" />
          </Form.Item>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="price" label="基础价格" rules={[{ required: true, message: '请输入价格' }, { min: 0, message: '价格不能为负数' }]}>
                <InputNumber placeholder="请输入价格" style={{ width: '100%' }} prefix="¥" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="roomsId" label="最大人数" rules={[{ min: 1, message: '最大人数至少为1' }]}>
                <InputNumber placeholder="请输入最大人数" style={{ width: '100%' }} />
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
          <Form.Item name="tag" label="描述">
            <Input.TextArea placeholder="请输入描述信息" />
          </Form.Item>
          <Form.Item name="price" label="基础价格">
            <InputNumber placeholder="请输入价格" style={{ width: '100%' }} prefix="¥" />
          </Form.Item>
          <Form.Item name="freeTime" label="免费时长">
            <InputNumber placeholder="请输入免费时长" style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select placeholder="请选择状态">
              <Select.Option value="1">启用</Select.Option>
              <Select.Option value="0">禁用</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      <Modal title="房间详情" open={detailVisible} onCancel={() => setDetailVisible(false)} okText="确定" cancelText="取消" footer={null} className="detail-modal">
        <Spin spinning={detailLoading}>
          {currentRoom && (
          <div>
            <div className="detail-header">
              <span className="detail-title">{currentRoom.name}</span>
              <Tag className="detail-tag" color={currentRoom.status === '1' ? 'green' : 'red'}>
                {currentRoom.status === '1' ? '启用' : '禁用'}
              </Tag>
            </div>
            <div className="detail-content">
              <div className="detail-grid">
              <div className="detail-item">
                <div className="detail-item-label">房间名称</div>
                <div className="detail-item-value">{currentRoom.name}</div>
              </div>
              <div className="detail-item">
                <div className="detail-item-label">基础价格</div>
                <div className="detail-item-value">¥{currentRoom.price}/小时</div>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">ID</div>
              <div className="detail-value">{currentRoom.id}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">描述</div>
              <div className="detail-value">{currentRoom.tag || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">房号</div>
              <div className="detail-value">{currentRoom.boardNo}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">锁号</div>
              <div className="detail-value">{currentRoom.lockNo}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">免费时长</div>
              <div className="detail-value">{currentRoom.freeTime} 分钟</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">状态</div>
              <div className="detail-value">
                <Tag color={currentRoom.status === '1' ? 'green' : 'red'}>
                  {currentRoom.status === '1' ? '启用' : '禁用'}
                </Tag>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">创建时间</div>
              <div className="detail-value">{currentRoom.createdAt}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">更新时间</div>
              <div className="detail-value">{currentRoom.updatedAt}</div>
            </div>
            </div>
          </div>
        )}
        </Spin>
      </Modal>
    </div>
  );
};