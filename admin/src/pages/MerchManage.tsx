import { Table, Button, Modal, Form, Input, Select, message, Tag, Space } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import type { Merch } from '@/types';
import {
  getMerchList,
  getMerchDetail,
  createMerch,
  updateMerch,
  deleteMerch,
  batchDeleteMerch,
} from '@/api';
import { SearchBar } from '@/components/SearchBar';
import { BatchActions } from '@/components/BatchActions';

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
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [searchParams, setSearchParams] = useState<Record<string, any>>({});
  const [modalVisible, setModalVisible] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const [currentMerch, setCurrentMerch] = useState<Merch | null>(null);
  const [detailVisible, setDetailVisible] = useState(false);
  const [form] = Form.useForm();

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
    { title: '上次登录', dataIndex: 'logAt', key: 'logAt' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: Merch) => (
        <Space>
          <Button icon={<EyeOutlined />} onClick={() => viewDetail(record)} size="small">详情</Button>
          <Button icon={<EditOutlined />} onClick={() => editMerch(record)} size="small">编辑</Button>
          <Button icon={<DeleteOutlined />} danger onClick={() => deleteMerchItem(record.id)} size="small">删除</Button>
        </Space>
      ),
    },
  ];

  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getMerchList({
        ...searchParams,
        page,
        page_size: pageSize,
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
  }, [page, pageSize, searchParams]);

  const handleSearch = (values: Record<string, any>) => {
    setSearchParams(values);
    setPage(1);
  };

  const viewDetail = async (merch: Merch) => {
    try {
      const res = await getMerchDetail(merch.id);
      setCurrentMerch(res.data);
      setDetailVisible(true);
    } catch (error) {
      console.error(error);
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
    try {
      await deleteMerch(id);
      message.success('删除成功');
      loadData();
    } catch (error) {
      console.error(error);
    }
  };

  const handleBatchDelete = async () => {
    try {
      await batchDeleteMerch({ ids: selectedRowKeys.map(String) });
      message.success('批量删除成功');
      setSelectedRowKeys([]);
      loadData();
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
        <h2>商家管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={addMerch}>添加商家</Button>
      </div>
      <SearchBar
        onSearch={handleSearch}
        filters={[
          { type: 'input', key: 'account', label: '账号', placeholder: '请输入账号' },
          { type: 'input', key: 'phone', label: '手机号', placeholder: '请输入手机号' },
          {
            type: 'select',
            key: 'role',
            label: '角色',
            options: [
              { label: '商家', value: '商家' },
              { label: '管理者', value: '管理者' },
              { label: '代理商', value: '代理商' },
            ],
          },
          {
            type: 'select',
            key: 'status',
            label: '状态',
            options: [
              { label: '白名单', value: '0' },
              { label: '黑名单', value: '1' },
            ],
          },
        ]}
      />
      <BatchActions
        selectedRowKeys={selectedRowKeys}
        onBatchDelete={handleBatchDelete}
        data={data}
        columns={columns}
      />
      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          onChange: (p, ps) => { setPage(p); setPageSize(ps); },
        }}
        rowSelection={{
          selectedRowKeys,
          onChange: setSelectedRowKeys,
        }}
      />

      <Modal
        title={isEdit ? '编辑商家' : '添加商家'}
        visible={modalVisible}
        onOk={saveMerch}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="account" label="账号" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: !isEdit }]}>
            <Input.Password />
          </Form.Item>
          <Form.Item name="email" label="邮箱">
            <Input />
          </Form.Item>
          <Form.Item name="phone" label="手机号">
            <Input />
          </Form.Item>
          <Form.Item name="role" label="角色">
            <Select options={[{ label: '商家', value: '商家' }, { label: '管理者', value: '管理者' }, { label: '代理商', value: '代理商' }]} />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select options={[{ label: '白名单', value: '0' }, { label: '黑名单', value: '1' }]} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal title="商家详情" visible={detailVisible} onCancel={() => setDetailVisible(false)} footer={null}>
        {currentMerch && (
          <div>
            <p><strong>ID:</strong> {currentMerch.id}</p>
            <p><strong>账号:</strong> {currentMerch.account}</p>
            <p><strong>邮箱:</strong> {currentMerch.email}</p>
            <p><strong>手机号:</strong> {currentMerch.phone}</p>
            <p><strong>角色:</strong> {currentMerch.role}</p>
            <p><strong>状态:</strong> {currentMerch.status === '0' ? '白名单' : '黑名单'}</p>
            <p><strong>上次登录:</strong> {currentMerch.logAt}</p>
            <p><strong>创建时间:</strong> {currentMerch.createdAt}</p>
          </div>
        )}
      </Modal>
    </div>
  );
};
