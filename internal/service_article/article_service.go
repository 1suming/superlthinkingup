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

package service_article

import (
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-answer/internal/service/event_queue"
	"strings"
	"time"

	"github.com/apache/incubator-answer/internal/base/constant"
	"github.com/apache/incubator-answer/internal/base/handler"
	"github.com/apache/incubator-answer/internal/base/pager"
	"github.com/apache/incubator-answer/internal/base/reason"
	"github.com/apache/incubator-answer/internal/base/translator"
	"github.com/apache/incubator-answer/internal/base/validator"
	"github.com/apache/incubator-answer/internal/entity"
	"github.com/apache/incubator-answer/internal/schema"
	"github.com/apache/incubator-answer/internal/service/activity"
	"github.com/apache/incubator-answer/internal/service/activity_common"
	"github.com/apache/incubator-answer/internal/service/activity_queue"
	answercommon "github.com/apache/incubator-answer/internal/service/answer_common"
	articlecommon "github.com/apache/incubator-answer/internal/service/article_common"
	collectioncommon "github.com/apache/incubator-answer/internal/service/collection_common"
	"github.com/apache/incubator-answer/internal/service/config"
	"github.com/apache/incubator-answer/internal/service/export"
	metacommon "github.com/apache/incubator-answer/internal/service/meta_common"
	"github.com/apache/incubator-answer/internal/service/notice_queue"
	"github.com/apache/incubator-answer/internal/service/notification"
	"github.com/apache/incubator-answer/internal/service/permission"
	"github.com/apache/incubator-answer/internal/service/review"
	"github.com/apache/incubator-answer/internal/service/revision_common"
	"github.com/apache/incubator-answer/internal/service/role"
	"github.com/apache/incubator-answer/internal/service/siteinfo_common"
	"github.com/apache/incubator-answer/internal/service/tag"
	tagcommon "github.com/apache/incubator-answer/internal/service/tag_common"
	usercommon "github.com/apache/incubator-answer/internal/service/user_common"
	"github.com/apache/incubator-answer/pkg/checker"
	"github.com/apache/incubator-answer/pkg/converter"
	"github.com/apache/incubator-answer/pkg/htmltext"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
	"golang.org/x/net/context"
)

// ArticleRepo article repository

// ArticleService user service
type ArticleService struct {
	activityRepo                     activity_common.ActivityRepo
	articleRepo                      articlecommon.ArticleRepo
	answerRepo                       answercommon.AnswerRepo
	tagCommon                        *tagcommon.TagCommonService
	tagService                       *tag.TagService
	articlecommon                    *articlecommon.ArticleCommon
	userCommon                       *usercommon.UserCommon
	userRepo                         usercommon.UserRepo
	userRoleRelService               *role.UserRoleRelService
	revisionService                  *revision_common.RevisionService
	metaService                      *metacommon.MetaCommonService
	collectionCommon                 *collectioncommon.CollectionCommon
	answerActivityService            *activity.AnswerActivityService
	emailService                     *export.EmailService
	notificationQueueService         notice_queue.NotificationQueueService
	externalNotificationQueueService notice_queue.ExternalNotificationQueueService
	activityQueueService             activity_queue.ActivityQueueService
	siteInfoService                  siteinfo_common.SiteInfoCommonService
	newArticleNotificationService    *notification.ExternalNotificationService
	reviewService                    *review.ReviewService
	configService                    *config.ConfigService
	eventQueueService                event_queue.EventQueueService
}

func NewArticleService(
	activityRepo activity_common.ActivityRepo,
	articleRepo articlecommon.ArticleRepo,
	answerRepo answercommon.AnswerRepo,
	tagCommon *tagcommon.TagCommonService,
	tagService *tag.TagService,
	articlecommon *articlecommon.ArticleCommon,
	userCommon *usercommon.UserCommon,
	userRepo usercommon.UserRepo,
	userRoleRelService *role.UserRoleRelService,
	revisionService *revision_common.RevisionService,
	metaService *metacommon.MetaCommonService,
	collectionCommon *collectioncommon.CollectionCommon,
	answerActivityService *activity.AnswerActivityService,
	emailService *export.EmailService,
	notificationQueueService notice_queue.NotificationQueueService,
	externalNotificationQueueService notice_queue.ExternalNotificationQueueService,
	activityQueueService activity_queue.ActivityQueueService,
	siteInfoService siteinfo_common.SiteInfoCommonService,
	newArticleNotificationService *notification.ExternalNotificationService,
	reviewService *review.ReviewService,
	configService *config.ConfigService,
	eventQueueService event_queue.EventQueueService,
) *ArticleService {
	return &ArticleService{
		activityRepo:                     activityRepo,
		articleRepo:                      articleRepo,
		answerRepo:                       answerRepo,
		tagCommon:                        tagCommon,
		tagService:                       tagService,
		articlecommon:                    articlecommon,
		userCommon:                       userCommon,
		userRepo:                         userRepo,
		userRoleRelService:               userRoleRelService,
		revisionService:                  revisionService,
		metaService:                      metaService,
		collectionCommon:                 collectionCommon,
		answerActivityService:            answerActivityService,
		emailService:                     emailService,
		notificationQueueService:         notificationQueueService,
		externalNotificationQueueService: externalNotificationQueueService,
		activityQueueService:             activityQueueService,
		siteInfoService:                  siteInfoService,
		newArticleNotificationService:    newArticleNotificationService,
		reviewService:                    reviewService,
		configService:                    configService,
		eventQueueService:                eventQueueService,
	}
}

func (qs *ArticleService) CloseArticle(ctx context.Context, req *schema.CloseArticleReq) error {
	articleInfo, has, err := qs.articleRepo.GetArticle(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	cf, err := qs.configService.GetConfigByID(ctx, req.CloseType)
	if err != nil || cf == nil {
		return errors.BadRequest(reason.ReportNotFound)
	}
	if cf.Key == constant.ReasonADuplicate && !checker.IsURL(req.CloseMsg) {
		return errors.BadRequest(reason.InvalidURLError)
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
	err = qs.metaService.AddMeta(ctx, req.ID, entity.ArticleCloseReasonKey, string(closeMeta))
	if err != nil {
		return err
	}

	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           req.UserID,
		ObjectID:         articleInfo.ID,
		OriginalObjectID: articleInfo.ID,
		ActivityTypeKey:  constant.ActArticleClosed,
	})
	return nil
}

// ReopenArticle reopen article
func (qs *ArticleService) ReopenArticle(ctx context.Context, req *schema.ReopenArticleReq) error {
	articleInfo, has, err := qs.articleRepo.GetArticle(ctx, req.ArticleID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	articleInfo.Status = entity.ArticleStatusAvailable
	err = qs.articleRepo.UpdateArticleStatus(ctx, articleInfo.ID, articleInfo.Status)
	if err != nil {
		return err
	}
	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           req.UserID,
		ObjectID:         articleInfo.ID,
		OriginalObjectID: articleInfo.ID,
		ActivityTypeKey:  constant.ActArticleReopened,
	})
	return nil
}

