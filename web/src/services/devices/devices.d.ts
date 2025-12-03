declare namespace Device {
  interface Device {
    id: number;
    name: string;
    device_type: string;
    status: string;
    ip: string;
    created_at: string;
    updated_at: string;
  }

  interface Devices {
    devices: Device[];
    total: number;
  }

  interface ListDevicesResponse {
    code: number;
    message: string;
    data: Devices;
  }

  interface CreateDeviceResponse {
    code: number;
    message: string;
    data?: Device;
  }

  interface UpdateDeviceResponse {
    code: number;
    message: string;
    data?: Device;
  }

  interface DeleteDeviceResponse {
    code: number;
    message: string;
  }

  // 请求参数类型
  interface ListDevicesParams {
    page?: number;
    page_size?: number;
    name?: string;
    device_type?: string;
    status?: string;
  }

  interface CreateDeviceParams {
    name: string;
    device_type: string;
    ip: string;
  }

  interface UpdateDeviceParams {
    id: number;
    name?: string;
    device_type?: string;
    ip?: string;
    status?: string;
  }
}
