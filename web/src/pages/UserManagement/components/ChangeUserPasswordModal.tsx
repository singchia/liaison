import { changeUserPassword } from '@/services/iam';
import { Form, Input, message, Modal } from 'antd';
import React from 'react';

interface ChangeUserPasswordModalProps {
  visible: boolean;
  user?: IAM.User;
  onCancel: () => void;
  onSuccess: () => void;
}

const ChangeUserPasswordModal: React.FC<ChangeUserPasswordModalProps> = ({
  visible,
  user,
  onCancel,
  onSuccess,
}) => {
  const [form] = Form.useForm();

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const { newPassword, confirmPassword } = values;

      if (newPassword !== confirmPassword) {
        message.error('两次输入的密码不一致');
        return;
      }

      if (!user) {
        message.error('用户信息不存在');
        return;
      }

      await changeUserPassword({
        id: user.id,
        newPassword,
      });

      message.success('密码修改成功');
      form.resetFields();
      onSuccess();
    } catch (error) {
      console.error('密码修改失败', error);
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title={`修改用户 "${user?.username}" 的密码`}
      open={visible}
      onCancel={handleCancel}
      onOk={handleSubmit}
      width={500}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          label="新密码"
          name="newPassword"
          rules={[
            { required: true, message: '请输入新密码' },
            { min: 6, message: '密码至少6位' },
          ]}
        >
          <Input.Password placeholder="请输入新密码" />
        </Form.Item>

        <Form.Item
          label="确认密码"
          name="confirmPassword"
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
          <Input.Password placeholder="请再次输入新密码" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default ChangeUserPasswordModal;
