
drop table if exists ta_article;
create table ta_article (
                            ID bigint unsigned NOT NULL ,
                            user_id bigint(20) unsigned NOT NULL default '0' comment '作者ID',



                            post_date datetime NOT NULL  comment '发布时间',
                            post_update_time datetime  comment '文章修改时间',

                             `original_text` mediumtext COLLATE utf8mb4_unicode_ci NOT NULL,
                          `parsed_text` mediumtext COLLATE utf8mb4_unicode_ci NOT NULL,


                            title text NOT NULL,
                            excerpt text NOT NULL comment '摘录',
                            status tinyint(2) NOT NULL default '0',
                            comment_status  tinyint(2) NOT NULL default '0' comment '评论状态（open/closed）',

                            password varchar(255) NOT NULL default '' comment '密码',
                            slug_name varchar(200) NOT NULL default '' comment '文章缩略名',

                            content_filtered longtext NOT NULL,
                            menu_order int(11) NOT NULL default '0' comment '排序ID',
                            comment_count bigint(20) NOT NULL default '0' comment '评论总数',

  `revision_id` bigint NOT NULL DEFAULT '0' comment '修订号',


                            `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `updated_at` timestamp not null  DEFAULT   CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,


      `pin` int NOT NULL DEFAULT '1',
  `show` int NOT NULL DEFAULT '1',
 `collection_count` int NOT NULL DEFAULT '0',
  `follow_count` int NOT NULL DEFAULT '0',
  `view_count` int NOT NULL DEFAULT '0',
      `hot_score` int NOT NULL DEFAULT '0',
      `unique_view_count` int NOT NULL DEFAULT '0',
  `vote_count` int NOT NULL DEFAULT '0',

  `thumbnails` varchar(255) comment '文章缩略图',


                            PRIMARY KEY  (ID),
                            KEY     idx_title (    title(191)),
                            KEY idx_status_date (    status,    post_date),
                            KEY     idx_user_id (    user_id)
) comment  "文章,参考wordpress的wp_articles表 ,ta_前缀表示table-article";

alter table ta_article add  column `view_count` int NOT NULL DEFAULT '0',

alter table ta_article add  column `hot_score` int NOT NULL DEFAULT '0',
alter table ta_article add  column `unique_view_count` int NOT NULL DEFAULT '0',
alter table ta_article add  column `vote_count` int NOT NULL DEFAULT '0'


	KEY article_parent (article_parent),


-- alter table ta_article add  column `thumbnails` json comment '文章缩略图';
alter table ta_article  modify  column `thumbnails` varchar(256)  not null default '' comment '文章缩略图';

alter table ta_article add  column `original_text_format` tinyint NOT NULL DEFAULT '0' comment 'text原始格式,0:markdown,1:html';


	-- article_type varchar(20) NOT NULL default 'article' comment '：文章类型（post/page等）',
	-- article_mime_type varchar(100) NOT NULL default '',

	-- guid varchar(255) NOT NULL default '' comment 'GUID是文章或页面的唯一标识符，通常用于RSS馈送。',
	-- article_parent bigint(20) unsigned NOT NULL default '0', post_parent：父文章，主要用于PAGE

--
--
-- CREATE TABLE `ta_tag` (
--   `id` bigint NOT NULL COMMENT 'tag_id',
--   `created_at` timestamp NULL DEFAULT NULL,
--   `updated_at` timestamp NULL DEFAULT NULL,
--   `main_tag_id` bigint NOT NULL DEFAULT '0',
--   `main_tag_slug_name` varchar(35) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
--   `slug_name` varchar(35) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
--   `display_name` varchar(35) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
--   `original_text` mediumtext COLLATE utf8mb4_unicode_ci NOT NULL,
--   `parsed_text` mediumtext COLLATE utf8mb4_unicode_ci NOT NULL,
--   `article_count` int NOT NULL DEFAULT '0',
--   `status` int NOT NULL DEFAULT '1',
--   PRIMARY KEY (`id`),
--   UNIQUE KEY `UQE_tag_slug_name` (`slug_name`)
-- ) comment 'tag表,参考answer的tag';
--
-- create table ta_tag_relation (
--     ID bigint(20) unsigned NOT NULL auto_increment,
--     article_id bigint(20) unsigned NOT NULL default '0',
--     tag_id bigint(20) unsigned NOT NULL default '0',
--     `created_at` timestamp NULL DEFAULT NULL,
--    `updated_at` timestamp NULL DEFAULT NULL,
-- ) comment "文章tag关系";
--

	tag复用apache answer之前的

//@ms: 参考 answer的 comment表
CREATE TABLE `ta_comment` (
                           `id` bigint NOT NULL AUTO_INCREMENT,
                           `created_at` timestamp NULL DEFAULT NULL,
                           `updated_at` timestamp NULL DEFAULT NULL,
                           `user_id` bigint NOT NULL DEFAULT '0',
                           `reply_user_id` bigint DEFAULT NULL,
                           `reply_comment_id` bigint DEFAULT NULL,
                           `object_id` bigint NOT NULL DEFAULT '0',
                           `article_id` bigint NOT NULL DEFAULT '0',
                           `vote_count` int NOT NULL DEFAULT '0',
                           `status` tinyint NOT NULL DEFAULT '0',
                           `original_text` mediumtext COLLATE utf8mb4_unicode_ci NOT NULL,
                           `parsed_text` mediumtext COLLATE utf8mb4_unicode_ci NOT NULL,
                           PRIMARY KEY (`id`),
                           KEY `IDX_comment_object_id` (`object_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

alter table user add  column `article_count` int NOT NULL DEFAULT '0' comment '文章数';


INSERT INTO `config` ( `key`, `value`) VALUES ('rank.article.add', '1');
INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.audit', '1');
INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.close', '-1');
INSERT INTO `config` ( `key`, `value`) VALUES ('rank.article.delete', '-1');
INSERT INTO `config` ( `key`, `value`) VALUES ('rank.article.edit', '1');
INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.edit_without_review', '1');
INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.hide', '-1');
INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.pin', '-1');
INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.reopen', '-1');
INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.show', '-1');
INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.undeleted', '-1');
INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.unpin', '-1');
INSERT INTO `config` ( `key`, `value`) VALUES ('rank.article.vote_down', '125');
INSERT INTO `config` ( `key`, `value`) VALUES ('rank.article.vote_up', '1');


再原本的tag上增加tag_type类型，和 order_

alter table tag add  column  tag_type tinyint NOT NULL DEFAULT '0' comment '0:默认,1:Article类型的';
alter table tag add  column  tag_sort int NOT NULL DEFAULT '0' comment 'tag的顺序，越小越靠前';

