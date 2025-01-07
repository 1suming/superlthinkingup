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

// QuotePieceController QuotePiece controller
type QuotePieceController struct {
	QuotePieceService   *service_quote.QuotePieceService
	answerService       *content.AnswerService
	rankService         *rank.RankService
	siteInfoService     siteinfo_common.SiteInfoCommonService
	actionService       *action.CaptchaService
	rateLimitMiddleware *middleware.RateLimitMiddleware
}

// NewQuotePieceController new controller
func NewQuotePieceController(
	QuotePieceService *service_quote.QuotePieceService,
	answerService *content.AnswerService,
	rankService *rank.RankService,
	siteInfoService siteinfo_common.SiteInfoCommonService,
	actionService *action.CaptchaService,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
) *QuotePieceController {
	return &QuotePieceController{
		QuotePieceService:   QuotePieceService,
		answerService:       answerService,
		rankService:         rankService,
		siteInfoService:     siteInfoService,
		actionService:       actionService,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}

// RemoveQuotePiece delete QuotePiece
// @Summary delete QuotePiece
// @Description delete QuotePiece
// @Tags QuotePiece
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.RemoveQuotePieceReq true "QuotePiece"
// @Success 200 {object} handler.RespBody
// @Router  /answer/api/v1/QuotePiece [delete]
func (qc *QuotePieceController) RemoveQuotePiece(ctx *gin.Context) {
	req := &schema.RemoveQuotePieceReq{}
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

	can, err := qc.rankService.CheckOperationPermission(ctx, req.UserID, permission.QuotePieceDelete, req.ID)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	err = qc.QuotePieceService.RemoveQuotePiece(ctx, req)
	if !isAdmin {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionDelete, req.UserID)
	}
	handler.HandleResponse(ctx, err, nil)
}

