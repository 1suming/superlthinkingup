create table ta_articles (
    ID bigint(20) unsigned NOT NULL auto_increment,
	article_author bigint(20) unsigned NOT NULL default '0',
	article_date datetime NOT NULL  ,
	article_content longtext NOT NULL,
	article_title text NOT NULL,
	article_excerpt text NOT NULL,
	article_status tinyint(2) NOT NULL default '0',
	comment_status  tinyint(2) NOT NULL default '0',
 
 	article_password varchar(255) NOT NULL default '',
	article_name varchar(200) NOT NULL default '',

	article_modified datetime NOT NULL  ,
	article_modified_gmt datetime NOT NULL  ,
	article_content_filtered longtext NOT NULL,
	menu_order int(11) NOT NULL default '0' comment '排序ID',
	comment_count bigint(20) NOT NULL default '0',


    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp not null  DEFAULT   CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,


	PRIMARY KEY  (ID),
	KEY article_name (article_name(191)),
	KEY idx_status_date (article_status,article_date),
	KEY article_author (article_author)
) comment  "文章,参考wordpress的wp_articles表 ,ta_前缀表示table-article";

	KEY article_parent (article_parent),

	-- article_type varchar(20) NOT NULL default 'article' comment '：文章类型（post/page等）',
	-- article_mime_type varchar(100) NOT NULL default '',

	-- guid varchar(255) NOT NULL default '' comment 'GUID是文章或页面的唯一标识符，通常用于RSS馈送。',
	-- article_parent bigint(20) unsigned NOT NULL default '0', post_parent：父文章，主要用于PAGE



CREATE TABLE `ta_tag` (
  `id` bigint NOT NULL COMMENT 'tag_id',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `main_tag_id` bigint NOT NULL DEFAULT '0',
  `main_tag_slug_name` varchar(35) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `slug_name` varchar(35) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `display_name` varchar(35) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `original_text` mediumtext COLLATE utf8mb4_unicode_ci NOT NULL,
  `parsed_text` mediumtext COLLATE utf8mb4_unicode_ci NOT NULL,
  `article_count` int NOT NULL DEFAULT '0',
  `status` int NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UQE_tag_slug_name` (`slug_name`)
) comment 'tag表,参考answer的tag';

create table ta_tag_relation (
    ID bigint(20) unsigned NOT NULL auto_increment,
    article_id bigint(20) unsigned NOT NULL default '0',
    tag_id bigint(20) unsigned NOT NULL default '0',
    `created_at` timestamp NULL DEFAULT NULL,
   `updated_at` timestamp NULL DEFAULT NULL,
) comment "文章tag关系";