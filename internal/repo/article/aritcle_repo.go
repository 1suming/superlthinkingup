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

package article

import (
	"context"
	"encoding/json"
	"fmt"
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
	articlecommon "github.com/apache/incubator-answer/internal/service/article_common"
	"github.com/apache/incubator-answer/internal/service/unique"
	"github.com/apache/incubator-answer/pkg/htmltext"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/apache/incubator-answer/plugin"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
	"xorm.io/builder"
	"xorm.io/xorm"
)

// articleRepo article repository
type articleRepo struct {
	data         *data.Data
	uniqueIDRepo unique.UniqueIDRepo
}

// NewArticleRepo new repository
func NewArticleRepo(
	data *data.Data,
	uniqueIDRepo unique.UniqueIDRepo,
) articlecommon.ArticleRepo {
	return &articleRepo{
		data:         data,
		uniqueIDRepo: uniqueIDRepo,
	}
}

// AddArticle add article
func (qr *articleRepo) AddArticle(ctx context.Context, article *entity.Article) (err error) {
	article.ID, err = qr.uniqueIDRepo.GenUniqueIDStr(ctx, article.TableName())
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_, err = qr.data.DB.Context(ctx).Insert(article)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		article.ID = uid.EnShortID(article.ID)
	}
	return
}

// RemoveArticle delete article
func (qr *articleRepo) RemoveArticle(ctx context.Context, id string) (err error) {
	id = uid.DeShortID(id)
	_, err = qr.data.DB.Context(ctx).Where("id =?", id).Delete(&entity.Article{})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// UpdateArticle update article
func (qr *articleRepo) UpdateArticle(ctx context.Context, article *entity.Article, Cols []string) (err error) {
	article.ID = uid.DeShortID(article.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", article.ID).Cols(Cols...).Update(article)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		article.ID = uid.EnShortID(article.ID)
	}
	_ = qr.UpdateSearch(ctx, article.ID)
	return
}

func (qr *articleRepo) UpdatePvCount(ctx context.Context, articleID string) (err error) {
	articleID = uid.DeShortID(articleID)
	article := &entity.Article{}
	_, err = qr.data.DB.Context(ctx).Where("id =?", articleID).Incr("view_count", 1).Update(article)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, article.ID)
	return nil
}

func (qr *articleRepo) UpdateAnswerCount(ctx context.Context, articleID string, num int) (err error) {
	articleID = uid.DeShortID(articleID)
	article := &entity.Article{}
	//article.AnswerCount = num
	_, err = qr.data.DB.Context(ctx).Where("id =?", articleID).Cols("answer_count").Update(article)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, article.ID)
	return nil
}

func (qr *articleRepo) UpdateCollectionCount(ctx context.Context, articleID string) (count int64, err error) {
	articleID = uid.DeShortID(articleID)
	_, err = qr.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)
		count, err = session.Count(&entity.Collection{ObjectID: articleID})
		if err != nil {
			return nil, err
		}

		article := &entity.Article{CollectionCount: int(count)}
		_, err = session.ID(articleID).MustCols("collection_count").Update(article)
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

func (qr *articleRepo) UpdateArticleStatus(ctx context.Context, articleID string, status int) (err error) {
	articleID = uid.DeShortID(articleID)
	_, err = qr.data.DB.Context(ctx).ID(articleID).Cols("status").Update(&entity.Article{Status: status})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, articleID)
	return nil
}

func (qr *articleRepo) UpdateArticleStatusWithOutUpdateTime(ctx context.Context, article *entity.Article) (err error) {
	article.ID = uid.DeShortID(article.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", article.ID).Cols("status").Update(article)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, article.ID)
	return nil
}

func (qr *articleRepo) RecoverArticle(ctx context.Context, articleID string) (err error) {
	articleID = uid.DeShortID(articleID)
	_, err = qr.data.DB.Context(ctx).ID(articleID).Cols("status").Update(&entity.Article{Status: entity.ArticleStatusAvailable})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, articleID)
	return nil
}

