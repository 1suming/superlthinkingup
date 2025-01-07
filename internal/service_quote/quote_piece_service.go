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

package service_quote

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
	quote_common "github.com/apache/incubator-answer/internal/service_quote/quote_common"
	"github.com/apache/incubator-answer/pkg/checker"
	"github.com/apache/incubator-answer/pkg/converter"
	"github.com/apache/incubator-answer/pkg/htmltext"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
	"golang.org/x/net/context"
)

// QuotePieceRepo quote repository

// QuotePieceService user service
type QuotePieceService struct {
	activityRepo                     activity_common.ActivityRepo
	QuotePieceRepo                   quote_common.QuotePieceRepo
	answerRepo                       answercommon.AnswerRepo
	tagCommon                        *tagcommon.TagCommonService
	tagService                       *tag.TagService
	QuotePieceCommon                 *quote_common.QuotePieceCommon
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
	newQuotePieceNotificationService *notification.ExternalNotificationService
	reviewService                    *review.ReviewService
	configService                    *config.ConfigService
	eventQueueService                event_queue.EventQueueService
}

func NewQuotePieceService(
	activityRepo activity_common.ActivityRepo,
	QuotePieceRepo quote_common.QuotePieceRepo,
	answerRepo answercommon.AnswerRepo,
	tagCommon *tagcommon.TagCommonService,
	tagService *tag.TagService,
	QuotePieceCommon *quote_common.QuotePieceCommon,
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
	newQuotePieceNotificationService *notification.ExternalNotificationService,
	reviewService *review.ReviewService,
	configService *config.ConfigService,
	eventQueueService event_queue.EventQueueService,
) *QuotePieceService {
	return &QuotePieceService{
		activityRepo:                     activityRepo,
		QuotePieceRepo:                   QuotePieceRepo,
		answerRepo:                       answerRepo,
		tagCommon:                        tagCommon,
		tagService:                       tagService,
		QuotePieceCommon:                 QuotePieceCommon,
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
		newQuotePieceNotificationService: newQuotePieceNotificationService,
		reviewService:                    reviewService,
		configService:                    configService,
		eventQueueService:                eventQueueService,
	}
}

func (qs *QuotePieceService) CloseQuotePiece(ctx context.Context, req *schema.CloseQuotePieceReq) error {
	quoteInfo, has, err := qs.QuotePieceRepo.GetQuotePiece(ctx, req.ID)
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

	quoteInfo.Status = entity.QuotePieceStatusClosed
	err = qs.QuotePieceRepo.UpdateQuotePieceStatus(ctx, quoteInfo.ID, quoteInfo.Status)
	if err != nil {
		return err
	}

	closeMeta, _ := json.Marshal(schema.CloseQuotePieceMeta{
		CloseType: req.CloseType,
		CloseMsg:  req.CloseMsg,
	})
	err = qs.metaService.AddMeta(ctx, req.ID, entity.QuotePieceCloseReasonKey, string(closeMeta))
	if err != nil {
		return err
	}

	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           req.UserID,
		ObjectID:         quoteInfo.ID,
		OriginalObjectID: quoteInfo.ID,
		ActivityTypeKey:  constant.ActQuotePieceClosed,
	})
	return nil
}

// ReopenQuotePiece reopen quote
func (qs *QuotePieceService) ReopenQuotePiece(ctx context.Context, req *schema.ReopenQuotePieceReq) error {
	quoteInfo, has, err := qs.QuotePieceRepo.GetQuotePiece(ctx, req.QuotePieceID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	quoteInfo.Status = entity.QuotePieceStatusAvailable
	err = qs.QuotePieceRepo.UpdateQuotePieceStatus(ctx, quoteInfo.ID, quoteInfo.Status)
	if err != nil {
		return err
	}
	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           req.UserID,
		ObjectID:         quoteInfo.ID,
		OriginalObjectID: quoteInfo.ID,
		ActivityTypeKey:  constant.ActQuotePieceReopened,
	})
	return nil
}

