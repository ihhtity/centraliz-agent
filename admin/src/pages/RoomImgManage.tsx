import { Table, Button, Modal, Form, Input, message, Space, Spin, Image } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined, SearchOutlined, UploadOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';
import type { RoomImage } from '@/types';
import {
  getRoomImageList,
  getRoomImageDetail,
  createRoomImage,
  updateRoomImage,
  deleteRoomImage,
  batchDeleteRoomImage,
  batchUpdateRoomImage,
} from '@/api';
import { CustomPagination } from '@/components/CustomPagination';
import { ExportButton } from '@/components/ExportButton';

export const RoomImgManage = () => {
  const [data, setData] = useState<RoomImage[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [showSearch, setShowSearch] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const [currentItem, setCurrentItem] = useState<RoomImage | null>(null);
  const [detailVisible, setDetailVisible] = useState(false);
  const [batchEditVisible, setBatchEditVisible] = useState(false);
  const [form] = Form.useForm();
  const [searchForm] = Form.useForm();
  const [batchEditForm] = Form.useForm();

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '图片名称', dataIndex: 'name', key: 'name' },
    {
      title: '图片预览',
      dataIndex: 'image',
      key: 'image',
      render: (image: string) => (
        <Image src={image} width={80} height={60} />
      ),
    },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (time: string) => formatTime(time) },
    { title: '更新时间', dataIndex: 'updatedAt', key: 'updatedAt', render: (time: string) => formatTime(time) },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: RoomImage) => (
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
      const res = await getRoomImageList({
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

  const viewDetail = async (item: RoomImage) => {
    setDetailLoading(true);
    try {
      const res = await getRoomImageDetail(item.id);
      setCurrentItem(res.data);
      setDetailVisible(true);
    } catch (error) {
      console.error(error);
      message.error('获取详情失败');
    } finally {
      setDetailLoading(false);
    }
  };

  const editItem = (item: RoomImage) => {
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

  const saveItem = async () => {
    try {
      const values = await form.validateFields();
      if (isEdit && currentItem) {
        await updateRoomImage(currentItem.id, values);
        message.success('更新成功');
      } else {
        await createRoomImage(values);
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
      content: '确定要删除这个房间图片吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await deleteRoomImage(id);
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
      message.warning('请选择要删除的房间图片');
      return;
    }
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.length} 张图片吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await batchDeleteRoomImage({ ids: selectedRowKeys.map(String) });
          message.success('批量删除成功');
          setSelectedRowKeys([]);
          loadData();
        } catch (error) {
          message.error('批量删除失败');
        }
      },
    });
  };

  const handleBatchEdit = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要编辑的房间图片');
      return;
    }
    batchEditForm.resetFields();
    setBatchEditVisible(true);
  };

  const saveBatchEdit = async () => {
    try {
      const values = await batchEditForm.validateFields();
      await batchUpdateRoomImage({ ids: selectedRowKeys.map(String), data: values });
      message.success('批量更新成功');
      setBatchEditVisible(false);
      setSelectedRowKeys([]);
      loadData();
    } catch (error) {
      console.error(error);
    }
  };

  const exportHeaders: Record<string, string> = {
    id: 'ID',
    name: '图片名称',
    image: '图片地址',
    createdAt: '创建时间',
    updatedAt: '更新时间',
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
        const response = await fetch('/admin/roomimg/import', {
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
        <h2>房间图片管理</h2>
        <div className="action-buttons">
          <Button className="action-btn-search" size="small" icon={<SearchOutlined />} onClick={() => setShowSearch(!showSearch)}>{showSearch ? '收起搜索' : '搜索'}</Button>
          <ExportButton data={data} filename="房间图片列表" headers={exportHeaders} />
          <Button className="action-btn-import" size="small" icon={<UploadOutlined />} onClick={handleImport}>导入</Button>
          <Button className="action-btn-edit" size="small" icon={<EditOutlined />} onClick={handleBatchEdit} disabled={selectedRowKeys.length === 0}>编辑({selectedRowKeys.length})</Button>
          <Button className="action-btn-delete" size="small" icon={<DeleteOutlined />} onClick={handleBatchDelete} disabled={selectedRowKeys.length === 0}>删除({selectedRowKeys.length})</Button>
          <Button className="action-btn-add" size="small" icon={<PlusOutlined />} onClick={addItem}>添加</Button>
        </div>
      </div>

      {showSearch && (
        <Form className="search-form" form={searchForm} layout="inline">
          <Form.Item name="name">
            <Input placeholder="图片名称" prefix={<SearchOutlined />} />
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
        title={isEdit ? '编辑房间图片' : '添加房间图片'}
        open={modalVisible}
        onOk={saveItem}
        onCancel={() => setModalVisible(false)}
        okText="确定"
        cancelText="取消"
        className="form-modal"
        width="80%"
        style={{ top: '5%' }}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="图片名称" rules={[{ required: true, message: '请输入图片名称' }, { max: 255, message: '图片名称长度不超过255' }]}>
            <Input placeholder="请输入图片名称" />
          </Form.Item>
          <Form.Item name="image" label="图片地址">
            <Input placeholder="请输入图片URL地址" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="批量编辑房间图片"
        open={batchEditVisible}
        onOk={saveBatchEdit}
        onCancel={() => setBatchEditVisible(false)}
        okText="确定"
        cancelText="取消"
        className="form-modal"
      >
        <Form form={batchEditForm} layout="vertical">
          <Form.Item name="name" label="图片名称">
            <Input placeholder="请输入图片名称" />
          </Form.Item>
          <Form.Item name="image" label="图片地址">
            <Input placeholder="请输入图片URL地址" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal title="房间图片详情" open={detailVisible} onCancel={() => setDetailVisible(false)} okText="确定" cancelText="取消" footer={null} className="detail-modal">
        <Spin spinning={detailLoading}>
          {currentItem && (
            <div>
              <div className="detail-header">
                <span className="detail-title">{currentItem.name}</span>
              </div>
              <div className="detail-content">
                <div className="detail-grid">
                  <div className="detail-item">
                    <div className="detail-item-label">图片名称</div>
                    <div className="detail-item-value">{currentItem.name}</div>
                  </div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">ID</div>
                  <div className="detail-value">{currentItem.id}</div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">图片预览</div>
                  <div className="detail-value">
                    <Image src={currentItem.image} width={300} />
                  </div>
                </div>
                <div className="detail-row">
                  <div className="detail-label">图片地址</div>
                  <div className="detail-value">{currentItem.image}</div>
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