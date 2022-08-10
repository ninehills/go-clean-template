package entity

// 通用的查询和返回结构

// 分页查询.
type PageQuery struct {
	PageNo   int64 `json:"pageNo" example:"1"`
	PageSize int64 `json:"pageSize" example:"10"`
}

// 分页结果返回.
type PageResult struct {
	PageNo     int64 `json:"pageNo" example:"1"`
	PageSize   int64 `json:"pageSize" example:"10"`
	TotalCount int64 `json:"totalCount" example:"100"`
}

// 排序查询.
type OrderQuery struct {
	Order   string `json:"order" example:"asc"`
	OrderBy string `json:"orderBy" example:"id"`
}