func (qs *QuotePieceService) AddQuotePieceCheckTags(ctx context.Context, Tags []*entity.Tag) ([]string, error) {
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
func (qs *QuotePieceService) CheckAddQuotePiece(ctx context.Context, req *schema.QuotePieceAdd) (errorlist any, err error) {
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
	if !req.QuotePiecePermission.CanUseReservedTag {
		taglist, err := qs.AddQuotePieceCheckTags(ctx, Tags)
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
func (qs *QuotePieceService) HasNewTag(ctx context.Context, tags []*schema.TagItem) (bool, error) {
	return qs.tagCommon.HasNewTag(ctx, tags)
}

// AddQuotePiece add quote
func (qs *QuotePieceService) AddQuotePiece(ctx context.Context, req *schema.QuotePieceAdd) (quoteInfo any, err error) {

	//if len(req.Tags) == 0 {
	//	errorlist := make([]*validator.FormErrorField, 0)
	//	errorlist = append(errorlist, &validator.FormErrorField{
	//		ErrorField: "tags",
	//		ErrorMsg:   translator.Tr(handler.GetLangByCtx(ctx), reason.TagNotFound),
	//	})
	//	err = errors.BadRequest(reason.RecommendTagEnter)
	//	return errorlist, err
	//}

	if len(req.Tags) != 0 { //@cws
		recommendExist, err := qs.tagCommon.ExistRecommend(ctx, req.Tags)
		if err != nil {
			return nil, err
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
	}
	//@cws 允许tags为0

	tagNameList := make([]string, 0)
	for _, tag := range req.Tags {
		tag.SlugName = strings.ReplaceAll(tag.SlugName, " ", "-")
		tagNameList = append(tagNameList, tag.SlugName)
	}
	tags, tagerr := qs.tagCommon.GetTagListByNames(ctx, tagNameList)
	if tagerr != nil {
		return quoteInfo, tagerr
	}
	if !req.QuotePiecePermission.CanUseReservedTag {
		taglist, err := qs.AddQuotePieceCheckTags(ctx, tags)
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

	quote := &entity.QuotePiece{}
	now := time.Now()
	quote.UserID = req.UserID
	quote.Title = req.Title
	//quote.OriginalText = req.Content
	//quote.Bio = req.Content

	quote.OriginalText = req.Content
	quote.ParsedText = req.HTML

	log.Infof("addQuotePiece content:%s", req.Content)
	//quote.ParsedText = req.HTML
	//quote.AcceptedAnswerID = "0"
	//quote.LastAnswerID = "0"
	//quote.LastEditUserID = "0"
	//quote.PostUpdateTime = nil
	quote.Status = entity.QuotePieceStatusPending
	//quote.RevisionID = "0"
	quote.CreatedAt = now
	//quote.PostUpdateTime = now
	quote.Pin = entity.QuotePieceUnPin
	quote.Show = entity.QuotePieceShow
	quote.UpdatedAt = now

	//@cws
	//quote.PostDate = now
	//quote.OriginalTextFormat = req.ContentFormat

	err = qs.QuotePieceRepo.AddQuotePiece(ctx, quote)
	if err != nil {
		return
	}
	quote.Status = qs.reviewService.AddQuotePieceReview(ctx, quote, req.Tags, req.IP, req.UserAgent)
	if err := qs.QuotePieceRepo.UpdateQuotePieceStatus(ctx, quote.ID, quote.Status); err != nil {
		return nil, err
	}
	objectTagData := schema.TagChange{}
	objectTagData.ObjectID = quote.ID
	objectTagData.Tags = req.Tags
	objectTagData.UserID = req.UserID
	err = qs.ChangeTag(ctx, &objectTagData)
	if err != nil {
		return
	}
	_ = qs.QuotePieceRepo.UpdateSearch(ctx, quote.ID)
	//
	//revisionDTO := &schema.AddRevisionDTO{
	//	UserID:   quote.UserID,
	//	ObjectID: quote.ID,
	//	Title:    quote.Title,
	//}
	//
	//quoteWithTagsRevision, err := qs.changeQuotePieceToRevision(ctx, quote, tags)
	//if err != nil {
	//	return nil, err
	//}
	//infoJSON, _ := json.Marshal(quoteWithTagsRevision)
	//revisionDTO.Content = string(infoJSON)
	//revisionID, err := qs.revisionService.AddRevision(ctx, revisionDTO, true)
	//if err != nil {
	//	return
	//}

	//// user add quote count
	//userQuotePieceCount, err := qs.QuotePieceCommon.GetUserQuotePieceCount(ctx, quote.UserID)
	//if err != nil {
	//	log.Errorf("get user quote count error %v", err)
	//} else {
	//	err = qs.userCommon.UpdateQuoteCount(ctx, quote.UserID, userQuotePieceCount)
	//	if err != nil {
	//		log.Errorf("update user quote count error %v", err)
	//	}
	//}
	//
	//qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
	//	UserID:           quote.UserID,
	//	ObjectID:         quote.ID,
	//	OriginalObjectID: quote.ID,
	//	ActivityTypeKey:  constant.ActQuotePieceAsked,
	//	RevisionID:       revisionID,
	//})
	//
	//if quote.Status == entity.QuotePieceStatusAvailable {
	//	qs.externalNotificationQueueService.Send(ctx,
	//		schema.CreateNewQuotePieceNotificationMsg(quote.ID, quote.Title, quote.UserID, tags))
	//}
	//qs.eventQueueService.Send(ctx, schema.NewEvent(constant.EventQuoteCreate, req.UserID).TID(quote.ID).
	//	QID(quote.ID, quote.UserID))

	quoteInfo, err = qs.GetQuotePiece(ctx, quote.ID, quote.UserID, req.QuotePiecePermission) //
	return
}

// OperationQuotePiece
func (qs *QuotePieceService) OperationQuotePiece(ctx context.Context, req *schema.OperationQuotePieceReq) (err error) {
	quoteInfo, has, err := qs.QuotePieceRepo.GetQuotePiece(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	// Hidden quote cannot be placed at the top
	if quoteInfo.Show == entity.QuotePieceHide && req.Operation == schema.QuotePieceOperationPin {
		return nil
	}
	// QuotePiece cannot be hidden when they are at the top
	if quoteInfo.Pin == entity.QuotePiecePin && req.Operation == schema.QuotePieceOperationHide {
		return nil
	}

	switch req.Operation {
	case schema.QuotePieceOperationHide:
		quoteInfo.Show = entity.QuotePieceHide
		err = qs.tagCommon.HideTagRelListByObjectID(ctx, req.ID)
		if err != nil {
			return err
		}
		err = qs.tagCommon.RefreshTagCountByArticleID(ctx, req.ID)
		if err != nil {
			return err
		}
	case schema.QuotePieceOperationShow:
		quoteInfo.Show = entity.QuotePieceShow
		err = qs.tagCommon.ShowTagRelListByObjectID(ctx, req.ID)
		if err != nil {
			return err
		}
		err = qs.tagCommon.RefreshTagCountByQuoteID(ctx, req.ID)
		if err != nil {
			return err
		}
	case schema.QuotePieceOperationPin:
		quoteInfo.Pin = entity.QuotePiecePin
	case schema.QuotePieceOperationUnPin:
		quoteInfo.Pin = entity.QuotePieceUnPin
	}

	err = qs.QuotePieceRepo.UpdateQuotePieceOperation(ctx, quoteInfo)
	if err != nil {
		return err
	}

	actMap := make(map[string]constant.ActivityTypeKey)
	actMap[schema.QuotePieceOperationPin] = constant.ActQuotePiecePin
	actMap[schema.QuotePieceOperationUnPin] = constant.ActQuotePieceUnPin
	actMap[schema.QuotePieceOperationHide] = constant.ActQuotePieceHide
	actMap[schema.QuotePieceOperationShow] = constant.ActQuotePieceShow
	_, ok := actMap[req.Operation]
	if ok {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           req.UserID,
			ObjectID:         quoteInfo.ID,
			OriginalObjectID: quoteInfo.ID,
			ActivityTypeKey:  actMap[req.Operation],
		})
	}

	return nil
}

// RemoveQuotePiece delete quote
func (qs *QuotePieceService) RemoveQuotePiece(ctx context.Context, req *schema.RemoveQuotePieceReq) (err error) {
	quoteInfo, has, err := qs.QuotePieceRepo.GetQuotePiece(ctx, req.ID)
	if err != nil {
		return err
	}
	//if the status is deleted, return directly
	if quoteInfo.Status == entity.QuotePieceStatusDeleted {
		return nil
	}
	if !has {
		return nil
	}
	if !req.IsAdmin {
		if quoteInfo.UserID != req.UserID {
			return errors.BadRequest(reason.ArticleCannotDeleted)
		}
		//
		//if quoteInfo.AcceptedAnswerID != "0" {
		//	return errors.BadRequest(reason.QuotePieceCannotDeleted)
		//}
		//if quoteInfo.AnswerCount > 1 {
		//	return errors.BadRequest(reason.QuotePieceCannotDeleted)
		//}

		//if quoteInfo.AnswerCount == 1 {
		//	answersearch := &entity.AnswerSearch{}
		//	answersearch.QuotePieceID = req.ID
		//	answerList, _, err := qs.QuotePieceCommon.AnswerCommon.Search(ctx, answersearch)
		//	if err != nil {
		//		return err
		//	}
		//	for _, answer := range answerList {
		//		if answer.VoteCount > 0 {
		//			return errors.BadRequest(reason.QuotePieceCannotDeleted)
		//		}
		//	}
		//}
	}

	quoteInfo.Status = entity.QuotePieceStatusDeleted
	err = qs.QuotePieceRepo.UpdateQuotePieceStatusWithOutUpdateTime(ctx, quoteInfo)
	if err != nil {
		return err
	}

	//userQuotePieceCount, err := qs.QuotePieceCommon.GetUserQuotePieceCount(ctx, quoteInfo.UserID)
	//if err != nil {
	//	log.Error("user GetUserQuotePieceCount error", err.Error())
	//} else {
	//	err = qs.userCommon.UpdateQuotePieceCount(ctx, quoteInfo.UserID, userQuotePieceCount)
	//	if err != nil {
	//		log.Error("user IncreaseQuotePieceCount error", err.Error())
	//	}
	//}

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
	err = qs.tagCommon.RefreshTagQuoteCount(ctx, tagIDs)
	if err != nil {
		log.Error("efreshTagQuotePieceCount error", err.Error())
	}

	// #2372 In order to simplify the process and complexity, as well as to consider if it is in-house,
	// facing the problem of recovery.
	// err = qs.answerActivityService.DeleteQuotePiece(ctx, quoteInfo.ID, quoteInfo.CreatedAt, quoteInfo.VoteCount)
	// if err != nil {
	// 	 log.Errorf("user DeleteQuotePiece rank rollback error %s", err.Error())
	// }
	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           quoteInfo.UserID,
		TriggerUserID:    converter.StringToInt64(req.UserID),
		ObjectID:         quoteInfo.ID,
		OriginalObjectID: quoteInfo.ID,
		ActivityTypeKey:  constant.ActQuotePieceDeleted,
	})
	qs.eventQueueService.Send(ctx, schema.NewEvent(constant.EventQuotePieceDelete, req.UserID).TID(quoteInfo.ID).
		QID(quoteInfo.ID, quoteInfo.UserID))
	return nil
}

func (qs *QuotePieceService) UpdateQuotePieceCheckTags(ctx context.Context, req *schema.QuotePieceUpdate) (errorlist []*validator.FormErrorField, err error) {
	dbinfo, has, err := qs.QuotePieceRepo.GetQuotePiece(ctx, req.ID)
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

func (qs *QuotePieceService) RecoverQuotePiece(ctx context.Context, req *schema.QuotePieceRecoverReq) (err error) {
	quoteInfo, exist, err := qs.QuotePieceRepo.GetQuotePiece(ctx, req.QuotePieceID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.QuotePieceNotFound)
	}
	if quoteInfo.Status != entity.QuotePieceStatusDeleted {
		return nil
	}

	err = qs.QuotePieceRepo.RecoverQuotePiece(ctx, req.QuotePieceID)
	if err != nil {
		return err
	}

	// update user's quote count
	//userQuotePieceCount, err := qs.QuotePieceCommon.GetUserQuotePieceCount(ctx, quoteInfo.UserID)
	//if err != nil {
	//	log.Error("user GetUserQuotePieceCount error", err.Error())
	//} else {
	//	err = qs.userCommon.UpdateQuotePieceCount(ctx, quoteInfo.UserID, userQuotePieceCount)
	//	if err != nil {
	//		log.Error("user IncreaseQuotePieceCount error", err.Error())
	//	}
	//}

	// update tag's quote count
	if err = qs.tagCommon.RecoverTagRelListByObjectID(ctx, quoteInfo.ID); err != nil {
		log.Errorf("remove tag rel list by object id error %v", err)
	}

	tagIDs := make([]string, 0)
	tags, err := qs.tagCommon.GetObjectEntityTag(ctx, quoteInfo.ID)
	if err != nil {
		return err
	}
	for _, v := range tags {
		tagIDs = append(tagIDs, v.ID)
	}
	if len(tagIDs) > 0 {
		if err = qs.tagCommon.RefreshTagQuestionCount(ctx, tagIDs); err != nil {
			log.Errorf("update tag's quote count failed, %v", err)
		}
	}

	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           req.UserID,
		TriggerUserID:    converter.StringToInt64(req.UserID),
		ObjectID:         quoteInfo.ID,
		OriginalObjectID: quoteInfo.ID,
		ActivityTypeKey:  constant.ActQuotePieceUndeleted,
	})
	return nil
}

