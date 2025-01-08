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

package quote_piece

import (
	"context"
	"encoding/json"
	"fmt"
	quote_common "github.com/apache/incubator-answer/internal/service_quote/quote_common"
	"strings"
	"time"
	"unicode"

	"github.com/apache/incubator-answer/internal/base/constant"
	"github.com/apache/incubator-answer/internal/base/data"
	"github.com/apache/incubator-answer/internal/base/handler"
	"github.com/apache/incubator-answer/internal/base/pager"
	"github.com/apache/incubator-answer/internal/base/reason"
	"github.com/apache/incubator-answer/internal/entity"
	"github.com/apache/incubator-answer/internal/schema"
	"github.com/apache/incubator-answer/internal/service/unique"
	"github.com/apache/incubator-answer/pkg/htmltext"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/apache/incubator-answer/plugin"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
	"xorm.io/builder"
	"xorm.io/xorm"
)

// quotePieceRepo quotePiece repository
type quotePieceRepo struct {
	data         *data.Data
	uniqueIDRepo unique.UniqueIDRepo
}

// NewQuotePieceRepo new repository
func NewQuotePieceRepo(
	data *data.Data,
	uniqueIDRepo unique.UniqueIDRepo,
) quote_common.QuotePieceRepo {
	return &quotePieceRepo{
		data:         data,
		uniqueIDRepo: uniqueIDRepo,
	}
}

// AddQuotePiece add quotePiece
func (qr *quotePieceRepo) AddQuotePiece(ctx context.Context, quotePiece *entity.QuotePiece) (err error) {
	quotePiece.ID, err = qr.uniqueIDRepo.GenUniqueIDStr(ctx, quotePiece.TableName())
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_, err = qr.data.DB.Context(ctx).Insert(quotePiece)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		quotePiece.ID = uid.EnShortID(quotePiece.ID)
	}
	return
}

// RemoveQuotePiece delete quotePiece
func (qr *quotePieceRepo) RemoveQuotePiece(ctx context.Context, id string) (err error) {
	id = uid.DeShortID(id)
	_, err = qr.data.DB.Context(ctx).Where("id =?", id).Delete(&entity.QuotePiece{})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// UpdateQuotePiece update quotePiece
func (qr *quotePieceRepo) UpdateQuotePiece(ctx context.Context, quotePiece *entity.QuotePiece, Cols []string) (err error) {
	quotePiece.ID = uid.DeShortID(quotePiece.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quotePiece.ID).Cols(Cols...).Update(quotePiece)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		quotePiece.ID = uid.EnShortID(quotePiece.ID)
	}
	_ = qr.UpdateSearch(ctx, quotePiece.ID)
	return
}

func (qr *quotePieceRepo) UpdatePvCount(ctx context.Context, quotePieceID string) (err error) {
	quotePieceID = uid.DeShortID(quotePieceID)
	quotePiece := &entity.QuotePiece{}
	_, err = qr.data.DB.Context(ctx).Where("id =?", quotePieceID).Incr("view_count", 1).Update(quotePiece)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quotePiece.ID)
	return nil
}

//func (qr *quotePieceRepo) UpdateAnswerCount(ctx context.Context, quotePieceID string, num int) (err error) {
//	quotePieceID = uid.DeShortID(quotePieceID)
//	quotePiece := &entity.QuotePiece{}
//	quotePiece.AnswerCount = num
//	_, err = qr.data.DB.Context(ctx).Where("id =?", quotePieceID).Cols("answer_count").Update(quotePiece)
//	if err != nil {
//		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
//	}
//	_ = qr.UpdateSearch(ctx, quotePiece.ID)
//	return nil
//}

func (qr *quotePieceRepo) UpdateCollectionCount(ctx context.Context, quotePieceID string) (count int64, err error) {
	quotePieceID = uid.DeShortID(quotePieceID)
	_, err = qr.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)
		count, err = session.Count(&entity.Collection{ObjectID: quotePieceID})
		if err != nil {
			return nil, err
		}

		quotePiece := &entity.QuotePiece{CollectionCount: int(count)}
		_, err = session.ID(quotePieceID).MustCols("collection_count").Update(quotePiece)
		if err != nil {
			return nil, err
		}
		return
	})
	if err != nil {
		return 0, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return count, nil
}

func (qr *quotePieceRepo) UpdateQuotePieceStatus(ctx context.Context, quotePieceID string, status int) (err error) {
	quotePieceID = uid.DeShortID(quotePieceID)
	_, err = qr.data.DB.Context(ctx).ID(quotePieceID).Cols("status").Update(&entity.QuotePiece{Status: status})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quotePieceID)
	return nil
}

