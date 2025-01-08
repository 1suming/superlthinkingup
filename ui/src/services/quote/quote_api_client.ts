import useSWR from 'swr';
import qs from 'qs';

import request from '@/utils/request';
import type * as Type from '@/common/interface';


export const useHotQuotes = (
    params: Type.QueryQuotesReq = {
      page: 1,
      page_size: 6,
      order: 'hot',
      in_days: 7,
    },
  ) => {
    const apiUrl = `/answer/api/v1/quote/page?${qs.stringify(params)}`;
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

  export const useSimilarQuote = (params: {
    quote_id: string;
    page_size: number;
  }) => {
    const apiUrl = `/answer/api/v1/quote/similar/tag?${qs.stringify(params)}`;
  
    const { data, error } = useSWR<Type.ListResult, Error>(
      params.quote_id ? apiUrl : null,
      request.instance.get,
    );
    return {
      data,
      isLoading: !data && !error,
      error,
    };
  };
  