func (qs *QuotePieceService) UpdateQuotePieceInviteUser(ctx context.Context, req *schema.QuotePieceUpdateInviteUser) (err error) {
	return nil
	//originQuotePiece, exist, err := qs.QuotePieceRepo.GetQuotePiece(ctx, req.ID)
	//if err != nil {
	//	return err
	//}
	//if !exist {
	//	return errors.BadRequest(reason.QuotePieceNotFound)
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
	//quote := &entity.QuotePiece{}
	//quote.ID = uid.DeShortID(req.ID)
	////quote.InviteUserID = inviteUserStr
	//_ = inviteUserStr
	//
	//saveerr := qs.QuotePieceRepo.UpdateQuotePiece(ctx, quote, []string{"invite_user_id"})
	//if saveerr != nil {
	//	return saveerr
	//}
	////send notification
	//oldInviteUserIDsStr := originQuotePiece.InviteUserID
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
	//go qs.notificationInviteUser(ctx, needSendNotificationUserIDs, originQuotePiece.ID, originQuotePiece.Title, req.UserID)
	//
	//return nil
}

func (qs *QuotePieceService) notificationInviteUser(
	ctx context.Context, invitedUserIDs []string, quoteID, QuotePieceName, quoteUserID string) {
	return
	//inviter, exist, err := qs.userCommon.GetUserBasicInfoByID(ctx, quoteUserID)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//if !exist {
	//	log.Warnf("user %s not found", quoteUserID)
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
	//		TriggerUserID:  quoteUserID,
	//		Type:           schema.NotificationTypeInbox,
	//		ObjectID:       quoteID,
	//	}
	//	msg.ObjectType = constant.QuotePieceObjectType
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
	//		QuotePieceTitle:       QuotePieceName,
	//		QuotePieceID:          quoteID,
	//		UnsubscribeCode:    token.GenerateToken(),
	//	}
	//	externalNotificationMsg.NewInviteAnswerTemplateRawData = rawData
	//	qs.externalNotificationQueueService.Send(ctx, externalNotificationMsg)
	//}
}