func (qr *quotePieceRepo) UpdateQuotePieceStatusWithOutUpdateTime(ctx context.Context, quotePiece *entity.QuotePiece) (err error) {
	quotePiece.ID = uid.DeShortID(quotePiece.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quotePiece.ID).Cols("status").Update(quotePiece)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quotePiece.ID)
	return nil
}

func (qr *quotePieceRepo) RecoverQuotePiece(ctx context.Context, quotePieceID string) (err error) {
	quotePieceID = uid.DeShortID(quotePieceID)
	_, err = qr.data.DB.Context(ctx).ID(quotePieceID).Cols("status").Update(&entity.QuotePiece{Status: entity.QuotePieceStatusAvailable})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quotePieceID)
	return nil
}

func (qr *quotePieceRepo) UpdateQuotePieceOperation(ctx context.Context, quotePiece *entity.QuotePiece) (err error) {
	quotePiece.ID = uid.DeShortID(quotePiece.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quotePiece.ID).Cols("pin", "show").Update(quotePiece)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (qr *quotePieceRepo) UpdateAccepted(ctx context.Context, quotePiece *entity.QuotePiece) (err error) {
	quotePiece.ID = uid.DeShortID(quotePiece.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quotePiece.ID).Cols("accepted_answer_id").Update(quotePiece)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quotePiece.ID)
	return nil
}

func (qr *quotePieceRepo) UpdateLastAnswer(ctx context.Context, quotePiece *entity.QuotePiece) (err error) {
	quotePiece.ID = uid.DeShortID(quotePiece.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quotePiece.ID).Cols("last_answer_id").Update(quotePiece)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quotePiece.ID)
	return nil
}

// GetQuotePiece get quotePiece one
func (qr *quotePieceRepo) GetQuotePiece(ctx context.Context, id string) (
	quotePiece *entity.QuotePiece, exist bool, err error,
) {
	id = uid.DeShortID(id)
	quotePiece = &entity.QuotePiece{}
	quotePiece.ID = id
	exist, err = qr.data.DB.Context(ctx).Where("id = ?", id).Get(quotePiece)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		quotePiece.ID = uid.EnShortID(quotePiece.ID)
	}
	return
}

func (qr *quotePieceRepo) GetQuotePieceSimple(ctx context.Context, id string) (
	quotePieceBasicInfo *schema.QuotePieceBasicInfo, exist bool, err error,
) {
	id = uid.DeShortID(id)
	quotePiece := &entity.QuotePiece{}
	quotePieceBasicInfo = &schema.QuotePieceBasicInfo{}
	//quotePiece.ID = id
	exist, err = qr.data.DB.Context(ctx).Table(quotePiece.TableName()).Where("id = ?", id).Get(quotePieceBasicInfo)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		quotePieceBasicInfo.ID = uid.EnShortID(quotePieceBasicInfo.ID)
	}
	return
}

