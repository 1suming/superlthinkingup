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
     key: "site_disclaim_info",
  };
  const { data,error  } =  useSiteInfoValByKey(reqParams);

 
//   useEffect(() => {
   
//   }, []);

 

  return (
    <Container style={{ paddingTop: '4rem', paddingBottom: '5rem' }}>
        
      
        <Col className="mx-auto" md={12} lg={12} xl={12}>
          <h3>免责声明</h3>
          <article  className="fmt text-break text-wrap mt-4 " dangerouslySetInnerHTML={{ __html: (data?.content)?(data?.content):"" }}   />
 
        </Col>
      

      
    </Container>
  );
};

export default React.memo(Index);
