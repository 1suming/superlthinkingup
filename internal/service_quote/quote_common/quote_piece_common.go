package quote_common

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
//
//package QuotePiececommon

import (
	"context"
	"encoding/json"
	"github.com/apache/incubator-answer/internal/base/constant"
	"github.com/apache/incubator-answer/internal/base/data"
	"github.com/apache/incubator-answer/internal/base/handler"
	"github.com/apache/incubator-answer/internal/base/reason"
	"github.com/apache/incubator-answer/internal/entity"
	"github.com/apache/incubator-answer/internal/schema"
	"github.com/apache/incubator-answer/internal/service/activity_common"
	"github.com/apache/incubator-answer/internal/service/activity_queue"
	answercommon "github.com/apache/incubator-answer/internal/service/answer_common"
	collectioncommon "github.com/apache/incubator-answer/internal/service/collection_common"
	"github.com/apache/incubator-answer/internal/service/config"
	metacommon "github.com/apache/incubator-answer/internal/service/meta_common"
	"github.com/apache/incubator-answer/internal/service/revision"
	tagcommon "github.com/apache/incubator-answer/internal/service/tag_common"
	usercommon "github.com/apache/incubator-answer/internal/service/user_common"
	"github.com/apache/incubator-answer/pkg/checker"
	"github.com/apache/incubator-answer/pkg/htmltext"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
	"math"
)

// QuotePieceRepo QuotePiece repository
type QuotePieceRepo interface {
	AddQuotePiece(ctx context.Context, QuotePiece *entity.QuotePiece) (err error)
	RemoveQuotePiece(ctx context.Context, id string) (err error)
	UpdateQuotePiece(ctx context.Context, QuotePiece *entity.QuotePiece, Cols []string) (err error)
	GetQuotePiece(ctx context.Context, id string) (QuotePiece *entity.QuotePiece, exist bool, err error)
	GetQuotePieceSimple(ctx context.Context, id string) (quotePieceBasicInfo *schema.QuotePieceBasicInfo, exist bool, err error)

	GetQuotePieceList(ctx context.Context, QuotePiece *entity.QuotePiece) (QuotePieces []*entity.QuotePiece, err error)
	GetQuotePiecePage(ctx context.Context, page, pageSize int, tagIDs []string, userID, orderCond string, inDays int, showHidden, showPending bool) (
		QuotePieceList []*entity.QuotePiece, total int64, err error)
	GetRecommendQuotePiecePageByTags(ctx context.Context, userID string, tagIDs, followedQuotePieceIDs []string, page, pageSize int) (QuotePieceList []*entity.QuotePiece, total int64, err error)
	UpdateQuotePieceStatus(ctx context.Context, QuotePieceID string, status int) (err error)
	UpdateQuotePieceStatusWithOutUpdateTime(ctx context.Context, QuotePiece *entity.QuotePiece) (err error)
	RecoverQuotePiece(ctx context.Context, QuotePieceID string) (err error)
	UpdateQuotePieceOperation(ctx context.Context, QuotePiece *entity.QuotePiece) (err error)
	GetQuotePiecesByTitle(ctx context.Context, title string, pageSize int) (QuotePieceList []*entity.QuotePiece, err error)
	GetQuotePieceByTitle(ctx context.Context, title string) (QuotePiece *entity.QuotePiece, err error)

	UpdatePvCount(ctx context.Context, QuotePieceID string) (err error)
	//UpdateAnswerCount(ctx context.Context, QuotePieceID string, num int) (err error)
	UpdateCollectionCount(ctx context.Context, QuotePieceID string) (count int64, err error)
	UpdateAccepted(ctx context.Context, QuotePiece *entity.QuotePiece) (err error)
	UpdateLastAnswer(ctx context.Context, QuotePiece *entity.QuotePiece) (err error)
	FindByID(ctx context.Context, id []string) (QuotePieceList []*entity.QuotePiece, err error)
	AdminQuotePiecePage(ctx context.Context, search *schema.AdminQuotePiecePageReq) ([]*entity.QuotePiece, int64, error)
	GetQuotePieceCount(ctx context.Context) (count int64, err error)
	GetUserQuotePieceCount(ctx context.Context, userID string, show int) (count int64, err error)
	SitemapQuotePieces(ctx context.Context, page, pageSize int) (QuotePieceIDList []*schema.SiteMapQuotePieceInfo, err error)
	RemoveAllUserQuotePiece(ctx context.Context, userID string) (err error)
	UpdateSearch(ctx context.Context, QuotePieceID string) (err error)
}

