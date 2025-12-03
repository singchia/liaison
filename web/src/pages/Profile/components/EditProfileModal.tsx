import { updateProfile } from '@/services/iam';
import { Form, Input, message, Modal } from 'antd';
import React, { useEffect } from 'react';

interface EditProfileModalProps {
  visible: boolean;
  user: IAM.User | undefined;
  onCancel: () => void;
  onSuccess: () => void;
}

const EditProfileModal: React.FC<EditProfileModalProps> = ({
  visible,
  user,
  onCancel,
  onSuccess,
}) => {
  const [form] = Form.useForm();

  useEffect(() => {
    if (visible && user) {
      form.setFieldsValue({
        username: user.username,
        email: user.email,
        avatar: user.avatar,
      });
    }
  }, [visible, user, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const updatedUser = await updateProfile(values);
      message.success('更新成功');

      // 更新本地存储
      localStorage.setItem('username', updatedUser.username);
      localStorage.setItem('currentUser', JSON.stringify(updatedUser));

      form.resetFields();
      onSuccess();
    } catch (error) {
      console.error('更新失败', error);
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title="编辑资料"
      open={visible}
      onCancel={handleCancel}
      onOk={handleSubmit}
      width={600}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          label="用户名"
          name="username"
          rules={[{ required: true, message: '请输入用户名' }]}
        >
          <Input placeholder="请输入用户名" />
        </Form.Item>

        <Form.Item
          label="邮箱"
          name="email"
          rules={[
            { required: true, message: '请输入邮箱' },
            { type: 'email', message: '请输入有效的邮箱地址' },
          ]}
        >
          <Input placeholder="请输入邮箱" />
        </Form.Item>

        <Form.Item label="头像地址" name="avatar">
          <Input placeholder="请输入头像URL" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default EditProfileModal;
