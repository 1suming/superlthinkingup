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

import React, { FormEvent, useState, useEffect } from 'react';
import { Container, Form, Button, Col } from 'react-bootstrap';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { Trans, useTranslation } from 'react-i18next';

import { usePageTags } from '@/hooks';
// import type { LoginReqParams, FormDataType ,QuerySiteInfoKeyValReq,QuerySiteInfoKeyValResp} from '@/common/interface';
import * as Type from '@/common/interface';

import { Unactivate, WelcomeTitle, PluginRender } from '@/components';
import {
  loggedUserInfoStore,
  loginSettingStore,
  userCenterStore,
} from '@/stores';
import {
  floppyNavigation,
  guard,
  handleFormError,
  userCenter,
  scrollToElementTop,
} from '@/utils';
import { PluginType, useCaptchaPlugin } from '@/utils/pluginKit';
import { login, UcAgent ,useSiteInfoValByKey} from '@/services';
import { setupAppTheme } from '@/utils/localize';

const Index: React.FC = () => {
  const { t } = useTranslation('translation', { keyPrefix: 'login' });
  usePageTags({
    title: t('site_about', { keyPrefix: 'page_title' }),
  });

  const reqParams: Type.QuerySiteInfoKeyValReq = {
     key: "site_about_info",
  };
  const { data,error  } =  useSiteInfoValByKey(reqParams);

 
//   useEffect(() => {
   
//   }, []);

 

  return (
    <Container style={{ paddingTop: '4rem', paddingBottom: '5rem' }}>
        
      
        <Col className="mx-auto text-center" md={12} lg={12} xl={12}>
          <h3>关于超维社</h3>
          <article  className="fmt text-break text-wrap mt-4 " dangerouslySetInnerHTML={{ __html: (data?.content)?(data?.content):"" }}   />
 
        </Col>
        <div className="row mt-5">
          <div className="col-md-4 mb-4">
            <div className="text-center p-4 h-100" style={{ 
              background: 'linear-gradient(135deg, #6B8DE3 0%, #5E72EB 100%)',
              borderRadius: '20px',
              boxShadow: '0 4px 15px rgba(94,114,235,0.2)',
              transition: 'all 0.3s ease',
              cursor: 'pointer',
              border: '1px solid rgba(255,255,255,0.1)',
              color: 'white'
            }}
            onMouseOver={(e) => {
              e.currentTarget.style.transform = 'translateY(-5px)';
              e.currentTarget.style.boxShadow = '0 8px 25px rgba(94,114,235,0.3)';
            }}
            onFocus={(e) => {
              e.currentTarget.style.transform = 'translateY(-5px)';
              e.currentTarget.style.boxShadow = '0 8px 25px rgba(94,114,235,0.3)';
            }}
            onMouseOut={(e) => {
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = '0 4px 15px rgba(94,114,235,0.2)';
            }}
            onBlur={(e) => {
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = '0 4px 15px rgba(94,114,235,0.2)';
            }}>
              <div className="mb-3">
                <i className="bi bi-chat-dots" style={{ 
                  fontSize: '2.5rem', 
                  color: 'white',
                  filter: 'drop-shadow(0 2px 4px rgba(0,0,0,0.2))'
                }}></i>
              </div>
              <h4 className="mb-3 fw-bold">开放讨论</h4>
              <p style={{ color: 'rgba(255,255,255,0.9)' }}>
                提供一个开放、包容的讨论环境，让每个人都能自由表达观点，促进思维的碰撞与创新。
              </p>
            </div>
          </div>

          <div className="col-md-4 mb-4">
            <div className="text-center p-4 h-100" style={{ 
              background: 'linear-gradient(135deg, #FF6B6B 0%, #FF8E8E 100%)',
              borderRadius: '20px',
              boxShadow: '0 4px 15px rgba(255,107,107,0.2)',
              transition: 'all 0.3s ease',
              cursor: 'pointer',
              border: '1px solid rgba(255,255,255,0.1)',
              color: 'white'
            }}
            onMouseOver={(e) => {
              e.currentTarget.style.transform = 'translateY(-5px)';
              e.currentTarget.style.boxShadow = '0 8px 25px rgba(255,107,107,0.3)';
            }}
            onFocus={(e) => {
              e.currentTarget.style.transform = 'translateY(-5px)';
              e.currentTarget.style.boxShadow = '0 8px 25px rgba(255,107,107,0.3)';
            }}
            onMouseOut={(e) => {
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = '0 4px 15px rgba(255,107,107,0.2)';
            }}
            onBlur={(e) => {
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = '0 4px 15px rgba(255,107,107,0.2)';
            }}>
              <div className="mb-3">
                <i className="bi bi-book" style={{ 
                  fontSize: '2.5rem', 
                  color: 'white',
                  filter: 'drop-shadow(0 2px 4px rgba(0,0,0,0.2))'
                }}></i>
              </div>
              <h4 className="mb-3 fw-bold">知识分享</h4>
              <p style={{ color: 'rgba(255,255,255,0.9)' }}>
                鼓励用户分享专业知识和经验，构建高质量的知识库，助力彼此成长。
              </p>
            </div>
          </div>

          <div className="col-md-4 mb-4">
            <div className="text-center p-4 h-100" style={{ 
              background: 'linear-gradient(135deg, #4CAF50 0%, #45B649 100%)',
              borderRadius: '20px',
              boxShadow: '0 4px 15px rgba(76,175,80,0.2)',
              transition: 'all 0.3s ease',
              cursor: 'pointer',
              border: '1px solid rgba(255,255,255,0.1)',
              color: 'white'
            }}
            onMouseOver={(e) => {
              e.currentTarget.style.transform = 'translateY(-5px)';
              e.currentTarget.style.boxShadow = '0 8px 25px rgba(76,175,80,0.3)';
            }}
            onFocus={(e) => {
              e.currentTarget.style.transform = 'translateY(-5px)';
              e.currentTarget.style.boxShadow = '0 8px 25px rgba(76,175,80,0.3)';
            }}
            onMouseOut={(e) => {
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = '0 4px 15px rgba(76,175,80,0.2)';
            }}
            onBlur={(e) => {
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = '0 4px 15px rgba(76,175,80,0.2)';
            }}>
              <div className="mb-3">
                <i className="bi bi-people" style={{ 
                  fontSize: '2.5rem', 
                  color: 'white',
                  filter: 'drop-shadow(0 2px 4px rgba(0,0,0,0.2))'
                }}></i>
              </div>
              <h4 className="mb-3 fw-bold">社区互动</h4>
              <p style={{ color: 'rgba(255,255,255,0.9)' }}>
                打造积极向上的社区氛围，通过互动交流建立有价值的社交网络。
              </p>
            </div>
          </div>
        </div>

      
    </Container>
  );
};

export default React.memo(Index);
