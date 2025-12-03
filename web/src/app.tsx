import BasicLayout from '@/layouts/BasicLayout';
import HeaderRender from '@/layouts/HeaderRender';
import { getProfile } from '@/services/iam';
import type { Settings as LayoutSettings } from '@ant-design/pro-components';
import { history } from '@umijs/max';
import { message } from 'antd';
import defaultSettings from '../config/defaultSettings';

const loginPath = '/login';

export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  currentUser?: IAM.User;
  loading?: boolean;
  fetchUserInfo?: () => Promise<IAM.User | undefined>;
}> {
  const fetchUserInfo = async () => {
    try {
      const data = await getProfile();
      return data;
    } catch (_error) {
      history.push(loginPath);
    }
  };
  // 如果不是登录页面，执行
  const { location } = history;
  if (![loginPath, '/profile'].includes(location.pathname)) {
    const currentUser = await fetchUserInfo();
    return {
      fetchUserInfo,
      currentUser,
      settings: defaultSettings as Partial<LayoutSettings>,
    };
  }
  return {
    fetchUserInfo,
    settings: defaultSettings as Partial<LayoutSettings>,
  };
}

export const layout = ({ initialState }: any) => {
  return {
    ...defaultSettings,
    title: initialState?.settings?.title || 'App',
    layout: 'mix',
    menu: { locale: false },
    onPageChange: () => {
      // 简单示例：未登录跳转登录页
      const token =
        localStorage.getItem('token') || localStorage.getItem('jwt');
      const { pathname } = window.location;
      if (!token && pathname !== '/login') {
        window.location.href = '/login';
      }
    },
    headerRender: () => (
      <HeaderRender
        logo={initialState?.settings?.logo}
        title={initialState?.settings?.title}
      />
    ),
    childrenRender: (children: any, props: any) => {
      // 登录页不使用 BasicLayout
      if (props.location?.pathname === '/login') {
        return children;
      }
      // 其他页面自动包装 BasicLayout
      return <BasicLayout>{children}</BasicLayout>;
    },
  };
};

export const request = {
  timeout: 20000,
  errorConfig: {
    errorHandler(error: any) {
      const { response, data } = error || {};
      const code = data?.code || response?.status;
      if (code === 401) {
        message.error('登录状态已过期，请重新登录');
        localStorage.removeItem('token');
        localStorage.removeItem('jwt');
        window.location.href = '/login';
        return;
      }
      const msg = data?.message || error?.message || '请求异常';
      message.error(msg);
    },
  },
  requestInterceptors: [
    (config: any) => {
      const token =
        localStorage.getItem('token') || localStorage.getItem('jwt');
      if (token) {
        config.headers = {
          ...(config.headers || {}),
          Authorization: `Bearer ${token}`,
        } as any;
      }
      return config;
    },
  ],
  responseInterceptors: [(response: any) => response],
};
