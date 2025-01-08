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

package quote

import (
	"context"
	"encoding/json"
	"fmt"
	quotecommon "github.com/apache/incubator-answer/internal/service_quote/quote_common"
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

// quoteRepo quote repository
type quoteRepo struct {
	data         *data.Data
	uniqueIDRepo unique.UniqueIDRepo
}

// NewQuoteRepo new repository
func NewQuoteRepo(
	data *data.Data,
	uniqueIDRepo unique.UniqueIDRepo,
) quotecommon.QuoteRepo {
	return &quoteRepo{
		data:         data,
		uniqueIDRepo: uniqueIDRepo,
	}
}

// AddQuote add quote
func (qr *quoteRepo) AddQuote(ctx context.Context, quote *entity.Quote) (err error) {
	quote.ID, err = qr.uniqueIDRepo.GenUniqueIDStr(ctx, quote.TableName())
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_, err = qr.data.DB.Context(ctx).Insert(quote)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		quote.ID = uid.EnShortID(quote.ID)
	}
	return
}

// RemoveQuote delete quote
func (qr *quoteRepo) RemoveQuote(ctx context.Context, id string) (err error) {
	id = uid.DeShortID(id)
	_, err = qr.data.DB.Context(ctx).Where("id =?", id).Delete(&entity.Quote{})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// UpdateQuote update quote
func (qr *quoteRepo) UpdateQuote(ctx context.Context, quote *entity.Quote, Cols []string) (err error) {
	quote.ID = uid.DeShortID(quote.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quote.ID).Cols(Cols...).Update(quote)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		quote.ID = uid.EnShortID(quote.ID)
	}
	_ = qr.UpdateSearch(ctx, quote.ID)
	return
}

func (qr *quoteRepo) UpdatePvCount(ctx context.Context, quoteID string) (err error) {
	quoteID = uid.DeShortID(quoteID)
	quote := &entity.Quote{}
	_, err = qr.data.DB.Context(ctx).Where("id =?", quoteID).Incr("view_count", 1).Update(quote)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quote.ID)
	return nil
}

//func (qr *quoteRepo) UpdateAnswerCount(ctx context.Context, quoteID string, num int) (err error) {
//	quoteID = uid.DeShortID(quoteID)
//	quote := &entity.Quote{}
//	quote.AnswerCount = num
//	_, err = qr.data.DB.Context(ctx).Where("id =?", quoteID).Cols("answer_count").Update(quote)
//	if err != nil {
//		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
//	}
//	_ = qr.UpdateSearch(ctx, quote.ID)
//	return nil
//}

func (qr *quoteRepo) UpdateCollectionCount(ctx context.Context, quoteID string) (count int64, err error) {
	quoteID = uid.DeShortID(quoteID)
	_, err = qr.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)
		count, err = session.Count(&entity.Collection{ObjectID: quoteID})
		if err != nil {
			return nil, err
		}

		quote := &entity.Quote{CollectionCount: int(count)}
		_, err = session.ID(quoteID).MustCols("collection_count").Update(quote)
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

func (qr *quoteRepo) UpdateQuoteStatus(ctx context.Context, quoteID string, status int) (err error) {
	quoteID = uid.DeShortID(quoteID)
	_, err = qr.data.DB.Context(ctx).ID(quoteID).Cols("status").Update(&entity.Quote{Status: status})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quoteID)
	return nil
}

func (qr *quoteRepo) UpdateQuoteStatusWithOutUpdateTime(ctx context.Context, quote *entity.Quote) (err error) {
	quote.ID = uid.DeShortID(quote.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quote.ID).Cols("status").Update(quote)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quote.ID)
	return nil
}

func (qr *quoteRepo) RecoverQuote(ctx context.Context, quoteID string) (err error) {
	quoteID = uid.DeShortID(quoteID)
	_, err = qr.data.DB.Context(ctx).ID(quoteID).Cols("status").Update(&entity.Quote{Status: entity.QuoteStatusAvailable})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quoteID)
	return nil
}

func (qr *quoteRepo) UpdateQuoteOperation(ctx context.Context, quote *entity.Quote) (err error) {
	quote.ID = uid.DeShortID(quote.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quote.ID).Cols("pin", "show").Update(quote)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (qr *quoteRepo) UpdateAccepted(ctx context.Context, quote *entity.Quote) (err error) {
	quote.ID = uid.DeShortID(quote.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quote.ID).Cols("accepted_answer_id").Update(quote)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quote.ID)
	return nil
}

