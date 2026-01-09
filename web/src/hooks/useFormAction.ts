import { useCallback } from 'react';
import { executeAction } from '@/utils/request';

interface FormActionOptions {
  successMessage?: string;
  errorMessage?: string;
  onSuccess?: () => void;
}

/**
 * 创建通用的表单提交处理器
 */
export function useFormAction() {
  /**
   * 创建新建处理器
   */
  const createAddHandler = useCallback(
    <T,>(
      action: (values: T) => Promise<API.Response>,
      options: FormActionOptions,
    ) => {
      return async (values: T) => {
        return executeAction(() => action(values), {
          successMessage: options.successMessage || '创建成功',
          errorMessage: options.errorMessage || '创建失败',
          onSuccess: options.onSuccess,
        });
      };
    },
    [],
  );

  /**
   * 创建编辑处理器
   */
  const createEditHandler = useCallback(
    <T,>(
      action: (id: number, values: T) => Promise<API.Response>,
      id: number | undefined,
      options: FormActionOptions,
    ) => {
      return async (values: T) => {
        if (!id) return false;
        return executeAction(() => action(id, values), {
          successMessage: options.successMessage || '更新成功',
          errorMessage: options.errorMessage || '更新失败',
          onSuccess: options.onSuccess,
        });
      };
    },
    [],
  );

  /**
   * 创建删除处理器
   */
  const createDeleteHandler = useCallback(
    (
      action: (id: number) => Promise<API.Response>,
      options: FormActionOptions,
    ) => {
      return async (id: number) => {
        return executeAction(() => action(id), {
          successMessage: options.successMessage || '删除成功',
          errorMessage: options.errorMessage || '删除失败',
          onSuccess: options.onSuccess,
        });
      };
    },
    [],
  );

  return {
    createAddHandler,
    createEditHandler,
    createDeleteHandler,
  };
}
