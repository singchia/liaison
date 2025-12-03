// 设备相关的API服务
import { request } from '@/services/request';

/**
 * 获取设备列表
 */
export async function getDevices(
  params: Device.ListDevicesParams,
): Promise<Device.ListDevicesResponse['data']> {
  return request('/v1/devices', { params });
}

/**
 * 创建设备
 */
export async function createDevice(
  params: Device.CreateDeviceParams,
): Promise<Device.CreateDeviceResponse['data']> {
  return request('POST /v1/devices', {
    data: params,
  });
}

/**
 * 更新设备
 */
export async function updateDevice(
  params: Device.UpdateDeviceParams,
): Promise<Device.UpdateDeviceResponse['data']> {
  const { id, ...data } = params;
  return request('PUT /v1/devices/:id', {
    params: { id },
    data,
  });
}

/**
 * 删除设备
 */
export async function deleteDevice(
  id: number,
): Promise<Device.DeleteDeviceResponse['message']> {
  return request('DELETE /v1/devices/:id', { params: { id } });
}

/**
 * 获取设备详情
 */
export async function getDevice(
  id: number,
): Promise<Device.CreateDeviceResponse['data']> {
  return request('/v1/devices/:id', { params: { id } });
}