func (qs *ArticleService) AddArticleCheckTags(ctx context.Context, Tags []*entity.Tag) ([]string, error) {
	list := make([]string, 0)
	for _, tag := range Tags {
		if tag.Reserved {
			list = append(list, tag.DisplayName)
		}
	}
	if len(list) > 0 {
		return list, errors.BadRequest(reason.RequestFormatError)
	}
	return []string{}, nil
}
func (qs *ArticleService) CheckAddArticle(ctx context.Context, req *schema.ArticleAdd) (errorlist any, err error) {
	if len(req.Tags) == 0 {
		errorlist := make([]*validator.FormErrorField, 0)
		errorlist = append(errorlist, &validator.FormErrorField{
			ErrorField: "tags",
			ErrorMsg:   translator.Tr(handler.GetLangByCtx(ctx), reason.TagNotFound),
		})
		err = errors.BadRequest(reason.RecommendTagEnter)
		return errorlist, err
	}
	recommendExist, err := qs.tagCommon.ExistRecommend(ctx, req.Tags)
	if err != nil {
		return
	}
	if !recommendExist {
		errorlist := make([]*validator.FormErrorField, 0)
		errorlist = append(errorlist, &validator.FormErrorField{
			ErrorField: "tags",
			ErrorMsg:   translator.Tr(handler.GetLangByCtx(ctx), reason.RecommendTagEnter),
		})
		err = errors.BadRequest(reason.RecommendTagEnter)
		return errorlist, err
	}

	tagNameList := make([]string, 0)
	for _, tag := range req.Tags {
		tagNameList = append(tagNameList, tag.SlugName)
	}
	Tags, tagerr := qs.tagCommon.GetTagListByNames(ctx, tagNameList)
	if tagerr != nil {
		return errorlist, tagerr
	}
	if !req.ArticlePermission.CanUseReservedTag {
		taglist, err := qs.AddArticleCheckTags(ctx, Tags)
		errMsg := fmt.Sprintf(`"%s" can only be used by moderators.`,
			strings.Join(taglist, ","))
		if err != nil {
			errorlist := make([]*validator.FormErrorField, 0)
			errorlist = append(errorlist, &validator.FormErrorField{
				ErrorField: "tags",
				ErrorMsg:   errMsg,
			})
			err = errors.BadRequest(reason.RecommendTagEnter)
			return errorlist, err
		}
	}
	return nil, nil
}

// HasNewTag
func (qs *ArticleService) HasNewTag(ctx context.Context, tags []*schema.TagItem) (bool, error) {
	return qs.tagCommon.HasNewTag(ctx, tags)
}

// AddArticle add article
func (qs *ArticleService) AddArticle(ctx context.Context, req *schema.ArticleAdd) (articleInfo any, err error) {
	if len(req.Tags) == 0 {
		errorlist := make([]*validator.FormErrorField, 0)
		errorlist = append(errorlist, &validator.FormErrorField{
			ErrorField: "tags",
			ErrorMsg:   translator.Tr(handler.GetLangByCtx(ctx), reason.TagNotFound),
		})
		err = errors.BadRequest(reason.RecommendTagEnter)
		return errorlist, err
	}
	recommendExist, err := qs.tagCommon.ExistRecommend(ctx, req.Tags)
	if err != nil {
		return
	}
	if !recommendExist {
		errorlist := make([]*validator.FormErrorField, 0)
		errorlist = append(errorlist, &validator.FormErrorField{
			ErrorField: "tags",
			ErrorMsg:   translator.Tr(handler.GetLangByCtx(ctx), reason.RecommendTagEnter),
		})
		err = errors.BadRequest(reason.RecommendTagEnter)
		return errorlist, err
	}

	tagNameList := make([]string, 0)
	for _, tag := range req.Tags {
		tag.SlugName = strings.ReplaceAll(tag.SlugName, " ", "-")
		tagNameList = append(tagNameList, tag.SlugName)
	}
	tags, tagerr := qs.tagCommon.GetTagListByNames(ctx, tagNameList)
	if tagerr != nil {
		return articleInfo, tagerr
	}
	if !req.ArticlePermission.CanUseReservedTag {
		taglist, err := qs.AddArticleCheckTags(ctx, tags)
		errMsg := fmt.Sprintf(`"%s" can only be used by moderators.`,
			strings.Join(taglist, ","))
		if err != nil {
			errorlist := make([]*validator.FormErrorField, 0)
			errorlist = append(errorlist, &validator.FormErrorField{
				ErrorField: "tags",
				ErrorMsg:   errMsg,
			})
			err = errors.BadRequest(reason.RecommendTagEnter)
			return errorlist, err
		}
	}

	article := &entity.Article{}
	now := time.Now()
	article.UserID = req.UserID
	article.Title = req.Title
	article.OriginalText = req.Content
	log.Infof("addArticle content:%s", req.Content)
	article.ParsedText = req.HTML
	//article.AcceptedAnswerID = "0"
	//article.LastAnswerID = "0"
	//article.LastEditUserID = "0"
	//article.PostUpdateTime = nil
	article.Status = entity.ArticleStatusPending
	article.RevisionID = "0"
	article.CreatedAt = now
	article.PostUpdateTime = now
	article.Pin = entity.ArticleUnPin
	article.Show = entity.ArticleShow
	article.UpdatedAt = now

	//@cws
	article.PostDate = now
	article.OriginalTextFormat = req.ContentFormat

	err = qs.articleRepo.AddArticle(ctx, article)
	if err != nil {
		return
	}
	article.Status = qs.reviewService.AddArticleReview(ctx, article, req.Tags, req.IP, req.UserAgent)
	if err := qs.articleRepo.UpdateArticleStatus(ctx, article.ID, article.Status); err != nil {
		return nil, err
	}
	objectTagData := schema.TagChange{}
	objectTagData.ObjectID = article.ID
	objectTagData.Tags = req.Tags
	objectTagData.UserID = req.UserID
	err = qs.ChangeTag(ctx, &objectTagData)
	if err != nil {
		return
	}
	_ = qs.articleRepo.UpdateSearch(ctx, article.ID)

	revisionDTO := &schema.AddRevisionDTO{
		UserID:   article.UserID,
		ObjectID: article.ID,
		Title:    article.Title,
	}

	articleWithTagsRevision, err := qs.changeArticleToRevision(ctx, article, tags)
	if err != nil {
		return nil, err
	}
	infoJSON, _ := json.Marshal(articleWithTagsRevision)
	revisionDTO.Content = string(infoJSON)
	revisionID, err := qs.revisionService.AddRevision(ctx, revisionDTO, true)
	if err != nil {
		return
	}

	// user add article count
	userArticleCount, err := qs.articlecommon.GetUserArticleCount(ctx, article.UserID)
	if err != nil {
		log.Errorf("get user article count error %v", err)
	} else {
		err = qs.userCommon.UpdateArticleCount(ctx, article.UserID, userArticleCount)
		if err != nil {
			log.Errorf("update user article count error %v", err)
		}
	}

	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           article.UserID,
		ObjectID:         article.ID,
		OriginalObjectID: article.ID,
		ActivityTypeKey:  constant.ActArticleAsked,
		RevisionID:       revisionID,
	})

	if article.Status == entity.ArticleStatusAvailable {
		qs.externalNotificationQueueService.Send(ctx,
			schema.CreateNewArticleNotificationMsg(article.ID, article.Title, article.UserID, tags))
	}
	qs.eventQueueService.Send(ctx, schema.NewEvent(constant.EventArticleCreate, req.UserID).TID(article.ID).
		QID(article.ID, article.UserID))

	articleInfo, err = qs.GetArticle(ctx, article.ID, article.UserID, req.ArticlePermission)
	return
}

