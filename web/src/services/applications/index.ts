// 应用相关的API服务
import { request } from '@/services/request';

/**
 * 获取应用列表
 */
export async function getApplications(
  params: Application.ListApplicationsParams,
): Promise<Application.ListApplicationsResponse['data']> {
  return request('/v1/applications', { params });
}

/**
 * 创建应用
 */
export async function createApplication(
  params: Application.CreateApplicationParams,
): Promise<Application.CreateApplicationResponse['data']> {
  return request('POST /v1/applications', {
    data: params,
  });
}

/**
 * 更新应用
 */
export async function updateApplication(
  params: Application.UpdateApplicationParams,
): Promise<Application.UpdateApplicationResponse['data']> {
  const { id, ...data } = params;
  return request('PUT /v1/applications/:id', {
    params: { id },
    data,
  });
}

/**
 * 删除应用
 */
export async function deleteApplication(
  id: number,
): Promise<Application.DeleteApplicationResponse['message']> {
  return request('DELETE /v1/applications/:id', { params: { id } });
}

/**
 * 获取应用详情
 */
export async function getApplication(
  id: number,
): Promise<Application.CreateApplicationResponse['data']> {
  return request('/v1/applications/:id', { params: { id } });
}
