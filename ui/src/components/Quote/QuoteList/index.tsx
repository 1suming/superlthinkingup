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

import { FC } from 'react';
import { ListGroup } from 'react-bootstrap';
import { Link, NavLink, useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { pathFactory } from '@/router/pathFactory';
import {
  Tag,
  Pagination,
  FormatTime,
  Empty,
  BaseUserCard,
  BaseUserCardMy,
  QueryGroup,
  QuestionListLoader,
  ArticleCounts,
  Icon,
} from '@/components';



import * as Type from '@/common/interface';
import { useSkeletonControl } from '@/hooks';

export const QUESTION_ORDER_KEYS: Type.QuestionOrderBy[] = [
  'newest',
  'active',
  'hot',
  'score',
  'unanswered',
  'recommend',
];
interface Props {
  source: 'articles'|'questions' | 'tag';
  order?: Type.QuestionOrderBy;
  data;
  orderList?: Type.QuestionOrderBy[];
  isLoading: boolean;
}

const QuoteList: FC<Props> = ({
  source,
  order,
  data,
  orderList,
  isLoading = false,
}) => {
  const { t } = useTranslation('translation', { keyPrefix: 'article' });
  const [urlSearchParams] = useSearchParams();
  const { isSkeletonShow } = useSkeletonControl(isLoading);
  const curOrder =
    order || urlSearchParams.get('order') || QUESTION_ORDER_KEYS[0];
  const curPage = Number(urlSearchParams.get('page')) || 1;
  const pageSize = 20;
  const count = data?.count || 0;
  const orderKeys = orderList || QUESTION_ORDER_KEYS;

  console.log("quoteList data:",data)

  function genQuoteContent(li: any){
    // const len = li.thumbnails?.length
    // console.log("thumbnails len:",len)
    const quote_piece_basic_info=li.quote_piece_basic_info
    console.log("pathFactory.quoteLanding(li.id, li.url_title)",pathFactory.quoteLanding(li.id, li.url_title))
    if(quote_piece_basic_info.avatar==""){
        return (
           
              <div className="article-img-content-wrapper">
                    <h2 className="text-wrap text-break article-title">
                    {li.pin === 2 && (
                        <Icon
                        name="pin-fill"
                        className="me-1"
                        title={t('pinned', { keyPrefix: 'btns' })}
                        />
                    )}
                    {/* <NavLink
                        to={pathFactory.quoteLanding(li.id, li.url_title)}
                        className="link-dark">
                        {li.title}
                        {li.status === 2 ? ` [${t('closed')}]` : ''}
                    </NavLink> */}
                    </h2>
                  
                    <Link to={pathFactory.quoteLanding(li.id, li.url_title)}> 
                     <p className="quote-decription"  dangerouslySetInnerHTML={{ __html: li.description }} />
                     </Link>
                </div>
           
        )
    }

        const firstImg =quote_piece_basic_info.avatar ;// li.thumbnails[0]
        return (
            <div className="article-img-content-wrapper article-has-img">
                <div className="img-wrapper">
                    <a href={pathFactory.quoteLanding(li.id, li.url_title)} target="_blank" rel="noreferrer"> 
                        <img
                                src={firstImg}
                                width={216}
                                height="144"
                                className=""
                                alt=""
                    
                        />
                    </a> 
                </div>
                <div className="quote-content-wrapper">
                    <h2 className="text-wrap text-break article-title">
                        {li.pin === 2 && (
                            <Icon
                            name="pin-fill"
                            className="me-1"
                            title={t('pinned', { keyPrefix: 'btns' })}
                            />
                        )}
                        <NavLink
                            to={pathFactory.quoteLanding(li.id, li.url_title)}
                            target="_blank" rel="noreferrer"
                            className="link-dark">
                            {li.title}
                            {li.status === 2 ? ` [${t('closed')}]` : ''}
                        </NavLink>
                        </h2>
                        <p className="quote-decription" dangerouslySetInnerHTML={{ __html: li.description }} />
                </div>
            </div>


           
        )

   

  }
  return (
    <div>
      {/* <div className="mb-3 d-flex flex-wrap  justify-content-end">
        { <h5 className="fs-5 text-nowrap mb-3 mb-md-0">
          {source === 'questions'
            ? t('all_questions')
            : t('x_questions', { count })}
        </h5> }
        <QueryGroup
          data={orderKeys}
          currentSort={curOrder}
          pathname={source === 'questions' ? '/questions' : ''}
          i18nKeyPrefix="article"
        />
      </div> */}
      <ListGroup className="rounded-0">
        {isSkeletonShow ? (
          <QuestionListLoader />
        ) : (
          data?.list?.map((li) => {
            return (
              <ListGroup.Item
                key={li.id}
                className=" py-3 px-3 border-start-0 border-end-0  quote-item ">
                <div className="quote-content">
                  { genQuoteContent(li) }
                </div>
                <div className="d-flex">
                  <span>出自:</span>
                   <BaseUserCardMy
                      
                      display_name ={li.quote_author_basic_info.author_name }
                      link_to={`/authors/${li.quote_author_basic_info.id}`}
                      avatar={li.quote_author_basic_info.avatar}
                      showAvatar={false}
                      className="me-1"
                    />
                    - 
                    <div> 
                      <Link to={`/pieces/${li.quote_piece_basic_info.id}`} >
                        &laquo;{li.quote_piece_basic_info.title}&raquo;
                     </Link>
                    </div>
                </div>
                
                <div className="d-flex flex-wrap flex-column flex-md-row align-items-md-center small mb-2 text-secondary">
                  <div className="d-flex flex-wrap me-0 me-md-3" >
                    <BaseUserCard
                      data={li.operator}
                      showAvatar={false}
                      showReputation={false}
                      className="me-1"
                    />
                    •
                    <FormatTime
                      time={li.operated_at}
                      className="text-secondary ms-1 flex-shrink-0"
                      preFix={t(li.operation_type)}
                    />
                  </div>
                  <ArticleCounts
                    data={{
                      votes: li.vote_count,
                      answers: li.answer_count,
                      views: li.view_count,
                    }}
                    isAccepted={li.accepted_answer_id >= 1}
                    className="mt-2 mt-md-0"
                  />
                </div>
                {/* <div className="question-tags m-n1">
                  {Array.isArray(li.tags)
                    ? li.tags.map((tag) => {
                        return (
                          <Tag key={tag.slug_name} className="m-1" data={tag} tagType='article' />
                        );
                      })
                    : null}
                </div> */}
              </ListGroup.Item>
            );
          })
        )}
      </ListGroup>
      {count <= 0 && !isLoading && <Empty />}
      <div className="mt-4 mb-2 d-flex justify-content-center">
        <Pagination
          currentPage={curPage}
          totalSize={count}
          pageSize={pageSize}
          pathname={source === 'articles' ? '/articles' : ''}
        />
      </div>
    </div>
  );
};
export default QuoteList;