// QuotePieceCommon user service
type QuotePieceCommon struct {
	quotePieceRepo       QuotePieceRepo
	answerRepo           answercommon.AnswerRepo
	voteRepo             activity_common.VoteRepo
	followCommon         activity_common.FollowRepo
	tagCommon            *tagcommon.TagCommonService
	userCommon           *usercommon.UserCommon
	collectionCommon     *collectioncommon.CollectionCommon
	AnswerCommon         *answercommon.AnswerCommon
	metaCommonService    *metacommon.MetaCommonService
	configService        *config.ConfigService
	activityQueueService activity_queue.ActivityQueueService
	revisionRepo         revision.RevisionRepo
	data                 *data.Data
}

func NewQuotePieceCommon(quotePieceRepo QuotePieceRepo,
	answerRepo answercommon.AnswerRepo,
	voteRepo activity_common.VoteRepo,
	followCommon activity_common.FollowRepo,
	tagCommon *tagcommon.TagCommonService,
	userCommon *usercommon.UserCommon,
	collectionCommon *collectioncommon.CollectionCommon,
	answerCommon *answercommon.AnswerCommon,
	metaCommonService *metacommon.MetaCommonService,
	configService *config.ConfigService,
	activityQueueService activity_queue.ActivityQueueService,
	revisionRepo revision.RevisionRepo,
	data *data.Data,
) *QuotePieceCommon {
	return &QuotePieceCommon{
		quotePieceRepo:       quotePieceRepo,
		answerRepo:           answerRepo,
		voteRepo:             voteRepo,
		followCommon:         followCommon,
		tagCommon:            tagCommon,
		userCommon:           userCommon,
		collectionCommon:     collectionCommon,
		AnswerCommon:         answerCommon,
		metaCommonService:    metaCommonService,
		configService:        configService,
		activityQueueService: activityQueueService,
		revisionRepo:         revisionRepo,
		data:                 data,
	}
}

func (qs *QuotePieceCommon) GetUserQuotePieceCount(ctx context.Context, userID string) (count int64, err error) {
	return qs.quotePieceRepo.GetUserQuotePieceCount(ctx, userID, 0)
}

func (qs *QuotePieceCommon) GetPersonalUserQuotePieceCount(ctx context.Context, loginUserID, userID string, isAdmin bool) (count int64, err error) {
	show := entity.QuotePieceShow
	if loginUserID == userID || isAdmin {
		show = 0
	}
	return qs.quotePieceRepo.GetUserQuotePieceCount(ctx, userID, show)
}

func (qs *QuotePieceCommon) UpdatePv(ctx context.Context, QuotePieceID string) error {
	return qs.quotePieceRepo.UpdatePvCount(ctx, QuotePieceID)
}

//func (qs *QuotePieceCommon) UpdateAnswerCount(ctx context.Context, QuotePieceID string) error {
//	count, err := qs.answerRepo.GetCountByQuotePieceID(ctx, QuotePieceID)
//	if err != nil {
//		return err
//	}
//	if count == 0 {
//		err = qs.quotePieceRepo.UpdateLastAnswer(ctx, &entity.QuotePiece{
//			ID:           QuotePieceID,
//			LastAnswerID: "0",
//		})
//		if err != nil {
//			return err
//		}
//	}
//	return qs.quotePieceRepo.UpdateAnswerCount(ctx, QuotePieceID, int(count))
//}

func (qs *QuotePieceCommon) UpdateCollectionCount(ctx context.Context, QuotePieceID string) (count int64, err error) {
	return qs.quotePieceRepo.UpdateCollectionCount(ctx, QuotePieceID)
}

func (qs *QuotePieceCommon) UpdateAccepted(ctx context.Context, QuotePieceID, AnswerID string) error {
	QuotePiece := &entity.QuotePiece{}
	QuotePiece.ID = QuotePieceID
	//QuotePiece.AcceptedAnswerID = AnswerID
	return qs.quotePieceRepo.UpdateAccepted(ctx, QuotePiece)
}

func (qs *QuotePieceCommon) UpdateLastAnswer(ctx context.Context, QuotePieceID, AnswerID string) error {
	QuotePiece := &entity.QuotePiece{}
	QuotePiece.ID = QuotePieceID
	//QuotePiece.LastAnswerID = AnswerID
	return qs.quotePieceRepo.UpdateLastAnswer(ctx, QuotePiece)
}

