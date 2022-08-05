package dao

// 人工编写的查询语句，用于实现 sqlc 无法实现的功能，如 ORDER BY 的自定义排序
import (
	"context"
	"fmt"
)

type QueryUserParams struct {
	Offset   int64
	Limit    int64
	OrderBy  string
	Order    string
	Username string // 空字符串代表全部匹配
	Status   int32  // 0 则搜索全部状态
	Email    string // 空字符串代表全部匹配
}

func (q *Queries) QueryUser(ctx context.Context, arg QueryUserParams) ([]User, int64, error) {
	countSql := `SELECT COUNT(*) FROM user`
	querySql := `SELECT * FROM user`
	if arg.Username != "" || arg.Status != 0 || arg.Email != "" {
		countSql += ` WHERE `
		querySql += ` WHERE `
	}
	and := ""
	if arg.Username != "" {
		countSql += fmt.Sprintf(` %s username = '%s'`, and, arg.Username)
		querySql += fmt.Sprintf(` %s username = '%s'`, and, arg.Username)
		and = "and"
	}
	if arg.Status != 0 {
		countSql += fmt.Sprintf(` %s status = %d`, and, arg.Status)
		querySql += fmt.Sprintf(` %s status = %d`, and, arg.Status)
		and = "and"
	}
	if arg.Email != "" {
		countSql += fmt.Sprintf(` %s email = '%s'`, and, arg.Email)
		querySql += fmt.Sprintf(` %s email = '%s'`, and, arg.Email)
		and = "and"
	}
	if arg.OrderBy == "" {
		arg.OrderBy = "id"
	}
	if arg.Order == "" {
		arg.Order = "asc"
	}
	querySql += ` ORDER BY ` + arg.OrderBy + ` ` + arg.Order
	querySql += ` LIMIT ?, ?`

	// 计算 Count
	row := q.db.QueryRowContext(ctx, countSql)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return nil, 0, err
	}
	// 进行查询
	rows, err := q.db.QueryContext(ctx, querySql, arg.Offset, arg.Limit)
	if err != nil {
		return nil, count, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Status,
			&i.Email,
			&i.Password,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, count, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, count, err
	}
	if err := rows.Err(); err != nil {
		return nil, count, err
	}
	return items, count, nil
}
