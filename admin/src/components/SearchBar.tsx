import { Input, Select, Button, Form, Space } from 'antd';
import React from 'react';

interface SearchBarProps {
  onSearch: (values: Record<string, any>) => void;
  filters: {
    type: 'input' | 'select';
    key: string;
    label: string;
    options?: { label: string; value: string }[];
    placeholder?: string;
  }[];
}

export const SearchBar: React.FC<SearchBarProps> = ({ onSearch, filters }) => {
  const [form] = Form.useForm();

  const handleSearch = () => {
    const values = form.getFieldsValue();
    onSearch(values);
  };

  const handleReset = () => {
    form.resetFields();
    onSearch({});
  };

  return (
    <Form form={form} layout="inline" className="search-bar">
      <Space>
        {filters.map((filter) => {
          if (filter.type === 'input') {
            return (
              <Form.Item key={filter.key} name={filter.key} label={filter.label}>
                <Input placeholder={filter.placeholder} style={{ width: 180 }} />
              </Form.Item>
            );
          }
          return (
            <Form.Item key={filter.key} name={filter.key} label={filter.label}>
              <Select
                placeholder={filter.placeholder}
                options={filter.options}
                allowClear
                style={{ width: 180 }}
              />
            </Form.Item>
          );
        })}
        <Button type="primary" onClick={handleSearch}>
          搜索
        </Button>
        <Button onClick={handleReset}>重置</Button>
      </Space>
    </Form>
  );
};