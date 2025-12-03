import { useQueryParams } from '@/hooks/useQueryParams';
import { getApplications } from '@/services/applications';
import { SearchOutlined } from '@ant-design/icons';
import { useRequest } from 'ahooks';
import { Button, Card, Form, Input, Select, Space, Table, Tag } from 'antd';
import React, { useEffect, useState } from 'react';

const { Option } = Select;

const AppList: React.FC = () => {
  const [form] = Form.useForm();
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // 使用 hook 管理 URL query 参数
  const { params, updateParams } =
    useQueryParams<Application.ListApplicationsParams>({
      defaultParams: {
        page: 1,
        page_size: 10,
      },
    });

  // 使用 ahooks useRequest 进行请求
  const { data, loading } = useRequest(() => getApplications(params), {
    refreshDeps: [params],
  });

  // 从 params 回填表单
  useEffect(() => {
    form.setFieldsValue({
      name: params.name,
      status: params.status,
    });
  }, [params, form]);

  // 搜索处理
  const handleSearch = (values: any) => {
    updateParams({
      ...values,
      page: 1,
    });
  };

  // 分页处理
  const handleTableChange = (pagination: any) => {
    updateParams({
      page: pagination.current,
      page_size: pagination.pageSize,
    });
  };

  // 行选择处理
  const rowSelection = {
    selectedRowKeys,
    onChange: (keys: React.Key[]) => {
      setSelectedRowKeys(keys);
    },
  };

  // 状态标签渲染
  const renderStatus = (status: string) => {
    const statusConfig = {
      运行中: { color: 'blue', text: '运行中' },
      关闭: { color: 'red', text: '关闭' },
      已上线: { color: 'green', text: '已上线' },
      异常: { color: 'red', text: '异常' },
    };
    const config = statusConfig[status as keyof typeof statusConfig] || {
      color: 'default',
      text: status,
    };
    return <Tag color={config.color}>{config.text}</Tag>;
  };

  // 操作列渲染
  const renderActions = () => (
    <Space>
      <a>配置</a>
      <a>订阅警报</a>
    </Space>
  );

  // 表格列配置
  const columns = [
    {
      title: '应用名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '应用类型',
      dataIndex: 'application_type',
      key: 'application_type',
    },
    {
      title: 'IP地址',
      dataIndex: 'ip',
      key: 'ip',
    },
    {
      title: '端口',
      dataIndex: 'port',
      key: 'port',
      sorter: true,
    },
    {
      title: '设备',
      dataIndex: ['device', 'name'],
      key: 'device',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: renderStatus,
      filters: [
        { text: '运行中', value: '运行中' },
        { text: '关闭', value: '关闭' },
        { text: '已上线', value: '已上线' },
        { text: '异常', value: '异常' },
      ],
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      key: 'updated_at',
      sorter: true,
    },
    {
      title: '操作',
      key: 'action',
      render: renderActions,
    },
  ];

  // 计算总服务调用数

  return (
    <Card>
      {/* 搜索表单 */}
      <Form
        form={form}
        layout="inline"
        onFinish={handleSearch}
        style={{ marginBottom: 16 }}
      >
        <Form.Item name="name" label="应用名:">
          <Input placeholder="请输入" style={{ width: 200 }} />
        </Form.Item>
        <Form.Item name="status" label="状态:">
          <Select placeholder="请选择" style={{ width: 200 }} allowClear>
            <Option value="运行中">运行中</Option>
            <Option value="关闭">关闭</Option>
            <Option value="已上线">已上线</Option>
            <Option value="异常">异常</Option>
          </Select>
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
            查询
          </Button>
        </Form.Item>
      </Form>

      {/* 数据表格 */}
      <Table
        rowSelection={rowSelection}
        columns={columns}
        dataSource={data?.applications || []}
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
    </Card>
  );
};

export default AppList;
