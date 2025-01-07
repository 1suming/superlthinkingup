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

package controller_quote

import (
	"github.com/apache/incubator-answer/internal/service_quote"
	"net/http"

	"github.com/apache/incubator-answer/internal/base/handler"
	"github.com/apache/incubator-answer/internal/base/middleware"
	"github.com/apache/incubator-answer/internal/base/pager"
	"github.com/apache/incubator-answer/internal/base/reason"
	"github.com/apache/incubator-answer/internal/base/translator"
	"github.com/apache/incubator-answer/internal/base/validator"
	"github.com/apache/incubator-answer/internal/entity"
	"github.com/apache/incubator-answer/internal/schema"
	"github.com/apache/incubator-answer/internal/service/action"
	"github.com/apache/incubator-answer/internal/service/content"
	"github.com/apache/incubator-answer/internal/service/permission"
	"github.com/apache/incubator-answer/internal/service/rank"
	"github.com/apache/incubator-answer/internal/service/siteinfo_common"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"
)

// QuoteAuthorController QuoteAuthor controller
type QuoteAuthorController struct {
	QuoteAuthorService  *service_quote.QuoteAuthorService
	answerService       *content.AnswerService
	rankService         *rank.RankService
	siteInfoService     siteinfo_common.SiteInfoCommonService
	actionService       *action.CaptchaService
	rateLimitMiddleware *middleware.RateLimitMiddleware
}

// NewQuoteAuthorController new controller
func NewQuoteAuthorController(
	QuoteAuthorService *service_quote.QuoteAuthorService,
	answerService *content.AnswerService,
	rankService *rank.RankService,
	siteInfoService siteinfo_common.SiteInfoCommonService,
	actionService *action.CaptchaService,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
) *QuoteAuthorController {
	return &QuoteAuthorController{
		QuoteAuthorService:  QuoteAuthorService,
		answerService:       answerService,
		rankService:         rankService,
		siteInfoService:     siteInfoService,
		actionService:       actionService,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}

// RemoveQuoteAuthor delete QuoteAuthor
// @Summary delete QuoteAuthor
// @Description delete QuoteAuthor
// @Tags QuoteAuthor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.RemoveQuoteAuthorReq true "QuoteAuthor"
// @Success 200 {object} handler.RespBody
// @Router  /answer/api/v1/QuoteAuthor [delete]
func (qc *QuoteAuthorController) RemoveQuoteAuthor(ctx *gin.Context) {
	req := &schema.RemoveQuoteAuthorReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	req.IsAdmin = middleware.GetIsAdminFromContext(ctx)
	isAdmin := middleware.GetUserIsAdminModerator(ctx)
	if !isAdmin {
		captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionDelete, req.UserID, req.CaptchaID, req.CaptchaCode)
		if !captchaPass {
			errFields := append([]*validator.FormErrorField{}, &validator.FormErrorField{
				ErrorField: "captcha_code",
				ErrorMsg:   translator.Tr(handler.GetLang(ctx), reason.CaptchaVerificationFailed),
			})
			handler.HandleResponse(ctx, errors.BadRequest(reason.CaptchaVerificationFailed), errFields)
			return
		}
	}

	can, err := qc.rankService.CheckOperationPermission(ctx, req.UserID, permission.QuoteAuthorDelete, req.ID)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	err = qc.QuoteAuthorService.RemoveQuoteAuthor(ctx, req)
	if !isAdmin {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionDelete, req.UserID)
	}
	handler.HandleResponse(ctx, err, nil)
}

