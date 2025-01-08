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

import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Row, Col, Form, Button, Card } from 'react-bootstrap';
import { useParams, useNavigate, useSearchParams,NavLink } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import dayjs from 'dayjs';
import classNames from 'classnames';
import isEqual from 'lodash/isEqual';
import debounce from 'lodash/debounce';
import fm from 'front-matter';
import { AutoComplete } from "antd";
import type { AutoCompleteProps } from 'antd';

import { usePageTags, usePromptWithUnload } from '@/hooks';
import { Editor, EditorRef, TagSelector ,EditorTinyMCE} from '@/components';
import type * as Type from '@/common/interface';
import { DRAFT_QUESTION_STORAGE_KEY,EnumTinyMceToolbarType } from '@/common/constants';
import {
//   saveQuestion,
  saveQuote,//
  quoteDetail,
  modifyQuote,
  useQueryRevisions,
  queryQuoteByTitle,
  queryQuoteAuthorByTitle,
  queryQuotePieceByTitle,
  getTagsBySlugName,
} from '@/services';
import {
  handleFormError,
  SaveDraft,
  storageExpires,
  scrollToElementTop,
} from '@/utils';
import { pathFactory } from '@/router/pathFactory';
import { useCaptchaPlugin } from '@/utils/pluginKit';

import SearchQuestion from './components/SearchQuestion';


interface FormDataItem {
  title: Type.FormValue<string>;
  tags: Type.FormValue<Type.Tag[]>;
  content: Type.FormValue<string>;
  answer_content: Type.FormValue<string>;
  edit_summary: Type.FormValue<string>;

  content_format: Type.FormValue<number>;
  content_plain_text: Type.FormValue<string>;
  author: Type.FormValue<string>;//作者
  author_id: Type.FormValue<string>;//作者id
  piece: Type.FormValue<string>;//出处
  piece_id: Type.FormValue<string>;//出处id
}

const saveDraft = new SaveDraft({ type: 'question' });

