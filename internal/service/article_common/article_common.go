package articlecommon

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
//
//package articlecommon

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"time"

	"github.com/apache/incubator-answer/internal/base/constant"
	"github.com/apache/incubator-answer/internal/base/data"
	"github.com/apache/incubator-answer/internal/base/handler"
	"github.com/apache/incubator-answer/internal/base/reason"
	"github.com/apache/incubator-answer/internal/service/activity_common"
	"github.com/apache/incubator-answer/internal/service/activity_queue"
	"github.com/apache/incubator-answer/internal/service/config"
	metacommon "github.com/apache/incubator-answer/internal/service/meta_common"
	"github.com/apache/incubator-answer/internal/service/revision"
	"github.com/apache/incubator-answer/pkg/checker"
	"github.com/apache/incubator-answer/pkg/htmltext"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/segmentfault/pacman/errors"

	"github.com/apache/incubator-answer/internal/entity"
	"github.com/apache/incubator-answer/internal/schema"
	answercommon "github.com/apache/incubator-answer/internal/service/answer_common"
	collectioncommon "github.com/apache/incubator-answer/internal/service/collection_common"
	tagcommon "github.com/apache/incubator-answer/internal/service/tag_common"
	usercommon "github.com/apache/incubator-answer/internal/service/user_common"
	"github.com/segmentfault/pacman/log"
)

// ArticleRepo article repository
type ArticleRepo interface {
	AddArticle(ctx context.Context, article *entity.Article) (err error)
	RemoveArticle(ctx context.Context, id string) (err error)
	UpdateArticle(ctx context.Context, article *entity.Article, Cols []string) (err error)
	GetArticle(ctx context.Context, id string) (article *entity.Article, exist bool, err error)
	GetArticleList(ctx context.Context, article *entity.Article) (articles []*entity.Article, err error)
	GetArticlePage(ctx context.Context, page, pageSize int, tagIDs []string, userID, orderCond string, inDays int, showHidden, showPending bool) (
		articleList []*entity.Article, total int64, err error)
	GetRecommendArticlePageByTags(ctx context.Context, userID string, tagIDs, followedArticleIDs []string, page, pageSize int) (articleList []*entity.Article, total int64, err error)
	UpdateArticleStatus(ctx context.Context, articleID string, status int) (err error)
	UpdateArticleStatusWithOutUpdateTime(ctx context.Context, article *entity.Article) (err error)
	RecoverArticle(ctx context.Context, articleID string) (err error)
	UpdateArticleOperation(ctx context.Context, article *entity.Article) (err error)
	GetArticlesByTitle(ctx context.Context, title string, pageSize int) (articleList []*entity.Article, err error)
	UpdatePvCount(ctx context.Context, articleID string) (err error)
	UpdateAnswerCount(ctx context.Context, articleID string, num int) (err error)
	UpdateCollectionCount(ctx context.Context, articleID string) (count int64, err error)
	UpdateAccepted(ctx context.Context, article *entity.Article) (err error)
	UpdateLastAnswer(ctx context.Context, article *entity.Article) (err error)
	FindByID(ctx context.Context, id []string) (articleList []*entity.Article, err error)
	AdminArticlePage(ctx context.Context, search *schema.AdminArticlePageReq) ([]*entity.Article, int64, error)
	GetArticleCount(ctx context.Context) (count int64, err error)
	GetUserArticleCount(ctx context.Context, userID string, show int) (count int64, err error)
	SitemapArticles(ctx context.Context, page, pageSize int) (articleIDList []*schema.SiteMapArticleInfo, err error)
	RemoveAllUserArticle(ctx context.Context, userID string) (err error)
	UpdateSearch(ctx context.Context, articleID string) (err error)
}

// ArticleCommon user service
type ArticleCommon struct {
	articleRepo          ArticleRepo
	answerRepo           answercommon.AnswerRepo
	voteRepo             activity_common.VoteRepo
	followCommon         activity_common.FollowRepo
	tagCommon            *tagcommon.TagCommonService
	userCommon           *usercommon.UserCommon
	collectionCommon     *collectioncommon.CollectionCommon
	AnswerCommon         *answercommon.AnswerCommon
	metaCommonService    *metacommon.MetaCommonService
	configService        *config.ConfigService
	activityQueueService activity_queue.ActivityQueueService
	revisionRepo         revision.RevisionRepo
	data                 *data.Data
}

