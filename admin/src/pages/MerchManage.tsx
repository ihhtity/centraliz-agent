import { Table, Button, Modal, Form, Input, Select, message, Tag, Space, Spin, Row, Col } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined, SearchOutlined, ExportOutlined, UploadOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import * as XLSX from 'xlsx';
import type { Merch } from '@/types';
import {
  getMerchList,
  getMerchDetail,
  createMerch,
  updateMerch,
  deleteMerch,
  batchDeleteMerch,
  batchUpdateMerch,
} from '@/api';
import { CustomPagination } from '@/components/CustomPagination';

const roleColors: Record<string, string> = {
  '商家': 'blue',
  '管理者': 'green',
  '代理商': 'purple',
};

const statusColors: Record<string, string> = {
  '0': 'green',
  '1': 'red',
};

export const MerchManage = () => {
  const [data, setData] = useState<Merch[]>([]);
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
  const [currentMerch, setCurrentMerch] = useState<Merch | null>(null);
  const [detailVisible, setDetailVisible] = useState(false);
  const [form] = Form.useForm();
  const [batchForm] = Form.useForm();
  const [searchForm] = Form.useForm();

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '账号', dataIndex: 'account', key: 'account' },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    { title: '手机号', dataIndex: 'phone', key: 'phone' },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (role: string) => (
        <Tag color={roleColors[role] || 'default'}>{role}</Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={statusColors[status] || 'default'}>{status === '0' ? '白名单' : '黑名单'}</Tag>
      ),
    },
    { title: '上次登录', dataIndex: 'logAt', key: 'logAt', render: (time: string) => formatTime(time) },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (time: string) => formatTime(time) },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: Merch) => (
        <Space>
          <Button className="action-btn-detail" icon={<EyeOutlined />} onClick={() => viewDetail(record)} size="small">详情</Button>
          <Button className="action-btn-edit" icon={<EditOutlined />} onClick={() => editMerch(record)} size="small">编辑</Button>
          <Button className="action-btn-delete" icon={<DeleteOutlined />} onClick={() => deleteMerchItem(record.id)} size="small">删除</Button>
        </Space>
      ),
    },
  ];

  const loadData = async (params: Record<string, any> = {}) => {
    setLoading(true);
    try {
      const res = await getMerchList({
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

  const viewDetail = async (merch: Merch) => {
    setDetailLoading(true);
    try {
      const res = await getMerchDetail(merch.id);
      setCurrentMerch(res.data);
      setDetailVisible(true);
    } catch (error) {
      console.error(error);
      message.error('获取详情失败');
    } finally {
      setDetailLoading(false);
    }
  };

  const editMerch = (merch: Merch) => {
    setCurrentMerch(merch);
    setIsEdit(true);
    form.setFieldsValue(merch);
    setModalVisible(true);
  };

  const addMerch = () => {
    setIsEdit(false);
    setCurrentMerch(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleBatchEdit = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要编辑的商家');
      return;
    }
    batchForm.resetFields();
    setIsBatchModalVisible(true);
  };

  const handleBatchSubmit = async () => {
    try {
      const values = await batchForm.validateFields();
      await batchUpdateMerch({ ids: selectedRowKeys.map(k => k.toString()), data: values });
      message.success('批量更新成功');
      setIsBatchModalVisible(false);
      setSelectedRowKeys([]);
      loadData();
    } catch (error) {
      message.error('批量更新失败');
    }
  };

  const saveMerch = async () => {
    try {
      const values = await form.validateFields();
      if (isEdit && currentMerch) {
        await updateMerch(currentMerch.id, values);
        message.success('更新成功');
      } else {
        await createMerch(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      loadData();
    } catch (error) {
      console.error(error);
    }
  };

  const deleteMerchItem = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个商家吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await deleteMerch(id);
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
      message.warning('请选择要删除的商家');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 个商家吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteMerch({ ids: selectedRowKeys.map(String) });
          message.success('批量删除成功');
          setSelectedRowKeys([]);
          loadData();
        } catch (error) {
          message.error('批量删除失败');
        }
      },
    });
  };

  const handleExport = () => {
    if (data.length === 0) {
      message.warning('暂无数据可导出');
      return;
    }
    const exportData = data.map(item => ({
      ID: item.id,
      账号: item.account,
      邮箱: item.email,
      手机号: item.phone,
      角色: item.role,
      状态: item.status === '0' ? '白名单' : '黑名单',
      上次登录: item.logAt,
      创建时间: item.createdAt,
    }));
    const ws = XLSX.utils.json_to_sheet(exportData);
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, '商家');
    XLSX.writeFile(wb, `商家_${new Date().toLocaleDateString()}.xlsx`);
    message.success('导出成功');
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
        const response = await fetch('/admin/merch/import', {
          method: 'POST',
          headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` },
          body: formData,
        });
        const result = await response.json();
        if (result.code === 200) {
          message.success('导入成功');
          loadData();
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
        <h2>商家管理</h2>
        <div className="action-buttons">
          <Button className="action-btn-search" size="small" icon={<SearchOutlined />} onClick={() => setShowSearch(!showSearch)}>{showSearch ? '收起搜索' : '搜索'}</Button>
          <Button className="action-btn-export" size="small" icon={<ExportOutlined />} onClick={handleExport}>导出</Button>
          <Button className="action-btn-import" size="small" icon={<UploadOutlined />} onClick={handleImport}>导入</Button>
          <Button className="action-btn-edit" size="small" icon={<EditOutlined />} onClick={handleBatchEdit} disabled={selectedRowKeys.length === 0}>编辑({selectedRowKeys.length})</Button>
          <Button className="action-btn-delete" size="small" icon={<DeleteOutlined />} onClick={handleBatchDelete} disabled={selectedRowKeys.length === 0}>删除({selectedRowKeys.length})</Button>
          <Button className="action-btn-add" size="small" icon={<PlusOutlined />} onClick={addMerch}>添加</Button>
        </div>
      </div>

      {showSearch && (
        <Form className="search-form" form={searchForm} layout="inline">
          <Form.Item name="account">
            <Input placeholder="账号" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="email">
            <Input placeholder="邮箱" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="phone">
            <Input placeholder="手机号" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="role">
            <Select placeholder="角色" allowClear>
              <Select.Option value="商家">商家</Select.Option>
              <Select.Option value="管理者">管理者</Select.Option>
              <Select.Option value="代理商">代理商</Select.Option>
            </Select>
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
            onSelect: (record: Merch, selected: boolean) => {
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

      <Modal
        title={isEdit ? '编辑商家' : '添加商家'}
        open={modalVisible}
        onOk={saveMerch}
        onCancel={() => setModalVisible(false)}
        okText="确定"
        cancelText="取消"
        className="form-modal"
      >
        <Form form={form} layout="vertical">
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="account" label="账号" rules={[{ required: true, message: '请输入账号' }, { min: 3, max: 50, message: '账号长度在3-50之间' }]}>
                <Input placeholder="请输入账号" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="password" label="密码" rules={[{ required: !isEdit, message: '请输入密码' }, { min: 6, max: 50, message: '密码长度在6-50之间' }]}>
                <Input.Password placeholder="请输入密码" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="email" label="邮箱" rules={[{ type: 'email', message: '请输入正确的邮箱格式' }]}>
                <Input placeholder="请输入邮箱" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="phone" label="手机号" rules={[{ pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式' }]}>
                <Input placeholder="请输入手机号" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="role" label="角色" rules={[{ required: true, message: '请选择角色' }]}>
                <Select options={[{ label: '商家', value: '商家' }, { label: '管理者', value: '管理者' }, { label: '代理商', value: '代理商' }]} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="status" label="状态" rules={[{ required: true, message: '请选择状态' }]}>
                <Select options={[{ label: '白名单', value: '0' }, { label: '黑名单', value: '1' }]} />
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
              <Form.Item name="role" label="角色">
                <Select placeholder="请选择角色" allowClear>
                  <Select.Option value="商家">商家</Select.Option>
                  <Select.Option value="管理者">管理者</Select.Option>
                  <Select.Option value="代理商">代理商</Select.Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="status" label="状态">
                <Select placeholder="请选择状态" allowClear>
                  <Select.Option value="0">白名单</Select.Option>
                  <Select.Option value="1">黑名单</Select.Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>

      <Modal title="商家详情" open={detailVisible} onCancel={() => setDetailVisible(false)} okText="确定" cancelText="取消" footer={null} className="detail-modal">
        <Spin spinning={detailLoading}>
          {currentMerch && (
          <div>
            <div className="detail-header">
              <span className="detail-title">{currentMerch.account}</span>
              <Tag className="detail-tag" color={currentMerch.status === '0' ? 'green' : 'red'}>
                {currentMerch.status === '0' ? '白名单' : '黑名单'}
              </Tag>
            </div>
            <div className="detail-content">
              <div className="detail-grid">
              <div className="detail-item">
                <div className="detail-item-label">账号</div>
                <div className="detail-item-value">{currentMerch.account}</div>
              </div>
              <div className="detail-item">
                <div className="detail-item-label">角色</div>
                <div className="detail-item-value">{currentMerch.role}</div>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">ID</div>
              <div className="detail-value">{currentMerch.id}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">邮箱</div>
              <div className="detail-value">{currentMerch.email}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">手机号</div>
              <div className="detail-value">{currentMerch.phone}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">状态</div>
              <div className="detail-value">
                <Tag color={currentMerch.status === '0' ? 'green' : 'red'}>
                  {currentMerch.status === '0' ? '白名单' : '黑名单'}
                </Tag>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">上次登录</div>
              <div className="detail-value">{currentMerch.logAt || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">创建时间</div>
              <div className="detail-value">{currentMerch.createdAt}</div>
            </div>
            </div>
          </div>
        )}
        </Spin>
      </Modal>
    </div>
  );
};