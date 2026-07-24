import { Layout, Menu, Button, Modal } from 'antd';
import {
  HomeOutlined,
  AppstoreOutlined,
  BoxPlotOutlined,
  UnorderedListOutlined,
  FileTextOutlined,
  ShoppingCartOutlined,
  UserOutlined,
  ClockCircleOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  MessageOutlined,
} from '@ant-design/icons';
import { useState, useEffect, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import './Layout.scss';

const { Header, Sider, Content } = Layout;

// 菜单配置
const menuItems = [
  { key: '/dashboard', icon: <HomeOutlined />, label: '仪表盘' },
  { key: '/order', icon: <ShoppingCartOutlined />, label: '订单管理' },
  { key: '/user', icon: <FileTextOutlined />, label: '用户管理' },
  { key: '/merch', icon: <UserOutlined />, label: '商户管理' },
  { key: '/submerch', icon: <UserOutlined />, label: '子商户' },
  { key: '/rule', icon: <FileTextOutlined />, label: '规则管理' },
  { key: '/group', icon: <UnorderedListOutlined />, label: '分组管理' },
  { key: '/room', icon: <AppstoreOutlined />, label: '房间管理' },
  { key: '/device', icon: <BoxPlotOutlined />, label: '设备管理' },
  { key: '/devicelog', icon: <ClockCircleOutlined />, label: '设备日志' },
  { key: '/roomimg', icon: <ShoppingCartOutlined />, label: '房间图片' },
  { key: '/roomtag', icon: <UnorderedListOutlined />, label: '房间标签' },
  { key: '/huifu', icon: <AppstoreOutlined />, label: '汇付账号' },
  { key: '/merchpay', icon: <BoxPlotOutlined />, label: '商户支付' },
  { key: '/wxuser', icon: <ClockCircleOutlined />, label: '微信用户' },
  { key: '/dev-assistant', icon: <MessageOutlined />, label: '开发助手' },
];

export const AdminLayout = ({ children }: { children: React.ReactNode }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const [collapsed, setCollapsed] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  const handleLogout = () => {
    Modal.confirm({
      title: '退出登录',
      content: '确定要退出登录吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: () => {
        localStorage.removeItem('token');
        navigate('/login');
      },
    });
  };

  useEffect(() => {
    const scrollTop = localStorage.getItem('menuScrollTop');
    if (scrollTop && menuRef.current) {
      menuRef.current.scrollTop = parseInt(scrollTop, 10);
    }

    const activeMenuItem = document.querySelector('.ant-menu-item-selected');
    if (activeMenuItem) {
      activeMenuItem.scrollIntoView({ block: 'nearest' });
    }
  }, [location.pathname]);

  const handleMenuScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const target = e.target as HTMLDivElement;
    localStorage.setItem('menuScrollTop', target.scrollTop.toString());
  };

  return (
    <div className="layout-container">
      <Header className="layout-header">
        <div className="layout-header-left">
          <Button 
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />} 
            onClick={() => setCollapsed(!collapsed)}
            className="collapse-btn"
          />
          <span className="layout-logo">{collapsed ? '管' : '数据管理'}</span>
        </div>
        <div className="layout-header-right">
          <Button icon={<LogoutOutlined />} onClick={handleLogout}>退出登录</Button>
        </div>
      </Header>
      <div className="layout-content">
        <Sider className="layout-sider" collapsed={collapsed}>
          <div ref={menuRef} className="menu-scroll-container" onScroll={handleMenuScroll}>
            <Menu
              mode="inline"
              selectedKeys={[location.pathname]}
              items={menuItems}
              onClick={({ key }) => navigate(key)}
              theme="dark"
              inlineCollapsed={collapsed}
            />
          </div>
        </Sider>
        <Content className="layout-main">{children}</Content>
      </div>
    </div>
  );
};