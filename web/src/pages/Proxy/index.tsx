import {
  ActionType,
  ModalForm,
  PageContainer,
  ProColumns,
  ProFormDigit,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
  ProTable,
} from '@ant-design/pro-components';
import { Space, Tag, Typography, Switch, Alert, Tooltip, message, Button, Modal, Select, Spin, Popconfirm } from 'antd';
import { CheckCircleOutlined } from '@ant-design/icons';
import { useRef, useState, useEffect } from 'react';
import { useSearchParams, useLocation } from '@umijs/max';
import {
  getProxyList,
  createProxy,
  updateProxy,
  deleteProxy,
  getApplicationList,
  getProxyFirewall,
  upsertProxyFirewall,
  deleteProxyFirewall,
} from '@/services/api';
import { executeAction, tableRequest } from '@/utils/request';
import { CreateButton, EditLink, DeleteLink } from '@/components/TableButtons';
import { defaultPagination, defaultSearch, buildSearchParams } from '@/utils/tableConfig';
import { useI18n } from '@/i18n';

const { Text } = Typography;

const ProxyPage: React.FC = () => {
  const { tr } = useI18n();
  const actionRef = useRef<ActionType>();
  const createFormRef = useRef<any>();
  const location = useLocation();
  const [searchParams] = useSearchParams();
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [currentRow, setCurrentRow] = useState<API.Proxy>();
  const [initialApplicationId, setInitialApplicationId] = useState<number | undefined>();
  const [applicationOptions, setApplicationOptions] = useState<
    { label: string; value: number; application_type?: string }[]
  >([]);
  const [applicationMap, setApplicationMap] = useState<Map<number, API.Application>>(new Map());
  const [selectedApplicationId, setSelectedApplicationId] = useState<number | undefined>();
  const hasProcessedUrlRef = useRef(false); // 使用 ref 跟踪是否已处理过 URL 参数

  // 防火墙 Modal 状态
  const [firewallModalVisible, setFirewallModalVisible] = useState(false);
  const [firewallProxy, setFirewallProxy] = useState<API.Proxy | undefined>();
  const [firewallCidrs, setFirewallCidrs] = useState<string[]>([]);
  const [firewallUpdatedAt, setFirewallUpdatedAt] = useState<string>('');
  const [firewallLoading, setFirewallLoading] = useState(false);
  const [firewallSaving, setFirewallSaving] = useState(false);

  // CIDR 校验：v4 (e.g. 10.0.0.0/8) 或 v6 (包含冒号)。宽松匹配，后端会最终判定。
  const isValidCidr = (s: string): boolean => {
    if (!s) return false;
    const v4 = /^(\d{1,3}\.){3}\d{1,3}\/\d{1,2}$/;
    return v4.test(s) || (s.includes(':') && s.includes('/'));
  };

  const openFirewallModal = async (proxy: API.Proxy) => {
    setFirewallProxy(proxy);
    setFirewallCidrs([]);
    setFirewallUpdatedAt('');
    setFirewallModalVisible(true);
    setFirewallLoading(true);
    try {
      const res = await getProxyFirewall(proxy.id);
      if (res.code === 200 && res.data) {
        // 后端：没有规则时返回 ["0.0.0.0/0"] + updated_at 空，此时 UI 显示「未设置」
        const hasRule = !!res.data.updated_at;
        setFirewallCidrs(hasRule ? (res.data.allowed_cidrs || []) : []);
        setFirewallUpdatedAt(res.data.updated_at || '');
      }
    } catch (err: any) {
      message.error(err?.message || tr('加载防火墙规则失败', 'Failed to load firewall rule'));
    } finally {
      setFirewallLoading(false);
    }
  };

  const handleFirewallSave = async () => {
    if (!firewallProxy) return;
    const invalid = firewallCidrs.find((c) => !isValidCidr(c));
    if (invalid) {
      message.error(tr(`无效 CIDR：${invalid}`, `Invalid CIDR: ${invalid}`));
      return;
    }
    setFirewallSaving(true);
    await executeAction(
      () => upsertProxyFirewall(firewallProxy.id, { allowed_cidrs: firewallCidrs }),
      {
        successMessage:
          firewallCidrs.length === 0
            ? tr('已保存（空规则=拒绝全部）', 'Saved (empty rule = deny all)')
            : tr('已保存', 'Saved'),
        errorMessage: tr('保存失败', 'Failed to save'),
        onSuccess: () => {
          setFirewallModalVisible(false);
        },
      },
    );
    setFirewallSaving(false);
  };

  const handleFirewallReset = async () => {
    if (!firewallProxy) return;
    await executeAction(() => deleteProxyFirewall(firewallProxy.id), {
      successMessage: tr('已恢复为放行全部', 'Reset to allow-all'),
      errorMessage: tr('操作失败', 'Operation failed'),
      onSuccess: () => {
        setFirewallCidrs([]);
        setFirewallUpdatedAt('');
        setFirewallModalVisible(false);
      },
    });
  };

  // 从 URL 查询参数中读取 application_id、application_name 和 autoCreate（只执行一次）
  useEffect(() => {
    // 如果已经处理过 URL 参数，不再重复处理
    if (hasProcessedUrlRef.current) return;
    
    // 优先从 searchParams 读取
    let applicationId = searchParams.get('application_id');
    let applicationName = searchParams.get('application_name');
    let autoCreate = searchParams.get('autoCreate');
    
    // 如果 searchParams 没有，从 location.search 读取
    if (location.search) {
      const urlParams = new URLSearchParams(location.search);
      if (!applicationId) {
        applicationId = urlParams.get('application_id');
      }
      if (!applicationName) {
        applicationName = urlParams.get('application_name');
      }
      if (!autoCreate) {
        autoCreate = urlParams.get('autoCreate');
      }
    }
    
    // 如果 application_id 存在，设置初始值
    if (applicationId) {
      const id = parseInt(applicationId, 10);
      if (!isNaN(id)) {
        setInitialApplicationId(id);
        
        // 如果 URL 中有应用名称，立即添加到 options 中（避免等待列表加载）
        if (applicationName) {
          const decodedName = decodeURIComponent(applicationName);
          setApplicationOptions((prev) => {
            // 检查是否已经在 options 中
            const exists = prev.some(opt => opt.value === id);
            if (exists) {
              // 如果已存在，更新 label（使用URL中的名称）
              return prev.map(opt => 
                opt.value === id ? { ...opt, label: decodedName } : opt
              );
            }
            // 如果不存在，添加到 options 中
            return [...prev, { label: decodedName, value: id }];
          });
        }
        
        // 重新加载应用列表，确保包含完整信息（IP和端口）
        getApplicationList({ page_size: 100 }).then((res) => {
          const apps = res.data?.applications || [];
          const appMap = new Map<number, API.Application>();
          apps.forEach((app: API.Application) => {
            appMap.set(app.id, app);
          });
          setApplicationMap(appMap);
          
          const options =
            apps.map((item: API.Application) => ({
              label: `${item.name} (${item.ip}:${item.port})`,
              value: item.id,
              application_type: item.application_type,
            })) || [];
          // 如果URL中有应用名称，确保对应的选项存在（即使列表中没有）
          if (applicationName) {
            const decodedName = decodeURIComponent(applicationName);
            const exists = options.some(opt => opt.value === id);
            if (!exists) {
              // 尝试从 appMap 中获取应用类型，如果找不到则默认为空
              const app = appMap.get(id);
              options.push({ 
                label: decodedName, 
                value: id,
                application_type: app?.application_type || '',
              });
            }
          }
          setApplicationOptions(options);
        }).catch(() => {
          // 忽略错误，但如果URL中有应用名称，至少保留它
          if (applicationName) {
            const decodedName = decodeURIComponent(applicationName);
            setApplicationOptions((prev) => {
              const exists = prev.some(opt => opt.value === id);
              if (exists) return prev;
              return [...prev, { label: decodedName, value: id, application_type: '' }];
            });
          }
        });
      }
    }
    
    // 如果 autoCreate 为 true，自动打开对话框
    if (autoCreate === 'true') {
      hasProcessedUrlRef.current = true; // 标记为已处理
      setCreateModalVisible(true);
      // 清除 URL 中的查询参数
      window.history.replaceState({}, '', '/proxy');
    } else if (applicationId) {
      // 如果没有 autoCreate，但有 application_id，也打开对话框（向后兼容）
      hasProcessedUrlRef.current = true; // 标记为已处理
      setCreateModalVisible(true);
      // 清除 URL 中的查询参数
      window.history.replaceState({}, '', '/proxy');
    } else {
      // 如果没有相关参数，也标记为已处理，避免重复检查
      hasProcessedUrlRef.current = true;
    }
  }, [searchParams, location.search]);

  // 页面加载时就拉应用列表
  useEffect(() => {
    const loadApplications = async () => {
      try {
        const res = await getApplicationList({ page_size: 100 });
        const apps = res.data?.applications || [];
        const appMap = new Map<number, API.Application>();
        apps.forEach((app: API.Application) => {
          appMap.set(app.id, app);
        });
        setApplicationMap(appMap);
        
        const options =
          apps.map((item: API.Application) => ({
            label: `${item.name} (${item.ip}:${item.port})`,
            value: item.id,
            application_type: item.application_type,
          })) || [];
        setApplicationOptions(options);
      } catch {
        setApplicationOptions([]);
      }
    };

    loadApplications();
  }, []);



  const reload = () => actionRef.current?.reload();

  const handleAdd = async (values: any) => {
    const createPort = values.port || undefined;
    
    const result = await executeAction(
      () => createProxy({
        name: values.name,
        description: values.description,
        port: createPort,
        application_id: values.application_id,
      }),
      {
        successMessage: tr('创建成功', 'Created successfully'),
        errorMessage: tr('创建失败', 'Create failed'),
        onSuccess: () => {
          // 如果创建时端口为空，后端会动态分配端口并在响应中返回
          // 刷新列表即可显示动态分配的端口
        },
      },
    );
    
    setCreateModalVisible(false);
    reload();
    
    return result;
  };

  const handleEdit = async (values: any) => {
    if (!currentRow?.id) return false;
    return executeAction(
      () => updateProxy(currentRow.id, {
        name: values.name,
        description: values.description,
        port: values.port,
      }),
      {
        successMessage: tr('更新成功', 'Updated successfully'),
        errorMessage: tr('更新失败', 'Update failed'),
        onSuccess: () => {
          setEditModalVisible(false);
          reload();
        },
      },
    );
  };

  const handleDelete = async (id: number) => {
    await executeAction(() => deleteProxy(id), {
      successMessage: tr('删除成功', 'Deleted successfully'),
      errorMessage: tr('删除失败', 'Delete failed'),
      onSuccess: reload,
    });
  };

  const columns: ProColumns<API.Proxy>[] = [
    {
      title: tr('访问名称', 'Entry Name'),
      dataIndex: 'name',
      ellipsis: true,
      copyable: true,
      fieldProps: {
        placeholder: tr('请输入访问名称', 'Please input entry name'),
      },
    },
    {
      title: tr('描述', 'Description'),
      dataIndex: 'description',
      ellipsis: true,
      search: false,
    },
    {
      title: tr('公网端口', 'Public Port'),
      dataIndex: 'port',
      width: 100,
      search: false,
      render: (port) => <Tag color="blue">{port}</Tag>,
    },
    {
      title: tr('关联应用', 'Application'),
      dataIndex: ['application', 'name'],
      ellipsis: true,
      search: false,
      render: (_, record) =>
        record.application ? (
          <Space direction="vertical" size={0}>
            <Space>
              <Text>{record.application.name}</Text>
              {record.application.application_type === 'http' && (
                <Tooltip title={<span style={{ fontSize: '11px' }}>{tr('已开启 HTTPS', 'HTTPS enabled')}</span>}>
                  <CheckCircleOutlined style={{ color: '#52c41a', fontSize: 16 }} />
                </Tooltip>
              )}
            </Space>
            <Text type="secondary" className="text-xs">
              {record.application.ip}:{record.application.port}
            </Text>
          </Space>
        ) : (
          '-'
        ),
    },
    {
      title: tr('状态', 'Status'),
      dataIndex: 'status',
      width: 100,
      search: false,
      valueEnum: {
        running: { text: tr('运行中', 'Running'), status: 'Success' },
        stopped: { text: tr('已停止', 'Stopped'), status: 'Default' },
        error: { text: tr('异常', 'Error'), status: 'Error' },
      },
    },
    {
      title: tr('启用', 'Enabled'),
      dataIndex: 'enabled',
      width: 80,
      search: false,
      align: 'center',
      render: (_, record) => (
        <Switch
          checked={record.status === 'running'}
          onChange={async (checked) => {
            const newStatus = checked ? 'running' : 'stopped';
            await executeAction(
              () => updateProxy(record.id, {
                name: record.name,
                description: record.description,
                port: record.port,
                status: newStatus,
              }),
              {
                successMessage: checked ? tr('已启用', 'Enabled') : tr('已停用', 'Disabled'),
                errorMessage: tr('操作失败', 'Operation failed'),
                onSuccess: reload,
              },
            );
          }}
        />
      ),
    },
    {
      title: tr('创建时间', 'Created At'),
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
      search: false,
    },
    {
      title: tr('操作', 'Actions'),
      valueType: 'option',
      width: 180,
      fixed: 'right',
      align: 'center',
      render: (_, record) => {
        const accessUrl = record.access_url;
        const url = accessUrl && typeof accessUrl === 'string'
          ? (accessUrl.startsWith('http://') || accessUrl.startsWith('https://') 
              ? accessUrl 
              : `https://${accessUrl}`)
          : null;
        
        return (
          <Space>
            {url && (
              <Tooltip title={<span style={{ fontSize: '12px' }}>{accessUrl}</span>}>
                <Button
                  type="link"
                  size="small"
                  style={{ padding: 0, height: 'auto' }}
                  onClick={() => {
                    window.open(url, '_blank');
                  }}
                >
                  {tr('去访问', 'Open')}
                </Button>
              </Tooltip>
            )}
            <EditLink onClick={() => {
              setCurrentRow(record);
              setEditModalVisible(true);
            }} />
            <Button
              type="link"
              size="small"
              style={{ padding: 0, height: 'auto' }}
              onClick={() => openFirewallModal(record)}
            >
              {tr('防火墙', 'Firewall')}
            </Button>
            <DeleteLink
              title="确定要删除这个访问吗？"
              onConfirm={() => handleDelete(record.id)}
            />
          </Space>
        );
      },
    },
  ];

  return (
    <PageContainer title={tr('访问', 'Entries')}>
      <div className="table-search-wrapper">
        <ProTable<API.Proxy>
        headerTitle={tr('访问列表', 'Entries')}
        actionRef={actionRef}
        rowKey="id"
        columns={columns}
        request={async (params) => {
          const searchParams = buildSearchParams<API.ProxyListParams>(params, ['name']);
          return tableRequest(() => getProxyList(searchParams), 'proxies');
        }}
        toolBarRender={() => [
          <CreateButton key="create" onClick={() => setCreateModalVisible(true)}>
            {tr('新建访问', 'New Entry')}
          </CreateButton>,
        ]}
        pagination={defaultPagination}
        search={{
          ...defaultSearch,
          labelWidth: 'auto',
        }}
        scroll={{ x: 'max-content' }}
      />
      </div>

      <ModalForm
        key={initialApplicationId ?? 'create'}
        title={tr('新建访问', 'New Entry')}
        open={createModalVisible}
        formRef={createFormRef}
        initialValues={
          initialApplicationId
            ? { application_id: initialApplicationId }
            : undefined
        }
        onOpenChange={(visible) => {
          setCreateModalVisible(visible);
          if (!visible) {
            setInitialApplicationId(undefined);
            setSelectedApplicationId(undefined);
          } else if (initialApplicationId) {
            setSelectedApplicationId(initialApplicationId);
          }
        }}
        onFinish={handleAdd}
        modalProps={{ destroyOnClose: true }}
        width={500}
      >
        <ProFormText
          name="name"
          label={tr('访问名称', 'Entry Name')}
          placeholder={tr('请输入访问名称', 'Please input entry name')}
          rules={[{ required: true, message: tr('请输入访问名称', 'Please input entry name') }]}
        />
        <ProFormSelect
          name="application_id"
          label={tr('关联应用', 'Application')}
          placeholder={tr('请选择要访问的应用', 'Please select an application')}
          rules={[{ required: true, message: tr('请选择应用', 'Please select an application') }]}
          options={applicationOptions}
          fieldProps={{
            onChange: (value: number) => {
              setSelectedApplicationId(value);
            },
          }}
        />
        {selectedApplicationId && applicationMap.get(selectedApplicationId)?.application_type === 'http' && (
          <Alert
            message={<span style={{ fontSize: '11px', lineHeight: '16px', marginBottom: 0, display: 'block' }}>{tr('将开启 HTTPS', 'HTTPS will be enabled')}</span>}
            description={<span style={{ fontSize: '10px', lineHeight: '14px', marginTop: 0, display: 'block' }}>{tr('HTTP 应用将默认使用 HTTPS 协议访问，使用系统配置的 TLS 证书', 'HTTP applications will be exposed over HTTPS with configured TLS certificates')}</span>}
            type="info"
            icon={<CheckCircleOutlined style={{ color: '#52c41a', fontSize: '14px' }} />}
            style={{ marginBottom: 16, padding: '8px 12px' }}
          />
        )}
        <ProFormDigit
          name="port"
          label={tr('公网端口', 'Public Port')}
          placeholder={tr('留空自动分配', 'Leave empty for auto allocation')}
          min={1}
          max={65535}
          extra={tr('映射到公网的端口号，留空则自动分配', 'Mapped public port, empty means auto allocation')}
        />
        <ProFormTextArea
          name="description"
          label={tr('描述', 'Description')}
          placeholder={tr('请输入访问描述', 'Please input description')}
        />
      </ModalForm>

      <ModalForm
        title={tr('编辑访问', 'Edit Entry')}
        open={editModalVisible}
        onOpenChange={setEditModalVisible}
        onFinish={handleEdit}
        modalProps={{ destroyOnClose: true }}
        initialValues={currentRow}
        width={500}
      >
        <ProFormText
          name="name"
          label={tr('访问名称', 'Entry Name')}
          placeholder={tr('请输入访问名称', 'Please input entry name')}
          rules={[{ required: true, message: tr('请输入访问名称', 'Please input entry name') }]}
        />
        <ProFormDigit
          name="port"
          label={tr('公网端口', 'Public Port')}
          placeholder={tr('请输入端口', 'Please input port')}
          min={1}
          max={65535}
        />
        <ProFormTextArea
          name="description"
          label={tr('描述', 'Description')}
          placeholder={tr('请输入访问描述', 'Please input description')}
        />
      </ModalForm>

      <Modal
        title={
          firewallProxy
            ? tr(`防火墙 — ${firewallProxy.name}`, `Firewall — ${firewallProxy.name}`)
            : tr('防火墙', 'Firewall')
        }
        open={firewallModalVisible}
        onCancel={() => setFirewallModalVisible(false)}
        width={560}
        destroyOnClose
        footer={[
          <Popconfirm
            key="reset"
            title={tr('恢复为放行全部？', 'Reset to allow-all?')}
            description={tr(
              '删除规则后，任何来源 IP 都能访问此代理。',
              'After removal, any source IP can reach this proxy.',
            )}
            okText={tr('确认', 'Confirm')}
            cancelText={tr('取消', 'Cancel')}
            onConfirm={handleFirewallReset}
            disabled={!firewallUpdatedAt}
          >
            <Button danger disabled={!firewallUpdatedAt}>
              {tr('恢复默认', 'Reset')}
            </Button>
          </Popconfirm>,
          <Button key="cancel" onClick={() => setFirewallModalVisible(false)}>
            {tr('取消', 'Cancel')}
          </Button>,
          <Button
            key="save"
            type="primary"
            loading={firewallSaving}
            onClick={handleFirewallSave}
          >
            {tr('保存', 'Save')}
          </Button>,
        ]}
      >
        {firewallLoading ? (
          <div style={{ textAlign: 'center', padding: 40 }}>
            <Spin />
          </div>
        ) : (
          <>
            <Alert
              type="info"
              showIcon
              style={{ marginBottom: 16 }}
              message={
                firewallUpdatedAt
                  ? `${tr('规则上次修改：', 'Rule last updated: ')} ${firewallUpdatedAt}`
                  : tr('当前未设置规则 —— 放行全部来源。', 'No rule set — currently allowing all sources.')
              }
              description={
                <div style={{ fontSize: 12 }}>
                  {tr(
                    '填写允许访问此代理的源 IP CIDR。留空保存 = 拒绝全部；点「恢复默认」= 放行全部。HTTP / TCP 代理均在 Accept 时按 L4 过滤。',
                    'Enter source-IP CIDRs allowed to reach this proxy. Save with empty list = deny all; click Reset = allow all. Both HTTP and TCP proxies are filtered at Accept.',
                  )}
                </div>
              }
            />
            <Select
              mode="tags"
              style={{ width: '100%' }}
              value={firewallCidrs}
              onChange={(vs: string[]) => setFirewallCidrs(vs.map((v) => v.trim()).filter(Boolean))}
              placeholder={tr('例如：10.0.0.0/8', 'e.g. 10.0.0.0/8')}
              tokenSeparators={[',', ' ', '\n']}
              open={false}
            />
            <div style={{ marginTop: 8, fontSize: 12, color: 'rgba(0,0,0,0.45)' }}>
              {tr('输入后按回车或逗号分隔。', 'Press Enter or comma to separate.')}
            </div>
          </>
        )}
      </Modal>
    </PageContainer>
  );
};

export default ProxyPage;
