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
//package quotecommon

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"time"

	"github.com/apache/incubator-answer/internal/base/constant"
	"github.com/apache/incubator-answer/internal/base/data"
	"github.com/apache/incubator-answer/internal/base/handler"
	"github.com/apache/incubator-answer/internal/base/reason"
	"github.com/apache/incubator-answer/internal/service/activity_common"
	"github.com/apache/incubator-answer/internal/service/activity_queue"
	"github.com/apache/incubator-answer/internal/service/config"
	metacommon "github.com/apache/incubator-answer/internal/service/meta_common"
	"github.com/apache/incubator-answer/internal/service/revision"
	"github.com/apache/incubator-answer/pkg/checker"
	"github.com/apache/incubator-answer/pkg/htmltext"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/segmentfault/pacman/errors"

	"github.com/apache/incubator-answer/internal/entity"
	"github.com/apache/incubator-answer/internal/schema"
	answercommon "github.com/apache/incubator-answer/internal/service/answer_common"
	collectioncommon "github.com/apache/incubator-answer/internal/service/collection_common"
	tagcommon "github.com/apache/incubator-answer/internal/service/tag_common"
	usercommon "github.com/apache/incubator-answer/internal/service/user_common"
	"github.com/segmentfault/pacman/log"
)

// QuoteRepo quote repository
type QuoteRepo interface {
	AddQuote(ctx context.Context, quote *entity.Quote) (err error)
	RemoveQuote(ctx context.Context, id string) (err error)
	UpdateQuote(ctx context.Context, quote *entity.Quote, Cols []string) (err error)
	GetQuote(ctx context.Context, id string) (quote *entity.Quote, exist bool, err error)
	GetQuoteList(ctx context.Context, quote *entity.Quote) (quotes []*entity.Quote, err error)
	GetQuotePage(ctx context.Context, page, pageSize int, tagIDs []string, userID, orderCond string, inDays int, showHidden, showPending bool) (
		quoteList []*entity.Quote, total int64, err error)
	GetRecommendQuotePageByTags(ctx context.Context, userID string, tagIDs, followedQuoteIDs []string, page, pageSize int) (quoteList []*entity.Quote, total int64, err error)
	UpdateQuoteStatus(ctx context.Context, quoteID string, status int) (err error)
	UpdateQuoteStatusWithOutUpdateTime(ctx context.Context, quote *entity.Quote) (err error)
	RecoverQuote(ctx context.Context, quoteID string) (err error)
	UpdateQuoteOperation(ctx context.Context, quote *entity.Quote) (err error)
	GetQuotesByTitle(ctx context.Context, title string, pageSize int) (quoteList []*entity.Quote, err error)
	UpdatePvCount(ctx context.Context, quoteID string) (err error)
	//UpdateAnswerCount(ctx context.Context, quoteID string, num int) (err error)
	UpdateCollectionCount(ctx context.Context, quoteID string) (count int64, err error)
	UpdateAccepted(ctx context.Context, quote *entity.Quote) (err error)
	UpdateLastAnswer(ctx context.Context, quote *entity.Quote) (err error)
	FindByID(ctx context.Context, id []string) (quoteList []*entity.Quote, err error)
	AdminQuotePage(ctx context.Context, search *schema.AdminQuotePageReq) ([]*entity.Quote, int64, error)
	GetQuoteCount(ctx context.Context) (count int64, err error)
	GetUserQuoteCount(ctx context.Context, userID string, show int) (count int64, err error)
	SitemapQuotes(ctx context.Context, page, pageSize int) (quoteIDList []*schema.SiteMapQuoteInfo, err error)
	RemoveAllUserQuote(ctx context.Context, userID string) (err error)
	UpdateSearch(ctx context.Context, quoteID string) (err error)
}

// QuoteCommon user service
type QuoteCommon struct {
	quoteRepo            QuoteRepo
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

	quoteAuthorRepo QuoteAuthorRepo
	quotePieceRepo  QuotePieceRepo
}

