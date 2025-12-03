import { useSearchParams } from '@umijs/max';
import { useState } from 'react';

interface UseQueryParamsOptions<P> {
  // 默认参数
  defaultParams?: Partial<P>;
  // 参数转换函数（可选，用于处理参数格式）
  transformParams?: (params: any) => P;
}

/**
 * URL Query 参数管理 Hook
 * @param options 配置项
 * @returns 返回当前参数、更新函数和重置函数
 */
export function useQueryParams<P = any>({
  defaultParams = {},
  transformParams,
}: UseQueryParamsOptions<P> = {}) {
  const [searchParams, setSearchParams] = useSearchParams();
  const [params, setParams] = useState<P>(() => {
    // 从 URL query 初始化参数
    const queryParams: any = {};
    searchParams.forEach((value, key) => {
      // 尝试转换数字类型
      if (!isNaN(Number(value))) {
        queryParams[key] = Number(value);
      } else {
        queryParams[key] = value;
      }
    });

    // 合并默认参数和 query 参数
    const initialParams = {
      page: 1,
      page_size: 10,
      ...defaultParams,
      ...queryParams,
    } as P;

    return transformParams ? transformParams(initialParams) : initialParams;
  });

  // 更新搜索参数，同时同步到 URL query
  const updateParams = (newParams: Partial<P>) => {
    const mergedParams = { ...params, ...newParams } as P;
    setParams(mergedParams);

    // 更新 URL query
    const queryObj: any = {};
    Object.entries(mergedParams as any).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        queryObj[key] = String(value);
      }
    });
    setSearchParams(queryObj, { replace: true });
  };

  // 重置参数
  const resetParams = () => {
    const resetedParams = {
      page: 1,
      page_size: 10,
      ...defaultParams,
    } as P;
    setParams(resetedParams);
    setSearchParams({}, { replace: true });
  };

  return {
    params,
    updateParams,
    resetParams,
  };
}