// GetQuotePiecesByTitle get quotePiece list by title
func (qr *quotePieceRepo) GetQuotePiecesByTitle(ctx context.Context, title string, pageSize int) (
	quotePieceList []*entity.QuotePiece, err error) {
	quotePieceList = make([]*entity.QuotePiece, 0)
	session := qr.data.DB.Context(ctx)
	session.Where("status != ?", entity.QuotePieceStatusDeleted)
	session.Where("title like ?", "%"+title+"%")
	session.Limit(pageSize)
	err = session.Find(&quotePieceList)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range quotePieceList {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return
}
func (qr *quotePieceRepo) GetQuotePieceByTitle(ctx context.Context, title string) (
	quotePiece *entity.QuotePiece, err error) {
	//quotePieceList = make([]*entity.QuotePiece, 0)
	quotePieceList := make([]*entity.QuotePiece, 0)
	session := qr.data.DB.Context(ctx)
	session.Where("status != ?", entity.QuotePieceStatusDeleted)
	session.Where("title = ?", title) //不要用like，用等于
	//session.Limit(pageSize)
	err = session.Find(&quotePieceList)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range quotePieceList {
			item.ID = uid.EnShortID(item.ID)
			quotePiece = item
			return //@只要一条
		}
	}
	return
}
func (qr *quotePieceRepo) FindByID(ctx context.Context, id []string) (quotePieceList []*entity.QuotePiece, err error) {
	for key, itemID := range id {
		id[key] = uid.DeShortID(itemID)
	}
	quotePieceList = make([]*entity.QuotePiece, 0)
	err = qr.data.DB.Context(ctx).Table("quotePiece").In("id", id).Find(&quotePieceList)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range quotePieceList {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return
}

// GetQuotePieceList get quotePiece list all
func (qr *quotePieceRepo) GetQuotePieceList(ctx context.Context, quotePiece *entity.QuotePiece) (quotePieceList []*entity.QuotePiece, err error) {
	quotePiece.ID = uid.DeShortID(quotePiece.ID)
	quotePieceList = make([]*entity.QuotePiece, 0)
	err = qr.data.DB.Context(ctx).Find(quotePieceList, quotePiece)
	if err != nil {
		return quotePieceList, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	for _, item := range quotePieceList {
		item.ID = uid.DeShortID(item.ID)
	}
	return
}

func (qr *quotePieceRepo) GetQuotePieceCount(ctx context.Context) (count int64, err error) {
	session := qr.data.DB.Context(ctx)
	session.Where(builder.Lt{"status": entity.QuotePieceStatusDeleted})
	count, err = session.Count(&entity.QuotePiece{Show: entity.QuotePieceShow})
	if err != nil {
		return 0, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return count, nil
}

func (qr *quotePieceRepo) GetUserQuotePieceCount(ctx context.Context, userID string, show int) (count int64, err error) {
	session := qr.data.DB.Context(ctx)
	session.Where(builder.Lt{"status": entity.QuotePieceStatusDeleted})
	count, err = session.Count(&entity.QuotePiece{UserID: userID, Show: show})
	if err != nil {
		return count, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (qr *quotePieceRepo) SitemapQuotePieces(ctx context.Context, page, pageSize int) (
	quotePieceIDList []*schema.SiteMapQuotePieceInfo, err error) {
	page = page - 1
	quotePieceIDList = make([]*schema.SiteMapQuotePieceInfo, 0)

	// try to get sitemap data from cache
	cacheKey := fmt.Sprintf(constant.SiteMapQuotePieceCacheKeyPrefix, page)
	cacheData, exist, err := qr.data.Cache.GetString(ctx, cacheKey)
	if err == nil && exist {
		_ = json.Unmarshal([]byte(cacheData), &quotePieceIDList)
		return quotePieceIDList, nil
	}

	// get sitemap data from db
	rows := make([]*entity.QuotePiece, 0)
	session := qr.data.DB.Context(ctx)
	//session.Select("id,title,created_at,post_update_time")
	session.Select("id,title,created_at,update_at")

	session.Where("`show` = ?", entity.QuotePieceShow)
	session.Where("status = ? OR status = ?", entity.QuotePieceStatusAvailable, entity.QuotePieceStatusClosed)
	session.Limit(pageSize, page*pageSize)
	session.Asc("created_at")
	err = session.Find(&rows)
	if err != nil {
		return quotePieceIDList, err
	}

	// warp data
	for _, quotePiece := range rows {
		item := &schema.SiteMapQuotePieceInfo{ID: quotePiece.ID}
		if handler.GetEnableShortID(ctx) {
			item.ID = uid.EnShortID(quotePiece.ID)
		}
		item.Title = htmltext.UrlTitle(quotePiece.Title)
		//if quotePiece.PostUpdateTime.IsZero() {
		//	item.UpdateTime = quotePiece.CreatedAt.Format(time.RFC3339)
		//} else {
		//	item.UpdateTime = quotePiece.PostUpdateTime.Format(time.RFC3339)
		//}

		quotePieceIDList = append(quotePieceIDList, item)
	}

	// set sitemap data to cache
	cacheDataByte, _ := json.Marshal(quotePieceIDList)
	if err := qr.data.Cache.SetString(ctx, cacheKey, string(cacheDataByte), constant.SiteMapQuotePieceCacheTime); err != nil {
		log.Error(err)
	}
	return quotePieceIDList, nil
}

// GetQuotePiecePage query quotePiece page
func (qr *quotePieceRepo) GetQuotePiecePage(ctx context.Context, page, pageSize int,
	tagIDs []string, userID, orderCond string, inDays int, showHidden, showPending bool) (
	quotePieceList []*entity.QuotePiece, total int64, err error) {
	quotePieceList = make([]*entity.QuotePiece, 0)
	session := qr.data.DB.Context(ctx)
	status := []int{entity.QuotePieceStatusAvailable, entity.QuotePieceStatusClosed}
	if showPending {
		status = append(status, entity.QuotePieceStatusPending)
	}
	session.In("quotePiece.status", status)
	if len(tagIDs) > 0 {
		session.Join("LEFT", "tag_rel", "quotePiece.id = tag_rel.object_id")
		session.In("tag_rel.tag_id", tagIDs)
		session.And("tag_rel.status = ?", entity.TagRelStatusAvailable)
	}
	if len(userID) > 0 {
		session.And("quotePiece.user_id = ?", userID)
		if !showHidden {
			session.And("quotePiece.show = ?", entity.QuotePieceShow)
		}
	} else {
		session.And("quotePiece.show = ?", entity.QuotePieceShow)
	}
	if inDays > 0 {
		session.And("quotePiece.created_at > ?", time.Now().AddDate(0, 0, -inDays))
	}

	switch orderCond {
	case "newest":
		session.OrderBy("quotePiece.pin desc,quotePiece.created_at DESC")
	case "active":
		if inDays == 0 {
			session.And("quotePiece.created_at > ?", time.Now().AddDate(0, 0, -180))
		}
		session.And("quotePiece.post_update_time > ?", time.Now().AddDate(0, 0, -90))
		session.OrderBy("quotePiece.pin desc,quotePiece.post_update_time DESC, quotePiece.updated_at DESC")
	case "hot":
		session.OrderBy("quotePiece.pin desc,quotePiece.hot_score DESC")
	case "score":
		session.OrderBy("quotePiece.pin desc,quotePiece.vote_count DESC, quotePiece.view_count DESC")
	case "unanswered":
		session.Where("quotePiece.answer_count = 0")
		session.OrderBy("quotePiece.pin desc,quotePiece.created_at DESC")
	}

	total, err = pager.Help(page, pageSize, &quotePieceList, &entity.QuotePiece{}, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range quotePieceList {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return quotePieceList, total, err
}

// GetRecommendQuotePiecePageByTags get recommend quotePiece page by tags
func (qr *quotePieceRepo) GetRecommendQuotePiecePageByTags(ctx context.Context, userID string, tagIDs, followedQuotePieceIDs []string, page, pageSize int) (
	quotePieceList []*entity.QuotePiece, total int64, err error) {
	quotePieceList = make([]*entity.QuotePiece, 0)
	orderBySQL := "quotePiece.pin DESC, quotePiece.created_at DESC"

	// Please Make sure every quotePiece has at least one tag
	selectSQL := entity.QuotePiece{}.TableName() + ".*"
	if len(followedQuotePieceIDs) > 0 {
		idStr := "'" + strings.Join(followedQuotePieceIDs, "','") + "'"
		selectSQL += fmt.Sprintf(", CASE WHEN quotePiece.id IN (%s) THEN 0 ELSE 1 END AS order_priority", idStr)
		orderBySQL = "order_priority, " + orderBySQL
	}
	session := qr.data.DB.Context(ctx).Select(selectSQL)

	if len(tagIDs) > 0 {
		session.Where("quotePiece.user_id != ?", userID).
			And("quotePiece.id NOT IN (SELECT quotePiece_id FROM answer WHERE user_id = ?)", userID).
			Join("INNER", "tag_rel", "quotePiece.id = tag_rel.object_id").
			And("tag_rel.status = ?", entity.TagRelStatusAvailable).
			Join("INNER", "tag", "tag.id = tag_rel.tag_id").
			In("tag.id", tagIDs)
	} else if len(followedQuotePieceIDs) == 0 {
		return quotePieceList, 0, nil
	}

	if len(followedQuotePieceIDs) > 0 {
		if len(tagIDs) > 0 {
			// if tags provided, show followed quotePieces and tag quotePieces
			session.Or(builder.In("quotePiece.id", followedQuotePieceIDs))
		} else {
			// if no tags, only show followed quotePieces
			session.Where(builder.In("quotePiece.id", followedQuotePieceIDs))
		}
	}

	session.
		And("quotePiece.show = ? and quotePiece.status = ?", entity.QuotePieceShow, entity.QuotePieceStatusAvailable).
		Distinct("quotePiece.id").
		OrderBy(orderBySQL)

	total, err = pager.Help(page, pageSize, &quotePieceList, &entity.QuotePiece{}, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	if handler.GetEnableShortID(ctx) {
		for _, item := range quotePieceList {
			item.ID = uid.EnShortID(item.ID)
		}
	}

	return quotePieceList, total, err
}

func (qr *quotePieceRepo) AdminQuotePiecePage(ctx context.Context, search *schema.AdminQuotePiecePageReq) ([]*entity.QuotePiece, int64, error) {
	var (
		count   int64
		err     error
		session = qr.data.DB.Context(ctx).Table("quotePiece")
	)

	session.Where(builder.Eq{
		"status": search.Status,
	})

	rows := make([]*entity.QuotePiece, 0)
	if search.Page > 0 {
		search.Page = search.Page - 1
	} else {
		search.Page = 0
	}
	if search.PageSize == 0 {
		search.PageSize = constant.DefaultPageSize
	}

	// search by quotePiece title like or quotePiece id
	if len(search.Query) > 0 {
		// check id search
		var (
			idSearch = false
			id       = ""
		)

		if strings.Contains(search.Query, "quotePiece:") {
			idSearch = true
			id = strings.TrimSpace(strings.TrimPrefix(search.Query, "quotePiece:"))
			id = uid.DeShortID(id)
			for _, r := range id {
				if !unicode.IsDigit(r) {
					idSearch = false
					break
				}
			}
		}

		if idSearch {
			session.And(builder.Eq{
				"id": id,
			})
		} else {
			session.And(builder.Like{
				"title", search.Query,
			})
		}
	}

	offset := search.Page * search.PageSize

	session.OrderBy("created_at desc").
		Limit(search.PageSize, offset)
	count, err = session.FindAndCount(&rows)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		return rows, count, err
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range rows {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return rows, count, nil
}

// UpdateSearch update search, if search plugin not enable, do nothing
func (qr *quotePieceRepo) UpdateSearch(ctx context.Context, quotePieceID string) (err error) {
	// check search plugin
	var s plugin.Search
	_ = plugin.CallSearch(func(search plugin.Search) error {
		s = search
		return nil
	})
	if s == nil {
		return
	}
	quotePieceID = uid.DeShortID(quotePieceID)
	quotePiece, exist, err := qr.GetQuotePiece(ctx, quotePieceID)
	if !exist {
		return
	}
	if err != nil {
		return err
	}

	// get tags
	var (
		tagListList = make([]*entity.TagRel, 0)
		tags        = make([]string, 0)
	)
	session := qr.data.DB.Context(ctx).Where("object_id = ?", quotePieceID)
	session.Where("status = ?", entity.TagRelStatusAvailable)
	err = session.Find(&tagListList)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	for _, tag := range tagListList {
		tags = append(tags, tag.TagID)
	}
	content := &plugin.SearchContent{
		ObjectID:    quotePieceID,
		Title:       quotePiece.Title,
		Type:        constant.QuotePieceObjectType,
		Content:     quotePiece.OriginalText, //.OriginalText,
		Answers:     int64(0),                // int64(quotePiece.AnswerCount),
		Status:      plugin.SearchContentStatus(quotePiece.Status),
		Tags:        tags,
		QuestionID:  quotePieceID,
		UserID:      quotePiece.UserID,
		Views:       int64(quotePiece.ViewCount),
		Created:     quotePiece.CreatedAt.Unix(),
		Active:      quotePiece.UpdatedAt.Unix(),
		Score:       int64(quotePiece.VoteCount),
		HasAccepted: true, //quotePiece.AcceptedAnswerID != "" && quotePiece.AcceptedAnswerID != "0",
	}
	err = s.UpdateContent(ctx, content)
	return
}

func (qr *quotePieceRepo) RemoveAllUserQuotePiece(ctx context.Context, userID string) (err error) {
	// get all quotePiece id that need to be deleted
	quotePieceIDs := make([]string, 0)
	session := qr.data.DB.Context(ctx).Where("user_id = ?", userID)
	session.Where("status != ?", entity.QuotePieceStatusDeleted)
	err = session.Select("id").Table("quotePiece").Find(&quotePieceIDs)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if len(quotePieceIDs) == 0 {
		return nil
	}

	log.Infof("find %d quotePieces need to be deleted for user %s", len(quotePieceIDs), userID)

	// delete all quotePiece
	session = qr.data.DB.Context(ctx).Where("user_id = ?", userID)
	session.Where("status != ?", entity.QuotePieceStatusDeleted)
	_, err = session.Cols("status", "updated_at").Update(&entity.QuotePiece{
		UpdatedAt: time.Now(),
		Status:    entity.QuotePieceStatusDeleted,
	})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	// update search content
	for _, id := range quotePieceIDs {
		_ = qr.UpdateSearch(ctx, id)
	}
	return nil
}
