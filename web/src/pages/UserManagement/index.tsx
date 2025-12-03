import { useQueryParams } from '@/hooks/useQueryParams';
import { deleteUser, getUsers } from '@/services/iam';
import { DeleteOutlined, LockOutlined, PlusOutlined } from '@ant-design/icons';
import { useRequest } from 'ahooks';
import {
  Avatar,
  Button,
  Card,
  Form,
  Input,
  message,
  Modal,
  Select,
  Space,
  Table,
  Tag,
} from 'antd';
import React, { useEffect, useState } from 'react';
import ChangeUserPasswordModal from './components/ChangeUserPasswordModal';
import UserFormModal from './components/UserFormModal';

const { Option } = Select;

const UserManagement: React.FC = () => {
  const [form] = Form.useForm();
  const [userFormVisible, setUserFormVisible] = useState(false);
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);
  const [currentUser, setCurrentUser] = useState<IAM.User | undefined>(
    undefined,
  );

  // 使用 hook 管理 URL query 参数
  const { params, updateParams } = useQueryParams<IAM.ListUsersParams>({
    defaultParams: {
      page: 1,
      page_size: 10,
    },
  });

  // 使用 ahooks useRequest 进行请求
  const { data, loading, refresh } = useRequest(() => getUsers(params), {
    refreshDeps: [params],
  });

  // 从 params 回填表单
  useEffect(() => {
    form.setFieldsValue({
      username: params.username,
      email: params.email,
      role: params.role,
    });
  }, [params, form]);

  // 搜索处理
  const handleSearch = (values: any) => {
    updateParams({
      ...values,
      page: 1,
    });
  };

  // 新增用户
  const handleAdd = () => {
    setCurrentUser(undefined);
    setUserFormVisible(true);
  };

  // 编辑用户
  const handleEdit = (record: IAM.User) => {
    setCurrentUser(record);
    setUserFormVisible(true);
  };

  // 修改密码
  const handleChangePassword = (record: IAM.User) => {
    setCurrentUser(record);
    setPasswordModalVisible(true);
  };

  // 删除用户
  const handleDelete = (record: IAM.User) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除用户 "${record.username}" 吗？`,
      onOk: async () => {
        try {
          await deleteUser(record.id);
          message.success('删除成功');
          refresh();
        } catch (error) {
          console.error('删除失败', error);
        }
      },
    });
  };

  // 分页处理
  const handleTableChange = (pagination: any) => {
    updateParams({
      page: pagination.current,
      page_size: pagination.pageSize,
    });
  };

  // 表格列配置
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '头像',
      dataIndex: 'avatar',
      key: 'avatar',
      width: 80,
      render: (avatar: string, record: IAM.User) => (
        <Avatar src={avatar} icon={<DeleteOutlined />}>
          {record.username?.[0]?.toUpperCase()}
        </Avatar>
      ),
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (role: string) => {
        const roleConfig = {
          admin: { color: 'red', text: '管理员' },
          user: { color: 'blue', text: '普通用户' },
        };
        const config = roleConfig[role as keyof typeof roleConfig] || {
          color: 'default',
          text: role || '未设置',
        };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      key: 'updated_at',
      width: 180,
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: IAM.User) => (
        <Space>
          <a onClick={() => handleEdit(record)}>编辑</a>
          <a onClick={() => handleChangePassword(record)}>
            <LockOutlined /> 改密码
          </a>
          <a style={{ color: 'red' }} onClick={() => handleDelete(record)}>
            删除
          </a>
        </Space>
      ),
    },
  ];

  return (
    <Card>
      {/* 搜索表单 */}
      <Form
        form={form}
        layout="inline"
        onFinish={handleSearch}
        style={{ marginBottom: 16 }}
      >
        <Form.Item name="username" label="用户名">
          <Input placeholder="请输入" style={{ width: 200 }} />
        </Form.Item>
        <Form.Item name="email" label="邮箱">
          <Input placeholder="请输入" style={{ width: 200 }} />
        </Form.Item>
        <Form.Item name="role" label="角色">
          <Select placeholder="请选择" style={{ width: 150 }} allowClear>
            <Option value="admin">管理员</Option>
            <Option value="user">普通用户</Option>
          </Select>
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">
            查询
          </Button>
        </Form.Item>
      </Form>

      {/* 操作按钮 */}
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          新增用户
        </Button>
      </div>

      {/* 数据表格 */}
      <Table
        columns={columns}
        dataSource={data?.users || []}
        loading={loading}
        rowKey="id"
        pagination={{
          current: params.page,
          pageSize: params.page_size,
          total: data?.total || 0,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `总共 ${total} 条`,
          pageSizeOptions: ['10', '20', '50', '100'],
        }}
        onChange={handleTableChange}
      />

      {/* 新增/编辑用户弹框 */}
      <UserFormModal
        visible={userFormVisible}
        user={currentUser}
        onCancel={() => {
          setUserFormVisible(false);
          setCurrentUser(undefined);
        }}
        onSuccess={() => {
          setUserFormVisible(false);
          setCurrentUser(undefined);
          refresh();
        }}
      />

      {/* 修改密码弹框 */}
      <ChangeUserPasswordModal
        visible={passwordModalVisible}
        user={currentUser}
        onCancel={() => {
          setPasswordModalVisible(false);
          setCurrentUser(undefined);
        }}
        onSuccess={() => {
          setPasswordModalVisible(false);
          setCurrentUser(undefined);
        }}
      />
    </Card>
  );
};

export default UserManagement;