func NewArticleCommon(articleRepo ArticleRepo,
	answerRepo answercommon.AnswerRepo,
	voteRepo activity_common.VoteRepo,
	followCommon activity_common.FollowRepo,
	tagCommon *tagcommon.TagCommonService,
	userCommon *usercommon.UserCommon,
	collectionCommon *collectioncommon.CollectionCommon,
	answerCommon *answercommon.AnswerCommon,
	metaCommonService *metacommon.MetaCommonService,
	configService *config.ConfigService,
	activityQueueService activity_queue.ActivityQueueService,
	revisionRepo revision.RevisionRepo,
	data *data.Data,
) *ArticleCommon {
	return &ArticleCommon{
		articleRepo:          articleRepo,
		answerRepo:           answerRepo,
		voteRepo:             voteRepo,
		followCommon:         followCommon,
		tagCommon:            tagCommon,
		userCommon:           userCommon,
		collectionCommon:     collectionCommon,
		AnswerCommon:         answerCommon,
		metaCommonService:    metaCommonService,
		configService:        configService,
		activityQueueService: activityQueueService,
		revisionRepo:         revisionRepo,
		data:                 data,
	}
}

func (qs *ArticleCommon) GetUserArticleCount(ctx context.Context, userID string) (count int64, err error) {
	return qs.articleRepo.GetUserArticleCount(ctx, userID, 0)
}

func (qs *ArticleCommon) GetPersonalUserArticleCount(ctx context.Context, loginUserID, userID string, isAdmin bool) (count int64, err error) {
	show := entity.ArticleShow
	if loginUserID == userID || isAdmin {
		show = 0
	}
	return qs.articleRepo.GetUserArticleCount(ctx, userID, show)
}

func (qs *ArticleCommon) UpdatePv(ctx context.Context, articleID string) error {
	return qs.articleRepo.UpdatePvCount(ctx, articleID)
}

//func (qs *ArticleCommon) UpdateAnswerCount(ctx context.Context, articleID string) error {
//	count, err := qs.answerRepo.GetCountByArticleID(ctx, articleID)
//	if err != nil {
//		return err
//	}
//	if count == 0 {
//		err = qs.articleRepo.UpdateLastAnswer(ctx, &entity.Article{
//			ID:           articleID,
//			LastAnswerID: "0",
//		})
//		if err != nil {
//			return err
//		}
//	}
//	return qs.articleRepo.UpdateAnswerCount(ctx, articleID, int(count))
//}

func (qs *ArticleCommon) UpdateCollectionCount(ctx context.Context, articleID string) (count int64, err error) {
	return qs.articleRepo.UpdateCollectionCount(ctx, articleID)
}

func (qs *ArticleCommon) UpdateAccepted(ctx context.Context, articleID, AnswerID string) error {
	article := &entity.Article{}
	article.ID = articleID
	//article.AcceptedAnswerID = AnswerID
	return qs.articleRepo.UpdateAccepted(ctx, article)
}

func (qs *ArticleCommon) UpdateLastAnswer(ctx context.Context, articleID, AnswerID string) error {
	article := &entity.Article{}
	article.ID = articleID
	//article.LastAnswerID = AnswerID
	return qs.articleRepo.UpdateLastAnswer(ctx, article)
}

func (qs *ArticleCommon) UpdatePostTime(ctx context.Context, articleID string) error {
	articleinfo := &entity.Article{}
	now := time.Now()
	_ = now
	articleinfo.ID = articleID
	//articleinfo.PostUpdateTime = now
	return qs.articleRepo.UpdateArticle(ctx, articleinfo, []string{"post_update_time"})
}
func (qs *ArticleCommon) UpdatePostSetTime(ctx context.Context, articleID string, setTime time.Time) error {
	articleinfo := &entity.Article{}
	articleinfo.ID = articleID
	articleinfo.PostUpdateTime = setTime
	return qs.articleRepo.UpdateArticle(ctx, articleinfo, []string{"post_update_time"})
}

