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
	quotecommon "github.com/apache/incubator-answer/internal/service_quote/quote_common"
	"github.com/apache/incubator-answer/pkg/checker"
	"github.com/apache/incubator-answer/pkg/converter"
	"github.com/apache/incubator-answer/pkg/htmltext"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
	"golang.org/x/net/context"
)

// QuoteRepo quote repository

// QuoteService user service
type QuoteService struct {
	activityRepo                     activity_common.ActivityRepo
	quoteRepo                        quotecommon.QuoteRepo
	answerRepo                       answercommon.AnswerRepo
	tagCommon                        *tagcommon.TagCommonService
	tagService                       *tag.TagService
	quotecommon                      *quotecommon.QuoteCommon
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
	newQuoteNotificationService      *notification.ExternalNotificationService
	reviewService                    *review.ReviewService
	configService                    *config.ConfigService
	eventQueueService                event_queue.EventQueueService

	quoteAuthorService *QuoteAuthorService
	quotePieceService  *QuotePieceService
	quoteAuthorRepo    quotecommon.QuoteAuthorRepo
	quotePieceRepo     quotecommon.QuotePieceRepo
	quoteAuthorCommon  *quotecommon.QuoteAuthorCommon
	quotePieceCommon   *quotecommon.QuotePieceCommon
}

func NewQuoteService(
	activityRepo activity_common.ActivityRepo,
	quoteRepo quotecommon.QuoteRepo,
	answerRepo answercommon.AnswerRepo,
	tagCommon *tagcommon.TagCommonService,
	tagService *tag.TagService,
	quotecommon *quotecommon.QuoteCommon,
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
	newQuoteNotificationService *notification.ExternalNotificationService,
	reviewService *review.ReviewService,
	configService *config.ConfigService,
	eventQueueService event_queue.EventQueueService,

	quoteAuthorService *QuoteAuthorService, //@
	quotePieceService *QuotePieceService,
	quoteAuthorRepo quotecommon.QuoteAuthorRepo,
	quotePieceRepo quotecommon.QuotePieceRepo,

	quoteAuthorCommon *quotecommon.QuoteAuthorCommon,
	quotePieceCommon *quotecommon.QuotePieceCommon,
) *QuoteService {
	return &QuoteService{
		activityRepo:                     activityRepo,
		quoteRepo:                        quoteRepo,
		answerRepo:                       answerRepo,
		tagCommon:                        tagCommon,
		tagService:                       tagService,
		quotecommon:                      quotecommon,
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
		newQuoteNotificationService:      newQuoteNotificationService,
		reviewService:                    reviewService,
		configService:                    configService,
		eventQueueService:                eventQueueService,

		quoteAuthorService: quoteAuthorService,
		quotePieceService:  quotePieceService,

		quoteAuthorRepo: quoteAuthorRepo,
		quotePieceRepo:  quotePieceRepo,

		quoteAuthorCommon: quoteAuthorCommon,
		quotePieceCommon:  quotePieceCommon,
	}
}

func (qs *QuoteService) CloseQuote(ctx context.Context, req *schema.CloseQuoteReq) error {
	quoteInfo, has, err := qs.quoteRepo.GetQuote(ctx, req.ID)
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

	quoteInfo.Status = entity.QuoteStatusClosed
	err = qs.quoteRepo.UpdateQuoteStatus(ctx, quoteInfo.ID, quoteInfo.Status)
	if err != nil {
		return err
	}

	closeMeta, _ := json.Marshal(schema.CloseQuoteMeta{
		CloseType: req.CloseType,
		CloseMsg:  req.CloseMsg,
	})
	err = qs.metaService.AddMeta(ctx, req.ID, entity.QuoteCloseReasonKey, string(closeMeta))
	if err != nil {
		return err
	}

	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           req.UserID,
		ObjectID:         quoteInfo.ID,
		OriginalObjectID: quoteInfo.ID,
		ActivityTypeKey:  constant.ActQuoteClosed,
	})
	return nil
}

// ReopenQuote reopen quote
func (qs *QuoteService) ReopenQuote(ctx context.Context, req *schema.ReopenQuoteReq) error {
	quoteInfo, has, err := qs.quoteRepo.GetQuote(ctx, req.QuoteID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	quoteInfo.Status = entity.QuoteStatusAvailable
	err = qs.quoteRepo.UpdateQuoteStatus(ctx, quoteInfo.ID, quoteInfo.Status)
	if err != nil {
		return err
	}
	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           req.UserID,
		ObjectID:         quoteInfo.ID,
		OriginalObjectID: quoteInfo.ID,
		ActivityTypeKey:  constant.ActQuoteReopened,
	})
	return nil
}

