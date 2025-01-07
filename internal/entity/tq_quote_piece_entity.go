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
	QuotePieceStatusAvailable = 1
	QuotePieceStatusClosed    = 2
	QuotePieceStatusDeleted   = 10
	QuotePieceStatusPending   = 11
	QuotePieceUnPin           = 1
	QuotePiecePin             = 2
	QuotePieceShow            = 1
	QuotePieceHide            = 2
)

var AdminQuotePieceSearchStatus = map[string]int{
	"available": QuotePieceStatusAvailable,
	"closed":    QuotePieceStatusClosed,
	"deleted":   QuotePieceStatusDeleted,
	"pending":   QuotePieceStatusPending,
}

var AdminQuotePieceSearchStatusIntToString = map[int]string{
	QuotePieceStatusAvailable: "available",
	QuotePieceStatusClosed:    "closed",
	QuotePieceStatusDeleted:   "deleted",
	QuotePieceStatusPending:   "pending",
}

type QuotePiece struct {
	ID     string `json:"id" xorm:"id pk "`
	UserID string `json:"user_id" xorm:"user_id"` // 发布者ID

	Title string `json:"title" xorm:"title"`

	PublishDate time.Time `json:"publish_date" xorm:"publish_date"` // 出版日期

	OriginalText string `xorm:"not null MEDIUMTEXT original_text"`
	ParsedText   string `xorm:"not null MEDIUMTEXT parsed_text"`

	PieceType int `json:"piece_type" xorm:"piece_type"` // book'', ''article'', ''speech'', ''website'', ''other'' 来源类型，使用枚举约束来源类别

	Status int    `json:"status" xorm:"status"`
	Avatar string `json:"avatar" xorm:"avatar"`
	//Bio    string `json:"bio" xorm:"bio"`
	//BioHtml    string    `json:"bio_html" xorm:"bio_html"`
	CreatedAt time.Time `json:"created_at" xorm:"created_at"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated_at"`

	Pin  int `json:"pin" xorm:"pin"`
	Show int `json:"show" xorm:"show"`

	CollectionCount int `json:"collection_count" xorm:"collection_count"`
	FollowCount     int `json:"follow_count" xorm:"follow_count"`
	ViewCount       int `json:"view_count" xorm:"view_count"`
	HotScore        int `json:"hot_score" xorm:"hot_score"`
	UniqueViewCount int `json:"unique_view_count" xorm:"unique_view_count"`
	VoteCount       int `json:"vote_count" xorm:"vote_count"`
}

// TableName QuotePiece table name
func (QuotePiece) TableName() string {
	return "tq_quote_piece"
}

// // QuoteWithTagsRevision Quote
type QuotePieceWithTagsRevision struct {
	QuotePiece
	Tags []*TagSimpleInfoForRevision `json:"tags"`
}
