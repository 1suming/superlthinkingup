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

package schema

type SiteMapList struct {
	QuestionIDs []*SiteMapQuestionInfo `json:"question_ids"`
	MaxPageNum  []int                  `json:"max_page_num"`
}

type SiteMapPageList struct {
	PageData []*SiteMapQuestionInfo `json:"page_data"`
}

type SiteMapQuestionInfo struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	UpdateTime string `json:"time"`
}

// @ms:
type SiteMapArticleInfo struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	UpdateTime string `json:"time"`
}

type SiteMapQuoteInfo struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	UpdateTime string `json:"time"`
}
type SiteMapQuoteAuthorInfo struct {
	ID         string `json:"id"`
	AuthorName string `json:"author_name"`
	UpdateTime string `json:"time"`
}
type SiteMapQuotePieceInfo struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	UpdateTime string `json:"time"`
}