//func (qs *QuotePieceCommon) UpdatePostTime(ctx context.Context, QuotePieceID string) error {
//	QuotePieceinfo := &entity.QuotePiece{}
//	now := time.Now()
//	_ = now
//	QuotePieceinfo.ID = QuotePieceID
//	QuotePieceinfo.PostUpdateTime = now
//	return qs.quotePieceRepo.UpdateQuotePiece(ctx, QuotePieceinfo, []string{"post_update_time"})
//}
//func (qs *QuotePieceCommon) UpdatePostSetTime(ctx context.Context, QuotePieceID string, setTime time.Time) error {
//	QuotePieceinfo := &entity.QuotePiece{}
//	QuotePieceinfo.ID = QuotePieceID
//	QuotePieceinfo.PostUpdateTime = setTime
//	return qs.quotePieceRepo.UpdateQuotePiece(ctx, QuotePieceinfo, []string{"post_update_time"})
//}

func (qs *QuotePieceCommon) FindInfoByID(ctx context.Context, QuotePieceIDs []string, loginUserID string) (map[string]*schema.QuotePieceInfoResp, error) {
	list := make(map[string]*schema.QuotePieceInfoResp)
	QuotePieceList, err := qs.quotePieceRepo.FindByID(ctx, QuotePieceIDs)
	if err != nil {
		return list, err
	}
	QuotePieces, err := qs.FormatQuotePieces(ctx, QuotePieceList, loginUserID)
	if err != nil {
		return list, err
	}
	for _, item := range QuotePieces {
		list[item.ID] = item
	}
	return list, nil
}

func (qs *QuotePieceCommon) InviteUserInfo(ctx context.Context, QuotePieceID string) (inviteList []*schema.UserBasicInfo, err error) {
	return
	//InviteUserInfo := make([]*schema.UserBasicInfo, 0)
	//dbinfo, has, err := qs.quotePieceRepo.GetQuotePiece(ctx, QuotePieceID)
	//if err != nil {
	//	return InviteUserInfo, err
	//}
	//if !has {
	//	return InviteUserInfo, errors.NotFound(reason.QuotePieceNotFound)
	//}
	///@ms: InviteUser
	//if dbinfo.InviteUserID != "" {
	//	InviteUserIDs := make([]string, 0)
	//	err := json.Unmarshal([]byte(dbinfo.InviteUserID), &InviteUserIDs)
	//	if err == nil {
	//		inviteUserInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, InviteUserIDs)
	//		if err == nil {
	//			for _, userid := range InviteUserIDs {
	//				_, ok := inviteUserInfoMap[userid]
	//				if ok {
	//					InviteUserInfo = append(InviteUserInfo, inviteUserInfoMap[userid])
	//				}
	//			}
	//		}
	//	}
	//}
	//return InviteUserInfo, nil
}

