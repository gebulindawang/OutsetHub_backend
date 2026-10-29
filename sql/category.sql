/*
 起点社区 - 分类表
 Schema: outsethub
 注意:本脚本会 DROP TABLE,执行前确认已有数据是否需要保留
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
  `id`         int         NOT NULL AUTO_INCREMENT,
  `name`       varchar(50) NOT NULL                COMMENT '分类名,全局唯一',
  `sort`       int         NOT NULL DEFAULT 0      COMMENT '排序权重,越小越靠前',
  `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_sort` (`sort`)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- 初始板块
INSERT INTO `category` (`name`, `sort`) VALUES
  ('前端',   10),
  ('后端',   20),
  ('数据库', 30),
  ('算法',   40),
  ('闲聊',   50);

SET FOREIGN_KEY_CHECKS = 1;
