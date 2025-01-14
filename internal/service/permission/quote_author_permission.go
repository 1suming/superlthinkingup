package permission

import (
	"context"
	"github.com/apache/incubator-answer/internal/entity"

	"github.com/apache/incubator-answer/internal/base/handler"
	"github.com/apache/incubator-answer/internal/base/translator"
	"github.com/apache/incubator-answer/internal/schema"
)

// GetQuoteAuthorPermission get QuoteAuthor permission
func GetQuoteAuthorPermission(ctx context.Context, userID string, creatorUserID string, status int,
	canEdit, canDelete, canClose, canReopen, canPin, canHide, canUnPin, canShow, canRecover bool) (
	actions []*schema.PermissionMemberAction) {
	lang := handler.GetLangByCtx(ctx)
	actions = make([]*schema.PermissionMemberAction, 0)
	if len(userID) > 0 {
		actions = append(actions, &schema.PermissionMemberAction{
			Action: "report",
			Name:   translator.Tr(lang, reportActionName),
			Type:   "reason",
		})
	}
	if (canEdit || userID == creatorUserID) && status != entity.QuoteAuthorStatusDeleted {
		actions = append(actions, &schema.PermissionMemberAction{
			Action: "edit",
			Name:   translator.Tr(lang, editActionName),
			Type:   "edit",
		})
	}
	if canClose && status == entity.QuoteAuthorStatusAvailable {
		actions = append(actions, &schema.PermissionMemberAction{
			Action: "close",
			Name:   translator.Tr(lang, closeActionName),
			Type:   "confirm",
		})
	}
	if canReopen {
		actions = append(actions, &schema.PermissionMemberAction{
			Action: "reopen",
			Name:   translator.Tr(lang, reopenActionName),
			Type:   "confirm",
		})
	}
	if canPin {
		actions = append(actions, &schema.PermissionMemberAction{
			Action: "pin",
			Name:   translator.Tr(lang, pinActionName),
			Type:   "confirm",
		})
	}
	if canHide {
		actions = append(actions, &schema.PermissionMemberAction{
			Action: "hide",
			Name:   translator.Tr(lang, hideActionName),
			Type:   "confirm",
		})
	}

	if canUnPin {
		actions = append(actions, &schema.PermissionMemberAction{
			Action: "unpin",
			Name:   translator.Tr(lang, unpinActionName),
			Type:   "confirm",
		})
	}

	if canShow {
		actions = append(actions, &schema.PermissionMemberAction{
			Action: "show",
			Name:   translator.Tr(lang, showActionName),
			Type:   "confirm",
		})
	}

	if (canDelete || userID == creatorUserID) && status != entity.QuoteAuthorStatusDeleted {
		actions = append(actions, &schema.PermissionMemberAction{
			Action: "delete",
			Name:   translator.Tr(lang, deleteActionName),
			Type:   "confirm",
		})
	}

	if canRecover && status == entity.QuoteAuthorStatusDeleted {
		actions = append(actions, &schema.PermissionMemberAction{
			Action: "undelete",
			Name:   translator.Tr(lang, undeleteActionName),
			Type:   "confirm",
		})
	}
	return actions
}

// GetQuoteAuthorExtendsPermission get QuoteAuthor extends permission
func GetQuoteAuthorExtendsPermission(ctx context.Context, canInviteOtherToAnswer bool) (
	actions []*schema.PermissionMemberAction) {
	lang := handler.GetLangByCtx(ctx)
	actions = make([]*schema.PermissionMemberAction, 0)
	if canInviteOtherToAnswer {
		actions = append(actions, &schema.PermissionMemberAction{
			Action: "invite_other_to_answer",
			Name:   translator.Tr(lang, inviteSomeoneToAnswerActionName),
			Type:   "confirm",
		})
	}
	return actions
}