// OperationQuotePiece Operation QuotePiece
// @Summary Operation QuotePiece
// @Description Operation QuotePiece \n operation [pin unpin hide show]
// @Tags QuotePiece
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.OperationQuotePieceReq true "QuotePiece"
// @Success 200 {object} handler.RespBody
// @Router  /answer/api/v1/QuotePiece/operation [put]
func (qc *QuotePieceController) OperationQuotePiece(ctx *gin.Context) {
	req := &schema.OperationQuotePieceReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.QuotePiecePin,
		permission.QuotePieceUnPin,
		permission.QuotePieceHide,
		permission.QuotePieceShow,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	req.CanPin = canList[0]
	req.CanList = canList[1]
	if (req.Operation == schema.QuotePieceOperationPin || req.Operation == schema.QuotePieceOperationUnPin) && !req.CanPin {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	if (req.Operation == schema.QuotePieceOperationHide || req.Operation == schema.QuotePieceOperationShow) && !req.CanList {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	err = qc.QuotePieceService.OperationQuotePiece(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// CloseQuotePiece Close QuotePiece
// @Summary Close QuotePiece
// @Description Close QuotePiece
// @Tags QuotePiece
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.CloseQuotePieceReq true "QuotePiece"
// @Success 200 {object} handler.RespBody
// @Router  /answer/api/v1/QuotePiece/status [put]
func (qc *QuotePieceController) CloseQuotePiece(ctx *gin.Context) {
	req := &schema.CloseQuotePieceReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := qc.rankService.CheckOperationPermission(ctx, req.UserID, permission.QuotePieceClose, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = qc.QuotePieceService.CloseQuotePiece(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// ReopenQuotePiece reopen QuotePiece
// @Summary reopen QuotePiece
// @Description reopen QuotePiece
// @Tags QuotePiece
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.ReopenQuotePieceReq true "QuotePiece"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuotePiece/reopen [put]
func (qc *QuotePieceController) ReopenQuotePiece(ctx *gin.Context) {
	req := &schema.ReopenQuotePieceReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.QuotePieceID = uid.DeShortID(req.QuotePieceID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := qc.rankService.CheckOperationPermission(ctx, req.UserID, permission.QuotePieceReopen, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = qc.QuotePieceService.ReopenQuotePiece(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// GetQuotePiece get QuotePiece details
// @Summary get QuotePiece details
// @Description get QuotePiece details
// @Tags QuotePiece
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id query string true "QuotePiece TagID"  default(1)
// @Success 200 {string} string ""
// @Router /answer/api/v1/QuotePiece/info [get]
func (qc *QuotePieceController) GetQuotePiece(ctx *gin.Context) {
	id := ctx.Query("id")
	id = uid.DeShortID(id)
	userID := middleware.GetLoginUserIDFromContext(ctx)
	req := schema.QuotePiecePermission{}
	canList, err := qc.rankService.CheckOperationPermissions(ctx, userID, []string{
		permission.QuotePieceEdit,
		permission.QuotePieceDelete,
		permission.QuotePieceClose,
		permission.QuotePieceReopen,
		permission.QuotePiecePin,
		permission.QuotePieceUnPin,
		permission.QuotePieceHide,
		permission.QuotePieceShow,
		permission.AnswerInviteSomeoneToAnswer,
		permission.QuotePieceUnDelete,
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

	info, err := qc.QuotePieceService.GetQuotePieceAndAddPV(ctx, id, userID, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if handler.GetEnableShortID(ctx) {
		info.ID = uid.EnShortID(info.ID)
	}
	handler.HandleResponse(ctx, nil, info)
}

// GetQuotePieceInviteUserInfo get QuotePiece invite user info
// @Summary get QuotePiece invite user info
// @Description get QuotePiece invite user info
// @Tags QuotePiece
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id query string true "QuotePiece ID"  default(1)
// @Success 200 {string} string ""
// @Router /answer/api/v1/QuotePiece/invite [get]
func (qc *QuotePieceController) GetQuotePieceInviteUserInfo(ctx *gin.Context) {
	QuotePieceID := uid.DeShortID(ctx.Query("id"))
	resp, err := qc.QuotePieceService.InviteUserInfo(ctx, QuotePieceID)
	handler.HandleResponse(ctx, err, resp)

}

// SimilarQuotePiece godoc
// @Summary Search Similar QuotePiece
// @Description Search Similar QuotePiece
// @Tags QuotePiece
// @Accept  json
// @Produce  json
// @Param QuotePiece_id query string true "QuotePiece_id"  default()
// @Success 200 {string} string ""
// @Router /answer/api/v1/QuotePiece/similar/tag [get]
func (qc *QuotePieceController) SimilarQuotePiece(ctx *gin.Context) {
	QuotePieceID := ctx.Query("QuotePiece_id")
	QuotePieceID = uid.DeShortID(QuotePieceID)
	userID := middleware.GetLoginUserIDFromContext(ctx)
	list, count, err := qc.QuotePieceService.SimilarQuotePiece(ctx, QuotePieceID, userID)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, nil, gin.H{
		"list":  list,
		"count": count,
	})
}

// QuotePiecePage get QuotePieces by page
// @Summary get QuotePieces by page
// @Description get QuotePieces by page
// @Tags QuotePiece
// @Accept  json
// @Produce  json
// @Param data body schema.QuotePiecePageReq  true "QuotePiecePageReq"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.QuotePiecePageResp}}
// @Router /answer/api/v1/QuotePiece/page [get]
func (qc *QuotePieceController) QuotePiecePage(ctx *gin.Context) {
	req := &schema.QuotePiecePageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)

	QuotePieces, total, err := qc.QuotePieceService.GetQuotePiecePage(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, nil, pager.NewPageModel(total, QuotePieces))
}

// QuotePieceRecommendPage get recommend QuotePieces by page
// @Summary get recommend QuotePieces by page
// @Description get recommend QuotePieces by page
// @Tags QuotePiece
// @Accept  json
// @Produce  json
// @Param data body schema.QuotePiecePageReq  true "QuotePiecePageReq"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.QuotePiecePageResp}}
// @Router /answer/api/v1/QuotePiece/recommend/page [get]
func (qc *QuotePieceController) QuotePieceRecommendPage(ctx *gin.Context) {
	req := &schema.QuotePiecePageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)

	if req.LoginUserID == "" {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	QuotePieces, total, err := qc.QuotePieceService.GetRecommendQuotePiecePage(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, nil, pager.NewPageModel(total, QuotePieces))
}

// AddQuotePiece add QuotePiece
// @Summary add QuotePiece
// @Description add QuotePiece
// @Tags QuotePiece
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuotePieceAdd true "QuotePiece"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuotePiece [post]
func (qc *QuotePieceController) AddQuotePiece(ctx *gin.Context) {
	req := &schema.QuotePieceAdd{}
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
		permission.QuotePieceAdd,
		permission.QuotePieceEdit,
		permission.QuotePieceDelete,
		permission.QuotePieceClose,
		permission.QuotePieceReopen,
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
		captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionQuotePiece, req.UserID, req.CaptchaID, req.CaptchaCode)
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
	hasNewTag, err := qc.QuotePieceService.HasNewTag(ctx, req.Tags)
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

	errList, err := qc.QuotePieceService.CheckAddQuotePiece(ctx, req)
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

	resp, err := qc.QuotePieceService.AddQuotePiece(ctx, req)
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
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionQuotePiece, req.UserID)
	}
	handler.HandleResponse(ctx, err, resp)
}

// AddQuotePieceByAnswer add QuotePiece
// @Summary add QuotePiece and answer
// @Description add QuotePiece and answer
// @Tags QuotePiece
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuotePieceAddByAnswer true "QuotePiece"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuotePiece/answer [post]
func (qc *QuotePieceController) AddQuotePieceByAnswer(ctx *gin.Context) {
	req := &schema.QuotePieceAddByAnswer{}
	errFields := handler.BindAndCheckReturnErr(ctx, req)
	if ctx.IsAborted() {
		return
	}
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.QuotePieceAdd,
		permission.QuotePieceEdit,
		permission.QuotePieceDelete,
		permission.QuotePieceClose,
		permission.QuotePieceReopen,
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
		captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionQuotePiece, req.UserID, req.CaptchaID, req.CaptchaCode)
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
	QuotePieceReq := new(schema.QuotePieceAdd)
	err = copier.Copy(QuotePieceReq, req)
	if err != nil {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RequestFormatError), nil)
		return
	}
	errList, err := qc.QuotePieceService.CheckAddQuotePiece(ctx, QuotePieceReq)
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
	resp, err := qc.QuotePieceService.AddQuotePiece(ctx, QuotePieceReq)
	if err != nil {
		errlist, ok := resp.([]*validator.FormErrorField)
		if ok {
			errFields = append(errFields, errlist...)
		}
	}

	if !isAdmin || !linkUrlLimitUser {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionQuotePiece, req.UserID)
	}

	if len(errFields) > 0 {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		return
	}
	////add the QuotePiece id to the answer
	//QuotePieceInfo, ok := resp.(*schema.QuotePieceInfoResp)
	//if ok {
	//	answerReq := &schema.AnswerAddReq{}
	//	answerReq.QuotePieceID = uid.DeShortID(QuotePieceInfo.ID)
	//	answerReq.UserID = middleware.GetLoginUserIDFromContext(ctx)
	//	answerReq.Content = req.AnswerContent
	//	answerReq.HTML = req.AnswerHTML
	//	answerID, err := qc.answerService.Insert(ctx, answerReq)
	//	if err != nil {
	//		handler.HandleResponse(ctx, err, nil)
	//		return
	//	}
	//	info, QuotePieceInfo, has, err := qc.answerService.Get(ctx, answerID, req.UserID)
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
	//		"QuotePiece": QuotePieceInfo,
	//	})
	//	return
	//}

	handler.HandleResponse(ctx, err, resp)
}

// UpdateQuotePiece update QuotePiece
// @Summary update QuotePiece
// @Description update QuotePiece
// @Tags QuotePiece
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuotePieceUpdate true "QuotePiece"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuotePiece [put]
func (qc *QuotePieceController) UpdateQuotePiece(ctx *gin.Context) {
	req := &schema.QuotePieceUpdate{}
	errFields := handler.BindAndCheckReturnErr(ctx, req)
	if ctx.IsAborted() {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, requireRanks, err := qc.rankService.CheckOperationPermissionsForRanks(ctx, req.UserID, []string{
		permission.QuotePieceEdit,
		permission.QuotePieceDelete,
		permission.QuotePieceEditWithoutReview,
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

	errlist, err := qc.QuotePieceService.UpdateQuotePieceCheckTags(ctx, req)
	if err != nil {
		errFields = append(errFields, errlist...)
	}

	if len(errFields) > 0 {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		return
	}

	// can add tag
	hasNewTag, err := qc.QuotePieceService.HasNewTag(ctx, req.Tags)
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

	resp, err := qc.QuotePieceService.UpdateQuotePiece(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, resp)
		return
	}
	respInfo, ok := resp.(*schema.QuotePieceInfoResp)
	if !ok {
		handler.HandleResponse(ctx, err, resp)
		return
	}
	if !isAdmin || !linkUrlLimitUser {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionEdit, req.UserID)
	}
	handler.HandleResponse(ctx, nil, &schema.UpdateQuotePieceResp{UrlTitle: respInfo.UrlTitle, WaitForReview: !req.NoNeedReview})
}

// QuotePieceRecover recover deleted QuotePiece
// @Summary recover deleted QuotePiece
// @Description recover deleted QuotePiece
// @Tags QuotePiece
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuotePieceRecoverReq true "QuotePiece"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuotePiece/recover [post]
func (qc *QuotePieceController) QuotePieceRecover(ctx *gin.Context) {
	req := &schema.QuotePieceRecoverReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.QuotePieceID = uid.DeShortID(req.QuotePieceID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.QuotePieceUnDelete,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !canList[0] {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = qc.QuotePieceService.RecoverQuotePiece(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// UpdateQuotePieceInviteUser update QuotePiece invite user
// @Summary update QuotePiece invite user
// @Description update QuotePiece invite user
// @Tags QuotePiece
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.QuotePieceUpdateInviteUser true "QuotePiece"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuotePiece/invite [put]
func (qc *QuotePieceController) UpdateQuotePieceInviteUser(ctx *gin.Context) {
	req := &schema.QuotePieceUpdateInviteUser{}
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
	err = qc.QuotePieceService.UpdateQuotePieceInviteUser(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !isAdmin {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionInvitationAnswer, req.UserID)
	}
	handler.HandleResponse(ctx, nil, nil)
}

// GetSimilarQuotePieces fuzzy query similar QuotePieces based on title
// @Summary fuzzy query similar QuotePieces based on title
// @Description fuzzy query similar QuotePieces based on title
// @Tags QuotePiece
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param title query string true "title"  default(string)
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/QuotePiece/similar [get]
func (qc *QuotePieceController) GetSimilarQuotePieces(ctx *gin.Context) {
	title := ctx.Query("title")
	resp, err := qc.QuotePieceService.GetQuotePiecesByTitle(ctx, title)
	handler.HandleResponse(ctx, err, resp)
}

// UserTop godoc
// @Summary UserTop
// @Description UserTop
// @Tags QuotePiece
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username query string true "username"  default(string)
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/personal/qa/top [get]
func (qc *QuotePieceController) UserTop(ctx *gin.Context) {
	userName := ctx.Query("username")
	userID := middleware.GetLoginUserIDFromContext(ctx)
	QuotePieceList, answerList, err := qc.QuotePieceService.SearchUserTopList(ctx, userName, userID)
	handler.HandleResponse(ctx, err, gin.H{
		"QuotePiece": QuotePieceList,
		"answer":     answerList,
	})
}

// PersonalQuotePiecePage list personal QuotePieces
// @Summary list personal QuotePieces
// @Description list personal QuotePieces
// @Tags Personal
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username query string true "username"  default(string)
// @Param order query string true "order"  Enums(newest,score)
// @Param page query string true "page"  default(0)
// @Param page_size query string true "page_size" default(20)
// @Success 200 {object} handler.RespBody
// @Router /personal/QuotePiece/page [get]
func (qc *QuotePieceController) PersonalQuotePiecePage(ctx *gin.Context) {
	req := &schema.PersonalQuotePiecePageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	req.IsAdmin = middleware.GetUserIsAdminModerator(ctx)
	resp, err := qc.QuotePieceService.PersonalQuotePiecePage(ctx, req)
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
func (qc *QuotePieceController) PersonalAnswerPage(ctx *gin.Context) {
	req := &schema.PersonalAnswerPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	req.IsAdmin = middleware.GetUserIsAdminModerator(ctx)
	resp, err := qc.QuotePieceService.PersonalAnswerPage(ctx, req)
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
func (qc *QuotePieceController) PersonalCollectionPage(ctx *gin.Context) {
	req := &schema.PersonalCollectionPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	resp, err := qc.QuotePieceService.PersonalCollectionPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// AdminQuotePiecePage admin QuotePiece page
// @Summary AdminQuotePiecePage admin QuotePiece page
// @Description Status:[available,closed,deleted,pending]
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "page size"
// @Param page_size query int false "page size"
// @Param status query string false "user status" Enums(available, closed, deleted, pending)
// @Param query query string false "QuotePiece id or title"
// @Success 200 {object} handler.RespBody
// @Router /answer/admin/api/QuotePiece/page [get]
func (qc *QuotePieceController) AdminQuotePiecePage(ctx *gin.Context) {
	req := &schema.AdminQuotePiecePageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	resp, err := qc.QuotePieceService.AdminQuotePiecePage(ctx, req)
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
// @Param query query string false "answer id or QuotePiece title"
// @Param QuotePiece_id query string false "QuotePiece id"
// @Success 200 {object} handler.RespBody
// @Router /answer/admin/api/answer/page [get]
func (qc *QuotePieceController) AdminAnswerPage(ctx *gin.Context) {
	req := &schema.AdminAnswerPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	resp, err := qc.QuotePieceService.AdminAnswerPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// AdminUpdateQuotePieceStatus update QuotePiece status
// @Summary update QuotePiece status
// @Description update QuotePiece status
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.AdminUpdateQuotePieceStatusReq true "AdminUpdateQuotePieceStatusReq"
// @Success 200 {object} handler.RespBody
// @Router /answer/admin/api/QuotePiece/status [put]
func (qc *QuotePieceController) AdminUpdateQuotePieceStatus(ctx *gin.Context) {
	req := &schema.AdminUpdateQuotePieceStatusReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.QuotePieceID = uid.DeShortID(req.QuotePieceID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	err := qc.QuotePieceService.AdminSetQuotePieceStatus(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}
