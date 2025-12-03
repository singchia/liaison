// 代理相关的API服务
import { request } from '@/services/request';

/**
 * 获取代理列表
 */
export async function getProxies(
  params: Proxy.ListProxiesParams,
): Promise<Proxy.ListProxiesResponse['data']> {
  return request('/v1/proxies', { params });
}

/**
 * 创建代理
 */
export async function createProxy(
  params: Proxy.CreateProxyParams,
): Promise<Proxy.CreateProxyResponse['data']> {
  return request('POST /v1/proxies', {
    data: params,
  });
}

/**
 * 更新代理
 */
export async function updateProxy(
  params: Proxy.UpdateProxyParams,
): Promise<Proxy.UpdateProxyResponse['data']> {
  const { id, ...data } = params;
  return request('PUT /v1/proxies/:id', {
    params: { id },
    data,
  });
}

/**
 * 删除代理
 */
export async function deleteProxy(
  id: number,
): Promise<Proxy.DeleteProxyResponse['message']> {
  return request('DELETE /v1/proxies/:id', { params: { id } });
}

/**
 * 获取代理详情
 */
export async function getProxy(
  id: number,
): Promise<Proxy.CreateProxyResponse['data']> {
  return request('/v1/proxies/:id', { params: { id } });
}