func (qs *QuotePieceCommon) Info(ctx context.Context, QuotePieceID string, loginUserID string) (resp *schema.QuotePieceInfoResp, err error) {
	QuotePieceInfo, has, err := qs.quotePieceRepo.GetQuotePiece(ctx, QuotePieceID)
	if err != nil {
		return resp, err
	}
	QuotePieceInfo.ID = uid.DeShortID(QuotePieceInfo.ID)
	if !has {
		return resp, errors.NotFound(reason.QuoteNotFound)
	}
	resp = qs.ShowFormat(ctx, QuotePieceInfo)
	if resp.Status == entity.QuotePieceStatusClosed {
		metaInfo, err := qs.metaCommonService.GetMetaByObjectIdAndKey(ctx, QuotePieceInfo.ID, entity.QuotePieceCloseReasonKey)
		if err != nil {
			log.Error(err)
		} else {
			closeMsg := &schema.CloseQuotePieceMeta{}
			err = json.Unmarshal([]byte(metaInfo.Value), closeMsg)
			if err != nil {
				log.Error("json.Unmarshal CloseQuotePieceMeta error", err.Error())
			} else {
				cfg, err := qs.configService.GetConfigByID(ctx, closeMsg.CloseType)
				if err != nil {
					log.Error("json.Unmarshal QuotePieceCloseJson error", err.Error())
				} else {
					reasonItem := &schema.ReasonItem{}
					_ = json.Unmarshal(cfg.GetByteValue(), reasonItem)
					reasonItem.Translate(cfg.Key, handler.GetLangByCtx(ctx))
					operation := &schema.Operation{}
					operation.Type = reasonItem.Name
					operation.Description = reasonItem.Description
					operation.Msg = closeMsg.CloseMsg
					operation.Time = metaInfo.CreatedAt.Unix()
					operation.Level = schema.OperationLevelInfo
					resp.Operation = operation
				}
			}
		}
	}

	if resp.Status != entity.QuotePieceStatusDeleted {
		if resp.Tags, err = qs.tagCommon.GetObjectTag(ctx, QuotePieceID); err != nil {
			return resp, err
		}
	} else {
		revisionInfo, exist, err := qs.revisionRepo.GetLastRevisionByObjectID(ctx, QuotePieceID)
		if err != nil {
			log.Errorf("get revision error %s", err)
		}
		if exist {
			QuotePieceWithTagsRevision := &entity.QuotePieceWithTagsRevision{}
			if err = json.Unmarshal([]byte(revisionInfo.Content), QuotePieceWithTagsRevision); err != nil {
				log.Errorf("revision parsing error %s", err)
				return resp, nil
			}
			for _, tag := range QuotePieceWithTagsRevision.Tags {
				resp.Tags = append(resp.Tags, &schema.TagResp{
					ID:              tag.ID,
					SlugName:        tag.SlugName,
					DisplayName:     tag.DisplayName,
					MainTagSlugName: tag.MainTagSlugName,
					Recommend:       tag.Recommend,
					Reserved:        tag.Reserved,
				})

			}
		}
	}
	for _, v := range resp.Tags {
		log.Infof("my_tags id:%+v", *v)
	}

	userIds := make([]string, 0)
	if checker.IsNotZeroString(QuotePieceInfo.UserID) {
		userIds = append(userIds, QuotePieceInfo.UserID)
	}
	//if checker.IsNotZeroString(QuotePieceInfo.LastEditUserID) {
	//	userIds = append(userIds, QuotePieceInfo.LastEditUserID)
	//}
	if checker.IsNotZeroString(resp.LastAnsweredUserID) {
		userIds = append(userIds, resp.LastAnsweredUserID)
	}
	userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIds)
	if err != nil {
		return resp, err
	}
	resp.UserInfo = userInfoMap[QuotePieceInfo.UserID]
	//resp.UpdateUserInfo = userInfoMap[QuotePieceInfo.LastEditUserID]
	resp.LastAnsweredUserInfo = userInfoMap[resp.LastAnsweredUserID]
	if len(loginUserID) == 0 {
		return resp, nil
	}

	resp.VoteStatus = qs.voteRepo.GetVoteStatus(ctx, QuotePieceID, loginUserID)
	resp.IsFollowed, _ = qs.followCommon.IsFollowed(ctx, loginUserID, QuotePieceID)

	ids, err := qs.AnswerCommon.SearchAnswerIDs(ctx, loginUserID, QuotePieceInfo.ID)
	if err != nil {
		log.Error("AnswerFunc.SearchAnswerIDs", err)
	}
	resp.Answered = len(ids) > 0
	if resp.Answered {
		resp.FirstAnswerId = ids[0]
	}

	collectedMap, err := qs.collectionCommon.SearchObjectCollected(ctx, loginUserID, []string{QuotePieceInfo.ID})
	if err != nil {
		return nil, err
	}
	if len(collectedMap) > 0 {
		resp.Collected = true
	}
	return resp, nil
}

