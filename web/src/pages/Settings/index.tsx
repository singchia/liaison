import { PageContainer } from '@ant-design/pro-components';
import {
  Card,
  Tabs,
  Form,
  Input,
  Button,
  App,
  Descriptions,
  Avatar,
  Typography,
  Divider,
} from 'antd';
import {
  UserOutlined,
  LockOutlined,
  SafetyOutlined,
} from '@ant-design/icons';
import { useState } from 'react';
import { useModel } from '@umijs/max';
import { changePassword } from '@/services/api';
import { executeAction } from '@/utils/request';
import './index.less';

const { Title, Text } = Typography;

const SettingsPage: React.FC = () => {
  const { message } = App.useApp();
  const { initialState } = useModel('@@initialState');
  const [passwordLoading, setPasswordLoading] = useState(false);
  const [passwordForm] = Form.useForm();

  const handleChangePassword = async (values: {
    oldPassword: string;
    newPassword: string;
    confirmPassword: string;
  }) => {
    if (values.newPassword !== values.confirmPassword) {
      message.error('两次输入的新密码不一致');
      return;
    }

    setPasswordLoading(true);
    await executeAction(
      () =>
        changePassword({
          old_password: values.oldPassword,
          new_password: values.newPassword,
        }),
      {
        successMessage: '密码修改成功',
        errorMessage: '密码修改失败',
        onSuccess: () => passwordForm.resetFields(),
      },
    );
    setPasswordLoading(false);
  };

  const items = [
    {
      key: 'account',
      label: (
        <span>
          <UserOutlined />
          账户信息
        </span>
      ),
      children: (
        <div className="settings-section">
          <Card variant="borderless">
            <div className="user-profile">
              <Avatar
                size={80}
                icon={<UserOutlined />}
                src="https://api.dicebear.com/7.x/bottts/svg?seed=Liaison"
              />
              <div className="user-info">
                <Title level={4}>
                  {initialState?.currentUser?.name || 'Admin'}
                </Title>
                <Text type="secondary">
                  {initialState?.currentUser?.email || 'default@liaison.local'}
                </Text>
              </div>
            </div>
            
            <Divider />
            
            <Descriptions
              column={{ xs: 1, sm: 1, md: 2 }}
              styles={{ label: { fontWeight: 500 } }}
            >
              <Descriptions.Item label="用户名">
                {initialState?.currentUser?.name || 'Admin'}
              </Descriptions.Item>
              <Descriptions.Item label="邮箱">
                {initialState?.currentUser?.email || 'default@liaison.local'}
              </Descriptions.Item>
              <Descriptions.Item label="角色">
                {initialState?.currentUser?.role || '管理员'}
              </Descriptions.Item>
              <Descriptions.Item label="注册时间">
                {initialState?.currentUser?.created_at || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="最后登录">
                {initialState?.currentUser?.last_login || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="登录IP">
                {initialState?.currentUser?.login_ip || '-'}
              </Descriptions.Item>
            </Descriptions>
          </Card>
        </div>
      ),
    },
    {
      key: 'password',
      label: (
        <span>
          <LockOutlined />
          修改密码
        </span>
      ),
      children: (
        <div className="settings-section">
          <Card variant="borderless">
            <div className="password-tips">
              <SafetyOutlined className="text-blue-500 text-xl mr-2" />
              <div>
                <Text strong>密码安全提示</Text>
                <br />
                <Text type="secondary">
                  建议定期修改密码，密码长度至少8位，包含字母和数字
                </Text>
              </div>
            </div>
            
            <Divider />
            
            <Form
              form={passwordForm}
              layout="vertical"
              onFinish={handleChangePassword}
              className="password-form"
              requiredMark={false}
            >
              <Form.Item
                name="oldPassword"
                label="当前密码"
                rules={[{ required: true, message: '请输入当前密码' }]}
              >
                <Input.Password
                  prefix={<LockOutlined />}
                  placeholder="请输入当前密码"
                />
              </Form.Item>

              <Form.Item
                name="newPassword"
                label="新密码"
                rules={[
                  { required: true, message: '请输入新密码' },
                  { min: 8, message: '密码长度至少8位' },
                  {
                    pattern: /^(?=.*[A-Za-z])(?=.*\d)/,
                    message: '密码必须包含字母和数字',
                  },
                ]}
              >
                <Input.Password
                  prefix={<LockOutlined />}
                  placeholder="请输入新密码"
                />
              </Form.Item>

              <Form.Item
                name="confirmPassword"
                label="确认新密码"
                dependencies={['newPassword']}
                rules={[
                  { required: true, message: '请确认新密码' },
                  ({ getFieldValue }) => ({
                    validator(_, value) {
                      if (!value || getFieldValue('newPassword') === value) {
                        return Promise.resolve();
                      }
                      return Promise.reject(new Error('两次输入的密码不一致'));
                    },
                  }),
                ]}
              >
                <Input.Password
                  prefix={<LockOutlined />}
                  placeholder="请再次输入新密码"
                />
              </Form.Item>

              <Form.Item>
                <Button
                  type="primary"
                  htmlType="submit"
                  loading={passwordLoading}
                >
                  修改密码
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </div>
      ),
    },
  ];

  return (
    <PageContainer>
      <Card variant="borderless">
        <Tabs
          items={items}
          tabPosition="left"
          className="settings-tabs"
        />
      </Card>
    </PageContainer>
  );
};

export default SettingsPage;
