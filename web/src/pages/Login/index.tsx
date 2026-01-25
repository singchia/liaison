import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { LoginForm, ProFormText } from '@ant-design/pro-components';
import { history, useModel } from '@umijs/max';
import { App } from 'antd';
import { login } from '@/services/api';
import { APP_NAME } from '@/constants';
import './index.less';

const Login: React.FC = () => {
  const { message } = App.useApp();
  const { setInitialState } = useModel('@@initialState');

  const handleSubmit = async (values: { email: string; password: string }) => {
    try {
      const result = await login(values);
      if (result.code === 200 && result.data?.token) {
        localStorage.setItem('token', result.data.token);
        message.success('登录成功！');

        setInitialState((s) => ({
          ...s,
          currentUser: result.data?.user,
        }));

        const urlParams = new URL(window.location.href).searchParams;
        history.push(urlParams.get('redirect') || '/');
        return;
      }
      message.error(result.message || '登录失败');
    } catch (error: any) {
      message.error(error?.message || '登录失败，请重试！');
    }
  };

  return (
    <div className="login-container">
      <div className="login-content">
        <div className="login-header">
          <img 
            src="/liaison.png" 
            alt="Liaison" 
            className="login-logo-img"
          />
          <span className="login-title">{APP_NAME}</span>
        </div>
        
        <LoginForm
          style={{
            minWidth: 280,
            maxWidth: '75vw',
          }}
          submitter={{
            searchConfig: {
              submitText: '登录',
            },
          }}
          onFinish={handleSubmit}
        >
          <ProFormText
            name="email"
            fieldProps={{
              size: 'large',
              prefix: <UserOutlined className="prefixIcon" />,
            }}
            placeholder="邮箱"
            initialValue=""
            rules={[
              {
                required: true,
                message: '请输入邮箱!',
              },
              {
                type: 'email',
                message: '请输入有效的邮箱地址!',
              },
            ]}
          />
          <ProFormText.Password
            name="password"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined className="prefixIcon" />,
            }}
            placeholder="密码"
            initialValue=""
            rules={[
              {
                required: true,
                message: '请输入密码！',
              },
            ]}
          />
        </LoginForm>
        
        <div className="login-footer">
          <p>© 2026 {APP_NAME}. All rights reserved.</p>
        </div>
      </div>
    </div>
  );
};

export default Login;