func (qs *ArticleCommon) FindInfoByID(ctx context.Context, articleIDs []string, loginUserID string) (map[string]*schema.ArticleInfoResp, error) {
	list := make(map[string]*schema.ArticleInfoResp)
	articleList, err := qs.articleRepo.FindByID(ctx, articleIDs)
	if err != nil {
		return list, err
	}
	articles, err := qs.FormatArticles(ctx, articleList, loginUserID)
	if err != nil {
		return list, err
	}
	for _, item := range articles {
		list[item.ID] = item
	}
	return list, nil
}

func (qs *ArticleCommon) InviteUserInfo(ctx context.Context, articleID string) (inviteList []*schema.UserBasicInfo, err error) {
	return
	//InviteUserInfo := make([]*schema.UserBasicInfo, 0)
	//dbinfo, has, err := qs.articleRepo.GetArticle(ctx, articleID)
	//if err != nil {
	//	return InviteUserInfo, err
	//}
	//if !has {
	//	return InviteUserInfo, errors.NotFound(reason.ArticleNotFound)
	//}
	///@ms: InviteUser
	//if dbinfo.InviteUserID != "" {
	//	InviteUserIDs := make([]string, 0)
	//	err := json.Unmarshal([]byte(dbinfo.InviteUserID), &InviteUserIDs)
	//	if err == nil {
	//		inviteUserInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, InviteUserIDs)
	//		if err == nil {
	//			for _, userid := range InviteUserIDs {
	//				_, ok := inviteUserInfoMap[userid]
	//				if ok {
	//					InviteUserInfo = append(InviteUserInfo, inviteUserInfoMap[userid])
	//				}
	//			}
	//		}
	//	}
	//}
	//return InviteUserInfo, nil
}

func (qs *ArticleCommon) Info(ctx context.Context, articleID string, loginUserID string) (resp *schema.ArticleInfoResp, err error) {
	articleInfo, has, err := qs.articleRepo.GetArticle(ctx, articleID)
	if err != nil {
		return resp, err
	}
	articleInfo.ID = uid.DeShortID(articleInfo.ID)
	if !has {
		return resp, errors.NotFound(reason.ArticleNotFound)
	}
	resp = qs.ShowFormat(ctx, articleInfo)
	if resp.Status == entity.ArticleStatusClosed {
		metaInfo, err := qs.metaCommonService.GetMetaByObjectIdAndKey(ctx, articleInfo.ID, entity.ArticleCloseReasonKey)
		if err != nil {
			log.Error(err)
		} else {
			closeMsg := &schema.CloseArticleMeta{}
			err = json.Unmarshal([]byte(metaInfo.Value), closeMsg)
			if err != nil {
				log.Error("json.Unmarshal CloseArticleMeta error", err.Error())
			} else {
				cfg, err := qs.configService.GetConfigByID(ctx, closeMsg.CloseType)
				if err != nil {
					log.Error("json.Unmarshal ArticleCloseJson error", err.Error())
				} else {
					reasonItem := &schema.ReasonItem{}
					_ = json.Unmarshal(cfg.GetByteValue(), reasonItem)
					reasonItem.Translate(cfg.Key, handler.GetLangByCtx(ctx))
					operation := &schema.Operation{}
					operation.Type = reasonItem.Name
					operation.Description = reasonItem.Description
					operation.Msg = closeMsg.CloseMsg
					operation.Time = metaInfo.CreatedAt.Unix()
					operation.Level = schema.OperationLevelInfo
					resp.Operation = operation
				}
			}
		}
	}

	if resp.Status != entity.ArticleStatusDeleted {
		if resp.Tags, err = qs.tagCommon.GetObjectTag(ctx, articleID); err != nil {
			return resp, err
		}
	} else {
		revisionInfo, exist, err := qs.revisionRepo.GetLastRevisionByObjectID(ctx, articleID)
		if err != nil {
			log.Errorf("get revision error %s", err)
		}
		if exist {
			articleWithTagsRevision := &entity.ArticleWithTagsRevision{}
			if err = json.Unmarshal([]byte(revisionInfo.Content), articleWithTagsRevision); err != nil {
				log.Errorf("revision parsing error %s", err)
				return resp, nil
			}
			for _, tag := range articleWithTagsRevision.Tags {
				resp.Tags = append(resp.Tags, &schema.TagResp{
					ID:              tag.ID,
					SlugName:        tag.SlugName,
					DisplayName:     tag.DisplayName,
					MainTagSlugName: tag.MainTagSlugName,
					Recommend:       tag.Recommend,
					Reserved:        tag.Reserved,
				})

			}
		}
	}
	for _, v := range resp.Tags {
		log.Infof("my_tags id:%+v", *v)
	}

	userIds := make([]string, 0)
	if checker.IsNotZeroString(articleInfo.UserID) {
		userIds = append(userIds, articleInfo.UserID)
	}
	//if checker.IsNotZeroString(articleInfo.LastEditUserID) {
	//	userIds = append(userIds, articleInfo.LastEditUserID)
	//}
	if checker.IsNotZeroString(resp.LastAnsweredUserID) {
		userIds = append(userIds, resp.LastAnsweredUserID)
	}
	userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIds)
	if err != nil {
		return resp, err
	}
	resp.UserInfo = userInfoMap[articleInfo.UserID]
	//resp.UpdateUserInfo = userInfoMap[articleInfo.LastEditUserID]
	resp.LastAnsweredUserInfo = userInfoMap[resp.LastAnsweredUserID]
	if len(loginUserID) == 0 {
		return resp, nil
	}

	resp.VoteStatus = qs.voteRepo.GetVoteStatus(ctx, articleID, loginUserID)
	resp.IsFollowed, _ = qs.followCommon.IsFollowed(ctx, loginUserID, articleID)

	ids, err := qs.AnswerCommon.SearchAnswerIDs(ctx, loginUserID, articleInfo.ID)
	if err != nil {
		log.Error("AnswerFunc.SearchAnswerIDs", err)
	}
	resp.Answered = len(ids) > 0
	if resp.Answered {
		resp.FirstAnswerId = ids[0]
	}

	collectedMap, err := qs.collectionCommon.SearchObjectCollected(ctx, loginUserID, []string{articleInfo.ID})
	if err != nil {
		return nil, err
	}
	if len(collectedMap) > 0 {
		resp.Collected = true
	}
	return resp, nil
}