// OperationArticle
func (qs *ArticleService) OperationArticle(ctx context.Context, req *schema.OperationArticleReq) (err error) {
	articleInfo, has, err := qs.articleRepo.GetArticle(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	// Hidden article cannot be placed at the top
	if articleInfo.Show == entity.ArticleHide && req.Operation == schema.ArticleOperationPin {
		return nil
	}
	// Article cannot be hidden when they are at the top
	if articleInfo.Pin == entity.ArticlePin && req.Operation == schema.ArticleOperationHide {
		return nil
	}

	switch req.Operation {
	case schema.ArticleOperationHide:
		articleInfo.Show = entity.ArticleHide
		err = qs.tagCommon.HideTagRelListByObjectID(ctx, req.ID)
		if err != nil {
			return err
		}
		err = qs.tagCommon.RefreshTagCountByArticleID(ctx, req.ID)
		if err != nil {
			return err
		}
	case schema.ArticleOperationShow:
		articleInfo.Show = entity.ArticleShow
		err = qs.tagCommon.ShowTagRelListByObjectID(ctx, req.ID)
		if err != nil {
			return err
		}
		err = qs.tagCommon.RefreshTagCountByArticleID(ctx, req.ID)
		if err != nil {
			return err
		}
	case schema.ArticleOperationPin:
		articleInfo.Pin = entity.ArticlePin
	case schema.ArticleOperationUnPin:
		articleInfo.Pin = entity.ArticleUnPin
	}

	err = qs.articleRepo.UpdateArticleOperation(ctx, articleInfo)
	if err != nil {
		return err
	}

	actMap := make(map[string]constant.ActivityTypeKey)
	actMap[schema.ArticleOperationPin] = constant.ActArticlePin
	actMap[schema.ArticleOperationUnPin] = constant.ActArticleUnPin
	actMap[schema.ArticleOperationHide] = constant.ActArticleHide
	actMap[schema.ArticleOperationShow] = constant.ActArticleShow
	_, ok := actMap[req.Operation]
	if ok {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           req.UserID,
			ObjectID:         articleInfo.ID,
			OriginalObjectID: articleInfo.ID,
			ActivityTypeKey:  actMap[req.Operation],
		})
	}

	return nil
}

// RemoveArticle delete article
func (qs *ArticleService) RemoveArticle(ctx context.Context, req *schema.RemoveArticleReq) (err error) {
	articleInfo, has, err := qs.articleRepo.GetArticle(ctx, req.ID)
	if err != nil {
		return err
	}
	//if the status is deleted, return directly
	if articleInfo.Status == entity.ArticleStatusDeleted {
		return nil
	}
	if !has {
		return nil
	}
	if !req.IsAdmin {
		if articleInfo.UserID != req.UserID {
			return errors.BadRequest(reason.ArticleCannotDeleted)
		}
		//
		//if articleInfo.AcceptedAnswerID != "0" {
		//	return errors.BadRequest(reason.ArticleCannotDeleted)
		//}
		//if articleInfo.AnswerCount > 1 {
		//	return errors.BadRequest(reason.ArticleCannotDeleted)
		//}

		//if articleInfo.AnswerCount == 1 {
		//	answersearch := &entity.AnswerSearch{}
		//	answersearch.ArticleID = req.ID
		//	answerList, _, err := qs.articlecommon.AnswerCommon.Search(ctx, answersearch)
		//	if err != nil {
		//		return err
		//	}
		//	for _, answer := range answerList {
		//		if answer.VoteCount > 0 {
		//			return errors.BadRequest(reason.ArticleCannotDeleted)
		//		}
		//	}
		//}
	}

	articleInfo.Status = entity.ArticleStatusDeleted
	err = qs.articleRepo.UpdateArticleStatusWithOutUpdateTime(ctx, articleInfo)
	if err != nil {
		return err
	}

	userArticleCount, err := qs.articlecommon.GetUserArticleCount(ctx, articleInfo.UserID)
	if err != nil {
		log.Error("user GetUserArticleCount error", err.Error())
	} else {
		err = qs.userCommon.UpdateArticleCount(ctx, articleInfo.UserID, userArticleCount)
		if err != nil {
			log.Error("user IncreaseArticleCount error", err.Error())
		}
	}

	//tag count
	tagIDs := make([]string, 0)
	Tags, tagerr := qs.tagCommon.GetObjectEntityTag(ctx, req.ID)
	if tagerr != nil {
		log.Error("GetObjectEntityTag error", tagerr)
		return nil
	}
	for _, v := range Tags {
		tagIDs = append(tagIDs, v.ID)
	}
	err = qs.tagCommon.RemoveTagRelListByObjectID(ctx, req.ID)
	if err != nil {
		log.Error("RemoveTagRelListByObjectID error", err.Error())
	}
	err = qs.tagCommon.RefreshTagArticleCount(ctx, tagIDs)
	if err != nil {
		log.Error("efreshTagArticleCount error", err.Error())
	}

	// #2372 In order to simplify the process and complexity, as well as to consider if it is in-house,
	// facing the problem of recovery.
	// err = qs.answerActivityService.DeleteArticle(ctx, articleInfo.ID, articleInfo.CreatedAt, articleInfo.VoteCount)
	// if err != nil {
	// 	 log.Errorf("user DeleteArticle rank rollback error %s", err.Error())
	// }
	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           articleInfo.UserID,
		TriggerUserID:    converter.StringToInt64(req.UserID),
		ObjectID:         articleInfo.ID,
		OriginalObjectID: articleInfo.ID,
		ActivityTypeKey:  constant.ActArticleDeleted,
	})
	qs.eventQueueService.Send(ctx, schema.NewEvent(constant.EventArticleDelete, req.UserID).TID(articleInfo.ID).
		QID(articleInfo.ID, articleInfo.UserID))
	return nil
}

func (qs *ArticleService) UpdateArticleCheckTags(ctx context.Context, req *schema.ArticleUpdate) (errorlist []*validator.FormErrorField, err error) {
	dbinfo, has, err := qs.articleRepo.GetArticle(ctx, req.ID)
	if err != nil {
		return
	}
	if !has {
		return
	}

	oldTags, tagerr := qs.tagCommon.GetObjectEntityTag(ctx, req.ID)
	if tagerr != nil {
		log.Error("GetObjectEntityTag error", tagerr)
		return nil, nil
	}

	tagNameList := make([]string, 0)
	oldtagNameList := make([]string, 0)
	for _, tag := range req.Tags {
		tagNameList = append(tagNameList, tag.SlugName)
	}
	for _, tag := range oldTags {
		oldtagNameList = append(oldtagNameList, tag.SlugName)
	}

	isChange := qs.tagCommon.CheckTagsIsChange(ctx, tagNameList, oldtagNameList)

	//If the content is the same, ignore it
	if dbinfo.Title == req.Title && dbinfo.OriginalText == req.Content && !isChange {
		return
	}

	Tags, tagerr := qs.tagCommon.GetTagListByNames(ctx, tagNameList)
	if tagerr != nil {
		log.Error("GetTagListByNames error", tagerr)
		return nil, nil
	}

	// if user can not use reserved tag, old reserved tag can not be removed and new reserved tag can not be added.
	if !req.CanUseReservedTag {
		CheckOldTag, CheckNewTag, CheckOldTaglist, CheckNewTaglist := qs.CheckChangeReservedTag(ctx, oldTags, Tags)
		if !CheckOldTag {
			errMsg := fmt.Sprintf(`The reserved tag "%s" must be present.`,
				strings.Join(CheckOldTaglist, ","))
			errorlist := make([]*validator.FormErrorField, 0)
			errorlist = append(errorlist, &validator.FormErrorField{
				ErrorField: "tags",
				ErrorMsg:   errMsg,
			})
			err = errors.BadRequest(reason.RequestFormatError).WithMsg(errMsg)
			return errorlist, err
		}
		if !CheckNewTag {
			errMsg := fmt.Sprintf(`"%s" can only be used by moderators.`,
				strings.Join(CheckNewTaglist, ","))
			errorlist := make([]*validator.FormErrorField, 0)
			errorlist = append(errorlist, &validator.FormErrorField{
				ErrorField: "tags",
				ErrorMsg:   errMsg,
			})
			err = errors.BadRequest(reason.RequestFormatError).WithMsg(errMsg)
			return errorlist, err
		}
	}
	return nil, nil
}

