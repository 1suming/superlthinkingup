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
//package quoteAuthorcommon

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

// QuoteAuthorRepo quoteAuthor repository
type QuoteAuthorRepo interface {
	AddQuoteAuthor(ctx context.Context, quoteAuthor *entity.QuoteAuthor) (err error)
	RemoveQuoteAuthor(ctx context.Context, id string) (err error)
	UpdateQuoteAuthor(ctx context.Context, quoteAuthor *entity.QuoteAuthor, Cols []string) (err error)
	GetQuoteAuthor(ctx context.Context, id string) (quoteAuthor *entity.QuoteAuthor, exist bool, err error)
	GetQuoteAuthorList(ctx context.Context, quoteAuthor *entity.QuoteAuthor) (quoteAuthors []*entity.QuoteAuthor, err error)
	GetQuoteAuthorPage(ctx context.Context, page, pageSize int, tagIDs []string, userID, orderCond string, inDays int, showHidden, showPending bool) (
		quoteAuthorList []*entity.QuoteAuthor, total int64, err error)
	GetRecommendQuoteAuthorPageByTags(ctx context.Context, userID string, tagIDs, followedQuoteAuthorIDs []string, page, pageSize int) (quoteAuthorList []*entity.QuoteAuthor, total int64, err error)
	UpdateQuoteAuthorStatus(ctx context.Context, quoteAuthorID string, status int) (err error)
	UpdateQuoteAuthorStatusWithOutUpdateTime(ctx context.Context, quoteAuthor *entity.QuoteAuthor) (err error)
	RecoverQuoteAuthor(ctx context.Context, quoteAuthorID string) (err error)
	UpdateQuoteAuthorOperation(ctx context.Context, quoteAuthor *entity.QuoteAuthor) (err error)
	GetQuoteAuthorsByAuthorName(ctx context.Context, title string, pageSize int) (quoteAuthorList []*entity.QuoteAuthor, err error)
	GetQuoteAuthorByAuthorName(ctx context.Context, title string) (quoteAuthor *entity.QuoteAuthor, err error)

	UpdatePvCount(ctx context.Context, quoteAuthorID string) (err error)
	//UpdateAnswerCount(ctx context.Context, quoteAuthorID string, num int) (err error)
	UpdateCollectionCount(ctx context.Context, quoteAuthorID string) (count int64, err error)
	UpdateAccepted(ctx context.Context, quoteAuthor *entity.QuoteAuthor) (err error)
	UpdateLastAnswer(ctx context.Context, quoteAuthor *entity.QuoteAuthor) (err error)
	FindByID(ctx context.Context, id []string) (quoteAuthorList []*entity.QuoteAuthor, err error)
	AdminQuoteAuthorPage(ctx context.Context, search *schema.AdminQuoteAuthorPageReq) ([]*entity.QuoteAuthor, int64, error)
	GetQuoteAuthorCount(ctx context.Context) (count int64, err error)
	GetUserQuoteAuthorCount(ctx context.Context, userID string, show int) (count int64, err error)
	SitemapQuoteAuthors(ctx context.Context, page, pageSize int) (quoteAuthorIDList []*schema.SiteMapQuoteAuthorInfo, err error)
	RemoveAllUserQuoteAuthor(ctx context.Context, userID string) (err error)
	UpdateSearch(ctx context.Context, quoteAuthorID string) (err error)
}

