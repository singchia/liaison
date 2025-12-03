declare namespace Application {
  interface Device {
    id: number;
    name: string;
    created_at: string;
    updated_at: string;
  }

  interface Application {
    id: number;
    name: string;
    application_type: string;
    ip: string;
    port: number;
    edge_id: number;
    device: Device;
    created_at: string;
    updated_at: string;
  }

  interface Applications {
    applications: Application[];
    total: number;
  }

  interface ListApplicationsResponse {
    code: number;
    message: string;
    data: Applications;
  }

  interface CreateApplicationResponse {
    code: number;
    message: string;
    data?: Application;
  }

  interface UpdateApplicationResponse {
    code: number;
    message: string;
    data?: Application;
  }

  interface DeleteApplicationResponse {
    code: number;
    message: string;
  }

  // 请求参数类型
  interface ListApplicationsParams {
    device_id?: number;
    page?: number;
    page_size?: number;
    name?: string;
    application_type?: string;
    status?: string;
  }

  interface CreateApplicationParams {
    application_type: string;
    name: string;
    ip: string;
    port: number;
    device_id: number;
  }

  interface UpdateApplicationParams {
    id: number;
    name?: string;
    ip?: string;
    port?: number;
    application_type?: string;
  }
}