// UpdateQuotePiece update quote
func (qs *QuotePieceService) UpdateQuotePiece(ctx context.Context, req *schema.QuotePieceUpdate) (quoteInfo any, err error) {
	var canUpdate bool
	quoteInfo = &schema.QuotePieceInfoResp{}

	_, existUnreviewed, err := qs.revisionService.ExistUnreviewedByObjectID(ctx, req.ID)
	if err != nil {
		return
	}
	if existUnreviewed {
		err = errors.BadRequest(reason.QuotePieceCannotUpdate)
		return
	}

	dbinfo, has, err := qs.QuotePieceRepo.GetQuotePiece(ctx, req.ID)
	if err != nil {
		return
	}
	if !has {
		return
	}
	if dbinfo.Status == entity.QuotePieceStatusDeleted {
		err = errors.BadRequest(reason.QuotePieceCannotUpdate)
		return nil, err
	}
	log.Infof("UpdateQuotePiece b:%s", req.Content)
	log.Infof("UpdateQuotePiece b html:%s", req.HTML)
	now := time.Now()
	quote := &entity.QuotePiece{}
	quote.Title = req.Title
	//quote.OriginalText = req.Content
	quote.OriginalText = req.Content
	//quote.ParsedText = req.HTML
	quote.ID = uid.DeShortID(req.ID)
	quote.UpdatedAt = now
	//quote.PostUpdateTime = now
	quote.UserID = dbinfo.UserID
	//quote.LastEditUserID = req.UserID

	oldTags, tagerr := qs.tagCommon.GetObjectEntityTag(ctx, quote.ID)
	if tagerr != nil {
		return quoteInfo, tagerr
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
		return quoteInfo, tagerr
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
		UserID:   quote.UserID,
		ObjectID: quote.ID,
		Title:    quote.Title,
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
		//update quote to db
		saveerr := qs.QuotePieceRepo.UpdateQuotePiece(ctx, quote, []string{"title", "original_text", "parsed_text", "updated_at", "post_update_time", "last_edit_user_id"})
		if saveerr != nil {
			return quoteInfo, saveerr
		}
		objectTagData := schema.TagChange{}
		objectTagData.ObjectID = quote.ID
		objectTagData.Tags = req.Tags
		objectTagData.UserID = req.UserID
		tagerr := qs.ChangeTag(ctx, &objectTagData)
		if err != nil {
			return quoteInfo, tagerr
		}
	}

	quoteWithTagsRevision, err := qs.changeQuotePieceToRevision(ctx, quote, Tags)
	if err != nil {
		return nil, err
	}
	infoJSON, _ := json.Marshal(quoteWithTagsRevision)
	revisionDTO.Content = string(infoJSON)
	revisionID, err := qs.revisionService.AddRevision(ctx, revisionDTO, true)
	log.Infof("AddRevision revisionID:%+v err:%+v", revisionID, err)
	if err != nil {
		return
	}
	if canUpdate {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           req.UserID,
			ObjectID:         quote.ID,
			ActivityTypeKey:  constant.ActQuotePieceEdited,
			RevisionID:       revisionID,
			OriginalObjectID: quote.ID,
		})
		qs.eventQueueService.Send(ctx, schema.NewEvent(constant.EventQuotePieceUpdate, req.UserID).TID(quote.ID).
			QID(quote.ID, quote.UserID))
	}

	quoteInfo, err = qs.GetQuotePiece(ctx, quote.ID, quote.UserID, req.QuotePiecePermission)
	return
}