const Ask = () => {
  const initFormData = {
    title: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    tags: {
      value: [],
      isInvalid: false,
      errorMsg: '',
    },
    content: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    answer_content: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    edit_summary: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    //markdown或html
    content_format: {
        value: 0,
        isInvalid: false,
        errorMsg: '',
      },
    content_plain_text: {
        value: '',
        isInvalid: false,
        errorMsg: '',
      },
    author: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    author_id: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    piece: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    piece_id: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
  };
  const { t } = useTranslation('translation', { keyPrefix: 'write_quote' });
  const [formData, setFormData] = useState<FormDataItem>(initFormData);
  const [immData, setImmData] = useState<FormDataItem>(initFormData);
  const [checked, setCheckState] = useState(false);
  const [blockState, setBlockState] = useState(false);
  const [focusType, setForceType] = useState('');
  const [hasDraft, setHasDraft] = useState(false);

  const [authorOptions, setAuthorOptions] = React.useState<AutoCompleteProps['options']>([]);
  const [pieceOptions, setPieceOptions] = React.useState<AutoCompleteProps['options']>([]);

  const [authorValue, setAuthorValue] = useState('');
  const [pieceValue, setPieceValue] = useState('');
  const [pieceId, setPieceId] = useState('');
  const [authorId, setAuthorId] = useState('');
  
  enum EnumEditorType{
    Default=1,
    TinyMCE=2,
    
}
enum ArticleContentFormat {
	MARKDOWN  = 0,
	HTML      = 1,
};


const [editorType,setEditorType] = useState(EnumEditorType.TinyMCE);//@cws编辑器类型


  const resetForm = () => {
    setFormData(initFormData);
    setCheckState(false);
    setForceType('');
  };
  const [similarQuestions, setSimilarQuestions] = useState([]);
  const [similarQuestionsAuthor, setSimilarQuestionsAuthor] = useState([]);

  const editorRef = useRef<EditorRef>({
    getHtml: () => '',
  });
  const editorRef2 = useRef<EditorRef>({
    getHtml: () => '',
  });

  const { qid } = useParams();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const updateTags = (tags: string) => {
    getTagsBySlugName(tags).then((resp) => {
      // eslint-disable-next-line
      handleTagsChange(resp);
    });
  };

  const isEdit = qid !== undefined; //@ms:是否编辑

  const saveCaptcha = useCaptchaPlugin('question');
  const editCaptcha = useCaptchaPlugin('edit');

  const removeDraft = () => {
    saveDraft.save.cancel();
    saveDraft.remove();
    setHasDraft(false);
  };

  useEffect(() => {
    if (!qid) {
      // order: 1. tags query. 2. prefill query. 3. draft
      const queryTags = searchParams.get('tags');
      if (queryTags) {
        updateTags(queryTags);
      }
      const draft = storageExpires.get(DRAFT_QUESTION_STORAGE_KEY);

      const prefill = searchParams.get('prefill');
      if (prefill || draft) {
        if (prefill) {
          const file = fm<any>(decodeURIComponent(prefill));
          formData.title.value = file.attributes?.title;
          formData.content.value = file.body;
          if (!queryTags && file.attributes?.tags) {
            // Remove spaces in file.attributes.tags
            const filterTags = file.attributes.tags
              .split(',')
              .map((tag) => tag.trim())
              .join(',');
            updateTags(filterTags);
          }
        } else if (draft) {
          formData.title.value = draft.title;
          formData.content.value = draft.content;
          formData.tags.value = draft.tags;
          formData.answer_content.value = draft.answer_content;
          setCheckState(Boolean(draft.answer_content));
          setHasDraft(true);
        }
        setFormData({ ...formData });
      } else {
        resetForm();
      }
    }

    return () => {
      resetForm();
    };
  }, [qid]);

  useEffect(() => {
    const { title, tags, content, answer_content } = formData;
    const { title: editTitle, tags: editTags, content: editContent } = immData;

    // edited
    if (qid) {
      if (
        editTitle.value !== title.value ||
        editContent.value !== content.value ||
        !isEqual(
          editTags.value.map((v) => v.slug_name),
          tags.value.map((v) => v.slug_name),
        )
      ) {
        setBlockState(true);
      } else {
        setBlockState(false);
      }
      return;
    }
    // write
    if (
      title.value ||
      tags.value.length > 0 ||
      content.value ||
      answer_content.value
    ) {
      // save draft
      saveDraft.save({
        params: {
          title: title.value,
          tags: tags.value,
          content: content.value,
          answer_content: answer_content.value,
        },
        callback: () => setHasDraft(true),
      });
      setBlockState(true);
    } else {
      removeDraft();
      setBlockState(false);
    }
  }, [formData]);

  usePromptWithUnload({
    when: blockState,
  });

  const { data: revisions = [] } = useQueryRevisions(qid);
//   console.log("http revisions:",revisions)

  useEffect(() => {
    if (!isEdit) {
      return;
    }
    quoteDetail(qid).then((res) => {
        console.log("http res: quoteDetail",res)
      formData.title.value = res.title;
      formData.content_format.value =res.content_format 
      if(res.content_format == ArticleContentFormat.MARKDOWN){
        setEditorType(EnumEditorType.Default)
        formData.content.value = res.content;
      }else if(res.content_format == ArticleContentFormat.HTML){
        setEditorType(EnumEditorType.TinyMCE)
        formData.content.value=res.html; //如果是html格式的，存在html里面。
      }
     
      formData.tags.value = res.tags.map((item) => {
        return {
          ...item,
          parsed_text: '',
          original_text: '',
        };
      });

      

      setImmData({ ...formData });
      setFormData({ ...formData });
    });
  }, [qid]);

  const querySimilarQuotes = useCallback(
    debounce((title) => {
      queryQuoteByTitle(title).then((res) => {
        setSimilarQuestions(res);
      });
    }, 400),
    [],
  );
  const querySimilarQuotesAuthor = useCallback(
    debounce((title) => {
      queryQuoteAuthorByTitle(title).then((res) => {
        setSimilarQuestionsAuthor(res);
      });
    }, 400),
    [],
  );


  const handleTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      title: { value: e.currentTarget.value, errorMsg: '', isInvalid: false },
    });
    if (e.currentTarget.value.length >= 3) {
      querySimilarQuotes(e.currentTarget.value);
    }
    if (e.currentTarget.value.length === 0) {
      setSimilarQuestions([]);
    }
  };
  const handleContentChange = (value: string,plainValue:string) => {
     console.log("handleContentChange value:",value)
     console.log("handleContentChange plainValue:",plainValue)
     value=plainValue //@只使用纯文本
    setFormData({
      ...formData,
      content: { value, errorMsg: '', isInvalid: false },
      // content_plain_text: { value, errorMsg: '', isInvalid: false },
    });
  };
  const handleTagsChange = (value) =>
    setFormData({
      ...formData,
      tags: { value, errorMsg: '', isInvalid: false },
    });

  const handleAnswerChange = (value: string) =>
    setFormData({
      ...formData,
      answer_content: { value, errorMsg: '', isInvalid: false },
    });

  const handleSummaryChange = (evt: React.ChangeEvent<HTMLInputElement>) =>
    setFormData({
      ...formData,
      edit_summary: {
        ...formData.edit_summary,
        value: evt.currentTarget.value,
      },
    });

  const deleteDraft = () => {
    const res = window.confirm(t('discard_confirm', { keyPrefix: 'draft' }));
    if (res) {
      removeDraft();
      resetForm();
    }
  };

  const submitModifyQuote = (params) => {
    setBlockState(false);
    const ep = {
      ...params,
      id: qid,
      edit_summary: formData.edit_summary.value,
    };
    console.log("submitModifyQuote final ep:",ep)
    const imgCode = editCaptcha?.getCaptcha();
    if (imgCode?.verify) {
      ep.captcha_code = imgCode.captcha_code;
      ep.captcha_id = imgCode.captcha_id;
    }
    modifyQuote(ep)
      .then(async (res) => {
        await editCaptcha?.close();
        navigate(pathFactory.quoteLanding(qid, res?.url_title), {
          state: { isReview: res?.wait_for_review },
        });
      })
      .catch((err) => {
        if (err.isError) {
          editCaptcha?.handleCaptchaError(err.list);
          const data = handleFormError(err, formData);
          setFormData({ ...data });
          const ele = document.getElementById(err.list[0].error_field);
          scrollToElementTop(ele);
        }
      });
  };

  //@form 提交
