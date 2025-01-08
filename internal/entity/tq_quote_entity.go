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

package entity

import (
	"time"
)

const (
	QuoteStatusAvailable = 1
	QuoteStatusClosed    = 2
	QuoteStatusDeleted   = 10
	QuoteStatusPending   = 11
	QuoteUnPin           = 1
	QuotePin             = 2
	QuoteShow            = 1
	QuoteHide            = 2
)

var AdminQuoteSearchStatus = map[string]int{
	"available": QuoteStatusAvailable,
	"closed":    QuoteStatusClosed,
	"deleted":   QuoteStatusDeleted,
	"pending":   QuoteStatusPending,
}

var AdminQuoteSearchStatusIntToString = map[int]string{
	QuoteStatusAvailable: "available",
	QuoteStatusClosed:    "closed",
	QuoteStatusDeleted:   "deleted",
	QuoteStatusPending:   "pending",
}

// Quote Quote

type Quote struct {
	ID            string `json:"ID" xorm:"ID pk"`
	UserID        string `json:"user_id" xorm:"user_id"`                 // 发布者ID
	QuoteAuthorId string `json:"quote_author_id" xorm:"quote_author_id"` // 发布者ID
	QuotePieceId  string `json:"quote_piece_id" xorm:"quote_piece_id"`   // 作品（来源出处)

	Title string `json:"title" xorm:"title"`

	OriginalText    string    `xorm:"not null MEDIUMTEXT original_text"`
	ParsedText      string    `xorm:"not null MEDIUMTEXT parsed_text"`
	Status          int       `json:"status" xorm:"status"`
	CommentCount    int       `json:"comment_count" xorm:"comment_count"` // 评论总数
	RevisionID      string    `json:"revision_id" xorm:"revision_id"`     // 修订号
	CreatedAt       time.Time `json:"created_at" xorm:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" xorm:"updated_at"`
	Pin             int       `json:"pin" xorm:"pin"`
	Show            int       `json:"show" xorm:"show"`
	CollectionCount int       `json:"collection_count" xorm:"collection_count"`
	FollowCount     int       `json:"follow_count" xorm:"follow_count"`
	ViewCount       int       `json:"view_count" xorm:"view_count"`
	HotScore        int       `json:"hot_score" xorm:"hot_score"`
	UniqueViewCount int       `json:"unique_view_count" xorm:"unique_view_count"`
	VoteCount       int       `json:"vote_count" xorm:"vote_count"`

	PostUpdateTime time.Time `xorm:"post_update_time TIMESTAMP"`
}

// TableName Quote table name
func (Quote) TableName() string {
	return "tq_quote"
}
func QuoteGetAlias() string {
	return "quote"
}

// // QuoteWithTagsRevision Quote
type QuoteWithTagsRevision struct {
	Quote
	Tags []*TagSimpleInfoForRevision `json:"tags"`
}

//// TagSimpleInfoForRevision tag simple info for revision
//type TagSimpleInfoForRevision struct {
//	ID              string `xorm:"not null pk comment('tag_id') BIGINT(20) id"`
//	MainTagID       int64  `xorm:"not null default 0 BIGINT(20) main_tag_id"`
//	MainTagSlugName string `xorm:"not null default '' VARCHAR(35) main_tag_slug_name"`
//	SlugName        string `xorm:"not null default '' unique VARCHAR(35) slug_name"`
//	DisplayName     string `xorm:"not null default '' VARCHAR(35) display_name"`
//	Recommend       bool   `xorm:"not null default false BOOL recommend"`
//	Reserved        bool   `xorm:"not null default false BOOL reserved"`
//	RevisionID      string `xorm:"not null default 0 BIGINT(20) revision_id"`
//}