// GetQuotePiece get quote one
func (qs *QuotePieceService) GetQuotePiece(ctx context.Context, quoteID, userID string,
	per schema.QuotePiecePermission) (resp *schema.QuotePieceInfoResp, err error) {
	quote, err := qs.QuotePieceCommon.Info(ctx, quoteID, userID)
	if err != nil {
		return
	}
	// If the quote is deleted or pending, only the administrator and the author can view it
	if (quote.Status == entity.QuotePieceStatusDeleted ||
		quote.Status == entity.QuotePieceStatusPending) && !per.CanReopen && quote.UserID != userID {
		return nil, errors.NotFound(reason.QuotePieceNotFound)
	}
	if quote.Status != entity.QuotePieceStatusClosed {
		per.CanReopen = false
	}
	if quote.Status == entity.QuotePieceStatusClosed {
		per.CanClose = false
	}
	if quote.Pin == entity.QuotePiecePin {
		per.CanPin = false
		per.CanHide = false
	}
	if quote.Pin == entity.QuotePieceUnPin {
		per.CanUnPin = false
	}
	if quote.Show == entity.QuotePieceShow {
		per.CanShow = false
	}
	if quote.Show == entity.QuotePieceHide {
		per.CanHide = false
		per.CanPin = false
	}

	if quote.Status == entity.QuotePieceStatusDeleted {
		operation := &schema.Operation{}
		operation.Msg = translator.Tr(handler.GetLangByCtx(ctx), reason.QuotePieceAlreadyDeleted)
		operation.Level = schema.OperationLevelDanger
		quote.Operation = operation
	}
	if quote.Status == entity.QuotePieceStatusPending {
		operation := &schema.Operation{}
		operation.Msg = translator.Tr(handler.GetLangByCtx(ctx), reason.QuotePieceUnderReview)
		operation.Level = schema.OperationLevelSecondary
		quote.Operation = operation
	}

	quote.Description = htmltext.FetchExcerpt(quote.HTML, "...", 240)
	quote.MemberActions = permission.GetQuotePiecePermission(ctx, userID, quote.UserID, quote.Status,
		per.CanEdit, per.CanDelete,
		per.CanClose, per.CanReopen, per.CanPin, per.CanHide, per.CanUnPin, per.CanShow,
		per.CanRecover)
	quote.ExtendsActions = permission.GetQuotePieceExtendsPermission(ctx, per.CanInviteOtherToAnswer)
	return quote, nil
}

// GetQuotePieceAndAddPV get quote one
func (qs *QuotePieceService) GetQuotePieceAndAddPV(ctx context.Context, quoteID, loginUserID string,
	per schema.QuotePiecePermission) (
	resp *schema.QuotePieceInfoResp, err error) {
	err = qs.QuotePieceCommon.UpdatePv(ctx, quoteID)
	if err != nil {
		log.Error(err)
	}
	return qs.GetQuotePiece(ctx, quoteID, loginUserID, per)
}

func (qs *QuotePieceService) InviteUserInfo(ctx context.Context, quoteID string) (inviteList []*schema.UserBasicInfo, err error) {
	return qs.QuotePieceCommon.InviteUserInfo(ctx, quoteID)
}

func (qs *QuotePieceService) ChangeTag(ctx context.Context, objectTagData *schema.TagChange) error {
	return qs.tagCommon.ObjectChangeTag(ctx, objectTagData)
}

func (qs *QuotePieceService) CheckChangeReservedTag(ctx context.Context, oldobjectTagData, objectTagData []*entity.Tag) (bool, bool, []string, []string) {
	return qs.tagCommon.CheckChangeReservedTag(ctx, oldobjectTagData, objectTagData)
}

// PersonalQuotePiecePage get quote list by user
func (qs *QuotePieceService) PersonalQuotePiecePage(ctx context.Context, req *schema.PersonalQuotePiecePageReq) (
	pageModel *pager.PageModel, err error) {

	userinfo, exist, err := qs.userCommon.GetUserBasicInfoByUserName(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.BadRequest(reason.UserNotFound)
	}
	search := &schema.QuotePiecePageReq{}
	search.OrderCond = req.OrderCond
	search.Page = req.Page
	search.PageSize = req.PageSize
	search.UserIDBeSearched = userinfo.ID
	search.LoginUserID = req.LoginUserID
	// Only author and administrator can view the pending quote
	if req.LoginUserID == userinfo.ID || req.IsAdmin {
		search.ShowPending = true
	}
	quoteList, total, err := qs.GetQuotePiecePage(ctx, search)
	if err != nil {
		return nil, err
	}
	userQuotePieceInfoList := make([]*schema.UserQuotePieceInfo, 0)
	for _, item := range quoteList {
		info := &schema.UserQuotePieceInfo{}
		_ = copier.Copy(info, item)
		status, ok := entity.AdminQuotePieceSearchStatusIntToString[item.Status]
		if ok {
			info.Status = status
		}
		userQuotePieceInfoList = append(userQuotePieceInfoList, info)
	}
	return pager.NewPageModel(total, userQuotePieceInfoList), nil
}

func (qs *QuotePieceService) PersonalAnswerPage(ctx context.Context, req *schema.PersonalAnswerPageReq) (
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
	//quoteIDs := make([]string, 0)
	//answerList, total, err := qs.QuotePieceCommon.AnswerCommon.PersonalAnswerPage(ctx, cond)
	//if err != nil {
	//	return nil, err
	//}
	//
	//answerlist := make([]*schema.AnswerInfo, 0)
	//userAnswerlist := make([]*schema.UserAnswerInfo, 0)
	//for _, item := range answerList {
	//	answerinfo := qs.QuotePieceCommon.AnswerCommon.ShowFormat(ctx, item)
	//	answerlist = append(answerlist, answerinfo)
	//	quoteIDs = append(quoteIDs, uid.DeShortID(item.QuotePieceID))
	//}
	//quoteMaps, err := qs.QuotePieceCommon.FindInfoByID(ctx, quoteIDs, req.LoginUserID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, item := range answerlist {
	//	_, ok := quoteMaps[item.QuotePieceID]
	//	if ok {
	//		item.QuotePieceInfo = quoteMaps[item.QuotePieceID]
	//	} else {
	//		continue
	//	}
	//	info := &schema.UserAnswerInfo{}
	//	_ = copier.Copy(info, item)
	//	info.AnswerID = item.ID
	//	info.QuotePieceID = item.QuotePieceID
	//	if item.QuotePieceInfo.Status == entity.QuotePieceStatusDeleted {
	//		info.QuotePieceInfo.Title = "Deleted quote"
	//
	//	}
	//	userAnswerlist = append(userAnswerlist, info)
	//}
	//
	//return pager.NewPageModel(total, userAnswerlist), nil
}

