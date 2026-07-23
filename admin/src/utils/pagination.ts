export const getPaginationConfig = (total: number, current: number, pageSize: number) => {
  return {
    current,
    pageSize,
    total,
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (total: number, range: [number, number]) => {
      return `显示第 ${range[0]} 到第 ${range[1]} 条记录，总共 ${total} 条记录`;
    },
    pageSizeOptions: ['10', '20', '50', '100', '200'],
    showSizeChange: true,
    size: 'default' as const,
  };
};