func (qs *QuoteService) AddQuoteCheckTags(ctx context.Context, Tags []*entity.Tag) ([]string, error) {
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
func (qs *QuoteService) CheckAddQuote(ctx context.Context, req *schema.QuoteAdd) (errorlist any, err error) {
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
	if !req.QuotePermission.CanUseReservedTag {
		taglist, err := qs.AddQuoteCheckTags(ctx, Tags)
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
func (qs *QuoteService) HasNewTag(ctx context.Context, tags []*schema.TagItem) (bool, error) {
	return qs.tagCommon.HasNewTag(ctx, tags)
}

// 参考 func (qc *QuoteAuthorController) AddQuoteAuthor(ctx *gin.Context) {
func (qs *QuoteService) AddQuoteAuthor(ctx context.Context, req *schema.QuoteAdd) (quoteInfo any, err error) {

	//req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	//canList, requireRanks, err := qc.rankService.CheckOperationPermissionsForRanks(ctx, req.UserID, []string{
	//	permission.QuoteAuthorAdd,
	//	permission.QuoteAuthorEdit,
	//	permission.QuoteAuthorDelete,
	//	permission.QuoteAuthorClose,
	//	permission.QuoteAuthorReopen,
	//	permission.TagUseReservedTag,
	//	permission.TagAdd,
	//	permission.LinkUrlLimit,
	//})
	//if err != nil {
	//	handler.HandleResponse(ctx, err, nil)
	//	return
	//}
	//linkUrlLimitUser := canList[7]
	//isAdmin := middleware.GetUserIsAdminModerator(ctx)
	//if !isAdmin || !linkUrlLimitUser {
	//	captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionQuoteAuthor, req.UserID, req.CaptchaID, req.CaptchaCode)
	//	if !captchaPass {
	//		errFields := append([]*validator.FormErrorField{}, &validator.FormErrorField{
	//			ErrorField: "captcha_code",
	//			ErrorMsg:   translator.Tr(handler.GetLang(ctx), reason.CaptchaVerificationFailed),
	//		})
	//		handler.HandleResponse(ctx, errors.BadRequest(reason.CaptchaVerificationFailed), errFields)
	//		return
	//	}
	//}

	//req.CanAdd = canList[0]
	//req.CanEdit = canList[1]
	//req.CanDelete = canList[2]
	//req.CanClose = canList[3]
	//req.CanReopen = canList[4]
	//req.CanUseReservedTag = canList[5]
	//req.CanAddTag = canList[6]
	//if !req.CanAdd {
	//	handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
	//	return
	//}

	//// can add tag
	//hasNewTag, err := qc.QuoteAuthorService.HasNewTag(ctx, req.Tags)
	//if err != nil {
	//	handler.HandleResponse(ctx, err, nil)
	//	return
	//}
	//if !req.CanAddTag && hasNewTag {
	//	lang := handler.GetLang(ctx)
	//	msg := translator.TrWithData(lang, reason.NoEnoughRankToOperate, &schema.PermissionTrTplData{Rank: requireRanks[6]})
	//	handler.HandleResponse(ctx, errors.Forbidden(reason.NoEnoughRankToOperate).WithMsg(msg), nil)
	//	return
	//}

	//errList, err := qc.QuoteAuthorService.CheckAddQuoteAuthor(ctx, req)
	//if err != nil {
	//	errlist, ok := errList.([]*validator.FormErrorField)
	//	if ok {
	//		errFields = append(errFields, errlist...)
	//	}
	//}
	//
	//if len(errFields) > 0 {
	//	handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
	//	return
	//}
	//
	//req.UserAgent = ctx.GetHeader("User-Agent")
	//req.IP = ctx.ClientIP()

	//resp, err := qc.QuoteAuthorService.AddQuoteAuthor(ctx, req)
	req_QuoteAuthorAdd := &schema.QuoteAuthorAdd{
		UserID:     req.UserID,
		AuthorName: req.Author,
		//Content: req.Content,
		//ContentFormat: req.ContentFormat,
		//Title: req.Title,
		//Tags: req.Tags,
		Tags: make([]*schema.TagItem, 0),
	}
	var errFields []*validator.FormErrorField
	errFields = make([]*validator.FormErrorField, 0)
	resp, err := qs.quoteAuthorService.AddQuoteAuthor(ctx, req_QuoteAuthorAdd)
	if err != nil {

		errlist, ok := resp.([]*validator.FormErrorField)
		if ok {
			errFields = append(errFields, errlist...)
		}

	}

	if len(errFields) > 0 {
		//handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		//return
		err := errors.BadRequest(reason.RequestFormatError)
		return errFields, err
	}

	quoteInfo = resp
	log.Infof("resp:%+v", resp)

	return
	//if !isAdmin || !linkUrlLimitUser {
	//	qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionQuoteAuthor, req.UserID)
	//}
	//handler.HandleResponse(ctx, err, resp)
}

// 参考 func (qc *QuoteAuthorController) AddQuoteAuthor(ctx *gin.Context) {
func (qs *QuoteService) AddQuotePiece(ctx context.Context, req *schema.QuoteAdd) (quoteInfo any, err error) {

	//req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	//canList, requireRanks, err := qc.rankService.CheckOperationPermissionsForRanks(ctx, req.UserID, []string{
	//	permission.QuoteAuthorAdd,
	//	permission.QuoteAuthorEdit,
	//	permission.QuoteAuthorDelete,
	//	permission.QuoteAuthorClose,
	//	permission.QuoteAuthorReopen,
	//	permission.TagUseReservedTag,
	//	permission.TagAdd,
	//	permission.LinkUrlLimit,
	//})
	//if err != nil {
	//	handler.HandleResponse(ctx, err, nil)
	//	return
	//}
	//linkUrlLimitUser := canList[7]
	//isAdmin := middleware.GetUserIsAdminModerator(ctx)
	//if !isAdmin || !linkUrlLimitUser {
	//	captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionQuoteAuthor, req.UserID, req.CaptchaID, req.CaptchaCode)
	//	if !captchaPass {
	//		errFields := append([]*validator.FormErrorField{}, &validator.FormErrorField{
	//			ErrorField: "captcha_code",
	//			ErrorMsg:   translator.Tr(handler.GetLang(ctx), reason.CaptchaVerificationFailed),
	//		})
	//		handler.HandleResponse(ctx, errors.BadRequest(reason.CaptchaVerificationFailed), errFields)
	//		return
	//	}
	//}

	//req.CanAdd = canList[0]
	//req.CanEdit = canList[1]
	//req.CanDelete = canList[2]
	//req.CanClose = canList[3]
	//req.CanReopen = canList[4]
	//req.CanUseReservedTag = canList[5]
	//req.CanAddTag = canList[6]
	//if !req.CanAdd {
	//	handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
	//	return
	//}

	//// can add tag
	//hasNewTag, err := qc.QuoteAuthorService.HasNewTag(ctx, req.Tags)
	//if err != nil {
	//	handler.HandleResponse(ctx, err, nil)
	//	return
	//}
	//if !req.CanAddTag && hasNewTag {
	//	lang := handler.GetLang(ctx)
	//	msg := translator.TrWithData(lang, reason.NoEnoughRankToOperate, &schema.PermissionTrTplData{Rank: requireRanks[6]})
	//	handler.HandleResponse(ctx, errors.Forbidden(reason.NoEnoughRankToOperate).WithMsg(msg), nil)
	//	return
	//}

	//errList, err := qc.QuoteAuthorService.CheckAddQuoteAuthor(ctx, req)
	//if err != nil {
	//	errlist, ok := errList.([]*validator.FormErrorField)
	//	if ok {
	//		errFields = append(errFields, errlist...)
	//	}
	//}
	//
	//if len(errFields) > 0 {
	//	handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
	//	return
	//}
	//
	//req.UserAgent = ctx.GetHeader("User-Agent")
	//req.IP = ctx.ClientIP()

	//resp, err := qc.QuoteAuthorService.AddQuoteAuthor(ctx, req)
	req_QuotePieceAdd := &schema.QuotePieceAdd{
		UserID: req.UserID,
		//AuthorName: req.Author,
		//Content: req.Content,
		//ContentFormat: req.ContentFormat,
		Title: req.PieceName,
		//Tags: req.Tags,
		Tags: make([]*schema.TagItem, 0),
	}
	var errFields []*validator.FormErrorField
	errFields = make([]*validator.FormErrorField, 0)
	resp, err := qs.quotePieceService.AddQuotePiece(ctx, req_QuotePieceAdd)
	if err != nil {

		errlist, ok := resp.([]*validator.FormErrorField)
		if ok {
			errFields = append(errFields, errlist...)
		}

	}

	if len(errFields) > 0 {
		//handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		//return
		err := errors.BadRequest(reason.RequestFormatError)
		return errFields, err
	}

	quoteInfo = resp
	log.Infof("resp:%+v", resp)

	return
	//if !isAdmin || !linkUrlLimitUser {
	//	qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionQuoteAuthor, req.UserID)
	//}
	//handler.HandleResponse(ctx, err, resp)
}

// AddQuote add quote
func (qs *QuoteService) AddQuote(ctx context.Context, req *schema.QuoteAdd) (quoteInfo any, err error) {
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
		return quoteInfo, tagerr
	}
	if !req.QuotePermission.CanUseReservedTag {
		taglist, err := qs.AddQuoteCheckTags(ctx, tags)
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

	quote := &entity.Quote{}
	now := time.Now()
	quote.UserID = req.UserID
	quote.Title = req.Title
	quote.OriginalText = req.Content
	log.Infof("addQuote content:%s", req.Content)
	quote.ParsedText = req.HTML
	//quote.AcceptedAnswerID = "0"
	//quote.LastAnswerID = "0"
	//quote.LastEditUserID = "0"
	//quote.PostUpdateTime = nil
	quote.Status = entity.QuoteStatusPending
	quote.RevisionID = "0"
	quote.CreatedAt = now
	quote.PostUpdateTime = now
	quote.Pin = entity.QuoteUnPin
	quote.Show = entity.QuoteShow
	quote.UpdatedAt = now

	req.AuthorId = strings.TrimSpace(req.AuthorId)
	req.Author = strings.TrimSpace(req.Author)
	//@cws
	//quote.PostDate = now
	//quote.OriginalTextFormat = req.ContentFormat
	if req.AuthorId != "" {
		quote.QuoteAuthorId = uid.DeShortID(req.AuthorId)
		log.Infof("req.AuthorId  has val:%s", uid.DeShortID(req.AuthorId))
	} else { //根据输入文字新建作者
		//先根据authorName搜索ID ，如果没有再新建

		log.Infof(" req.Author val:%s", req.Author)
		if req.Author == "" {
			req.Author = "佚名"
		}

		quoteAuthorBaseInfo, err := qs.quoteAuthorService.GetQuoteAuthorByAuthorName(ctx, req.Author)
		if err != nil {
			err = errors.BadRequest(reason.UnknownError)
			return nil, err
		}
		if quoteAuthorBaseInfo != nil {
			log.Infof("quoteAuthorBaseInfo exit:%+v", quoteAuthorBaseInfo)
			quote.QuoteAuthorId = uid.DeShortID(quoteAuthorBaseInfo.ID)
		} else {
			//TODO author可以再增加个同义词功能，
			log.Infof("add quote author :%+v", req.Author)
			quoteAuthorAny, err := qs.AddQuoteAuthor(ctx, req)
			if err != nil {
				//errlist, _ := quoteAuthorAny.([]*validator.FormErrorField)
				return quoteAuthorAny, err
			}
			quoteAuthorInfoResp, ok := quoteAuthorAny.(*schema.QuoteAuthorInfoResp)
			if !ok {
				log.Errorf("quoteAuthorInfoResp error")
				return quoteAuthorAny, errors.BadRequest(reason.RequestFormatError)
			}
			log.Infof("quoteAuthorInfoResp.ID:%+v", quoteAuthorInfoResp.ID)
			//resp里面的ID是encode的，要decode 		info.ID = uid.EnShortID(data.ID)
			quote.QuoteAuthorId = uid.DeShortID(quoteAuthorInfoResp.ID)
		}

		log.Infof("quote.QuoteAuthorId:%s", quote.QuoteAuthorId)
		log.Infof("req.AuthorId  no val:%s", req.AuthorId)

	}

	req.PieceId = strings.TrimSpace(req.PieceId)
	req.PieceName = strings.TrimSpace(req.PieceName)

	if req.PieceId != "" {
		quote.QuotePieceId = uid.DeShortID(req.PieceId)
		log.Infof("req.PieceId  has val:%s", uid.DeShortID(req.PieceId))
	} else { //根据输入文字新建作者
		//先根据authorName搜索ID ，如果没有再新建

		log.Infof(" req.PieceName val:%s", req.PieceName)
		if req.PieceName == "" {
			req.PieceName = "未名"
		}

		quotePieceBaseInfo, err := qs.quotePieceService.GetQuotePieceByTitle(ctx, req.PieceName)
		if err != nil {
			err = errors.BadRequest(reason.UnknownError)
			return nil, err
		}
		if quotePieceBaseInfo != nil {
			log.Infof("quoteAuthorBaseInfo exit:%+v", quotePieceBaseInfo)
			quote.QuotePieceId = uid.DeShortID(quotePieceBaseInfo.ID)
		} else {
			log.Infof("add quote piece :%+v", req.Title)
			quotePieceAny, err := qs.AddQuotePiece(ctx, req)
			if err != nil {
				//errlist, _ := quoteAuthorAny.([]*validator.FormErrorField)
				return quotePieceAny, err
			}
			quotePieceInfoResp, ok := quotePieceAny.(*schema.QuotePieceInfoResp)
			if !ok {
				log.Errorf("quoteAuthorInfoResp error")
				return quotePieceAny, errors.BadRequest(reason.RequestFormatError)
			}
			log.Infof("quoteAuthorInfoResp.ID:%+v", quotePieceInfoResp.ID)
			//resp里面的ID是encode的，要decode 		info.ID = uid.EnShortID(data.ID)
			quote.QuotePieceId = uid.DeShortID(quotePieceInfoResp.ID)
		}

		log.Infof("quote.QuotePieceId:%s", quote.QuotePieceId)
		log.Infof("req.PieceId  no val:%s", req.PieceId)

	}

	err = qs.quoteRepo.AddQuote(ctx, quote)
	if err != nil {
		return
	}
	quote.Status = qs.reviewService.AddQuoteReview(ctx, quote, req.Tags, req.IP, req.UserAgent)
	if err := qs.quoteRepo.UpdateQuoteStatus(ctx, quote.ID, quote.Status); err != nil {
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
	_ = qs.quoteRepo.UpdateSearch(ctx, quote.ID)

	revisionDTO := &schema.AddRevisionDTO{
		UserID:   quote.UserID,
		ObjectID: quote.ID,
		Title:    quote.Title,
	}

	quoteWithTagsRevision, err := qs.changeQuoteToRevision(ctx, quote, tags)
	if err != nil {
		return nil, err
	}
	infoJSON, _ := json.Marshal(quoteWithTagsRevision)
	revisionDTO.Content = string(infoJSON)
	revisionID, err := qs.revisionService.AddRevision(ctx, revisionDTO, true)
	if err != nil {
		return
	}

	// user add quote count
	userQuoteCount, err := qs.quotecommon.GetUserQuoteCount(ctx, quote.UserID)
	if err != nil {
		log.Errorf("get user quote count error %v", err)
	} else {
		err = qs.userCommon.UpdateQuoteCount(ctx, quote.UserID, userQuoteCount)
		if err != nil {
			log.Errorf("update user quote count error %v", err)
		}
	}

	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           quote.UserID,
		ObjectID:         quote.ID,
		OriginalObjectID: quote.ID,
		ActivityTypeKey:  constant.ActQuoteAsked,
		RevisionID:       revisionID,
	})

	if quote.Status == entity.QuoteStatusAvailable {
		qs.externalNotificationQueueService.Send(ctx,
			schema.CreateNewQuoteNotificationMsg(quote.ID, quote.Title, quote.UserID, tags))
	}
	qs.eventQueueService.Send(ctx, schema.NewEvent(constant.EventQuoteCreate, req.UserID).TID(quote.ID).
		QID(quote.ID, quote.UserID))

	quoteInfo, err = qs.GetQuote(ctx, quote.ID, quote.UserID, req.QuotePermission)
	return
}

// OperationQuote
func (qs *QuoteService) OperationQuote(ctx context.Context, req *schema.OperationQuoteReq) (err error) {
	quoteInfo, has, err := qs.quoteRepo.GetQuote(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	// Hidden quote cannot be placed at the top
	if quoteInfo.Show == entity.QuoteHide && req.Operation == schema.QuoteOperationPin {
		return nil
	}
	// Quote cannot be hidden when they are at the top
	if quoteInfo.Pin == entity.QuotePin && req.Operation == schema.QuoteOperationHide {
		return nil
	}

	switch req.Operation {
	case schema.QuoteOperationHide:
		quoteInfo.Show = entity.QuoteHide
		err = qs.tagCommon.HideTagRelListByObjectID(ctx, req.ID)
		if err != nil {
			return err
		}
		err = qs.tagCommon.RefreshTagCountByArticleID(ctx, req.ID)
		if err != nil {
			return err
		}
	case schema.QuoteOperationShow:
		quoteInfo.Show = entity.QuoteShow
		err = qs.tagCommon.ShowTagRelListByObjectID(ctx, req.ID)
		if err != nil {
			return err
		}
		err = qs.tagCommon.RefreshTagCountByQuoteID(ctx, req.ID)
		if err != nil {
			return err
		}
	case schema.QuoteOperationPin:
		quoteInfo.Pin = entity.QuotePin
	case schema.QuoteOperationUnPin:
		quoteInfo.Pin = entity.QuoteUnPin
	}

	err = qs.quoteRepo.UpdateQuoteOperation(ctx, quoteInfo)
	if err != nil {
		return err
	}

	actMap := make(map[string]constant.ActivityTypeKey)
	actMap[schema.QuoteOperationPin] = constant.ActQuotePin
	actMap[schema.QuoteOperationUnPin] = constant.ActQuoteUnPin
	actMap[schema.QuoteOperationHide] = constant.ActQuoteHide
	actMap[schema.QuoteOperationShow] = constant.ActQuoteShow
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

// RemoveQuote delete quote
func (qs *QuoteService) RemoveQuote(ctx context.Context, req *schema.RemoveQuoteReq) (err error) {
	quoteInfo, has, err := qs.quoteRepo.GetQuote(ctx, req.ID)
	if err != nil {
		return err
	}
	//if the status is deleted, return directly
	if quoteInfo.Status == entity.QuoteStatusDeleted {
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
		//	return errors.BadRequest(reason.QuoteCannotDeleted)
		//}
		//if quoteInfo.AnswerCount > 1 {
		//	return errors.BadRequest(reason.QuoteCannotDeleted)
		//}

		//if quoteInfo.AnswerCount == 1 {
		//	answersearch := &entity.AnswerSearch{}
		//	answersearch.QuoteID = req.ID
		//	answerList, _, err := qs.quotecommon.AnswerCommon.Search(ctx, answersearch)
		//	if err != nil {
		//		return err
		//	}
		//	for _, answer := range answerList {
		//		if answer.VoteCount > 0 {
		//			return errors.BadRequest(reason.QuoteCannotDeleted)
		//		}
		//	}
		//}
	}

	quoteInfo.Status = entity.QuoteStatusDeleted
	err = qs.quoteRepo.UpdateQuoteStatusWithOutUpdateTime(ctx, quoteInfo)
	if err != nil {
		return err
	}

	userQuoteCount, err := qs.quotecommon.GetUserQuoteCount(ctx, quoteInfo.UserID)
	if err != nil {
		log.Error("user GetUserQuoteCount error", err.Error())
	} else {
		err = qs.userCommon.UpdateQuoteCount(ctx, quoteInfo.UserID, userQuoteCount)
		if err != nil {
			log.Error("user IncreaseQuoteCount error", err.Error())
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
	err = qs.tagCommon.RefreshTagQuoteCount(ctx, tagIDs)
	if err != nil {
		log.Error("efreshTagQuoteCount error", err.Error())
	}

	// #2372 In order to simplify the process and complexity, as well as to consider if it is in-house,
	// facing the problem of recovery.
	// err = qs.answerActivityService.DeleteQuote(ctx, quoteInfo.ID, quoteInfo.CreatedAt, quoteInfo.VoteCount)
	// if err != nil {
	// 	 log.Errorf("user DeleteQuote rank rollback error %s", err.Error())
	// }
	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           quoteInfo.UserID,
		TriggerUserID:    converter.StringToInt64(req.UserID),
		ObjectID:         quoteInfo.ID,
		OriginalObjectID: quoteInfo.ID,
		ActivityTypeKey:  constant.ActQuoteDeleted,
	})
	qs.eventQueueService.Send(ctx, schema.NewEvent(constant.EventQuoteDelete, req.UserID).TID(quoteInfo.ID).
		QID(quoteInfo.ID, quoteInfo.UserID))
	return nil
}

func (qs *QuoteService) UpdateQuoteCheckTags(ctx context.Context, req *schema.QuoteUpdate) (errorlist []*validator.FormErrorField, err error) {
	dbinfo, has, err := qs.quoteRepo.GetQuote(ctx, req.ID)
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

func (qs *QuoteService) RecoverQuote(ctx context.Context, req *schema.QuoteRecoverReq) (err error) {
	quoteInfo, exist, err := qs.quoteRepo.GetQuote(ctx, req.QuoteID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.QuoteNotFound)
	}
	if quoteInfo.Status != entity.QuoteStatusDeleted {
		return nil
	}

	err = qs.quoteRepo.RecoverQuote(ctx, req.QuoteID)
	if err != nil {
		return err
	}

	// update user's quote count
	userQuoteCount, err := qs.quotecommon.GetUserQuoteCount(ctx, quoteInfo.UserID)
	if err != nil {
		log.Error("user GetUserQuoteCount error", err.Error())
	} else {
		err = qs.userCommon.UpdateQuoteCount(ctx, quoteInfo.UserID, userQuoteCount)
		if err != nil {
			log.Error("user IncreaseQuoteCount error", err.Error())
		}
	}

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
		ActivityTypeKey:  constant.ActQuoteUndeleted,
	})
	return nil
}

func (qs *QuoteService) UpdateQuoteInviteUser(ctx context.Context, req *schema.QuoteUpdateInviteUser) (err error) {
	return nil
	//originQuote, exist, err := qs.quoteRepo.GetQuote(ctx, req.ID)
	//if err != nil {
	//	return err
	//}
	//if !exist {
	//	return errors.BadRequest(reason.QuoteNotFound)
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
	//quote := &entity.Quote{}
	//quote.ID = uid.DeShortID(req.ID)
	////quote.InviteUserID = inviteUserStr
	//_ = inviteUserStr
	//
	//saveerr := qs.quoteRepo.UpdateQuote(ctx, quote, []string{"invite_user_id"})
	//if saveerr != nil {
	//	return saveerr
	//}
	////send notification
	//oldInviteUserIDsStr := originQuote.InviteUserID
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
	//go qs.notificationInviteUser(ctx, needSendNotificationUserIDs, originQuote.ID, originQuote.Title, req.UserID)
	//
	//return nil
}

func (qs *QuoteService) notificationInviteUser(
	ctx context.Context, invitedUserIDs []string, quoteID, quoteTitle, quoteUserID string) {
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
	//	msg.ObjectType = constant.QuoteObjectType
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
	//		QuoteTitle:       quoteTitle,
	//		QuoteID:          quoteID,
	//		UnsubscribeCode:    token.GenerateToken(),
	//	}
	//	externalNotificationMsg.NewInviteAnswerTemplateRawData = rawData
	//	qs.externalNotificationQueueService.Send(ctx, externalNotificationMsg)
	//}
}

// UpdateQuote update quote
func (qs *QuoteService) UpdateQuote(ctx context.Context, req *schema.QuoteUpdate) (quoteInfo any, err error) {
	var canUpdate bool
	quoteInfo = &schema.QuoteInfoResp{}

	_, existUnreviewed, err := qs.revisionService.ExistUnreviewedByObjectID(ctx, req.ID)
	if err != nil {
		return
	}
	if existUnreviewed {
		err = errors.BadRequest(reason.QuoteCannotUpdate)
		return
	}

	dbinfo, has, err := qs.quoteRepo.GetQuote(ctx, req.ID)
	if err != nil {
		return
	}
	if !has {
		return
	}
	if dbinfo.Status == entity.QuoteStatusDeleted {
		err = errors.BadRequest(reason.QuoteCannotUpdate)
		return nil, err
	}
	log.Infof("UpdateQuote b:%s", req.Content)
	log.Infof("UpdateQuote b html:%s", req.HTML)
	now := time.Now()
	quote := &entity.Quote{}
	quote.Title = req.Title
	quote.OriginalText = req.Content
	quote.ParsedText = req.HTML
	quote.ID = uid.DeShortID(req.ID)
	quote.UpdatedAt = now
	quote.PostUpdateTime = now
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
		saveerr := qs.quoteRepo.UpdateQuote(ctx, quote, []string{"title", "original_text", "parsed_text", "updated_at", "post_update_time", "last_edit_user_id"})
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

	quoteWithTagsRevision, err := qs.changeQuoteToRevision(ctx, quote, Tags)
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
			ActivityTypeKey:  constant.ActQuoteEdited,
			RevisionID:       revisionID,
			OriginalObjectID: quote.ID,
		})
		qs.eventQueueService.Send(ctx, schema.NewEvent(constant.EventQuoteUpdate, req.UserID).TID(quote.ID).
			QID(quote.ID, quote.UserID))
	}

	quoteInfo, err = qs.GetQuote(ctx, quote.ID, quote.UserID, req.QuotePermission)
	return
}