func (qs *QuotePieceCommon) FormatQuotePiecesPage(
	ctx context.Context, QuotePieceList []*entity.QuotePiece, loginUserID string, orderCond string) (
	formattedQuotePieces []*schema.QuotePiecePageResp, err error) {
	formattedQuotePieces = make([]*schema.QuotePiecePageResp, 0)
	QuotePieceIDs := make([]string, 0)
	userIDs := make([]string, 0)

	for _, QuotePieceInfo := range QuotePieceList {

		//thumbnails := GetQuotePieceThumbnails(QuotePieceInfo)

		t := &schema.QuotePiecePageResp{
			ID:              QuotePieceInfo.ID,
			CreatedAt:       QuotePieceInfo.CreatedAt.Unix(),
			Title:           QuotePieceInfo.Title,
			UrlTitle:        htmltext.UrlTitle(QuotePieceInfo.Title),
			Description:     htmltext.FetchExcerpt(QuotePieceInfo.ParsedText, "...", 80), //240),
			Status:          QuotePieceInfo.Status,
			ViewCount:       QuotePieceInfo.ViewCount,
			UniqueViewCount: QuotePieceInfo.UniqueViewCount,
			VoteCount:       QuotePieceInfo.VoteCount,
			//AnswerCount:      QuotePieceInfo.AnswerCount,
			CollectionCount: QuotePieceInfo.CollectionCount,
			FollowCount:     QuotePieceInfo.FollowCount,
			//AcceptedAnswerID: QuotePieceInfo.AcceptedAnswerID,
			//LastAnswerID: QuotePieceInfo.LastAnswerID,
			Pin:  QuotePieceInfo.Pin,
			Show: QuotePieceInfo.Show,

			//Thumbnails: thumbnails,
		}

		QuotePieceIDs = append(QuotePieceIDs, QuotePieceInfo.ID)
		userIDs = append(userIDs, QuotePieceInfo.UserID)
		haveEdited, haveAnswered := false, false
		//if checker.IsNotZeroString(QuotePieceInfo.LastEditUserID) {
		//	haveEdited = true
		//	userIDs = append(userIDs, QuotePieceInfo.LastEditUserID)
		//}
		//if checker.IsNotZeroString(QuotePieceInfo.LastAnswerID) {
		//	haveAnswered = true
		//
		//	answerInfo, exist, err := qs.answerRepo.GetAnswer(ctx, QuotePieceInfo.LastAnswerID)
		//	if err == nil && exist {
		//		if answerInfo.LastEditUserID != "0" {
		//			t.LastAnsweredUserID = answerInfo.LastEditUserID
		//		} else {
		//			t.LastAnsweredUserID = answerInfo.UserID
		//		}
		//		t.LastAnsweredAt = answerInfo.CreatedAt
		//		userIDs = append(userIDs, t.LastAnsweredUserID)
		//	}
		//}

		// if order condition is newest or nobody edited or nobody answered, only show QuotePiece author
		if orderCond == schema.QuotePieceOrderCondNewest || (!haveEdited && !haveAnswered) {
			t.OperationType = schema.QuotePiecePageRespOperationTypeAsked
			t.OperatedAt = QuotePieceInfo.CreatedAt.Unix()
			t.Operator = &schema.QuotePiecePageRespOperator{ID: QuotePieceInfo.UserID}
		} else {
			// if no one
			if haveEdited {
				t.OperationType = schema.QuotePiecePageRespOperationTypeModified
				t.OperatedAt = QuotePieceInfo.UpdatedAt.Unix()
				//t.Operator = &schema.QuotePiecePageRespOperator{ID: QuotePieceInfo.LastEditUserID}
			}

			if haveAnswered {
				if t.LastAnsweredAt.Unix() > t.OperatedAt {
					t.OperationType = schema.QuotePiecePageRespOperationTypeAnswered
					t.OperatedAt = t.LastAnsweredAt.Unix()
					t.Operator = &schema.QuotePiecePageRespOperator{ID: t.LastAnsweredUserID}
				}
			}
		}
		formattedQuotePieces = append(formattedQuotePieces, t)
	}

	tagsMap, err := qs.tagCommon.BatchGetObjectTag(ctx, QuotePieceIDs)
	if err != nil {
		return formattedQuotePieces, err
	}
	userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIDs)
	if err != nil {
		return formattedQuotePieces, err
	}

	for _, item := range formattedQuotePieces {
		tags, ok := tagsMap[item.ID]
		if ok {
			item.Tags = tags
		} else {
			item.Tags = make([]*schema.TagResp, 0)
		}
		userInfo, ok := userInfoMap[item.Operator.ID]
		if ok {
			if userInfo != nil {
				item.Operator.DisplayName = userInfo.DisplayName
				item.Operator.Username = userInfo.Username
				item.Operator.Rank = userInfo.Rank
				item.Operator.Status = userInfo.Status
			}
		}

	}
	return formattedQuotePieces, nil
}