// QuoteAuthorCommon user service
type QuoteAuthorCommon struct {
	quoteAuthorRepo      QuoteAuthorRepo
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

func NewQuoteAuthorCommon(quoteAuthorRepo QuoteAuthorRepo,
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
) *QuoteAuthorCommon {
	return &QuoteAuthorCommon{
		quoteAuthorRepo:      quoteAuthorRepo,
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

func (qs *QuoteAuthorCommon) GetUserQuoteAuthorCount(ctx context.Context, userID string) (count int64, err error) {
	return qs.quoteAuthorRepo.GetUserQuoteAuthorCount(ctx, userID, 0)
}

func (qs *QuoteAuthorCommon) GetPersonalUserQuoteAuthorCount(ctx context.Context, loginUserID, userID string, isAdmin bool) (count int64, err error) {
	show := entity.QuoteAuthorShow
	if loginUserID == userID || isAdmin {
		show = 0
	}
	return qs.quoteAuthorRepo.GetUserQuoteAuthorCount(ctx, userID, show)
}

func (qs *QuoteAuthorCommon) UpdatePv(ctx context.Context, quoteAuthorID string) error {
	return qs.quoteAuthorRepo.UpdatePvCount(ctx, quoteAuthorID)
}

//func (qs *QuoteAuthorCommon) UpdateAnswerCount(ctx context.Context, quoteAuthorID string) error {
//	count, err := qs.answerRepo.GetCountByQuoteAuthorID(ctx, quoteAuthorID)
//	if err != nil {
//		return err
//	}
//	if count == 0 {
//		err = qs.quoteAuthorRepo.UpdateLastAnswer(ctx, &entity.QuoteAuthor{
//			ID:           quoteAuthorID,
//			LastAnswerID: "0",
//		})
//		if err != nil {
//			return err
//		}
//	}
//	return qs.quoteAuthorRepo.UpdateAnswerCount(ctx, quoteAuthorID, int(count))
//}

func (qs *QuoteAuthorCommon) UpdateCollectionCount(ctx context.Context, quoteAuthorID string) (count int64, err error) {
	return qs.quoteAuthorRepo.UpdateCollectionCount(ctx, quoteAuthorID)
}

func (qs *QuoteAuthorCommon) UpdateAccepted(ctx context.Context, quoteAuthorID, AnswerID string) error {
	quoteAuthor := &entity.QuoteAuthor{}
	quoteAuthor.ID = quoteAuthorID
	//quoteAuthor.AcceptedAnswerID = AnswerID
	return qs.quoteAuthorRepo.UpdateAccepted(ctx, quoteAuthor)
}

func (qs *QuoteAuthorCommon) UpdateLastAnswer(ctx context.Context, quoteAuthorID, AnswerID string) error {
	quoteAuthor := &entity.QuoteAuthor{}
	quoteAuthor.ID = quoteAuthorID
	//quoteAuthor.LastAnswerID = AnswerID
	return qs.quoteAuthorRepo.UpdateLastAnswer(ctx, quoteAuthor)
}

//func (qs *QuoteAuthorCommon) UpdatePostTime(ctx context.Context, quoteAuthorID string) error {
//	quoteAuthorinfo := &entity.QuoteAuthor{}
//	now := time.Now()
//	_ = now
//	quoteAuthorinfo.ID = quoteAuthorID
//	quoteAuthorinfo.PostUpdateTime = now
//	return qs.quoteAuthorRepo.UpdateQuoteAuthor(ctx, quoteAuthorinfo, []string{"post_update_time"})
//}
//func (qs *QuoteAuthorCommon) UpdatePostSetTime(ctx context.Context, quoteAuthorID string, setTime time.Time) error {
//	quoteAuthorinfo := &entity.QuoteAuthor{}
//	quoteAuthorinfo.ID = quoteAuthorID
//	quoteAuthorinfo.PostUpdateTime = setTime
//	return qs.quoteAuthorRepo.UpdateQuoteAuthor(ctx, quoteAuthorinfo, []string{"post_update_time"})
//}

func (qs *QuoteAuthorCommon) FindInfoByID(ctx context.Context, quoteAuthorIDs []string, loginUserID string) (map[string]*schema.QuoteAuthorInfoResp, error) {
	list := make(map[string]*schema.QuoteAuthorInfoResp)
	quoteAuthorList, err := qs.quoteAuthorRepo.FindByID(ctx, quoteAuthorIDs)
	if err != nil {
		return list, err
	}
	quoteAuthors, err := qs.FormatQuoteAuthors(ctx, quoteAuthorList, loginUserID)
	if err != nil {
		return list, err
	}
	for _, item := range quoteAuthors {
		list[item.ID] = item
	}
	return list, nil
}

func (qs *QuoteAuthorCommon) InviteUserInfo(ctx context.Context, quoteAuthorID string) (inviteList []*schema.UserBasicInfo, err error) {
	return
	//InviteUserInfo := make([]*schema.UserBasicInfo, 0)
	//dbinfo, has, err := qs.quoteAuthorRepo.GetQuoteAuthor(ctx, quoteAuthorID)
	//if err != nil {
	//	return InviteUserInfo, err
	//}
	//if !has {
	//	return InviteUserInfo, errors.NotFound(reason.QuoteAuthorNotFound)
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

func (qs *QuoteAuthorCommon) Info(ctx context.Context, quoteAuthorID string, loginUserID string) (resp *schema.QuoteAuthorInfoResp, err error) {
	quoteAuthorInfo, has, err := qs.quoteAuthorRepo.GetQuoteAuthor(ctx, quoteAuthorID)
	if err != nil {
		return resp, err
	}
	quoteAuthorInfo.ID = uid.DeShortID(quoteAuthorInfo.ID)
	if !has {
		return resp, errors.NotFound(reason.QuoteNotFound)
	}
	resp = qs.ShowFormat(ctx, quoteAuthorInfo)
	if resp.Status == entity.QuoteAuthorStatusClosed {
		metaInfo, err := qs.metaCommonService.GetMetaByObjectIdAndKey(ctx, quoteAuthorInfo.ID, entity.QuoteAuthorCloseReasonKey)
		if err != nil {
			log.Error(err)
		} else {
			closeMsg := &schema.CloseQuoteAuthorMeta{}
			err = json.Unmarshal([]byte(metaInfo.Value), closeMsg)
			if err != nil {
				log.Error("json.Unmarshal CloseQuoteAuthorMeta error", err.Error())
			} else {
				cfg, err := qs.configService.GetConfigByID(ctx, closeMsg.CloseType)
				if err != nil {
					log.Error("json.Unmarshal QuoteAuthorCloseJson error", err.Error())
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

	if resp.Status != entity.QuoteAuthorStatusDeleted {
		if resp.Tags, err = qs.tagCommon.GetObjectTag(ctx, quoteAuthorID); err != nil {
			return resp, err
		}
	} else {
		revisionInfo, exist, err := qs.revisionRepo.GetLastRevisionByObjectID(ctx, quoteAuthorID)
		if err != nil {
			log.Errorf("get revision error %s", err)
		}
		if exist {
			quoteAuthorWithTagsRevision := &entity.QuoteAuthorWithTagsRevision{}
			if err = json.Unmarshal([]byte(revisionInfo.Content), quoteAuthorWithTagsRevision); err != nil {
				log.Errorf("revision parsing error %s", err)
				return resp, nil
			}
			for _, tag := range quoteAuthorWithTagsRevision.Tags {
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
	if checker.IsNotZeroString(quoteAuthorInfo.UserID) {
		userIds = append(userIds, quoteAuthorInfo.UserID)
	}
	//if checker.IsNotZeroString(quoteAuthorInfo.LastEditUserID) {
	//	userIds = append(userIds, quoteAuthorInfo.LastEditUserID)
	//}
	if checker.IsNotZeroString(resp.LastAnsweredUserID) {
		userIds = append(userIds, resp.LastAnsweredUserID)
	}
	userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIds)
	if err != nil {
		return resp, err
	}
	resp.UserInfo = userInfoMap[quoteAuthorInfo.UserID]
	//resp.UpdateUserInfo = userInfoMap[quoteAuthorInfo.LastEditUserID]
	resp.LastAnsweredUserInfo = userInfoMap[resp.LastAnsweredUserID]
	if len(loginUserID) == 0 {
		return resp, nil
	}

	resp.VoteStatus = qs.voteRepo.GetVoteStatus(ctx, quoteAuthorID, loginUserID)
	resp.IsFollowed, _ = qs.followCommon.IsFollowed(ctx, loginUserID, quoteAuthorID)

	ids, err := qs.AnswerCommon.SearchAnswerIDs(ctx, loginUserID, quoteAuthorInfo.ID)
	if err != nil {
		log.Error("AnswerFunc.SearchAnswerIDs", err)
	}
	resp.Answered = len(ids) > 0
	if resp.Answered {
		resp.FirstAnswerId = ids[0]
	}

	collectedMap, err := qs.collectionCommon.SearchObjectCollected(ctx, loginUserID, []string{quoteAuthorInfo.ID})
	if err != nil {
		return nil, err
	}
	if len(collectedMap) > 0 {
		resp.Collected = true
	}
	return resp, nil
}

func (qs *QuoteAuthorCommon) FormatQuoteAuthorsPage(
	ctx context.Context, quoteAuthorList []*entity.QuoteAuthor, loginUserID string, orderCond string) (
	formattedQuoteAuthors []*schema.QuoteAuthorPageResp, err error) {
	formattedQuoteAuthors = make([]*schema.QuoteAuthorPageResp, 0)
	quoteAuthorIDs := make([]string, 0)
	userIDs := make([]string, 0)

	for _, quoteAuthorInfo := range quoteAuthorList {

		//thumbnails := GetQuoteAuthorThumbnails(quoteAuthorInfo)

		t := &schema.QuoteAuthorPageResp{
			ID:              quoteAuthorInfo.ID,
			CreatedAt:       quoteAuthorInfo.CreatedAt.Unix(),
			AuthorName:      quoteAuthorInfo.AuthorName,
			UrlAuthorName:   htmltext.UrlTitle(quoteAuthorInfo.AuthorName),
			Description:     htmltext.FetchExcerpt(quoteAuthorInfo.Bio, "...", 80), //240),
			Status:          quoteAuthorInfo.Status,
			ViewCount:       quoteAuthorInfo.ViewCount,
			UniqueViewCount: quoteAuthorInfo.UniqueViewCount,
			VoteCount:       quoteAuthorInfo.VoteCount,
			//AnswerCount:      quoteAuthorInfo.AnswerCount,
			CollectionCount: quoteAuthorInfo.CollectionCount,
			FollowCount:     quoteAuthorInfo.FollowCount,
			//AcceptedAnswerID: quoteAuthorInfo.AcceptedAnswerID,
			//LastAnswerID: quoteAuthorInfo.LastAnswerID,
			Pin:  quoteAuthorInfo.Pin,
			Show: quoteAuthorInfo.Show,

			//Thumbnails: thumbnails,
		}

		quoteAuthorIDs = append(quoteAuthorIDs, quoteAuthorInfo.ID)
		userIDs = append(userIDs, quoteAuthorInfo.UserID)
		haveEdited, haveAnswered := false, false
		//if checker.IsNotZeroString(quoteAuthorInfo.LastEditUserID) {
		//	haveEdited = true
		//	userIDs = append(userIDs, quoteAuthorInfo.LastEditUserID)
		//}
		//if checker.IsNotZeroString(quoteAuthorInfo.LastAnswerID) {
		//	haveAnswered = true
		//
		//	answerInfo, exist, err := qs.answerRepo.GetAnswer(ctx, quoteAuthorInfo.LastAnswerID)
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

		// if order condition is newest or nobody edited or nobody answered, only show quoteAuthor author
		if orderCond == schema.QuoteAuthorOrderCondNewest || (!haveEdited && !haveAnswered) {
			t.OperationType = schema.QuoteAuthorPageRespOperationTypeAsked
			t.OperatedAt = quoteAuthorInfo.CreatedAt.Unix()
			t.Operator = &schema.QuoteAuthorPageRespOperator{ID: quoteAuthorInfo.UserID}
		} else {
			// if no one
			if haveEdited {
				t.OperationType = schema.QuoteAuthorPageRespOperationTypeModified
				t.OperatedAt = quoteAuthorInfo.UpdatedAt.Unix()
				//t.Operator = &schema.QuoteAuthorPageRespOperator{ID: quoteAuthorInfo.LastEditUserID}
			}

			if haveAnswered {
				if t.LastAnsweredAt.Unix() > t.OperatedAt {
					t.OperationType = schema.QuoteAuthorPageRespOperationTypeAnswered
					t.OperatedAt = t.LastAnsweredAt.Unix()
					t.Operator = &schema.QuoteAuthorPageRespOperator{ID: t.LastAnsweredUserID}
				}
			}
		}
		formattedQuoteAuthors = append(formattedQuoteAuthors, t)
	}

	tagsMap, err := qs.tagCommon.BatchGetObjectTag(ctx, quoteAuthorIDs)
	if err != nil {
		return formattedQuoteAuthors, err
	}
	userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIDs)
	if err != nil {
		return formattedQuoteAuthors, err
	}

	for _, item := range formattedQuoteAuthors {
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
	return formattedQuoteAuthors, nil
}

func (qs *QuoteAuthorCommon) FormatQuoteAuthors(ctx context.Context, quoteAuthorList []*entity.QuoteAuthor, loginUserID string) ([]*schema.QuoteAuthorInfoResp, error) {
	list := make([]*schema.QuoteAuthorInfoResp, 0)
	objectIds := make([]string, 0)
	userIds := make([]string, 0)

	for _, quoteAuthorInfo := range quoteAuthorList {
		item := qs.ShowFormat(ctx, quoteAuthorInfo)
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

// RemoveQuoteAuthor delete quoteAuthor
func (qs *QuoteAuthorCommon) RemoveQuoteAuthor(ctx context.Context, req *schema.RemoveQuoteAuthorReq) (err error) {
	quoteAuthorInfo, has, err := qs.quoteAuthorRepo.GetQuoteAuthor(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	if quoteAuthorInfo.Status == entity.QuoteAuthorStatusDeleted {
		return nil
	}

	quoteAuthorInfo.Status = entity.QuoteAuthorStatusDeleted
	err = qs.quoteAuthorRepo.UpdateQuoteAuthorStatus(ctx, quoteAuthorInfo.ID, quoteAuthorInfo.Status)
	if err != nil {
		return err
	}

	userQuoteAuthorCount, err := qs.GetUserQuoteAuthorCount(ctx, quoteAuthorInfo.UserID)
	if err != nil {
		log.Error("user GetUserQuoteAuthorCount error", err.Error())
	} else {
		//@ms:TODO	err = qs.userCommon.UpdateQuoteAuthorCount(ctx, quoteAuthorInfo.UserID, userQuoteAuthorCount)
		_ = userQuoteAuthorCount
		if err != nil {
			log.Error("user IncreaseQuoteAuthorCount error", err.Error())
		}
	}

	return nil
}

func (qs *QuoteAuthorCommon) CloseQuoteAuthor(ctx context.Context, req *schema.CloseQuoteAuthorReq) error {
	quoteAuthorInfo, has, err := qs.quoteAuthorRepo.GetQuoteAuthor(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	quoteAuthorInfo.Status = entity.QuoteAuthorStatusClosed
	err = qs.quoteAuthorRepo.UpdateQuoteAuthorStatus(ctx, quoteAuthorInfo.ID, quoteAuthorInfo.Status)
	if err != nil {
		return err
	}

	closeMeta, _ := json.Marshal(schema.CloseQuoteAuthorMeta{
		CloseType: req.CloseType,
		CloseMsg:  req.CloseMsg,
	})
	err = qs.metaCommonService.AddMeta(ctx, req.ID, entity.QuoteAuthorCloseReasonKey, string(closeMeta))
	if err != nil {
		return err
	}

	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           quoteAuthorInfo.UserID,
		ObjectID:         quoteAuthorInfo.ID,
		OriginalObjectID: quoteAuthorInfo.ID,
		ActivityTypeKey:  constant.ActQuoteAuthorClosed,
	})
	return nil
}

//// RemoveAnswer delete answer
//func (qs *QuoteAuthorCommon) RemoveAnswer(ctx context.Context, id string) (err error) {
//	answerinfo, has, err := qs.answerRepo.GetByID(ctx, id)
//	if err != nil {
//		return err
//	}
//	if !has {
//		return nil
//	}
//
//	// user add quoteAuthor count
//
//	err = qs.UpdateAnswerCount(ctx, answerinfo.QuoteAuthorID)
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

func (qs *QuoteAuthorCommon) SitemapCron(ctx context.Context) {
	quoteAuthorNum, err := qs.quoteAuthorRepo.GetQuoteAuthorCount(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	if quoteAuthorNum <= constant.SitemapMaxSize {
		_, err = qs.quoteAuthorRepo.SitemapQuoteAuthors(ctx, 1, int(quoteAuthorNum))
		if err != nil {
			log.Errorf("get site map quoteAuthor error: %v", err)
		}
		return
	}

	totalPages := int(math.Ceil(float64(quoteAuthorNum) / float64(constant.SitemapMaxSize)))
	for i := 1; i <= totalPages; i++ {
		_, err = qs.quoteAuthorRepo.SitemapQuoteAuthors(ctx, i, constant.SitemapMaxSize)
		if err != nil {
			log.Errorf("get site map quoteAuthor error: %v", err)
			return
		}
	}
}

func (qs *QuoteAuthorCommon) SetCache(ctx context.Context, cachekey string, info interface{}) error {
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

func (qs *QuoteAuthorCommon) ShowListFormat(ctx context.Context, data *entity.QuoteAuthor) *schema.QuoteAuthorInfoResp {
	return qs.ShowFormat(ctx, data)
}

func (qs *QuoteAuthorCommon) ShowFormat(ctx context.Context, data *entity.QuoteAuthor) *schema.QuoteAuthorInfoResp {
	info := schema.QuoteAuthorInfoResp{}
	info.ID = data.ID
	if handler.GetEnableShortID(ctx) {
		info.ID = uid.EnShortID(data.ID)
	}
	info.AuthorName = data.AuthorName
	info.UrlAuthorName = htmltext.UrlTitle(data.AuthorName)
	info.Content = data.Bio
	//info.HTML =  data.ParsedText
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
	info.QuoteAuthorUpdateTime = data.UpdatedAt.Unix()
	if data.UpdatedAt.Unix() < 1 {
		info.QuoteAuthorUpdateTime = 0
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
func (qs *QuoteAuthorCommon) ShowFormatWithTag(ctx context.Context, data *entity.QuoteAuthorWithTagsRevision) *schema.QuoteAuthorInfoResp {
	info := qs.ShowFormat(ctx, &data.QuoteAuthor)
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
