package dao

// 人工编写的查询语句，用于实现 sqlc 无法实现的功能，如 ORDER BY 的自定义排序.
import (
	"context"
	"fmt"
)

const (
	and = ` AND `
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

func (q *Queries) QueryUser(ctx context.Context, arg QueryUserParams) (items []User, count int64, err error) {
	countSQL := buildUserQuerySQL(arg, true)
	querySQL := buildUserQuerySQL(arg, false)

	// 计算 Count
	row := q.db.QueryRowContext(ctx, countSQL)

	err = row.Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	// 进行查询
	rows, err := q.db.QueryContext(ctx, querySQL, arg.Offset, arg.Limit)
	if err != nil {
		return nil, count, err
	}
	defer rows.Close()

	items = []User{}

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

func buildUserQuerySQL(arg QueryUserParams, isCount bool) (sql string) {
	if isCount {
		sql = `SELECT COUNT(*) FROM user`
	} else {
		sql = `SELECT * FROM user`
	}

	if arg.Username != "" || arg.Status != 0 || arg.Email != "" {
		sql += ` WHERE `
	}

	connect := ""
	if arg.Username != "" {
		sql += fmt.Sprintf(` %s username = '%s'`, connect, arg.Username)
		connect = and
	}

	if arg.Status != 0 {
		sql += fmt.Sprintf(` %s status = %d`, connect, arg.Status)
		connect = and
	}

	if arg.Email != "" {
		sql += fmt.Sprintf(` %s email = '%s'`, connect, arg.Email)
	}

	if arg.OrderBy == "" {
		arg.OrderBy = "id"
	}

	if arg.Order == "" {
		arg.Order = "asc"
	}

	if !isCount {
		sql += ` ORDER BY ` + arg.OrderBy + ` ` + arg.Order
		sql += ` LIMIT ?, ?`
	}

	return sql
}