// GetQuote get quote one
func (qs *QuoteService) GetQuote(ctx context.Context, quoteID, userID string,
	per schema.QuotePermission) (resp *schema.QuoteInfoResp, err error) {
	quote, err := qs.quotecommon.Info(ctx, quoteID, userID)
	if err != nil {
		return
	}
	// If the quote is deleted or pending, only the administrator and the author can view it
	if (quote.Status == entity.QuoteStatusDeleted ||
		quote.Status == entity.QuoteStatusPending) && !per.CanReopen && quote.UserID != userID {
		return nil, errors.NotFound(reason.QuoteNotFound)
	}
	if quote.Status != entity.QuoteStatusClosed {
		per.CanReopen = false
	}
	if quote.Status == entity.QuoteStatusClosed {
		per.CanClose = false
	}
	if quote.Pin == entity.QuotePin {
		per.CanPin = false
		per.CanHide = false
	}
	if quote.Pin == entity.QuoteUnPin {
		per.CanUnPin = false
	}
	if quote.Show == entity.QuoteShow {
		per.CanShow = false
	}
	if quote.Show == entity.QuoteHide {
		per.CanHide = false
		per.CanPin = false
	}

	if quote.Status == entity.QuoteStatusDeleted {
		operation := &schema.Operation{}
		operation.Msg = translator.Tr(handler.GetLangByCtx(ctx), reason.QuoteAlreadyDeleted)
		operation.Level = schema.OperationLevelDanger
		quote.Operation = operation
	}
	if quote.Status == entity.QuoteStatusPending {
		operation := &schema.Operation{}
		operation.Msg = translator.Tr(handler.GetLangByCtx(ctx), reason.QuoteUnderReview)
		operation.Level = schema.OperationLevelSecondary
		quote.Operation = operation
	}

	quote.Description = htmltext.FetchExcerpt(quote.HTML, "...", 120)
	//@cws
	if quote.Title == "" {
		quote.Title = quote.Description //如果title为空
	}
	quote.MemberActions = permission.GetQuotePermission(ctx, userID, quote.UserID, quote.Status,
		per.CanEdit, per.CanDelete,
		per.CanClose, per.CanReopen, per.CanPin, per.CanHide, per.CanUnPin, per.CanShow,
		per.CanRecover)
	quote.ExtendsActions = permission.GetQuoteExtendsPermission(ctx, per.CanInviteOtherToAnswer)

	//@cws 增加author和piece信息
	//quote.Author, quote.Piece, err = qs.quotecommon.GetAuthorAndPiece(ctx, quote.UserID, quote.PieceID)

	//authorInfo, err := qs.quoteAuthorCommon.Info(ctx, quoteID, userID)
	quoteAuthorInfo, has, err := qs.quoteAuthorRepo.GetQuoteAuthor(ctx, quote.QuoteAuthorId)
	if err != nil {
		return
	}
	if !has {
		return
	}
	quote.QuoteAuthorBasicInfo = &schema.QuoteAuthorBasicInfo{
		ID:         quoteAuthorInfo.ID,
		AuthorName: quoteAuthorInfo.AuthorName,
		Avatar:     quoteAuthorInfo.Avatar,
	}
	if handler.GetEnableShortID(ctx) {
		quote.QuoteAuthorBasicInfo.ID = uid.EnShortID(quote.QuoteAuthorBasicInfo.ID)
	}

	quotePieceInfo, has, err := qs.quotePieceRepo.GetQuotePiece(ctx, quote.QuotePieceId)
	if err != nil {
		return
	}
	if !has {
		return
	}
	quote.QuotePieceBasicInfo = &schema.QuotePieceBasicInfo{
		ID:     quotePieceInfo.ID,
		Title:  quotePieceInfo.Title,
		Avatar: quotePieceInfo.Avatar,
	}
	if handler.GetEnableShortID(ctx) {
		quote.QuotePieceBasicInfo.ID = uid.EnShortID(quote.QuotePieceBasicInfo.ID)
	}

	return quote, nil
}

