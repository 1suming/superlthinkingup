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

package constant

type ActivityTypeKey string

const (
	ActEdited    = "edited"
	ActClosed    = "closed"
	ActVotedDown = "voted_down"
	ActVotedUp   = "voted_up"
	ActVoteDown  = "vote_down"
	ActVoteUp    = "vote_up"
	ActUpVote    = "upvote"
	ActDownVote  = "downvote"
	ActFollow    = "follow"
	ActAccepted  = "accepted"
	ActAccept    = "accept"
	ActPin       = "pin"
	ActUnPin     = "unpin"
	ActShow      = "show"
	ActHide      = "hide"
)

const (
	ActQuestionAsked     ActivityTypeKey = "question.asked"
	ActQuestionClosed    ActivityTypeKey = "question.closed"
	ActQuestionReopened  ActivityTypeKey = "question.reopened"
	ActQuestionAnswered  ActivityTypeKey = "question.answered"
	ActQuestionCommented ActivityTypeKey = "question.commented"
	ActQuestionAccept    ActivityTypeKey = "question.accept"
	ActQuestionUpvote    ActivityTypeKey = "question.upvote"
	ActQuestionDownVote  ActivityTypeKey = "question.downvote"
	ActQuestionEdited    ActivityTypeKey = "question.edited"
	ActQuestionRollback  ActivityTypeKey = "question.rollback"
	ActQuestionDeleted   ActivityTypeKey = "question.deleted"
	ActQuestionUndeleted ActivityTypeKey = "question.undeleted"
	ActQuestionPin       ActivityTypeKey = "question.pin"
	ActQuestionUnPin     ActivityTypeKey = "question.unpin"
	ActQuestionHide      ActivityTypeKey = "question.hide"
	ActQuestionShow      ActivityTypeKey = "question.show"
)

const (
	ActAnswerAnswered  ActivityTypeKey = "answer.answered"
	ActAnswerCommented ActivityTypeKey = "answer.commented"
	ActAnswerAccept    ActivityTypeKey = "answer.accept"
	ActAnswerUpvote    ActivityTypeKey = "answer.upvote"
	ActAnswerDownVote  ActivityTypeKey = "answer.downvote"
	ActAnswerEdited    ActivityTypeKey = "answer.edited"
	ActAnswerRollback  ActivityTypeKey = "answer.rollback"
	ActAnswerDeleted   ActivityTypeKey = "answer.deleted"
	ActAnswerUndeleted ActivityTypeKey = "answer.undeleted"
)

const (
	ActTagCreated   ActivityTypeKey = "tag.created"
	ActTagEdited    ActivityTypeKey = "tag.edited"
	ActTagRollback  ActivityTypeKey = "tag.rollback"
	ActTagDeleted   ActivityTypeKey = "tag.deleted"
	ActTagUndeleted ActivityTypeKey = "tag.undeleted"
)

// @ms:
const (
	ActArticleAsked     ActivityTypeKey = "question.asked"
	ActArticleClosed    ActivityTypeKey = "question.closed"
	ActArticleReopened  ActivityTypeKey = "question.reopened"
	ActArticleAnswered  ActivityTypeKey = "question.answered"
	ActArticleCommented ActivityTypeKey = "question.commented"
	ActArticleAccept    ActivityTypeKey = "question.accept"
	ActArticleUpvote    ActivityTypeKey = "question.upvote"
	ActArticleDownVote  ActivityTypeKey = "question.downvote"
	ActArticleEdited    ActivityTypeKey = "question.edited"
	ActArticleRollback  ActivityTypeKey = "question.rollback"
	ActArticleDeleted   ActivityTypeKey = "question.deleted"
	ActArticleUndeleted ActivityTypeKey = "question.undeleted"
	ActArticlePin       ActivityTypeKey = "question.pin"
	ActArticleUnPin     ActivityTypeKey = "question.unpin"
	ActArticleHide      ActivityTypeKey = "question.hide"
	ActArticleShow      ActivityTypeKey = "question.show"
)

