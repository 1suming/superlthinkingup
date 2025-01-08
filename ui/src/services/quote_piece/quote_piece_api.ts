
import qs from 'qs';
import useSWR from 'swr';

import request from '@/utils/request';
import type * as Type from '@/common/interface';

export const quotePieceDetail = (id: string) => {
    // return request.get<Type.QuotePieceDetailRes>(
    //     `/answer/api/v1/question/info?id=${id}`,
    //     { allow404: true },
    //   );

    return request.get<Type.QuotePieceDetailRes>(
      `/answer/api/v1/quote/piece/info?id=${id}`,
      { allow404: true },
    );
  };
  


  export const saveQuotePiece = (params: Type.QuotePieceParams) => {
    return request.post('/answer/api/v1/quote/piece', params);
  };
  
  export const modifyQuotePiece = (
    params: Type.QuotePieceParams & { id: string; edit_summary: string },
  ) => {
    return request.put(`/answer/api/v1/quote`, params);
  };

  export const deleteQuotePiece = (params: {
    id: string;
    captcha_code?: string;
    captcha_id?: string;
  }) => {
    return request.delete('/answer/api/v1/quote/piece', params);
  };

  export const unDeleteQuotePiece = (qid) => {
    return request.post('/answer/api/v1/quote/piece/recover', {
      question_id: qid,
    });
  };

  export const quotePieceOperation = (params: Type.QuotePieceOperationReq) => {
    return request.put('/answer/api/v1/quote/piece/operation', params);
  };

  export const useQuotePieceList = (params: Type.QueryQuotePiecesReq) => {
    const apiUrl = `/answer/api/v1/quote/piece/page?${qs.stringify(params)}`;
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

export const queryQuotePieceByTitle = (title: string) => {
    return request.get(`/answer/api/v1/quote/piece/similar?title=${title}`);
};
    