// GetQuoteAndAddPV get quote one
func (qs *QuoteService) GetQuoteAndAddPV(ctx context.Context, quoteID, loginUserID string,
	per schema.QuotePermission) (
	resp *schema.QuoteInfoResp, err error) {
	err = qs.quotecommon.UpdatePv(ctx, quoteID)
	if err != nil {
		log.Error(err)
	}
	return qs.GetQuote(ctx, quoteID, loginUserID, per)
}

func (qs *QuoteService) InviteUserInfo(ctx context.Context, quoteID string) (inviteList []*schema.UserBasicInfo, err error) {
	return qs.quotecommon.InviteUserInfo(ctx, quoteID)
}

func (qs *QuoteService) ChangeTag(ctx context.Context, objectTagData *schema.TagChange) error {
	return qs.tagCommon.ObjectChangeTag(ctx, objectTagData)
}

func (qs *QuoteService) CheckChangeReservedTag(ctx context.Context, oldobjectTagData, objectTagData []*entity.Tag) (bool, bool, []string, []string) {
	return qs.tagCommon.CheckChangeReservedTag(ctx, oldobjectTagData, objectTagData)
}

// PersonalQuotePage get quote list by user
func (qs *QuoteService) PersonalQuotePage(ctx context.Context, req *schema.PersonalQuotePageReq) (
	pageModel *pager.PageModel, err error) {

	userinfo, exist, err := qs.userCommon.GetUserBasicInfoByUserName(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.BadRequest(reason.UserNotFound)
	}
	search := &schema.QuotePageReq{}
	search.OrderCond = req.OrderCond
	search.Page = req.Page
	search.PageSize = req.PageSize
	search.UserIDBeSearched = userinfo.ID
	search.LoginUserID = req.LoginUserID
	// Only author and administrator can view the pending quote
	if req.LoginUserID == userinfo.ID || req.IsAdmin {
		search.ShowPending = true
	}
	quoteList, total, err := qs.GetQuotePage(ctx, search)
	if err != nil {
		return nil, err
	}
	userQuoteInfoList := make([]*schema.UserQuoteInfo, 0)
	for _, item := range quoteList {
		info := &schema.UserQuoteInfo{}
		_ = copier.Copy(info, item)
		status, ok := entity.AdminQuoteSearchStatusIntToString[item.Status]
		if ok {
			info.Status = status
		}
		userQuoteInfoList = append(userQuoteInfoList, info)
	}
	return pager.NewPageModel(total, userQuoteInfoList), nil
}

