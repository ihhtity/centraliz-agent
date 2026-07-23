import React from 'react';
import { Button, Dropdown, MenuProps } from 'antd';
import { ExportOutlined } from '@ant-design/icons';
import { exportData, ExportFormat } from '../utils/export';

interface ExportButtonProps {
  data: any[];
  filename: string;
  headers?: Record<string, string>;
  disabled?: boolean;
}

const formatOptions: { value: ExportFormat; label: string }[] = [
  { value: 'json', label: 'JSON' },
  { value: 'xml', label: 'XML' },
  { value: 'csv', label: 'CSV' },
  { value: 'txt', label: 'TXT' },
  { value: 'word', label: 'MS-WORD' },
  { value: 'excel', label: 'MS-Excel' },
];

export const ExportButton: React.FC<ExportButtonProps> = ({ 
  data, 
  filename, 
  headers,
  disabled = false 
}) => {
  const menuItems: MenuProps['items'] = formatOptions.map(option => ({
    key: option.value,
    label: option.label,
  }));

  const handleExport: MenuProps['onClick'] = ({ key }) => {
    exportData({
      format: key as ExportFormat,
      data,
      filename,
      headers,
    });
  };

  return (
    <Dropdown 
      menu={{ items: menuItems, onClick: handleExport }}
      disabled={disabled || data.length === 0}
    >
      <Button className="action-btn-export" size="small" icon={<ExportOutlined />}>
        导出
      </Button>
    </Dropdown>
  );
};

export default ExportButton;