const (
	ActQuoteAsked     ActivityTypeKey = "quote.asked"
	ActQuoteClosed    ActivityTypeKey = "quote.closed"
	ActQuoteReopened  ActivityTypeKey = "quote.reopened"
	ActQuoteAnswered  ActivityTypeKey = "quote.answered"
	ActQuoteCommented ActivityTypeKey = "quote.commented"
	ActQuoteAccept    ActivityTypeKey = "quote.accept"
	ActQuoteUpvote    ActivityTypeKey = "quote.upvote"
	ActQuoteDownVote  ActivityTypeKey = "quote.downvote"
	ActQuoteEdited    ActivityTypeKey = "quote.edited"
	ActQuoteRollback  ActivityTypeKey = "quote.rollback"
	ActQuoteDeleted   ActivityTypeKey = "quote.deleted"
	ActQuoteUndeleted ActivityTypeKey = "quote.undeleted"
	ActQuotePin       ActivityTypeKey = "quote.pin"
	ActQuoteUnPin     ActivityTypeKey = "quote.unpin"
	ActQuoteHide      ActivityTypeKey = "quote.hide"
	ActQuoteShow      ActivityTypeKey = "quote.show"
)
const (
	ActQuoteAuthorAsked     ActivityTypeKey = "quote_author.asked"
	ActQuoteAuthorClosed    ActivityTypeKey = "quote_author.closed"
	ActQuoteAuthorReopened  ActivityTypeKey = "quote_author.reopened"
	ActQuoteAuthorAnswered  ActivityTypeKey = "quote_author.answered"
	ActQuoteAuthorCommented ActivityTypeKey = "quote_author.commented"
	ActQuoteAuthorAccept    ActivityTypeKey = "quote_author.accept"
	ActQuoteAuthorUpvote    ActivityTypeKey = "quote_author.upvote"
	ActQuoteAuthorDownVote  ActivityTypeKey = "quote_author.downvote"
	ActQuoteAuthorEdited    ActivityTypeKey = "quote_author.edited"
	ActQuoteAuthorRollback  ActivityTypeKey = "quote_author.rollback"
	ActQuoteAuthorDeleted   ActivityTypeKey = "quote_author.deleted"
	ActQuoteAuthorUndeleted ActivityTypeKey = "quote_author.undeleted"
	ActQuoteAuthorPin       ActivityTypeKey = "quote_author.pin"
	ActQuoteAuthorUnPin     ActivityTypeKey = "quote_author.unpin"
	ActQuoteAuthorHide      ActivityTypeKey = "quote_author.hide"
	ActQuoteAuthorShow      ActivityTypeKey = "quote_author.show"
)

const (
	ActQuotePieceAsked     ActivityTypeKey = "quote_piece.asked"
	ActQuotePieceClosed    ActivityTypeKey = "quote_piece.closed"
	ActQuotePieceReopened  ActivityTypeKey = "quote_piece.reopened"
	ActQuotePieceAnswered  ActivityTypeKey = "quote_piece.answered"
	ActQuotePieceCommented ActivityTypeKey = "quote_piece.commented"
	ActQuotePieceAccept    ActivityTypeKey = "quote_piece.accept"
	ActQuotePieceUpvote    ActivityTypeKey = "quote_piece.upvote"
	ActQuotePieceDownVote  ActivityTypeKey = "quote_piece.downvote"
	ActQuotePieceEdited    ActivityTypeKey = "quote_piece.edited"
	ActQuotePieceRollback  ActivityTypeKey = "quote_piece.rollback"
	ActQuotePieceDeleted   ActivityTypeKey = "quote_piece.deleted"
	ActQuotePieceUndeleted ActivityTypeKey = "quote_piece.undeleted"
	ActQuotePiecePin       ActivityTypeKey = "quote_piece.pin"
	ActQuotePieceUnPin     ActivityTypeKey = "quote_piece.unpin"
	ActQuotePieceHide      ActivityTypeKey = "quote_piece.hide"
	ActQuotePieceShow      ActivityTypeKey = "quote_piece.show"
)