func (qs *QuoteService) PersonalAnswerPage(ctx context.Context, req *schema.PersonalAnswerPageReq) (
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
	//answerList, total, err := qs.quotecommon.AnswerCommon.PersonalAnswerPage(ctx, cond)
	//if err != nil {
	//	return nil, err
	//}
	//
	//answerlist := make([]*schema.AnswerInfo, 0)
	//userAnswerlist := make([]*schema.UserAnswerInfo, 0)
	//for _, item := range answerList {
	//	answerinfo := qs.quotecommon.AnswerCommon.ShowFormat(ctx, item)
	//	answerlist = append(answerlist, answerinfo)
	//	quoteIDs = append(quoteIDs, uid.DeShortID(item.QuoteID))
	//}
	//quoteMaps, err := qs.quotecommon.FindInfoByID(ctx, quoteIDs, req.LoginUserID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, item := range answerlist {
	//	_, ok := quoteMaps[item.QuoteID]
	//	if ok {
	//		item.QuoteInfo = quoteMaps[item.QuoteID]
	//	} else {
	//		continue
	//	}
	//	info := &schema.UserAnswerInfo{}
	//	_ = copier.Copy(info, item)
	//	info.AnswerID = item.ID
	//	info.QuoteID = item.QuoteID
	//	if item.QuoteInfo.Status == entity.QuoteStatusDeleted {
	//		info.QuoteInfo.Title = "Deleted quote"
	//
	//	}
	//	userAnswerlist = append(userAnswerlist, info)
	//}
	//
	//return pager.NewPageModel(total, userAnswerlist), nil
}

