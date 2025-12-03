// 代理相关的类型定义

declare namespace Proxy {
  interface Proxy {
    id: number;
    name: string;
    description: string;
    port: number;
    application: Application.Application;
    created_at: string;
    updated_at: string;
  }

  interface Proxies {
    proxies: Proxy[];
    total: number;
  }

  interface ListProxiesResponse {
    code: number;
    message: string;
    data: Proxies;
  }

  interface CreateProxyResponse {
    code: number;
    message: string;
    data?: Proxy;
  }

  interface UpdateProxyResponse {
    code: number;
    message: string;
    data?: Proxy;
  }

  interface DeleteProxyResponse {
    code: number;
    message: string;
  }

  // 请求参数类型
  interface ListProxiesParams {
    page?: number;
    page_size?: number;
    name?: string;
    application_id?: number;
  }

  interface CreateProxyParams {
    name: string;
    description: string;
    port: number;
    application_id: number;
  }

  interface UpdateProxyParams {
    id: number;
    name?: string;
    description?: string;
    port?: number;
    application_id?: number;
  }
}
