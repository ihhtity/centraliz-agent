import { Layout, Menu, Button } from 'antd';
import {
  HomeOutlined,
  AppstoreOutlined,
  BoxPlotOutlined,
  TeamOutlined,
  FileTextOutlined,
  ShoppingCartOutlined,
  UserOutlined,
  ClockCircleOutlined,
  LogoutOutlined,
  MenuOutlined,
} from '@ant-design/icons';
import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import './Layout.scss';

const { Header, Sider, Content } = Layout;

const menuItems = [
  { key: '/dashboard', icon: <HomeOutlined />, label: '仪表盘' },
  { key: '/room', icon: <AppstoreOutlined />, label: '包间类型' },
  { key: '/device', icon: <BoxPlotOutlined />, label: '包间管理' },
  { key: '/order', icon: <ShoppingCartOutlined />, label: '订单管理' },
  { key: '/group', icon: <TeamOutlined />, label: '会员管理' },
  { key: '/merch', icon: <UserOutlined />, label: '用户管理' },
  { key: '/rule', icon: <FileTextOutlined />, label: '规则管理' },
  { key: '/devicelog', icon: <ClockCircleOutlined />, label: '设备日志' },
];

export const AdminLayout = ({ children }: { children: React.ReactNode }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const [collapsed, setCollapsed] = useState(false);

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <div className="layout-container">
      <Header className="layout-header">
        <div className="layout-header-left">
          <span className="layout-logo">棋牌室管理</span>
        </div>
        <div className="layout-header-right">
          <Button icon={<LogoutOutlined />} onClick={handleLogout}>退出登录</Button>
        </div>
      </Header>
      <div className="layout-content">
        <Sider className="layout-sider">
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
            theme="dark"
          />
        </Sider>
        <Content className="layout-main">{children}</Content>
      </div>
    </div>
  );
};