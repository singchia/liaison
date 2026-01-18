import {
  ActionType,
  ModalForm,
  PageContainer,
  ProColumns,
  ProFormDigit,
  ProFormSelect,
  ProFormText,
  ProTable,
} from '@ant-design/pro-components';
import { Space, Tag, Typography, Input } from 'antd';
import { LinkOutlined, ApiOutlined, EditOutlined } from '@ant-design/icons';
import { useRef, useState } from 'react';
import {
  getApplicationList,
  createApplication,
  updateApplication,
  deleteApplication,
  getEdgeList,
  createProxy,
} from '@/services/api';
import { executeAction, tableRequest } from '@/utils/request';
import { CreateButton, DeleteLink } from '@/components/TableButtons';
import { defaultPagination, defaultSearch, buildSearchParams } from '@/utils/tableConfig';

const { Text } = Typography;

const AppPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [proxyModalVisible, setProxyModalVisible] = useState(false);
  const [currentRow, setCurrentRow] = useState<API.Application>();

  const reload = () => actionRef.current?.reload();

  const handleAdd = async (values: any) => {
    return executeAction(
      () =>
        createApplication({
          name: values.name,
          application_type: values.application_type,
          ip: values.ip,
          port: values.port,
          edge_id: values.edge_id,
          device_id: values.device_id,
        }),
      {
        successMessage: '创建成功',
        errorMessage: '创建失败',
        onSuccess: () => {
          setCreateModalVisible(false);
          reload();
        },
      },
    );
  };

  const handleEdit = async (values: any) => {
    if (!currentRow?.id) return false;
    return executeAction(
      () => updateApplication(currentRow.id, { name: values.name }),
      {
        successMessage: '更新成功',
        errorMessage: '更新失败',
        onSuccess: () => {
          setEditModalVisible(false);
          reload();
        },
      },
    );
  };

  const handleDelete = async (id: number) => {
    await executeAction(() => deleteApplication(id), {
      successMessage: '删除成功',
      errorMessage: '删除失败',
      onSuccess: reload,
    });
  };

  const handleCreateProxy = async (values: any) => {
    if (!currentRow?.id) return false;
    return executeAction(
      () =>
        createProxy({
          name: values.name || currentRow.name,
          description: values.description,
          port: values.port,
          application_id: currentRow.id,
        }),
      {
        successMessage: '代理创建成功',
        errorMessage: '代理创建失败',
        onSuccess: () => {
          setProxyModalVisible(false);
          reload();
        },
      },
    );
  };

  const columns: ProColumns<API.Application>[] = [
    {
      title: '应用名称',
      dataIndex: 'name',
      ellipsis: true,
      render: (_, record) => (
        <Space>
          <ApiOutlined />
          <span>{record.name}</span>
        </Space>
      ),
    },
    {
      title: '类型',
      dataIndex: 'application_type',
      width: 100,
      search: false,
      valueEnum: {
        web: { text: 'Web' },
        tcp: { text: 'TCP' },
        udp: { text: 'UDP' },
        ssh: { text: 'SSH' },
        rdp: { text: 'RDP' },
        database: { text: '数据库' },
      },
    },
    {
      title: 'IP 地址',
      dataIndex: 'ip',
      width: 140,
      search: false,
      render: (ip) => <Text code>{ip}</Text>,
    },
    {
      title: '端口',
      dataIndex: 'port',
      width: 80,
      search: false,
      render: (port) => <Tag>{port}</Tag>,
    },
    {
      title: '所在设备',
      dataIndex: 'device_name',
      ellipsis: true,
      width: 150,
      render: (_, record) => record.device?.name || '-',
      renderFormItem: () => {
        return <Input placeholder="请输入设备名称" />;
      },
    },
    {
      title: '已关联代理',
      dataIndex: 'proxy',
      ellipsis: true,
      width: 150,
      search: false,
      render: (_, record) => {
        if (record.proxy) {
          return (
            <Tag color="blue">
              <LinkOutlined /> {record.proxy.name}:{record.proxy.port}
            </Tag>
          );
        }
        return <Tag>未关联</Tag>;
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 170,
      search: false,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 200,
      render: (_, record) => (
        <Space>
          <a onClick={() => {
            setCurrentRow(record);
            setProxyModalVisible(true);
          }}>
            <LinkOutlined /> 创建代理
          </a>
          <a onClick={() => {
            setCurrentRow(record);
            setEditModalVisible(true);
          }}>
            <EditOutlined /> 编辑
          </a>
          <DeleteLink
            title="确定要删除这个应用吗？"
            onConfirm={() => handleDelete(record.id)}
          />
        </Space>
      ),
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.Application>
        headerTitle="应用列表"
        actionRef={actionRef}
        rowKey="id"
        columns={columns}
        request={async (params) => {
          const searchParams = buildSearchParams<API.ApplicationListParams>(params, ['name', 'device_name']);
          return tableRequest(() => getApplicationList(searchParams), 'applications');
        }}
        toolBarRender={() => [
          <CreateButton key="create" onClick={() => setCreateModalVisible(true)}>
            新建应用
          </CreateButton>,
        ]}
        pagination={defaultPagination}
        search={defaultSearch}
        scroll={{ x: 'max-content' }}
      />

      <ModalForm
        title="新建应用"
        open={createModalVisible}
        onOpenChange={setCreateModalVisible}
        onFinish={handleAdd}
        modalProps={{ destroyOnClose: true }}
        width={500}
      >
        <ProFormText
          name="name"
          label="应用名称"
          placeholder="请输入应用名称"
          rules={[{ required: true, message: '请输入应用名称' }]}
        />
        <ProFormSelect
          name="application_type"
          label="应用类型"
          placeholder="请选择应用类型"
          options={[
            { label: 'Web', value: 'web' },
            { label: 'TCP', value: 'tcp' },
            { label: 'UDP', value: 'udp' },
            { label: 'SSH', value: 'ssh' },
            { label: 'RDP', value: 'rdp' },
            { label: '数据库', value: 'database' },
          ]}
          rules={[{ required: true, message: '请选择应用类型' }]}
        />
        <ProFormText
          name="ip"
          label="IP 地址"
          placeholder="请输入应用 IP 地址，如 192.168.1.100"
          rules={[{ required: true, message: '请输入 IP 地址' }]}
        />
        <ProFormDigit
          name="port"
          label="端口"
          placeholder="请输入端口号"
          min={1}
          max={65535}
          rules={[{ required: true, message: '请输入端口号' }]}
        />
        <ProFormSelect
          name="edge_id"
          label="连接器"
          placeholder="请选择连接器"
          rules={[{ required: true, message: '请选择连接器' }]}
          request={async () => {
            try {
              const res = await getEdgeList({ page_size: 100 });
              return (
                res.data?.edges?.map((item: API.Edge) => ({
                  label: item.name,
                  value: item.id,
                })) || []
              );
            } catch {
              return [];
            }
          }}
        />
      </ModalForm>

      <ModalForm
        title="编辑应用"
        open={editModalVisible}
        onOpenChange={setEditModalVisible}
        onFinish={handleEdit}
        modalProps={{ destroyOnClose: true }}
        initialValues={currentRow}
        width={500}
      >
        <ProFormText
          name="name"
          label="应用名称"
          placeholder="请输入应用名称"
          rules={[{ required: true, message: '请输入应用名称' }]}
        />
      </ModalForm>

      <ModalForm
        title="为应用创建代理"
        open={proxyModalVisible}
        onOpenChange={setProxyModalVisible}
        onFinish={handleCreateProxy}
        modalProps={{ destroyOnClose: true }}
        width={500}
      >
        <ProFormText
          name="name"
          label="代理名称"
          placeholder="请输入代理名称"
          initialValue={currentRow?.name}
          rules={[{ required: true, message: '请输入代理名称' }]}
        />
        <ProFormDigit
          name="port"
          label="公网端口"
          placeholder="留空自动分配"
          min={1}
          max={65535}
          extra="映射到公网的端口号，留空则自动分配"
        />
        <ProFormText
          name="description"
          label="描述"
          placeholder="请输入描述"
        />
      </ModalForm>
    </PageContainer>
  );
};

export default AppPage;
