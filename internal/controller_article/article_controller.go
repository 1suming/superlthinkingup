package controller_article

import (
	"github.com/apache/incubator-answer/internal/base/handler"
	"github.com/apache/incubator-answer/internal/service_article"
	"github.com/gin-gonic/gin"
)

//@csw

type ArticleController struct {
	articleService *service_article.ArticleService
}

func NewArticleController(articleService *service_article.ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}
func (ac *ArticleController) Get()    {}
func (ac *ArticleController) Post()   {}
func (ac *ArticleController) Put()    {}
func (ac *ArticleController) Delete() {}

func (ac *ArticleController) ArticlePage(ctx *gin.Context) {
	handler.HandleResponse(ctx, nil, "success")
}
