// 边缘节点相关的API服务
import { request } from '@/services/request';

/**
 * 获取边缘节点列表
 */
export async function getEdges(
  params: Edge.ListEdgesParams,
): Promise<Edge.ListEdgesResponse['data']> {
  return request('/v1/edges', { params });
}

/**
 * 创建边缘节点
 */
export async function createEdge(
  params: Edge.CreateEdgeParams,
): Promise<Edge.CreateEdgeResponse['data']> {
  return request('POST /v1/edges', {
    data: params,
  });
}

/**
 * 更新边缘节点
 */
export async function updateEdge(
  params: Edge.UpdateEdgeParams,
): Promise<Edge.UpdateEdgeResponse['data']> {
  const { id, ...data } = params;
  return request('PUT /v1/edges/:id', {
    params: { id },
    data,
  });
}

/**
 * 删除边缘节点
 */
export async function deleteEdge(
  id: number,
): Promise<Edge.DeleteEdgeResponse['message']> {
  return request('DELETE /v1/edges/:id', { params: { id } });
}

/**
 * 获取边缘节点详情
 */
export async function getEdge(
  id: number,
): Promise<Edge.CreateEdgeResponse['data']> {
  return request('/v1/edges/:id', { params: { id } });
}
