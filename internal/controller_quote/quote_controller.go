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
	"github.com/apache/incubator-answer/internal/service_quote"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"
)

// QuoteController Quote controller
type QuoteController struct {
	quoteService        *service_quote.QuoteService
	answerService       *content.AnswerService
	rankService         *rank.RankService
	siteInfoService     siteinfo_common.SiteInfoCommonService
	actionService       *action.CaptchaService
	rateLimitMiddleware *middleware.RateLimitMiddleware
}

// NewQuoteController new controller
func NewQuoteController(
	quoteService *service_quote.QuoteService,
	answerService *content.AnswerService,
	rankService *rank.RankService,
	siteInfoService siteinfo_common.SiteInfoCommonService,
	actionService *action.CaptchaService,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
) *QuoteController {
	return &QuoteController{
		quoteService:        quoteService,
		answerService:       answerService,
		rankService:         rankService,
		siteInfoService:     siteInfoService,
		actionService:       actionService,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}

// RemoveQuote delete Quote
// @Summary delete Quote
// @Description delete Quote
// @Tags Quote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.RemoveQuoteReq true "Quote"
// @Success 200 {object} handler.RespBody
// @Router  /answer/api/v1/Quote [delete]
func (qc *QuoteController) RemoveQuote(ctx *gin.Context) {
	req := &schema.RemoveQuoteReq{}
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

	can, err := qc.rankService.CheckOperationPermission(ctx, req.UserID, permission.QuoteDelete, req.ID)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	err = qc.quoteService.RemoveQuote(ctx, req)
	if !isAdmin {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionDelete, req.UserID)
	}
	handler.HandleResponse(ctx, err, nil)
}

