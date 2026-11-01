/*
 起点社区 - 评论表
 Schema: outsethub
 注意:本脚本会 DROP TABLE,执行前确认已有数据是否需要保留
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
  `id`               bigint unsigned NOT NULL AUTO_INCREMENT,
  `article_id`       bigint unsigned NOT NULL             COMMENT '所属文章 ID,对应 article.id',
  `user_id`          int             NOT NULL             COMMENT '评论人 ID,对应 user.id',
  `parent_id`        bigint unsigned NOT NULL DEFAULT 0   COMMENT '顶级评论为 0;楼中楼回复指向顶级评论 ID,固定两级不做无限嵌套',
  `reply_to_user_id` int             NOT NULL DEFAULT 0   COMMENT '被回复人 ID,用于楼中楼里显示"回复 @某人";回复顶级评论时为 0',
  `content`          varchar(500)    NOT NULL             COMMENT '评论内容',
  `like_count`       int unsigned    NOT NULL DEFAULT 0   COMMENT '点赞数,冗余字段',
  `status`           tinyint         NOT NULL DEFAULT 1   COMMENT '0隐藏 1正常 2已删除',
  `created_at`       datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_article_created` (`article_id`, `created_at`) COMMENT '文章内按时间排序',
  KEY `idx_parent` (`parent_id`)                         COMMENT '楼中楼拉取子回复',
  KEY `idx_user` (`user_id`)                             COMMENT '我的评论'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- 测试数据:2 条顶级评论 + 楼中楼回复,演示 parent_id / reply_to_user_id 的用法
-- user_id 对应 user 表,article_id 对应 article 表
INSERT INTO `comment` (`article_id`, `user_id`, `parent_id`, `reply_to_user_id`, `content`, `like_count`, `created_at`) VALUES
(1, 2, 0, 0, '写得很清楚,补充一点:ref 包对象时内部其实也是 Proxy,只是要求 .value 整体替换。', 12, NOW() - INTERVAL 1 HOUR),
(1, 1, 1, 2, '对,所以我的习惯是组件内一律 ref,需要替换就整体换,不混用。', 3, NOW() - INTERVAL 30 MINUTE),
(2, 2, 0, 0, '三个页面就开始重复请求了,直接上 Pinia 吧,越拖越疼。', 8, NOW() - INTERVAL 2 HOUR),
(2, 1, 3, 2, '主要顾虑是上了之后滥用,什么状态都往 store 里塞。', 2, NOW() - INTERVAL 90 MINUTE),
(3, 1, 0, 0, 'ShouldBindJSON 漏 return 那个我上周刚踩过,症状就是响应变成非法 JSON。', 15, NOW() - INTERVAL 5 HOUR),
(5, 1, 0, 0, 'EXPLAIN 是真好用,type 从 ALL 变成 ref 的那一刻很有成就感。', 21, NOW() - INTERVAL 1 DAY),
(6, 2, 0, 0, '左闭右开记一套就够 +1,两种写法来回切才是 bug 之源。', 9, NOW() - INTERVAL 2 DAY);

SET FOREIGN_KEY_CHECKS = 1;
