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
	//类似：	r.GET("/question/page", a.articleController.ArticlePage)

	// article
	r.GET("/article/info", a.articleController.GetArticle) //详情
	r.GET("/article/invite", a.articleController.GetArticleInviteUserInfo)
	r.GET("/article/page", a.articleController.ArticlePage) //列表
	r.GET("/article/recommend/page", a.articleController.ArticleRecommendPage)
	r.GET("/article/similar/tag", a.articleController.SimilarArticle)
	//r.GET("/personal/qa/top", a.questionController.UserTop)
	r.GET("/personal/article/page", a.articleController.PersonalArticlePage)

}
