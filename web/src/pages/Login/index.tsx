import { login } from '@/services/iam';
import { Button, Card, Form, Input, message } from 'antd';
import React, { useState } from 'react';

const Login: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  const onFinish = async (values: { username: string; password: string }) => {
    const { username, password } = values || ({} as any);
    if (!username || !password) {
      message.error('请输入用户名和密码');
      return;
    }

    try {
      setLoading(true);
      const data = await login({ username, password });

      // 保存 token 和用户信息
      localStorage.setItem('token', data.token);
      localStorage.setItem('username', data.user.username);
      localStorage.setItem('currentUser', JSON.stringify(data.user));

      message.success('登录成功');
      window.location.href = '/home';
    } catch (error) {
      // 错误已在 request 封装中处理
      console.error('登录失败', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: '#f5f5f5',
        padding: 24,
      }}
    >
      <Card title="登录" style={{ width: 360 }}>
        <Form form={form} layout="vertical" onFinish={onFinish}>
          <Form.Item
            label="用户名"
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="用户名" allowClear />
          </Form.Item>
          <Form.Item
            label="密码"
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password placeholder="密码" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading}>
              登录
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default Login;