// \b 是一个特殊的元字符，表示单词边界。它匹配一个单词的开始或结束的位置
func RepImages(htmls string, needSize int) []string {
	var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	imgs := imgRE.FindAllStringSubmatch(htmls, needSize) //imgRE.FindAllStringSubmatch(htmls, -1)
	out := make([]string, len(imgs))
	for i := range out {
		out[i] = imgs[i][1]
		//fmt.Println(strconv.Itoa(i), out[i])
	}
	return out
}
func RegFindImages(htmls string, needSize int) []string {
	var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	imgs := imgRE.FindStringSubmatch(htmls) //imgRE.FindAllStringSubmatch(htmls, -1)
	//fmt.Printf("imgs:%v \n", imgs)
	log.Infof("RegFindImages   imgs:%v \n", imgs)

	out := make([]string, 0)

	if len(imgs) > 1 {
		out = append(out, imgs[1])
	}
	return out
}
func GetArticleThumbnails(articleInfo *entity.Article) []schema.ArticleThumbnail {
	/*
		1. 首先是否设置了文章自己本身的特色图片。
		2. 如果没有，就检查下这篇文章的标签是否有特色图片。
		3. 如果没有，就检查下这篇文章是否有图片，有就获取第一张。
		4. 如果还是没有，就检查下这篇文章所在的分类是否有特色图片。
	*/
	thumbnails := make([]schema.ArticleThumbnail, 0)
	if articleInfo.Thumbnails != "" {
		err := json.Unmarshal([]byte(articleInfo.Thumbnails), &thumbnails)
		if err != nil {
			log.Errorf("GetArticleThumbnails err:%v \n", err)
		}
	}
	if len(thumbnails) != 0 {
		log.Infof("thumbnails1:%+v", thumbnails)
		return thumbnails
	}

	images := RegFindImages(articleInfo.ParsedText, 1)

	for _, val := range images {
		thumbnail := schema.ArticleThumbnail{
			Url: val,
		}
		thumbnails = append(thumbnails, thumbnail)
	}

	if len(thumbnails) != 0 {
		log.Infof("thumbnails2:%+v", thumbnails)
		return thumbnails
	}
	randomNum := rand.Intn(41) //[0,n)之
	imageUrl := fmt.Sprintf("%s%d%s", constant.OSS_BUCKET_NAME+"/images/postthumbnail/", randomNum, ".jpg")
	thumbnail := schema.ArticleThumbnail{
		Url: imageUrl,
	}

	thumbnails = append(thumbnails, thumbnail)
	log.Infof("thumbnails3:%+v", thumbnails)
	return thumbnails
}
func (qs *ArticleCommon) FormatArticlesPage(
	ctx context.Context, articleList []*entity.Article, loginUserID string, orderCond string) (
	formattedArticles []*schema.ArticlePageResp, err error) {
	formattedArticles = make([]*schema.ArticlePageResp, 0)
	articleIDs := make([]string, 0)
	userIDs := make([]string, 0)

	for _, articleInfo := range articleList {

		thumbnails := GetArticleThumbnails(articleInfo)

		t := &schema.ArticlePageResp{
			ID:              articleInfo.ID,
			CreatedAt:       articleInfo.CreatedAt.Unix(),
			Title:           articleInfo.Title,
			UrlTitle:        htmltext.UrlTitle(articleInfo.Title),
			Description:     htmltext.FetchExcerpt(articleInfo.ParsedText, "...", 80), //240),
			Status:          articleInfo.Status,
			ViewCount:       articleInfo.ViewCount,
			UniqueViewCount: articleInfo.UniqueViewCount,
			VoteCount:       articleInfo.VoteCount,
			//AnswerCount:      articleInfo.AnswerCount,
			CollectionCount: articleInfo.CollectionCount,
			FollowCount:     articleInfo.FollowCount,
			//AcceptedAnswerID: articleInfo.AcceptedAnswerID,
			//LastAnswerID: articleInfo.LastAnswerID,
			Pin:  articleInfo.Pin,
			Show: articleInfo.Show,

			Thumbnails: thumbnails,
		}

		articleIDs = append(articleIDs, articleInfo.ID)
		userIDs = append(userIDs, articleInfo.UserID)
		haveEdited, haveAnswered := false, false
		//if checker.IsNotZeroString(articleInfo.LastEditUserID) {
		//	haveEdited = true
		//	userIDs = append(userIDs, articleInfo.LastEditUserID)
		//}
		//if checker.IsNotZeroString(articleInfo.LastAnswerID) {
		//	haveAnswered = true
		//
		//	answerInfo, exist, err := qs.answerRepo.GetAnswer(ctx, articleInfo.LastAnswerID)
		//	if err == nil && exist {
		//		if answerInfo.LastEditUserID != "0" {
		//			t.LastAnsweredUserID = answerInfo.LastEditUserID
		//		} else {
		//			t.LastAnsweredUserID = answerInfo.UserID
		//		}
		//		t.LastAnsweredAt = answerInfo.CreatedAt
		//		userIDs = append(userIDs, t.LastAnsweredUserID)
		//	}
		//}

		// if order condition is newest or nobody edited or nobody answered, only show article author
		if orderCond == schema.ArticleOrderCondNewest || (!haveEdited && !haveAnswered) {
			t.OperationType = schema.ArticlePageRespOperationTypeAsked
			t.OperatedAt = articleInfo.CreatedAt.Unix()
			t.Operator = &schema.ArticlePageRespOperator{ID: articleInfo.UserID}
		} else {
			// if no one
			if haveEdited {
				t.OperationType = schema.ArticlePageRespOperationTypeModified
				t.OperatedAt = articleInfo.UpdatedAt.Unix()
				//t.Operator = &schema.ArticlePageRespOperator{ID: articleInfo.LastEditUserID}
			}

			if haveAnswered {
				if t.LastAnsweredAt.Unix() > t.OperatedAt {
					t.OperationType = schema.ArticlePageRespOperationTypeAnswered
					t.OperatedAt = t.LastAnsweredAt.Unix()
					t.Operator = &schema.ArticlePageRespOperator{ID: t.LastAnsweredUserID}
				}
			}
		}
		formattedArticles = append(formattedArticles, t)
	}

	tagsMap, err := qs.tagCommon.BatchGetObjectTag(ctx, articleIDs)
	if err != nil {
		return formattedArticles, err
	}
	userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIDs)
	if err != nil {
		return formattedArticles, err
	}

	for _, item := range formattedArticles {
		tags, ok := tagsMap[item.ID]
		if ok {
			item.Tags = tags
		} else {
			item.Tags = make([]*schema.TagResp, 0)
		}
		userInfo, ok := userInfoMap[item.Operator.ID]
		if ok {
			if userInfo != nil {
				item.Operator.DisplayName = userInfo.DisplayName
				item.Operator.Username = userInfo.Username
				item.Operator.Rank = userInfo.Rank
				item.Operator.Status = userInfo.Status
			}
		}

	}
	return formattedArticles, nil
}

