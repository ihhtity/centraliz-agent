import { Table, Button, Modal, Form, Input, Select, message, Tag, Space, Spin, InputNumber } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined, SearchOutlined, ExportOutlined, UploadOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import * as XLSX from 'xlsx';
import type { HuifuAccount } from '@/types';
import {
  getHuifuAccountList,
  getHuifuAccountDetail,
  createHuifuAccount,
  updateHuifuAccount,
  deleteHuifuAccount,
  batchDeleteHuifuAccount,
  batchUpdateHuifuAccount,
} from '@/api';
import { CustomPagination } from '@/components/CustomPagination';

const chooseColors: Record<string, string> = {
  '0': 'gray',
  '1': 'green',
};

const shareColors: Record<string, string> = {
  '0': 'gray',
  '1': 'blue',
};

export const HuifuAccountManage = () => {
  const [data, setData] = useState<HuifuAccount[]>([]);
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
  const [currentItem, setCurrentItem] = useState<HuifuAccount | null>(null);
  const [detailVisible, setDetailVisible] = useState(false);
  const [form] = Form.useForm();
  const [batchForm] = Form.useForm();
  const [searchForm] = Form.useForm();

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '汇付编码', dataIndex: 'code', key: 'code' },
    { title: '账号', dataIndex: 'account', key: 'account' },
    { title: '姓名', dataIndex: 'name', key: 'name' },
    { title: '手机号', dataIndex: 'phone', key: 'phone' },
    { title: '店名', dataIndex: 'storename', key: 'storename' },
    { title: '账号类型', dataIndex: 'type', key: 'type' },
    {
      title: '选择状态',
      dataIndex: 'choose',
      key: 'choose',
      render: (choose: string) => (
        <Tag color={chooseColors[choose] || 'default'}>{choose === '0' ? '未选择' : '已选择'}</Tag>
      ),
    },
    {
      title: '分账状态',
      dataIndex: 'share',
      key: 'share',
      render: (share: string) => (
        <Tag color={shareColors[share] || 'default'}>{share === '0' ? '关闭分账' : '开启分账'}</Tag>
      ),
    },
    { title: '分账比率', dataIndex: 'rate', key: 'rate' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (time: string) => formatTime(time) },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: HuifuAccount) => (
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
      const res = await getHuifuAccountList({
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

  const viewDetail = async (item: HuifuAccount) => {
    setDetailLoading(true);
    try {
      const res = await getHuifuAccountDetail(item.id);
      setCurrentItem(res.data);
      setDetailVisible(true);
    } catch (error) {
      console.error(error);
      message.error('获取详情失败');
    } finally {
      setDetailLoading(false);
    }
  };

  const editItem = (item: HuifuAccount) => {
    setCurrentItem(item);
    setIsEdit(true);
    form.setFieldsValue(item);
    setModalVisible(true);
  };

  const addItem = () => {
    setIsEdit(false);
    setCurrentItem(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleBatchEdit = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要编辑的汇付账号');
      return;
    }
    batchForm.resetFields();
    setIsBatchModalVisible(true);
  };

  const handleBatchSubmit = async () => {
    try {
      const values = await batchForm.validateFields();
      await batchUpdateHuifuAccount({ ids: selectedRowKeys.map(k => k.toString()), data: values });
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
        await updateHuifuAccount(currentItem.id, values);
        message.success('更新成功');
      } else {
        await createHuifuAccount(values);
        message.success('创建成功');
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
      content: '确定要删除这个汇付账号吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await deleteHuifuAccount(id);
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
      message.warning('请选择要删除的汇付账号');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 个汇付账号吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteHuifuAccount({ ids: selectedRowKeys.map(String) });
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
      汇付编码: item.code,
      账号: item.account,
      姓名: item.name,
      手机号: item.phone,
      店名: item.storename,
      账号类型: item.type,
      选择状态: item.choose === '0' ? '未选择' : '已选择',
      分账状态: item.share === '0' ? '关闭分账' : '开启分账',
      分账比率: item.rate,
      创建时间: item.createdAt,
    }));
    const ws = XLSX.utils.json_to_sheet(exportData);
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, '汇付账号');
    XLSX.writeFile(wb, `汇付账号_${new Date().toLocaleDateString()}.xlsx`);
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
        const response = await fetch('/admin/huifu/import', {
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
        <h2>汇付账号管理</h2>
        <div className="action-buttons">
          <Button className="action-btn-search" size="small" icon={<SearchOutlined />} onClick={() => setShowSearch(!showSearch)}>{showSearch ? '收起搜索' : '搜索'}</Button>
          <Button className="action-btn-export" size="small" icon={<ExportOutlined />} onClick={handleExport}>导出</Button>
          <Button className="action-btn-import" size="small" icon={<UploadOutlined />} onClick={handleImport}>导入</Button>
          <Button className="action-btn-edit" size="small" icon={<EditOutlined />} onClick={handleBatchEdit} disabled={selectedRowKeys.length === 0}>编辑({selectedRowKeys.length})</Button>
          <Button className="action-btn-delete" size="small" icon={<DeleteOutlined />} onClick={handleBatchDelete} disabled={selectedRowKeys.length === 0}>删除({selectedRowKeys.length})</Button>
          <Button className="action-btn-add" size="small" icon={<PlusOutlined />} onClick={addItem}>添加</Button>
        </div>
      </div>

      {showSearch && (
        <Form className="search-form" form={searchForm} layout="inline">
          <Form.Item name="code">
            <Input placeholder="汇付编码" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="account">
            <Input placeholder="账号" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="name">
            <Input placeholder="姓名" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="phone">
            <Input placeholder="手机号" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="storename">
            <Input placeholder="店名" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="type">
            <Select placeholder="账号类型" allowClear />
          </Form.Item>
          <Form.Item name="choose">
            <Select placeholder="选择状态" allowClear>
              <Select.Option value="0">未选择</Select.Option>
              <Select.Option value="1">已选择</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="share">
            <Select placeholder="分账状态" allowClear>
              <Select.Option value="0">关闭分账</Select.Option>
              <Select.Option value="1">开启分账</Select.Option>
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
        title={isEdit ? '编辑汇付账号' : '添加汇付账号'}
        open={modalVisible}
        onOk={saveItem}
        onCancel={() => setModalVisible(false)}
        okText="确定"
        cancelText="取消"
        className="form-modal"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="code" label="汇付编码" rules={[{ max: 255, message: '汇付编码长度不超过255' }]}>
            <Input placeholder="请输入汇付编码" />
          </Form.Item>
          <Form.Item name="account" label="账号" rules={[{ max: 255, message: '账号长度不超过255' }]}>
            <Input placeholder="请输入账号" />
          </Form.Item>
          <Form.Item name="name" label="姓名" rules={[{ max: 255, message: '姓名长度不超过255' }]}>
            <Input placeholder="请输入姓名" />
          </Form.Item>
          <Form.Item name="phone" label="手机号" rules={[{ pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式' }]}>
            <Input placeholder="请输入手机号" />
          </Form.Item>
          <Form.Item name="identity" label="身份证" rules={[{ pattern: /^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$/, message: '请输入正确的身份证格式' }]}>
            <Input placeholder="请输入身份证" />
          </Form.Item>
          <Form.Item name="card" label="银行卡" rules={[{ max: 255, message: '银行卡号长度不超过255' }]}>
            <Input placeholder="请输入银行卡" />
          </Form.Item>
          <Form.Item name="storename" label="店名" rules={[{ max: 255, message: '店名长度不超过255' }]}>
            <Input placeholder="请输入店名" />
          </Form.Item>
          <Form.Item name="area" label="经营地址">
            <Input.TextArea placeholder="请输入经营地址" />
          </Form.Item>
          <Form.Item name="type" label="账号类型">
            <Select placeholder="请选择账号类型" />
          </Form.Item>
          <Form.Item name="choose" label="选择状态">
            <Select options={[{ label: '未选择', value: '0' }, { label: '已选择', value: '1' }]} />
          </Form.Item>
          <Form.Item name="share" label="分账状态">
            <Select options={[{ label: '关闭分账', value: '0' }, { label: '开启分账', value: '1' }]} />
          </Form.Item>
          <Form.Item name="rate" label="分账比率" rules={[{ min: 0, max: 100, message: '分账比率在0-100之间' }]}>
            <InputNumber placeholder="请输入分账比率" style={{ width: '100%' }} />
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
          <Form.Item name="type" label="账号类型">
            <Select placeholder="请选择账号类型" allowClear />
          </Form.Item>
          <Form.Item name="choose" label="选择状态">
            <Select placeholder="请选择选择状态" allowClear>
              <Select.Option value="0">未选择</Select.Option>
              <Select.Option value="1">已选择</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="share" label="分账状态">
            <Select placeholder="请选择分账状态" allowClear>
              <Select.Option value="0">关闭分账</Select.Option>
              <Select.Option value="1">开启分账</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>

        <Modal title="汇付账号详情" open={detailVisible} onCancel={() => setDetailVisible(false)} okText="确定" cancelText="取消" footer={null} className="detail-modal">
        <Spin spinning={detailLoading}>
          {currentItem && (
            <div>
              <div className="detail-header">
                <span className="detail-title">{currentItem.account}</span>
                <Tag className="detail-tag" color={currentItem.choose === '0' ? 'gray' : 'green'}>
                  {currentItem.choose === '0' ? '未选择' : '已选择'}
                </Tag>
              </div>
              <div className="detail-content">
                <div className="detail-grid">
                  <div className="detail-item">
                    <div className="detail-item-label">汇付编码</div>
                    <div className="detail-item-value">{currentItem.code}</div>
                  </div>
                  <div className="detail-item">
                    <div className="detail-item-label">姓名</div>
                    <div className="detail-item-value">{currentItem.name}</div>
                  </div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">ID</div>
                  <div className="detail-value">{currentItem.id}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">手机号</div>
                  <div className="detail-value">{currentItem.phone}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">身份证</div>
                  <div className="detail-value">{currentItem.identity || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">银行卡</div>
                  <div className="detail-value">{currentItem.card || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">店名</div>
                  <div className="detail-value">{currentItem.storename || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">经营地址</div>
                  <div className="detail-value">{currentItem.area || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">账号类型</div>
                  <div className="detail-value">{currentItem.type || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">分账状态</div>
                  <div className="detail-value">
                    <Tag color={currentItem.share === '0' ? 'gray' : 'blue'}>
                      {currentItem.share === '0' ? '关闭分账' : '开启分账'}
                    </Tag>
                  </div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">分账比率</div>
                  <div className="detail-value">{currentItem.rate}</div>
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
    </div>
  );
};