func (qs *ArticleService) RecoverArticle(ctx context.Context, req *schema.ArticleRecoverReq) (err error) {
	articleInfo, exist, err := qs.articleRepo.GetArticle(ctx, req.ArticleID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.ArticleNotFound)
	}
	if articleInfo.Status != entity.ArticleStatusDeleted {
		return nil
	}

	err = qs.articleRepo.RecoverArticle(ctx, req.ArticleID)
	if err != nil {
		return err
	}

	// update user's article count
	userArticleCount, err := qs.articlecommon.GetUserArticleCount(ctx, articleInfo.UserID)
	if err != nil {
		log.Error("user GetUserArticleCount error", err.Error())
	} else {
		err = qs.userCommon.UpdateArticleCount(ctx, articleInfo.UserID, userArticleCount)
		if err != nil {
			log.Error("user IncreaseArticleCount error", err.Error())
		}
	}

	// update tag's article count
	if err = qs.tagCommon.RecoverTagRelListByObjectID(ctx, articleInfo.ID); err != nil {
		log.Errorf("remove tag rel list by object id error %v", err)
	}

	tagIDs := make([]string, 0)
	tags, err := qs.tagCommon.GetObjectEntityTag(ctx, articleInfo.ID)
	if err != nil {
		return err
	}
	for _, v := range tags {
		tagIDs = append(tagIDs, v.ID)
	}
	if len(tagIDs) > 0 {
		if err = qs.tagCommon.RefreshTagQuestionCount(ctx, tagIDs); err != nil {
			log.Errorf("update tag's article count failed, %v", err)
		}
	}

	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           req.UserID,
		TriggerUserID:    converter.StringToInt64(req.UserID),
		ObjectID:         articleInfo.ID,
		OriginalObjectID: articleInfo.ID,
		ActivityTypeKey:  constant.ActArticleUndeleted,
	})
	return nil
}

func (qs *ArticleService) UpdateArticleInviteUser(ctx context.Context, req *schema.ArticleUpdateInviteUser) (err error) {
	return nil
	//originArticle, exist, err := qs.articleRepo.GetArticle(ctx, req.ID)
	//if err != nil {
	//	return err
	//}
	//if !exist {
	//	return errors.BadRequest(reason.ArticleNotFound)
	//}
	//
	////verify invite user
	//inviteUserInfoList, err := qs.userCommon.BatchGetUserBasicInfoByUserNames(ctx, req.InviteUser)
	//if err != nil {
	//	log.Error("BatchGetUserBasicInfoByUserNames error", err.Error())
	//}
	//inviteUserIDs := make([]string, 0)
	//for _, item := range req.InviteUser {
	//	_, ok := inviteUserInfoList[item]
	//	if ok {
	//		inviteUserIDs = append(inviteUserIDs, inviteUserInfoList[item].ID)
	//	}
	//}
	//inviteUserStr := ""
	//inviteUserByte, err := json.Marshal(inviteUserIDs)
	//if err != nil {
	//	log.Error("json.Marshal error", err.Error())
	//	inviteUserStr = "[]"
	//} else {
	//	inviteUserStr = string(inviteUserByte)
	//}
	//article := &entity.Article{}
	//article.ID = uid.DeShortID(req.ID)
	////article.InviteUserID = inviteUserStr
	//_ = inviteUserStr
	//
	//saveerr := qs.articleRepo.UpdateArticle(ctx, article, []string{"invite_user_id"})
	//if saveerr != nil {
	//	return saveerr
	//}
	////send notification
	//oldInviteUserIDsStr := originArticle.InviteUserID
	//oldInviteUserIDs := make([]string, 0)
	//needSendNotificationUserIDs := make([]string, 0)
	//if oldInviteUserIDsStr != "" {
	//	err = json.Unmarshal([]byte(oldInviteUserIDsStr), &oldInviteUserIDs)
	//	if err == nil {
	//		needSendNotificationUserIDs = converter.ArrayNotInArray(oldInviteUserIDs, inviteUserIDs)
	//	}
	//} else {
	//	needSendNotificationUserIDs = inviteUserIDs
	//}
	//go qs.notificationInviteUser(ctx, needSendNotificationUserIDs, originArticle.ID, originArticle.Title, req.UserID)
	//
	//return nil
}

func (qs *ArticleService) notificationInviteUser(
	ctx context.Context, invitedUserIDs []string, articleID, articleTitle, articleUserID string) {
	return
	//inviter, exist, err := qs.userCommon.GetUserBasicInfoByID(ctx, articleUserID)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//if !exist {
	//	log.Warnf("user %s not found", articleUserID)
	//	return
	//}
	//
	//users, err := qs.userRepo.BatchGetByID(ctx, invitedUserIDs)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//invitee := make(map[string]*entity.User, len(users))
	//for _, user := range users {
	//	invitee[user.ID] = user
	//}
	//for _, userID := range invitedUserIDs {
	//	msg := &schema.NotificationMsg{
	//		ReceiverUserID: userID,
	//		TriggerUserID:  articleUserID,
	//		Type:           schema.NotificationTypeInbox,
	//		ObjectID:       articleID,
	//	}
	//	msg.ObjectType = constant.ArticleObjectType
	//	msg.NotificationAction = constant.NotificationInvitedYouToAnswer
	//	qs.notificationQueueService.Send(ctx, msg)
	//
	//	receiverUserInfo, ok := invitee[userID]
	//	if !ok {
	//		log.Warnf("user %s not found", userID)
	//		return
	//	}
	//	externalNotificationMsg := &schema.ExternalNotificationMsg{
	//		ReceiverUserID: receiverUserInfo.ID,
	//		ReceiverEmail:  receiverUserInfo.EMail,
	//		ReceiverLang:   receiverUserInfo.Language,
	//	}
	//	rawData := &schema.NewInviteAnswerTemplateRawData{
	//		InviterDisplayName: inviter.DisplayName,
	//		ArticleTitle:       articleTitle,
	//		ArticleID:          articleID,
	//		UnsubscribeCode:    token.GenerateToken(),
	//	}
	//	externalNotificationMsg.NewInviteAnswerTemplateRawData = rawData
	//	qs.externalNotificationQueueService.Send(ctx, externalNotificationMsg)
	//}
}

