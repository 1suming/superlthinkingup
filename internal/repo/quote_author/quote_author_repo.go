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

package quote_author

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

// quoteAuthorRepo quoteAuthor repository
type quoteAuthorRepo struct {
	data         *data.Data
	uniqueIDRepo unique.UniqueIDRepo
}

// NewQuoteAuthorRepo new repository
func NewQuoteAuthorRepo(
	data *data.Data,
	uniqueIDRepo unique.UniqueIDRepo,
) quote_common.QuoteAuthorRepo {
	return &quoteAuthorRepo{
		data:         data,
		uniqueIDRepo: uniqueIDRepo,
	}
}

// AddQuoteAuthor add quoteAuthor
func (qr *quoteAuthorRepo) AddQuoteAuthor(ctx context.Context, quoteAuthor *entity.QuoteAuthor) (err error) {
	quoteAuthor.ID, err = qr.uniqueIDRepo.GenUniqueIDStr(ctx, quoteAuthor.TableName())
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_, err = qr.data.DB.Context(ctx).Insert(quoteAuthor)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		quoteAuthor.ID = uid.EnShortID(quoteAuthor.ID)
	}
	return
}

// RemoveQuoteAuthor delete quoteAuthor
func (qr *quoteAuthorRepo) RemoveQuoteAuthor(ctx context.Context, id string) (err error) {
	id = uid.DeShortID(id)
	_, err = qr.data.DB.Context(ctx).Where("id =?", id).Delete(&entity.QuoteAuthor{})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// UpdateQuoteAuthor update quoteAuthor
func (qr *quoteAuthorRepo) UpdateQuoteAuthor(ctx context.Context, quoteAuthor *entity.QuoteAuthor, Cols []string) (err error) {
	quoteAuthor.ID = uid.DeShortID(quoteAuthor.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quoteAuthor.ID).Cols(Cols...).Update(quoteAuthor)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		quoteAuthor.ID = uid.EnShortID(quoteAuthor.ID)
	}
	_ = qr.UpdateSearch(ctx, quoteAuthor.ID)
	return
}

func (qr *quoteAuthorRepo) UpdatePvCount(ctx context.Context, quoteAuthorID string) (err error) {
	quoteAuthorID = uid.DeShortID(quoteAuthorID)
	quoteAuthor := &entity.QuoteAuthor{}
	_, err = qr.data.DB.Context(ctx).Where("id =?", quoteAuthorID).Incr("view_count", 1).Update(quoteAuthor)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quoteAuthor.ID)
	return nil
}

//func (qr *quoteAuthorRepo) UpdateAnswerCount(ctx context.Context, quoteAuthorID string, num int) (err error) {
//	quoteAuthorID = uid.DeShortID(quoteAuthorID)
//	quoteAuthor := &entity.QuoteAuthor{}
//	quoteAuthor.AnswerCount = num
//	_, err = qr.data.DB.Context(ctx).Where("id =?", quoteAuthorID).Cols("answer_count").Update(quoteAuthor)
//	if err != nil {
//		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
//	}
//	_ = qr.UpdateSearch(ctx, quoteAuthor.ID)
//	return nil
//}

func (qr *quoteAuthorRepo) UpdateCollectionCount(ctx context.Context, quoteAuthorID string) (count int64, err error) {
	quoteAuthorID = uid.DeShortID(quoteAuthorID)
	_, err = qr.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)
		count, err = session.Count(&entity.Collection{ObjectID: quoteAuthorID})
		if err != nil {
			return nil, err
		}

		quoteAuthor := &entity.QuoteAuthor{CollectionCount: int(count)}
		_, err = session.ID(quoteAuthorID).MustCols("collection_count").Update(quoteAuthor)
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

func (qr *quoteAuthorRepo) UpdateQuoteAuthorStatus(ctx context.Context, quoteAuthorID string, status int) (err error) {
	quoteAuthorID = uid.DeShortID(quoteAuthorID)
	_, err = qr.data.DB.Context(ctx).ID(quoteAuthorID).Cols("status").Update(&entity.QuoteAuthor{Status: status})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quoteAuthorID)
	return nil
}

func (qr *quoteAuthorRepo) UpdateQuoteAuthorStatusWithOutUpdateTime(ctx context.Context, quoteAuthor *entity.QuoteAuthor) (err error) {
	quoteAuthor.ID = uid.DeShortID(quoteAuthor.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quoteAuthor.ID).Cols("status").Update(quoteAuthor)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quoteAuthor.ID)
	return nil
}

