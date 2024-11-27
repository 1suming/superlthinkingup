import useSWR from 'swr';
import qs from 'qs';

import request from '@/utils/request';
import type * as Type from '@/common/interface';


export const useHotArticles = (
    params: Type.QueryArticlesReq = {
      page: 1,
      page_size: 6,
      order: 'hot',
      in_days: 7,
    },
  ) => {
    const apiUrl = `/answer/api/v1/article/page?${qs.stringify(params)}`;
    const { data, error } = useSWR<Type.ListResult, Error>(
      [apiUrl],
      request.instance.get,
    );
    return {
      data,
      isLoading: !data && !error,
      error,
    };
  };

  export const useSimilarArticle = (params: {
    article_id: string;
    page_size: number;
  }) => {
    const apiUrl = `/answer/api/v1/article/similar/tag?${qs.stringify(params)}`;
  
    const { data, error } = useSWR<Type.ListResult, Error>(
      params.article_id ? apiUrl : null,
      request.instance.get,
    );
    return {
      data,
      isLoading: !data && !error,
      error,
    };
  };
  