// PersonalCollectionPage get collection list by user
func (qs *QuotePieceService) PersonalCollectionPage(ctx context.Context, req *schema.PersonalCollectionPageReq) (
	pageModel *pager.PageModel, err error) {
	list := make([]*schema.QuotePieceInfoResp, 0)
	collectionSearch := &entity.CollectionSearch{}
	collectionSearch.UserID = req.UserID
	collectionSearch.Page = req.Page
	collectionSearch.PageSize = req.PageSize
	collectionList, total, err := qs.collectionCommon.SearchList(ctx, collectionSearch)
	if err != nil {
		return nil, err
	}
	quoteIDs := make([]string, 0)
	for _, item := range collectionList {
		quoteIDs = append(quoteIDs, item.ObjectID)
	}

	quoteMaps, err := qs.QuotePieceCommon.FindInfoByID(ctx, quoteIDs, req.UserID)
	if err != nil {
		return nil, err
	}
	for _, id := range quoteIDs {
		if handler.GetEnableShortID(ctx) {
			id = uid.EnShortID(id)
		}
		_, ok := quoteMaps[id]
		if ok {
			quoteMaps[id].LastAnsweredUserInfo = nil
			quoteMaps[id].UpdateUserInfo = nil
			quoteMaps[id].Content = ""
			quoteMaps[id].HTML = ""
			if quoteMaps[id].Status == entity.QuotePieceStatusDeleted {
				quoteMaps[id].Title = "Deleted quote"
			}
			list = append(list, quoteMaps[id])
		}
	}

	return pager.NewPageModel(total, list), nil
}

func (qs *QuotePieceService) SearchUserTopList(ctx context.Context, userName string, loginUserID string) ([]*schema.UserQuotePieceInfo, []*schema.UserAnswerInfo, error) {

	answerlist := make([]*schema.AnswerInfo, 0)

	userAnswerlist := make([]*schema.UserAnswerInfo, 0)
	userQuotePiecelist := make([]*schema.UserQuotePieceInfo, 0)
	_ = answerlist
	return userQuotePiecelist, userAnswerlist, nil
	//
	//userinfo, Exist, err := qs.userCommon.GetUserBasicInfoByUserName(ctx, userName)
	//if err != nil {
	//	return userQuotePiecelist, userAnswerlist, err
	//}
	//if !Exist {
	//	return userQuotePiecelist, userAnswerlist, nil
	//}
	//search := &schema.QuotePiecePageReq{}
	//search.OrderCond = "score"
	//search.Page = 0
	//search.PageSize = 5
	//search.UserIDBeSearched = userinfo.ID
	//search.LoginUserID = loginUserID
	//quotelist, _, err := qs.GetQuotePiecePage(ctx, search)
	//if err != nil {
	//	return userQuotePiecelist, userAnswerlist, err
	//}
	//answersearch := &entity.AnswerSearch{}
	//answersearch.UserID = userinfo.ID
	//answersearch.PageSize = 5
	//answersearch.Order = entity.AnswerSearchOrderByVote
	//quoteIDs := make([]string, 0)
	//answerList, _, err := qs.QuotePieceCommon.AnswerCommon.Search(ctx, answersearch)
	//if err != nil {
	//	return userQuotePiecelist, userAnswerlist, err
	//}
	//for _, item := range answerList {
	//	answerinfo := qs.QuotePieceCommon.AnswerCommon.ShowFormat(ctx, item)
	//	answerlist = append(answerlist, answerinfo)
	//	quoteIDs = append(quoteIDs, item.QuotePieceID)
	//}
	//quoteMaps, err := qs.QuotePieceCommon.FindInfoByID(ctx, quoteIDs, loginUserID)
	//if err != nil {
	//	return userQuotePiecelist, userAnswerlist, err
	//}
	//for _, item := range answerlist {
	//	_, ok := quoteMaps[item.QuotePieceID]
	//	if ok {
	//		item.QuotePieceInfo = quoteMaps[item.QuotePieceID]
	//	}
	//}
	//
	//for _, item := range quotelist {
	//	info := &schema.UserQuotePieceInfo{}
	//	_ = copier.Copy(info, item)
	//	info.UrlTitle = htmltext.UrlTitle(info.Title)
	//	userQuotePiecelist = append(userQuotePiecelist, info)
	//}
	//
	//for _, item := range answerlist {
	//	info := &schema.UserAnswerInfo{}
	//	_ = copier.Copy(info, item)
	//	info.AnswerID = item.ID
	//	info.QuotePieceID = item.QuotePieceID
	//	info.QuotePieceInfo.UrlTitle = htmltext.UrlTitle(info.QuotePieceInfo.Title)
	//	userAnswerlist = append(userAnswerlist, info)
	//}
	//
	//return userQuotePiecelist, userAnswerlist, nil
}

// GetQuotePiecesByTitle get quotes by title
func (qs *QuotePieceService) GetQuotePiecesByTitle(ctx context.Context, title string) (
	resp []*schema.QuotePieceBaseInfo, err error) {
	resp = make([]*schema.QuotePieceBaseInfo, 0)
	if len(title) == 0 {
		return resp, nil
	}
	quotes, err := qs.QuotePieceRepo.GetQuotePiecesByTitle(ctx, title, 10)
	if err != nil {
		return resp, err
	}
	for _, quote := range quotes {
		item := &schema.QuotePieceBaseInfo{}
		item.ID = quote.ID
		item.Title = quote.Title
		item.UrlTitle = htmltext.UrlTitle(quote.Title)
		item.ViewCount = quote.ViewCount
		//item.AnswerCount = quote.AnswerCount
		item.CollectionCount = quote.CollectionCount
		item.FollowCount = quote.FollowCount
		status, ok := entity.AdminQuotePieceSearchStatusIntToString[quote.Status]
		if ok {
			item.Status = status
		}
		//if quote.AcceptedAnswerID != "0" {
		//	item.AcceptedAnswer = true
		//}
		resp = append(resp, item)
	}
	return resp, nil
}