// PersonalCollectionPage get collection list by user
func (qs *QuoteService) PersonalCollectionPage(ctx context.Context, req *schema.PersonalCollectionPageReq) (
	pageModel *pager.PageModel, err error) {
	list := make([]*schema.QuoteInfoResp, 0)
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

	quoteMaps, err := qs.quotecommon.FindInfoByID(ctx, quoteIDs, req.UserID)
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
			if quoteMaps[id].Status == entity.QuoteStatusDeleted {
				quoteMaps[id].Title = "Deleted quote"
			}
			list = append(list, quoteMaps[id])
		}
	}

	return pager.NewPageModel(total, list), nil
}

func (qs *QuoteService) SearchUserTopList(ctx context.Context, userName string, loginUserID string) ([]*schema.UserQuoteInfo, []*schema.UserAnswerInfo, error) {

	answerlist := make([]*schema.AnswerInfo, 0)

	userAnswerlist := make([]*schema.UserAnswerInfo, 0)
	userQuotelist := make([]*schema.UserQuoteInfo, 0)
	_ = answerlist
	return userQuotelist, userAnswerlist, nil
	//
	//userinfo, Exist, err := qs.userCommon.GetUserBasicInfoByUserName(ctx, userName)
	//if err != nil {
	//	return userQuotelist, userAnswerlist, err
	//}
	//if !Exist {
	//	return userQuotelist, userAnswerlist, nil
	//}
	//search := &schema.QuotePageReq{}
	//search.OrderCond = "score"
	//search.Page = 0
	//search.PageSize = 5
	//search.UserIDBeSearched = userinfo.ID
	//search.LoginUserID = loginUserID
	//quotelist, _, err := qs.GetQuotePage(ctx, search)
	//if err != nil {
	//	return userQuotelist, userAnswerlist, err
	//}
	//answersearch := &entity.AnswerSearch{}
	//answersearch.UserID = userinfo.ID
	//answersearch.PageSize = 5
	//answersearch.Order = entity.AnswerSearchOrderByVote
	//quoteIDs := make([]string, 0)
	//answerList, _, err := qs.quotecommon.AnswerCommon.Search(ctx, answersearch)
	//if err != nil {
	//	return userQuotelist, userAnswerlist, err
	//}
	//for _, item := range answerList {
	//	answerinfo := qs.quotecommon.AnswerCommon.ShowFormat(ctx, item)
	//	answerlist = append(answerlist, answerinfo)
	//	quoteIDs = append(quoteIDs, item.QuoteID)
	//}
	//quoteMaps, err := qs.quotecommon.FindInfoByID(ctx, quoteIDs, loginUserID)
	//if err != nil {
	//	return userQuotelist, userAnswerlist, err
	//}
	//for _, item := range answerlist {
	//	_, ok := quoteMaps[item.QuoteID]
	//	if ok {
	//		item.QuoteInfo = quoteMaps[item.QuoteID]
	//	}
	//}
	//
	//for _, item := range quotelist {
	//	info := &schema.UserQuoteInfo{}
	//	_ = copier.Copy(info, item)
	//	info.UrlTitle = htmltext.UrlTitle(info.Title)
	//	userQuotelist = append(userQuotelist, info)
	//}
	//
	//for _, item := range answerlist {
	//	info := &schema.UserAnswerInfo{}
	//	_ = copier.Copy(info, item)
	//	info.AnswerID = item.ID
	//	info.QuoteID = item.QuoteID
	//	info.QuoteInfo.UrlTitle = htmltext.UrlTitle(info.QuoteInfo.Title)
	//	userAnswerlist = append(userAnswerlist, info)
	//}
	//
	//return userQuotelist, userAnswerlist, nil
}