func (qs *ArticleCommon) FormatArticles(ctx context.Context, articleList []*entity.Article, loginUserID string) ([]*schema.ArticleInfoResp, error) {
	list := make([]*schema.ArticleInfoResp, 0)
	objectIds := make([]string, 0)
	userIds := make([]string, 0)

	for _, articleInfo := range articleList {
		item := qs.ShowFormat(ctx, articleInfo)
		list = append(list, item)
		objectIds = append(objectIds, item.ID)
		userIds = append(userIds, item.UserID, item.LastEditUserID, item.LastAnsweredUserID)
	}
	tagsMap, err := qs.tagCommon.BatchGetObjectTag(ctx, objectIds)
	if err != nil {
		return list, err
	}

	userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIds)
	if err != nil {
		return list, err
	}

	for _, item := range list {
		item.Tags = tagsMap[item.ID]
		item.UserInfo = userInfoMap[item.UserID]
		item.UpdateUserInfo = userInfoMap[item.LastEditUserID]
		item.LastAnsweredUserInfo = userInfoMap[item.LastAnsweredUserID]
	}
	if loginUserID == "" {
		return list, nil
	}

	collectedMap, err := qs.collectionCommon.SearchObjectCollected(ctx, loginUserID, objectIds)
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		item.Collected = collectedMap[item.ID]
	}
	return list, nil
}