// SimilarQuotePiece
func (qs *QuotePieceService) SimilarQuotePiece(ctx context.Context, quoteID string, loginUserID string) ([]*schema.QuotePiecePageResp, int64, error) {
	quote, err := qs.QuotePieceCommon.Info(ctx, quoteID, loginUserID)
	if err != nil {
		return nil, 0, nil
	}
	tagNames := make([]string, 0, len(quote.Tags))
	for _, tag := range quote.Tags {
		tagNames = append(tagNames, tag.SlugName)
	}
	search := &schema.QuotePiecePageReq{}
	search.OrderCond = "hot"
	search.Page = 0
	search.PageSize = 6
	if len(tagNames) > 0 {
		search.Tag = tagNames[0]
	}
	search.LoginUserID = loginUserID
	similarQuotePieces, _, err := qs.GetQuotePiecePage(ctx, search)
	if err != nil {
		return nil, 0, err
	}
	var result []*schema.QuotePiecePageResp
	for _, v := range similarQuotePieces {
		if uid.DeShortID(v.ID) != quoteID {
			result = append(result, v)
		}
	}
	return result, int64(len(result)), nil
}

// GetQuotePiecePage query quotes page
func (qs *QuotePieceService) GetQuotePiecePage(ctx context.Context, req *schema.QuotePiecePageReq) (
	quotes []*schema.QuotePiecePageResp, total int64, err error) {
	quotes = make([]*schema.QuotePiecePageResp, 0)
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
			return quotes, 0, nil
		}
		req.UserIDBeSearched = userinfo.ID
	}

	if req.OrderCond == schema.QuotePieceOrderCondHot {
		req.InDays = schema.HotInDays
	}

	quoteList, total, err := qs.QuotePieceRepo.GetQuotePiecePage(ctx, req.Page, req.PageSize,
		tagIDs, req.UserIDBeSearched, req.OrderCond, req.InDays, showHidden, req.ShowPending)
	if err != nil {
		return nil, 0, err
	}
	quotes, err = qs.QuotePieceCommon.FormatQuotePiecesPage(ctx, quoteList, req.LoginUserID, req.OrderCond)
	if err != nil {
		return nil, 0, err
	}
	return quotes, total, nil
}

// GetRecommendQuotePiecePage retrieves recommended quote page based on following tags and quotes.
func (qs *QuotePieceService) GetRecommendQuotePiecePage(ctx context.Context, req *schema.QuotePiecePageReq) (
	quotes []*schema.QuotePiecePageResp, total int64, err error) {
	followingTagsResp, err := qs.tagService.GetFollowingTags(ctx, req.LoginUserID)
	if err != nil {
		return nil, 0, err
	}
	tagIDs := make([]string, 0, len(followingTagsResp))
	for _, tag := range followingTagsResp {
		tagIDs = append(tagIDs, tag.TagID)
	}

	activityType, err := qs.activityRepo.GetActivityTypeByObjectType(ctx, constant.QuotePieceObjectType, "follow")
	if err != nil {
		return nil, 0, err
	}
	activities, err := qs.activityRepo.GetUserActivitysByActivityType(ctx, req.LoginUserID, activityType)
	if err != nil {
		return nil, 0, err
	}

	followedQuotePieceIDs := make([]string, 0, len(activities))
	for _, activity := range activities {
		if activity.Cancelled == entity.ActivityCancelled {
			continue
		}
		followedQuotePieceIDs = append(followedQuotePieceIDs, activity.ObjectID)
	}
	quoteList, total, err := qs.QuotePieceRepo.GetRecommendQuotePiecePageByTags(ctx, req.LoginUserID, tagIDs, followedQuotePieceIDs, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	quotes, err = qs.QuotePieceCommon.FormatQuotePiecesPage(ctx, quoteList, req.LoginUserID, "frequent")
	if err != nil {
		return nil, 0, err
	}

	return quotes, total, nil
}

func (qs *QuotePieceService) AdminSetQuotePieceStatus(ctx context.Context, req *schema.AdminUpdateQuotePieceStatusReq) error {
	setStatus, ok := entity.AdminQuotePieceSearchStatus[req.Status]
	if !ok {
		return errors.BadRequest(reason.RequestFormatError)
	}
	quoteInfo, exist, err := qs.QuotePieceRepo.GetQuotePiece(ctx, req.QuotePieceID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.QuotePieceNotFound)
	}
	err = qs.QuotePieceRepo.UpdateQuotePieceStatus(ctx, quoteInfo.ID, setStatus)
	if err != nil {
		return err
	}

	msg := &schema.NotificationMsg{}
	if setStatus == entity.QuotePieceStatusDeleted {
		// #2372 In order to simplify the process and complexity, as well as to consider if it is in-house,
		// facing the problem of recovery.
		//err = qs.answerActivityService.DeleteQuotePiece(ctx, quoteInfo.ID, quoteInfo.CreatedAt, quoteInfo.VoteCount)
		//if err != nil {
		//	log.Errorf("admin delete quote then rank rollback error %s", err.Error())
		//}
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           quoteInfo.UserID,
			TriggerUserID:    converter.StringToInt64(req.UserID),
			ObjectID:         quoteInfo.ID,
			OriginalObjectID: quoteInfo.ID,
			ActivityTypeKey:  constant.ActQuotePieceDeleted,
		})
		msg.NotificationAction = constant.NotificationYourQuotePieceIsClosed
	}
	if setStatus == entity.QuotePieceStatusAvailable && quoteInfo.Status == entity.QuotePieceStatusClosed {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           quoteInfo.UserID,
			TriggerUserID:    converter.StringToInt64(req.UserID),
			ObjectID:         quoteInfo.ID,
			OriginalObjectID: quoteInfo.ID,
			ActivityTypeKey:  constant.ActQuotePieceReopened,
		})
	}
	if setStatus == entity.QuotePieceStatusClosed && quoteInfo.Status != entity.QuotePieceStatusClosed {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           quoteInfo.UserID,
			TriggerUserID:    converter.StringToInt64(req.UserID),
			ObjectID:         quoteInfo.ID,
			OriginalObjectID: quoteInfo.ID,
			ActivityTypeKey:  constant.ActQuotePieceClosed,
		})
		msg.NotificationAction = constant.NotificationYourQuotePieceIsClosed
	}
	// recover
	if setStatus == entity.QuotePieceStatusAvailable && quoteInfo.Status == entity.QuotePieceStatusDeleted {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           req.UserID,
			TriggerUserID:    converter.StringToInt64(req.UserID),
			ObjectID:         quoteInfo.ID,
			OriginalObjectID: quoteInfo.ID,
			ActivityTypeKey:  constant.ActQuotePieceUndeleted,
		})
	}

	if len(msg.NotificationAction) > 0 {
		msg.ObjectID = quoteInfo.ID
		msg.Type = schema.NotificationTypeInbox
		msg.ReceiverUserID = quoteInfo.UserID
		msg.TriggerUserID = req.UserID
		msg.ObjectType = constant.QuotePieceObjectType
		qs.notificationQueueService.Send(ctx, msg)
	}
	return nil
}

