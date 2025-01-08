
import qs from 'qs';
import useSWR from 'swr';

import request from '@/utils/request';
import type * as Type from '@/common/interface';

export const quoteAuthorDetail = (id: string) => {
    // return request.get<Type.QuoteAuthorDetailRes>(
    //     `/answer/api/v1/question/info?id=${id}`,
    //     { allow404: true },
    //   );

    return request.get<Type.QuoteAuthorDetailRes>(
      `/answer/api/v1/quote/author/info?id=${id}`,
      { allow404: true },
    );
  };
  


  export const saveQuoteAuthor = (params: Type.QuoteAuthorParams) => {
    return request.post('/answer/api/v1/quote/author', params);
  };
  
  export const modifyQuoteAuthor = (
    params: Type.QuoteAuthorParams & { id: string; edit_summary: string },
  ) => {
    return request.put(`/answer/api/v1/quote`, params);
  };

  export const deleteQuoteAuthor = (params: {
    id: string;
    captcha_code?: string;
    captcha_id?: string;
  }) => {
    return request.delete('/answer/api/v1/quote/author', params);
  };

  export const unDeleteQuoteAuthor = (qid) => {
    return request.post('/answer/api/v1/quote/author/recover', {
      question_id: qid,
    });
  };

  export const quoteAuthorOperation = (params: Type.QuoteAuthorOperationReq) => {
    return request.put('/answer/api/v1/quote/author/operation', params);
  };

  export const useQuoteAuthorList = (params: Type.QueryQuoteAuthorsReq) => {
    const apiUrl = `/answer/api/v1/quote/author/page?${qs.stringify(params)}`;
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

export const queryQuoteAuthorByTitle = (title: string) => {
    return request.get(`/answer/api/v1/quote/author/similar?title=${title}`);
};
    
