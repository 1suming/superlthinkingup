package constant

const OSS_BUCKET_NAME = "https://superthinking.oss-cn-shenzhen.aliyuncs.com"

// /images/postthumbnail/1.jpg"
const (
	TAG_TYPE_ARTICLE int8 = 1
)
const (
	DeletedArticleTitleTrKey = "article.deleted_title"
	ArticlesTitleTrKey       = "article.articles_title"
	//TagsListTitleTrKey        = "tag.tags_title"
	//TagHasNoDescription       = "tag.no_description"
)

// //父级id为0的。//为-2表示 选择 parent_tag_id为0的
const TAG_PARENT_TAG_ID_IS_ZERO = -2
