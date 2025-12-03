import { history, request as umiRequest } from '@umijs/max';
import { message } from 'antd';

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

export interface RequestOptions {
  method?: HttpMethod;
  params?: Record<string, any>;
  data?: any;
  headers?: Record<string, string>;
}

function parseUrlAndMethod(
  input: string,
  method?: HttpMethod,
): { url: string; method: HttpMethod } {
  const trimmed = input.trim();
  const match = /^(GET|POST|PUT|DELETE|PATCH)\s+(.+)/i.exec(trimmed);
  if (match) {
    return { url: match[2], method: match[1].toUpperCase() as HttpMethod };
  }
  return { url: trimmed, method: (method || 'GET') as HttpMethod };
}

function getAuthToken(): string | null {
  // 约定本地存储token键名，可按需调整
  return localStorage.getItem('token') || localStorage.getItem('jwt');
}

function interpolateUrl(
  url: string,
  params: Record<string, any> | undefined,
): { url: string; restParams: Record<string, any> | undefined } {
  if (!params) return { url, restParams: params };
  let newUrl = url;
  const rest: Record<string, any> = { ...params };
  // 将 /:id /:name 替换为 params 中对应的值
  newUrl = newUrl.replace(/:(\w+)/g, (_m, key) => {
    if (rest[key] !== undefined && rest[key] !== null) {
      const v = encodeURIComponent(rest[key]);
      delete rest[key];
      return v;
    }
    return _m; // 未提供则保持原样，交由后端或调用方处理
  });
  return {
    url: newUrl,
    restParams: Object.keys(rest).length ? rest : undefined,
  };
}

export async function request<T = any>(
  input: string,
  options: RequestOptions = {},
): Promise<T> {
  const { url, method } = parseUrlAndMethod(input, options.method);

  // 支持 /v1/devices/:id 形式
  const { url: finalUrl, restParams } = interpolateUrl(url, options.params);

  const token = getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  };
  if (token) headers.Authorization = `Bearer ${token}`;

  try {
    const resp = await umiRequest(finalUrl, {
      method,
      params: restParams,
      data: options.data,
      headers,
    });

    // 约定后端返回 { code, data, message }
    const code: number | undefined = resp?.code;

    if (code === 200 || code === 0) {
      // 兼容部分接口返回0为成功
      return resp?.data as T;
    }

    if (code === 401) {
      message.error('登录状态已过期，请重新登录');
      localStorage.removeItem('token');
      localStorage.removeItem('jwt');
      // 跳转登录
      try {
        history.replace('/login');
      } catch (_) {
        window.location.href = '/login';
      }
      throw new Error('Unauthorized');
    }

    const errMsg = resp?.message || '请求失败';
    message.error(errMsg);
    throw new Error(errMsg);
  } catch (error: any) {
    if (error?.response?.status === 401) {
      message.error('登录状态已过期，请重新登录');
      localStorage.removeItem('token');
      localStorage.removeItem('jwt');
      try {
        history.replace('/login');
      } catch (_) {
        window.location.href = '/login';
      }
    } else if (error?.message) {
      message.error(error.message);
    } else {
      message.error('网络异常，请稍后重试');
    }
    throw error;
  }
}