func NewQuoteCommon(quoteRepo QuoteRepo,
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

	quoteAuthorRepo QuoteAuthorRepo,
	quotePieceRepo QuotePieceRepo,
) *QuoteCommon {
	return &QuoteCommon{
		quoteRepo:            quoteRepo,
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

		quoteAuthorRepo: quoteAuthorRepo,
		quotePieceRepo:  quotePieceRepo,
	}
}

func (qs *QuoteCommon) GetUserQuoteCount(ctx context.Context, userID string) (count int64, err error) {
	return qs.quoteRepo.GetUserQuoteCount(ctx, userID, 0)
}

func (qs *QuoteCommon) GetPersonalUserQuoteCount(ctx context.Context, loginUserID, userID string, isAdmin bool) (count int64, err error) {
	show := entity.QuoteShow
	if loginUserID == userID || isAdmin {
		show = 0
	}
	return qs.quoteRepo.GetUserQuoteCount(ctx, userID, show)
}

func (qs *QuoteCommon) UpdatePv(ctx context.Context, quoteID string) error {
	return qs.quoteRepo.UpdatePvCount(ctx, quoteID)
}

//func (qs *QuoteCommon) UpdateAnswerCount(ctx context.Context, quoteID string) error {
//	count, err := qs.answerRepo.GetCountByQuoteID(ctx, quoteID)
//	if err != nil {
//		return err
//	}
//	if count == 0 {
//		err = qs.quoteRepo.UpdateLastAnswer(ctx, &entity.Quote{
//			ID:           quoteID,
//			LastAnswerID: "0",
//		})
//		if err != nil {
//			return err
//		}
//	}
//	return qs.quoteRepo.UpdateAnswerCount(ctx, quoteID, int(count))
//}

func (qs *QuoteCommon) UpdateCollectionCount(ctx context.Context, quoteID string) (count int64, err error) {
	return qs.quoteRepo.UpdateCollectionCount(ctx, quoteID)
}

func (qs *QuoteCommon) UpdateAccepted(ctx context.Context, quoteID, AnswerID string) error {
	quote := &entity.Quote{}
	quote.ID = quoteID
	//quote.AcceptedAnswerID = AnswerID
	return qs.quoteRepo.UpdateAccepted(ctx, quote)
}

func (qs *QuoteCommon) UpdateLastAnswer(ctx context.Context, quoteID, AnswerID string) error {
	quote := &entity.Quote{}
	quote.ID = quoteID
	//quote.LastAnswerID = AnswerID
	return qs.quoteRepo.UpdateLastAnswer(ctx, quote)
}

func (qs *QuoteCommon) UpdatePostTime(ctx context.Context, quoteID string) error {
	quoteinfo := &entity.Quote{}
	now := time.Now()
	_ = now
	quoteinfo.ID = quoteID
	quoteinfo.PostUpdateTime = now
	return qs.quoteRepo.UpdateQuote(ctx, quoteinfo, []string{"post_update_time"})
}
func (qs *QuoteCommon) UpdatePostSetTime(ctx context.Context, quoteID string, setTime time.Time) error {
	quoteinfo := &entity.Quote{}
	quoteinfo.ID = quoteID
	quoteinfo.PostUpdateTime = setTime
	return qs.quoteRepo.UpdateQuote(ctx, quoteinfo, []string{"post_update_time"})
}

func (qs *QuoteCommon) FindInfoByID(ctx context.Context, quoteIDs []string, loginUserID string) (map[string]*schema.QuoteInfoResp, error) {
	list := make(map[string]*schema.QuoteInfoResp)
	quoteList, err := qs.quoteRepo.FindByID(ctx, quoteIDs)
	if err != nil {
		return list, err
	}
	quotes, err := qs.FormatQuotes(ctx, quoteList, loginUserID)
	if err != nil {
		return list, err
	}
	for _, item := range quotes {
		list[item.ID] = item
	}
	return list, nil
}

