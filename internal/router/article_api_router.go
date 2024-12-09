package router

import (
	"github.com/apache/incubator-answer/internal/controller_article"
	"github.com/gin-gonic/gin"
)

type ArticleAPIRouter struct {
	articleController *controller_article.ArticleController
}

func NewArticleAPIRouter(
	articleController *controller_article.ArticleController) *ArticleAPIRouter {
	return &ArticleAPIRouter{
		articleController: articleController,
	}
}

// 未授权也可以访问
func (a *ArticleAPIRouter) RegisterUnAuthArticleAPIRouter(r *gin.RouterGroup) {
	//类似：	r.GET("/article/page", a.articleController.ArticlePage)

	// article
	r.GET("/article/info", a.articleController.GetArticle) //详情
	r.GET("/article/invite", a.articleController.GetArticleInviteUserInfo)
	r.GET("/article/page", a.articleController.ArticlePage) //列表
	r.GET("/article/recommend/page", a.articleController.ArticleRecommendPage)
	r.GET("/article/similar/tag", a.articleController.SimilarArticle)
	//r.GET("/personal/qa/top", a.articleController.UserTop)
	r.GET("/personal/article/page", a.articleController.PersonalArticlePage)

}
func (a *ArticleAPIRouter) RegisterArticleAPIRouter(r *gin.RouterGroup) {
	// article
	r.POST("/article", a.articleController.AddArticle)
	//r.POST("/article/answer", a.articleController.AddArticleByAnswer)
	r.PUT("/article", a.articleController.UpdateArticle)
	//r.PUT("/article/invite", a.articleController.UpdateArticleInviteUser)
	r.DELETE("/article", a.articleController.RemoveArticle)
	r.PUT("/article/status", a.articleController.CloseArticle)
	r.PUT("/article/operation", a.articleController.OperationArticle)
	//r.PUT("/article/reopen", a.articleController.ReopenArticle)
	r.GET("/article/similar", a.articleController.GetSimilarArticles)
	r.POST("/article/recover", a.articleController.ArticleRecover)
}