// RemoveArticle delete article
func (qs *ArticleCommon) RemoveArticle(ctx context.Context, req *schema.RemoveArticleReq) (err error) {
	articleInfo, has, err := qs.articleRepo.GetArticle(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	if articleInfo.Status == entity.ArticleStatusDeleted {
		return nil
	}

	articleInfo.Status = entity.ArticleStatusDeleted
	err = qs.articleRepo.UpdateArticleStatus(ctx, articleInfo.ID, articleInfo.Status)
	if err != nil {
		return err
	}

	userArticleCount, err := qs.GetUserArticleCount(ctx, articleInfo.UserID)
	if err != nil {
		log.Error("user GetUserArticleCount error", err.Error())
	} else {
		//@ms:TODO	err = qs.userCommon.UpdateArticleCount(ctx, articleInfo.UserID, userArticleCount)
		_ = userArticleCount
		if err != nil {
			log.Error("user IncreaseArticleCount error", err.Error())
		}
	}

	return nil
}

func (qs *ArticleCommon) CloseArticle(ctx context.Context, req *schema.CloseArticleReq) error {
	articleInfo, has, err := qs.articleRepo.GetArticle(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	articleInfo.Status = entity.ArticleStatusClosed
	err = qs.articleRepo.UpdateArticleStatus(ctx, articleInfo.ID, articleInfo.Status)
	if err != nil {
		return err
	}

	closeMeta, _ := json.Marshal(schema.CloseArticleMeta{
		CloseType: req.CloseType,
		CloseMsg:  req.CloseMsg,
	})
	err = qs.metaCommonService.AddMeta(ctx, req.ID, entity.ArticleCloseReasonKey, string(closeMeta))
	if err != nil {
		return err
	}

	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           articleInfo.UserID,
		ObjectID:         articleInfo.ID,
		OriginalObjectID: articleInfo.ID,
		ActivityTypeKey:  constant.ActArticleClosed,
	})
	return nil
}

//// RemoveAnswer delete answer
//func (qs *ArticleCommon) RemoveAnswer(ctx context.Context, id string) (err error) {
//	answerinfo, has, err := qs.answerRepo.GetByID(ctx, id)
//	if err != nil {
//		return err
//	}
//	if !has {
//		return nil
//	}
//
//	// user add article count
//
//	err = qs.UpdateAnswerCount(ctx, answerinfo.ArticleID)
//	if err != nil {
//		log.Error("UpdateAnswerCount error", err.Error())
//	}
//	userAnswerCount, err := qs.answerRepo.GetCountByUserID(ctx, answerinfo.UserID)
//	if err != nil {
//		log.Error("GetCountByUserID error", err.Error())
//	}
//	err = qs.userCommon.UpdateAnswerCount(ctx, answerinfo.UserID, int(userAnswerCount))
//	if err != nil {
//		log.Error("user UpdateAnswerCount error", err.Error())
//	}
//
//	return qs.answerRepo.RemoveAnswer(ctx, id)
//}