func (qs *QuoteCommon) InviteUserInfo(ctx context.Context, quoteID string) (inviteList []*schema.UserBasicInfo, err error) {
	return
	//InviteUserInfo := make([]*schema.UserBasicInfo, 0)
	//dbinfo, has, err := qs.quoteRepo.GetQuote(ctx, quoteID)
	//if err != nil {
	//	return InviteUserInfo, err
	//}
	//if !has {
	//	return InviteUserInfo, errors.NotFound(reason.QuoteNotFound)
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

//包含author和作品

func (qs *QuoteCommon) Info(ctx context.Context, quoteID string, loginUserID string) (resp *schema.QuoteInfoResp, err error) {
	quoteInfo, has, err := qs.quoteRepo.GetQuote(ctx, quoteID)
	if err != nil {
		return resp, err
	}
	quoteInfo.ID = uid.DeShortID(quoteInfo.ID)
	if !has {
		return resp, errors.NotFound(reason.QuoteNotFound)
	}
	resp = qs.ShowFormat(ctx, quoteInfo)
	if resp.Status == entity.QuoteStatusClosed {
		metaInfo, err := qs.metaCommonService.GetMetaByObjectIdAndKey(ctx, quoteInfo.ID, entity.QuoteCloseReasonKey)
		if err != nil {
			log.Error(err)
		} else {
			closeMsg := &schema.CloseQuoteMeta{}
			err = json.Unmarshal([]byte(metaInfo.Value), closeMsg)
			if err != nil {
				log.Error("json.Unmarshal CloseQuoteMeta error", err.Error())
			} else {
				cfg, err := qs.configService.GetConfigByID(ctx, closeMsg.CloseType)
				if err != nil {
					log.Error("json.Unmarshal QuoteCloseJson error", err.Error())
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

	if resp.Status != entity.QuoteStatusDeleted {
		if resp.Tags, err = qs.tagCommon.GetObjectTag(ctx, quoteID); err != nil {
			return resp, err
		}
	} else {
		revisionInfo, exist, err := qs.revisionRepo.GetLastRevisionByObjectID(ctx, quoteID)
		if err != nil {
			log.Errorf("get revision error %s", err)
		}
		if exist {
			quoteWithTagsRevision := &entity.QuoteWithTagsRevision{}
			if err = json.Unmarshal([]byte(revisionInfo.Content), quoteWithTagsRevision); err != nil {
				log.Errorf("revision parsing error %s", err)
				return resp, nil
			}
			for _, tag := range quoteWithTagsRevision.Tags {
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
	if checker.IsNotZeroString(quoteInfo.UserID) {
		userIds = append(userIds, quoteInfo.UserID)
	}
	//if checker.IsNotZeroString(quoteInfo.LastEditUserID) {
	//	userIds = append(userIds, quoteInfo.LastEditUserID)
	//}
	if checker.IsNotZeroString(resp.LastAnsweredUserID) {
		userIds = append(userIds, resp.LastAnsweredUserID)
	}
	userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIds)
	if err != nil {
		return resp, err
	}
	resp.UserInfo = userInfoMap[quoteInfo.UserID]
	//resp.UpdateUserInfo = userInfoMap[quoteInfo.LastEditUserID]
	resp.LastAnsweredUserInfo = userInfoMap[resp.LastAnsweredUserID]
	if len(loginUserID) == 0 {
		return resp, nil
	}

	resp.VoteStatus = qs.voteRepo.GetVoteStatus(ctx, quoteID, loginUserID)
	resp.IsFollowed, _ = qs.followCommon.IsFollowed(ctx, loginUserID, quoteID)

	ids, err := qs.AnswerCommon.SearchAnswerIDs(ctx, loginUserID, quoteInfo.ID)
	if err != nil {
		log.Error("AnswerFunc.SearchAnswerIDs", err)
	}
	resp.Answered = len(ids) > 0
	if resp.Answered {
		resp.FirstAnswerId = ids[0]
	}

	collectedMap, err := qs.collectionCommon.SearchObjectCollected(ctx, loginUserID, []string{quoteInfo.ID})
	if err != nil {
		return nil, err
	}
	if len(collectedMap) > 0 {
		resp.Collected = true
	}
	return resp, nil
}

// \b 是一个特殊的元字符，表示单词边界。它匹配一个单词的开始或结束的位置
func RepImages(htmls string, needSize int) []string {
	var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	imgs := imgRE.FindAllStringSubmatch(htmls, needSize) //imgRE.FindAllStringSubmatch(htmls, -1)
	out := make([]string, len(imgs))
	for i := range out {
		out[i] = imgs[i][1]
		//fmt.Println(strconv.Itoa(i), out[i])
	}
	return out
}
func RegFindImages(htmls string, needSize int) []string {
	var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	imgs := imgRE.FindStringSubmatch(htmls) //imgRE.FindAllStringSubmatch(htmls, -1)
	//fmt.Printf("imgs:%v \n", imgs)
	log.Infof("RegFindImages   imgs:%v \n", imgs)

	out := make([]string, 0)

	if len(imgs) > 1 {
		out = append(out, imgs[1])
	}
	return out
}
func GetQuoteThumbnails(quoteInfo *entity.Quote) []schema.QuoteThumbnail {
	/*
		1. 首先是否设置了文章自己本身的特色图片。
		2. 如果没有，就检查下这篇文章的标签是否有特色图片。
		3. 如果没有，就检查下这篇文章是否有图片，有就获取第一张。
		4. 如果还是没有，就检查下这篇文章所在的分类是否有特色图片。
	*/
	thumbnails := make([]schema.QuoteThumbnail, 0)
	//if quoteInfo.Thumbnails != "" {
	//	err := json.Unmarshal([]byte(quoteInfo.Thumbnails), &thumbnails)
	//	if err != nil {
	//		log.Errorf("GetQuoteThumbnails err:%v \n", err)
	//	}
	//}
	//if len(thumbnails) != 0 {
	//	log.Infof("thumbnails1:%+v", thumbnails)
	//	return thumbnails
	//}

	images := RegFindImages(quoteInfo.ParsedText, 1)

	for _, val := range images {
		thumbnail := schema.QuoteThumbnail{
			Url: val,
		}
		thumbnails = append(thumbnails, thumbnail)
	}

	if len(thumbnails) != 0 {
		log.Infof("thumbnails2:%+v", thumbnails)
		return thumbnails
	}
	randomNum := rand.Intn(41) //[0,n)之
	imageUrl := fmt.Sprintf("%s%d%s", constant.OSS_BUCKET_NAME+"/images/postthumbnail/", randomNum, ".jpg")
	thumbnail := schema.QuoteThumbnail{
		Url: imageUrl,
	}

	thumbnails = append(thumbnails, thumbnail)
	log.Infof("thumbnails3:%+v", thumbnails)
	return thumbnails
}
func (qs *QuoteCommon) FormatQuotesPage(
	ctx context.Context, quoteList []*entity.Quote, loginUserID string, orderCond string) (
	formattedQuotes []*schema.QuotePageResp, err error) {
	formattedQuotes = make([]*schema.QuotePageResp, 0)
	quoteIDs := make([]string, 0)
	userIDs := make([]string, 0)

	for _, quoteInfo := range quoteList {

		//thumbnails := GetQuoteThumbnails(quoteInfo)

		t := &schema.QuotePageResp{
			ID:              quoteInfo.ID,
			CreatedAt:       quoteInfo.CreatedAt.Unix(),
			Title:           quoteInfo.Title,
			UrlTitle:        htmltext.UrlTitle(quoteInfo.Title),
			Description:     htmltext.FetchExcerpt(quoteInfo.ParsedText, "...", 80), //240),
			Status:          quoteInfo.Status,
			ViewCount:       quoteInfo.ViewCount,
			UniqueViewCount: quoteInfo.UniqueViewCount,
			VoteCount:       quoteInfo.VoteCount,
			//AnswerCount:      quoteInfo.AnswerCount,
			CollectionCount: quoteInfo.CollectionCount,
			FollowCount:     quoteInfo.FollowCount,
			//AcceptedAnswerID: quoteInfo.AcceptedAnswerID,
			//LastAnswerID: quoteInfo.LastAnswerID,
			Pin:  quoteInfo.Pin,
			Show: quoteInfo.Show,

			//Thumbnails: thumbnails,
			QuoteAuthorId: quoteInfo.QuoteAuthorId,
			QuotePieceId:  quoteInfo.QuotePieceId,
		}
		if t.Title == "" {
			t.Title = t.Description //@cws 如果title为空，则使用描述代替
		}

		quoteIDs = append(quoteIDs, quoteInfo.ID)
		userIDs = append(userIDs, quoteInfo.UserID)
		haveEdited, haveAnswered := false, false
		//if checker.IsNotZeroString(quoteInfo.LastEditUserID) {
		//	haveEdited = true
		//	userIDs = append(userIDs, quoteInfo.LastEditUserID)
		//}
		//if checker.IsNotZeroString(quoteInfo.LastAnswerID) {
		//	haveAnswered = true
		//
		//	answerInfo, exist, err := qs.answerRepo.GetAnswer(ctx, quoteInfo.LastAnswerID)
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
		log.Info("@quoteInfo.QuoteAuthorId:", quoteInfo.QuoteAuthorId)
		if checker.IsNotZeroString(quoteInfo.QuoteAuthorId) {
			quoteAuthorBaseInfoInfo, exist, err := qs.quoteAuthorRepo.GetQuoteAuthorSimple(ctx, quoteInfo.QuoteAuthorId)
			log.Infof("quoteAuthorBaseInfoInfo :%+v ,%+v", err, exist)
			if err == nil && exist {
				log.Info("has val autho basic info")
				t.QuoteAuthorBasicInfo = quoteAuthorBaseInfoInfo
			}
		}
		if checker.IsNotZeroString(quoteInfo.QuotePieceId) {
			quotePieceBaseInfo, exist, err := qs.quotePieceRepo.GetQuotePieceSimple(ctx, quoteInfo.QuotePieceId)
			if err == nil && exist {
				t.QuotePieceBasicInfo = quotePieceBaseInfo
			}
		}

		// if order condition is newest or nobody edited or nobody answered, only show quote author
		if orderCond == schema.QuoteOrderCondNewest || (!haveEdited && !haveAnswered) {
			t.OperationType = schema.QuotePageRespOperationTypeAsked
			t.OperatedAt = quoteInfo.CreatedAt.Unix()
			t.Operator = &schema.QuotePageRespOperator{ID: quoteInfo.UserID}
		} else {
			// if no one
			if haveEdited {
				t.OperationType = schema.QuotePageRespOperationTypeModified
				t.OperatedAt = quoteInfo.UpdatedAt.Unix()
				//t.Operator = &schema.QuotePageRespOperator{ID: quoteInfo.LastEditUserID}
			}

			if haveAnswered {
				if t.LastAnsweredAt.Unix() > t.OperatedAt {
					t.OperationType = schema.QuotePageRespOperationTypeAnswered
					t.OperatedAt = t.LastAnsweredAt.Unix()
					t.Operator = &schema.QuotePageRespOperator{ID: t.LastAnsweredUserID}
				}
			}
		}
		formattedQuotes = append(formattedQuotes, t)
	}

	tagsMap, err := qs.tagCommon.BatchGetObjectTag(ctx, quoteIDs)
	if err != nil {
		return formattedQuotes, err
	}
	userInfoMap, err := qs.userCommon.BatchUserBasicInfoByID(ctx, userIDs)
	if err != nil {
		return formattedQuotes, err
	}

	for _, item := range formattedQuotes {
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
	return formattedQuotes, nil
}

func (qs *QuoteCommon) FormatQuotes(ctx context.Context, quoteList []*entity.Quote, loginUserID string) ([]*schema.QuoteInfoResp, error) {
	list := make([]*schema.QuoteInfoResp, 0)
	objectIds := make([]string, 0)
	userIds := make([]string, 0)

	for _, quoteInfo := range quoteList {
		item := qs.ShowFormat(ctx, quoteInfo)
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

// RemoveQuote delete quote
func (qs *QuoteCommon) RemoveQuote(ctx context.Context, req *schema.RemoveQuoteReq) (err error) {
	quoteInfo, has, err := qs.quoteRepo.GetQuote(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	if quoteInfo.Status == entity.QuoteStatusDeleted {
		return nil
	}

	quoteInfo.Status = entity.QuoteStatusDeleted
	err = qs.quoteRepo.UpdateQuoteStatus(ctx, quoteInfo.ID, quoteInfo.Status)
	if err != nil {
		return err
	}

	userQuoteCount, err := qs.GetUserQuoteCount(ctx, quoteInfo.UserID)
	if err != nil {
		log.Error("user GetUserQuoteCount error", err.Error())
	} else {
		//@ms:TODO	err = qs.userCommon.UpdateQuoteCount(ctx, quoteInfo.UserID, userQuoteCount)
		_ = userQuoteCount
		if err != nil {
			log.Error("user IncreaseQuoteCount error", err.Error())
		}
	}

	return nil
}

func (qs *QuoteCommon) CloseQuote(ctx context.Context, req *schema.CloseQuoteReq) error {
	quoteInfo, has, err := qs.quoteRepo.GetQuote(ctx, req.ID)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	quoteInfo.Status = entity.QuoteStatusClosed
	err = qs.quoteRepo.UpdateQuoteStatus(ctx, quoteInfo.ID, quoteInfo.Status)
	if err != nil {
		return err
	}

	closeMeta, _ := json.Marshal(schema.CloseQuoteMeta{
		CloseType: req.CloseType,
		CloseMsg:  req.CloseMsg,
	})
	err = qs.metaCommonService.AddMeta(ctx, req.ID, entity.QuoteCloseReasonKey, string(closeMeta))
	if err != nil {
		return err
	}

	qs.activityQueueService.Send(ctx, &schema.ActivityMsg{
		UserID:           quoteInfo.UserID,
		ObjectID:         quoteInfo.ID,
		OriginalObjectID: quoteInfo.ID,
		ActivityTypeKey:  constant.ActQuoteClosed,
	})
	return nil
}

//// RemoveAnswer delete answer
//func (qs *QuoteCommon) RemoveAnswer(ctx context.Context, id string) (err error) {
//	answerinfo, has, err := qs.answerRepo.GetByID(ctx, id)
//	if err != nil {
//		return err
//	}
//	if !has {
//		return nil
//	}
//
//	// user add quote count
//
//	err = qs.UpdateAnswerCount(ctx, answerinfo.QuoteID)
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

func (qs *QuoteCommon) SitemapCron(ctx context.Context) {
	quoteNum, err := qs.quoteRepo.GetQuoteCount(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	if quoteNum <= constant.SitemapMaxSize {
		_, err = qs.quoteRepo.SitemapQuotes(ctx, 1, int(quoteNum))
		if err != nil {
			log.Errorf("get site map quote error: %v", err)
		}
		return
	}

	totalPages := int(math.Ceil(float64(quoteNum) / float64(constant.SitemapMaxSize)))
	for i := 1; i <= totalPages; i++ {
		_, err = qs.quoteRepo.SitemapQuotes(ctx, i, constant.SitemapMaxSize)
		if err != nil {
			log.Errorf("get site map quote error: %v", err)
			return
		}
	}
}

func (qs *QuoteCommon) SetCache(ctx context.Context, cachekey string, info interface{}) error {
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

func (qs *QuoteCommon) ShowListFormat(ctx context.Context, data *entity.Quote) *schema.QuoteInfoResp {
	return qs.ShowFormat(ctx, data)
}

func (qs *QuoteCommon) ShowFormat(ctx context.Context, data *entity.Quote) *schema.QuoteInfoResp {
	info := schema.QuoteInfoResp{}
	info.ID = data.ID
	if handler.GetEnableShortID(ctx) {
		info.ID = uid.EnShortID(data.ID)
	}
	info.Title = data.Title
	info.UrlTitle = htmltext.UrlTitle(data.Title)
	info.Content = data.OriginalText
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
	info.PostUpdateTime = data.PostUpdateTime.Unix()
	if data.PostUpdateTime.Unix() < 1 {
		info.PostUpdateTime = 0
	}
	info.QuoteUpdateTime = data.UpdatedAt.Unix()
	if data.UpdatedAt.Unix() < 1 {
		info.QuoteUpdateTime = 0
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

	info.QuoteAuthorId = data.QuoteAuthorId
	info.QuotePieceId = data.QuotePieceId

	return &info
}
func (qs *QuoteCommon) ShowFormatWithTag(ctx context.Context, data *entity.QuoteWithTagsRevision) *schema.QuoteInfoResp {
	info := qs.ShowFormat(ctx, &data.Quote)
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