// UpdateArticle update article
func (qs *ArticleService) UpdateArticle(ctx context.Context, req *schema.ArticleUpdate) (articleInfo any, err error) {
	var canUpdate bool
	articleInfo = &schema.ArticleInfoResp{}

	_, existUnreviewed, err := qs.revisionService.ExistUnreviewedByObjectID(ctx, req.ID)
	if err != nil {
		return
	}
	if existUnreviewed {
		err = errors.BadRequest(reason.ArticleCannotUpdate)
		return
	}

	dbinfo, has, err := qs.articleRepo.GetArticle(ctx, req.ID)
	if err != nil {
		return
	}
	if !has {
		return
	}
	if dbinfo.Status == entity.ArticleStatusDeleted {
		err = errors.BadRequest(reason.ArticleCannotUpdate)
		return nil, err
	}
	log.Infof("UpdateArticle b:%s", req.Content)
	log.Infof("UpdateArticle b html:%s", req.HTML)
	now := time.Now()
	article := &entity.Article{}
	article.Title = req.Title
	article.OriginalText = req.Content
	article.ParsedText = req.HTML
	article.ID = uid.DeShortID(req.ID)
	article.UpdatedAt = now
	article.PostUpdateTime = now
	article.UserID = dbinfo.UserID
	//article.LastEditUserID = req.UserID

	oldTags, tagerr := qs.tagCommon.GetObjectEntityTag(ctx, article.ID)
	if tagerr != nil {
		return articleInfo, tagerr
	}

	tagNameList := make([]string, 0)
	oldtagNameList := make([]string, 0)
	for _, tag := range req.Tags {
		tag.SlugName = strings.ReplaceAll(tag.SlugName, " ", "-")
		tagNameList = append(tagNameList, tag.SlugName)
	}
	for _, tag := range oldTags {
		oldtagNameList = append(oldtagNameList, tag.SlugName)
	}

	isChange := qs.tagCommon.CheckTagsIsChange(ctx, tagNameList, oldtagNameList)

	//If the content is the same, ignore it
	if dbinfo.Title == req.Title && dbinfo.OriginalText == req.Content && !isChange {
		return
	}

	Tags, tagerr := qs.tagCommon.GetTagListByNames(ctx, tagNameList)
	if tagerr != nil {
		return articleInfo, tagerr
	}

	// if user can not use reserved tag, old reserved tag can not be removed and new reserved tag can not be added.
	if !req.CanUseReservedTag {
		CheckOldTag, CheckNewTag, CheckOldTaglist, CheckNewTaglist := qs.CheckChangeReservedTag(ctx, oldTags, Tags)
		if !CheckOldTag {
			errMsg := fmt.Sprintf(`The reserved tag "%s" must be present.`,
				strings.Join(CheckOldTaglist, ","))
			errorlist := make([]*validator.FormErrorField, 0)
			errorlist = append(errorlist, &validator.FormErrorField{
				ErrorField: "tags",
				ErrorMsg:   errMsg,
			})
			err = errors.BadRequest(reason.RequestFormatError).WithMsg(errMsg)
			return errorlist, err
		}
		if !CheckNewTag {
			errMsg := fmt.Sprintf(`"%s" can only be used by moderators.`,
				strings.Join(CheckNewTaglist, ","))
			errorlist := make([]*validator.FormErrorField, 0)
			errorlist = append(errorlist, &validator.FormErrorField{
				ErrorField: "tags",
				ErrorMsg:   errMsg,
			})
			err = errors.BadRequest(reason.RequestFormatError).WithMsg(errMsg)
			return errorlist, err
		}
	}
	// Check whether mandatory labels are selected
	recommendExist, err := qs.tagCommon.ExistRecommend(ctx, req.Tags)
	if err != nil {
		return
	}
	if !recommendExist {
		errorlist := make([]*validator.FormErrorField, 0)
		errorlist = append(errorlist, &validator.FormErrorField{
			ErrorField: "tags",
			ErrorMsg:   translator.Tr(handler.GetLangByCtx(ctx), reason.RecommendTagEnter),
		})
		err = errors.BadRequest(reason.RecommendTagEnter)
		return errorlist, err
	}

	//Administrators and themselves do not need to be audited

	revisionDTO := &schema.AddRevisionDTO{
		UserID:   article.UserID,
		ObjectID: article.ID,
		Title:    article.Title,
		Log:      req.EditSummary,
	}

	if req.NoNeedReview {
		canUpdate = true
	}

	// It's not you or the administrator that needs to be reviewed
	if !canUpdate {
		revisionDTO.Status = entity.RevisionUnreviewedStatus
		revisionDTO.UserID = req.UserID //use revision userid
	} else {
		//Direct modification
		revisionDTO.Status = entity.RevisionReviewPassStatus
		//update article to db
		saveerr := qs.articleRepo.UpdateArticle(ctx, article, []string{"title", "original_text", "parsed_text", "updated_at", "post_update_time", "last_edit_user_id"})
		if saveerr != nil {
			return articleInfo, saveerr
		}
		objectTagData := schema.TagChange{}
		objectTagData.ObjectID = article.ID
		objectTagData.Tags = req.Tags
		objectTagData.UserID = req.UserID
		tagerr := qs.ChangeTag(ctx, &objectTagData)
		if err != nil {
			return articleInfo, tagerr
		}
	}

	articleWithTagsRevision, err := qs.changeArticleToRevision(ctx, article, Tags)
	if err != nil {
		return nil, err
	}
	infoJSON, _ := json.Marshal(articleWithTagsRevision)
	revisionDTO.Content = string(infoJSON)
	revisionID, err := qs.revisionService.AddRevision(ctx, revisionDTO, true)
	log.Infof("AddRevision revisionID:%+v err:%+v", revisionID, err)
	if err != nil {
		return
	}
	if canUpdate {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           req.UserID,
			ObjectID:         article.ID,
			ActivityTypeKey:  constant.ActArticleEdited,
			RevisionID:       revisionID,
			OriginalObjectID: article.ID,
		})
		qs.eventQueueService.Send(ctx, schema.NewEvent(constant.EventArticleUpdate, req.UserID).TID(article.ID).
			QID(article.ID, article.UserID))
	}

	articleInfo, err = qs.GetArticle(ctx, article.ID, article.UserID, req.ArticlePermission)
	return
}

// GetArticle get article one
func (qs *ArticleService) GetArticle(ctx context.Context, articleID, userID string,
	per schema.ArticlePermission) (resp *schema.ArticleInfoResp, err error) {
	article, err := qs.articlecommon.Info(ctx, articleID, userID)
	if err != nil {
		return
	}
	// If the article is deleted or pending, only the administrator and the author can view it
	if (article.Status == entity.ArticleStatusDeleted ||
		article.Status == entity.ArticleStatusPending) && !per.CanReopen && article.UserID != userID {
		return nil, errors.NotFound(reason.ArticleNotFound)
	}
	if article.Status != entity.ArticleStatusClosed {
		per.CanReopen = false
	}
	if article.Status == entity.ArticleStatusClosed {
		per.CanClose = false
	}
	if article.Pin == entity.ArticlePin {
		per.CanPin = false
		per.CanHide = false
	}
	if article.Pin == entity.ArticleUnPin {
		per.CanUnPin = false
	}
	if article.Show == entity.ArticleShow {
		per.CanShow = false
	}
	if article.Show == entity.ArticleHide {
		per.CanHide = false
		per.CanPin = false
	}

	if article.Status == entity.ArticleStatusDeleted {
		operation := &schema.Operation{}
		operation.Msg = translator.Tr(handler.GetLangByCtx(ctx), reason.ArticleAlreadyDeleted)
		operation.Level = schema.OperationLevelDanger
		article.Operation = operation
	}
	if article.Status == entity.ArticleStatusPending {
		operation := &schema.Operation{}
		operation.Msg = translator.Tr(handler.GetLangByCtx(ctx), reason.ArticleUnderReview)
		operation.Level = schema.OperationLevelSecondary
		article.Operation = operation
	}

	article.Description = htmltext.FetchExcerpt(article.HTML, "...", 240)
	article.MemberActions = permission.GetArticlePermission(ctx, userID, article.UserID, article.Status,
		per.CanEdit, per.CanDelete,
		per.CanClose, per.CanReopen, per.CanPin, per.CanHide, per.CanUnPin, per.CanShow,
		per.CanRecover)
	article.ExtendsActions = permission.GetArticleExtendsPermission(ctx, per.CanInviteOtherToAnswer)
	return article, nil
}

