package entity

// 通用的查询和返回结构

// 分页查询.
type PageQuery struct {
	PageNo   int64 `example:"1"  json:"pageNo"`
	PageSize int64 `example:"10" json:"pageSize"`
}

// 分页结果返回.
type PageResult struct {
	PageNo     int64 `example:"1"   json:"pageNo"`
	PageSize   int64 `example:"10"  json:"pageSize"`
	TotalCount int64 `example:"100" json:"totalCount"`
}

// 排序查询.
type OrderQuery struct {
	Order   string `example:"asc" json:"order"`
	OrderBy string `example:"id"  json:"orderBy"`
}
