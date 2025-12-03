import { logout } from '@/services/iam';
import { UserOutlined } from '@ant-design/icons';
import { Avatar, Dropdown, MenuProps, Space, Typography, message } from 'antd';
import React, { useMemo } from 'react';

const RightContent: React.FC = () => {
  const username = useMemo(() => {
    return (
      localStorage.getItem('username') ||
      JSON.parse(localStorage.getItem('currentUser') || '{}')?.username ||
      '用户名'
    );
  }, []);

  const onLogout = async () => {
    try {
      await logout();
      localStorage.removeItem('token');
      localStorage.removeItem('jwt');
      localStorage.removeItem('username');
      localStorage.removeItem('currentUser');
      message.success('退出成功');
      window.location.href = '/login';
    } catch (error) {
      // 即使接口失败也清除本地数据
      localStorage.removeItem('token');
      localStorage.removeItem('jwt');
      localStorage.removeItem('username');
      localStorage.removeItem('currentUser');
      window.location.href = '/login';
    }
  };

  const items: MenuProps['items'] = [
    {
      key: 'profile',
      label: '个人中心',
      onClick: () => {
        window.location.href = '/profile';
      },
    },
    { type: 'divider' },
    { key: 'logout', label: '退出登录', onClick: onLogout },
  ];

  return (
    <Dropdown menu={{ items }} placement="bottomRight">
      <Space style={{ cursor: 'pointer' }}>
        <Avatar size={28} icon={<UserOutlined />} />
        <Typography.Text>{username}</Typography.Text>
      </Space>
    </Dropdown>
  );
};

export default RightContent;
