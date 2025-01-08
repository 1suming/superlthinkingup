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

import { useState, useEffect, memo } from 'react';
import { Button, Form } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

import classNames from 'classnames';

import { TextArea, Mentions,ArticleUserCard ,Avatar, Icon } from '@/components';
import { usePageUsers, usePromptWithUnload } from '@/hooks';
import { parseEditMentionUser } from '@/utils';

import {
    loggedUserInfoStore,
    
  } from '@/stores';

import './index.scss';


const Index = ({
  className = '',
  value: initialValue = '',
  onSendReply,
  type = '',
  onCancel,
  mode,
}) => {
  const [value, setValue] = useState('');
  const [immData, setImmData] = useState('');
  const pageUsers = usePageUsers();
  const { t } = useTranslation('translation', { keyPrefix: 'comment' });
  const [validationErrorMsg, setValidationErrorMsg] = useState('');

  const { user, clear: clearUserStore } = loggedUserInfoStore();
console.log("user:",user)

  usePromptWithUnload({
    when: type === 'edit' ? immData !== value : Boolean(value),
  });
  useEffect(() => {
    if (!initialValue) {
      return;
    }
    setImmData(initialValue);
    setValue(initialValue);
  }, [initialValue]);

  const handleChange = (e) => {
    setValue(e.target.value);
  };
  const handleSelected = (val) => {
    setValue(val);
  };
  const handleSendReply = () => {
    onSendReply(value).catch((ex) => {
      if (ex.isError) {
        setValidationErrorMsg(ex.msg);
      }
    }).then(() => {
       //清空@cws
       setValue("");
    //    console.log("成功情况value")
    });
  };

  return (
    <>
    <div
      className={classNames(
        'd-flex align-items-start flex-column flex-md-row mb-2',
        className,
      )}>
      <div className="w-100">

      <div className="d-block d-flex flex-wrap ">
        <div className="comment-user-head">
                <Avatar
                    size="50px"
                    avatar={user?.avatar}
                    alt={user?.display_name}
                    searchStr="s=96"
                />
        </div>
        
        <div
            className={classNames('custom-form-control comment-input-textarea-wrapper', {
                'is-invalid': validationErrorMsg,
            })}>
            <Mentions
                pageUsers={pageUsers.getUsers()}
                onSelected={handleSelected}>
                <TextArea
                size="sm"
                value={type === 'edit' ? parseEditMentionUser(value) : value}
                onChange={handleChange}
                isInvalid={validationErrorMsg !== '' }
                rows={3}
                autoFocus={false}
                placeholder="相信你的评论一定很精彩，可以使用@通知某人哦"

                
                />
            </Mentions>
          
         </div>

      </div>

        
        <Form.Control.Feedback type="invalid">
          {validationErrorMsg}
        </Form.Control.Feedback>
      </div>

      
      
    </div>
    <div className="comment-input-btn d-flex justify-content-end " >
         <Button
          size="lg"
          className="text-nowrap ms-0 ms-md-2 mt-2 mt-md-0 comment-input-btn-send"
          onClick={() => handleSendReply()}>
            发布 
        </Button>
    </div>
    </>

  );
};

export default memo(Index);
