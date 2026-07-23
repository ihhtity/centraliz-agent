import { Table, Button, Modal, Form, Input, Select, message, Tag, Space, Spin } from 'antd';
import { SearchOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import { getDeviceLogList, getDeviceLogDetail, batchDeleteDeviceLog, importDeviceLog } from '@/api';
import { CustomPagination } from '@/components/CustomPagination';
import { ExportButton } from '@/components/ExportButton';
import type { DeviceLog } from '@/types';

const { Option } = Select;

const controlColors: Record<string, string> = {
  'open': 'green',
  'close': 'red',
};

const statusColors: Record<string, string> = {
  'success': 'green',
  'fail': 'red',
};

export const DeviceLogManage = () => {
  const [data, setData] = useState<DeviceLog[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [showSearch, setShowSearch] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [detailVisible, setDetailVisible] = useState(false);
  const [currentLog, setCurrentLog] = useState<DeviceLog | null>(null);

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '设备编码', dataIndex: 'code', key: 'code' },
    { title: '设备名称', dataIndex: 'deviceName', key: 'deviceName' },
    { title: '房间名称', dataIndex: 'roomName', key: 'roomName' },
    { title: '类型', dataIndex: 'type', key: 'type' },
    {
      title: '控制',
      dataIndex: 'control',
      key: 'control',
      render: (control: string) => (
        <Tag color={controlColors[control] || 'default'}>{control === 'open' ? '开锁' : '关锁'}</Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={statusColors[status] || 'default'}>{status === 'success' ? '成功' : '失败'}</Tag>
      ),
    },
    { title: '占用者', dataIndex: 'occupant', key: 'occupant' },
    { title: '电话', dataIndex: 'phone', key: 'phone' },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (time: string) => formatTime(time) },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: DeviceLog) => (
        <Space>
          <Button className="action-btn-detail" icon={<EyeOutlined />} onClick={() => viewDetail(record)} size="small">详情</Button>
          <Button className="action-btn-delete" icon={<DeleteOutlined />} onClick={() => deleteLogItem(record.id)} size="small">删除</Button>
        </Space>
      ),
    },
  ];

  const fetchData = async (params: { page?: number; page_size?: number; deviceName?: string; status?: string; control?: string; type?: string; code?: string } = {}) => {
    setLoading(true);
    try {
      const res = await getDeviceLogList({ 
        page: params.page || currentPage, 
        page_size: params.page_size || pageSize, 
        device_name: params.deviceName, 
        status: params.status,
        control: params.control,
        type: params.type,
        code: params.code,
      });
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

  const viewDetail = async (log: DeviceLog) => {
    setDetailLoading(true);
    try {
      const res = await getDeviceLogDetail(log.id);
      setCurrentLog(res.data);
      setDetailVisible(true);
    } catch (error) {
      console.error(error);
      message.error('获取详情失败');
    } finally {
      setDetailLoading(false);
    }
  };

  const deleteLogItem = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个日志吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteDeviceLog({ ids: [id.toString()] });
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
      message.warning('请选择要删除的日志');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 条日志吗？`,
      okText: '确定',
      cancelText: '取消',
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

  const handleImport = async (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    try {
      await importDeviceLog(formData);
      message.success('导入成功');
      fetchData();
    } catch (error) {
      message.error('导入失败');
    }
  };

  const exportHeaders: Record<string, string> = {
    id: 'ID',
    code: '设备编码',
    deviceName: '设备名称',
    roomName: '房间名称',
    type: '类型',
    control: '控制',
    status: '状态',
    occupant: '占用者',
    phone: '电话',
    createdAt: '创建时间',
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
        <h2>设备日志</h2>
        <div className="action-buttons">
          <Button className="action-btn-search" size="small" icon={<SearchOutlined />} onClick={() => setShowSearch(!showSearch)}>{showSearch ? '收起搜索' : '搜索'}</Button>
          <ExportButton data={data} filename="设备日志列表" headers={exportHeaders} />
          <Button className="action-btn-import" size="small" onClick={() => document.getElementById('device-log-import')?.click()}>导入</Button>
          <input type="file" id="device-log-import" style={{ display: 'none' }} onChange={(e) => e.target.files?.[0] && handleImport(e.target.files[0])} />
          <Button className="action-btn-delete" size="small" icon={<DeleteOutlined />} onClick={handleBatchDelete} disabled={selectedRowKeys.length === 0}>删除({selectedRowKeys.length})</Button>
        </div>
      </div>

      {showSearch && (
        <Form className="search-form" form={searchForm} layout="inline">
          <Form.Item name="code">
            <Input placeholder="设备编码" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="deviceName">
            <Input placeholder="设备名称" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="roomName">
            <Input placeholder="房间名称" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="occupant">
            <Input placeholder="使用者" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="phone">
            <Input placeholder="使用者电话" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="type">
            <Select placeholder="设备类型" allowClear>
              <Option value="lock">门锁</Option>
              <Option value="device">设备</Option>
            </Select>
          </Form.Item>
          <Form.Item name="control">
            <Select placeholder="控制类型" allowClear>
              <Option value="open">开锁</Option>
              <Option value="close">关锁</Option>
            </Select>
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
      )}

      <div className="table-container">
        <Table
          rowSelection={rowSelection}
          columns={columns}
          dataSource={data}
          pagination={false}
          loading={loading}
          rowKey="id"
          scroll={{ x: 1200 }}
        />
        <CustomPagination
          total={total}
          current={currentPage}
          pageSize={pageSize}
          onChange={(page, pageSize) => handleTableChange({ current: page, pageSize })}
        />
      </div>

      <Modal title="设备日志详情" open={detailVisible} onCancel={() => setDetailVisible(false)}
        okText="确定" cancelText="取消" footer={null} className="detail-modal">
        <Spin spinning={detailLoading}>
          {currentLog && (
          <div>
            <div className="detail-header">
              <span className="detail-title">{currentLog.deviceName}</span>
              <Tag className="detail-tag" color={statusColors[currentLog.status]}>
                {currentLog.status === 'success' ? '成功' : '失败'}
              </Tag>
            </div>
            <div className="detail-content">
              <div className="detail-grid">
              <div className="detail-item">
                <div className="detail-item-label">设备编码</div>
                <div className="detail-item-value">{currentLog.code}</div>
              </div>
              <div className="detail-item">
                <div className="detail-item-label">设备名称</div>
                <div className="detail-item-value">{currentLog.deviceName}</div>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">ID</div>
              <div className="detail-value">{currentLog.id}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">房间名称</div>
              <div className="detail-value">{currentLog.roomName || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">类型</div>
              <div className="detail-value">{currentLog.type === 'lock' ? '门锁' : '设备'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">控制</div>
              <div className="detail-value">
                <Tag color={controlColors[currentLog.control]}>
                  {currentLog.control === 'open' ? '开锁' : '关锁'}
                </Tag>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">状态</div>
              <div className="detail-value">
                <Tag color={statusColors[currentLog.status]}>
                  {currentLog.status === 'success' ? '成功' : '失败'}
                </Tag>
              </div>
            </div>
            <div className="detail-row">
              <div className="detail-label">占用者</div>
              <div className="detail-value">{currentLog.occupant || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">电话</div>
              <div className="detail-value">{currentLog.phone || '-'}</div>
            </div>
            <div className="detail-row">
              <div className="detail-label">创建时间</div>
              <div className="detail-value">{currentLog.createdAt}</div>
            </div>
            </div>
          </div>
        )}
        </Spin>
      </Modal>
    </div>
  );
};