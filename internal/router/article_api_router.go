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
	//类似：	r.GET("/question/page", a.questionController.QuestionPage)

	r.GET("/article/page", a.articleController.ArticlePage)
}