// OperationQuote Operation Quote
// @Summary Operation Quote
// @Description Operation Quote \n operation [pin unpin hide show]
// @Tags Quote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.OperationQuoteReq true "Quote"
// @Success 200 {object} handler.RespBody
// @Router  /answer/api/v1/Quote/operation [put]
func (qc *QuoteController) OperationQuote(ctx *gin.Context) {
	req := &schema.OperationQuoteReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.QuotePin,
		permission.QuoteUnPin,
		permission.QuoteHide,
		permission.QuoteShow,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	req.CanPin = canList[0]
	req.CanList = canList[1]
	if (req.Operation == schema.QuoteOperationPin || req.Operation == schema.QuoteOperationUnPin) && !req.CanPin {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	if (req.Operation == schema.QuoteOperationHide || req.Operation == schema.QuoteOperationShow) && !req.CanList {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	err = qc.quoteService.OperationQuote(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// CloseQuote Close Quote
// @Summary Close Quote
// @Description Close Quote
// @Tags Quote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.CloseQuoteReq true "Quote"
// @Success 200 {object} handler.RespBody
// @Router  /answer/api/v1/Quote/status [put]
func (qc *QuoteController) CloseQuote(ctx *gin.Context) {
	req := &schema.CloseQuoteReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := qc.rankService.CheckOperationPermission(ctx, req.UserID, permission.QuoteClose, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = qc.quoteService.CloseQuote(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// ReopenQuote reopen Quote
// @Summary reopen Quote
// @Description reopen Quote
// @Tags Quote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.ReopenQuoteReq true "Quote"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/Quote/reopen [put]
func (qc *QuoteController) ReopenQuote(ctx *gin.Context) {
	req := &schema.ReopenQuoteReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.QuoteID = uid.DeShortID(req.QuoteID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := qc.rankService.CheckOperationPermission(ctx, req.UserID, permission.QuoteReopen, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = qc.quoteService.ReopenQuote(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// GetQuote get Quote details
// @Summary get Quote details
// @Description get Quote details
// @Tags Quote
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id query string true "Quote TagID"  default(1)
// @Success 200 {string} string ""
// @Router /answer/api/v1/Quote/info [get]
func (qc *QuoteController) GetQuote(ctx *gin.Context) {
	id := ctx.Query("id")
	id = uid.DeShortID(id)
	userID := middleware.GetLoginUserIDFromContext(ctx)
	req := schema.QuotePermission{}
	canList, err := qc.rankService.CheckOperationPermissions(ctx, userID, []string{
		permission.QuoteEdit,
		permission.QuoteDelete,
		permission.QuoteClose,
		permission.QuoteReopen,
		permission.QuotePin,
		permission.QuoteUnPin,
		permission.QuoteHide,
		permission.QuoteShow,
		permission.AnswerInviteSomeoneToAnswer,
		permission.QuoteUnDelete,
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

	info, err := qc.quoteService.GetQuoteAndAddPV(ctx, id, userID, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if handler.GetEnableShortID(ctx) {
		info.ID = uid.EnShortID(info.ID)
	}
	handler.HandleResponse(ctx, nil, info)
}

// GetQuoteInviteUserInfo get Quote invite user info
// @Summary get Quote invite user info
// @Description get Quote invite user info
// @Tags Quote
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id query string true "Quote ID"  default(1)
// @Success 200 {string} string ""
// @Router /answer/api/v1/Quote/invite [get]
func (qc *QuoteController) GetQuoteInviteUserInfo(ctx *gin.Context) {
	QuoteID := uid.DeShortID(ctx.Query("id"))
	resp, err := qc.quoteService.InviteUserInfo(ctx, QuoteID)
	handler.HandleResponse(ctx, err, resp)

}

// SimilarQuote godoc
// @Summary Search Similar Quote
// @Description Search Similar Quote
// @Tags Quote
// @Accept  json
// @Produce  json
// @Param Quote_id query string true "Quote_id"  default()
// @Success 200 {string} string ""
// @Router /answer/api/v1/Quote/similar/tag [get]
func (qc *QuoteController) SimilarQuote(ctx *gin.Context) {
	QuoteID := ctx.Query("Quote_id")
	QuoteID = uid.DeShortID(QuoteID)
	userID := middleware.GetLoginUserIDFromContext(ctx)
	list, count, err := qc.quoteService.SimilarQuote(ctx, QuoteID, userID)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, nil, gin.H{
		"list":  list,
		"count": count,
	})
}

// QuotePage get Quotes by page
// @Summary get Quotes by page
// @Description get Quotes by page
// @Tags Quote
// @Accept  json
// @Produce  json
// @Param data body schema.QuotePageReq  true "QuotePageReq"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.QuotePageResp}}
// @Router /answer/api/v1/Quote/page [get]
func (qc *QuoteController) QuotePage(ctx *gin.Context) {
	req := &schema.QuotePageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)

	Quotes, total, err := qc.quoteService.GetQuotePage(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, nil, pager.NewPageModel(total, Quotes))
}

// QuoteRecommendPage get recommend Quotes by page
// @Summary get recommend Quotes by page
// @Description get recommend Quotes by page
// @Tags Quote
// @Accept  json
// @Produce  json
// @Param data body schema.QuotePageReq  true "QuotePageReq"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.QuotePageResp}}
// @Router /answer/api/v1/Quote/recommend/page [get]
func (qc *QuoteController) QuoteRecommendPage(ctx *gin.Context) {
	req := &schema.QuotePageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)

	if req.LoginUserID == "" {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	Quotes, total, err := qc.quoteService.GetRecommendQuotePage(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, nil, pager.NewPageModel(total, Quotes))
}

// AddQuote add Quote
// @Summary add Quote
// @Description add Quote
// @Tags Quote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuoteAdd true "Quote"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/Quote [post]
func (qc *QuoteController) AddQuote(ctx *gin.Context) {
	req := &schema.QuoteAdd{}
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
		permission.QuoteAdd,
		permission.QuoteEdit,
		permission.QuoteDelete,
		permission.QuoteClose,
		permission.QuoteReopen,
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
		captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionQuote, req.UserID, req.CaptchaID, req.CaptchaCode)
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
	hasNewTag, err := qc.quoteService.HasNewTag(ctx, req.Tags)
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

	errList, err := qc.quoteService.CheckAddQuote(ctx, req)
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

	resp, err := qc.quoteService.AddQuote(ctx, req)
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
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionQuote, req.UserID)
	}
	handler.HandleResponse(ctx, err, resp)
}

// AddQuoteByAnswer add Quote
// @Summary add Quote and answer
// @Description add Quote and answer
// @Tags Quote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuoteAddByAnswer true "Quote"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/Quote/answer [post]
func (qc *QuoteController) AddQuoteByAnswer(ctx *gin.Context) {
	req := &schema.QuoteAddByAnswer{}
	errFields := handler.BindAndCheckReturnErr(ctx, req)
	if ctx.IsAborted() {
		return
	}
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.QuoteAdd,
		permission.QuoteEdit,
		permission.QuoteDelete,
		permission.QuoteClose,
		permission.QuoteReopen,
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
		captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionQuote, req.UserID, req.CaptchaID, req.CaptchaCode)
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
	QuoteReq := new(schema.QuoteAdd)
	err = copier.Copy(QuoteReq, req)
	if err != nil {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RequestFormatError), nil)
		return
	}
	errList, err := qc.quoteService.CheckAddQuote(ctx, QuoteReq)
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
	resp, err := qc.quoteService.AddQuote(ctx, QuoteReq)
	if err != nil {
		errlist, ok := resp.([]*validator.FormErrorField)
		if ok {
			errFields = append(errFields, errlist...)
		}
	}

	if !isAdmin || !linkUrlLimitUser {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionQuote, req.UserID)
	}

	if len(errFields) > 0 {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		return
	}
	////add the Quote id to the answer
	//QuoteInfo, ok := resp.(*schema.QuoteInfoResp)
	//if ok {
	//	answerReq := &schema.AnswerAddReq{}
	//	answerReq.QuoteID = uid.DeShortID(QuoteInfo.ID)
	//	answerReq.UserID = middleware.GetLoginUserIDFromContext(ctx)
	//	answerReq.Content = req.AnswerContent
	//	answerReq.HTML = req.AnswerHTML
	//	answerID, err := qc.answerService.Insert(ctx, answerReq)
	//	if err != nil {
	//		handler.HandleResponse(ctx, err, nil)
	//		return
	//	}
	//	info, QuoteInfo, has, err := qc.answerService.Get(ctx, answerID, req.UserID)
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
	//		"Quote": QuoteInfo,
	//	})
	//	return
	//}

	handler.HandleResponse(ctx, err, resp)
}

// UpdateQuote update Quote
// @Summary update Quote
// @Description update Quote
// @Tags Quote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuoteUpdate true "Quote"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/Quote [put]
func (qc *QuoteController) UpdateQuote(ctx *gin.Context) {
	req := &schema.QuoteUpdate{}
	errFields := handler.BindAndCheckReturnErr(ctx, req)
	if ctx.IsAborted() {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, requireRanks, err := qc.rankService.CheckOperationPermissionsForRanks(ctx, req.UserID, []string{
		permission.QuoteEdit,
		permission.QuoteDelete,
		permission.QuoteEditWithoutReview,
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

	errlist, err := qc.quoteService.UpdateQuoteCheckTags(ctx, req)
	if err != nil {
		errFields = append(errFields, errlist...)
	}

	if len(errFields) > 0 {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		return
	}

	// can add tag
	hasNewTag, err := qc.quoteService.HasNewTag(ctx, req.Tags)
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

	resp, err := qc.quoteService.UpdateQuote(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, resp)
		return
	}
	respInfo, ok := resp.(*schema.QuoteInfoResp)
	if !ok {
		handler.HandleResponse(ctx, err, resp)
		return
	}
	if !isAdmin || !linkUrlLimitUser {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionEdit, req.UserID)
	}
	handler.HandleResponse(ctx, nil, &schema.UpdateQuoteResp{UrlTitle: respInfo.UrlTitle, WaitForReview: !req.NoNeedReview})
}

// QuoteRecover recover deleted Quote
// @Summary recover deleted Quote
// @Description recover deleted Quote
// @Tags Quote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuoteRecoverReq true "Quote"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/Quote/recover [post]
func (qc *QuoteController) QuoteRecover(ctx *gin.Context) {
	req := &schema.QuoteRecoverReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.QuoteID = uid.DeShortID(req.QuoteID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.QuoteUnDelete,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !canList[0] {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = qc.quoteService.RecoverQuote(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// UpdateQuoteInviteUser update Quote invite user
// @Summary update Quote invite user
// @Description update Quote invite user
// @Tags Quote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuoteUpdateInviteUser true "Quote"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/Quote/invite [put]
func (qc *QuoteController) UpdateQuoteInviteUser(ctx *gin.Context) {
	req := &schema.QuoteUpdateInviteUser{}
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
	err = qc.quoteService.UpdateQuoteInviteUser(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !isAdmin {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionInvitationAnswer, req.UserID)
	}
	handler.HandleResponse(ctx, nil, nil)
}

// GetSimilarQuotes fuzzy query similar Quotes based on title
// @Summary fuzzy query similar Quotes based on title
// @Description fuzzy query similar Quotes based on title
// @Tags Quote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param title query string true "title"  default(string)
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/Quote/similar [get]
func (qc *QuoteController) GetSimilarQuotes(ctx *gin.Context) {
	title := ctx.Query("title")
	resp, err := qc.quoteService.GetQuotesByTitle(ctx, title)
	handler.HandleResponse(ctx, err, resp)
}

// UserTop godoc
// @Summary UserTop
// @Description UserTop
// @Tags Quote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username query string true "username"  default(string)
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/personal/qa/top [get]
func (qc *QuoteController) UserTop(ctx *gin.Context) {
	userName := ctx.Query("username")
	userID := middleware.GetLoginUserIDFromContext(ctx)
	QuoteList, answerList, err := qc.quoteService.SearchUserTopList(ctx, userName, userID)
	handler.HandleResponse(ctx, err, gin.H{
		"Quote":  QuoteList,
		"answer": answerList,
	})
}

// PersonalQuotePage list personal Quotes
// @Summary list personal Quotes
// @Description list personal Quotes
// @Tags Personal
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username query string true "username"  default(string)
// @Param order query string true "order"  Enums(newest,score)
// @Param page query string true "page"  default(0)
// @Param page_size query string true "page_size" default(20)
// @Success 200 {object} handler.RespBody
// @Router /personal/Quote/page [get]
func (qc *QuoteController) PersonalQuotePage(ctx *gin.Context) {
	req := &schema.PersonalQuotePageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	req.IsAdmin = middleware.GetUserIsAdminModerator(ctx)
	resp, err := qc.quoteService.PersonalQuotePage(ctx, req)
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
func (qc *QuoteController) PersonalAnswerPage(ctx *gin.Context) {
	req := &schema.PersonalAnswerPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	req.IsAdmin = middleware.GetUserIsAdminModerator(ctx)
	resp, err := qc.quoteService.PersonalAnswerPage(ctx, req)
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
func (qc *QuoteController) PersonalCollectionPage(ctx *gin.Context) {
	req := &schema.PersonalCollectionPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	resp, err := qc.quoteService.PersonalCollectionPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// AdminQuotePage admin Quote page
// @Summary AdminQuotePage admin Quote page
// @Description Status:[available,closed,deleted,pending]
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "page size"
// @Param page_size query int false "page size"
// @Param status query string false "user status" Enums(available, closed, deleted, pending)
// @Param query query string false "Quote id or title"
// @Success 200 {object} handler.RespBody
// @Router /answer/admin/api/Quote/page [get]
func (qc *QuoteController) AdminQuotePage(ctx *gin.Context) {
	req := &schema.AdminQuotePageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	resp, err := qc.quoteService.AdminQuotePage(ctx, req)
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
// @Param query query string false "answer id or Quote title"
// @Param Quote_id query string false "Quote id"
// @Success 200 {object} handler.RespBody
// @Router /answer/admin/api/answer/page [get]
func (qc *QuoteController) AdminAnswerPage(ctx *gin.Context) {
	req := &schema.AdminAnswerPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	resp, err := qc.quoteService.AdminAnswerPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// AdminUpdateQuoteStatus update Quote status
// @Summary update Quote status
// @Description update Quote status
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.AdminUpdateQuoteStatusReq true "AdminUpdateQuoteStatusReq"
// @Success 200 {object} handler.RespBody
// @Router /answer/admin/api/Quote/status [put]
func (qc *QuoteController) AdminUpdateQuoteStatus(ctx *gin.Context) {
	req := &schema.AdminUpdateQuoteStatusReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.QuoteID = uid.DeShortID(req.QuoteID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	err := qc.quoteService.AdminSetQuoteStatus(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}
