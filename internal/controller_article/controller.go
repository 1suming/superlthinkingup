package controller_article

import "github.com/google/wire"

// ProviderSetController is controller providers.
var ProviderSetController = wire.NewSet(

	NewArticleController,
)