func (qr *quoteAuthorRepo) RecoverQuoteAuthor(ctx context.Context, quoteAuthorID string) (err error) {
	quoteAuthorID = uid.DeShortID(quoteAuthorID)
	_, err = qr.data.DB.Context(ctx).ID(quoteAuthorID).Cols("status").Update(&entity.QuoteAuthor{Status: entity.QuoteAuthorStatusAvailable})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quoteAuthorID)
	return nil
}

func (qr *quoteAuthorRepo) UpdateQuoteAuthorOperation(ctx context.Context, quoteAuthor *entity.QuoteAuthor) (err error) {
	quoteAuthor.ID = uid.DeShortID(quoteAuthor.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quoteAuthor.ID).Cols("pin", "show").Update(quoteAuthor)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (qr *quoteAuthorRepo) UpdateAccepted(ctx context.Context, quoteAuthor *entity.QuoteAuthor) (err error) {
	quoteAuthor.ID = uid.DeShortID(quoteAuthor.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quoteAuthor.ID).Cols("accepted_answer_id").Update(quoteAuthor)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quoteAuthor.ID)
	return nil
}

func (qr *quoteAuthorRepo) UpdateLastAnswer(ctx context.Context, quoteAuthor *entity.QuoteAuthor) (err error) {
	quoteAuthor.ID = uid.DeShortID(quoteAuthor.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quoteAuthor.ID).Cols("last_answer_id").Update(quoteAuthor)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quoteAuthor.ID)
	return nil
}

// GetQuoteAuthor get quoteAuthor one
func (qr *quoteAuthorRepo) GetQuoteAuthor(ctx context.Context, id string) (
	quoteAuthor *entity.QuoteAuthor, exist bool, err error,
) {
	id = uid.DeShortID(id)
	quoteAuthor = &entity.QuoteAuthor{}
	quoteAuthor.ID = id
	exist, err = qr.data.DB.Context(ctx).Where("id = ?", id).Get(quoteAuthor)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		quoteAuthor.ID = uid.EnShortID(quoteAuthor.ID)
	}
	return
}

