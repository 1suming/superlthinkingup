
import qs from 'qs';
import useSWR from 'swr';

import request from '@/utils/request';
import type * as Type from '@/common/interface';

export const quoteDetail = (id: string) => {
    // return request.get<Type.QuoteDetailRes>(
    //     `/answer/api/v1/question/info?id=${id}`,
    //     { allow404: true },
    //   );

    return request.get<Type.QuoteDetailRes>(
      `/answer/api/v1/quote/info?id=${id}`,
      { allow404: true },
    );
  };
  


  export const saveQuote = (params: Type.QuoteParams) => {
    return request.post('/answer/api/v1/quote', params);
  };
  
  export const modifyQuote = (
    params: Type.QuoteParams & { id: string; edit_summary: string },
  ) => {
    return request.put(`/answer/api/v1/quote`, params);
  };

  export const deleteQuote = (params: {
    id: string;
    captcha_code?: string;
    captcha_id?: string;
  }) => {
    return request.delete('/answer/api/v1/quote', params);
  };

  export const unDeleteQuote = (qid) => {
    return request.post('/answer/api/v1/quote/recover', {
      question_id: qid,
    });
  };

  export const quoteOperation = (params: Type.QuoteOperationReq) => {
    return request.put('/answer/api/v1/quote/operation', params);
  };

  export const useQuoteList = (params: Type.QueryQuotesReq) => {
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

export const queryQuoteByTitle = (title: string) => {
    return request.get(`/answer/api/v1/quote/similar?title=${title}`);
};
    