// GetQuotesByTitle get quotes by title
func (qs *QuoteService) GetQuotesByTitle(ctx context.Context, title string) (
	resp []*schema.QuoteBaseInfo, err error) {
	resp = make([]*schema.QuoteBaseInfo, 0)
	if len(title) == 0 {
		return resp, nil
	}
	quotes, err := qs.quoteRepo.GetQuotesByTitle(ctx, title, 10)
	if err != nil {
		return resp, err
	}
	for _, quote := range quotes {
		item := &schema.QuoteBaseInfo{}
		item.ID = quote.ID
		item.Title = quote.Title
		item.UrlTitle = htmltext.UrlTitle(quote.Title)
		item.ViewCount = quote.ViewCount
		//item.AnswerCount = quote.AnswerCount
		item.CollectionCount = quote.CollectionCount
		item.FollowCount = quote.FollowCount
		status, ok := entity.AdminQuoteSearchStatusIntToString[quote.Status]
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

// SimilarQuote
func (qs *QuoteService) SimilarQuote(ctx context.Context, quoteID string, loginUserID string) ([]*schema.QuotePageResp, int64, error) {
	quote, err := qs.quotecommon.Info(ctx, quoteID, loginUserID)
	if err != nil {
		return nil, 0, nil
	}
	tagNames := make([]string, 0, len(quote.Tags))
	for _, tag := range quote.Tags {
		tagNames = append(tagNames, tag.SlugName)
	}
	search := &schema.QuotePageReq{}
	search.OrderCond = "hot"
	search.Page = 0
	search.PageSize = 6
	if len(tagNames) > 0 {
		search.Tag = tagNames[0]
	}
	search.LoginUserID = loginUserID
	similarQuotes, _, err := qs.GetQuotePage(ctx, search)
	if err != nil {
		return nil, 0, err
	}
	var result []*schema.QuotePageResp
	for _, v := range similarQuotes {
		if uid.DeShortID(v.ID) != quoteID {
			result = append(result, v)
		}
	}
	return result, int64(len(result)), nil
}

// GetQuotePage query quotes page
func (qs *QuoteService) GetQuotePage(ctx context.Context, req *schema.QuotePageReq) (
	quotes []*schema.QuotePageResp, total int64, err error) {
	quotes = make([]*schema.QuotePageResp, 0)
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

	if req.OrderCond == schema.QuoteOrderCondHot {
		req.InDays = schema.HotInDays
	}

	quoteList, total, err := qs.quoteRepo.GetQuotePage(ctx, req.Page, req.PageSize,
		tagIDs, req.UserIDBeSearched, req.OrderCond, req.InDays, showHidden, req.ShowPending)
	if err != nil {
		return nil, 0, err
	}
	quotes, err = qs.quotecommon.FormatQuotesPage(ctx, quoteList, req.LoginUserID, req.OrderCond)
	if err != nil {
		return nil, 0, err
	}
	return quotes, total, nil
}

// GetRecommendQuotePage retrieves recommended quote page based on following tags and quotes.
func (qs *QuoteService) GetRecommendQuotePage(ctx context.Context, req *schema.QuotePageReq) (
	quotes []*schema.QuotePageResp, total int64, err error) {
	followingTagsResp, err := qs.tagService.GetFollowingTags(ctx, req.LoginUserID)
	if err != nil {
		return nil, 0, err
	}
	tagIDs := make([]string, 0, len(followingTagsResp))
	for _, tag := range followingTagsResp {
		tagIDs = append(tagIDs, tag.TagID)
	}

	activityType, err := qs.activityRepo.GetActivityTypeByObjectType(ctx, constant.QuoteObjectType, "follow")
	if err != nil {
		return nil, 0, err
	}
	activities, err := qs.activityRepo.GetUserActivitysByActivityType(ctx, req.LoginUserID, activityType)
	if err != nil {
		return nil, 0, err
	}

	followedQuoteIDs := make([]string, 0, len(activities))
	for _, activity := range activities {
		if activity.Cancelled == entity.ActivityCancelled {
			continue
		}
		followedQuoteIDs = append(followedQuoteIDs, activity.ObjectID)
	}
	quoteList, total, err := qs.quoteRepo.GetRecommendQuotePageByTags(ctx, req.LoginUserID, tagIDs, followedQuoteIDs, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	quotes, err = qs.quotecommon.FormatQuotesPage(ctx, quoteList, req.LoginUserID, "frequent")
	if err != nil {
		return nil, 0, err
	}

	return quotes, total, nil
}

func (qs *QuoteService) AdminSetQuoteStatus(ctx context.Context, req *schema.AdminUpdateQuoteStatusReq) error {
	setStatus, ok := entity.AdminQuoteSearchStatus[req.Status]
	if !ok {
		return errors.BadRequest(reason.RequestFormatError)
	}
	quoteInfo, exist, err := qs.quoteRepo.GetQuote(ctx, req.QuoteID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.QuoteNotFound)
	}
	err = qs.quoteRepo.UpdateQuoteStatus(ctx, quoteInfo.ID, setStatus)
	if err != nil {
		return err
	}

	msg := &schema.NotificationMsg{}
	if setStatus == entity.QuoteStatusDeleted {
		// #2372 In order to simplify the process and complexity, as well as to consider if it is in-house,
		// facing the problem of recovery.
		//err = qs.answerActivityService.DeleteQuote(ctx, quoteInfo.ID, quoteInfo.CreatedAt, quoteInfo.VoteCount)
		//if err != nil {
		//	log.Errorf("admin delete quote then rank rollback error %s", err.Error())
		//}
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           quoteInfo.UserID,
			TriggerUserID:    converter.StringToInt64(req.UserID),
			ObjectID:         quoteInfo.ID,
			OriginalObjectID: quoteInfo.ID,
			ActivityTypeKey:  constant.ActQuoteDeleted,
		})
		msg.NotificationAction = constant.NotificationYourQuoteIsClosed
	}
	if setStatus == entity.QuoteStatusAvailable && quoteInfo.Status == entity.QuoteStatusClosed {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           quoteInfo.UserID,
			TriggerUserID:    converter.StringToInt64(req.UserID),
			ObjectID:         quoteInfo.ID,
			OriginalObjectID: quoteInfo.ID,
			ActivityTypeKey:  constant.ActQuoteReopened,
		})
	}
	if setStatus == entity.QuoteStatusClosed && quoteInfo.Status != entity.QuoteStatusClosed {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           quoteInfo.UserID,
			TriggerUserID:    converter.StringToInt64(req.UserID),
			ObjectID:         quoteInfo.ID,
			OriginalObjectID: quoteInfo.ID,
			ActivityTypeKey:  constant.ActQuoteClosed,
		})
		msg.NotificationAction = constant.NotificationYourQuoteIsClosed
	}
	// recover
	if setStatus == entity.QuoteStatusAvailable && quoteInfo.Status == entity.QuoteStatusDeleted {
		qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
			UserID:           req.UserID,
			TriggerUserID:    converter.StringToInt64(req.UserID),
			ObjectID:         quoteInfo.ID,
			OriginalObjectID: quoteInfo.ID,
			ActivityTypeKey:  constant.ActQuoteUndeleted,
		})
	}

	if len(msg.NotificationAction) > 0 {
		msg.ObjectID = quoteInfo.ID
		msg.Type = schema.NotificationTypeInbox
		msg.ReceiverUserID = quoteInfo.UserID
		msg.TriggerUserID = req.UserID
		msg.ObjectType = constant.QuoteObjectType
		qs.notificationQueueService.Send(ctx, msg)
	}
	return nil
}

