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

import { FC,useState } from 'react';
import { Row, Col ,Nav} from 'react-bootstrap';
import { useMatch, Link, useSearchParams,NavLink } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { usePageTags } from '@/hooks';


import classNames from 'classnames';

import {
  FollowingTags,
 // QuestionList,
 ArticleList,
//   HotQuestions,
  HotArticles,
  CustomSidebar,
} from '@/components';
import {
  siteInfoStore,
  loggedUserInfoStore,
  loginSettingStore,
  myGlobalInfoStore,
} from '@/stores';
import { useQuestionList, useQuestionRecommendList ,useArticleList, useQueryTags} from '@/services';
import * as Type from '@/common/interface';
import { userCenter, floppyNavigation } from '@/utils';
import { QUESTION_ORDER_KEYS } from '@/components/Article/ArticleList';




const Questions: FC = () => {
  const { t } = useTranslation('translation', { keyPrefix: 'question' });
  const { t: t2 } = useTranslation('translation');
  const { user: loggedUser } = loggedUserInfoStore((_) => _);
  const [urlSearchParams] = useSearchParams();
  const curPage = Number(urlSearchParams.get('page')) || 1;
  const curOrder = (urlSearchParams.get('order') ||
    QUESTION_ORDER_KEYS[0]) as Type.QuestionOrderBy;

    const TAG_TYPE_ARTICLE=1;

    //---tag
    const tag_page=0
    const tag_pageSize = 20;
    const {
        data: tags,
        mutate,
        isLoading,
        
    } = useQueryTags({
        tag_page,
      page_size: tag_pageSize,
      tag_type: TAG_TYPE_ARTICLE,
     
    });

    console.log("my tags:",tags)

   
 const [selectedTagId,setselectedTagId] = useState("")

  const reqParams: Type.QueryArticlesReq = {
    page_size: 20,
    page: curPage,
    order: curOrder as Type.ArticleOrderBy,
    
    
    tag_id: selectedTagId,//@cws
   // tag: routeParams.tagName,
  };
  const { data: listData, isLoading: listLoading } =
    // curOrder === 'recommend'
    //   ? useQuestionRecommendList(reqParams)
      useArticleList(reqParams);



  const isIndexPage = useMatch('/');
  let pageTitle = t('questions', { keyPrefix: 'page_title' });
  let slogan = '';
  const { siteInfo } = siteInfoStore();
  if (isIndexPage) {
    pageTitle = `${siteInfo.name}`;
    slogan = `${siteInfo.short_description}`;
  }
  const { login: loginSetting } = loginSettingStore();

  usePageTags({ title: pageTitle, subtitle: slogan });



  
  const {isSideNavSticky,sideNavStickyTop}= myGlobalInfoStore()

 
    const  sideNavStickTopStyle={
        'position':"sticky",
        'top': sideNavStickyTop,
        'zIndex':"1019", //最顶部的top是1020，要比它小，不然会把top的box-shadow挡住

        'background':'#fff',
 
    } as React.CSSProperties ; // 'background':'#fff',
    const emptyStyle={} as React.CSSProperties;


    const handleTagSelected= (e,tag_id)=>{
        e.preventDefault();
        setselectedTagId(tag_id)
        console.log("click tag_id",tag_id)
    }

// <NavLink to="/questions" className="nav-link">全部 </NavLink>
//                 <a className="nav-link" href="#">思维模型</a>
//                 <a className="nav-link" href="#">个人提升</a>
  return (
    <>
     <Row className="flex-fill" style={ sideNavStickTopStyle }>
         <Col  >
            <nav id="second-article-sideNav" className="nav"   >
                
                <a className={  classNames("nav-link",{"active": selectedTagId==="" }) } href="#" onClick={ event=> handleTagSelected(event, "")}  >全部</a>
                { tags?.list?.map((tag) => (

                    <a className={  classNames("nav-link",{"active": selectedTagId=== tag.tag_id }) } href="#" key={tag.tag_id} data-key={tag.tag_id} onClick={ event=> handleTagSelected(event,tag.tag_id)}>{tag.slug_name}</a>
                ))}

            </nav>
        </Col>
       
    </Row>
    <Row className="pt-4 mb-5">
      <Col className="page-main flex-auto">
        <ArticleList
          source="questions"
          data={listData}
          order={curOrder}
          orderList={
            loggedUser.username
              ? QUESTION_ORDER_KEYS
              : QUESTION_ORDER_KEYS.filter((key) => key !== 'recommend')
          }
          isLoading={listLoading}
        />
      </Col>
      <Col className="page-right-side mt-4 mt-xl-0">
        <CustomSidebar />
        {!loggedUser.username && (
          <div className="card mb-4">
            <div className="card-body">
              <h5 className="card-title">
                {t2('website_welcome', {
                  site_name: siteInfo.name,
                })}
              </h5>
              <p className="card-text">{siteInfo.description}</p>
              <Link
                to={userCenter.getLoginUrl()}
                className="btn btn-primary"
                onClick={floppyNavigation.handleRouteLinkClick}>
                {t('login', { keyPrefix: 'btns' })}
              </Link>
              {loginSetting.allow_new_registrations ? (
                <Link
                  to={userCenter.getSignUpUrl()}
                  className="btn btn-link ms-2"
                  onClick={floppyNavigation.handleRouteLinkClick}>
                  {t('signup', { keyPrefix: 'btns' })}
                </Link>
              ) : null}
            </div>
          </div>
        )}
        {/* {loggedUser.access_token && <FollowingTags />} */}
        <HotArticles />
      </Col>
    </Row>
    </>
  );
};

export default Questions;
