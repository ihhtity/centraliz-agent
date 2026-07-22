import { Button, Space, Popconfirm } from 'antd';
import { DeleteOutlined, DownloadOutlined } from '@ant-design/icons';
import React from 'react';
import * as XLSX from 'xlsx';
import { saveAs } from 'file-saver';

interface BatchActionsProps {
  selectedRowKeys: React.Key[];
  onBatchDelete: () => void;
  data: any[];
  columns: { title: string; key: string }[];
}

export const BatchActions: React.FC<BatchActionsProps> = ({
  selectedRowKeys,
  onBatchDelete,
  data,
  columns,
}) => {
  const handleExport = () => {
    const exportData = data.map((item) => {
      const row: Record<string, any> = {};
      columns.forEach((col) => {
        row[col.title] = item[col.key];
      });
      return row;
    });

    const worksheet = XLSX.utils.json_to_sheet(exportData);
    const workbook = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(workbook, worksheet, '数据');
    const excelBuffer = XLSX.write(workbook, { bookType: 'xlsx', type: 'array' });
    const blob = new Blob([excelBuffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
    saveAs(blob, '导出数据.xlsx');
  };

  if (selectedRowKeys.length === 0) return null;

  return (
    <Space>
      <span>已选择 {selectedRowKeys.length} 条记录</span>
      <Popconfirm
        title={`确定删除选中的 ${selectedRowKeys.length} 条记录吗？`}
        onConfirm={onBatchDelete}
        okText="确定"
        cancelText="取消"
      >
        <Button danger icon={<DeleteOutlined />}>
          批量删除
        </Button>
      </Popconfirm>
      <Button icon={<DownloadOutlined />} onClick={handleExport}>
        导出
      </Button>
    </Space>
  );
};