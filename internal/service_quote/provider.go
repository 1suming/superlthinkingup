package service_quote

import (
	"github.com/apache/incubator-answer/internal/service_quote/quote_common"
	"github.com/google/wire"
)

var ProviderSetService = wire.NewSet(

	NewQuoteService,
	quote_common.NewQuoteCommon,

	NewQuoteAuthorService,
	quote_common.NewQuoteAuthorCommon,
	quote_common.NewQuotePieceCommon,
	NewQuotePieceService,
)
