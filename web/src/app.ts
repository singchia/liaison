import { history, RequestConfig } from '@umijs/max';
import { Dropdown, message, Typography, Tooltip } from 'antd';
import { LogoutOutlined, SettingOutlined, GithubOutlined, BugOutlined } from '@ant-design/icons';
import React from 'react';
import { getCurrentUser, logout } from '@/services/api';
import { APP_NAME, GITHUB_URL } from '@/constants';
import './global.less';

const { Text, Link } = Typography;
const GITHUB_ISSUES_URL = 'https://github.com/singchia/liaison/issues/new';

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
export const layout = ({ initialState }: any) => {
  const dropdownMenuItems = [
    {
      key: 'settings',
      icon: React.createElement(SettingOutlined),
      label: '个人设置',
      onClick: () => history.push('/settings'),
    },
    {
      key: 'logout',
      icon: React.createElement(LogoutOutlined),
      label: '退出登录',
      onClick: handleLogout,
    },
    {
      type: 'divider' as const,
    },
    {
      key: 'report',
      icon: React.createElement(BugOutlined),
      label: '报告问题',
      onClick: () => {
        window.open(GITHUB_ISSUES_URL, '_blank');
      },
    },
  ];

  return {
    logo: React.createElement('img', {
      src: '/liaison.png',
      alt: 'Liaison',
      style: {
        height: 52,
      },
    }),
    menu: {
      locale: false,
    },
    layout: 'mix',
    splitMenus: false,
    fixedHeader: true,
    fixSiderbar: true,
    navTheme: 'light',
    contentWidth: 'Fluid',
    colorPrimary: '#1677ff', // 更现代的蓝色
    siderWidth: 220,
    title: 'Liaison',
    token: {
      // 现代化的设计 token
      colorBgContainer: '#ffffff',
      colorBgElevated: '#fafafa',
      borderRadius: 8,
      wireframe: false,
      // 更好的间距
      sizeStep: 4,
      sizeUnit: 4,
    },
    avatarProps: {
      src: '/avatar.svg',
      size: 'small',
      title: initialState?.currentUser?.name || 'Admin',
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
      fontSize: 14,
      fontColor: 'rgba(0, 0, 0, 0.04)',
      gapX: 120,
      gapY: 120,
      rotate: -22,
      fontStyle: 'normal',
      fontWeight: 'normal',
    },
    footerRender: () => {
      return React.createElement(
        'div',
        {
          className: 'global-footer',
          style: {
            position: 'fixed',
            bottom: 0,
            left: 0,
            right: 0,
            textAlign: 'center',
            padding: '16px 24px',
            backgroundColor: 'transparent',
            zIndex: 100,
            color: 'rgba(0, 0, 0, 0.45)',
          },
        },
        React.createElement(
          'div',
          {
            style: {
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
              gap: 8,
              flexWrap: 'wrap',
            },
          },
          React.createElement(Text, { type: 'secondary', style: { fontSize: 13 } }, `© 2026 ${APP_NAME}. All rights reserved.`),
          React.createElement(
            Link,
            {
              href: GITHUB_URL,
              target: '_blank',
              rel: 'noopener noreferrer',
              style: {
                display: 'inline-flex',
                alignItems: 'center',
                gap: 6,
                color: 'rgba(0, 0, 0, 0.45)',
                fontSize: 13,
                transition: 'color 0.3s',
                textDecoration: 'none',
              },
              onMouseEnter: (e: any) => {
                e.currentTarget.style.color = '#1677ff';
              },
              onMouseLeave: (e: any) => {
                e.currentTarget.style.color = 'rgba(0, 0, 0, 0.45)';
              },
            },
            React.createElement('img', {
              src: '/github.svg',
              alt: 'GitHub',
              style: {
                width: 18,
                height: 18,
                verticalAlign: 'middle',
              },
            }),
            React.createElement('span', null, 'Github')
          ),
          React.createElement(
            Tooltip,
            {
              title: React.createElement('img', {
                src: '/wechat.png',
                alt: '微信二维码',
                style: {
                  width: 200,
                  height: 200,
                  display: 'block',
                },
              }),
              placement: 'top',
              overlayStyle: {
                padding: 0,
              },
              overlayInnerStyle: {
                padding: 8,
                backgroundColor: '#fff',
                borderRadius: 8,
                boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
              },
            },
            React.createElement(
              'div',
              {
                style: {
                  display: 'inline-flex',
                  alignItems: 'center',
                  gap: 6,
                  cursor: 'pointer',
                  color: 'rgba(0, 0, 0, 0.45)',
                  fontSize: 13,
                },
              },
              React.createElement('img', {
                src: '/wechat.svg',
                alt: 'WeChat',
                style: {
                  width: 18,
                  height: 18,
                  verticalAlign: 'middle',
                  opacity: 0.7,
                  transition: 'opacity 0.3s',
                },
                onMouseEnter: (e: any) => {
                  e.currentTarget.style.opacity = '1';
                },
                onMouseLeave: (e: any) => {
                  e.currentTarget.style.opacity = '0.7';
                },
              }),
              React.createElement('span', null, 'WeChat')
            )
          )
        )
      );
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
      if (response?.status === 401 || response?.status === 403) {
        localStorage.removeItem('token');
        message.error('登录已过期，请重新登录');
        history.push('/login');
      } else if (response?.status === 500) {
        // 检查是否是认证相关的500错误
        const errorMessage = error?.response?.data?.message || error?.message || '';
        if (errorMessage.includes('authentication') || errorMessage.includes('token') || errorMessage.includes('unauthorized')) {
          localStorage.removeItem('token');
          message.error('登录已过期，请重新登录');
          history.push('/login');
        } else {
          message.error('服务器错误');
        }
      } else if (!response) {
        message.error('网络异常');
      }
    },
  },
};
