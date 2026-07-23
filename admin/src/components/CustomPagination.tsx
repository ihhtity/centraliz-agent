import { Button, Select, Input } from 'antd';
import { LeftOutlined, RightOutlined } from '@ant-design/icons';
import { useState, useEffect } from 'react';

interface CustomPaginationProps {
  total: number;
  current: number;
  pageSize: number;
  pageSizeOptions?: string[];
  onChange: (page: number, pageSize: number) => void;
}

export const CustomPagination = ({
  total,
  current,
  pageSize,
  pageSizeOptions = ['10', '20', '50', '100', '200'],
  onChange,
}: CustomPaginationProps) => {
  const [jumpPage, setJumpPage] = useState<string>(current.toString());
  const totalPages = Math.ceil(total / pageSize);
  const start = total > 0 ? (current - 1) * pageSize + 1 : 0;
  const end = Math.min(current * pageSize, total);

  useEffect(() => {
    setJumpPage(current.toString());
  }, [current]);

  const handlePageSizeChange = (value: string) => {
    const newPageSize = parseInt(value, 10);
    const newCurrent = Math.min(current, Math.ceil(total / newPageSize));
    onChange(newCurrent, newPageSize);
  };

  const handlePageChange = (page: number) => {
    if (page >= 1 && page <= totalPages) {
      onChange(page, pageSize);
    }
  };

  const handleJump = () => {
    const page = parseInt(jumpPage, 10);
    if (!isNaN(page) && page >= 1 && page <= totalPages) {
      onChange(page, pageSize);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleJump();
    }
  };

  const getPageNumbers = () => {
    const pages: (number | string)[] = [];
    if (totalPages <= 7) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      if (current <= 3) {
        pages.push(1, 2, 3, 4, '...', totalPages);
      } else if (current >= totalPages - 2) {
        pages.push(1, '...', totalPages - 3, totalPages - 2, totalPages - 1, totalPages);
      } else {
        pages.push(1, '...', current - 1, current, current + 1, '...', totalPages);
      }
    }
    return pages;
  };

  return (
    <div className="custom-pagination">
      <div className="pagination-left">
        <span className="pagination-info">
          {total > 0 ? `显示第 ${start} 到第 ${end} 条记录，总共 ${total} 条记录` : '暂无记录'}
        </span>
        <span className="pagination-page-size">
          每页显示
          <Select
            value={pageSize.toString()}
            onChange={handlePageSizeChange}
            className="page-size-select"
            options={pageSizeOptions.map((size) => ({
              value: size,
              label: size,
            }))}
          />
          条记录
        </span>
      </div>
      <div className="pagination-right">
        <Button
          icon={<LeftOutlined />}
          onClick={() => handlePageChange(current - 1)}
          disabled={current <= 1}
          className="pagination-btn"
        >
          上一页
        </Button>
        <div className="pagination-numbers">
          {getPageNumbers().map((page, index) => (
            <span key={index}>
              {typeof page === 'number' ? (
                <Button
                  onClick={() => handlePageChange(page)}
                  className={`pagination-number ${page === current ? 'active' : ''}`}
                >
                  {page}
                </Button>
              ) : (
                <span className="pagination-ellipsis">...</span>
              )}
            </span>
          ))}
        </div>
        <Button
          icon={<RightOutlined />}
          onClick={() => handlePageChange(current + 1)}
          disabled={current >= totalPages}
          className="pagination-btn"
        >
          下一页
        </Button>
        <Input
          type="number"
          min={1}
          max={totalPages}
          value={jumpPage}
          onChange={(e) => setJumpPage(e.target.value)}
          onKeyPress={handleKeyPress}
          className="pagination-input"
          placeholder="页码"
        />
        <Button onClick={handleJump} className="pagination-jump-btn">
          跳转
        </Button>
      </div>
    </div>
  );
};