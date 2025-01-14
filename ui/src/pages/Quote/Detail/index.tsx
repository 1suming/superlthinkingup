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

import { useEffect, useState } from 'react';
import { Row, Col } from 'react-bootstrap';
import {
  useParams,
  useSearchParams,
  useNavigate,
  useLocation,
} from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { Pagination, CustomSidebar,CommonRelatedObjects } from '@/components';
import { loggedUserInfoStore, toastStore } from '@/stores';
import { scrollToElementTop, scrollToDocTop } from '@/utils';
import { usePageTags, usePageUsers, useSkeletonControl } from '@/hooks';
import type {
  ListResult,
  QuoteDetailRes,
  AnswerItem,
} from '@/common/interface';
import { quoteDetail, getAnswers, useSimilarQuote } from '@/services';

import {
  Article,
  Answer,
  AnswerHead,
  RelatedArticles,
  WriteAnswer,
  Alert,
  ContentLoader,
  InviteToAnswer,
} from './components';
import { pathFactory } from '@/router/pathFactory';

import './index.scss';

const Index = () => {
  const navigate = useNavigate();
  const { t } = useTranslation('translation');
  const { qid = '', slugPermalink = '' } = useParams();
  /**
   * Note: Compatible with Permalink
   */
  let { aid = '' } = useParams();
  if (!aid && slugPermalink) {
    aid = slugPermalink;
  }

  const [urlSearch] = useSearchParams();
  const page = Number(urlSearch.get('page') || 0);
  const order = urlSearch.get('order') || '';
  const [question, setQuestion] = useState<QuoteDetailRes | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const { isSkeletonShow } = useSkeletonControl(isLoading);
  const [answers, setAnswers] = useState<ListResult<AnswerItem>>({
    count: -1,
    list: [],
  });
  const { setUsers } = usePageUsers();
  const userInfo = loggedUserInfoStore((state) => state.user);
  const isAuthor = userInfo?.username === question?.user_info?.username;
  const isAdmin = userInfo?.role_id === 2;
  const isModerator = userInfo?.role_id === 3;
  const isLogged = Boolean(userInfo?.access_token);
  const loggedUserRank = userInfo?.rank;
  const location = useLocation();

  useEffect(() => {
    if (location.state?.isReview) {
      toastStore.getState().show({
        msg: t('review', { keyPrefix: 'toast' }),
        variant: 'warning',
      });

      // remove state isReview
      const newLocation = { ...location };
      delete newLocation.state;
      window.history.replaceState(null, '', newLocation.pathname);
    }
  }, [location.state]);

//   const requestAnswers = async () => {
//     const res = await getAnswers({
//       order: order === 'updated' || order === 'created' ? order : 'default',
//       question_id: qid,
//       page: 1,
//       page_size: 999,
//     });

//     if (res) {
//       res.list = res.list?.filter((v) => {
//         // delete answers only show to author and admin and has search params aid
//         if (v.status === 10) {
//           if (
//             (v?.user_info?.username === userInfo?.username || isAdmin) &&
//             aid === v.id
//           ) {
//             return v;
//           }
//           return null;
//         }

//         return v;
//       });

//       setAnswers({ ...res, count: res.list.length });
//       if (page > 0 || order) {
//         // scroll into view;
//         const element = document.getElementById('answerHeader');
//         scrollToElementTop(element);
//       }

//       res.list.forEach((item) => {
//         setUsers([
//           {
//             displayName: item.user_info?.display_name,
//             userName: item.user_info?.username,
//           },
//           {
//             displayName: item?.update_user_info?.display_name,
//             userName: item?.update_user_info?.username,
//           },
//         ]);
//       });
//     }
//   };

  const getDetail = async () => {
    console.log("getDetail func start")
    setIsLoading(true);
    try {
      const res = await quoteDetail(qid);
      console.log("api response:",res)
      if (res) {
        setUsers([
          {
            id: res.user_info?.id,
            displayName: res.user_info?.display_name,
            userName: res.user_info?.username,
            avatar_url: res.user_info?.avatar,
          },
          {
            id: res?.update_user_info?.id,
            displayName: res?.update_user_info?.display_name,
            userName: res?.update_user_info?.username,
            avatar_url: res?.update_user_info?.avatar,
          },
          {
            id: res?.last_answered_user_info?.id,
            displayName: res?.last_answered_user_info?.display_name,
            userName: res?.last_answered_user_info?.username,
            avatar_url: res?.last_answered_user_info?.avatar,
          },
        ]);
        setQuestion(res);//@ms:set question
      }
      setIsLoading(false);
    } catch (e) {
      setIsLoading(false);
    }
  };

  const initPage = (type: string) => {
    console.log("initPage func start")
    if (type === 'delete_question') {
      setTimeout(() => {
        navigate('/', { replace: true });
      }, 1000);
      return;
    }
    if (type === 'default') {
      scrollToDocTop();
      console.log("default getDetail")
      getDetail();
      return;
    }
    if (type === 'delete_answer') {
      getDetail();
    }
   //@ms: requestAnswers();
  };

  const writeAnswerCallback = (obj: AnswerItem) => {
    setAnswers({
      count: answers.count + 1,
      list: [...answers.list, obj],
    });

    if (question) {
      setQuestion({
        ...question,
        answered: true,
        first_answer_id: question.first_answer_id
          ? question.first_answer_id
          : obj.id,
      });
    }
  };

  useEffect(() => {
    if (!qid) {
      return;
    }
    console.log("useEffect getDetail")
    getDetail();
    //@ms>:requestAnswers();
  }, [qid]);

//   useEffect(() => {
//     if (page || order) {
//       requestAnswers();
//     }
//   }, [page, order]);
  usePageTags({
    title: question?.title,
    description: question?.description,
    keywords: question?.tags.map((_) => _.slug_name).join(','),
  });

  let showInviteToAnswer = question?.id;
  showInviteToAnswer ="";
  let canInvitePeople = false;
//   if (showInviteToAnswer && Array.isArray(question.extends_actions)) {
//     const inviteAct = question.extends_actions.find((op) => {
//       return op.action === 'invite_other_to_answer';
//     });
//     if (inviteAct) {
//       canInvitePeople = true;
//     }
//   }

const get_similar_quote_func=(id)=>{
    return function(){
      return useSimilarQuote({
        quote_id : id,
        page_size: 5,
      });
    };
};
const gen_link_func_similar_quote= (item: any)=>{
  return pathFactory.articleLanding(item.id, item.url_title);
};



  return (
    <Row className="quoteDetailPage pt-1 mb-5">
      <Col className="page-main flex-auto">
        {question?.operation?.level && <Alert data={question.operation} />}
        {isSkeletonShow ? (
          <ContentLoader />
        ) : (
          <Article
            data={question}
            initPage={initPage}
            hasAnswer={answers.count > 0}
            isLogged={isLogged}
          />
        )}
    
      
      </Col>
      <Col className="page-right-side mt-4 mt-xl-0">
        <CustomSidebar />
       
        <CommonRelatedObjects  
          head_title={t('title', { keyPrefix: 'related_quote' })}
          get_data_func={get_similar_quote_func(question?.id || '')}  
          gen_link_func={gen_link_func_similar_quote}
        
       />
      </Col>
    
    </Row>
  );
};

export default Index;
