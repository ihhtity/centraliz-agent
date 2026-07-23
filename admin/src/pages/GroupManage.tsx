import { Table, Button, Modal, Form, Input, Select, message, Row, Col, Tag, Space, Spin } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined, ExportOutlined, EyeOutlined, UploadOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import * as XLSX from 'xlsx';
import { getGroupList, getGroupDetail, createGroup, updateGroup, deleteGroup, batchDeleteGroup, batchUpdateGroup } from '@/api';
import { CustomPagination } from '@/components/CustomPagination';
import type { Group } from '@/types';

const { Option } = Select;

export const GroupManage = () => {
  const [data, setData] = useState<Group[]>([]);
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
  const [currentGroup, setCurrentGroup] = useState<Group | null>(null);

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
    {
      title: '消费推送',
      dataIndex: 'consumePush',
      key: 'consumePush',
      render: (val: string) => (
        <Tag color={val === '1' ? 'green' : 'red'}>{val === '1' ? '开启' : '关闭'}</Tag>
      ),
    },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (time: string) => formatTime(time) },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Group) => (
        <Space>
          <Button className="action-btn-detail" icon={<EyeOutlined />} onClick={() => viewDetail(record)} size="small">详情</Button>
          <Button className="action-btn-edit" icon={<EditOutlined />} onClick={() => handleEdit(record)} size="small">编辑</Button>
          <Button className="action-btn-delete" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)} size="small">删除</Button>
        </Space>
      ),
    },
  ];

  const fetchData = async (params: { page?: number; page_size?: number; name?: string; type?: string } = {}) => {
    setLoading(true);
    try {
      const res = await getGroupList({ page: params.page || currentPage, page_size: params.page_size || pageSize, name: params.name, type: params.type });
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

  const viewDetail = async (group: Group) => {
    setDetailLoading(true);
    try {
      const res = await getGroupDetail(group.id);
      setCurrentGroup(res.data);
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

  const handleEdit = (record: Group) => {
    setEditId(record.id);
    form.setFieldsValue(record);
    setIsModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个分组吗？',
      okText: '确定',
      cancelText: '取消',
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
      okText: '确定',
      cancelText: '取消',
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

  const handleBatchEdit = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要编辑的分组');
      return;
    }
    batchForm.resetFields();
    setIsBatchModalVisible(true);
  };

  const handleBatchSubmit = async () => {
    try {
      const values = await batchForm.validateFields();
      await batchUpdateGroup({ ids: selectedRowKeys.map(k => k.toString()), data: values });
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

  const handleExport = () => {
    if (data.length === 0) {
      message.warning('暂无数据可导出');
      return;
    }
    const exportData = data.map(item => ({
      ID: item.id,
      分组名称: item.name,
      电话: item.phone,
      房间数量: item.roomCount,
      类型: item.type,
      位置: item.location,
      地址: item.address,
      规则名称: item.ruleName,
      绑定编号: item.bindNumber,
      消费推送: item.consumePush,
      创建时间: item.createdAt,
    }));
    const ws = XLSX.utils.json_to_sheet(exportData);
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, '分组');
    XLSX.writeFile(wb, `分组_${new Date().toLocaleDateString()}.xlsx`);
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
        const response = await fetch('/admin/group/import', {
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
        <h2>分组管理</h2>
        <div className="action-buttons">
          <Button className="action-btn-search" size="small" icon={<SearchOutlined />} onClick={() => setShowSearch(!showSearch)}>{showSearch ? '收起搜索' : '搜索'}</Button>
          <Button className="action-btn-export" size="small" icon={<ExportOutlined />} onClick={handleExport}>导出</Button>
          <Button className="action-btn-import" size="small" icon={<UploadOutlined />} onClick={handleImport}>导入</Button>
          <Button className="action-btn-edit" size="small" icon={<EditOutlined />} onClick={handleBatchEdit} disabled={selectedRowKeys.length === 0}>编辑({selectedRowKeys.length})</Button>
          <Button className="action-btn-delete" size="small" icon={<DeleteOutlined />} onClick={handleBatchDelete} disabled={selectedRowKeys.length === 0}>删除({selectedRowKeys.length})</Button>
          <Button className="action-btn-add" size="small" icon={<PlusOutlined />} onClick={handleAdd}>添加</Button>
        </div>
      </div>

      {showSearch && (
        <Form className="search-form" form={searchForm} layout="inline">
          <Form.Item name="name">
            <Input placeholder="分组名称" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="phone">
            <Input placeholder="电话" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="type">
            <Select placeholder="分组类型" allowClear>
              <Option value="hotel">酒店</Option>
              <Option value="office">办公室</Option>
            </Select>
          </Form.Item>
          <Form.Item name="location">
            <Input placeholder="位置" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="address">
            <Input placeholder="地址" prefix={<SearchOutlined />} />
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
          scroll={{ x: 1300 }}
        />
        <CustomPagination
          total={total}
          current={currentPage}
          pageSize={pageSize}
          onChange={(page, pageSize) => handleTableChange({ current: page, pageSize })}
        />
      </div>

      <Modal
        title={editId ? '编辑分组' : '新增分组'}
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
              <Form.Item name="name" label="分组名称" rules={[{ required: true, message: '请输入分组名称' }, { max: 100, message: '名称长度不超过100' }]}>
                <Input placeholder="请输入分组名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="phone" label="电话" rules={[{ pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式' }]}>
                <Input placeholder="请输入电话" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="roomCount" label="房间数量" rules={[{ min: 0, message: '房间数量不能为负数' }]}>
                <Input type="number" placeholder="请输入房间数量" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="type" label="类型" rules={[{ max: 50, message: '类型长度不超过50' }]}>
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
          <Form.Item name="type" label="类型">
            <Input placeholder="请输入类型" />
          </Form.Item>
          <Form.Item name="location" label="位置">
            <Input placeholder="请输入位置" />
          </Form.Item>
          <Form.Item name="consumePush" label="消费推送">
            <Select placeholder="请选择消费推送" allowClear>
              <Option value="1">开启</Option>
              <Option value="0">关闭</Option>
            </Select>
          </Form.Item>
          <Form.Item name="ruleName" label="规则名称">
            <Input placeholder="请输入规则名称" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal title="分组详情" open={detailVisible} onCancel={() => setDetailVisible(false)} okText="确定" cancelText="取消" footer={null} className="detail-modal">
        <Spin spinning={detailLoading}>
          {currentGroup && (
          <div>
            <div className="detail-header">
              <span className="detail-title">{currentGroup.name}</span>
              <Tag className="detail-tag" color={currentGroup.consumePush === '1' ? 'green' : 'red'}>
                {currentGroup.consumePush === '1' ? '消费推送开启' : '消费推送关闭'}
              </Tag>
            </div>
            <div className="detail-content">
              <div className="detail-grid">
              <div className="detail-item">
                <div className="detail-item-label">分组名称</div>
                <div className="detail-item-value">{currentGroup.name}</div>
              </div>
              <div className="detail-item">
                <div className="detail-item-label">房间数量</div>
                <div className="detail-item-value">{currentGroup.roomCount}</div>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">ID</div>
              <div className="detail-value">{currentGroup.id}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">电话</div>
              <div className="detail-value">{currentGroup.phone || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">类型</div>
              <div className="detail-value">{currentGroup.type || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">位置</div>
              <div className="detail-value">{currentGroup.location || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">地址</div>
              <div className="detail-value">{currentGroup.address || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">规则名称</div>
              <div className="detail-value">{currentGroup.ruleName || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">绑定编号</div>
              <div className="detail-value">{currentGroup.bindNumber || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">消费推送</div>
              <div className="detail-value">
                <Tag color={currentGroup.consumePush === '1' ? 'green' : 'red'}>
                  {currentGroup.consumePush === '1' ? '开启' : '关闭'}
                </Tag>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">创建时间</div>
              <div className="detail-value">{currentGroup.createdAt}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">更新时间</div>
              <div className="detail-value">{currentGroup.updatedAt}</div>
            </div>
            </div>
          </div>
        )}
        </Spin>
      </Modal>
    </div>
  );
};