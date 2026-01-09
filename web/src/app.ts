import { history, RequestConfig } from '@umijs/max';
import { Dropdown, message } from 'antd';
import { LogoutOutlined, SettingOutlined } from '@ant-design/icons';
import React from 'react';
import { getCurrentUser, logout } from '@/services/api';

if (process.env.NODE_ENV === 'development') {
  const filterMessages = [
    'findDOMNode is deprecated',
    'Static function can not consume context',
  ];
  
  const shouldFilter = (args: any[]) => {
    const msg = args[0];
    if (typeof msg === 'string') {
      return filterMessages.some(filter => msg.includes(filter));
    }
    return false;
  };

  const originalWarn = console.warn;
  const originalError = console.error;
  
  console.warn = (...args) => {
    if (shouldFilter(args)) return;
    originalWarn.apply(console, args);
  };
  
  console.error = (...args) => {
    if (shouldFilter(args)) return;
    originalError.apply(console, args);
  };
}

export async function getInitialState(): Promise<{
  currentUser?: API.CurrentUser;
  fetchUserInfo?: () => Promise<API.CurrentUser | undefined>;
}> {
  const fetchUserInfo = async () => {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        return undefined;
      }
      const res = await getCurrentUser();
      if (res.code === 200 && res.data) {
        return res.data;
      }
      return undefined;
    } catch (error) {
      return undefined;
    }
  };

  const { location } = history;
  if (location.pathname !== '/login') {
    const currentUser = await fetchUserInfo();
    if (!currentUser) {
      history.push('/login');
    }
    return {
      fetchUserInfo,
      currentUser,
    };
  }
  return {
    fetchUserInfo,
  };
}

const handleLogout = async () => {
  try {
    await logout();
  } catch (e) {}
  localStorage.removeItem('token');
  message.success('已退出登录');
  history.push('/login');
};

// 布局配置
export const layout = () => {
  const dropdownMenuItems = [
    {
      key: 'settings',
      icon: React.createElement(SettingOutlined),
      label: '个人设置',
      onClick: () => history.push('/settings'),
    },
    {
      type: 'divider' as const,
    },
    {
      key: 'logout',
      icon: React.createElement(LogoutOutlined),
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  return {
    logo: false, // 不显示 logo 图标
    menu: {
      locale: false,
    },
    layout: 'mix',
    splitMenus: false,
    fixedHeader: true,
    fixSiderbar: true,
    navTheme: 'light',
    contentWidth: 'Fluid',
    colorPrimary: '#1890ff',
    siderWidth: 208,
    title: 'Liaison',
    titleRender: () => {
      return React.createElement(
        'div',
        { style: { display: 'flex', alignItems: 'center', color: '#1890ff', fontWeight: 600, fontSize: '18px' } },
        'Liaison'
      );
    },
    avatarProps: {
      src: 'https://gw.alipayobjects.com/zos/antfincdn/efFD%24IOql2/weixintupian_20170331104822.jpg',
      size: 'small',
      title: 'Admin',
      render: (_props: any, avatarChildren: React.ReactNode) => {
        return React.createElement(
          Dropdown,
          { menu: { items: dropdownMenuItems } },
          avatarChildren
        );
      },
    },
    waterMarkProps: {
      content: 'Liaison',
    },
  };
};

// 请求配置
export const request: RequestConfig = {
  timeout: 30000,
  requestInterceptors: [
    (config: any) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers = {
          ...config.headers,
          Authorization: `Bearer ${token}`,
        };
      }
      return config;
    },
  ],
  responseInterceptors: [
    (response: any) => {
      const { data } = response;
      if (data && data.code === 401) {
        localStorage.removeItem('token');
        history.push('/login');
        return response;
      }
      return response;
    },
  ],
  errorConfig: {
    errorHandler: (error: any) => {
      const { response } = error;
      if (response?.status === 401) {
        localStorage.removeItem('token');
        message.error('登录已过期，请重新登录');
        history.push('/login');
      } else if (response?.status === 403) {
        message.error('没有权限访问');
      } else if (response?.status === 500) {
        message.error('服务器错误');
      } else if (!response) {
        message.error('网络异常');
      }
    },
  },
};
