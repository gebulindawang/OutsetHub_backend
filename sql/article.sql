/*
 起点社区 - 文章表
 Schema: outsethub
 注意:本脚本会 DROP TABLE,执行前确认已有数据是否需要保留
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `article`;
CREATE TABLE `article` (
  `id`            bigint unsigned NOT NULL AUTO_INCREMENT,
  `title`         varchar(200)    NOT NULL                COMMENT '标题',
  `summary`       varchar(300)    NOT NULL DEFAULT ''     COMMENT '列表页摘要,留空则由正文截取',
  `content`       mediumtext      NOT NULL                COMMENT '正文',
  `user_id`       int             NOT NULL                COMMENT '作者 ID,对应 user.id',
  `category_id`   int             NOT NULL                COMMENT '分类 ID,对应 category.id',
  `view_count`    int unsigned    NOT NULL DEFAULT 0      COMMENT '浏览数',
  `like_count`    int unsigned    NOT NULL DEFAULT 0      COMMENT '点赞数,冗余字段',
  `comment_count` int unsigned    NOT NULL DEFAULT 0      COMMENT '评论数,冗余字段',
  `status`        tinyint         NOT NULL DEFAULT 1      COMMENT '0草稿 1已发布 2已删除',
  `created_at`    datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`    datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_cat_created` (`category_id`, `created_at` DESC) COMMENT '板块内按时间排序',
  KEY `idx_user` (`user_id`)                               COMMENT '我的文章',
  KEY `idx_created` (`created_at` DESC)                    COMMENT '全站最新 / 热榜'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- 测试数据
-- user_id 需要对应 user 表里真实存在的 id,先执行 SELECT id, username FROM `user`;
-- category_id: 1前端 2后端 3数据库 4算法 5闲聊
-- ----------------------------
INSERT INTO `article`
  (`title`, `summary`, `content`, `user_id`, `category_id`,
   `view_count`, `like_count`, `comment_count`, `status`, `created_at`)
VALUES
(
  'Vue3 的 ref 和 reactive 到底该怎么选',
  '一开始我也分不清,直到踩了把 ref 直接传给 axios 导致循环引用这个坑。这篇把我理解的边界条件写清楚。',
  '## 先说结论\n\n原始类型只能用 ref,对象两种都行。但 ref 在模板里会自动解包,reactive 不会,这是最容易混的地方。\n\n## 一个真实的坑\n\n把 ref 对象直接传给 axios,JSON.stringify 会走到 dep 上形成循环引用,报 Converting circular structure to JSON。这个异常在发请求之前就抛了,前端只能看到一个不存在的网络错误。\n\n## 我的习惯\n\n组件内部的状态一律 ref,需要整体替换时直接 user.value = {...};跨组件共享的数据放 Pinia,store 里返回的解构要用 storeToRefs。',
  1, 1, 1284, 96, 23, 1, NOW() - INTERVAL 2 HOUR
),
(
  '论坛项目做到一半,该不该上 Pinia',
  '现在是 localStorage 存 token、每个页面自己 ref 存用户信息。三个页面共用数据时开始重复请求了。',
  '现在的状态很简单:token 放 localStorage,每个页面进来自查一次用户信息。\n\n问题是页面一多,同一个 /user/me 会被打三次。想听听大家的实际经验,是继续用组合式函数抽一个 useUser,还是直接上 Pinia。',
  1, 1, 634, 41, 52, 1, NOW() - INTERVAL 26 HOUR
),
(
  'Gin + GORM 从零搭一个能跑的用户模块',
  '包含注册、登录、JWT 中间件、统一错误返回。讲清楚 bcrypt 为什么不能明文存,以及 ShouldBindJSON 漏写 return 会出什么事。',
  '## 分层\n\nrouter 只做路由注册,handler 做参数校验和响应组装,dao 只碰数据库。中间件单独一层,负责解析 token 并把 user_id 塞进 context。\n\n## 三个容易漏的点\n\n1. ShouldBindJSON 失败后必须 return,否则会带着零值继续查库。\n2. 一个 handler 里只能写一次 c.JSON,第二次调用会报 superfluous response.WriteHeader,而且 body 会被追加写成非法 JSON。\n3. handler 里不要 panic,gin.Recovery 兜住之后客户端只拿到 500 空响应,排查时没有任何线索。',
  2, 2, 876, 152, 41, 1, NOW() - INTERVAL 5 HOUR
),
(
  'Go 里的 panic 和 error 边界在哪里',
  'handler 里出错到底该返回 500 还是 panic?panic 之后客户端看到的是什么?顺便说说日志该怎么打。',
  '我的理解:error 是可预期的失败(参数错、找不到记录、密码不对),panic 是程序写错了(数组越界、空指针)。\n\n前者应该显式返回并带上可读信息,后者应该让进程直接崩掉或者被 Recovery 兜住并告警。\n\n把数据库查不到记录 panic 掉是最常见的误用,那明明是 404。',
  2, 2, 921, 118, 28, 1, NOW() - INTERVAL 2 DAY
),
(
  '为什么你的 JOIN 查询越写越慢',
  '三个真实案例:索引建在区分度最低的列上、隐式类型转换让索引失效、在 WHERE 里对字段做函数运算。',
  '案例一:给 status 建了索引,全表只有两种状态,优化器直接不走索引。\n\n案例二:user_id 是 int,查询时传了字符串,MySQL 做了隐式转换,索引失效。\n\n案例三:WHERE DATE(created_at) = 今天,字段被函数包住,索引失效,改成范围查询。\n\n三个都能用 EXPLAIN 一眼看出来。',
  2, 3, 2103, 310, 67, 1, NOW() - INTERVAL 30 HOUR
),
(
  '二分查找写了五年还在写错边界',
  '左闭右开还是左闭右闭?mid 会不会溢出?循环条件用小于还是小于等于?把这些变量固定下来。',
  '我的固定写法:左闭右开 [left, right),循环条件 left < right,mid = left + (right - left) / 2。\n\n这样每轮区间严格缩小,不会死循环;返回 left 就是第一个满足条件的位置。\n\n记住一套就够,不要在两种写法之间来回横跳。',
  1, 4, 1542, 207, 35, 1, NOW() - INTERVAL 32 HOUR
),
(
  '写了十年代码,最后还是回归了最朴素的调试方式',
  'console.log、二分注释、最小复现。工具越花哨,越容易在工具本身出问题的时候束手无策。',
  '断点调试、性能分析、链路追踪都很好用,但真正救命的往往是二分注释:把一半代码注释掉,看问题还在不在。\n\n以及最小复现:能在一个空文件里重现的问题,基本五分钟就能定位。',
  1, 5, 3480, 486, 129, 1, NOW() - INTERVAL 3 DAY
),
(
  '草稿:文章列表接口的分页设计',
  '还没写完,先记一下 cursor 分页和 offset 分页的取舍。',
  'offset 分页在数据量大时越翻越慢,而且新数据插入会导致重复展示。\n\ncursor 分页用 created_at + id 作为游标,稳定但只能顺序翻页,不能跳到第 N 页。',
  2, 2, 0, 0, 0, 0, NOW() - INTERVAL 10 MINUTE
);

SET FOREIGN_KEY_CHECKS = 1;
