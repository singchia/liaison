import React from 'react';
import { Button, Popconfirm } from 'antd';
import { PlusOutlined, ReloadOutlined, EditOutlined } from '@ant-design/icons';

interface RefreshButtonProps {
  onClick: () => void;
  loading?: boolean;
}

/**
 * 刷新按钮
 */
export const RefreshButton: React.FC<RefreshButtonProps> = ({ onClick, loading }) => (
  <Button key="refresh" icon={<ReloadOutlined />} onClick={onClick} loading={loading}>
    刷新
  </Button>
);

interface CreateButtonProps {
  onClick: () => void;
  children?: React.ReactNode;
}

/**
 * 新建按钮
 */
export const CreateButton: React.FC<CreateButtonProps> = ({ onClick, children = '新建' }) => (
  <Button key="create" type="primary" icon={<PlusOutlined />} onClick={onClick}>
    {children}
  </Button>
);

interface EditLinkProps {
  onClick: () => void;
  showIcon?: boolean;
}

/**
 * 编辑链接
 */
export const EditLink: React.FC<EditLinkProps> = ({ onClick, showIcon = false }) => (
  <a onClick={onClick}>
    {showIcon && <EditOutlined />} 编辑
  </a>
);

interface DeleteLinkProps {
  onConfirm: () => void;
  title?: string;
  description?: string;
}

/**
 * 删除链接（带确认）
 */
export const DeleteLink: React.FC<DeleteLinkProps> = ({
  onConfirm,
  title = '确定要删除吗？',
  description,
}) => (
  <Popconfirm title={title} description={description} onConfirm={onConfirm}>
    <a className="text-red-500">删除</a>
  </Popconfirm>
);
