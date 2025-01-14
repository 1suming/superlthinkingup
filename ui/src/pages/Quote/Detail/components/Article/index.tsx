/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import { memo, FC, useState, useEffect, useRef } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Button, OverlayTrigger, Tooltip } from 'react-bootstrap';

import {
  Tag,
//  Actions,
 // Operate,
//  UserCard,
//  Comment,
  FormatTime,
  htmlRender,
  Icon,
  ImgViewer,

  ArticleActions,
  ArticleOperate,
  ArticleComment,
  ArticleUserCard,
  BaseUserCardMy,

} from '@/components';


import { useRenderHtmlPlugin } from '@/utils/pluginKit';
import { formatCount, guard } from '@/utils';
import { following } from '@/services';
import { pathFactory } from '@/router/pathFactory';


interface Props {
  data: any;
  hasAnswer: boolean;
  isLogged: boolean;
  initPage: (type: string) => void;
}

const Index: FC<Props> = ({ data, initPage, hasAnswer, isLogged }) => {
    console.log("article data:",data)
  const { t } = useTranslation('translation', {
    // keyPrefix: 'question_detail',
    keyPrefix: 'article_detail',
  });
  const [searchParams] = useSearchParams();
  const [followed, setFollowed] = useState(data?.is_followed);
  const ref = useRef<HTMLDivElement>(null);

  useRenderHtmlPlugin(ref.current);

  const handleFollow = (e) => {
    e.preventDefault();
    if (!guard.tryNormalLogged(true)) {
      return;
    }
    following({
      object_id: data?.id,
      is_cancel: followed,
    }).then((res) => {
      setFollowed(res.is_followed);
    });
  };

  useEffect(() => {
    if (data) {
      setFollowed(data?.is_followed);
    }
  }, [data]);

  useEffect(() => {
    if (!ref.current) {
      return;
    }

    htmlRender(ref.current);
  }, [ref.current]);

  if (!data?.id) {
    return null;
  }

  return (
    <div className="quote-all-wrapper px-3">
   

      <div className="d-flex flex-wrap align-items-center small mb-3 text-secondary">
        <FormatTime
          time={data.create_time}
          preFix={t('created')}
          className="me-3"
        />

        <FormatTime
          time={data.update_time}
          preFix={t('updated')}
          className="me-3"
        />
        {data?.view_count > 0 && (
          <div className="me-3">
            {t('Views')} {formatCount(data.view_count)}
          </div>
        )}
        
      </div>
      <div className="m-n1">
        {data?.tags?.map((item: any) => {
          return <Tag className="m-1" key={item.slug_name} data={item}  tagType='article'/>;
        })}
      </div>
      <ImgViewer>
        <article
          ref={ref}
          className="fmt text-break text-wrap mt-4 quote-body-content "
          dangerouslySetInnerHTML={{ __html: data?.html }}
        />
          <div className="d-flex justify-content-end text-pink-500 text-lg">
        
              <BaseUserCardMy
                      
                      display_name ={data.quote_author_basic_info.author_name }
                      link_to={`/authors/${data.quote_author_basic_info.id}`}
                      avatar={data.quote_author_basic_info.avatar}
                      
                      avatarClass="me-2 d-block"
                      avatarSearchStr=''

                     
                    />
                    - 
                    <div> 
                      <Link to={`/pieces/${data.quote_piece_basic_info.id}`} >
                        &laquo;{data.quote_piece_basic_info.title}&raquo;
                     </Link>
                    </div>
                    
          </div>
      </ImgViewer>

      <ArticleActions
        className="mt-4"
        source="question"
        data={{
          id: data?.id,
          isHate: data?.vote_status === 'vote_down',
          isLike: data?.vote_status === 'vote_up',
          votesCount: data?.vote_count,
          collected: data?.collected,
          collectCount: data?.collection_count,
          username: data.user_info?.username,
        }}
      />


      <div className="d-block d-md-flex flex-wrap quote-usercard">
         <div className="mb-3 mb-md-0 me-4 flex-grow-1">
          <ArticleOperate
            qid={data?.id}
            type="article"
            memberActions={data?.member_actions}
            title={data.title}
            hasAnswer={hasAnswer}
            isAccepted={Boolean(data?.accepted_answer_id)}
            callback={initPage}
          />
        </div>


         <ArticleUserCard
            data={data?.user_info}
            time={data.create_time}
            preFix={t('asked')}
            isLogged={isLogged}
            timelinePath={`/posts/${data.id}/timeline`}
          />

      </div>

      <ArticleComment
        objectId={data?.id}
        mode="question"
        commentId={searchParams.get('commentId')}
      />
    </div>
  );
};

export default memo(Index);