// GetArticleAndAddPV get article one
func (qs *ArticleService) GetArticleAndAddPV(ctx context.Context, articleID, loginUserID string,
	per schema.ArticlePermission) (
	resp *schema.ArticleInfoResp, err error) {
	err = qs.articlecommon.UpdatePv(ctx, articleID)
	if err != nil {
		log.Error(err)
	}
	return qs.GetArticle(ctx, articleID, loginUserID, per)
}

func (qs *ArticleService) InviteUserInfo(ctx context.Context, articleID string) (inviteList []*schema.UserBasicInfo, err error) {
	return qs.articlecommon.InviteUserInfo(ctx, articleID)
}

func (qs *ArticleService) ChangeTag(ctx context.Context, objectTagData *schema.TagChange) error {
	return qs.tagCommon.ObjectChangeTag(ctx, objectTagData)
}

func (qs *ArticleService) CheckChangeReservedTag(ctx context.Context, oldobjectTagData, objectTagData []*entity.Tag) (bool, bool, []string, []string) {
	return qs.tagCommon.CheckChangeReservedTag(ctx, oldobjectTagData, objectTagData)
}

// PersonalArticlePage get article list by user
func (qs *ArticleService) PersonalArticlePage(ctx context.Context, req *schema.PersonalArticlePageReq) (
	pageModel *pager.PageModel, err error) {

	userinfo, exist, err := qs.userCommon.GetUserBasicInfoByUserName(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.BadRequest(reason.UserNotFound)
	}
	search := &schema.ArticlePageReq{}
	search.OrderCond = req.OrderCond
	search.Page = req.Page
	search.PageSize = req.PageSize
	search.UserIDBeSearched = userinfo.ID
	search.LoginUserID = req.LoginUserID
	// Only author and administrator can view the pending article
	if req.LoginUserID == userinfo.ID || req.IsAdmin {
		search.ShowPending = true
	}
	articleList, total, err := qs.GetArticlePage(ctx, search)
	if err != nil {
		return nil, err
	}
	userArticleInfoList := make([]*schema.UserArticleInfo, 0)
	for _, item := range articleList {
		info := &schema.UserArticleInfo{}
		_ = copier.Copy(info, item)
		status, ok := entity.AdminArticleSearchStatusIntToString[item.Status]
		if ok {
			info.Status = status
		}
		userArticleInfoList = append(userArticleInfoList, info)
	}
	return pager.NewPageModel(total, userArticleInfoList), nil
}

func (qs *ArticleService) PersonalAnswerPage(ctx context.Context, req *schema.PersonalAnswerPageReq) (
	pageModel *pager.PageModel, err error) {
	return
	//userinfo, exist, err := qs.userCommon.GetUserBasicInfoByUserName(ctx, req.Username)
	//if err != nil {
	//	return nil, err
	//}
	//if !exist {
	//	return nil, errors.BadRequest(reason.UserNotFound)
	//}
	//cond := &entity.PersonalAnswerPageQueryCond{}
	//cond.UserID = userinfo.ID
	//cond.Page = req.Page
	//cond.PageSize = req.PageSize
	//cond.ShowPending = req.IsAdmin || req.LoginUserID == cond.UserID
	//if req.OrderCond == "newest" {
	//	cond.Order = entity.AnswerSearchOrderByTime
	//} else {
	//	cond.Order = entity.AnswerSearchOrderByDefault
	//}
	//articleIDs := make([]string, 0)
	//answerList, total, err := qs.articlecommon.AnswerCommon.PersonalAnswerPage(ctx, cond)
	//if err != nil {
	//	return nil, err
	//}
	//
	//answerlist := make([]*schema.AnswerInfo, 0)
	//userAnswerlist := make([]*schema.UserAnswerInfo, 0)
	//for _, item := range answerList {
	//	answerinfo := qs.articlecommon.AnswerCommon.ShowFormat(ctx, item)
	//	answerlist = append(answerlist, answerinfo)
	//	articleIDs = append(articleIDs, uid.DeShortID(item.ArticleID))
	//}
	//articleMaps, err := qs.articlecommon.FindInfoByID(ctx, articleIDs, req.LoginUserID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, item := range answerlist {
	//	_, ok := articleMaps[item.ArticleID]
	//	if ok {
	//		item.ArticleInfo = articleMaps[item.ArticleID]
	//	} else {
	//		continue
	//	}
	//	info := &schema.UserAnswerInfo{}
	//	_ = copier.Copy(info, item)
	//	info.AnswerID = item.ID
	//	info.ArticleID = item.ArticleID
	//	if item.ArticleInfo.Status == entity.ArticleStatusDeleted {
	//		info.ArticleInfo.Title = "Deleted article"
	//
	//	}
	//	userAnswerlist = append(userAnswerlist, info)
	//}
	//
	//return pager.NewPageModel(total, userAnswerlist), nil
}

// PersonalCollectionPage get collection list by user
func (qs *ArticleService) PersonalCollectionPage(ctx context.Context, req *schema.PersonalCollectionPageReq) (
	pageModel *pager.PageModel, err error) {
	list := make([]*schema.ArticleInfoResp, 0)
	collectionSearch := &entity.CollectionSearch{}
	collectionSearch.UserID = req.UserID
	collectionSearch.Page = req.Page
	collectionSearch.PageSize = req.PageSize
	collectionList, total, err := qs.collectionCommon.SearchList(ctx, collectionSearch)
	if err != nil {
		return nil, err
	}
	articleIDs := make([]string, 0)
	for _, item := range collectionList {
		articleIDs = append(articleIDs, item.ObjectID)
	}

	articleMaps, err := qs.articlecommon.FindInfoByID(ctx, articleIDs, req.UserID)
	if err != nil {
		return nil, err
	}
	for _, id := range articleIDs {
		if handler.GetEnableShortID(ctx) {
			id = uid.EnShortID(id)
		}
		_, ok := articleMaps[id]
		if ok {
			articleMaps[id].LastAnsweredUserInfo = nil
			articleMaps[id].UpdateUserInfo = nil
			articleMaps[id].Content = ""
			articleMaps[id].HTML = ""
			if articleMaps[id].Status == entity.ArticleStatusDeleted {
				articleMaps[id].Title = "Deleted article"
			}
			list = append(list, articleMaps[id])
		}
	}

	return pager.NewPageModel(total, list), nil
}

