import { createUser, updateUser } from '@/services/iam';
import { Form, Input, message, Modal, Select } from 'antd';
import React, { useEffect } from 'react';

const { Option } = Select;

interface UserFormModalProps {
  visible: boolean;
  user?: IAM.User;
  onCancel: () => void;
  onSuccess: () => void;
}

const UserFormModal: React.FC<UserFormModalProps> = ({
  visible,
  user,
  onCancel,
  onSuccess,
}) => {
  const [form] = Form.useForm();
  const isEdit = !!user;

  useEffect(() => {
    if (visible) {
      if (user) {
        // 编辑模式，填充数据
        form.setFieldsValue({
          username: user.username,
          email: user.email,
          role: user.role,
          avatar: user.avatar,
        });
      } else {
        // 新增模式，清空表单
        form.resetFields();
      }
    }
  }, [visible, user, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();

      if (isEdit) {
        // 编辑用户
        await updateUser({ id: user.id, ...values });
        message.success('用户信息更新成功');
      } else {
        // 新增用户
        await createUser(values);
        message.success('用户创建成功');
      }

      form.resetFields();
      onSuccess();
    } catch (error) {
      console.error('操作失败', error);
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title={isEdit ? '编辑用户' : '新增用户'}
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

        {!isEdit && (
          <Form.Item
            label="密码"
            name="password"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 6, message: '密码至少6位' },
            ]}
          >
            <Input.Password placeholder="请输入密码" />
          </Form.Item>
        )}

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

        <Form.Item label="角色" name="role">
          <Select placeholder="请选择角色">
            <Option value="admin">管理员</Option>
            <Option value="user">普通用户</Option>
          </Select>
        </Form.Item>

        <Form.Item label="头像地址" name="avatar">
          <Input placeholder="请输入头像URL" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default UserFormModal;
