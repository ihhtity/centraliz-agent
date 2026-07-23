import { Table, Button, Modal, Form, Input, Select, message, Tag, Space, Spin } from 'antd';
import { EditOutlined, DeleteOutlined, EyeOutlined, SearchOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import type { User } from '@/types';
import {
  getUserList,
  getUserDetail,
  updateUser,
  deleteUser,
  batchDeleteUser,
  batchUpdateUser,
  importUser,
} from '@/api';
import { CustomPagination } from '@/components/CustomPagination';
import { ExportButton } from '@/components/ExportButton';

const statusColors: Record<string, string> = {
  '0': 'green',
  '1': 'red',
};

export const UserManage = () => {
  const [data, setData] = useState<User[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [showSearch, setShowSearch] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [isBatchModalVisible, setIsBatchModalVisible] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const [currentItem, setCurrentItem] = useState<User | null>(null);
  const [detailVisible, setDetailVisible] = useState(false);
  const [form] = Form.useForm();
  const [batchForm] = Form.useForm();
  const [searchForm] = Form.useForm();

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '用户名', dataIndex: 'name', key: 'name' },
    { title: '账号', dataIndex: 'account', key: 'account' },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    { title: '手机号', dataIndex: 'phone', key: 'phone' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string | null) => (
        <Tag color={statusColors[status || ''] || 'default'}>{status === '0' ? '白名单' : status === '1' ? '黑名单' : '-'}</Tag>
      ),
    },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (val: string) => formatTime(val) },
    { title: '更新时间', dataIndex: 'updatedAt', key: 'updatedAt', render: (val: string) => formatTime(val) },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: User) => (
        <Space>
          <Button className="action-btn-detail" icon={<EyeOutlined />} onClick={() => viewDetail(record)} size="small">详情</Button>
          <Button className="action-btn-edit" icon={<EditOutlined />} onClick={() => editItem(record)} size="small">编辑</Button>
          <Button className="action-btn-delete" icon={<DeleteOutlined />} onClick={() => deleteItem(record.id)} size="small">删除</Button>
        </Space>
      ),
    },
  ];

  const loadData = async (params: Record<string, any> = {}) => {
    setLoading(true);
    try {
      const res = await getUserList({
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

  const viewDetail = async (item: User) => {
    setDetailLoading(true);
    try {
      const res = await getUserDetail(item.id);
      setCurrentItem(res.data);
      setDetailVisible(true);
    } catch (error) {
      console.error(error);
      message.error('获取详情失败');
    } finally {
      setDetailLoading(false);
    }
  };

  const editItem = (item: User) => {
    setCurrentItem(item);
    setIsEdit(true);
    form.setFieldsValue(item);
    setModalVisible(true);
  };

  const handleBatchEdit = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要编辑的用户');
      return;
    }
    batchForm.resetFields();
    setIsBatchModalVisible(true);
  };

  const handleBatchSubmit = async () => {
    try {
      const values = await batchForm.validateFields();
      await batchUpdateUser({ ids: selectedRowKeys.map(k => k.toString()), data: values });
      message.success('批量更新成功');
      setIsBatchModalVisible(false);
      setSelectedRowKeys([]);
      loadData();
    } catch (error) {
      message.error('批量更新失败');
    }
  };

  const saveItem = async () => {
    try {
      const values = await form.validateFields();
      if (isEdit && currentItem) {
        await updateUser(currentItem.id, values);
        message.success('更新成功');
      }
      setModalVisible(false);
      loadData();
    } catch (error) {
      console.error(error);
    }
  };

  const deleteItem = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个用户吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await deleteUser(id);
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
      message.warning('请选择要删除的用户');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 个用户吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteUser({ ids: selectedRowKeys.map(String) });
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
      await importUser(formData);
      message.success('导入成功');
      loadData();
    } catch (error) {
      message.error('导入失败');
    }
  };

  const exportHeaders: Record<string, string> = {
    id: 'ID',
    name: '用户名',
    account: '账号',
    email: '邮箱',
    phone: '手机号',
    status: '状态',
    privacy: '隐私政策',
    openid: '微信OpenID',
    unionId: '微信UnionID',
    createdAt: '创建时间',
    updatedAt: '更新时间',
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
        <h2>用户管理</h2>
        <div className="action-buttons">
          <Button className="action-btn-search" size="small" icon={<SearchOutlined />} onClick={() => setShowSearch(!showSearch)}>{showSearch ? '收起搜索' : '搜索'}</Button>
          <ExportButton data={data} filename="用户列表" headers={exportHeaders} />
          <Button className="action-btn-import" size="small" onClick={() => document.getElementById('user-import')?.click()}>导入</Button>
          <input type="file" id="user-import" style={{ display: 'none' }} onChange={(e) => e.target.files?.[0] && handleImport(e.target.files[0])} />
          <Button className="action-btn-edit" size="small" icon={<EditOutlined />} onClick={handleBatchEdit} disabled={selectedRowKeys.length === 0}>编辑({selectedRowKeys.length})</Button>
          <Button className="action-btn-delete" size="small" icon={<DeleteOutlined />} onClick={handleBatchDelete} disabled={selectedRowKeys.length === 0}>删除({selectedRowKeys.length})</Button>
        </div>
      </div>

      {showSearch && (
        <Form className="search-form" form={searchForm} layout="inline">
          <Form.Item name="name">
            <Input placeholder="用户名" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="account">
            <Input placeholder="账号" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="email">
            <Input placeholder="邮箱" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="phone">
            <Input placeholder="手机号" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="openid">
            <Input placeholder="OpenID" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="status">
            <Select placeholder="状态" allowClear>
              <Select.Option value="0">白名单</Select.Option>
              <Select.Option value="1">黑名单</Select.Option>
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
          }}
        />
        <CustomPagination
          total={total}
          current={currentPage}
          pageSize={pageSize}
          onChange={(page, pageSize) => handleTableChange({ current: page, pageSize })}
        />
      </div>

      <Modal
        title="编辑用户"
        open={modalVisible}
        onOk={saveItem}
        onCancel={() => setModalVisible(false)}
        okText="确定"
        cancelText="取消"
        width={600}
        className="form-modal"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="用户名" rules={[{ required: true, message: '请输入用户名' }, { max: 64, message: '用户名长度不超过64' }]}>
            <Input placeholder="请输入用户名" />
          </Form.Item>
          <Form.Item name="email" label="邮箱" rules={[{ type: 'email', message: '请输入正确的邮箱格式' }]}>
            <Input placeholder="请输入邮箱" />
          </Form.Item>
          <Form.Item name="phone" label="手机号" rules={[{ pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式' }]}>
            <Input placeholder="请输入手机号" />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select options={[{ label: '白名单', value: '0' }, { label: '黑名单', value: '1' }]} />
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
        width={600}
        className="form-modal"
      >
        <Form form={batchForm} layout="vertical">
          <Form.Item name="status" label="状态">
            <Select placeholder="请选择状态" allowClear>
              <Select.Option value="0">白名单</Select.Option>
              <Select.Option value="1">黑名单</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      <Modal title="用户详情" open={detailVisible} onCancel={() => setDetailVisible(false)} okText="确定" cancelText="取消" footer={null} className="detail-modal">
        <Spin spinning={detailLoading}>
          {currentItem && (
            <div>
              <div className="detail-header">
                <span className="detail-title">{currentItem.name}</span>
                <Tag className="detail-tag" color={currentItem.status === '0' ? 'green' : currentItem.status === '1' ? 'red' : 'default'}>
                  {currentItem.status === '0' ? '白名单' : currentItem.status === '1' ? '黑名单' : '-'}
                </Tag>
              </div>
              <div className="detail-content">
                <div className="detail-grid">
                  <div className="detail-item">
                    <div className="detail-item-label">用户名</div>
                    <div className="detail-item-value">{currentItem.name}</div>
                  </div>
                  <div className="detail-item">
                    <div className="detail-item-label">账号</div>
                    <div className="detail-item-value">{currentItem.account}</div>
                  </div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">ID</div>
                  <div className="detail-value">{currentItem.id}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">邮箱</div>
                  <div className="detail-value">{currentItem.email || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">手机号</div>
                  <div className="detail-value">{currentItem.phone || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">状态</div>
                  <div className="detail-value">
                    <Tag color={currentItem.status === '0' ? 'green' : currentItem.status === '1' ? 'red' : 'default'}>
                      {currentItem.status === '0' ? '白名单' : currentItem.status === '1' ? '黑名单' : '-'}
                    </Tag>
                  </div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">隐私政策</div>
                  <div className="detail-value">{currentItem.privacy === '0' ? '拒绝' : '同意'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">微信OpenID</div>
                  <div className="detail-value">{currentItem.openid || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">微信UnionID</div>
                  <div className="detail-value">{currentItem.unionId || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">创建时间</div>
                  <div className="detail-value">{currentItem.createdAt || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">更新时间</div>
                  <div className="detail-value">{currentItem.updatedAt || '-'}</div>
                </div>
              </div>
            </div>
          )}
        </Spin>
      </Modal>
    </div>
  );
};