func (qs *ArticleService) SearchUserTopList(ctx context.Context, userName string, loginUserID string) ([]*schema.UserArticleInfo, []*schema.UserAnswerInfo, error) {

	answerlist := make([]*schema.AnswerInfo, 0)

	userAnswerlist := make([]*schema.UserAnswerInfo, 0)
	userArticlelist := make([]*schema.UserArticleInfo, 0)
	_ = answerlist
	return userArticlelist, userAnswerlist, nil
	//
	//userinfo, Exist, err := qs.userCommon.GetUserBasicInfoByUserName(ctx, userName)
	//if err != nil {
	//	return userArticlelist, userAnswerlist, err
	//}
	//if !Exist {
	//	return userArticlelist, userAnswerlist, nil
	//}
	//search := &schema.ArticlePageReq{}
	//search.OrderCond = "score"
	//search.Page = 0
	//search.PageSize = 5
	//search.UserIDBeSearched = userinfo.ID
	//search.LoginUserID = loginUserID
	//articlelist, _, err := qs.GetArticlePage(ctx, search)
	//if err != nil {
	//	return userArticlelist, userAnswerlist, err
	//}
	//answersearch := &entity.AnswerSearch{}
	//answersearch.UserID = userinfo.ID
	//answersearch.PageSize = 5
	//answersearch.Order = entity.AnswerSearchOrderByVote
	//articleIDs := make([]string, 0)
	//answerList, _, err := qs.articlecommon.AnswerCommon.Search(ctx, answersearch)
	//if err != nil {
	//	return userArticlelist, userAnswerlist, err
	//}
	//for _, item := range answerList {
	//	answerinfo := qs.articlecommon.AnswerCommon.ShowFormat(ctx, item)
	//	answerlist = append(answerlist, answerinfo)
	//	articleIDs = append(articleIDs, item.ArticleID)
	//}
	//articleMaps, err := qs.articlecommon.FindInfoByID(ctx, articleIDs, loginUserID)
	//if err != nil {
	//	return userArticlelist, userAnswerlist, err
	//}
	//for _, item := range answerlist {
	//	_, ok := articleMaps[item.ArticleID]
	//	if ok {
	//		item.ArticleInfo = articleMaps[item.ArticleID]
	//	}
	//}
	//
	//for _, item := range articlelist {
	//	info := &schema.UserArticleInfo{}
	//	_ = copier.Copy(info, item)
	//	info.UrlTitle = htmltext.UrlTitle(info.Title)
	//	userArticlelist = append(userArticlelist, info)
	//}
	//
	//for _, item := range answerlist {
	//	info := &schema.UserAnswerInfo{}
	//	_ = copier.Copy(info, item)
	//	info.AnswerID = item.ID
	//	info.ArticleID = item.ArticleID
	//	info.ArticleInfo.UrlTitle = htmltext.UrlTitle(info.ArticleInfo.Title)
	//	userAnswerlist = append(userAnswerlist, info)
	//}
	//
	//return userArticlelist, userAnswerlist, nil
}

// GetArticlesByTitle get articles by title
func (qs *ArticleService) GetArticlesByTitle(ctx context.Context, title string) (
	resp []*schema.ArticleBaseInfo, err error) {
	resp = make([]*schema.ArticleBaseInfo, 0)
	if len(title) == 0 {
		return resp, nil
	}
	articles, err := qs.articleRepo.GetArticlesByTitle(ctx, title, 10)
	if err != nil {
		return resp, err
	}
	for _, article := range articles {
		item := &schema.ArticleBaseInfo{}
		item.ID = article.ID
		item.Title = article.Title
		item.UrlTitle = htmltext.UrlTitle(article.Title)
		item.ViewCount = article.ViewCount
		//item.AnswerCount = article.AnswerCount
		item.CollectionCount = article.CollectionCount
		item.FollowCount = article.FollowCount
		status, ok := entity.AdminArticleSearchStatusIntToString[article.Status]
		if ok {
			item.Status = status
		}
		//if article.AcceptedAnswerID != "0" {
		//	item.AcceptedAnswer = true
		//}
		resp = append(resp, item)
	}
	return resp, nil
}

// SimilarArticle
func (qs *ArticleService) SimilarArticle(ctx context.Context, articleID string, loginUserID string) ([]*schema.ArticlePageResp, int64, error) {
	article, err := qs.articlecommon.Info(ctx, articleID, loginUserID)
	if err != nil {
		return nil, 0, nil
	}
	tagNames := make([]string, 0, len(article.Tags))
	for _, tag := range article.Tags {
		tagNames = append(tagNames, tag.SlugName)
	}
	search := &schema.ArticlePageReq{}
	search.OrderCond = "hot"
	search.Page = 0
	search.PageSize = 6
	if len(tagNames) > 0 {
		search.Tag = tagNames[0]
	}
	search.LoginUserID = loginUserID
	similarArticles, _, err := qs.GetArticlePage(ctx, search)
	if err != nil {
		return nil, 0, err
	}
	var result []*schema.ArticlePageResp
	for _, v := range similarArticles {
		if uid.DeShortID(v.ID) != articleID {
			result = append(result, v)
		}
	}
	return result, int64(len(result)), nil
}