// OperationQuoteAuthor Operation QuoteAuthor
// @Summary Operation QuoteAuthor
// @Description Operation QuoteAuthor \n operation [pin unpin hide show]
// @Tags QuoteAuthor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.OperationQuoteAuthorReq true "QuoteAuthor"
// @Success 200 {object} handler.RespBody
// @Router  /answer/api/v1/QuoteAuthor/operation [put]
func (qc *QuoteAuthorController) OperationQuoteAuthor(ctx *gin.Context) {
	req := &schema.OperationQuoteAuthorReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.QuoteAuthorPin,
		permission.QuoteAuthorUnPin,
		permission.QuoteAuthorHide,
		permission.QuoteAuthorShow,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	req.CanPin = canList[0]
	req.CanList = canList[1]
	if (req.Operation == schema.QuoteAuthorOperationPin || req.Operation == schema.QuoteAuthorOperationUnPin) && !req.CanPin {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	if (req.Operation == schema.QuoteAuthorOperationHide || req.Operation == schema.QuoteAuthorOperationShow) && !req.CanList {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	err = qc.QuoteAuthorService.OperationQuoteAuthor(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// CloseQuoteAuthor Close QuoteAuthor
// @Summary Close QuoteAuthor
// @Description Close QuoteAuthor
// @Tags QuoteAuthor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.CloseQuoteAuthorReq true "QuoteAuthor"
// @Success 200 {object} handler.RespBody
// @Router  /answer/api/v1/QuoteAuthor/status [put]
func (qc *QuoteAuthorController) CloseQuoteAuthor(ctx *gin.Context) {
	req := &schema.CloseQuoteAuthorReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := qc.rankService.CheckOperationPermission(ctx, req.UserID, permission.QuoteAuthorClose, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = qc.QuoteAuthorService.CloseQuoteAuthor(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// ReopenQuoteAuthor reopen QuoteAuthor
// @Summary reopen QuoteAuthor
// @Description reopen QuoteAuthor
// @Tags QuoteAuthor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.ReopenQuoteAuthorReq true "QuoteAuthor"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuoteAuthor/reopen [put]
func (qc *QuoteAuthorController) ReopenQuoteAuthor(ctx *gin.Context) {
	req := &schema.ReopenQuoteAuthorReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.QuoteAuthorID = uid.DeShortID(req.QuoteAuthorID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := qc.rankService.CheckOperationPermission(ctx, req.UserID, permission.QuoteAuthorReopen, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = qc.QuoteAuthorService.ReopenQuoteAuthor(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// GetQuoteAuthor get QuoteAuthor details
// @Summary get QuoteAuthor details
// @Description get QuoteAuthor details
// @Tags QuoteAuthor
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id query string true "QuoteAuthor TagID"  default(1)
// @Success 200 {string} string ""
// @Router /answer/api/v1/QuoteAuthor/info [get]
func (qc *QuoteAuthorController) GetQuoteAuthor(ctx *gin.Context) {
	id := ctx.Query("id")
	id = uid.DeShortID(id)
	userID := middleware.GetLoginUserIDFromContext(ctx)
	req := schema.QuoteAuthorPermission{}
	canList, err := qc.rankService.CheckOperationPermissions(ctx, userID, []string{
		permission.QuoteAuthorEdit,
		permission.QuoteAuthorDelete,
		permission.QuoteAuthorClose,
		permission.QuoteAuthorReopen,
		permission.QuoteAuthorPin,
		permission.QuoteAuthorUnPin,
		permission.QuoteAuthorHide,
		permission.QuoteAuthorShow,
		permission.AnswerInviteSomeoneToAnswer,
		permission.QuoteAuthorUnDelete,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	objectOwner := qc.rankService.CheckOperationObjectOwner(ctx, userID, id)

	req.CanEdit = canList[0] || objectOwner
	req.CanDelete = canList[1]
	req.CanClose = canList[2]
	req.CanReopen = canList[3]
	req.CanPin = canList[4]
	req.CanUnPin = canList[5]
	req.CanHide = canList[6]
	req.CanShow = canList[7]
	req.CanInviteOtherToAnswer = canList[8]
	req.CanRecover = canList[9]

	info, err := qc.QuoteAuthorService.GetQuoteAuthorAndAddPV(ctx, id, userID, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if handler.GetEnableShortID(ctx) {
		info.ID = uid.EnShortID(info.ID)
	}
	handler.HandleResponse(ctx, nil, info)
}

// GetQuoteAuthorInviteUserInfo get QuoteAuthor invite user info
// @Summary get QuoteAuthor invite user info
// @Description get QuoteAuthor invite user info
// @Tags QuoteAuthor
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id query string true "QuoteAuthor ID"  default(1)
// @Success 200 {string} string ""
// @Router /answer/api/v1/QuoteAuthor/invite [get]
func (qc *QuoteAuthorController) GetQuoteAuthorInviteUserInfo(ctx *gin.Context) {
	QuoteAuthorID := uid.DeShortID(ctx.Query("id"))
	resp, err := qc.QuoteAuthorService.InviteUserInfo(ctx, QuoteAuthorID)
	handler.HandleResponse(ctx, err, resp)

}

// SimilarQuoteAuthor godoc
// @Summary Search Similar QuoteAuthor
// @Description Search Similar QuoteAuthor
// @Tags QuoteAuthor
// @Accept  json
// @Produce  json
// @Param QuoteAuthor_id query string true "QuoteAuthor_id"  default()
// @Success 200 {string} string ""
// @Router /answer/api/v1/QuoteAuthor/similar/tag [get]
func (qc *QuoteAuthorController) SimilarQuoteAuthor(ctx *gin.Context) {
	QuoteAuthorID := ctx.Query("QuoteAuthor_id")
	QuoteAuthorID = uid.DeShortID(QuoteAuthorID)
	userID := middleware.GetLoginUserIDFromContext(ctx)
	list, count, err := qc.QuoteAuthorService.SimilarQuoteAuthor(ctx, QuoteAuthorID, userID)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, nil, gin.H{
		"list":  list,
		"count": count,
	})
}

// QuoteAuthorPage get QuoteAuthors by page
// @Summary get QuoteAuthors by page
// @Description get QuoteAuthors by page
// @Tags QuoteAuthor
// @Accept  json
// @Produce  json
// @Param data body schema.QuoteAuthorPageReq  true "QuoteAuthorPageReq"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.QuoteAuthorPageResp}}
// @Router /answer/api/v1/QuoteAuthor/page [get]
func (qc *QuoteAuthorController) QuoteAuthorPage(ctx *gin.Context) {
	req := &schema.QuoteAuthorPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)

	QuoteAuthors, total, err := qc.QuoteAuthorService.GetQuoteAuthorPage(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, nil, pager.NewPageModel(total, QuoteAuthors))
}

// QuoteAuthorRecommendPage get recommend QuoteAuthors by page
// @Summary get recommend QuoteAuthors by page
// @Description get recommend QuoteAuthors by page
// @Tags QuoteAuthor
// @Accept  json
// @Produce  json
// @Param data body schema.QuoteAuthorPageReq  true "QuoteAuthorPageReq"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.QuoteAuthorPageResp}}
// @Router /answer/api/v1/QuoteAuthor/recommend/page [get]
func (qc *QuoteAuthorController) QuoteAuthorRecommendPage(ctx *gin.Context) {
	req := &schema.QuoteAuthorPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)

	if req.LoginUserID == "" {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	QuoteAuthors, total, err := qc.QuoteAuthorService.GetRecommendQuoteAuthorPage(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, nil, pager.NewPageModel(total, QuoteAuthors))
}

// AddQuoteAuthor add QuoteAuthor
// @Summary add QuoteAuthor
// @Description add QuoteAuthor
// @Tags QuoteAuthor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuoteAuthorAdd true "QuoteAuthor"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuoteAuthor [post]
func (qc *QuoteAuthorController) AddQuoteAuthor(ctx *gin.Context) {
	req := &schema.QuoteAuthorAdd{}
	errFields := handler.BindAndCheckReturnErr(ctx, req)
	if ctx.IsAborted() {
		return
	}
	reject, rejectKey := qc.rateLimitMiddleware.DuplicateRequestRejection(ctx, req)
	if reject {
		return
	}
	defer func() {
		// If status is not 200 means that the bad request has been returned, so the record should be cleared
		if ctx.Writer.Status() != http.StatusOK {
			qc.rateLimitMiddleware.DuplicateRequestClear(ctx, rejectKey)
		}
	}()

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, requireRanks, err := qc.rankService.CheckOperationPermissionsForRanks(ctx, req.UserID, []string{
		permission.QuoteAuthorAdd,
		permission.QuoteAuthorEdit,
		permission.QuoteAuthorDelete,
		permission.QuoteAuthorClose,
		permission.QuoteAuthorReopen,
		permission.TagUseReservedTag,
		permission.TagAdd,
		permission.LinkUrlLimit,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	linkUrlLimitUser := canList[7]
	isAdmin := middleware.GetUserIsAdminModerator(ctx)
	if !isAdmin || !linkUrlLimitUser {
		captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionQuoteAuthor, req.UserID, req.CaptchaID, req.CaptchaCode)
		if !captchaPass {
			errFields := append([]*validator.FormErrorField{}, &validator.FormErrorField{
				ErrorField: "captcha_code",
				ErrorMsg:   translator.Tr(handler.GetLang(ctx), reason.CaptchaVerificationFailed),
			})
			handler.HandleResponse(ctx, errors.BadRequest(reason.CaptchaVerificationFailed), errFields)
			return
		}
	}

	req.CanAdd = canList[0]
	req.CanEdit = canList[1]
	req.CanDelete = canList[2]
	req.CanClose = canList[3]
	req.CanReopen = canList[4]
	req.CanUseReservedTag = canList[5]
	req.CanAddTag = canList[6]
	if !req.CanAdd {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	// can add tag
	hasNewTag, err := qc.QuoteAuthorService.HasNewTag(ctx, req.Tags)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !req.CanAddTag && hasNewTag {
		lang := handler.GetLang(ctx)
		msg := translator.TrWithData(lang, reason.NoEnoughRankToOperate, &schema.PermissionTrTplData{Rank: requireRanks[6]})
		handler.HandleResponse(ctx, errors.Forbidden(reason.NoEnoughRankToOperate).WithMsg(msg), nil)
		return
	}

	errList, err := qc.QuoteAuthorService.CheckAddQuoteAuthor(ctx, req)
	if err != nil {
		errlist, ok := errList.([]*validator.FormErrorField)
		if ok {
			errFields = append(errFields, errlist...)
		}
	}

	if len(errFields) > 0 {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		return
	}

	req.UserAgent = ctx.GetHeader("User-Agent")
	req.IP = ctx.ClientIP()

	resp, err := qc.QuoteAuthorService.AddQuoteAuthor(ctx, req)
	if err != nil {
		errlist, ok := resp.([]*validator.FormErrorField)
		if ok {
			errFields = append(errFields, errlist...)
		}
	}

	if len(errFields) > 0 {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		return
	}
	if !isAdmin || !linkUrlLimitUser {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionQuoteAuthor, req.UserID)
	}
	handler.HandleResponse(ctx, err, resp)
}

// AddQuoteAuthorByAnswer add QuoteAuthor
// @Summary add QuoteAuthor and answer
// @Description add QuoteAuthor and answer
// @Tags QuoteAuthor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuoteAuthorAddByAnswer true "QuoteAuthor"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuoteAuthor/answer [post]
func (qc *QuoteAuthorController) AddQuoteAuthorByAnswer(ctx *gin.Context) {
	req := &schema.QuoteAuthorAddByAnswer{}
	errFields := handler.BindAndCheckReturnErr(ctx, req)
	if ctx.IsAborted() {
		return
	}
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.QuoteAuthorAdd,
		permission.QuoteAuthorEdit,
		permission.QuoteAuthorDelete,
		permission.QuoteAuthorClose,
		permission.QuoteAuthorReopen,
		permission.TagUseReservedTag,
		permission.LinkUrlLimit,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}

	linkUrlLimitUser := canList[6]
	isAdmin := middleware.GetUserIsAdminModerator(ctx)
	if !isAdmin || !linkUrlLimitUser {
		captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionQuoteAuthor, req.UserID, req.CaptchaID, req.CaptchaCode)
		if !captchaPass {
			errFields := append([]*validator.FormErrorField{}, &validator.FormErrorField{
				ErrorField: "captcha_code",
				ErrorMsg:   translator.Tr(handler.GetLang(ctx), reason.CaptchaVerificationFailed),
			})
			handler.HandleResponse(ctx, errors.BadRequest(reason.CaptchaVerificationFailed), errFields)
			return
		}
	}
	req.CanAdd = canList[0]
	req.CanEdit = canList[1]
	req.CanDelete = canList[2]
	req.CanClose = canList[3]
	req.CanReopen = canList[4]
	req.CanUseReservedTag = canList[5]
	if !req.CanAdd {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	QuoteAuthorReq := new(schema.QuoteAuthorAdd)
	err = copier.Copy(QuoteAuthorReq, req)
	if err != nil {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RequestFormatError), nil)
		return
	}
	errList, err := qc.QuoteAuthorService.CheckAddQuoteAuthor(ctx, QuoteAuthorReq)
	if err != nil {
		errlist, ok := errList.([]*validator.FormErrorField)
		if ok {
			errFields = append(errFields, errlist...)
		}
	}

	if len(errFields) > 0 {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		return
	}

	req.UserAgent = ctx.GetHeader("User-Agent")
	req.IP = ctx.ClientIP()
	resp, err := qc.QuoteAuthorService.AddQuoteAuthor(ctx, QuoteAuthorReq)
	if err != nil {
		errlist, ok := resp.([]*validator.FormErrorField)
		if ok {
			errFields = append(errFields, errlist...)
		}
	}

	if !isAdmin || !linkUrlLimitUser {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionQuoteAuthor, req.UserID)
	}

	if len(errFields) > 0 {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		return
	}
	////add the QuoteAuthor id to the answer
	//QuoteAuthorInfo, ok := resp.(*schema.QuoteAuthorInfoResp)
	//if ok {
	//	answerReq := &schema.AnswerAddReq{}
	//	answerReq.QuoteAuthorID = uid.DeShortID(QuoteAuthorInfo.ID)
	//	answerReq.UserID = middleware.GetLoginUserIDFromContext(ctx)
	//	answerReq.Content = req.AnswerContent
	//	answerReq.HTML = req.AnswerHTML
	//	answerID, err := qc.answerService.Insert(ctx, answerReq)
	//	if err != nil {
	//		handler.HandleResponse(ctx, err, nil)
	//		return
	//	}
	//	info, QuoteAuthorInfo, has, err := qc.answerService.Get(ctx, answerID, req.UserID)
	//	if err != nil {
	//		handler.HandleResponse(ctx, err, nil)
	//		return
	//	}
	//	if !has {
	//		handler.HandleResponse(ctx, nil, nil)
	//		return
	//	}
	//	handler.HandleResponse(ctx, err, gin.H{
	//		"info":  info,
	//		"QuoteAuthor": QuoteAuthorInfo,
	//	})
	//	return
	//}

	handler.HandleResponse(ctx, err, resp)
}

// UpdateQuoteAuthor update QuoteAuthor
// @Summary update QuoteAuthor
// @Description update QuoteAuthor
// @Tags QuoteAuthor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuoteAuthorUpdate true "QuoteAuthor"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuoteAuthor [put]
func (qc *QuoteAuthorController) UpdateQuoteAuthor(ctx *gin.Context) {
	req := &schema.QuoteAuthorUpdate{}
	errFields := handler.BindAndCheckReturnErr(ctx, req)
	if ctx.IsAborted() {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, requireRanks, err := qc.rankService.CheckOperationPermissionsForRanks(ctx, req.UserID, []string{
		permission.QuoteAuthorEdit,
		permission.QuoteAuthorDelete,
		permission.QuoteAuthorEditWithoutReview,
		permission.TagUseReservedTag,
		permission.TagAdd,
		permission.LinkUrlLimit,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	linkUrlLimitUser := canList[5]
	isAdmin := middleware.GetUserIsAdminModerator(ctx)
	if !isAdmin || !linkUrlLimitUser {
		captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionEdit, req.UserID, req.CaptchaID, req.CaptchaCode)
		if !captchaPass {
			errFields := append([]*validator.FormErrorField{}, &validator.FormErrorField{
				ErrorField: "captcha_code",
				ErrorMsg:   translator.Tr(handler.GetLang(ctx), reason.CaptchaVerificationFailed),
			})
			handler.HandleResponse(ctx, errors.BadRequest(reason.CaptchaVerificationFailed), errFields)
			return
		}
	}

	objectOwner := qc.rankService.CheckOperationObjectOwner(ctx, req.UserID, req.ID)
	req.CanEdit = canList[0] || objectOwner
	req.CanDelete = canList[1]
	req.NoNeedReview = canList[2] || objectOwner
	req.CanUseReservedTag = canList[3]
	req.CanAddTag = canList[4]
	if !req.CanEdit {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	errlist, err := qc.QuoteAuthorService.UpdateQuoteAuthorCheckTags(ctx, req)
	if err != nil {
		errFields = append(errFields, errlist...)
	}

	if len(errFields) > 0 {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		return
	}

	// can add tag
	hasNewTag, err := qc.QuoteAuthorService.HasNewTag(ctx, req.Tags)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !req.CanAddTag && hasNewTag {
		lang := handler.GetLang(ctx)
		msg := translator.TrWithData(lang, reason.NoEnoughRankToOperate, &schema.PermissionTrTplData{Rank: requireRanks[4]})
		handler.HandleResponse(ctx, errors.Forbidden(reason.NoEnoughRankToOperate).WithMsg(msg), nil)
		return
	}

	resp, err := qc.QuoteAuthorService.UpdateQuoteAuthor(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, resp)
		return
	}
	respInfo, ok := resp.(*schema.QuoteAuthorInfoResp)
	if !ok {
		handler.HandleResponse(ctx, err, resp)
		return
	}
	if !isAdmin || !linkUrlLimitUser {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionEdit, req.UserID)
	}
	handler.HandleResponse(ctx, nil, &schema.UpdateQuoteAuthorResp{UrlAuthorName: respInfo.UrlAuthorName, WaitForReview: !req.NoNeedReview})
}

// QuoteAuthorRecover recover deleted QuoteAuthor
// @Summary recover deleted QuoteAuthor
// @Description recover deleted QuoteAuthor
// @Tags QuoteAuthor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuoteAuthorRecoverReq true "QuoteAuthor"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuoteAuthor/recover [post]
func (qc *QuoteAuthorController) QuoteAuthorRecover(ctx *gin.Context) {
	req := &schema.QuoteAuthorRecoverReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.QuoteAuthorID = uid.DeShortID(req.QuoteAuthorID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.QuoteAuthorUnDelete,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !canList[0] {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = qc.QuoteAuthorService.RecoverQuoteAuthor(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// UpdateQuoteAuthorInviteUser update QuoteAuthor invite user
// @Summary update QuoteAuthor invite user
// @Description update QuoteAuthor invite user
// @Tags QuoteAuthor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuoteAuthorUpdateInviteUser true "QuoteAuthor"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuoteAuthor/invite [put]
func (qc *QuoteAuthorController) UpdateQuoteAuthorInviteUser(ctx *gin.Context) {
	req := &schema.QuoteAuthorUpdateInviteUser{}
	errFields := handler.BindAndCheckReturnErr(ctx, req)
	if ctx.IsAborted() {
		return
	}
	if len(errFields) > 0 {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	isAdmin := middleware.GetUserIsAdminModerator(ctx)
	if !isAdmin {
		captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionInvitationAnswer, req.UserID, req.CaptchaID, req.CaptchaCode)
		if !captchaPass {
			errFields := append([]*validator.FormErrorField{}, &validator.FormErrorField{
				ErrorField: "captcha_code",
				ErrorMsg:   translator.Tr(handler.GetLang(ctx), reason.CaptchaVerificationFailed),
			})
			handler.HandleResponse(ctx, errors.BadRequest(reason.CaptchaVerificationFailed), errFields)
			return
		}
	}

	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.AnswerInviteSomeoneToAnswer,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}

	req.CanInviteOtherToAnswer = canList[0]
	if !req.CanInviteOtherToAnswer {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	err = qc.QuoteAuthorService.UpdateQuoteAuthorInviteUser(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !isAdmin {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionInvitationAnswer, req.UserID)
	}
	handler.HandleResponse(ctx, nil, nil)
}

// GetSimilarQuoteAuthors fuzzy query similar QuoteAuthors based on title
// @Summary fuzzy query similar QuoteAuthors based on title
// @Description fuzzy query similar QuoteAuthors based on title
// @Tags QuoteAuthor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param title query string true "title"  default(string)
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuoteAuthor/similar [get]
func (qc *QuoteAuthorController) GetSimilarQuoteAuthors(ctx *gin.Context) {
	title := ctx.Query("title")
	resp, err := qc.QuoteAuthorService.GetQuoteAuthorsByAuthorName(ctx, title)
	handler.HandleResponse(ctx, err, resp)
}

// UserTop godoc
// @Summary UserTop
// @Description UserTop
// @Tags QuoteAuthor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username query string true "username"  default(string)
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/personal/qa/top [get]
func (qc *QuoteAuthorController) UserTop(ctx *gin.Context) {
	userName := ctx.Query("username")
	userID := middleware.GetLoginUserIDFromContext(ctx)
	QuoteAuthorList, answerList, err := qc.QuoteAuthorService.SearchUserTopList(ctx, userName, userID)
	handler.HandleResponse(ctx, err, gin.H{
		"QuoteAuthor": QuoteAuthorList,
		"answer":      answerList,
	})
}

// PersonalQuoteAuthorPage list personal QuoteAuthors
// @Summary list personal QuoteAuthors
// @Description list personal QuoteAuthors
// @Tags Personal
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username query string true "username"  default(string)
// @Param order query string true "order"  Enums(newest,score)
// @Param page query string true "page"  default(0)
// @Param page_size query string true "page_size" default(20)
// @Success 200 {object} handler.RespBody
// @Router /personal/QuoteAuthor/page [get]
func (qc *QuoteAuthorController) PersonalQuoteAuthorPage(ctx *gin.Context) {
	req := &schema.PersonalQuoteAuthorPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	req.IsAdmin = middleware.GetUserIsAdminModerator(ctx)
	resp, err := qc.QuoteAuthorService.PersonalQuoteAuthorPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// PersonalAnswerPage list personal answers
// @Summary list personal answers
// @Description list personal answers
// @Tags Personal
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username query string true "username"  default(string)
// @Param order query string true "order"  Enums(newest,score)
// @Param page query string true "page"  default(0)
// @Param page_size query string true "page_size"  default(20)
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/personal/answer/page [get]
func (qc *QuoteAuthorController) PersonalAnswerPage(ctx *gin.Context) {
	req := &schema.PersonalAnswerPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	req.IsAdmin = middleware.GetUserIsAdminModerator(ctx)
	resp, err := qc.QuoteAuthorService.PersonalAnswerPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// PersonalCollectionPage list personal collections
// @Summary list personal collections
// @Description list personal collections
// @Tags Collection
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query string true "page"  default(0)
// @Param page_size query string true "page_size"  default(20)
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/personal/collection/page [get]
func (qc *QuoteAuthorController) PersonalCollectionPage(ctx *gin.Context) {
	req := &schema.PersonalCollectionPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	resp, err := qc.QuoteAuthorService.PersonalCollectionPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// AdminQuoteAuthorPage admin QuoteAuthor page
// @Summary AdminQuoteAuthorPage admin QuoteAuthor page
// @Description Status:[available,closed,deleted,pending]
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "page size"
// @Param page_size query int false "page size"
// @Param status query string false "user status" Enums(available, closed, deleted, pending)
// @Param query query string false "QuoteAuthor id or title"
// @Success 200 {object} handler.RespBody
// @Router /answer/admin/api/QuoteAuthor/page [get]
func (qc *QuoteAuthorController) AdminQuoteAuthorPage(ctx *gin.Context) {
	req := &schema.AdminQuoteAuthorPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	resp, err := qc.QuoteAuthorService.AdminQuoteAuthorPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// AdminAnswerPage admin answer page
// @Summary AdminAnswerPage admin answer page
// @Description Status:[available,deleted,pending]
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "page size"
// @Param page_size query int false "page size"
// @Param status query string false "user status" Enums(available,deleted,pending)
// @Param query query string false "answer id or QuoteAuthor title"
// @Param QuoteAuthor_id query string false "QuoteAuthor id"
// @Success 200 {object} handler.RespBody
// @Router /answer/admin/api/answer/page [get]
func (qc *QuoteAuthorController) AdminAnswerPage(ctx *gin.Context) {
	req := &schema.AdminAnswerPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	resp, err := qc.QuoteAuthorService.AdminAnswerPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// AdminUpdateQuoteAuthorStatus update QuoteAuthor status
// @Summary update QuoteAuthor status
// @Description update QuoteAuthor status
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.AdminUpdateQuoteAuthorStatusReq true "AdminUpdateQuoteAuthorStatusReq"
// @Success 200 {object} handler.RespBody
// @Router /answer/admin/api/QuoteAuthor/status [put]
func (qc *QuoteAuthorController) AdminUpdateQuoteAuthorStatus(ctx *gin.Context) {
	req := &schema.AdminUpdateQuoteAuthorStatusReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.QuoteAuthorID = uid.DeShortID(req.QuoteAuthorID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	err := qc.QuoteAuthorService.AdminSetQuoteAuthorStatus(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}