// GetQuoteAuthorsByTitle get quoteAuthor list by title
func (qr *quoteAuthorRepo) GetQuoteAuthorsByAuthorName(ctx context.Context, title string, pageSize int) (
	quoteAuthorList []*entity.QuoteAuthor, err error) {
	quoteAuthorList = make([]*entity.QuoteAuthor, 0)
	session := qr.data.DB.Context(ctx)
	session.Where("status != ?", entity.QuoteAuthorStatusDeleted)
	session.Where("author_name like ?", "%"+title+"%")
	session.Limit(pageSize)
	err = session.Find(&quoteAuthorList)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range quoteAuthorList {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return
}
func (qr *quoteAuthorRepo) GetQuoteAuthorByAuthorName(ctx context.Context, title string) (
	quoteAuthor *entity.QuoteAuthor, err error) {
	//quoteAuthorList = make([]*entity.QuoteAuthor, 0)
	quoteAuthorList := make([]*entity.QuoteAuthor, 0)
	session := qr.data.DB.Context(ctx)
	session.Where("status != ?", entity.QuoteAuthorStatusDeleted)
	session.Where("author_name = ?", title) //不要用like，用等于
	//session.Limit(pageSize)
	err = session.Find(&quoteAuthorList)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range quoteAuthorList {
			item.ID = uid.EnShortID(item.ID)
			quoteAuthor = item
			return //@只要一条
		}
	}
	return
}
func (qr *quoteAuthorRepo) FindByID(ctx context.Context, id []string) (quoteAuthorList []*entity.QuoteAuthor, err error) {
	for key, itemID := range id {
		id[key] = uid.DeShortID(itemID)
	}
	quoteAuthorList = make([]*entity.QuoteAuthor, 0)
	err = qr.data.DB.Context(ctx).Table("quoteAuthor").In("id", id).Find(&quoteAuthorList)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range quoteAuthorList {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return
}

// GetQuoteAuthorList get quoteAuthor list all
func (qr *quoteAuthorRepo) GetQuoteAuthorList(ctx context.Context, quoteAuthor *entity.QuoteAuthor) (quoteAuthorList []*entity.QuoteAuthor, err error) {
	quoteAuthor.ID = uid.DeShortID(quoteAuthor.ID)
	quoteAuthorList = make([]*entity.QuoteAuthor, 0)
	err = qr.data.DB.Context(ctx).Find(quoteAuthorList, quoteAuthor)
	if err != nil {
		return quoteAuthorList, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	for _, item := range quoteAuthorList {
		item.ID = uid.DeShortID(item.ID)
	}
	return
}

func (qr *quoteAuthorRepo) GetQuoteAuthorCount(ctx context.Context) (count int64, err error) {
	session := qr.data.DB.Context(ctx)
	session.Where(builder.Lt{"status": entity.QuoteAuthorStatusDeleted})
	count, err = session.Count(&entity.QuoteAuthor{Show: entity.QuoteAuthorShow})
	if err != nil {
		return 0, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return count, nil
}

func (qr *quoteAuthorRepo) GetUserQuoteAuthorCount(ctx context.Context, userID string, show int) (count int64, err error) {
	session := qr.data.DB.Context(ctx)
	session.Where(builder.Lt{"status": entity.QuoteAuthorStatusDeleted})
	count, err = session.Count(&entity.QuoteAuthor{UserID: userID, Show: show})
	if err != nil {
		return count, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (qr *quoteAuthorRepo) SitemapQuoteAuthors(ctx context.Context, page, pageSize int) (
	quoteAuthorIDList []*schema.SiteMapQuoteAuthorInfo, err error) {
	page = page - 1
	quoteAuthorIDList = make([]*schema.SiteMapQuoteAuthorInfo, 0)

	// try to get sitemap data from cache
	cacheKey := fmt.Sprintf(constant.SiteMapQuoteAuthorCacheKeyPrefix, page)
	cacheData, exist, err := qr.data.Cache.GetString(ctx, cacheKey)
	if err == nil && exist {
		_ = json.Unmarshal([]byte(cacheData), &quoteAuthorIDList)
		return quoteAuthorIDList, nil
	}

	// get sitemap data from db
	rows := make([]*entity.QuoteAuthor, 0)
	session := qr.data.DB.Context(ctx)
	//session.Select("id,title,created_at,post_update_time")
	session.Select("id,title,created_at,update_at")

	session.Where("`show` = ?", entity.QuoteAuthorShow)
	session.Where("status = ? OR status = ?", entity.QuoteAuthorStatusAvailable, entity.QuoteAuthorStatusClosed)
	session.Limit(pageSize, page*pageSize)
	session.Asc("created_at")
	err = session.Find(&rows)
	if err != nil {
		return quoteAuthorIDList, err
	}

	// warp data
	for _, quoteAuthor := range rows {
		item := &schema.SiteMapQuoteAuthorInfo{ID: quoteAuthor.ID}
		if handler.GetEnableShortID(ctx) {
			item.ID = uid.EnShortID(quoteAuthor.ID)
		}
		item.AuthorName = htmltext.UrlTitle(quoteAuthor.AuthorName)
		//if quoteAuthor.PostUpdateTime.IsZero() {
		//	item.UpdateTime = quoteAuthor.CreatedAt.Format(time.RFC3339)
		//} else {
		//	item.UpdateTime = quoteAuthor.PostUpdateTime.Format(time.RFC3339)
		//}

		quoteAuthorIDList = append(quoteAuthorIDList, item)
	}

	// set sitemap data to cache
	cacheDataByte, _ := json.Marshal(quoteAuthorIDList)
	if err := qr.data.Cache.SetString(ctx, cacheKey, string(cacheDataByte), constant.SiteMapQuoteAuthorCacheTime); err != nil {
		log.Error(err)
	}
	return quoteAuthorIDList, nil
}

// GetQuoteAuthorPage query quoteAuthor page
func (qr *quoteAuthorRepo) GetQuoteAuthorPage(ctx context.Context, page, pageSize int,
	tagIDs []string, userID, orderCond string, inDays int, showHidden, showPending bool) (
	quoteAuthorList []*entity.QuoteAuthor, total int64, err error) {
	quoteAuthorList = make([]*entity.QuoteAuthor, 0)
	session := qr.data.DB.Context(ctx)
	status := []int{entity.QuoteAuthorStatusAvailable, entity.QuoteAuthorStatusClosed}
	if showPending {
		status = append(status, entity.QuoteAuthorStatusPending)
	}
	session.In("quoteAuthor.status", status)
	if len(tagIDs) > 0 {
		session.Join("LEFT", "tag_rel", "quoteAuthor.id = tag_rel.object_id")
		session.In("tag_rel.tag_id", tagIDs)
		session.And("tag_rel.status = ?", entity.TagRelStatusAvailable)
	}
	if len(userID) > 0 {
		session.And("quoteAuthor.user_id = ?", userID)
		if !showHidden {
			session.And("quoteAuthor.show = ?", entity.QuoteAuthorShow)
		}
	} else {
		session.And("quoteAuthor.show = ?", entity.QuoteAuthorShow)
	}
	if inDays > 0 {
		session.And("quoteAuthor.created_at > ?", time.Now().AddDate(0, 0, -inDays))
	}

	switch orderCond {
	case "newest":
		session.OrderBy("quoteAuthor.pin desc,quoteAuthor.created_at DESC")
	case "active":
		if inDays == 0 {
			session.And("quoteAuthor.created_at > ?", time.Now().AddDate(0, 0, -180))
		}
		session.And("quoteAuthor.post_update_time > ?", time.Now().AddDate(0, 0, -90))
		session.OrderBy("quoteAuthor.pin desc,quoteAuthor.post_update_time DESC, quoteAuthor.updated_at DESC")
	case "hot":
		session.OrderBy("quoteAuthor.pin desc,quoteAuthor.hot_score DESC")
	case "score":
		session.OrderBy("quoteAuthor.pin desc,quoteAuthor.vote_count DESC, quoteAuthor.view_count DESC")
	case "unanswered":
		session.Where("quoteAuthor.answer_count = 0")
		session.OrderBy("quoteAuthor.pin desc,quoteAuthor.created_at DESC")
	}

	total, err = pager.Help(page, pageSize, &quoteAuthorList, &entity.QuoteAuthor{}, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range quoteAuthorList {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return quoteAuthorList, total, err
}

// GetRecommendQuoteAuthorPageByTags get recommend quoteAuthor page by tags
func (qr *quoteAuthorRepo) GetRecommendQuoteAuthorPageByTags(ctx context.Context, userID string, tagIDs, followedQuoteAuthorIDs []string, page, pageSize int) (
	quoteAuthorList []*entity.QuoteAuthor, total int64, err error) {
	quoteAuthorList = make([]*entity.QuoteAuthor, 0)
	orderBySQL := "quoteAuthor.pin DESC, quoteAuthor.created_at DESC"

	// Please Make sure every quoteAuthor has at least one tag
	selectSQL := entity.QuoteAuthor{}.TableName() + ".*"
	if len(followedQuoteAuthorIDs) > 0 {
		idStr := "'" + strings.Join(followedQuoteAuthorIDs, "','") + "'"
		selectSQL += fmt.Sprintf(", CASE WHEN quoteAuthor.id IN (%s) THEN 0 ELSE 1 END AS order_priority", idStr)
		orderBySQL = "order_priority, " + orderBySQL
	}
	session := qr.data.DB.Context(ctx).Select(selectSQL)

	if len(tagIDs) > 0 {
		session.Where("quoteAuthor.user_id != ?", userID).
			And("quoteAuthor.id NOT IN (SELECT quoteAuthor_id FROM answer WHERE user_id = ?)", userID).
			Join("INNER", "tag_rel", "quoteAuthor.id = tag_rel.object_id").
			And("tag_rel.status = ?", entity.TagRelStatusAvailable).
			Join("INNER", "tag", "tag.id = tag_rel.tag_id").
			In("tag.id", tagIDs)
	} else if len(followedQuoteAuthorIDs) == 0 {
		return quoteAuthorList, 0, nil
	}

	if len(followedQuoteAuthorIDs) > 0 {
		if len(tagIDs) > 0 {
			// if tags provided, show followed quoteAuthors and tag quoteAuthors
			session.Or(builder.In("quoteAuthor.id", followedQuoteAuthorIDs))
		} else {
			// if no tags, only show followed quoteAuthors
			session.Where(builder.In("quoteAuthor.id", followedQuoteAuthorIDs))
		}
	}

	session.
		And("quoteAuthor.show = ? and quoteAuthor.status = ?", entity.QuoteAuthorShow, entity.QuoteAuthorStatusAvailable).
		Distinct("quoteAuthor.id").
		OrderBy(orderBySQL)

	total, err = pager.Help(page, pageSize, &quoteAuthorList, &entity.QuoteAuthor{}, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	if handler.GetEnableShortID(ctx) {
		for _, item := range quoteAuthorList {
			item.ID = uid.EnShortID(item.ID)
		}
	}

	return quoteAuthorList, total, err
}

func (qr *quoteAuthorRepo) AdminQuoteAuthorPage(ctx context.Context, search *schema.AdminQuoteAuthorPageReq) ([]*entity.QuoteAuthor, int64, error) {
	var (
		count   int64
		err     error
		session = qr.data.DB.Context(ctx).Table("quoteAuthor")
	)

	session.Where(builder.Eq{
		"status": search.Status,
	})

	rows := make([]*entity.QuoteAuthor, 0)
	if search.Page > 0 {
		search.Page = search.Page - 1
	} else {
		search.Page = 0
	}
	if search.PageSize == 0 {
		search.PageSize = constant.DefaultPageSize
	}

	// search by quoteAuthor author_name like or quoteAuthor id
	if len(search.Query) > 0 {
		// check id search
		var (
			idSearch = false
			id       = ""
		)

		if strings.Contains(search.Query, "quoteAuthor:") {
			idSearch = true
			id = strings.TrimSpace(strings.TrimPrefix(search.Query, "quoteAuthor:"))
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
func (qr *quoteAuthorRepo) UpdateSearch(ctx context.Context, quoteAuthorID string) (err error) {
	// check search plugin
	var s plugin.Search
	_ = plugin.CallSearch(func(search plugin.Search) error {
		s = search
		return nil
	})
	if s == nil {
		return
	}
	quoteAuthorID = uid.DeShortID(quoteAuthorID)
	quoteAuthor, exist, err := qr.GetQuoteAuthor(ctx, quoteAuthorID)
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
	session := qr.data.DB.Context(ctx).Where("object_id = ?", quoteAuthorID)
	session.Where("status = ?", entity.TagRelStatusAvailable)
	err = session.Find(&tagListList)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	for _, tag := range tagListList {
		tags = append(tags, tag.TagID)
	}
	content := &plugin.SearchContent{
		ObjectID:    quoteAuthorID,
		Title:       quoteAuthor.AuthorName,
		Type:        constant.QuoteAuthorObjectType,
		Content:     quoteAuthor.Bio, //.OriginalText,
		Answers:     int64(0),        // int64(quoteAuthor.AnswerCount),
		Status:      plugin.SearchContentStatus(quoteAuthor.Status),
		Tags:        tags,
		QuestionID:  quoteAuthorID,
		UserID:      quoteAuthor.UserID,
		Views:       int64(quoteAuthor.ViewCount),
		Created:     quoteAuthor.CreatedAt.Unix(),
		Active:      quoteAuthor.UpdatedAt.Unix(),
		Score:       int64(quoteAuthor.VoteCount),
		HasAccepted: true, //quoteAuthor.AcceptedAnswerID != "" && quoteAuthor.AcceptedAnswerID != "0",
	}
	err = s.UpdateContent(ctx, content)
	return
}

func (qr *quoteAuthorRepo) RemoveAllUserQuoteAuthor(ctx context.Context, userID string) (err error) {
	// get all quoteAuthor id that need to be deleted
	quoteAuthorIDs := make([]string, 0)
	session := qr.data.DB.Context(ctx).Where("user_id = ?", userID)
	session.Where("status != ?", entity.QuoteAuthorStatusDeleted)
	err = session.Select("id").Table("quoteAuthor").Find(&quoteAuthorIDs)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if len(quoteAuthorIDs) == 0 {
		return nil
	}

	log.Infof("find %d quoteAuthors need to be deleted for user %s", len(quoteAuthorIDs), userID)

	// delete all quoteAuthor
	session = qr.data.DB.Context(ctx).Where("user_id = ?", userID)
	session.Where("status != ?", entity.QuoteAuthorStatusDeleted)
	_, err = session.Cols("status", "updated_at").Update(&entity.QuoteAuthor{
		UpdatedAt: time.Now(),
		Status:    entity.QuoteAuthorStatusDeleted,
	})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	// update search content
	for _, id := range quoteAuthorIDs {
		_ = qr.UpdateSearch(ctx, id)
	}
	return nil
}
