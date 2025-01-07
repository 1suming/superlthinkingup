package controller_quote

import "github.com/google/wire"

// ProviderSetController is controller providers.
var ProviderSetController = wire.NewSet(

	NewQuoteController,
	NewQuoteAuthorController,
	NewQuotePieceController,
)
