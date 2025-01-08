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
 QuoteList,
//   HotQuestions,
  HotQuotes,
  CustomSidebar,
} from '@/components';
import {
  siteInfoStore,
  loggedUserInfoStore,
  loginSettingStore,
  myGlobalInfoStore,
} from '@/stores';
import { useQuestionList, useQuestionRecommendList ,useQuoteList, useQueryTags} from '@/services';
import * as Type from '@/common/interface';
import { userCenter, floppyNavigation } from '@/utils';
import { QUESTION_ORDER_KEYS } from '@/components/Quote/QuoteList';




const Questions: FC = () => {
  const { t } = useTranslation('translation', { keyPrefix: 'question' });
  const { t: t2 } = useTranslation('translation');
  const { user: loggedUser } = loggedUserInfoStore((_) => _);
  const [urlSearchParams] = useSearchParams();

  const [urlSearchParamsDup,setUrlSearchParamsDup] = useSearchParams();


  const curPage = Number(urlSearchParams.get('page')) || 1;
  const curOrder = (urlSearchParams.get('order') ||
    QUESTION_ORDER_KEYS[0]) as Type.QuestionOrderBy;

const querySelectedTagId=urlSearchParams.get('tag_id') || "";

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
    //   tag_type: TAG_TYPE_ARTICLE,
      is_article_module_menu:1,
    });

    console.log("my tags:",tags)

   
//  const [selectedTagId,setselectedTagId] = useState(querySelectedTagId)
console.log("querySelectedTagId:",querySelectedTagId)
  const reqParams: Type.QueryQuotesReq = {
    page_size: 20,
    page: curPage,
    order: curOrder as Type.QuoteOrderBy,
    
    
    tag_id: querySelectedTagId,//@cws
   // tag: routeParams.tagName,
  };
  const { data: listData, isLoading: listLoading } =
    // curOrder === 'recommend'
    //   ? useQuestionRecommendList(reqParams)
      useQuoteList(reqParams);



  const isIndexPage = useMatch('/');
  let pageTitle = t('articles', { keyPrefix: 'page_title' });
  let slogan = '';
  const { siteInfo } = siteInfoStore();
  if (isIndexPage) {
    pageTitle = `${siteInfo.name}`;
    slogan = `${siteInfo.short_description}`;
  }
  const { login: loginSetting } = loginSettingStore();

  usePageTags({ title: pageTitle, subtitle: slogan });

//参考其他的
 



  
  const {isSideNavSticky,sideNavStickyTop}= myGlobalInfoStore()

 
    const  sideNavStickTopStyle={
        'position':"sticky",
        'top': sideNavStickyTop,
        'zIndex':"1019", //最顶部的top是1020，要比它小，不然会把top的box-shadow挡住

        'background':'#fff',
 
    } as React.CSSProperties ; // 'background':'#fff',
    const emptyStyle={} as React.CSSProperties;

    const handleParams = (selectedTag): string => {
        urlSearchParamsDup.delete('page'); //筛选tag时，删除url中的page，这个很合理
        if(selectedTag===""){
            urlSearchParamsDup.delete('tag_id');
        }else{
            urlSearchParamsDup.set("tag_id", selectedTag);
        }
        
        const searchStr = urlSearchParamsDup.toString();
        // console.log("handleParams:",searchStr)
         
        return `?${searchStr}`;
      };
    const handleTagSelected= (e,tag_id)=>{
        // e.preventDefault();
        // setselectedTagId(tag_id)
        // console.log("click tag_id",tag_id)

        const str = handleParams(tag_id);
        console.log("handleTagSelected str",str);
        // if (floppyNavigation.shouldProcessLinkClick(e)) {
            e.preventDefault();
        //   if (pathname) {
        //     navigate(`${pathname}${str}`);
        //   } else {
            setUrlSearchParamsDup(str); //排查bug浪费了一个小时，不能用 urlSearchParams， 如果用urlSearchParams，,修改了那么urlSearchParams.get('tag_id')就会里面获取到值，就会请求
            //urlSearchParams(str);是
        // }

    }

// <NavLink to="/questions" className="nav-link">全部 </NavLink>
//                 <a className="nav-link" href="#">思维模型</a>
//                 <a className="nav-link" href="#">个人提升</a>
  return (
    <>
     <Row className="flex-fill" style={ sideNavStickTopStyle }>
         <Col  >
            <nav id="second-article-sideNav" className="nav"   >
                
                <a className={  classNames("nav-link",{"active": querySelectedTagId==="" }) } href="/" onClick={ event=> handleTagSelected(event, "")}  >全部</a>
                { tags?.list?.map((tag) => (

                    <a className={  classNames("nav-link",{"active": querySelectedTagId=== tag.tag_id }) } href={handleParams(tag.tag_id)} key={tag.tag_id} data-key={tag.tag_id} onClick={ event=> handleTagSelected(event,tag.tag_id)}>{tag.slug_name}</a>
                ))}
 

            </nav>
        </Col>
       
    </Row>
    <Row className="pt-4 mb-5">
      <Col className="page-main flex-auto">
        <QuoteList
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
        <HotQuotes />
      </Col>
    </Row>
    </>
  );
};

export default Questions;