func (qr *quoteRepo) UpdateLastAnswer(ctx context.Context, quote *entity.Quote) (err error) {
	quote.ID = uid.DeShortID(quote.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", quote.ID).Cols("last_answer_id").Update(quote)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, quote.ID)
	return nil
}

// GetQuote get quote one
func (qr *quoteRepo) GetQuote(ctx context.Context, id string) (
	quote *entity.Quote, exist bool, err error,
) {
	id = uid.DeShortID(id)
	quote = &entity.Quote{}
	quote.ID = id
	exist, err = qr.data.DB.Context(ctx).Where("id = ?", id).Get(quote)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		quote.ID = uid.EnShortID(quote.ID)
	}
	return
}

// GetQuotesByTitle get quote list by title
func (qr *quoteRepo) GetQuotesByTitle(ctx context.Context, title string, pageSize int) (
	quoteList []*entity.Quote, err error) {
	quoteList = make([]*entity.Quote, 0)
	session := qr.data.DB.Context(ctx)
	session.Where("status != ?", entity.QuoteStatusDeleted)
	session.Where("title like ?", "%"+title+"%")
	session.Limit(pageSize)
	err = session.Find(&quoteList)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range quoteList {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return
}

func (qr *quoteRepo) FindByID(ctx context.Context, id []string) (quoteList []*entity.Quote, err error) {
	for key, itemID := range id {
		id[key] = uid.DeShortID(itemID)
	}
	quoteList = make([]*entity.Quote, 0)
	err = qr.data.DB.Context(ctx).Table("quote").In("id", id).Find(&quoteList)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range quoteList {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return
}

// GetQuoteList get quote list all
func (qr *quoteRepo) GetQuoteList(ctx context.Context, quote *entity.Quote) (quoteList []*entity.Quote, err error) {
	quote.ID = uid.DeShortID(quote.ID)
	quoteList = make([]*entity.Quote, 0)
	err = qr.data.DB.Context(ctx).Find(quoteList, quote)
	if err != nil {
		return quoteList, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	for _, item := range quoteList {
		item.ID = uid.DeShortID(item.ID)
	}
	return
}

func (qr *quoteRepo) GetQuoteCount(ctx context.Context) (count int64, err error) {
	session := qr.data.DB.Context(ctx)
	session.Where(builder.Lt{"status": entity.QuoteStatusDeleted})
	count, err = session.Count(&entity.Quote{Show: entity.QuoteShow})
	if err != nil {
		return 0, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return count, nil
}

func (qr *quoteRepo) GetUserQuoteCount(ctx context.Context, userID string, show int) (count int64, err error) {
	session := qr.data.DB.Context(ctx)
	session.Where(builder.Lt{"status": entity.QuoteStatusDeleted})
	count, err = session.Count(&entity.Quote{UserID: userID, Show: show})
	if err != nil {
		return count, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (qr *quoteRepo) SitemapQuotes(ctx context.Context, page, pageSize int) (
	quoteIDList []*schema.SiteMapQuoteInfo, err error) {
	page = page - 1
	quoteIDList = make([]*schema.SiteMapQuoteInfo, 0)

	// try to get sitemap data from cache
	cacheKey := fmt.Sprintf(constant.SiteMapQuoteCacheKeyPrefix, page)
	cacheData, exist, err := qr.data.Cache.GetString(ctx, cacheKey)
	if err == nil && exist {
		_ = json.Unmarshal([]byte(cacheData), &quoteIDList)
		return quoteIDList, nil
	}

	// get sitemap data from db
	rows := make([]*entity.Quote, 0)
	session := qr.data.DB.Context(ctx)
	//session.Select("id,title,created_at,post_update_time")
	session.Select("id,title,created_at,update_at")

	session.Where("`show` = ?", entity.QuoteShow)
	session.Where("status = ? OR status = ?", entity.QuoteStatusAvailable, entity.QuoteStatusClosed)
	session.Limit(pageSize, page*pageSize)
	session.Asc("created_at")
	err = session.Find(&rows)
	if err != nil {
		return quoteIDList, err
	}

	// warp data
	for _, quote := range rows {
		item := &schema.SiteMapQuoteInfo{ID: quote.ID}
		if handler.GetEnableShortID(ctx) {
			item.ID = uid.EnShortID(quote.ID)
		}
		item.Title = htmltext.UrlTitle(quote.Title)
		if quote.PostUpdateTime.IsZero() {
			item.UpdateTime = quote.CreatedAt.Format(time.RFC3339)
		} else {
			item.UpdateTime = quote.PostUpdateTime.Format(time.RFC3339)
		}

		quoteIDList = append(quoteIDList, item)
	}

	// set sitemap data to cache
	cacheDataByte, _ := json.Marshal(quoteIDList)
	if err := qr.data.Cache.SetString(ctx, cacheKey, string(cacheDataByte), constant.SiteMapQuoteCacheTime); err != nil {
		log.Error(err)
	}
	return quoteIDList, nil
}

// GetQuotePage query quote page
func (qr *quoteRepo) GetQuotePage(ctx context.Context, page, pageSize int,
	tagIDs []string, userID, orderCond string, inDays int, showHidden, showPending bool) (
	quoteList []*entity.Quote, total int64, err error) {
	quoteList = make([]*entity.Quote, 0)

	session := qr.data.DB.Context(ctx)
	session.Alias(entity.QuoteGetAlias()) //@ms:alias
	status := []int{entity.QuoteStatusAvailable, entity.QuoteStatusClosed}
	if showPending {
		status = append(status, entity.QuoteStatusPending)
	}
	session.In("quote.status", status)
	if len(tagIDs) > 0 {
		session.Join("LEFT", "tag_rel", "quote.id = tag_rel.object_id")
		session.In("tag_rel.tag_id", tagIDs)
		session.And("tag_rel.status = ?", entity.TagRelStatusAvailable)
	}
	session.Join("LEFT", "tq_quote_author", "quote.quote_author_id = tq_quote_author.id")
	session.Join("LEFT", "tq_quote_piece", "quote.quote_piece_id = tq_quote_piece.id")

	if len(userID) > 0 {
		session.And("quote.user_id = ?", userID)
		if !showHidden {
			session.And("quote.show = ?", entity.QuoteShow)
		}
	} else {
		session.And("quote.show = ?", entity.QuoteShow)
	}
	if inDays > 0 {
		session.And("quote.created_at > ?", time.Now().AddDate(0, 0, -inDays))
	}

	switch orderCond {
	case "newest":
		session.OrderBy("quote.pin desc,quote.created_at DESC")
	case "active":
		if inDays == 0 {
			session.And("quote.created_at > ?", time.Now().AddDate(0, 0, -180))
		}
		session.And("quote.post_update_time > ?", time.Now().AddDate(0, 0, -90))
		session.OrderBy("quote.pin desc,quote.post_update_time DESC, quote.updated_at DESC")
	case "hot":
		session.OrderBy("quote.pin desc,quote.hot_score DESC")
	case "score":
		session.OrderBy("quote.pin desc,quote.vote_count DESC, quote.view_count DESC")
	case "unanswered":
		session.Where("quote.answer_count = 0")
		session.OrderBy("quote.pin desc,quote.created_at DESC")
	}

	total, err = pager.Help(page, pageSize, &quoteList, &entity.Quote{}, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range quoteList {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return quoteList, total, err
}

// GetRecommendQuotePageByTags get recommend quote page by tags
func (qr *quoteRepo) GetRecommendQuotePageByTags(ctx context.Context, userID string, tagIDs, followedQuoteIDs []string, page, pageSize int) (
	quoteList []*entity.Quote, total int64, err error) {
	quoteList = make([]*entity.Quote, 0)
	orderBySQL := "quote.pin DESC, quote.created_at DESC"

	// Please Make sure every quote has at least one tag
	selectSQL := entity.Quote{}.TableName() + ".*"
	if len(followedQuoteIDs) > 0 {
		idStr := "'" + strings.Join(followedQuoteIDs, "','") + "'"
		selectSQL += fmt.Sprintf(", CASE WHEN quote.id IN (%s) THEN 0 ELSE 1 END AS order_priority", idStr)
		orderBySQL = "order_priority, " + orderBySQL
	}
	session := qr.data.DB.Context(ctx).Select(selectSQL)

	if len(tagIDs) > 0 {
		session.Where("quote.user_id != ?", userID).
			And("quote.id NOT IN (SELECT quote_id FROM answer WHERE user_id = ?)", userID).
			Join("INNER", "tag_rel", "quote.id = tag_rel.object_id").
			And("tag_rel.status = ?", entity.TagRelStatusAvailable).
			Join("INNER", "tag", "tag.id = tag_rel.tag_id").
			In("tag.id", tagIDs)
	} else if len(followedQuoteIDs) == 0 {
		return quoteList, 0, nil
	}

	if len(followedQuoteIDs) > 0 {
		if len(tagIDs) > 0 {
			// if tags provided, show followed quotes and tag quotes
			session.Or(builder.In("quote.id", followedQuoteIDs))
		} else {
			// if no tags, only show followed quotes
			session.Where(builder.In("quote.id", followedQuoteIDs))
		}
	}

	session.
		And("quote.show = ? and quote.status = ?", entity.QuoteShow, entity.QuoteStatusAvailable).
		Distinct("quote.id").
		OrderBy(orderBySQL)

	total, err = pager.Help(page, pageSize, &quoteList, &entity.Quote{}, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	if handler.GetEnableShortID(ctx) {
		for _, item := range quoteList {
			item.ID = uid.EnShortID(item.ID)
		}
	}

	return quoteList, total, err
}

func (qr *quoteRepo) AdminQuotePage(ctx context.Context, search *schema.AdminQuotePageReq) ([]*entity.Quote, int64, error) {
	var (
		count   int64
		err     error
		session = qr.data.DB.Context(ctx).Table("quote")
	)

	session.Where(builder.Eq{
		"status": search.Status,
	})

	rows := make([]*entity.Quote, 0)
	if search.Page > 0 {
		search.Page = search.Page - 1
	} else {
		search.Page = 0
	}
	if search.PageSize == 0 {
		search.PageSize = constant.DefaultPageSize
	}

	// search by quote title like or quote id
	if len(search.Query) > 0 {
		// check id search
		var (
			idSearch = false
			id       = ""
		)

		if strings.Contains(search.Query, "quote:") {
			idSearch = true
			id = strings.TrimSpace(strings.TrimPrefix(search.Query, "quote:"))
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
func (qr *quoteRepo) UpdateSearch(ctx context.Context, quoteID string) (err error) {
	// check search plugin
	var s plugin.Search
	_ = plugin.CallSearch(func(search plugin.Search) error {
		s = search
		return nil
	})
	if s == nil {
		return
	}
	quoteID = uid.DeShortID(quoteID)
	quote, exist, err := qr.GetQuote(ctx, quoteID)
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
	session := qr.data.DB.Context(ctx).Where("object_id = ?", quoteID)
	session.Where("status = ?", entity.TagRelStatusAvailable)
	err = session.Find(&tagListList)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	for _, tag := range tagListList {
		tags = append(tags, tag.TagID)
	}
	content := &plugin.SearchContent{
		ObjectID:    quoteID,
		Title:       quote.Title,
		Type:        constant.QuoteObjectType,
		Content:     quote.OriginalText,
		Answers:     int64(0), // int64(quote.AnswerCount),
		Status:      plugin.SearchContentStatus(quote.Status),
		Tags:        tags,
		QuestionID:  quoteID,
		UserID:      quote.UserID,
		Views:       int64(quote.ViewCount),
		Created:     quote.CreatedAt.Unix(),
		Active:      quote.UpdatedAt.Unix(),
		Score:       int64(quote.VoteCount),
		HasAccepted: true, //quote.AcceptedAnswerID != "" && quote.AcceptedAnswerID != "0",
	}
	err = s.UpdateContent(ctx, content)
	return
}

func (qr *quoteRepo) RemoveAllUserQuote(ctx context.Context, userID string) (err error) {
	// get all quote id that need to be deleted
	quoteIDs := make([]string, 0)
	session := qr.data.DB.Context(ctx).Where("user_id = ?", userID)
	session.Where("status != ?", entity.QuoteStatusDeleted)
	err = session.Select("id").Table("quote").Find(&quoteIDs)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if len(quoteIDs) == 0 {
		return nil
	}

	log.Infof("find %d quotes need to be deleted for user %s", len(quoteIDs), userID)

	// delete all quote
	session = qr.data.DB.Context(ctx).Where("user_id = ?", userID)
	session.Where("status != ?", entity.QuoteStatusDeleted)
	_, err = session.Cols("status", "updated_at").Update(&entity.Quote{
		UpdatedAt: time.Now(),
		Status:    entity.QuoteStatusDeleted,
	})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	// update search content
	for _, id := range quoteIDs {
		_ = qr.UpdateSearch(ctx, id)
	}
	return nil
}