// const [contentError,setContentError] = useState(false);
  const submitQuote = async (params) => {
    setBlockState(false);
    const imgCode = saveCaptcha?.getCaptcha();
    if (imgCode?.verify) {
      params.captcha_code = imgCode.captcha_code;
      params.captcha_id = imgCode.captcha_id;
    }
    let res;
    // if (checked) {
    //     console.log("saveQuestionWithAnswer true");
    //   res = await saveQuestionWithAnswer({
    //     ...params,
    //     answer_content: formData.answer_content.value,
    //   }).catch((err) => {
    //     if (err.isError) {
    //       const captchaErr = saveCaptcha?.handleCaptchaError(err.list);
    //       if (!(captchaErr && err.list.length === 1)) {
    //         const data = handleFormError(err, formData);
    //         setFormData({ ...data });
    //         const ele = document.getElementById(err.list[0].error_field);
    //         scrollToElementTop(ele);
    //       }
    //     }
    //   });
    // } else {
        //保存 @cws
      res = await saveQuote(params).catch((err) => {
        console.log("saveQuote err:",err)
        if (err.isError) {
          const captchaErr = saveCaptcha?.handleCaptchaError(err.list);
          console.log("captchaErr:",captchaErr)
          if (!(captchaErr && err.list.length === 1)) {
            console.log("handleFormError err:",err)
            const data = handleFormError(err, formData);
            console.log("handleFormError data:",data)
            setFormData({ ...data });
            // if(err.list[0].error_field == "content"){//@cws:content 错误，设置contentError为true
            //   setContentError(true);
            // }
            const ele = document.getElementById(err.list[0].error_field);
            console.log("ele:",ele)
            scrollToElementTop(ele);
          }
        }
      });
    // }
    console.log("saveQuote res:",res)

    const id = res?.id || res?.question?.id;
    if (id) {
      await saveCaptcha?.close();
      if (checked) {
        navigate(pathFactory.quoteLanding(id, res?.question?.url_title));
      } else {
        navigate(pathFactory.quoteLanding(id, res?.url_title));
      }
    }
    removeDraft();
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    event.stopPropagation();

    let  content_format:ArticleContentFormat = ArticleContentFormat.MARKDOWN
    if(editorType == EnumEditorType.TinyMCE ){
        content_format =ArticleContentFormat.HTML
    } else if(editorType == EnumEditorType.Default ){
        content_format=ArticleContentFormat.MARKDOWN
    }
    console.log("content_format:",content_format)
   
    const params: Type.QuoteParams = {
      title: formData.title.value,
      content: formData.content.value,
      tags: formData.tags.value,

      content_format: content_format,
      author: authorValue,
      author_id: authorId,
      piece_id: pieceId,
      piece_name: pieceValue,

    };
    
    console.log("submitQuote params:",params)

    if (isEdit) {
      if (!editCaptcha) {
        submitModifyQuote(params);
        return;
      }
      editCaptcha.check(() => submitModifyQuote(params));
    } else {
      if (!saveCaptcha) {
        submitQuote(params);
        return;
      }
      saveCaptcha?.check(async () => {
        submitQuote(params);
      });
    }
  };
  const backPage = () => {
    navigate(-1);
  };

  const handleSelectedRevision = (e) => {
//    console.log("handleSelectedRevision选择指定版本 revisions:",revisions)
    const index = e.target.value;
    const revision = revisions[index];
    console.log("handleSelectedRevision选择指定版本 revision",revision)
    formData.content.value = revision.content?.content || '';
    // console.log("revision.content:", formData.content.value )

    setImmData({ ...formData });
    setFormData({ ...formData });
  };
  const bool = similarQuestions.length > 0 && !isEdit;
  let pageTitle = t('ask_a_question', { keyPrefix: 'page_title' });
  if (isEdit) {
    pageTitle = t('edit_question', { keyPrefix: 'page_title' });
  }
  usePageTags({
    title: pageTitle,
  });
  
  const handleChangeEditorType = (evt)=>{
     evt.preventDefault();
     setEditorType(editorType==EnumEditorType.Default? EnumEditorType.TinyMCE:EnumEditorType.Default);
  }

  const handleAuthorSearch =async (value: string) => {
    console.log("handleAuthorSearch value:",value)
    
      // if (!value || value.includes('@')) {
      //   return [];
      // }
      if (value.length === 0) {
        // setSimilarQuestions([]);
        setAuthorOptions([]);
        return ;
      }
      //https://github.com/ant-design/ant-design/issues/22004 为什么autocomplete 
      if (value.length >= 1) {
        let res =await queryQuoteAuthorByTitle(value).catch((err)=>{
           console.log("queryQuoteAuthorByTitle err:",err)
        });
        console.log("queryQuoteAuthorByTitle res:",res)
        let options = res.map((item)=>{
          return {
            label: item.author_name,
            value: item.author_name,
            my_id: item.id,
          }
        });
        setAuthorOptions(options);

      }
      
      
      // return ['gmail.com', '163.com', 'qq.com'].map((domain) => ({
      //   label: `${value}@${domain}`,
      //   value: `${value}@${domain}`,
      // }));
     
  };
  const handleAuthorSearchDebounce = useCallback(debounce(handleAuthorSearch, 400), []);

  const handleAuthorSelect = (value: string, option: any) => {
    console.log("handleAuthorSelect value:",value)
    setAuthorValue(value);
    setAuthorId(option.my_id); //option.my_id 是作者id
    console.log("handleAuthorSelect option:",option)
    setFormData({
      ...formData,
      author: { value, errorMsg: '', isInvalid: false },
    });
  };

  const handlePieceSearchDebounce = useCallback(
    debounce(async (value) => {
      if (value.length === 0) {
        setPieceOptions([]);
        return;
      }
      if (value.length >= 1) {
        const res = await queryQuotePieceByTitle(value);
        const options = res.map((item) => ({
          label: item.title,
          value: item.title,
          my_id: item.id,
        }));
        setPieceOptions(options);
      }
    }, 400),
    [],
  );

  const handlePieceSelect = (value: string, option: any) => {
    console.log("handlePieceSelect value:",value)
    setPieceValue(value);
    setPieceId(option.my_id); //option.my_id 是出处id
    console.log("handlePieceSelect option:",option)
    setFormData({
      ...formData,
      piece: { value, errorMsg: '', isInvalid: false },
    });
  };
  return (
    <div className="pt-4 mb-5">
      <h3 className="mb-4">{isEdit ? t('edit_title') : t('title')}</h3>
      <Row>
        <Col className="page-main-not flex-auto">
          <Form noValidate onSubmit={handleSubmit}>
            {isEdit && (
              <Form.Group controlId="revision" className="mb-3">
                <Form.Label>{t('form.fields.revision.label')}</Form.Label>
                <Form.Select onChange={handleSelectedRevision}>
                  {revisions.map(({ reason, create_at, user_info }, index) => {
                    const date = dayjs(create_at * 1000)
                      .tz()
                      .format(t('long_date_with_time', { keyPrefix: 'dates' }));
                    return (
                      <option key={`${create_at}`} value={index}>
                        {`${date} - ${user_info.display_name} - ${
                          reason ||
                          (index === revisions.length - 1
                            ? t('default_first_reason')
                            : t('default_reason'))
                        }`}
                      </option>
                    );
                  })}
                </Form.Select>
              </Form.Group>
            )}

        

         

            <Form.Group controlId="content">
                <div className="d-flex">
                     <Form.Label>{t('form.fields.body.label')}</Form.Label>
                     {/* <NavLink
                        className=" ms-3 "
                            to="#"
                            onClick={handleChangeEditorType}
                            >
                            <span id="switchEditorTypeBtn">切换为 { editorType==EnumEditorType.Default? "可视化编辑器 ":"MarkDown编辑器" } </span>
                        </NavLink> */}
                 </div>
              {editorType==EnumEditorType.TinyMCE ? (
                    <EditorTinyMCE 
                        id="content"
                        onChange={handleContentChange}
                        editorPlaceholder=""
                        value={formData.content.value}
                        menubar={false}
                        toolbarType={EnumTinyMceToolbarType.Simple}
                        min_height={200}

                        className={classNames(
                          'form-control ',
                          focusType === 'content' && 'focus',
                          formData.content.isInvalid && 'is-invalid',
                          
                        )}

                        onFocus={() => {
                          setForceType('content');
                        }}
                        onBlur={() => {
                          setForceType('');
                        }}

                    />
              ):(
              <Editor
                value={formData.content.value}
                onChange={handleContentChange}
                className={classNames(
                  'form-control ',
                  focusType === 'content' && 'focus',
                  formData.content.isInvalid && 'is-invalid',
                  "article-editor"
                )}
                onFocus={() => {
                  setForceType('content');
                }}
                onBlur={() => {
                  setForceType('');
                }}
                ref={editorRef}
              />
             )}
              <Form.Control.Feedback type="invalid">
                {formData.content.errorMsg}
              </Form.Control.Feedback>
            </Form.Group>

            <Form.Group controlId="title" className="mb-3">
              <Form.Label>总结标题</Form.Label>
              <Form.Control
                type="text"
                value={formData.title.value}
                isInvalid={formData.title.isInvalid}
                onChange={handleTitleChange}
                placeholder="金句总结，可为空。太长建议填写"
                autoFocus
                contentEditable
              />
              <Form.Control.Feedback type="invalid">
                {formData.title.errorMsg}
              </Form.Control.Feedback>
              {bool && <SearchQuestion similarQuestions={similarQuestions} />}
            </Form.Group>

            <Form.Group controlId="author" className="mb-3">
              <Form.Label>作者</Form.Label>
              <div>
                <AutoComplete
                  style={{ width: 200 }}
                  value={authorValue}
                  onSelect={handleAuthorSelect}
                  onChange={(value) => setAuthorValue(value)}
                  onSearch={handleAuthorSearchDebounce}
                  placeholder="作者"
                  options={authorOptions}
                />
              </div>
             
              <Form.Control.Feedback type="invalid">
                {formData.author.errorMsg}
              </Form.Control.Feedback>
               
            </Form.Group>
            <Form.Group controlId="piece" className="mb-3">
              <Form.Label>出处</Form.Label>
              <div>
                <AutoComplete
                  style={{ width: 200 }}
                  value={pieceValue}
                  onSelect={handlePieceSelect}
                  onChange={(value) => setPieceValue(value)}
                  onSearch={handlePieceSearchDebounce}
                  placeholder="出处，如书名、电影名..."
                  options={pieceOptions}
                />
              </div>
             
              <Form.Control.Feedback type="invalid">
                {formData.piece.errorMsg}
              </Form.Control.Feedback>
               
            </Form.Group>


            <Form.Group controlId="tags" className="my-3">
              <Form.Label>{t('form.fields.tags.label')}</Form.Label>
              <TagSelector
                value={formData.tags.value}
                onChange={handleTagsChange}
                showRequiredTag
                maxTagLength={5}
                isInvalid={formData.tags.isInvalid}
                errMsg={formData.tags.errorMsg}
              />
            </Form.Group>

            {!isEdit && (
              <>
                {/* <Form.Switch
                  checked={checked}
                  type="switch"
                  label={t('answer_question')}
                  onChange={(e) => setCheckState(e.target.checked)}
                  id="radio-answer"
                /> */}
                {checked && (
                  <Form.Group controlId="answer" className="mt-3">
                    <Form.Label>{t('form.fields.answer.label')}</Form.Label>
                    <Editor
                      value={formData.answer_content.value}
                      onChange={handleAnswerChange}
                      ref={editorRef2}
                      className={classNames(
                        'form-control p-0',
                        focusType === 'answer' && 'focus',
                        formData.answer_content.isInvalid && 'is-invalid',
                      )}
                      onFocus={() => {
                        setForceType('answer');
                      }}
                      onBlur={() => {
                        setForceType('');
                      }}
                    />
                    <Form.Control
                      type="text"
                      isInvalid={formData.answer_content.isInvalid}
                      hidden
                    />
                    <Form.Control.Feedback type="invalid">
                      {formData.answer_content.errorMsg}
                    </Form.Control.Feedback>
                  </Form.Group>
                )}
              </>
            )}

            {isEdit && (
              <Form.Group controlId="edit_summary" className="my-3">
                <Form.Label>{t('form.fields.edit_summary.label')}</Form.Label>
                <Form.Control
                  type="text"
                  defaultValue={formData.edit_summary.value}
                  isInvalid={formData.edit_summary.isInvalid}
                  placeholder={t('form.fields.edit_summary.placeholder')}
                  onChange={handleSummaryChange}
                  contentEditable
                />
                <Form.Control.Feedback type="invalid">
                  {formData.edit_summary.errorMsg}
                </Form.Control.Feedback>
              </Form.Group>
            )}
            {!checked && (
              <div className="mt-3">
                <Button type="submit" className="me-2">
                  {isEdit ? t('btn_save_edits') : t('btn_post_question')}
                </Button>
                {isEdit && (
                  <Button variant="link" onClick={backPage}>
                    {t('cancel', { keyPrefix: 'btns' })}
                  </Button>
                )}

                {hasDraft && (
                  <Button variant="link" onClick={deleteDraft}>
                    {t('discard_draft', { keyPrefix: 'btns' })}
                  </Button>
                )}
              </div>
            )}
            {checked && (
              <div className="mt-3">
                <Button type="submit">{t('post_question&answer')}</Button>
                {hasDraft && (
                  <Button variant="link" className="ms-2" onClick={deleteDraft}>
                    {t('discard_draft', { keyPrefix: 'btns' })}
                  </Button>
                )}
              </div>
            )}
          </Form>
        </Col>
      
      </Row>
    </div>
  );
};

export default Ask;
