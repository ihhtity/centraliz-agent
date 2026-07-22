import { Table, Button, Modal, Form, Input, Select, message } from 'antd';
import { SearchOutlined, DeleteOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import { getDeviceLogList, batchDeleteDeviceLog } from '@/api';
import type { DeviceLog } from '@/types';

const { Option } = Select;

export const DeviceLogManage = () => {
  const [data, setData] = useState<DeviceLog[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '设备编码', dataIndex: 'code', key: 'code' },
    { title: '设备名称', dataIndex: 'deviceName', key: 'deviceName' },
    { title: '房间名称', dataIndex: 'roomName', key: 'roomName' },
    { title: '类型', dataIndex: 'type', key: 'type' },
    { title: '控制', dataIndex: 'control', key: 'control' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '占用者', dataIndex: 'occupant', key: 'occupant' },
    { title: '电话', dataIndex: 'phone', key: 'phone' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt' },
  ];

  const fetchData = async (params: { page?: number; page_size?: number; deviceName?: string; status?: string } = {}) => {
    setLoading(true);
    try {
      const res = await getDeviceLogList({ page: params.page || 1, page_size: params.page_size || 10, deviceName: params.deviceName, status: params.status });
      setData(res.data.data);
      setTotal(res.data.total);
    } catch (error) {
      message.error('获取设备日志列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleBatchDelete = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要删除的日志');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 条日志吗？`,
      onOk: async () => {
        try {
          await batchDeleteDeviceLog({ ids: selectedRowKeys.map(k => k.toString()) });
          message.success('批量删除成功');
          setSelectedRowKeys([]);
          fetchData();
        } catch (error) {
          message.error('批量删除失败');
        }
      },
    });
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
        <h2>设备日志</h2>
        <div className="action-buttons">
          <Button danger icon={<DeleteOutlined />} onClick={handleBatchDelete}>批量删除</Button>
        </div>
      </div>

      <Form className="search-form" form={searchForm} layout="inline">
        <Form.Item name="deviceName">
          <Input placeholder="设备名称" prefix={<SearchOutlined />} />
        </Form.Item>
        <Form.Item name="status">
          <Select placeholder="状态" allowClear>
            <Option value="success">成功</Option>
            <Option value="fail">失败</Option>
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
          scroll={{ x: 1200 }}
        />
      </div>
    </div>
  );
};