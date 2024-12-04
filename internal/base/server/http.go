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

package server

import (
	brotli "github.com/anargu/gin-brotli"
	"github.com/apache/incubator-answer/internal/base/middleware"
	"github.com/apache/incubator-answer/internal/router"
	"github.com/apache/incubator-answer/plugin"
	"github.com/apache/incubator-answer/ui"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
)

// NewHTTPServer new http server.
func NewHTTPServer(debug bool,
	staticRouter *router.StaticRouter,
	answerRouter *router.AnswerAPIRouter,
	swaggerRouter *router.SwaggerRouter,
	viewRouter *router.UIRouter,
	authUserMiddleware *middleware.AuthUserMiddleware,
	avatarMiddleware *middleware.AvatarMiddleware,
	shortIDMiddleware *middleware.ShortIDMiddleware,
	templateRouter *router.TemplateRouter,
	pluginAPIRouter *router.PluginAPIRouter,
	uiConf *UI,
	articleRouter *router.ArticleAPIRouter, //@csw
) *gin.Engine {

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	//Gin框架，body参数只能读取一次 https://blog.csdn.net/impressionw/article/details/84194783

	//formatter := func(p gin.LogFormatterParams) string {
	//	return fmt.Sprintf("[HTTP] %s %s %s %d %s\n",
	//		p.TimeStamp.Format("2006-01-02_15:04:05"),
	//		p.Method,
	//		p.Path,
	//		p.StatusCode,
	//		p.ClientIP,
	//
	//	)
	//}
	//conf := gin.LoggerConfig{
	//	SkipPaths: []string{},
	//	Output:    os.Stderr,
	//	Formatter: formatter,
	//}
	r.Use(middleware.LoggerWithConfig()) //@cws

	r.Use(brotli.Brotli(brotli.DefaultCompression), middleware.ExtractAndSetAcceptLanguage, shortIDMiddleware.SetShortIDFlag())
	r.GET("/healthz", func(ctx *gin.Context) { ctx.String(200, "OK") })

	html, _ := fs.Sub(ui.Template, "template")
	htmlTemplate := template.Must(template.New("").Funcs(funcMap).ParseFS(html, "*"))
	r.SetHTMLTemplate(htmlTemplate)
	r.Use(middleware.HeadersByRequestURI())
	viewRouter.Register(r, uiConf.BaseURL)

	rootGroup := r.Group("")
	swaggerRouter.Register(rootGroup)
	static := r.Group("")
	static.Use(avatarMiddleware.AvatarThumb(), authUserMiddleware.VisitAuth())
	staticRouter.RegisterStaticRouter(static)

	// The route must be available without logging in
	mustUnAuthV1 := r.Group("/answer/api/v1")
	answerRouter.RegisterMustUnAuthAnswerAPIRouter(authUserMiddleware, mustUnAuthV1)

	// register api that no need to login
	unAuthV1 := r.Group("/answer/api/v1")
	unAuthV1.Use(authUserMiddleware.Auth(), authUserMiddleware.EjectUserBySiteInfo())
	answerRouter.RegisterUnAuthAnswerAPIRouter(unAuthV1)

	// register api that must be authenticated but no need to check account status
	authWithoutStatusV1 := r.Group("/answer/api/v1")
	authWithoutStatusV1.Use(authUserMiddleware.MustAuthWithoutAccountAvailable())
	answerRouter.RegisterAuthUserWithAnyStatusAnswerAPIRouter(authWithoutStatusV1)

	// register api that must be authenticated
	authV1 := r.Group("/answer/api/v1")
	authV1.Use(authUserMiddleware.MustAuthAndAccountAvailable())
	answerRouter.RegisterAnswerAPIRouter(authV1)

	adminauthV1 := r.Group("/answer/admin/api")
	adminauthV1.Use(authUserMiddleware.AdminAuth())
	answerRouter.RegisterAnswerAdminAPIRouter(adminauthV1)

	templateRouter.RegisterTemplateRouter(rootGroup, uiConf.BaseURL)

	// plugin routes
	pluginAPIRouter.RegisterUnAuthConnectorRouter(mustUnAuthV1)
	pluginAPIRouter.RegisterAuthUserConnectorRouter(authV1)
	pluginAPIRouter.RegisterAuthAdminConnectorRouter(adminauthV1)

	_ = plugin.CallAgent(func(agent plugin.Agent) error {
		agent.RegisterUnAuthRouter(mustUnAuthV1)
		agent.RegisterAuthUserRouter(authV1)
		agent.RegisterAuthAdminRouter(adminauthV1)
		return nil
	})
	//@csw
	// register api that no need to login
	article_unAuthV1 := r.Group("/answer/api/v1")
	article_unAuthV1.Use(authUserMiddleware.Auth(), authUserMiddleware.EjectUserBySiteInfo())
	articleRouter.RegisterUnAuthArticleAPIRouter(article_unAuthV1)

	// register api that must be authenticated
	article_authV1 := r.Group("/answer/api/v1")
	article_authV1.Use(authUserMiddleware.MustAuthAndAccountAvailable())
	articleRouter.RegisterArticleAPIRouter(article_authV1)

	return r
}
