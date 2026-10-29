/*
 起点社区 - 用户表
 Schema: outsethub
 注意:本脚本会 DROP TABLE,已有账号会丢失

 若不想丢数据,改用下面的 ALTER 方式:
   ALTER TABLE `user`
     MODIFY `ID` int NOT NULL AUTO_INCREMENT,
     MODIFY `username` varchar(50) NOT NULL,
     MODIFY `password` varchar(255) NOT NULL,
     MODIFY `nickname` varchar(50) NOT NULL DEFAULT '',
     MODIFY `role` tinyint NOT NULL DEFAULT 0,
     ADD UNIQUE KEY `uk_username` (`username`),
     ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
     ADD COLUMN `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;
   -- 主键的 DESC 需先查看 SHOW CREATE TABLE `user` 再决定是否重建
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id`         int         NOT NULL AUTO_INCREMENT,
  `username`   varchar(50) NOT NULL                 COMMENT '登录名,唯一',
  `password`   varchar(255) NOT NULL                COMMENT 'bcrypt 哈希,固定 60 字符',
  `nickname`   varchar(50) NOT NULL DEFAULT ''      COMMENT '展示名',
  `role`       tinyint     NOT NULL DEFAULT 0       COMMENT '0普通用户 1管理员',
  `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
