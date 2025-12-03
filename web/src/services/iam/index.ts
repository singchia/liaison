// IAM（身份认证管理）相关的API服务
import { request } from '@/services/request';

/**
 * 用户登录
 */
export async function login(
  params: IAM.LoginParams,
): Promise<IAM.LoginResponse['data']> {
  return request('POST /v1/iam/login', {
    data: params,
  });
}

/**
 * 用户登出
 */
export async function logout(): Promise<IAM.LogoutResponse['message']> {
  return request('POST /v1/iam/logout');
}

/**
 * 获取当前用户信息
 */
export async function getProfile(): Promise<IAM.ProfileResponse['data']> {
  return request('/v1/iam/profile');
}

/**
 * 更新当前用户信息
 */
export async function updateProfile(
  params: IAM.UpdateProfileParams,
): Promise<IAM.UpdateProfileResponse['data']> {
  return request('PUT /v1/iam/profile', {
    data: params,
  });
}

/**
 * 获取用户列表（管理员功能）
 */
export async function getUsers(
  params: IAM.ListUsersParams,
): Promise<IAM.ListUsersResponse['data']> {
  return request('/v1/iam/users', { params });
}

/**
 * 创建用户（管理员功能）
 */
export async function createUser(
  params: IAM.CreateUserParams,
): Promise<IAM.CreateUserResponse['data']> {
  return request('POST /v1/iam/users', {
    data: params,
  });
}

/**
 * 更新用户信息（管理员功能）
 */
export async function updateUser(
  params: IAM.UpdateUserParams,
): Promise<IAM.UpdateUserResponse['data']> {
  const { id, ...data } = params;
  return request('PUT /v1/iam/users/:id', {
    params: { id },
    data,
  });
}

/**
 * 删除用户（管理员功能）
 */
export async function deleteUser(
  id: number,
): Promise<IAM.DeleteUserResponse['message']> {
  return request('DELETE /v1/iam/users/:id', { params: { id } });
}

/**
 * 修改用户密码（管理员功能）
 */
export async function changeUserPassword(
  params: IAM.ChangeUserPasswordParams,
): Promise<IAM.ChangeUserPasswordResponse['message']> {
  const { id, ...data } = params;
  return request('PUT /v1/iam/users/:id/password', {
    params: { id },
    data,
  });
}
