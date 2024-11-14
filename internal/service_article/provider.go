package service_article

import "github.com/google/wire"

var ProviderSetService = wire.NewSet(
	NewArticleService,
)
