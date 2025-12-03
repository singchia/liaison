// IAM（身份认证管理）相关的类型定义

declare namespace IAM {
  interface User {
    id: number;
    username: string;
    email?: string;
    avatar?: string;
    role?: string;
    created_at: string;
    updated_at: string;
  }

  interface LoginParams {
    username: string;
    password: string;
  }

  interface LoginResponse {
    code: number;
    message: string;
    data: {
      token: string;
      user: User;
    };
  }

  interface LogoutResponse {
    code: number;
    message: string;
  }

  interface ProfileResponse {
    code: number;
    message: string;
    data: User;
  }

  interface UpdateProfileParams {
    username?: string;
    email?: string;
    avatar?: string;
  }

  interface UpdateProfileResponse {
    code: number;
    message: string;
    data: User;
  }

  // 用户管理相关类型
  interface ListUsersParams {
    page?: number;
    page_size?: number;
    username?: string;
    email?: string;
    role?: string;
  }

  interface Users {
    users: User[];
    total: number;
  }

  interface ListUsersResponse {
    code: number;
    message: string;
    data: Users;
  }

  interface CreateUserParams {
    username: string;
    password: string;
    email?: string;
    role?: string;
    avatar?: string;
  }

  interface CreateUserResponse {
    code: number;
    message: string;
    data: User;
  }

  interface UpdateUserParams {
    id: number;
    username?: string;
    email?: string;
    role?: string;
    avatar?: string;
  }

  interface UpdateUserResponse {
    code: number;
    message: string;
    data: User;
  }

  interface DeleteUserResponse {
    code: number;
    message: string;
  }

  interface ChangeUserPasswordParams {
    id: number;
    newPassword: string;
  }

  interface ChangeUserPasswordResponse {
    code: number;
    message: string;
  }
}
