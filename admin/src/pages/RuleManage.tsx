import { Table, Button, Modal, Form, Input, Select, InputNumber, Checkbox, message, Row, Col, Tag, Space, Spin } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined, ExportOutlined, EyeOutlined, UploadOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import * as XLSX from 'xlsx';
import { getRuleList, getRuleDetail, createRule, updateRule, deleteRule, batchDeleteRule, batchUpdateRule } from '@/api';
import { CustomPagination } from '@/components/CustomPagination';
import type { Rule } from '@/types';

const { Option } = Select;

export const RuleManage = () => {
  const [data, setData] = useState<Rule[]>([]);
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
  const [currentRule, setCurrentRule] = useState<Rule | null>(null);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '规则名称', dataIndex: 'name', key: 'name' },
    { title: '类型', dataIndex: 'type', key: 'type' },
    { title: '模式', dataIndex: 'mode', key: 'mode' },
    { title: '价格', dataIndex: 'price', key: 'price', render: (val: number) => <span>¥{val}</span> },
    { title: '押金', dataIndex: 'deposit', key: 'deposit', render: (val: number) => <span>¥{val}</span> },
    { title: '费率', dataIndex: 'rate', key: 'rate' },
    { title: '时长', dataIndex: 'duration', key: 'duration' },
    { title: '时长单位', dataIndex: 'durationUnit', key: 'durationUnit' },
    { title: '免费时长', dataIndex: 'freeTime', key: 'freeTime' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) => (
        <Tag color={status === 1 ? 'green' : 'red'}>{status === 1 ? '启用' : '禁用'}</Tag>
      ),
    },
    { title: '商家ID', dataIndex: 'merchsId', key: 'merchsId', width: 80 },
    { title: '排序', dataIndex: 'sort', key: 'sort', width: 60 },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (val: string) => formatTime(val) },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Rule) => (
        <Space>
          <Button className="action-btn-detail" icon={<EyeOutlined />} onClick={() => viewDetail(record)} size="small">详情</Button>
          <Button className="action-btn-edit" icon={<EditOutlined />} onClick={() => handleEdit(record)} size="small">编辑</Button>
          <Button className="action-btn-delete" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)} size="small">删除</Button>
        </Space>
      ),
    },
  ];

  const fetchData = async (params: { page?: number; page_size?: number; name?: string } = {}) => {
    setLoading(true);
    try {
      const res = await getRuleList({ page: params.page || currentPage, page_size: params.page_size || pageSize, name: params.name });
      setData(res.data.data);
      setTotal(res.data.total);
    } catch (error) {
      message.error('获取规则列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const viewDetail = async (rule: Rule) => {
    setDetailLoading(true);
    try {
      const res = await getRuleDetail(rule.id);
      setCurrentRule(res.data);
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

  const handleEdit = (record: Rule) => {
    setEditId(record.id);
    form.setFieldsValue(record);
    setIsModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个规则吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await deleteRule(id);
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
      message.warning('请选择要删除的规则');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 个规则吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteRule({ ids: selectedRowKeys.map(k => k.toString()) });
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
      message.warning('请选择要编辑的规则');
      return;
    }
    batchForm.resetFields();
    setIsBatchModalVisible(true);
  };

  const handleBatchSubmit = async () => {
    try {
      const values = await batchForm.validateFields();
      await batchUpdateRule({ ids: selectedRowKeys.map(k => k.toString()), data: values });
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
        await updateRule(editId, values);
        message.success('更新成功');
      } else {
        await createRule(values);
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
      规则名称: item.name,
      类型: item.type,
      模式: item.mode,
      价格: item.price,
      押金: item.deposit,
      费率: item.rate,
      时长: item.duration,
      时长单位: item.durationUnit,
      免费时长: item.freeTime,
      状态: item.status === 1 ? '启用' : '禁用',
    }));
    const ws = XLSX.utils.json_to_sheet(exportData);
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, '规则');
    XLSX.writeFile(wb, `规则_${new Date().toLocaleDateString()}.xlsx`);
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
        const response = await fetch('/admin/rule/import', {
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
    onSelect: (record: Rule, selected: boolean) => {
      if (selected) {
        setSelectedRowKeys([...selectedRowKeys, record.id]);
      } else {
        setSelectedRowKeys(selectedRowKeys.filter(k => k !== record.id));
      }
    },
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
        <h2>规则管理</h2>
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
            <Input placeholder="规则名称" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="tag">
            <Input placeholder="标签" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="type">
            <Select placeholder="类型" allowClear>
              <Option value="free">免费模式</Option>
              <Option value="charge">收费模式</Option>
            </Select>
          </Form.Item>
          <Form.Item name="mode">
            <Select placeholder="模式" allowClear>
              <Option value="single">单次开锁</Option>
              <Option value="deposit">一存一取</Option>
              <Option value="pay_single">单次付费</Option>
              <Option value="pay_deposit">先存后取</Option>
              <Option value="pay_hourly">按时付费</Option>
              <Option value="pay_time">预付费</Option>
            </Select>
          </Form.Item>
          <Form.Item name="status">
            <Select placeholder="状态" allowClear>
              <Option value={1}>启用</Option>
              <Option value={0}>禁用</Option>
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
          scroll={{ x: 1300 }}
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
        title={editId ? '编辑规则' : '新增规则'}
        open={isModalVisible}
        onOk={handleSubmit}
        onCancel={() => setIsModalVisible(false)}
        okText="确定"
        cancelText="取消"
        className="form-modal"
        width={600}
      >
        <Form form={form} layout="vertical">
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="name" label="规则名称" rules={[{ required: true, message: '请输入规则名称' }, { max: 100, message: '名称长度不超过100' }]}>
                <Input placeholder="请输入规则名称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="type" label="类型" initialValue="free">
                <Select placeholder="请选择类型">
                  <Option value="free">免费模式</Option>
                  <Option value="charge">收费模式</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="mode" label="模式" initialValue="single">
                <Select placeholder="请选择模式">
                  <Option value="single">单次开锁</Option>
                  <Option value="deposit">一存一取</Option>
                  <Option value="pay_single">单次付费</Option>
                  <Option value="pay_deposit">先存后取</Option>
                  <Option value="pay_hourly">按时付费</Option>
                  <Option value="pay_time">预付费</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="status" label="状态" initialValue={1}>
                <Select placeholder="请选择状态">
                  <Option value={1}>启用</Option>
                  <Option value={0}>禁用</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="price" label="价格" rules={[{ required: true, message: '请输入价格' }, { min: 0, message: '价格不能为负数' }]}>
                <InputNumber placeholder="请输入价格" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="deposit" label="押金" rules={[{ min: 0, message: '押金不能为负数' }]}>
                <InputNumber placeholder="请输入押金" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="rate" label="费率">
                <InputNumber placeholder="请输入费率" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="duration" label="时长">
                <InputNumber placeholder="请输入时长" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="durationUnit" label="时长单位" initialValue="hour">
                <Select>
                  <Option value="hour">小时</Option>
                  <Option value="day">天</Option>
                  <Option value="month">月</Option>
                  <Option value="minute">分钟</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="freeTime" label="免费时长(分钟)">
                <InputNumber placeholder="请输入免费时长" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="autoEndTime" label="自动结束时间">
                <InputNumber placeholder="请输入自动结束时间" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="sort" label="排序">
                <InputNumber placeholder="请输入排序" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="autoRefund" label="自动退款" valuePropName="checked">
                <Checkbox>自动退款</Checkbox>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="manualRenew" label="手动续费" valuePropName="checked">
                <Checkbox>手动续费</Checkbox>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="merchsId" label="商家外键">
                <InputNumber placeholder="请输入商家ID" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="tag" label="标签">
                <Input placeholder="请输入标签" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="description" label="描述">
                <Input.TextArea placeholder="请输入描述" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="timeOptions" label="时间选项(JSON)">
                <Input.TextArea placeholder="请输入时间选项JSON" rows={3} />
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
              <Form.Item name="status" label="状态">
                <Select placeholder="请选择状态" allowClear>
                  <Option value={1}>启用</Option>
                  <Option value={0}>禁用</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="type" label="类型">
                <Select placeholder="请选择类型" allowClear>
                  <Option value="free">免费模式</Option>
                  <Option value="charge">收费模式</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="mode" label="模式">
                <Select placeholder="请选择模式" allowClear>
                  <Option value="single">单次开锁</Option>
                  <Option value="deposit">一存一取</Option>
                  <Option value="pay_single">单次付费</Option>
                  <Option value="pay_deposit">先存后取</Option>
                  <Option value="pay_hourly">按时付费</Option>
                  <Option value="pay_time">预付费</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="price" label="价格">
                <InputNumber placeholder="请输入价格" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="deposit" label="押金">
                <InputNumber placeholder="请输入押金" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="sort" label="排序">
                <InputNumber placeholder="请输入排序" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>

      <Modal title="规则详情" open={detailVisible} onCancel={() => setDetailVisible(false)} okText="确定" cancelText="取消" footer={null} width={600} className="detail-modal">
        <Spin spinning={detailLoading}>
          {currentRule && (
          <div>
            <div className="detail-header">
              <span className="detail-title">{currentRule.name}</span>
              <Tag className="detail-tag" color={currentRule.status === 1 ? 'green' : 'red'}>
                {currentRule.status === 1 ? '启用' : '禁用'}
              </Tag>
            </div>
            <div className="detail-content">
              <div className="detail-grid">
              <div className="detail-item">
                <div className="detail-item-label">规则名称</div>
                <div className="detail-item-value">{currentRule.name}</div>
              </div>
              <div className="detail-item">
                <div className="detail-item-label">价格</div>
                <div className="detail-item-value">¥{currentRule.price}</div>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">ID</div>
              <div className="detail-value">{currentRule.id}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">类型</div>
              <div className="detail-value">{currentRule.type || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">模式</div>
              <div className="detail-value">{currentRule.mode || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">押金</div>
              <div className="detail-value">¥{currentRule.deposit}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">费率</div>
              <div className="detail-value">{currentRule.rate || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">时长</div>
              <div className="detail-value">{currentRule.duration} {currentRule.durationUnit}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">免费时长</div>
              <div className="detail-value">{currentRule.freeTime} 分钟</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">状态</div>
              <div className="detail-value">
                <Tag color={currentRule.status === 1 ? 'green' : 'red'}>
                  {currentRule.status === 1 ? '启用' : '禁用'}
                </Tag>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">排序</div>
              <div className="detail-value">{currentRule.sort}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">自动退款</div>
              <div className="detail-value">
                <Tag color={currentRule.autoRefund ? 'green' : 'red'}>
                  {currentRule.autoRefund ? '是' : '否'}
                </Tag>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">手动续费</div>
              <div className="detail-value">
                <Tag color={currentRule.manualRenew ? 'green' : 'red'}>
                  {currentRule.manualRenew ? '是' : '否'}
                </Tag>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">描述</div>
              <div className="detail-value">{currentRule.description || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">标签</div>
              <div className="detail-value">{currentRule.tag || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">创建时间</div>
              <div className="detail-value">{currentRule.createdAt}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">更新时间</div>
              <div className="detail-value">{currentRule.updatedAt}</div>
            </div>
            </div>
          </div>
        )}
        </Spin>
      </Modal>
    </div>
  );
};