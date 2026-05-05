import React from 'react';
import { Layout as AntLayout, Menu, Dropdown, Avatar, Button, message } from 'antd';
import { 
  DashboardOutlined, 
  NotificationOutlined,
  HomeOutlined,
  TeamOutlined,
  UserOutlined,
  ToolOutlined,
  SettingOutlined,
  LogoutOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined
} from '@ant-design/icons';
import { useState, useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { logout, getCurrentUser, getRole } from '../utils/auth';

const { Header, Sider, Content } = AntLayout;

const Layout = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [user, setUser] = useState(null);
  const [role, setRole] = useState(null);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    setUser(getCurrentUser());
    setRole(getRole());
  }, []);

  const handleLogout = () => {
    logout();
    message.success('已退出登录');
    navigate('/login');
  };

  const getMenuItems = () => {
    const baseItems = [
      {
        key: '/profile',
        icon: <UserOutlined />,
        label: '个人中心',
      },
    ];

    if (role === 'admin') {
      return [
        {
          key: '/dashboard',
          icon: <DashboardOutlined />,
          label: '仪表盘',
        },
        {
          key: '/notices',
          icon: <NotificationOutlined />,
          label: '公告信息管理',
        },
        {
          key: '/dormitories',
          icon: <HomeOutlined />,
          label: '寝室管理',
        },
        {
          key: '/assignments',
          icon: <SettingOutlined />,
          label: '寝室安排',
        },
        {
          key: '/students',
          icon: <TeamOutlined />,
          label: '学生信息管理',
        },
        {
          key: '/workers',
          icon: <ToolOutlined />,
          label: '维修工管理',
        },
        {
          key: '/repair-requests',
          icon: <ToolOutlined />,
          label: '寝室报修管理',
        },
        {
          key: '/repair-records',
          icon: <SettingOutlined />,
          label: '维修情况管理',
        },
        {
          key: '/admins',
          icon: <UserOutlined />,
          label: '管理员管理',
        },
        ...baseItems,
      ];
    } else if (role === 'student') {
      return [
        {
          key: '/notices',
          icon: <NotificationOutlined />,
          label: '公告查询',
        },
        {
          key: '/my-assignment',
          icon: <HomeOutlined />,
          label: '寝室安排',
        },
        {
          key: '/my-repair-requests',
          icon: <ToolOutlined />,
          label: '寝室报修管理',
        },
        {
          key: '/my-repair-records',
          icon: <SettingOutlined />,
          label: '维修情况',
        },
        ...baseItems,
      ];
    } else if (role === 'worker') {
      return [
        {
          key: '/notices',
          icon: <NotificationOutlined />,
          label: '公告查询',
        },
        {
          key: '/my-repair-requests',
          icon: <ToolOutlined />,
          label: '寝室报修管理',
        },
        {
          key: '/my-repair-records',
          icon: <SettingOutlined />,
          label: '维修情况',
        },
        ...baseItems,
      ];
    }
    return baseItems;
  };

  const userMenu = (
    <Menu>
      <Menu.Item key="profile" icon={<UserOutlined />}>
        个人中心
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item key="logout" icon={<LogoutOutlined />} onClick={handleLogout}>
        退出登录
      </Menu.Item>
    </Menu>
  );

  const getRoleLabel = () => {
    switch (role) {
      case 'admin': return '管理员';
      case 'student': return '学生';
      case 'worker': return '维修工';
      default: return '';
    }
  };

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Sider trigger={null} collapsible collapsed={collapsed}>
        <div style={{ 
          height: 64, 
          display: 'flex', 
          alignItems: 'center', 
          justifyContent: 'center',
          background: 'rgba(255, 255, 255, 0.1)',
          margin: 16,
          borderRadius: 4
        }}>
          <span style={{ 
            color: 'white', 
            fontSize: collapsed ? 12 : 18, 
            fontWeight: 'bold' 
          }}>
            {collapsed ? '寝室' : '寝室管理系统'}
          </span>
        </div>
        
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={getMenuItems()}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      
      <AntLayout>
        <Header style={{ 
          padding: '0 24px', 
          background: '#fff', 
          display: 'flex', 
          alignItems: 'center',
          justifyContent: 'space-between',
          boxShadow: '0 1px 4px rgba(0,21,41,0.08)'
        }}>
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{ fontSize: '16px', width: 64, height: 64 }}
          />
          
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <span style={{ marginRight: 12, color: '#666' }}>
              {getRoleLabel()}
            </span>
            <Dropdown overlay={userMenu} placement="bottomRight">
              <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center' }}>
                <Avatar icon={<UserOutlined />} style={{ marginRight: 8 }} />
                <span>{user?.name || user?.username}</span>
              </div>
            </Dropdown>
          </div>
        </Header>
        
        <Content style={{ 
          margin: '24px', 
          padding: 24, 
          background: '#fff', 
          borderRadius: 4,
          minHeight: 280
        }}>
          <Outlet />
        </Content>
      </AntLayout>
    </AntLayout>
  );
};

export default Layout;