func (qr *articleRepo) UpdateArticleOperation(ctx context.Context, article *entity.Article) (err error) {
	article.ID = uid.DeShortID(article.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", article.ID).Cols("pin", "show").Update(article)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (qr *articleRepo) UpdateAccepted(ctx context.Context, article *entity.Article) (err error) {
	article.ID = uid.DeShortID(article.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", article.ID).Cols("accepted_answer_id").Update(article)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, article.ID)
	return nil
}

func (qr *articleRepo) UpdateLastAnswer(ctx context.Context, article *entity.Article) (err error) {
	article.ID = uid.DeShortID(article.ID)
	_, err = qr.data.DB.Context(ctx).Where("id =?", article.ID).Cols("last_answer_id").Update(article)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_ = qr.UpdateSearch(ctx, article.ID)
	return nil
}

// GetArticle get article one
func (qr *articleRepo) GetArticle(ctx context.Context, id string) (
	article *entity.Article, exist bool, err error,
) {
	id = uid.DeShortID(id)
	article = &entity.Article{}
	article.ID = id
	exist, err = qr.data.DB.Context(ctx).Where("id = ?", id).Get(article)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		article.ID = uid.EnShortID(article.ID)
	}
	return
}

// GetArticlesByTitle get article list by title
func (qr *articleRepo) GetArticlesByTitle(ctx context.Context, title string, pageSize int) (
	articleList []*entity.Article, err error) {
	articleList = make([]*entity.Article, 0)
	session := qr.data.DB.Context(ctx)
	session.Where("status != ?", entity.ArticleStatusDeleted)
	session.Where("title like ?", "%"+title+"%")
	session.Limit(pageSize)
	err = session.Find(&articleList)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range articleList {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return
}

func (qr *articleRepo) FindByID(ctx context.Context, id []string) (articleList []*entity.Article, err error) {
	for key, itemID := range id {
		id[key] = uid.DeShortID(itemID)
	}
	articleList = make([]*entity.Article, 0)
	err = qr.data.DB.Context(ctx).Table("article").In("id", id).Find(&articleList)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range articleList {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return
}

// GetArticleList get article list all
func (qr *articleRepo) GetArticleList(ctx context.Context, article *entity.Article) (articleList []*entity.Article, err error) {
	article.ID = uid.DeShortID(article.ID)
	articleList = make([]*entity.Article, 0)
	err = qr.data.DB.Context(ctx).Find(articleList, article)
	if err != nil {
		return articleList, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	for _, item := range articleList {
		item.ID = uid.DeShortID(item.ID)
	}
	return
}

func (qr *articleRepo) GetArticleCount(ctx context.Context) (count int64, err error) {
	session := qr.data.DB.Context(ctx)
	session.Where(builder.Lt{"status": entity.ArticleStatusDeleted})
	count, err = session.Count(&entity.Article{Show: entity.ArticleShow})
	if err != nil {
		return 0, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return count, nil
}

func (qr *articleRepo) GetUserArticleCount(ctx context.Context, userID string, show int) (count int64, err error) {
	session := qr.data.DB.Context(ctx)
	session.Where(builder.Lt{"status": entity.ArticleStatusDeleted})
	count, err = session.Count(&entity.Article{UserID: userID, Show: show})
	if err != nil {
		return count, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (qr *articleRepo) SitemapArticles(ctx context.Context, page, pageSize int) (
	articleIDList []*schema.SiteMapArticleInfo, err error) {
	page = page - 1
	articleIDList = make([]*schema.SiteMapArticleInfo, 0)

	// try to get sitemap data from cache
	cacheKey := fmt.Sprintf(constant.SiteMapArticleCacheKeyPrefix, page)
	cacheData, exist, err := qr.data.Cache.GetString(ctx, cacheKey)
	if err == nil && exist {
		_ = json.Unmarshal([]byte(cacheData), &articleIDList)
		return articleIDList, nil
	}

	// get sitemap data from db
	rows := make([]*entity.Article, 0)
	session := qr.data.DB.Context(ctx)
	session.Select("id,title,created_at,post_update_time")
	session.Where("`show` = ?", entity.ArticleShow)
	session.Where("status = ? OR status = ?", entity.ArticleStatusAvailable, entity.ArticleStatusClosed)
	session.Limit(pageSize, page*pageSize)
	session.Asc("created_at")
	err = session.Find(&rows)
	if err != nil {
		return articleIDList, err
	}

	// warp data
	for _, article := range rows {
		item := &schema.SiteMapArticleInfo{ID: article.ID}
		if handler.GetEnableShortID(ctx) {
			item.ID = uid.EnShortID(article.ID)
		}
		item.Title = htmltext.UrlTitle(article.Title)
		if article.PostUpdateTime.IsZero() {
			item.UpdateTime = article.CreatedAt.Format(time.RFC3339)
		} else {
			item.UpdateTime = article.PostUpdateTime.Format(time.RFC3339)
		}
		articleIDList = append(articleIDList, item)
	}

	// set sitemap data to cache
	cacheDataByte, _ := json.Marshal(articleIDList)
	if err := qr.data.Cache.SetString(ctx, cacheKey, string(cacheDataByte), constant.SiteMapArticleCacheTime); err != nil {
		log.Error(err)
	}
	return articleIDList, nil
}

// GetArticlePage query article page
func (qr *articleRepo) GetArticlePage(ctx context.Context, page, pageSize int,
	tagIDs []string, userID, orderCond string, inDays int, showHidden, showPending bool) (
	articleList []*entity.Article, total int64, err error) {
	articleList = make([]*entity.Article, 0)
	session := qr.data.DB.Context(ctx)
	status := []int{entity.ArticleStatusAvailable, entity.ArticleStatusClosed}
	if showPending {
		status = append(status, entity.ArticleStatusPending)
	}
	session.In("article.status", status)
	if len(tagIDs) > 0 {
		session.Join("LEFT", "tag_rel", "article.id = tag_rel.object_id")
		session.In("tag_rel.tag_id", tagIDs)
		session.And("tag_rel.status = ?", entity.TagRelStatusAvailable)
	}
	if len(userID) > 0 {
		session.And("article.user_id = ?", userID)
		if !showHidden {
			session.And("article.show = ?", entity.ArticleShow)
		}
	} else {
		session.And("article.show = ?", entity.ArticleShow)
	}
	if inDays > 0 {
		session.And("article.created_at > ?", time.Now().AddDate(0, 0, -inDays))
	}

	switch orderCond {
	case "newest":
		session.OrderBy("article.pin desc,article.created_at DESC")
	case "active":
		if inDays == 0 {
			session.And("article.created_at > ?", time.Now().AddDate(0, 0, -180))
		}
		session.And("article.post_update_time > ?", time.Now().AddDate(0, 0, -90))
		session.OrderBy("article.pin desc,article.post_update_time DESC, article.updated_at DESC")
	case "hot":
		session.OrderBy("article.pin desc,article.hot_score DESC")
	case "score":
		session.OrderBy("article.pin desc,article.vote_count DESC, article.view_count DESC")
	case "unanswered":
		session.Where("article.answer_count = 0")
		session.OrderBy("article.pin desc,article.created_at DESC")
	}

	total, err = pager.Help(page, pageSize, &articleList, &entity.Article{}, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if handler.GetEnableShortID(ctx) {
		for _, item := range articleList {
			item.ID = uid.EnShortID(item.ID)
		}
	}
	return articleList, total, err
}

// GetRecommendArticlePageByTags get recommend article page by tags
func (qr *articleRepo) GetRecommendArticlePageByTags(ctx context.Context, userID string, tagIDs, followedArticleIDs []string, page, pageSize int) (
	articleList []*entity.Article, total int64, err error) {
	articleList = make([]*entity.Article, 0)
	orderBySQL := "article.pin DESC, article.created_at DESC"

	// Please Make sure every article has at least one tag
	selectSQL := entity.Article{}.TableName() + ".*"
	if len(followedArticleIDs) > 0 {
		idStr := "'" + strings.Join(followedArticleIDs, "','") + "'"
		selectSQL += fmt.Sprintf(", CASE WHEN article.id IN (%s) THEN 0 ELSE 1 END AS order_priority", idStr)
		orderBySQL = "order_priority, " + orderBySQL
	}
	session := qr.data.DB.Context(ctx).Select(selectSQL)

	if len(tagIDs) > 0 {
		session.Where("article.user_id != ?", userID).
			And("article.id NOT IN (SELECT article_id FROM answer WHERE user_id = ?)", userID).
			Join("INNER", "tag_rel", "article.id = tag_rel.object_id").
			And("tag_rel.status = ?", entity.TagRelStatusAvailable).
			Join("INNER", "tag", "tag.id = tag_rel.tag_id").
			In("tag.id", tagIDs)
	} else if len(followedArticleIDs) == 0 {
		return articleList, 0, nil
	}

	if len(followedArticleIDs) > 0 {
		if len(tagIDs) > 0 {
			// if tags provided, show followed articles and tag articles
			session.Or(builder.In("article.id", followedArticleIDs))
		} else {
			// if no tags, only show followed articles
			session.Where(builder.In("article.id", followedArticleIDs))
		}
	}

	session.
		And("article.show = ? and article.status = ?", entity.ArticleShow, entity.ArticleStatusAvailable).
		Distinct("article.id").
		OrderBy(orderBySQL)

	total, err = pager.Help(page, pageSize, &articleList, &entity.Article{}, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	if handler.GetEnableShortID(ctx) {
		for _, item := range articleList {
			item.ID = uid.EnShortID(item.ID)
		}
	}

	return articleList, total, err
}

func (qr *articleRepo) AdminArticlePage(ctx context.Context, search *schema.AdminArticlePageReq) ([]*entity.Article, int64, error) {
	var (
		count   int64
		err     error
		session = qr.data.DB.Context(ctx).Table("article")
	)

	session.Where(builder.Eq{
		"status": search.Status,
	})

	rows := make([]*entity.Article, 0)
	if search.Page > 0 {
		search.Page = search.Page - 1
	} else {
		search.Page = 0
	}
	if search.PageSize == 0 {
		search.PageSize = constant.DefaultPageSize
	}

	// search by article title like or article id
	if len(search.Query) > 0 {
		// check id search
		var (
			idSearch = false
			id       = ""
		)

		if strings.Contains(search.Query, "article:") {
			idSearch = true
			id = strings.TrimSpace(strings.TrimPrefix(search.Query, "article:"))
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
func (qr *articleRepo) UpdateSearch(ctx context.Context, articleID string) (err error) {
	// check search plugin
	var s plugin.Search
	_ = plugin.CallSearch(func(search plugin.Search) error {
		s = search
		return nil
	})
	if s == nil {
		return
	}
	return
	//articleID = uid.DeShortID(articleID)
	//article, exist, err := qr.GetArticle(ctx, articleID)
	//if !exist {
	//	return
	//}
	//if err != nil {
	//	return err
	//}
	//
	//// get tags
	//var (
	//	tagListList = make([]*entity.TagRel, 0)
	//	tags        = make([]string, 0)
	//)
	//session := qr.data.DB.Context(ctx).Where("object_id = ?", articleID)
	//session.Where("status = ?", entity.TagRelStatusAvailable)
	//err = session.Find(&tagListList)
	//if err != nil {
	//	return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	//}
	//for _, tag := range tagListList {
	//	tags = append(tags, tag.TagID)
	//}
	//content := &plugin.SearchContent{
	//	ObjectID: articleID,
	//	Title:    article.Title,
	//	Type:     constant.ArticleObjectType,
	//	Content:  article.OriginalText,
	//	//Answers:     int64(article.AnswerCount),
	//	Status: plugin.SearchContentStatus(article.Status),
	//	Tags:   tags,
	//	//ArticleID:   articleID,
	//	UserID:      article.UserID,
	//	Views:       int64(article.ViewCount),
	//	Created:     article.CreatedAt.Unix(),
	//	Active:      article.UpdatedAt.Unix(),
	//	Score:       int64(article.VoteCount),
	//	HasAccepted: article.AcceptedAnswerID != "" && article.AcceptedAnswerID != "0",
	//}
	//err = s.UpdateContent(ctx, content)
	//return
}

func (qr *articleRepo) RemoveAllUserArticle(ctx context.Context, userID string) (err error) {
	// get all article id that need to be deleted
	articleIDs := make([]string, 0)
	session := qr.data.DB.Context(ctx).Where("user_id = ?", userID)
	session.Where("status != ?", entity.ArticleStatusDeleted)
	err = session.Select("id").Table("article").Find(&articleIDs)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if len(articleIDs) == 0 {
		return nil
	}

	log.Infof("find %d articles need to be deleted for user %s", len(articleIDs), userID)

	// delete all article
	session = qr.data.DB.Context(ctx).Where("user_id = ?", userID)
	session.Where("status != ?", entity.ArticleStatusDeleted)
	_, err = session.Cols("status", "updated_at").Update(&entity.Article{
		UpdatedAt: time.Now(),
		Status:    entity.ArticleStatusDeleted,
	})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	// update search content
	for _, id := range articleIDs {
		_ = qr.UpdateSearch(ctx, id)
	}
	return nil
}
