
import qs from 'qs';
import useSWR from 'swr';

import request from '@/utils/request';
import type * as Type from '@/common/interface';

export const articleDetail = (id: string) => {
    // return request.get<Type.ArticleDetailRes>(
    //     `/answer/api/v1/question/info?id=${id}`,
    //     { allow404: true },
    //   );

    return request.get<Type.ArticleDetailRes>(
      `/answer/api/v1/article/info?id=${id}`,
      { allow404: true },
    );
  };
  