func (qs *ArticleCommon) SitemapCron(ctx context.Context) {
	articleNum, err := qs.articleRepo.GetArticleCount(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	if articleNum <= constant.SitemapMaxSize {
		_, err = qs.articleRepo.SitemapArticles(ctx, 1, int(articleNum))
		if err != nil {
			log.Errorf("get site map article error: %v", err)
		}
		return
	}

	totalPages := int(math.Ceil(float64(articleNum) / float64(constant.SitemapMaxSize)))
	for i := 1; i <= totalPages; i++ {
		_, err = qs.articleRepo.SitemapArticles(ctx, i, constant.SitemapMaxSize)
		if err != nil {
			log.Errorf("get site map article error: %v", err)
			return
		}
	}
}

func (qs *ArticleCommon) SetCache(ctx context.Context, cachekey string, info interface{}) error {
	infoStr, err := json.Marshal(info)
	if err != nil {
		return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	err = qs.data.Cache.SetString(ctx, cachekey, string(infoStr), schema.DashboardCacheTime)
	if err != nil {
		return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	return nil
}

func (qs *ArticleCommon) ShowListFormat(ctx context.Context, data *entity.Article) *schema.ArticleInfoResp {
	return qs.ShowFormat(ctx, data)
}

func (qs *ArticleCommon) ShowFormat(ctx context.Context, data *entity.Article) *schema.ArticleInfoResp {
	info := schema.ArticleInfoResp{}
	info.ID = data.ID
	if handler.GetEnableShortID(ctx) {
		info.ID = uid.EnShortID(data.ID)
	}
	info.Title = data.Title
	info.UrlTitle = htmltext.UrlTitle(data.Title)
	info.Content = data.OriginalText
	info.HTML = data.ParsedText
	info.ViewCount = data.ViewCount
	info.UniqueViewCount = data.UniqueViewCount
	info.VoteCount = data.VoteCount
	//info.AnswerCount = data.AnswerCount
	info.CollectionCount = data.CollectionCount
	info.FollowCount = data.FollowCount
	//info.AcceptedAnswerID = data.AcceptedAnswerID
	//info.LastAnswerID = data.LastAnswerID
	info.CreateTime = data.CreatedAt.Unix()
	info.UpdateTime = data.UpdatedAt.Unix()
	info.PostUpdateTime = data.PostUpdateTime.Unix()
	if data.PostUpdateTime.Unix() < 1 {
		info.PostUpdateTime = 0
	}
	info.ArticleUpdateTime = data.UpdatedAt.Unix()
	if data.UpdatedAt.Unix() < 1 {
		info.ArticleUpdateTime = 0
	}
	info.Status = data.Status
	info.Pin = data.Pin
	info.Show = data.Show
	info.UserID = data.UserID
	//info.LastEditUserID = data.LastEditUserID
	//if data.LastAnswerID != "0" {
	//	answerInfo, exist, err := qs.answerRepo.GetAnswer(ctx, data.LastAnswerID)
	//	if err == nil && exist {
	//		if answerInfo.LastEditUserID != "0" {
	//			info.LastAnsweredUserID = answerInfo.LastEditUserID
	//		} else {
	//			info.LastAnsweredUserID = answerInfo.UserID
	//		}
	//	}
	//
	//}
	info.ContentFormat = data.OriginalTextFormat

	info.Tags = make([]*schema.TagResp, 0)
	return &info
}
func (qs *ArticleCommon) ShowFormatWithTag(ctx context.Context, data *entity.ArticleWithTagsRevision) *schema.ArticleInfoResp {
	info := qs.ShowFormat(ctx, &data.Article)
	Tags := make([]*schema.TagResp, 0)
	for _, tag := range data.Tags {
		item := &schema.TagResp{}
		item.SlugName = tag.SlugName
		item.DisplayName = tag.DisplayName
		item.Recommend = tag.Recommend
		item.Reserved = tag.Reserved
		Tags = append(Tags, item)
	}
	info.Tags = Tags
	return info
}
