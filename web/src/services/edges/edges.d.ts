declare namespace Edge {
  interface Edge {
    id: number;
    name: string;
    description: string;
    status: number;
    online: number;
    created_at: string;
    updated_at: string;
  }

  interface Edges {
    edges: Edge[];
    total: number;
  }

  interface ListEdgesResponse {
    code: number;
    message: string;
    data: Edges;
  }

  interface CreateEdgeResponse {
    code: number;
    message: string;
    data?: Edge;
  }

  interface UpdateEdgeResponse {
    code: number;
    message: string;
    data?: Edge;
  }

  interface DeleteEdgeResponse {
    code: number;
    message: string;
  }

  // 请求参数类型
  interface ListEdgesParams {
    page?: number;
    page_size?: number;
    name?: string;
    status?: number;
  }

  interface CreateEdgeParams {
    name: string;
    description: string;
    status?: number;
  }

  interface UpdateEdgeParams {
    id: number;
    name?: string;
    description?: string;
    status?: number;
  }
}
