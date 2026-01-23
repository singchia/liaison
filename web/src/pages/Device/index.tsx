import {
  ActionType,
  ModalForm,
  PageContainer,
  ProColumns,
  ProDescriptions,
  ProFormText,
  ProFormTextArea,
  ProTable,
} from '@ant-design/pro-components';
import { Drawer, Space, Tag, Popconfirm, Badge, Tooltip, message } from 'antd';
import { DesktopOutlined, InfoCircleOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useRef, useState } from 'react';
import { getDeviceList, getDeviceDetail, updateDevice, deleteDevice } from '@/services/api';
import { executeAction, tableRequest } from '@/utils/request';
import { defaultPagination, defaultSearch, buildSearchParams } from '@/utils/tableConfig';
import { formatMBSize } from '@/utils/format';

const DevicePage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [detailDrawerVisible, setDetailDrawerVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [currentRow, setCurrentRow] = useState<API.Device>();
  const [detailLoading, setDetailLoading] = useState(false);

  const reload = () => actionRef.current?.reload();

  const handleViewDetail = async (record: API.Device) => {
    setDetailDrawerVisible(true);
    setDetailLoading(true);
    try {
      const res = await getDeviceDetail(record.id);
      if (res.code === 200 && res.data) {
        setCurrentRow(res.data);
      } else {
        setCurrentRow(record);
      }
    } catch {
      setCurrentRow(record);
    } finally {
      setDetailLoading(false);
    }
  };

  const handleEdit = async (values: any) => {
    if (!currentRow?.id) return false;
    return executeAction(
      () => updateDevice(currentRow.id, {
        name: values.name,
        description: values.description,
      }),
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

  const handleDelete = async (id: number, online?: number) => {
    // 如果设备在线，不允许删除
    if (online === 1) {
      message.warning('在线设备不允许删除，请先断开连接');
      return;
    }
    await executeAction(() => deleteDevice(id), {
      successMessage: '删除成功',
      errorMessage: '删除失败',
      onSuccess: reload,
    });
  };

  const columns: ProColumns<API.Device>[] = [
    {
      title: '设备名称',
      dataIndex: 'name',
      ellipsis: true,
      fieldProps: {
        placeholder: '请输入设备名称',
        style: { width: 200 },
      },
      render: (_, record) => (
        <Space>
          <DesktopOutlined />
          <span>{record.name || `设备-${record.id}`}</span>
        </Space>
      ),
    },
    {
      title: '在线状态',
      dataIndex: 'online',
      width: 100,
      search: false,
      render: (_, record) => (
        <Badge
          status={record.online === 1 ? 'success' : 'default'}
          text={record.online === 1 ? '在线' : '离线'}
        />
      ),
    },
    {
      title: '操作系统',
      dataIndex: 'os',
      width: 120,
      search: false,
    },
    {
      title: '版本',
      dataIndex: 'version',
      width: 100,
      search: false,
    },
    {
      title: 'CPU',
      dataIndex: 'cpu',
      width: 80,
      search: false,
      render: (cpu) => <Tag>{cpu} 核</Tag>,
    },
    {
      title: '内存',
      dataIndex: 'memory',
      width: 100,
      search: false,
      render: (memory) => formatMBSize(memory as number),
    },
    {
      title: '磁盘',
      dataIndex: 'disk',
      width: 100,
      search: false,
      render: (disk) => formatMBSize(disk as number),
    },
    {
      title: '网卡',
      dataIndex: 'interfaces',
      search: false,
      ellipsis: true,
      render: (_, record) => {
        const interfaces = record.interfaces;
        if (!interfaces || !Array.isArray(interfaces) || interfaces.length === 0) return '-';
        return (
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
            {interfaces.map((iface, index) => {
              if (!iface || !iface.name) return null;
              // 过滤出 IPv4 地址
              const ipv4 = Array.isArray(iface.ip)
                ? iface.ip.filter((ip: string) => ip && !ip.includes(':'))
                : [];
              return (
                <Tag key={index} style={{ margin: 0 }}>
                  {iface.name}: {ipv4.length > 0 ? ipv4.join(', ') : '-'}
                </Tag>
              );
            }).filter(Boolean)}
          </div>
        );
      },
    },
    {
      title: '网卡IP',
      dataIndex: 'ip',
      hideInTable: true,
      tooltip: '支持搜索设备的任意网卡IP地址',
    },
    {
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
      search: false,
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      valueType: 'dateTime',
      width: 180,
      search: false,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 200,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <a onClick={() => handleViewDetail(record)}>
            <InfoCircleOutlined /> 详情
          </a>
          <a onClick={() => {
            setCurrentRow(record);
            setEditModalVisible(true);
          }}>
            <EditOutlined /> 编辑
          </a>
          {record.online === 1 ? (
            <Tooltip title="在线设备不允许删除，请先断开连接">
              <a style={{ color: '#d9d9d9', cursor: 'not-allowed' }}>
                <DeleteOutlined /> 删除
              </a>
            </Tooltip>
          ) : (
            <Popconfirm
              title="确定要删除这个设备吗？"
              description="删除后无法恢复，请谨慎操作"
              onConfirm={() => handleDelete(record.id, record.online)}
              okText="确定"
              cancelText="取消"
              okButtonProps={{ danger: true }}
            >
              <a style={{ color: '#ff4d4f' }}>
                <DeleteOutlined /> 删除
              </a>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.Device>
        headerTitle="设备列表"
        actionRef={actionRef}
        rowKey="id"
        columns={columns}
        request={async (params) => {
          const searchParams = buildSearchParams<API.DeviceListParams>(params, ['name', 'ip']);
          return tableRequest(() => getDeviceList(searchParams), 'devices');
        }}
        pagination={defaultPagination}
        search={defaultSearch}
        scroll={{ x: 'max-content' }}
      />

      <Drawer
        title="设备详情"
        width={600}
        open={detailDrawerVisible}
        onClose={() => setDetailDrawerVisible(false)}
        loading={detailLoading}
      >
        {currentRow && (
          <ProDescriptions<API.Device>
            column={1}
            dataSource={currentRow}
            columns={[
              { title: '设备 ID', dataIndex: 'id', copyable: true },
              { title: '设备名称', dataIndex: 'name' },
              { title: '操作系统', dataIndex: 'os' },
              { title: '版本', dataIndex: 'version' },
              { title: 'CPU', dataIndex: 'cpu', render: (cpu) => `${cpu} 核` },
              { title: '内存', dataIndex: 'memory', render: (memory) => formatMBSize(memory as number) },
              { title: '磁盘', dataIndex: 'disk', render: (disk) => formatMBSize(disk as number) },
              { 
                title: '网卡信息', 
                dataIndex: 'interfaces',
                render: (_, record) => {
                  const interfaces = record.interfaces;
                  if (!interfaces || !Array.isArray(interfaces) || interfaces.length === 0) return '-';
                  return (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                      {interfaces.map((iface, index) => {
                        if (!iface || !iface.name) return null;
                        // 过滤出 IPv4 地址
                        const ipv4 = Array.isArray(iface.ip) 
                          ? iface.ip.filter((ip: string) => ip && !ip.includes(':')) 
                          : [];
                        return (
                          <div key={index}>
                            <Tag>{iface.name}</Tag> {ipv4.length > 0 ? ipv4.join(', ') : '-'}
                            {iface.mac && <div style={{ fontSize: '12px', color: '#999' }}>MAC: {iface.mac}</div>}
                          </div>
                        );
                      }).filter(Boolean)}
                    </div>
                  );
                }
              },
              { title: '描述', dataIndex: 'description' },
              { title: '创建时间', dataIndex: 'created_at', valueType: 'dateTime' },
              { title: '更新时间', dataIndex: 'updated_at', valueType: 'dateTime' },
            ]}
          />
        )}
      </Drawer>

      <ModalForm
        title="编辑设备"
        open={editModalVisible}
        onOpenChange={setEditModalVisible}
        onFinish={handleEdit}
        initialValues={currentRow}
        modalProps={{ destroyOnClose: true }}
        width={500}
      >
        <ProFormText
          name="name"
          label="设备名称"
          placeholder="请输入设备名称"
          rules={[{ required: true, message: '请输入设备名称' }]}
        />
        <ProFormTextArea
          name="description"
          label="描述"
          placeholder="请输入设备描述"
        />
      </ModalForm>
    </PageContainer>
  );
};

export default DevicePage;
