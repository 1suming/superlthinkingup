package router

import (
	"github.com/apache/incubator-answer/internal/controller_quote"
	"github.com/gin-gonic/gin"
)

type QuoteAPIRouter struct {
	quoteController       *controller_quote.QuoteController
	quoteAuthorController *controller_quote.QuoteAuthorController
	quotePieceController  *controller_quote.QuotePieceController
}

func NewQuoteAPIRouter(
	quoteController *controller_quote.QuoteController,
	quoteAuthorController *controller_quote.QuoteAuthorController,
	quotePieceController *controller_quote.QuotePieceController,

) *QuoteAPIRouter {
	return &QuoteAPIRouter{
		quoteController:       quoteController,
		quoteAuthorController: quoteAuthorController,
		quotePieceController:  quotePieceController,
	}
}

// 未授权也可以访问
func (a *QuoteAPIRouter) RegisterUnAuthQuoteAPIRouter(r *gin.RouterGroup) {
	//类似：	r.GET("/quote/page", a.quoteController.QuotePage)

	// quote
	r.GET("/quote/info", a.quoteController.GetQuote) //详情
	r.GET("/quote/invite", a.quoteController.GetQuoteInviteUserInfo)
	r.GET("/quote/page", a.quoteController.QuotePage) //列表
	r.GET("/quote/recommend/page", a.quoteController.QuoteRecommendPage)
	r.GET("/quote/similar/tag", a.quoteController.SimilarQuote)
	//r.GET("/personal/qa/top", a.quoteController.UserTop)
	r.GET("/personal/quote/page", a.quoteController.PersonalQuotePage)

	r.GET("/quote/similar", a.quoteController.GetSimilarQuotes)

	//r.POST("/quote/author", a.adminUserController.AddUser)
	//r.PUT("/quote/author/info", a.userController.UserUpdateInfo)
	//r.GET("/quote/author/info/search", a.userController.SearchUserListByNam

	r.GET("/quote/author/similar", a.quoteAuthorController.GetSimilarQuoteAuthors)

	r.GET("/quote/piece/similar", a.quotePieceController.GetSimilarQuotePieces)

}
func (a *QuoteAPIRouter) RegisterQuoteAPIRouter(r *gin.RouterGroup) {
	// quote
	r.POST("/quote", a.quoteController.AddQuote)
	//r.POST("/quote/answer", a.quoteController.AddQuoteByAnswer)
	r.PUT("/quote", a.quoteController.UpdateQuote)
	//r.PUT("/quote/invite", a.quoteController.UpdateQuoteInviteUser)
	r.DELETE("/quote", a.quoteController.RemoveQuote)
	r.PUT("/quote/status", a.quoteController.CloseQuote)
	r.PUT("/quote/operation", a.quoteController.OperationQuote)
	//r.PUT("/quote/reopen", a.quoteController.ReopenQuote)
	r.POST("/quote/recover", a.quoteController.QuoteRecover)
}
