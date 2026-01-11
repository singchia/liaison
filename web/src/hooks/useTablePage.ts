import { ActionType } from '@ant-design/pro-components';
import { useRef, useState, useCallback } from 'react';
import { executeAction } from '@/utils/request';

interface UseTablePageOptions<T> {
  /** 删除操作 */
  deleteAction?: (id: number) => Promise<API.Response>;
  /** 删除成功消息 */
  deleteSuccessMessage?: string;
}

/**
 * 通用表格页面 Hook
 * @returns 表格页面常用状态和方法
 */
export function useTablePage<T extends { id: number }>(options?: UseTablePageOptions<T>) {
  const actionRef = useRef<ActionType>();
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [detailVisible, setDetailVisible] = useState(false);
  const [currentRow, setCurrentRow] = useState<T>();
  const [loading, setLoading] = useState(false);

  // 刷新表格
  const reload = useCallback(() => {
    actionRef.current?.reload();
  }, []);

  // 打开新建弹窗
  const openCreateModal = useCallback(() => {
    setCurrentRow(undefined);
    setCreateModalVisible(true);
  }, []);

  // 关闭新建弹窗
  const closeCreateModal = useCallback(() => {
    setCreateModalVisible(false);
  }, []);

  // 打开编辑弹窗
  const openEditModal = useCallback((record: T) => {
    setCurrentRow(record);
    setEditModalVisible(true);
  }, []);

  // 关闭编辑弹窗
  const closeEditModal = useCallback(() => {
    setEditModalVisible(false);
  }, []);

  // 打开详情
  const openDetail = useCallback((record: T) => {
    setCurrentRow(record);
    setDetailVisible(true);
  }, []);

  // 关闭详情
  const closeDetail = useCallback(() => {
    setDetailVisible(false);
  }, []);

  // 删除操作
  const handleDelete = useCallback(
    async (id: number) => {
      if (!options?.deleteAction) return;
      await executeAction(() => options.deleteAction!(id), {
        successMessage: options.deleteSuccessMessage || '删除成功',
        errorMessage: '删除失败',
        onSuccess: reload,
      });
    },
    [options, reload],
  );

  // 新建成功后的回调
  const onCreateSuccess = useCallback(() => {
    closeCreateModal();
    reload();
  }, [closeCreateModal, reload]);

  // 编辑成功后的回调
  const onEditSuccess = useCallback(() => {
    closeEditModal();
    reload();
  }, [closeEditModal, reload]);

  return {
    // refs
    actionRef,
    // 状态
    createModalVisible,
    editModalVisible,
    detailVisible,
    currentRow,
    loading,
    // 状态设置
    setCurrentRow,
    setLoading,
    // 方法
    reload,
    openCreateModal,
    closeCreateModal,
    openEditModal,
    closeEditModal,
    openDetail,
    closeDetail,
    handleDelete,
    onCreateSuccess,
    onEditSuccess,
  };
}
