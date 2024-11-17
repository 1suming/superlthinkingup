package controller_article

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
	"github.com/apache/incubator-answer/internal/service_article"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
)

//@csw

// ArticleController article controller
type ArticleController struct {
	articleService      *service_article.ArticleService
	answerService       *content.AnswerService
	rankService         *rank.RankService
	siteInfoService     siteinfo_common.SiteInfoCommonService
	actionService       *action.CaptchaService
	rateLimitMiddleware *middleware.RateLimitMiddleware
}

// NewArticleController new controller
func NewArticleController(
	articleService *service_article.ArticleService,
	answerService *content.AnswerService,
	rankService *rank.RankService,
	siteInfoService siteinfo_common.SiteInfoCommonService,
	actionService *action.CaptchaService,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
) *ArticleController {
	return &ArticleController{
		articleService:      articleService,
		answerService:       answerService,
		rankService:         rankService,
		siteInfoService:     siteInfoService,
		actionService:       actionService,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}

// RemoveArticle delete article
// @Summary delete article
// @Description delete article
// @Tags Article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.RemoveArticleReq true "article"
// @Success 200 {object} handler.RespBody
// @Router  /answer/api/v1/article [delete]
func (qc *ArticleController) RemoveArticle(ctx *gin.Context) {
	req := &schema.RemoveArticleReq{}
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

	can, err := qc.rankService.CheckOperationPermission(ctx, req.UserID, permission.ArticleDelete, req.ID)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	err = qc.articleService.RemoveArticle(ctx, req)
	if !isAdmin {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionDelete, req.UserID)
	}
	handler.HandleResponse(ctx, err, nil)
}