func (qs *QuoteService) AdminQuotePage(
	ctx context.Context, req *schema.AdminQuotePageReq) (
	resp *pager.PageModel, err error) {

	list := make([]*schema.AdminQuoteInfo, 0)
	quoteList, count, err := qs.quoteRepo.AdminQuotePage(ctx, req)
	if err != nil {
		return nil, err
	}

	userIds := make([]string, 0)
	for _, info := range quoteList {
		item := &schema.AdminQuoteInfo{}
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
func (qs *QuoteService) AdminAnswerPage(ctx context.Context, req *schema.AdminAnswerPageReq) (
	resp *pager.PageModel, err error) {
	return
	//
	//answerList, count, err := qs.quotecommon.AnswerCommon.AdminSearchList(ctx, req)
	//if err != nil {
	//	return nil, err
	//}
	//
	//quoteIDs := make([]string, 0)
	//userIds := make([]string, 0)
	//answerResp := make([]*schema.AdminAnswerInfo, 0)
	//for _, item := range answerList {
	//	answerInfo := qs.quotecommon.AnswerCommon.AdminShowFormat(ctx, item)
	//	answerResp = append(answerResp, answerInfo)
	//	quoteIDs = append(quoteIDs, item.QuoteID)
	//	userIds = append(userIds, item.UserID)
	//}
	//userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIds)
	//if err != nil {
	//	return nil, err
	//}
	//quoteMaps, err := qs.quotecommon.FindInfoByID(ctx, quoteIDs, req.LoginUserID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//for _, item := range answerResp {
	//	if q, ok := quoteMaps[item.QuoteID]; ok {
	//		item.QuoteInfo.Title = q.Title
	//	}
	//	if u, ok := userInfoMap[item.UserID]; ok {
	//		item.UserInfo = u
	//	}
	//}
	//return pager.NewPageModel(count, answerResp), nil
}

func (qs *QuoteService) changeQuoteToRevision(ctx context.Context, quoteInfo *entity.Quote, tags []*entity.Tag) (
	quoteRevision *entity.QuoteWithTagsRevision, err error) {
	quoteRevision = &entity.QuoteWithTagsRevision{}
	quoteRevision.Quote = *quoteInfo

	for _, tag := range tags {
		item := &entity.TagSimpleInfoForRevision{}
		_ = copier.Copy(item, tag)
		quoteRevision.Tags = append(quoteRevision.Tags, item)
	}
	return quoteRevision, nil
}

func (qs *QuoteService) SitemapCron(ctx context.Context) {
	siteSeo, err := qs.siteInfoService.GetSiteSeo(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	ctx = context.WithValue(ctx, constant.ShortIDFlag, siteSeo.IsShortLink())
	qs.quotecommon.SitemapCron(ctx)
}