func (qs *QuotePieceCommon) FormatQuotePieces(ctx context.Context, QuotePieceList []*entity.QuotePiece, loginUserID string) ([]*schema.QuotePieceInfoResp, error) {
	list := make([]*schema.QuotePieceInfoResp, 0)
	objectIds := make([]string, 0)
	userIds := make([]string, 0)

	for _, QuotePieceInfo := range QuotePieceList {
		item := qs.ShowFormat(ctx, QuotePieceInfo)
		list = append(list, item)
		objectIds = append(objectIds, item.ID)
		userIds = append(userIds, item.UserID, item.LastEditUserID, item.LastAnsweredUserID)
	}
	tagsMap, err := qs.tagCommon.BatchGetObjectTag(ctx, objectIds)
	if err != nil {
		return list, err
	}

	userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIds)
	if err != nil {
		return list, err
	}

	for _, item := range list {
		item.Tags = tagsMap[item.ID]
		item.UserInfo = userInfoMap[item.UserID]
		item.UpdateUserInfo = userInfoMap[item.LastEditUserID]
		item.LastAnsweredUserInfo = userInfoMap[item.LastAnsweredUserID]
	}
	if loginUserID == "" {
		return list, nil
	}

	collectedMap, err := qs.collectionCommon.SearchObjectCollected(ctx, loginUserID, objectIds)
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		item.Collected = collectedMap[item.ID]
	}
	return list, nil
}

// RemoveQuotePiece delete QuotePiece
func (qs *QuotePieceCommon) RemoveQuotePiece(ctx context.Context, req *schema.RemoveQuotePieceReq) (err error) {
	QuotePieceInfo, has, err := qs.quotePieceRepo.GetQuotePiece(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	if QuotePieceInfo.Status == entity.QuotePieceStatusDeleted {
		return nil
	}

	QuotePieceInfo.Status = entity.QuotePieceStatusDeleted
	err = qs.quotePieceRepo.UpdateQuotePieceStatus(ctx, QuotePieceInfo.ID, QuotePieceInfo.Status)
	if err != nil {
		return err
	}

	userQuotePieceCount, err := qs.GetUserQuotePieceCount(ctx, QuotePieceInfo.UserID)
	if err != nil {
		log.Error("user GetUserQuotePieceCount error", err.Error())
	} else {
		//@ms:TODO	err = qs.userCommon.UpdateQuotePieceCount(ctx, QuotePieceInfo.UserID, userQuotePieceCount)
		_ = userQuotePieceCount
		if err != nil {
			log.Error("user IncreaseQuotePieceCount error", err.Error())
		}
	}

	return nil
}

func (qs *QuotePieceCommon) CloseQuotePiece(ctx context.Context, req *schema.CloseQuotePieceReq) error {
	QuotePieceInfo, has, err := qs.quotePieceRepo.GetQuotePiece(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	QuotePieceInfo.Status = entity.QuotePieceStatusClosed
	err = qs.quotePieceRepo.UpdateQuotePieceStatus(ctx, QuotePieceInfo.ID, QuotePieceInfo.Status)
	if err != nil {
		return err
	}

	closeMeta, _ := json.Marshal(schema.CloseQuotePieceMeta{
		CloseType: req.CloseType,
		CloseMsg:  req.CloseMsg,
	})
	err = qs.metaCommonService.AddMeta(ctx, req.ID, entity.QuotePieceCloseReasonKey, string(closeMeta))
	if err != nil {
		return err
	}

	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           QuotePieceInfo.UserID,
		ObjectID:         QuotePieceInfo.ID,
		OriginalObjectID: QuotePieceInfo.ID,
		ActivityTypeKey:  constant.ActQuotePieceClosed,
	})
	return nil
}

//// RemoveAnswer delete answer
//func (qs *QuotePieceCommon) RemoveAnswer(ctx context.Context, id string) (err error) {
//	answerinfo, has, err := qs.answerRepo.GetByID(ctx, id)
//	if err != nil {
//		return err
//	}
//	if !has {
//		return nil
//	}
//
//	// user add QuotePiece count
//
//	err = qs.UpdateAnswerCount(ctx, answerinfo.QuotePieceID)
//	if err != nil {
//		log.Error("UpdateAnswerCount error", err.Error())
//	}
//	userAnswerCount, err := qs.answerRepo.GetCountByUserID(ctx, answerinfo.UserID)
//	if err != nil {
//		log.Error("GetCountByUserID error", err.Error())
//	}
//	err = qs.userCommon.UpdateAnswerCount(ctx, answerinfo.UserID, int(userAnswerCount))
//	if err != nil {
//		log.Error("user UpdateAnswerCount error", err.Error())
//	}
//
//	return qs.answerRepo.RemoveAnswer(ctx, id)
//}

