package entity

import "time"

// TaArticle 文章,参考wordpress的wp_articles表 ,ta_前缀表示table-article
// TaArticle 文章,参考wordpress的wp_articles表 ,ta_前缀表示table-article
// error: ID condition is error, expect 0 primarykeys, there are 1
//为你没有定义Id为主键，所以ID()函数不认

type Article struct {
	ID             string    `json:"ID"    xorm:"not null pk BIGINT(20) id" `
	UserID         string    `json:"user_id" xorm:"user_id"`
	PostDate       time.Time `json:"post_date" xorm:"post_date"`               // 发布时间
	PostUpdateTime time.Time `json:"post_update_time" xorm:"post_update_time"` // 修改时间

	OriginalText string `xorm:"not null MEDIUMTEXT original_text"`
	ParsedText   string `xorm:"not null MEDIUMTEXT parsed_text"`

	//Content         string    `json:"content" xorm:"content"`
	Title           string    `json:"title" xorm:"title"`
	Excerpt         string    `json:"excerpt" xorm:"excerpt"` // 摘录
	Status          int       `json:"status" xorm:"status"`
	CommentStatus   int8      `json:"comment_status" xorm:"comment_status"` // 评论状态（open/closed）
	Password        string    `json:"password" xorm:"password"`             // 密码
	SlugName        string    `json:"slug_name" xorm:"slug_name"`           // 文章缩略名
	ContentFiltered string    `json:"content_filtered" xorm:"content_filtered"`
	MenuOrder       int64     `json:"menu_order" xorm:"menu_order"`       // 排序ID
	CommentCount    int64     `json:"comment_count" xorm:"comment_count"` // 评论总数
	CreatedAt       time.Time `json:"created_at" xorm:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" xorm:"updated_at"`
	Pin             int       `json:"pin" xorm:"pin"`
	Show            int       `json:"show" xorm:"show"`

	HotScore        int `xorm:"not null default 0 INT(11) hot_score"`
	CollectionCount int `xorm:"not null default 0 INT(11) collection_count"`
	FollowCount     int `xorm:"not null default 0 INT(11) follow_count"`

	//LastEditUserID string `xorm:"not null default 0 BIGINT(20) last_edit_user_id"`

	ViewCount       int `xorm:"not null default 0 INT(11) view_count"`
	UniqueViewCount int `xorm:"not null default 0 INT(11) unique_view_count"`
	VoteCount       int `xorm:"not null default 0 INT(11) vote_count"`

	RevisionID string `xorm:"not null default 0 BIGINT(20) revision_id"`

	//缩略图
	Thumbnails string `json:"thumbnails" xorm:"thumbnails"`

	OriginalTextFormat int8 `json:"original_text_format" xorm:"original_text_format"`
}

//post_date datetime NOT NULL  comment '发布时间',
//                        post_modified datetime  comment '修改时间',

// id这里改成string
// TableName 表名称

// TableName 表名称
const ARTICLE_TABLE_NAME = "ta_article"

func (Article) TableName() string {
	return ARTICLE_TABLE_NAME
}
func (Article) KeyName() string {
	return "article"
}

const (
	ArticleStatusAvailable = 1
	ArticleStatusClosed    = 2
	ArticleStatusDeleted   = 10
	ArticleStatusPending   = 11
	ArticleUnPin           = 1
	ArticlePin             = 2
	ArticleShow            = 1
	ArticleHide            = 2
)

var AdminArticleSearchStatus = map[string]int{
	"available": ArticleStatusAvailable,
	"closed":    ArticleStatusClosed,
	"deleted":   ArticleStatusDeleted,
	"pending":   ArticleStatusPending,
}

var AdminArticleSearchStatusIntToString = map[int]string{
	ArticleStatusAvailable: "available",
	ArticleStatusClosed:    "closed",
	ArticleStatusDeleted:   "deleted",
	ArticleStatusPending:   "pending",
}

// ArticleWithTagsRevision question
type ArticleWithTagsRevision struct {
	Article
	Tags []*TagSimpleInfoForRevision `json:"tags"`
}

// TagSimpleInfoForRevision tag simple info for revision
//type TagSimpleInfoForRevision struct {
//	ID              string `xorm:"not null pk comment('tag_id') BIGINT(20) id"`
//	MainTagID       int64  `xorm:"not null default 0 BIGINT(20) main_tag_id"`
//	MainTagSlugName string `xorm:"not null default '' VARCHAR(35) main_tag_slug_name"`
//	SlugName        string `xorm:"not null default '' unique VARCHAR(35) slug_name"`
//	DisplayName     string `xorm:"not null default '' VARCHAR(35) display_name"`
//	Recommend       bool   `xorm:"not null default false BOOL recommend"`
//	Reserved        bool   `xorm:"not null default false BOOL reserved"`
//	RevisionID      string `xorm:"not null default 0 BIGINT(20) revision_id"`
//}
// QuestionWithTagsRevision question
