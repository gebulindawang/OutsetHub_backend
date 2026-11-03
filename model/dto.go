package model

// ============================================================
// 响应外壳
//
// 实际输出由 response 包的 Ok/OkMsg/Fail 统一构造,
// 这里的结构体只用于 swag 生成文档,字段必须与 response.Body 保持一致。
// 改动时两边一起改,否则文档与实现会漂移。
// ============================================================

// MsgResp 成功,无数据
type MsgResp struct {
	Code int    `json:"code" example:"0"`
	Msg  string `json:"msg" example:"操作成功"`
}

// ErrResp 失败,业务错误与中间件错误共用此结构
type ErrResp struct {
	Code int    `json:"code" example:"400"`
	Msg  string `json:"msg" example:"参数错误"`
}

// ============================================================
// 用户
// ============================================================

// RegisterReq 注册请求
type RegisterReq struct {
	Username string `json:"username" binding:"required" example:"lisi"`
	Password string `json:"password" binding:"required" example:"123456"`
	Nickname string `json:"nickname" example:"李四"`
}

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required" example:"lisi"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// LoginData 登录成功后返回的数据体
type LoginData struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxx"`
}

// LoginResp 登录响应
type LoginResp struct {
	Code int       `json:"code" example:"0"`
	Msg  string    `json:"msg" example:"登录成功"`
	Data LoginData `json:"data"`
}

// UserResp 当前用户信息
type UserResp struct {
	Code int    `json:"code" example:"0"`
	Msg  string `json:"msg" example:"ok"`
	Data User   `json:"data"`
}

// ============================================================
// 文章
// ============================================================

// CreateArticleReq 发布文章请求
type CreateArticleReq struct {
	Title      string `json:"title" binding:"required" example:"Vue3 的 ref 和 reactive 该怎么选"`
	Summary    string `json:"summary" example:"一句话摘要,列表页展示"`
	Content    string `json:"content" binding:"required" example:"## 正文\n\n支持 Markdown"`
	CategoryID int    `json:"category_id" binding:"required" example:"1"`
	Status     *int   `json:"status" example:"1"` // 不传默认 1 已发布,0 为草稿
}

// UpdateArticleReq 编辑文章请求,传哪个字段改哪个
type UpdateArticleReq struct {
	Title      string `json:"title" example:"新标题"`
	Summary    string `json:"summary" example:"新摘要"`
	Content    string `json:"content" example:"新正文"`
	CategoryID *int   `json:"category_id" example:"2"`
	Status     *int   `json:"status" example:"1"`
}

// ArticleListResp 文章列表(仅已发布,按创建时间倒序,固定 10 条)
type ArticleListResp struct {
	Code int       `json:"code" example:"0"`
	Msg  string    `json:"msg" example:"ok"`
	Data []Article `json:"data"`
}

// ArticleResp 文章详情
type ArticleResp struct {
	Code int     `json:"code" example:"0"`
	Msg  string  `json:"msg" example:"ok"`
	Data Article `json:"data"`
}

// ============================================================
// 板块
// ============================================================

// CreateCategoryReq 新建板块请求
type CreateCategoryReq struct {
	Name string `json:"name" binding:"required" example:"前端"`
	Sort int    `json:"sort" example:"10"`
}

// UpdateCategoryReq 修改板块请求,部分更新
type UpdateCategoryReq struct {
	Name string `json:"name" example:"前端开发"`
	Sort *int   `json:"sort" example:"20"`
}

// CategoryListResp 板块列表(按 sort 升序)
type CategoryListResp struct {
	Code int        `json:"code" example:"0"`
	Msg  string     `json:"msg" example:"ok"`
	Data []Category `json:"data"`
}

// ============================================================
// 评论
// ============================================================

// CreateCommentReq 发表评论请求
type CreateCommentReq struct {
	ArticleID     uint64 `json:"article_id" binding:"required" example:"1"`
	ParentID      uint64 `json:"parent_id" example:"0"` // 0 表示一级评论
	ReplyToUserID int    `json:"reply_to_user_id" example:"0"`
	Content       string `json:"content" binding:"required" example:"这段代码在 Safari 下会报错"`
}

// UpdateCommentReq 修改评论请求
type UpdateCommentReq struct {
	Content string `json:"content" example:"修改后的内容"`
}

// CommentListResp 指定文章下的评论列表(按创建时间升序)
type CommentListResp struct {
	Code int       `json:"code" example:"0"`
	Msg  string    `json:"msg" example:"ok"`
	Data []Comment `json:"data"`
}

// CommentResp 评论详情
type CommentResp struct {
	Code int     `json:"code" example:"0"`
	Msg  string  `json:"msg" example:"ok"`
	Data Comment `json:"data"`
}