func (qs *QuotePieceCommon) SitemapCron(ctx context.Context) {
	QuotePieceNum, err := qs.quotePieceRepo.GetQuotePieceCount(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	if QuotePieceNum <= constant.SitemapMaxSize {
		_, err = qs.quotePieceRepo.SitemapQuotePieces(ctx, 1, int(QuotePieceNum))
		if err != nil {
			log.Errorf("get site map QuotePiece error: %v", err)
		}
		return
	}

	totalPages := int(math.Ceil(float64(QuotePieceNum) / float64(constant.SitemapMaxSize)))
	for i := 1; i <= totalPages; i++ {
		_, err = qs.quotePieceRepo.SitemapQuotePieces(ctx, i, constant.SitemapMaxSize)
		if err != nil {
			log.Errorf("get site map QuotePiece error: %v", err)
			return
		}
	}
}

func (qs *QuotePieceCommon) SetCache(ctx context.Context, cachekey string, info interface{}) error {
	infoStr, err := json.Marshal(info)
	if err != nil {
		return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	err = qs.data.Cache.SetString(ctx, cachekey, string(infoStr), schema.DashboardCacheTime)
	if err != nil {
		return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	return nil
}

func (qs *QuotePieceCommon) ShowListFormat(ctx context.Context, data *entity.QuotePiece) *schema.QuotePieceInfoResp {
	return qs.ShowFormat(ctx, data)
}

func (qs *QuotePieceCommon) ShowFormat(ctx context.Context, data *entity.QuotePiece) *schema.QuotePieceInfoResp {
	info := schema.QuotePieceInfoResp{}
	info.ID = data.ID
	if handler.GetEnableShortID(ctx) {
		info.ID = uid.EnShortID(data.ID)
	}
	info.Title = data.Title
	info.UrlTitle = htmltext.UrlTitle(data.Title)
	//info.Content = data.Bio
	info.HTML = data.ParsedText
	info.ViewCount = data.ViewCount
	info.UniqueViewCount = data.UniqueViewCount
	info.VoteCount = data.VoteCount
	//info.AnswerCount = data.AnswerCount
	info.CollectionCount = data.CollectionCount
	info.FollowCount = data.FollowCount
	//info.AcceptedAnswerID = data.AcceptedAnswerID
	//info.LastAnswerID = data.LastAnswerID
	info.CreateTime = data.CreatedAt.Unix()
	info.UpdateTime = data.UpdatedAt.Unix()
	//info.PostUpdateTime = data.PostUpdateTime.Unix()
	//if data.PostUpdateTime.Unix() < 1 {
	//	info.PostUpdateTime = 0
	//}
	info.QuotePieceUpdateTime = data.UpdatedAt.Unix()
	if data.UpdatedAt.Unix() < 1 {
		info.QuotePieceUpdateTime = 0
	}
	info.Status = data.Status
	info.Pin = data.Pin
	info.Show = data.Show
	info.UserID = data.UserID
	//info.LastEditUserID = data.LastEditUserID
	//if data.LastAnswerID != "0" {
	//	answerInfo, exist, err := qs.answerRepo.GetAnswer(ctx, data.LastAnswerID)
	//	if err == nil && exist {
	//		if answerInfo.LastEditUserID != "0" {
	//			info.LastAnsweredUserID = answerInfo.LastEditUserID
	//		} else {
	//			info.LastAnsweredUserID = answerInfo.UserID
	//		}
	//	}
	//
	//}
	//info.ContentFormat = data.OriginalTextFormat

	info.Tags = make([]*schema.TagResp, 0)
	return &info
}
func (qs *QuotePieceCommon) ShowFormatWithTag(ctx context.Context, data *entity.QuotePieceWithTagsRevision) *schema.QuotePieceInfoResp {
	info := qs.ShowFormat(ctx, &data.QuotePiece)
	Tags := make([]*schema.TagResp, 0)
	for _, tag := range data.Tags {
		item := &schema.TagResp{}
		item.SlugName = tag.SlugName
		item.DisplayName = tag.DisplayName
		item.Recommend = tag.Recommend
		item.Reserved = tag.Reserved
		Tags = append(Tags, item)
	}
	info.Tags = Tags
	return info
}
