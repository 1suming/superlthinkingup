
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


INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.add');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.edit');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.edit_without_review');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.delete');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.close');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.reopen');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.vote_up');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.vote_down');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.audit');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.pin');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.hide');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.unpin');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.show');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (2, 'article.undeleted');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.add');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.edit');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.edit_without_review');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.delete');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.close');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.reopen');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.vote_up');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.vote_down');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.audit');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.pin');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.hide');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.unpin');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.show');
INSERT INTO role_power_rel (`role_id`, `power_type`) VALUES (3, 'article.undeleted');


INSERT INTO config (`key`, `value`) VALUES ('article.voted_up', '10');
INSERT INTO config (`key`, `value`) VALUES ('article.voted_down_cancel', '2');
INSERT INTO config (`key`, `value`) VALUES ('article.vote_down_cancel', '1');
INSERT INTO config (`key`, `value`) VALUES ('article.voted_down', '-2');
INSERT INTO config (`key`, `value`) VALUES ('article.voted_up_cancel', '-10');
INSERT INTO config (`key`, `value`) VALUES ('article.vote_down', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.vote_up', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.vote_up_cancel', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.follow', '0');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.add', '1');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.edit', '1');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.delete', '-1');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.vote_up', '1');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.vote_down', '125');
INSERT INTO config (`key`, `value`) VALUES ('article.flag.reasons', '[\"reason.spam\",\"reason.rude_or_abusive\",\"reason.something\",\"reason.a_duplicate\"]');
INSERT INTO config (`key`, `value`) VALUES ('article.close.reasons', '[\"reason.a_duplicate\",\"reason.community_specific\",\"reason.not_clarity\",\"reason.something\"]');
INSERT INTO config (`key`, `value`) VALUES ('article.status.reasons', '[\"reason.normal\",\"reason.closed\",\"reason.deleted\"]');
INSERT INTO config (`key`, `value`) VALUES ('article.review.reasons', '[\"reason.looks_ok\",\"reason.needs_edit\",\"reason.needs_close\",\"reason.needs_delete\"]');
INSERT INTO config (`key`, `value`) VALUES ('article.asked', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.closed', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.reopened', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.answered', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.commented', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.accept', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.edited', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.rollback', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.deleted', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.undeleted', '0');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.edit_without_review', '1');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.audit', '1');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.close', '-1');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.reopen', '-1');
INSERT INTO config (`key`, `value`) VALUES ('article.pin', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.unpin', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.show', '0');
INSERT INTO config (`key`, `value`) VALUES ('article.hide', '0');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.pin', '-1');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.unpin', '-1');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.show', '-1');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.hide', '-1');
INSERT INTO config (`key`, `value`) VALUES ('rank.article.undeleted', '-1');




-- INSERT INTO `config` ( `key`, `value`) VALUES ('rank.article.add', '1');
-- INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.audit', '1');
-- INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.close', '-1');
-- INSERT INTO `config` ( `key`, `value`) VALUES ('rank.article.delete', '-1');
-- INSERT INTO `config` ( `key`, `value`) VALUES ('rank.article.edit', '1');
-- INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.edit_without_review', '1');
-- INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.hide', '-1');
-- INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.pin', '-1');
-- INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.reopen', '-1');
-- INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.show', '-1');
-- INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.undeleted', '-1');
-- INSERT INTO `config` ( `key`, `value`) VALUES ( 'rank.article.unpin', '-1');
-- INSERT INTO `config` ( `key`, `value`) VALUES ('rank.article.vote_down', '125');
-- INSERT INTO `config` ( `key`, `value`) VALUES ('rank.article.vote_up', '1');


再原本的tag上增加tag_type类型，和 order_

alter table tag add  column  tag_type tinyint NOT NULL DEFAULT '0' comment '0:默认,1:Article类型的';
alter table tag add  column  tag_sort int NOT NULL DEFAULT '0' comment 'tag的顺序，越小越靠前';

alter table site_info add  column remark varchar(255) NOT NULL default ''   comment '';
INSERT INTO `site_info` (  `created_at`, `updated_at`, `type`, `content`, `status`, `remark`) VALUES (  NULL, NULL, 'site_about_info', '<p>超维社，我们坚信：<strong>超级思维改变人生</strong>。在这里，你将遇到众多志同道合的人，一起交流。&nbsp;</p>', 1, '');
INSERT INTO `site_info` (  `created_at`, `updated_at`, `type`, `content`, `status`, `remark`) VALUES (  NULL, NULL, 'site_disclaim_info', '部分文章取自网络，侵权请留言或发邮件到此邮箱：admin@superthinkingup.com。', 1, '');
INSERT INTO `site_info` (  `created_at`, `updated_at`, `type`, `content`, `status`, `remark`) VALUES (  NULL, NULL, 'site_contact_info', '超维社', 1, '');
