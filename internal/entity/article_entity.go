package entity

import "time"

// TaArticle 文章,参考wordpress的wp_articles表 ,ta_前缀表示table-article
// TaArticle 文章,参考wordpress的wp_articles表 ,ta_前缀表示table-article
type Article struct {
	ID             string    `json:"ID" gorm:"ID"`
	UserID         string    `json:"user_id" gorm:"user_id"`
	PostDate       time.Time `json:"post_date" gorm:"post_date"`               // 发布时间
	PostUpdateTime time.Time `json:"post_update_time" gorm:"post_update_time"` // 修改时间

	OriginalText string `xorm:"not null MEDIUMTEXT original_text"`
	ParsedText   string `xorm:"not null MEDIUMTEXT parsed_text"`

	//Content         string    `json:"content" gorm:"content"`
	Title           string    `json:"title" gorm:"title"`
	Excerpt         string    `json:"excerpt" gorm:"excerpt"` // 摘录
	Status          int       `json:"status" gorm:"status"`
	CommentStatus   int8      `json:"comment_status" gorm:"comment_status"` // 评论状态（open/closed）
	Password        string    `json:"password" gorm:"password"`             // 密码
	SlugName        string    `json:"slug_name" gorm:"slug_name"`           // 文章缩略名
	ContentFiltered string    `json:"content_filtered" gorm:"content_filtered"`
	MenuOrder       int64     `json:"menu_order" gorm:"menu_order"`       // 排序ID
	CommentCount    int64     `json:"comment_count" gorm:"comment_count"` // 评论总数
	CreatedAt       time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"updated_at"`
	Pin             int       `json:"pin" gorm:"pin"`
	Show            int       `json:"show" gorm:"show"`

	HotScore        int `xorm:"not null default 0 INT(11) hot_score"`
	CollectionCount int `xorm:"not null default 0 INT(11) collection_count"`
	FollowCount     int `xorm:"not null default 0 INT(11) follow_count"`

	//LastEditUserID string `xorm:"not null default 0 BIGINT(20) last_edit_user_id"`

	ViewCount       int `xorm:"not null default 0 INT(11) view_count"`
	UniqueViewCount int `xorm:"not null default 0 INT(11) unique_view_count"`
	VoteCount       int `xorm:"not null default 0 INT(11) vote_count"`

	RevisionID string `xorm:"not null default 0 BIGINT(20) revision_id"`
}

//post_date datetime NOT NULL  comment '发布时间',
//                        post_modified datetime  comment '修改时间',

// id这里改成string
// TableName 表名称

// TableName 表名称
func (Article) TableName() string {
	return "ta_article"
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
