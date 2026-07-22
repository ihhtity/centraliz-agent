import { Table, Button, Modal, Form, Input, Select, InputNumber, Checkbox, message, Row, Col } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import { getRuleList, createRule, updateRule, deleteRule } from '@/api';
import type { Rule } from '@/types';

const { Option } = Select;

export const RuleManage = () => {
  const [data, setData] = useState<Rule[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [editId, setEditId] = useState<number | null>(null);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '规则名称', dataIndex: 'name', key: 'name' },
    { title: '类型', dataIndex: 'type', key: 'type' },
    { title: '模式', dataIndex: 'mode', key: 'mode' },
    { title: '价格', dataIndex: 'price', key: 'price' },
    { title: '押金', dataIndex: 'deposit', key: 'deposit' },
    { title: '费率', dataIndex: 'rate', key: 'rate' },
    { title: '时长', dataIndex: 'duration', key: 'duration' },
    { title: '时长单位', dataIndex: 'durationUnit', key: 'durationUnit' },
    { title: '免费时长', dataIndex: 'freeTime', key: 'freeTime' },
    { title: '状态', dataIndex: 'status', key: 'status', render: (status: number) => (status === 1 ? '启用' : '禁用') },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Rule) => (
        <div className="action-buttons">
          <Button type="primary" ghost size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
          <Button danger size="small" icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>删除</Button>
        </div>
      ),
    },
  ];

  const fetchData = async (params: { page?: number; page_size?: number; name?: string } = {}) => {
    setLoading(true);
    try {
      const res = await getRuleList({ page: params.page || 1, page_size: params.page_size || 10, name: params.name });
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
        <h2>规则管理</h2>
        <div className="action-buttons">
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>新增规则</Button>
        </div>
      </div>

      <Form className="search-form" form={searchForm} layout="inline">
        <Form.Item name="name">
          <Input placeholder="规则名称" prefix={<SearchOutlined />} />
        </Form.Item>
        <Form.Item>
          <Button type="primary" onClick={handleSearch}>搜索</Button>
          <Button onClick={handleReset} style={{ marginLeft: 8 }}>重置</Button>
        </Form.Item>
      </Form>

      <div className="table-container">
        <Table
          columns={columns}
          dataSource={data}
          pagination={{ total, pageSize: 10, showSizeChanger: true, showQuickJumper: true }}
          loading={loading}
          rowKey="id"
          scroll={{ x: 1300 }}
        />
      </div>

      <Modal
        title={editId ? '编辑规则' : '新增规则'}
        visible={isModalVisible}
        onOk={handleSubmit}
        onCancel={() => setIsModalVisible(false)}
        className="form-modal"
        width={600}
      >
        <Form form={form} layout="vertical">
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="name" label="规则名称" rules={[{ required: true, message: '请输入规则名称' }]}>
                <Input placeholder="请输入规则名称" />
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
              <Form.Item name="mode" label="模式">
                <Input placeholder="请输入模式" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="status" label="状态">
                <Select defaultValue={1}>
                  <Option value={1}>启用</Option>
                  <Option value={0}>禁用</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name="price" label="价格">
                <InputNumber placeholder="请输入价格" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="deposit" label="押金">
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
              <Form.Item name="durationUnit" label="时长单位">
                <Input placeholder="请输入时长单位" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="freeTime" label="免费时长">
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
          <Form.Item name="description" label="描述">
            <Input.TextArea placeholder="请输入描述" />
          </Form.Item>
          <Form.Item name="timeOptions" label="时间选项">
            <Input placeholder="请输入时间选项" />
          </Form.Item>
          <Form.Item name="tag" label="标签">
            <Input placeholder="请输入标签" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};