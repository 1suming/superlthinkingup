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
import { Container, Form, Button, Col,Row } from 'react-bootstrap';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { Trans, useTranslation } from 'react-i18next';

import { usePageTags } from '@/hooks';
// import type { LoginReqParams, FormDataType ,QuerySiteInfoKeyValReq,QuerySiteInfoKeyValResp} from '@/common/interface';
import * as Type from '@/common/interface';

import "./AIIndex.scss"

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
        
       <Row>
            <Col className="mx-auto" md={12} lg={12} xl={12}>
            <div className=" text-center" >
            
            <h2 className="text-4xl">人工智能AI & chatGPT  </h2>
            <p style={{color:'rgb(71,85,105)'}}>一站式人工智能知识和工具分享平台</p>
            </div>
            
 
             </Col>
        </Row>
        <Row className="mx-auto">
            <Col  md={4} lg={4} xl={4}>
                    <div className="ai-list-item ">
                            <span className="title">chatGPT指令大全</span>
                            <p>集合数百个精炼过的指令，让你发挥 ChatGPT 的强大功能</p>
                    </div>
            
            </Col>
            <Col  md={4} lg={4} xl={4}>
                    <div className="ai-list-item">
                        <span className="title">ChatGPT 应用与教学</span>
                        <p>完整涵盖基础使用、实战案例、串接教学，让你快速上手</p>
                    </div>
            </Col>
            <Col  md={4} lg={4} xl={4}>
                    <div className="ai-list-item">
                        <span className="title">AI 资源 & 产业洞察</span>
                        <p>AI 学习资源、AI 产业中最值得了解的洞察</p>
                    </div>
            </Col>
        </Row>
       

        
      

      
    </Container>
  );
};

export default Index;