func (qs *QuotePieceService) AdminQuotePiecePage(
	ctx context.Context, req *schema.AdminQuotePiecePageReq) (
	resp *pager.PageModel, err error) {

	list := make([]*schema.AdminQuotePieceInfo, 0)
	quoteList, count, err := qs.QuotePieceRepo.AdminQuotePiecePage(ctx, req)
	if err != nil {
		return nil, err
	}

	userIds := make([]string, 0)
	for _, info := range quoteList {
		item := &schema.AdminQuotePieceInfo{}
		_ = copier.Copy(item, info)
		item.CreateTime = info.CreatedAt.Unix()
		//item.UpdateTime = info.PostUpdateTime.Unix()
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
func (qs *QuotePieceService) AdminAnswerPage(ctx context.Context, req *schema.AdminAnswerPageReq) (
	resp *pager.PageModel, err error) {
	return
	//
	//answerList, count, err := qs.QuotePieceCommon.AnswerCommon.AdminSearchList(ctx, req)
	//if err != nil {
	//	return nil, err
	//}
	//
	//quoteIDs := make([]string, 0)
	//userIds := make([]string, 0)
	//answerResp := make([]*schema.AdminAnswerInfo, 0)
	//for _, item := range answerList {
	//	answerInfo := qs.QuotePieceCommon.AnswerCommon.AdminShowFormat(ctx, item)
	//	answerResp = append(answerResp, answerInfo)
	//	quoteIDs = append(quoteIDs, item.QuotePieceID)
	//	userIds = append(userIds, item.UserID)
	//}
	//userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIds)
	//if err != nil {
	//	return nil, err
	//}
	//quoteMaps, err := qs.QuotePieceCommon.FindInfoByID(ctx, quoteIDs, req.LoginUserID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, item := range answerResp {
	//	if q, ok := quoteMaps[item.QuotePieceID]; ok {
	//		item.QuotePieceInfo.Title = q.Title
	//	}
	//	if u, ok := userInfoMap[item.UserID]; ok {
	//		item.UserInfo = u
	//	}
	//}
	//return pager.NewPageModel(count, answerResp), nil
}

func (qs *QuotePieceService) changeQuotePieceToRevision(ctx context.Context, quoteInfo *entity.QuotePiece, tags []*entity.Tag) (
	quoteRevision *entity.QuotePieceWithTagsRevision, err error) {
	quoteRevision = &entity.QuotePieceWithTagsRevision{}
	quoteRevision.QuotePiece = *quoteInfo

	for _, tag := range tags {
		item := &entity.TagSimpleInfoForRevision{}
		_ = copier.Copy(item, tag)
		quoteRevision.Tags = append(quoteRevision.Tags, item)
	}
	return quoteRevision, nil
}

func (qs *QuotePieceService) SitemapCron(ctx context.Context) {
	siteSeo, err := qs.siteInfoService.GetSiteSeo(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	ctx = context.WithValue(ctx, constant.ShortIDFlag, siteSeo.IsShortLink())
	qs.QuotePieceCommon.SitemapCron(ctx)
}

// 等于 根据authorname精确匹配
func (qs *QuotePieceService) GetQuotePieceByTitle(ctx context.Context, title string) (
	resp *schema.QuotePieceBaseInfo, err error) {
	//resp = make([]*schema.QuotePieceBaseInfo, 0)

	if len(title) == 0 {
		return resp, nil
	}
	quote, err := qs.QuotePieceRepo.GetQuotePieceByTitle(ctx, title)
	if err != nil {
		return resp, err
	}
	if quote == nil {
		return resp, nil
	}
	//for _, quote := range quotes {
	item := &schema.QuotePieceBaseInfo{}
	item.ID = quote.ID
	item.Title = quote.Title
	item.UrlTitle = htmltext.UrlTitle(quote.Title)
	item.ViewCount = quote.ViewCount
	//item.AnswerCount = quote.AnswerCount
	item.CollectionCount = quote.CollectionCount
	item.FollowCount = quote.FollowCount
	status, ok := entity.AdminQuotePieceSearchStatusIntToString[quote.Status]
	if ok {
		item.Status = status
	}
	//if quote.AcceptedAnswerID != "0" {
	//	item.AcceptedAnswer = true
	//}
	//	resp = append(resp, item)
	//}
	resp = item
	return resp, nil
}
