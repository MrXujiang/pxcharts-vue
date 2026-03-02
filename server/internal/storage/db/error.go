package db

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

const (
	ErrCodeUniqueViolation           = "23505" // 唯一约束冲突
	ErrCodeForeignKeyViolation       = "23503" // 外键约束冲突
	ErrCodeCheckViolation            = "23514" // check约束冲突
	ErrCodeNotNullViolation          = "23502" // 非空约束冲突
	ErrCodeStringDataRightTruncation = "22001" // 字符串长度超限
)

func GetPgErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

func IsUniqueViolation(err error) bool {
	return GetPgErrorCode(err) == ErrCodeUniqueViolation
}