// OperationArticle Operation article
// @Summary Operation article
// @Description Operation article \n operation [pin unpin hide show]
// @Tags Article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.OperationArticleReq true "article"
// @Success 200 {object} handler.RespBody
// @Router  /answer/api/v1/article/operation [put]
func (qc *ArticleController) OperationArticle(ctx *gin.Context) {
	req := &schema.OperationArticleReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.ArticlePin,
		permission.ArticleUnPin,
		permission.ArticleHide,
		permission.ArticleShow,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	req.CanPin = canList[0]
	req.CanList = canList[1]
	if (req.Operation == schema.ArticleOperationPin || req.Operation == schema.ArticleOperationUnPin) && !req.CanPin {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	if (req.Operation == schema.ArticleOperationHide || req.Operation == schema.ArticleOperationShow) && !req.CanList {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	err = qc.articleService.OperationArticle(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// CloseArticle Close article
// @Summary Close article
// @Description Close article
// @Tags Article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.CloseArticleReq true "article"
// @Success 200 {object} handler.RespBody
// @Router  /answer/api/v1/article/status [put]
func (qc *ArticleController) CloseArticle(ctx *gin.Context) {
	req := &schema.CloseArticleReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := qc.rankService.CheckOperationPermission(ctx, req.UserID, permission.ArticleClose, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = qc.articleService.CloseArticle(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// ReopenArticle reopen article
// @Summary reopen article
// @Description reopen article
// @Tags Article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.ReopenArticleReq true "article"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/article/reopen [put]
func (qc *ArticleController) ReopenArticle(ctx *gin.Context) {
	req := &schema.ReopenArticleReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.ArticleID = uid.DeShortID(req.ArticleID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := qc.rankService.CheckOperationPermission(ctx, req.UserID, permission.ArticleReopen, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = qc.articleService.ReopenArticle(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// GetArticle get article details
// @Summary get article details
// @Description get article details
// @Tags Article
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id query string true "Article TagID"  default(1)
// @Success 200 {string} string ""
// @Router /answer/api/v1/article/info [get]
func (qc *ArticleController) GetArticle(ctx *gin.Context) {
	id := ctx.Query("id")
	log.Info("GetArticle origin:id:%v", id)
	id = uid.DeShortID(id)
	log.Info("GetArticle  after id:%v", id)

	userID := middleware.GetLoginUserIDFromContext(ctx)
	req := schema.ArticlePermission{}
	canList, err := qc.rankService.CheckOperationPermissions(ctx, userID, []string{
		permission.ArticleEdit,
		permission.ArticleDelete,
		permission.ArticleClose,
		permission.ArticleReopen,
		permission.ArticlePin,
		permission.ArticleUnPin,
		permission.ArticleHide,
		permission.ArticleShow,
		permission.AnswerInviteSomeoneToAnswer,
		permission.ArticleUnDelete,
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

	info, err := qc.articleService.GetArticleAndAddPV(ctx, id, userID, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if handler.GetEnableShortID(ctx) {
		info.ID = uid.EnShortID(info.ID)
	}
	handler.HandleResponse(ctx, nil, info)
}

// GetArticleInviteUserInfo get article invite user info
// @Summary get article invite user info
// @Description get article invite user info
// @Tags Article
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id query string true "Article ID"  default(1)
// @Success 200 {string} string ""
// @Router /answer/api/v1/article/invite [get]
func (qc *ArticleController) GetArticleInviteUserInfo(ctx *gin.Context) {
	articleID := uid.DeShortID(ctx.Query("id"))
	resp, err := qc.articleService.InviteUserInfo(ctx, articleID)
	handler.HandleResponse(ctx, err, resp)

}

// SimilarArticle godoc
// @Summary Search Similar Article
// @Description Search Similar Article
// @Tags Article
// @Accept  json
// @Produce  json
// @Param article_id query string true "article_id"  default()
// @Success 200 {string} string ""
// @Router /answer/api/v1/article/similar/tag [get]
func (qc *ArticleController) SimilarArticle(ctx *gin.Context) {
	articleID := ctx.Query("article_id")
	articleID = uid.DeShortID(articleID)
	userID := middleware.GetLoginUserIDFromContext(ctx)
	list, count, err := qc.articleService.SimilarArticle(ctx, articleID, userID)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, nil, gin.H{
		"list":  list,
		"count": count,
	})
}

// ArticlePage get articles by page
// @Summary get articles by page
// @Description get articles by page
// @Tags Article
// @Accept  json
// @Produce  json
// @Param data body schema.ArticlePageReq  true "ArticlePageReq"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.ArticlePageResp}}
// @Router /answer/api/v1/article/page [get]
func (qc *ArticleController) ArticlePage(ctx *gin.Context) {
	req := &schema.ArticlePageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)

	articles, total, err := qc.articleService.GetArticlePage(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, nil, pager.NewPageModel(total, articles))
}

// ArticleRecommendPage get recommend articles by page
// @Summary get recommend articles by page
// @Description get recommend articles by page
// @Tags Article
// @Accept  json
// @Produce  json
// @Param data body schema.ArticlePageReq  true "ArticlePageReq"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.ArticlePageResp}}
// @Router /answer/api/v1/article/recommend/page [get]
func (qc *ArticleController) ArticleRecommendPage(ctx *gin.Context) {
	req := &schema.ArticlePageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)

	if req.LoginUserID == "" {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	articles, total, err := qc.articleService.GetRecommendArticlePage(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, nil, pager.NewPageModel(total, articles))
}

// AddArticle add article
// @Summary add article
// @Description add article
// @Tags Article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.ArticleAdd true "article"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/article [post]
func (qc *ArticleController) AddArticle(ctx *gin.Context) {
	req := &schema.ArticleAdd{}
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
		permission.ArticleAdd,
		permission.ArticleEdit,
		permission.ArticleDelete,
		permission.ArticleClose,
		permission.ArticleReopen,
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
		captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionArticle, req.UserID, req.CaptchaID, req.CaptchaCode)
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
	hasNewTag, err := qc.articleService.HasNewTag(ctx, req.Tags)
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

	errList, err := qc.articleService.CheckAddArticle(ctx, req)
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

	resp, err := qc.articleService.AddArticle(ctx, req)
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
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionArticle, req.UserID)
	}
	handler.HandleResponse(ctx, err, resp)
}

// AddArticleByAnswer add article
// @Summary add article and answer
// @Description add article and answer
// @Tags Article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.ArticleAddByAnswer true "article"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/article/answer [post]
func (qc *ArticleController) AddArticleByAnswer(ctx *gin.Context) {
	req := &schema.ArticleAddByAnswer{}
	errFields := handler.BindAndCheckReturnErr(ctx, req)
	if ctx.IsAborted() {
		return
	}
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.ArticleAdd,
		permission.ArticleEdit,
		permission.ArticleDelete,
		permission.ArticleClose,
		permission.ArticleReopen,
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
		captchaPass := qc.actionService.ActionRecordVerifyCaptcha(ctx, entity.CaptchaActionArticle, req.UserID, req.CaptchaID, req.CaptchaCode)
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
	articleReq := new(schema.ArticleAdd)
	err = copier.Copy(articleReq, req)
	if err != nil {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RequestFormatError), nil)
		return
	}
	errList, err := qc.articleService.CheckAddArticle(ctx, articleReq)
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
	resp, err := qc.articleService.AddArticle(ctx, articleReq)
	if err != nil {
		errlist, ok := resp.([]*validator.FormErrorField)
		if ok {
			errFields = append(errFields, errlist...)
		}
	}

	if !isAdmin || !linkUrlLimitUser {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionArticle, req.UserID)
	}

	if len(errFields) > 0 {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		return
	}
	////add the article id to the answer
	//articleInfo, ok := resp.(*schema.ArticleInfoResp)
	//if ok {
	//	answerReq := &schema.AnswerAddReq{}
	//	answerReq.ArticleID = uid.DeShortID(articleInfo.ID)
	//	answerReq.UserID = middleware.GetLoginUserIDFromContext(ctx)
	//	answerReq.Content = req.AnswerContent
	//	answerReq.HTML = req.AnswerHTML
	//	answerID, err := qc.answerService.Insert(ctx, answerReq)
	//	if err != nil {
	//		handler.HandleResponse(ctx, err, nil)
	//		return
	//	}
	//	info, articleInfo, has, err := qc.answerService.Get(ctx, answerID, req.UserID)
	//	if err != nil {
	//		handler.HandleResponse(ctx, err, nil)
	//		return
	//	}
	//	if !has {
	//		handler.HandleResponse(ctx, nil, nil)
	//		return
	//	}
	//	handler.HandleResponse(ctx, err, gin.H{
	//		"info":    info,
	//		"article": articleInfo,
	//	})
	//	return
	//}

	handler.HandleResponse(ctx, err, resp)
}

// UpdateArticle update article
// @Summary update article
// @Description update article
// @Tags Article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.ArticleUpdate true "article"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/article [put]
func (qc *ArticleController) UpdateArticle(ctx *gin.Context) {
	req := &schema.ArticleUpdate{}
	errFields := handler.BindAndCheckReturnErr(ctx, req)
	if ctx.IsAborted() {
		return
	}
	req.ID = uid.DeShortID(req.ID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, requireRanks, err := qc.rankService.CheckOperationPermissionsForRanks(ctx, req.UserID, []string{
		permission.ArticleEdit,
		permission.ArticleDelete,
		permission.ArticleEditWithoutReview,
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

	errlist, err := qc.articleService.UpdateArticleCheckTags(ctx, req)
	if err != nil {
		errFields = append(errFields, errlist...)
	}

	if len(errFields) > 0 {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), errFields)
		return
	}

	// can add tag
	hasNewTag, err := qc.articleService.HasNewTag(ctx, req.Tags)
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

	resp, err := qc.articleService.UpdateArticle(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, resp)
		return
	}
	respInfo, ok := resp.(*schema.ArticleInfoResp)
	if !ok {
		handler.HandleResponse(ctx, err, resp)
		return
	}
	if !isAdmin || !linkUrlLimitUser {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionEdit, req.UserID)
	}
	handler.HandleResponse(ctx, nil, &schema.UpdateArticleResp{UrlTitle: respInfo.UrlTitle, WaitForReview: !req.NoNeedReview})
}

// ArticleRecover recover deleted article
// @Summary recover deleted article
// @Description recover deleted article
// @Tags Article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.ArticleRecoverReq true "article"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/article/recover [post]
func (qc *ArticleController) ArticleRecover(ctx *gin.Context) {
	req := &schema.ArticleRecoverReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.ArticleID = uid.DeShortID(req.ArticleID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	canList, err := qc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		permission.ArticleUnDelete,
	})
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !canList[0] {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = qc.articleService.RecoverArticle(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// UpdateArticleInviteUser update article invite user
// @Summary update article invite user
// @Description update article invite user
// @Tags Article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.ArticleUpdateInviteUser true "article"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/article/invite [put]
func (qc *ArticleController) UpdateArticleInviteUser(ctx *gin.Context) {
	req := &schema.ArticleUpdateInviteUser{}
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
	err = qc.articleService.UpdateArticleInviteUser(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !isAdmin {
		qc.actionService.ActionRecordAdd(ctx, entity.CaptchaActionInvitationAnswer, req.UserID)
	}
	handler.HandleResponse(ctx, nil, nil)
}

// GetSimilarArticles fuzzy query similar articles based on title
// @Summary fuzzy query similar articles based on title
// @Description fuzzy query similar articles based on title
// @Tags Article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param title query string true "title"  default(string)
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/article/similar [get]
func (qc *ArticleController) GetSimilarArticles(ctx *gin.Context) {
	title := ctx.Query("title")
	resp, err := qc.articleService.GetArticlesByTitle(ctx, title)
	handler.HandleResponse(ctx, err, resp)
}

// UserTop godoc
// @Summary UserTop
// @Description UserTop
// @Tags Article
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username query string true "username"  default(string)
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/personal/qa/top [get]
func (qc *ArticleController) UserTop(ctx *gin.Context) {
	userName := ctx.Query("username")
	userID := middleware.GetLoginUserIDFromContext(ctx)
	articleList, answerList, err := qc.articleService.SearchUserTopList(ctx, userName, userID)
	handler.HandleResponse(ctx, err, gin.H{
		"article": articleList,
		"answer":  answerList,
	})
}

// PersonalArticlePage list personal articles
// @Summary list personal articles
// @Description list personal articles
// @Tags Personal
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username query string true "username"  default(string)
// @Param order query string true "order"  Enums(newest,score)
// @Param page query string true "page"  default(0)
// @Param page_size query string true "page_size" default(20)
// @Success 200 {object} handler.RespBody
// @Router /personal/article/page [get]
func (qc *ArticleController) PersonalArticlePage(ctx *gin.Context) {
	req := &schema.PersonalArticlePageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	req.IsAdmin = middleware.GetUserIsAdminModerator(ctx)
	resp, err := qc.articleService.PersonalArticlePage(ctx, req)
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
func (qc *ArticleController) PersonalAnswerPage(ctx *gin.Context) {
	req := &schema.PersonalAnswerPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	req.IsAdmin = middleware.GetUserIsAdminModerator(ctx)
	resp, err := qc.articleService.PersonalAnswerPage(ctx, req)
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
func (qc *ArticleController) PersonalCollectionPage(ctx *gin.Context) {
	req := &schema.PersonalCollectionPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	resp, err := qc.articleService.PersonalCollectionPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// AdminArticlePage admin article page
// @Summary AdminArticlePage admin article page
// @Description Status:[available,closed,deleted,pending]
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "page size"
// @Param page_size query int false "page size"
// @Param status query string false "user status" Enums(available, closed, deleted, pending)
// @Param query query string false "article id or title"
// @Success 200 {object} handler.RespBody
// @Router /answer/admin/api/article/page [get]
func (qc *ArticleController) AdminArticlePage(ctx *gin.Context) {
	req := &schema.AdminArticlePageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	resp, err := qc.articleService.AdminArticlePage(ctx, req)
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
// @Param query query string false "answer id or article title"
// @Param article_id query string false "article id"
// @Success 200 {object} handler.RespBody
// @Router /answer/admin/api/answer/page [get]
func (qc *ArticleController) AdminAnswerPage(ctx *gin.Context) {
	req := &schema.AdminAnswerPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.LoginUserID = middleware.GetLoginUserIDFromContext(ctx)
	resp, err := qc.articleService.AdminAnswerPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// AdminUpdateArticleStatus update article status
// @Summary update article status
// @Description update article status
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.AdminUpdateArticleStatusReq true "AdminUpdateArticleStatusReq"
// @Success 200 {object} handler.RespBody
// @Router /answer/admin/api/article/status [put]
func (qc *ArticleController) AdminUpdateArticleStatus(ctx *gin.Context) {
	req := &schema.AdminUpdateArticleStatusReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.ArticleID = uid.DeShortID(req.ArticleID)
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	err := qc.articleService.AdminSetArticleStatus(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}
