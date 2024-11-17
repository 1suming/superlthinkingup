package service_article

import (
	articlecommon "github.com/apache/incubator-answer/internal/service/article_common"
	"github.com/google/wire"
)

var ProviderSetService = wire.NewSet(
	NewArticleService,
	articlecommon.NewArticleCommon,
)
