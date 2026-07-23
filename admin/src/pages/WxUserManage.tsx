import { Table, Button, Modal, Form, Input, Select, message, Tag, Space, Spin, Image, Row, Col } from 'antd';
import { DeleteOutlined, EditOutlined, EyeOutlined, SearchOutlined, PlusOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import type { WechatUser } from '@/types';
import {
  getWxUserList,
  getWxUserDetail,
  updateWxUser,
  batchDeleteWxUser,
  importWxUser,
  createWxUser,
  batchUpdateWxUser,
} from '@/api';
import { CustomPagination } from '@/components/CustomPagination';
import { ExportButton } from '@/components/ExportButton';

const statusColors: Record<string, string> = {
  '0': 'green',
  '1': 'red',
};

const platformOptions: Record<string, string> = {
  'miniprogram': '小程序',
  'mp': '公众号',
};

export const WxUserManage = () => {
  const [data, setData] = useState<WechatUser[]>([]);
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
  const [currentItem, setCurrentItem] = useState<WechatUser | null>(null);
  const [detailVisible, setDetailVisible] = useState(false);
  const [form] = Form.useForm();
  const [batchForm] = Form.useForm();
  const [searchForm] = Form.useForm();

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '昵称', dataIndex: 'nickname', key: 'nickname' },
    { title: '微信OpenID', dataIndex: 'openId', key: 'openId' },
    { title: '微信UnionID', dataIndex: 'unionId', key: 'unionId' },
    {
      title: '平台类型',
      dataIndex: 'platform',
      key: 'platform',
      render: (platform: string) => (
        <Tag color={platform === 'miniprogram' ? 'blue' : 'purple'}>{platformOptions[platform] || platform}</Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={statusColors[status] || 'default'}>{status === '0' ? '正常' : '禁用'}</Tag>
      ),
    },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (val: string) => formatTime(val) },
    { title: '更新时间', dataIndex: 'updatedAt', key: 'updatedAt', render: (val: string) => formatTime(val) },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: WechatUser) => (
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
      const res = await getWxUserList({
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

  const viewDetail = async (item: WechatUser) => {
    setDetailLoading(true);
    try {
      const res = await getWxUserDetail(item.id);
      setCurrentItem(res.data);
      setDetailVisible(true);
    } catch (error) {
      console.error(error);
      message.error('获取详情失败');
    } finally {
      setDetailLoading(false);
    }
  };

  const handleAdd = () => {
    setIsEdit(false);
    setCurrentItem(null);
    form.resetFields();
    setModalVisible(true);
  };

  const editItem = (item: WechatUser) => {
    setIsEdit(true);
    setCurrentItem(item);
    form.setFieldsValue(item);
    setModalVisible(true);
  };

  const saveItem = async () => {
    try {
      const values = await form.validateFields();
      if (isEdit && currentItem) {
        await updateWxUser(currentItem.id, values);
        message.success('更新成功');
      } else {
        await createWxUser(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      loadData();
    } catch (error) {
      console.error(error);
    }
  };

  const handleBatchEdit = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要编辑的微信用户');
      return;
    }
    batchForm.resetFields();
    setIsBatchModalVisible(true);
  };

  const handleBatchSubmit = async () => {
    try {
      const values = await batchForm.validateFields();
      await batchUpdateWxUser({ ids: selectedRowKeys.map(k => k.toString()), data: values });
      message.success('批量更新成功');
      setIsBatchModalVisible(false);
      setSelectedRowKeys([]);
      loadData();
    } catch (error) {
      message.error('批量更新失败');
    }
  };

  const deleteItem = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个微信用户吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteWxUser({ ids: [id.toString()] });
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
      message.warning('请选择要删除的微信用户');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 个微信用户吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteWxUser({ ids: selectedRowKeys.map(String) });
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
      await importWxUser(formData);
      message.success('导入成功');
      loadData();
    } catch (error) {
      message.error('导入失败');
    }
  };

  const exportHeaders: Record<string, string> = {
    id: 'ID',
    nickname: '昵称',
    avatar: '头像',
    openId: '微信OpenID',
    gopenId: '微信公众号OpenID',
    unionId: '微信UnionID',
    platform: '平台类型',
    gender: '性别',
    country: '国家',
    province: '省份',
    city: '城市',
    status: '状态',
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
        <h2>微信用户管理</h2>
        <div className="action-buttons">
          <Button className="action-btn-search" size="small" icon={<SearchOutlined />} onClick={() => setShowSearch(!showSearch)}>{showSearch ? '收起搜索' : '搜索'}</Button>
          <ExportButton data={data} filename="微信用户列表" headers={exportHeaders} />
          <Button className="action-btn-import" size="small" onClick={() => document.getElementById('wx-user-import')?.click()}>导入</Button>
          <input type="file" id="wx-user-import" style={{ display: 'none' }} onChange={(e) => e.target.files?.[0] && handleImport(e.target.files[0])} />
          <Button className="action-btn-edit" size="small" icon={<EditOutlined />} onClick={handleBatchEdit} disabled={selectedRowKeys.length === 0}>编辑({selectedRowKeys.length})</Button>
          <Button className="action-btn-delete" size="small" icon={<DeleteOutlined />} onClick={handleBatchDelete} disabled={selectedRowKeys.length === 0}>删除({selectedRowKeys.length})</Button>
          <Button className="action-btn-add" size="small" icon={<PlusOutlined />} onClick={handleAdd}>添加</Button>
        </div>
      </div>

      {showSearch && (
        <Form className="search-form" form={searchForm} layout="inline">
          <Form.Item name="nickname">
            <Input placeholder="微信昵称" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="openid">
            <Input placeholder="微信OpenID" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="unionid">
            <Input placeholder="微信UnionID" prefix={<SearchOutlined />} />
          </Form.Item>
          <Form.Item name="platform">
            <Select placeholder="平台类型" allowClear>
              <Select.Option value="miniprogram">小程序</Select.Option>
              <Select.Option value="mp">公众号</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="status">
            <Select placeholder="状态" allowClear>
              <Select.Option value="0">正常</Select.Option>
              <Select.Option value="1">禁用</Select.Option>
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
            onSelect: (record: WechatUser, selected: boolean) => {
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
        title={isEdit ? '编辑微信用户' : '添加微信用户'}
        open={modalVisible}
        onOk={saveItem}
        onCancel={() => setModalVisible(false)}
        okText="确定"
        cancelText="取消"
        width={600}
        className="form-modal"
      >
        <Form form={form} layout="vertical">
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="nickname" label="微信昵称">
                <Input placeholder="请输入微信昵称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="openId" label="微信OpenID">
                <Input placeholder="请输入微信OpenID" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="gopenId" label="公众号OpenID">
                <Input placeholder="请输入公众号OpenID" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="unionId" label="微信UnionID">
                <Input placeholder="请输入微信UnionID" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="platform" label="平台类型">
                <Select options={[{ label: '小程序', value: 'miniprogram' }, { label: '公众号', value: 'mp' }]} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="status" label="状态">
                <Select options={[{ label: '正常', value: '0' }, { label: '禁用', value: '1' }]} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="gender" label="性别">
                <Select options={[{ label: '未知', value: 0 }, { label: '男', value: 1 }, { label: '女', value: 2 }]} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="avatar" label="头像">
                <Input placeholder="请输入头像URL" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={8}>
              <Form.Item name="country" label="国家">
                <Input placeholder="请输入国家" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="province" label="省份">
                <Input placeholder="请输入省份" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="city" label="城市">
                <Input placeholder="请输入城市" />
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
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="platform" label="平台类型">
                <Select placeholder="请选择平台类型" allowClear>
                  <Select.Option value="miniprogram">小程序</Select.Option>
                  <Select.Option value="mp">公众号</Select.Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="status" label="状态">
                <Select placeholder="请选择状态" allowClear>
                  <Select.Option value="0">正常</Select.Option>
                  <Select.Option value="1">禁用</Select.Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>

      <Modal title="微信用户详情" open={detailVisible} onCancel={() => setDetailVisible(false)} okText="确定" cancelText="取消" footer={null} className="detail-modal">
        <Spin spinning={detailLoading}>
          {currentItem && (
            <div>
              <div className="detail-header">
                <span className="detail-title">{currentItem.nickname}</span>
                <Tag className="detail-tag" color={currentItem.status === '0' ? 'green' : 'red'}>
                  {currentItem.status === '0' ? '正常' : '禁用'}
                </Tag>
              </div>
              <div className="detail-content">
                <div className="detail-grid">
                  <div className="detail-item">
                    <div className="detail-item-label">微信昵称</div>
                    <div className="detail-item-value">{currentItem.nickname}</div>
                  </div>
                  <div className="detail-item">
                    <div className="detail-item-label">平台类型</div>
                    <div className="detail-item-value">
                      <Tag color={currentItem.platform === 'miniprogram' ? 'blue' : 'purple'}>
                        {platformOptions[currentItem.platform] || currentItem.platform}
                      </Tag>
                    </div>
                  </div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">ID</div>
                  <div className="detail-value">{currentItem.id}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">头像</div>
                  <div className="detail-value">
                    {currentItem.avatar ? <Image src={currentItem.avatar} width={80} /> : '-'}
                  </div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">微信OpenID</div>
                  <div className="detail-value">{currentItem.openId}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">微信公众号GopenID</div>
                  <div className="detail-value">{currentItem.gopenId || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">微信UnionID</div>
                  <div className="detail-value">{currentItem.unionId || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">性别</div>
                  <div className="detail-value">{currentItem.gender === 1 ? '男' : currentItem.gender === 2 ? '女' : '未知'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">国家</div>
                  <div className="detail-value">{currentItem.country || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">省份</div>
                  <div className="detail-value">{currentItem.province || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">城市</div>
                  <div className="detail-value">{currentItem.city || '-'}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">状态</div>
                  <div className="detail-value">
                    <Tag color={currentItem.status === '0' ? 'green' : 'red'}>
                      {currentItem.status === '0' ? '正常' : '禁用'}
                    </Tag>
                  </div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">创建时间</div>
                  <div className="detail-value">{currentItem.createdAt}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">更新时间</div>
                  <div className="detail-value">{currentItem.updatedAt}</div>
                </div>
              </div>
            </div>
          )}
        </Spin>
      </Modal>
    </div>
  );
};