// GetArticlePage query articles page
func (qs *ArticleService) GetArticlePage(ctx context.Context, req *schema.ArticlePageReq) (
	articles []*schema.ArticlePageResp, total int64, err error) {
	articles = make([]*schema.ArticlePageResp, 0)
	// query by user role
	showHidden := false
	if req.LoginUserID != "" && req.UserIDBeSearched != "" {
		showHidden = req.LoginUserID == req.UserIDBeSearched
		if !showHidden {
			userRole, err := qs.userRoleRelService.GetUserRole(ctx, req.LoginUserID)
			if err != nil {
				return nil, 0, err
			}
			showHidden = userRole == role.RoleAdminID || userRole == role.RoleModeratorID
		}
	}
	// query by tag condition
	var tagIDs = make([]string, 0)
	//@cws，如果有tag_id，按tag_id查，否则按tag name查
	if len(req.TagId) > 0 {
		req.TagId = strings.TrimSpace(req.TagId)
	}
	if len(req.TagId) > 0 {
		tagIDs = append(tagIDs, req.TagId) //@cws
	} else if len(req.Tag) > 0 {
		tagInfo, exist, err := qs.tagCommon.GetTagBySlugName(ctx, strings.ToLower(req.Tag))
		if err != nil {
			return nil, 0, err
		}
		if exist {
			synTagIds, err := qs.tagCommon.GetTagIDsByMainTagID(ctx, tagInfo.ID)
			if err != nil {
				return nil, 0, err
			}
			tagIDs = append(synTagIds, tagInfo.ID)
		}
	}

	// query by user condition
	if req.Username != "" {
		userinfo, exist, err := qs.userCommon.GetUserBasicInfoByUserName(ctx, req.Username)
		if err != nil {
			return nil, 0, err
		}
		if !exist {
			return articles, 0, nil
		}
		req.UserIDBeSearched = userinfo.ID
	}

	if req.OrderCond == schema.ArticleOrderCondHot {
		req.InDays = schema.HotInDays
	}

	articleList, total, err := qs.articleRepo.GetArticlePage(ctx, req.Page, req.PageSize,
		tagIDs, req.UserIDBeSearched, req.OrderCond, req.InDays, showHidden, req.ShowPending)
	if err != nil {
		return nil, 0, err
	}
	articles, err = qs.articlecommon.FormatArticlesPage(ctx, articleList, req.LoginUserID, req.OrderCond)
	if err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

// GetRecommendArticlePage retrieves recommended article page based on following tags and articles.
func (qs *ArticleService) GetRecommendArticlePage(ctx context.Context, req *schema.ArticlePageReq) (
	articles []*schema.ArticlePageResp, total int64, err error) {
	followingTagsResp, err := qs.tagService.GetFollowingTags(ctx, req.LoginUserID)
	if err != nil {
		return nil, 0, err
	}
	tagIDs := make([]string, 0, len(followingTagsResp))
	for _, tag := range followingTagsResp {
		tagIDs = append(tagIDs, tag.TagID)
	}

	activityType, err := qs.activityRepo.GetActivityTypeByObjectType(ctx, constant.ArticleObjectType, "follow")
	if err != nil {
		return nil, 0, err
	}
	activities, err := qs.activityRepo.GetUserActivitysByActivityType(ctx, req.LoginUserID, activityType)
	if err != nil {
		return nil, 0, err
	}

	followedArticleIDs := make([]string, 0, len(activities))
	for _, activity := range activities {
		if activity.Cancelled == entity.ActivityCancelled {
			continue
		}
		followedArticleIDs = append(followedArticleIDs, activity.ObjectID)
	}
	articleList, total, err := qs.articleRepo.GetRecommendArticlePageByTags(ctx, req.LoginUserID, tagIDs, followedArticleIDs, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	articles, err = qs.articlecommon.FormatArticlesPage(ctx, articleList, req.LoginUserID, "frequent")
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (qs *ArticleService) AdminSetArticleStatus(ctx context.Context, req *schema.AdminUpdateArticleStatusReq) error {
	setStatus, ok := entity.AdminArticleSearchStatus[req.Status]
	if !ok {
		return errors.BadRequest(reason.RequestFormatError)
	}
	articleInfo, exist, err := qs.articleRepo.GetArticle(ctx, req.ArticleID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.ArticleNotFound)
	}
	err = qs.articleRepo.UpdateArticleStatus(ctx, articleInfo.ID, setStatus)
	if err != nil {
		return err
	}

	msg := &schema.NotificationMsg{}
	if setStatus == entity.ArticleStatusDeleted {
		// #2372 In order to simplify the process and complexity, as well as to consider if it is in-house,
		// facing the problem of recovery.
		//err = qs.answerActivityService.DeleteArticle(ctx, articleInfo.ID, articleInfo.CreatedAt, articleInfo.VoteCount)
		//if err != nil {
		//	log.Errorf("admin delete article then rank rollback error %s", err.Error())
		//}
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           articleInfo.UserID,
			TriggerUserID:    converter.StringToInt64(req.UserID),
			ObjectID:         articleInfo.ID,
			OriginalObjectID: articleInfo.ID,
			ActivityTypeKey:  constant.ActArticleDeleted,
		})
		msg.NotificationAction = constant.NotificationYourArticleIsClosed
	}
	if setStatus == entity.ArticleStatusAvailable && articleInfo.Status == entity.ArticleStatusClosed {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           articleInfo.UserID,
			TriggerUserID:    converter.StringToInt64(req.UserID),
			ObjectID:         articleInfo.ID,
			OriginalObjectID: articleInfo.ID,
			ActivityTypeKey:  constant.ActArticleReopened,
		})
	}
	if setStatus == entity.ArticleStatusClosed && articleInfo.Status != entity.ArticleStatusClosed {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           articleInfo.UserID,
			TriggerUserID:    converter.StringToInt64(req.UserID),
			ObjectID:         articleInfo.ID,
			OriginalObjectID: articleInfo.ID,
			ActivityTypeKey:  constant.ActArticleClosed,
		})
		msg.NotificationAction = constant.NotificationYourArticleIsClosed
	}
	// recover
	if setStatus == entity.ArticleStatusAvailable && articleInfo.Status == entity.ArticleStatusDeleted {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           req.UserID,
			TriggerUserID:    converter.StringToInt64(req.UserID),
			ObjectID:         articleInfo.ID,
			OriginalObjectID: articleInfo.ID,
			ActivityTypeKey:  constant.ActArticleUndeleted,
		})
	}

	if len(msg.NotificationAction) > 0 {
		msg.ObjectID = articleInfo.ID
		msg.Type = schema.NotificationTypeInbox
		msg.ReceiverUserID = articleInfo.UserID
		msg.TriggerUserID = req.UserID
		msg.ObjectType = constant.ArticleObjectType
		qs.notificationQueueService.Send(ctx, msg)
	}
	return nil
}

func (qs *ArticleService) AdminArticlePage(
	ctx context.Context, req *schema.AdminArticlePageReq) (
	resp *pager.PageModel, err error) {

	list := make([]*schema.AdminArticleInfo, 0)
	articleList, count, err := qs.articleRepo.AdminArticlePage(ctx, req)
	if err != nil {
		return nil, err
	}

	userIds := make([]string, 0)
	for _, info := range articleList {
		item := &schema.AdminArticleInfo{}
		_ = copier.Copy(item, info)
		item.CreateTime = info.CreatedAt.Unix()
		item.UpdateTime = info.PostUpdateTime.Unix()
		item.EditTime = info.UpdatedAt.Unix()
		list = append(list, item)
		userIds = append(userIds, info.UserID)
	}
	userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIds)
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		if u, ok := userInfoMap[item.UserID]; ok {
			item.UserInfo = u
		}
	}
	return pager.NewPageModel(count, list), nil
}

// AdminAnswerPage search answer list
func (qs *ArticleService) AdminAnswerPage(ctx context.Context, req *schema.AdminAnswerPageReq) (
	resp *pager.PageModel, err error) {
	return
	//
	//answerList, count, err := qs.articlecommon.AnswerCommon.AdminSearchList(ctx, req)
	//if err != nil {
	//	return nil, err
	//}
	//
	//articleIDs := make([]string, 0)
	//userIds := make([]string, 0)
	//answerResp := make([]*schema.AdminAnswerInfo, 0)
	//for _, item := range answerList {
	//	answerInfo := qs.articlecommon.AnswerCommon.AdminShowFormat(ctx, item)
	//	answerResp = append(answerResp, answerInfo)
	//	articleIDs = append(articleIDs, item.ArticleID)
	//	userIds = append(userIds, item.UserID)
	//}
	//userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIds)
	//if err != nil {
	//	return nil, err
	//}
	//articleMaps, err := qs.articlecommon.FindInfoByID(ctx, articleIDs, req.LoginUserID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, item := range answerResp {
	//	if q, ok := articleMaps[item.ArticleID]; ok {
	//		item.ArticleInfo.Title = q.Title
	//	}
	//	if u, ok := userInfoMap[item.UserID]; ok {
	//		item.UserInfo = u
	//	}
	//}
	//return pager.NewPageModel(count, answerResp), nil
}

func (qs *ArticleService) changeArticleToRevision(ctx context.Context, articleInfo *entity.Article, tags []*entity.Tag) (
	articleRevision *entity.ArticleWithTagsRevision, err error) {
	articleRevision = &entity.ArticleWithTagsRevision{}
	articleRevision.Article = *articleInfo

	for _, tag := range tags {
		item := &entity.TagSimpleInfoForRevision{}
		_ = copier.Copy(item, tag)
		articleRevision.Tags = append(articleRevision.Tags, item)
	}
	return articleRevision, nil
}

func (qs *ArticleService) SitemapCron(ctx context.Context) {
	siteSeo, err := qs.siteInfoService.GetSiteSeo(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	ctx = context.WithValue(ctx, constant.ShortIDFlag, siteSeo.IsShortLink())
	qs.articlecommon.SitemapCron(ctx)
}
